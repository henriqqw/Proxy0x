package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	netproxy "golang.org/x/net/proxy"
)

type ProxyStatus string

const (
	StatusOffline    ProxyStatus = "offline"
	StatusConnecting ProxyStatus = "connecting"
	StatusConnected  ProxyStatus = "connected"
	StatusError      ProxyStatus = "error"
)

type GeoInfo struct {
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	AS          string  `json:"as"`
}

type ProxyEntry struct {
	ID          string      `json:"id"`
	Host        string      `json:"host"`
	Port        string      `json:"port"`
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	Protocol    string      `json:"protocol"`
	Group       string      `json:"group"`
	Status      ProxyStatus `json:"status"`
	CurrentIP   string      `json:"currentIp"`
	Location    *GeoInfo    `json:"location"`
	Latency     int64       `json:"latency"`
	LastCheck   string      `json:"lastCheck"`
	Active      bool        `json:"active"`
	SessionID   string      `json:"sessionId"`
	Order       int         `json:"order"`
	ConnectedAt string      `json:"connectedAt"`
	TotalUptime int64       `json:"totalUptime"`
}

type LogEntry struct {
	Time    string `json:"time"`
	Action  string `json:"action"`
	ProxyID string `json:"proxyId"`
	Host    string `json:"host"`
	Detail  string `json:"detail"`
}

type App struct {
	ctx        context.Context
	proxies    map[string]*ProxyEntry
	clients    map[string]*http.Client
	mu         sync.RWMutex
	lpMu       sync.Mutex
	localProxy *LocalProxy
	stopCh     chan struct{}
	logs       []LogEntry
	logMu      sync.RWMutex
}

func NewApp() *App {
	return &App{
		proxies: make(map[string]*ProxyEntry),
		clients: make(map[string]*http.Client),
		stopCh:  make(chan struct{}),
		logs:    make([]LogEntry, 0, 200),
	}
}

func (a *App) addLog(action, proxyID, host, detail string) {
	entry := LogEntry{
		Time:    time.Now().Format(time.RFC3339),
		Action:  action,
		ProxyID: proxyID,
		Host:    host,
		Detail:  detail,
	}
	a.logMu.Lock()
	a.logs = append(a.logs, entry)
	if len(a.logs) > 200 {
		a.logs = a.logs[len(a.logs)-200:]
	}
	a.logMu.Unlock()
	runtime.EventsEmit(a.ctx, "log-entry", entry)
}

func (a *App) GetConnectionLog() []LogEntry {
	a.logMu.RLock()
	defer a.logMu.RUnlock()
	result := make([]LogEntry, len(a.logs))
	copy(result, a.logs)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func (a *App) dataDir() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "Proxy0x")
}

func (a *App) dataFile() string {
	return filepath.Join(a.dataDir(), "proxies.json")
}

func (a *App) configFile() string {
	return filepath.Join(a.dataDir(), "config.json")
}

type AppConfig struct {
	Language            string `json:"language"`
	Theme               string `json:"theme"`
	AutoRotateInterval  int    `json:"autoRotateInterval"`
	HealthCheckInterval int    `json:"healthCheckInterval"`
}

func (a *App) loadConfig() AppConfig {
	data, err := os.ReadFile(a.configFile())
	if err != nil {
		return AppConfig{Language: "en", Theme: "dark", HealthCheckInterval: 180}
	}
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{Language: "en", Theme: "dark", HealthCheckInterval: 180}
	}
	if cfg.Language == "" {
		cfg.Language = "en"
	}
	if cfg.Theme == "" {
		cfg.Theme = "dark"
	}
	if cfg.HealthCheckInterval <= 0 {
		cfg.HealthCheckInterval = 180
	}
	return cfg
}

func (a *App) saveConfig(cfg AppConfig) {
	os.MkdirAll(a.dataDir(), 0o700)
	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(a.configFile(), data, 0o600)
}

func (a *App) GetLanguage() string {
	return a.loadConfig().Language
}

func (a *App) SetLanguage(lang string) {
	cfg := a.loadConfig()
	cfg.Language = lang
	a.saveConfig(cfg)
	runtime.EventsEmit(a.ctx, "language-changed", lang)
}

func (a *App) GetTheme() string {
	return a.loadConfig().Theme
}

func (a *App) SetTheme(theme string) {
	cfg := a.loadConfig()
	cfg.Theme = theme
	a.saveConfig(cfg)
	runtime.EventsEmit(a.ctx, "theme-changed", theme)
}

func (a *App) GetAutoRotateInterval() int {
	return a.loadConfig().AutoRotateInterval
}

func (a *App) SetAutoRotateInterval(seconds int) {
	cfg := a.loadConfig()
	cfg.AutoRotateInterval = seconds
	a.saveConfig(cfg)
	runtime.EventsEmit(a.ctx, "config-changed", cfg)
}

func (a *App) GetHealthCheckInterval() int {
	return a.loadConfig().HealthCheckInterval
}

func (a *App) SetHealthCheckInterval(seconds int) {
	if seconds < 60 {
		seconds = 60
	}
	cfg := a.loadConfig()
	cfg.HealthCheckInterval = seconds
	a.saveConfig(cfg)
	runtime.EventsEmit(a.ctx, "config-changed", cfg)
}

func (a *App) loadProxies() {
	data, err := os.ReadFile(a.dataFile())
	if err != nil {
		return
	}
	var entries []*ProxyEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return
	}
	a.mu.Lock()
	for _, p := range entries {
		if dec, err := a.decrypt(p.Username); err == nil {
			p.Username = dec
		}
		if dec, err := a.decrypt(p.Password); err == nil {
			p.Password = dec
		}
		p.Status = StatusOffline
		p.Active = false
		p.CurrentIP = ""
		p.Latency = 0
		a.proxies[p.ID] = p
	}
	a.mu.Unlock()
}

func (a *App) saveProxies() {
	a.mu.RLock()
	entries := make([]*ProxyEntry, 0, len(a.proxies))
	for _, p := range a.proxies {
		clone := *p
		if enc, err := a.encrypt(clone.Username); err == nil {
			clone.Username = enc
		}
		if enc, err := a.encrypt(clone.Password); err == nil {
			clone.Password = enc
		}
		clone.Status = StatusOffline
		clone.Active = false
		clone.CurrentIP = ""
		clone.Latency = 0
		entries = append(entries, &clone)
	}
	a.mu.RUnlock()

	os.MkdirAll(a.dataDir(), 0o700)
	data, _ := json.MarshalIndent(entries, "", "  ")
	os.WriteFile(a.dataFile(), data, 0o600)
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	clearSystemProxy()
	a.loadProxies()
	go a.setupTray()
	go a.healthCheckLoop()
	go a.autoRotateLoop()
}

func (a *App) shutdown(ctx context.Context) {
	a.lpMu.Lock()
	if a.localProxy != nil {
		a.localProxy.stop()
		a.localProxy = nil
		clearSystemProxy()
	}
	a.lpMu.Unlock()
	close(a.stopCh)
}

func (a *App) MinimizeWindow() {
	runtime.WindowMinimise(a.ctx)
}

func (a *App) CloseWindow() {
	runtime.WindowHide(a.ctx)
}

func (a *App) QuitApp() {
	a.lpMu.Lock()
	if a.localProxy != nil {
		a.localProxy.stop()
		a.localProxy = nil
	}
	a.lpMu.Unlock()
	clearSystemProxy()
	runtime.Quit(a.ctx)
}

func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
}

func (a *App) ToggleMaximize() {
	runtime.WindowToggleMaximise(a.ctx)
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateSessionID() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func parseProxy(raw string) (*ProxyEntry, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty proxy string")
	}

	protocol := "http"
	if strings.HasPrefix(raw, "socks5://") {
		protocol = "socks5"
		raw = strings.TrimPrefix(raw, "socks5://")
	} else if strings.HasPrefix(raw, "http://") {
		raw = strings.TrimPrefix(raw, "http://")
	}

	parts := strings.SplitN(raw, ":", 4)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid proxy format, expected host:port or host:port:user:pass")
	}

	entry := &ProxyEntry{
		ID:        generateID(),
		Host:      parts[0],
		Port:      parts[1],
		Protocol:  protocol,
		Status:    StatusOffline,
		SessionID: generateSessionID(),
	}

	if len(parts) >= 4 {
		entry.Username = parts[2]
		entry.Password = parts[3]
	}

	return entry, nil
}

func (a *App) AddProxy(raw string) (*ProxyEntry, error) {
	proxy, err := parseProxy(raw)
	if err != nil {
		return nil, err
	}

	a.mu.Lock()
	a.proxies[proxy.ID] = proxy
	a.mu.Unlock()

	a.addLog("added", proxy.ID, proxy.Host+":"+proxy.Port, "")
	a.saveProxies()
	runtime.EventsEmit(a.ctx, "proxy-updated", proxy)
	return proxy, nil
}

func (a *App) AddProxiesBatch(raw string) ([]*ProxyEntry, error) {
	lines := strings.Split(raw, "\n")
	var added []*ProxyEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		proxy, err := parseProxy(line)
		if err != nil {
			continue
		}
		a.mu.Lock()
		a.proxies[proxy.ID] = proxy
		a.mu.Unlock()
		added = append(added, proxy)
	}

	a.addLog("batch-added", "", "", fmt.Sprintf("%d proxies", len(added)))
	a.saveProxies()
	runtime.EventsEmit(a.ctx, "proxies-changed")
	return added, nil
}

func (a *App) RemoveProxy(id string) {
	a.mu.Lock()
	p := a.proxies[id]
	host := ""
	if p != nil {
		host = p.Host + ":" + p.Port
	}
	delete(a.proxies, id)
	delete(a.clients, id)
	a.mu.Unlock()
	a.addLog("removed", id, host, "")
	a.saveProxies()
	runtime.EventsEmit(a.ctx, "proxies-changed")
}

func (a *App) GetProxies() []*ProxyEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make([]*ProxyEntry, 0, len(a.proxies))
	for _, p := range a.proxies {
		result = append(result, p)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Order < result[j].Order
	})
	return result
}

func (a *App) ReorderProxies(ids []string) {
	a.mu.Lock()
	for i, id := range ids {
		if p, ok := a.proxies[id]; ok {
			p.Order = i
		}
	}
	a.mu.Unlock()
	a.saveProxies()
}

func (a *App) ExportProxies() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var lines []string
	for _, p := range a.proxies {
		if p.Username != "" {
			lines = append(lines, fmt.Sprintf("%s:%s:%s:%s", p.Host, p.Port, p.Username, p.Password))
		} else {
			lines = append(lines, fmt.Sprintf("%s:%s", p.Host, p.Port))
		}
	}
	return strings.Join(lines, "\n")
}

func (a *App) buildProxyURL(proxy *ProxyEntry) string {
	scheme := "http"
	if proxy.Protocol == "socks5" {
		scheme = "socks5"
	}
	username := proxy.Username
	if username != "" && proxy.SessionID != "" {
		username = username + "_session-" + proxy.SessionID
	}
	if username != "" {
		return fmt.Sprintf("%s://%s:%s@%s:%s",
			scheme,
			url.QueryEscape(username),
			url.QueryEscape(proxy.Password),
			proxy.Host, proxy.Port)
	}
	return fmt.Sprintf("%s://%s:%s", scheme, proxy.Host, proxy.Port)
}

func (a *App) getOrCreateClient(proxy *ProxyEntry) *http.Client {
	clientKey := proxy.ID + ":" + proxy.SessionID

	if old, ok := a.clients[clientKey]; ok {
		return old
	}

	var transport *http.Transport

	if proxy.Protocol == "socks5" {
		username := proxy.Username
		if username != "" && proxy.SessionID != "" {
			username = username + "_session-" + proxy.SessionID
		}
		var auth *netproxy.Auth
		if username != "" {
			auth = &netproxy.Auth{User: username, Password: proxy.Password}
		}
		dialer, err := netproxy.SOCKS5("tcp", net.JoinHostPort(proxy.Host, proxy.Port), auth, netproxy.Direct)
		if err != nil {
			dialer = netproxy.Direct
		}
		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
			TLSHandshakeTimeout: 10 * time.Second,
			MaxIdleConns:        2,
			IdleConnTimeout:     90 * time.Second,
		}
	} else {
		proxyURL, _ := url.Parse(a.buildProxyURL(proxy))
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
			MaxIdleConns:        2,
			IdleConnTimeout:     90 * time.Second,
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   20 * time.Second,
	}
	a.clients[clientKey] = client
	return client
}

func (a *App) closeClient(key string) {
	if c, ok := a.clients[key]; ok {
		if t, ok := c.Transport.(*http.Transport); ok {
			t.CloseIdleConnections()
		}
		delete(a.clients, key)
	}
}

func (a *App) setLocalProxy(proxy *ProxyEntry, id string) {
	a.lpMu.Lock()
	defer a.lpMu.Unlock()

	if a.localProxy != nil {
		a.localProxy.stop()
		clearSystemProxy()
	}

	lp, err := newLocalProxy(proxy)
	if err == nil {
		a.localProxy = lp
		setSystemProxy(lp.port)
		a.addLog("system-proxy", id, proxy.Host+":"+proxy.Port, fmt.Sprintf("local proxy on port %d", lp.port))
	}
}

func (a *App) clearLocalProxyIf(id string) {
	a.lpMu.Lock()
	defer a.lpMu.Unlock()

	if a.localProxy != nil && a.localProxy.upstream.ID == id {
		a.localProxy.stop()
		a.localProxy = nil
		clearSystemProxy()
	}
}

func (a *App) ConnectProxy(id string) (*ProxyEntry, error) {
	// Disconnect any currently active proxy first
	a.mu.RLock()
	var activeID string
	for pid, p := range a.proxies {
		if p.Active && pid != id {
			activeID = pid
			break
		}
	}
	a.mu.RUnlock()
	if activeID != "" {
		a.DisconnectProxy(activeID)
	}

	a.mu.Lock()
	proxy, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	proxy.Status = StatusConnecting
	client := a.getOrCreateClient(proxy)
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "proxy-updated", proxy)

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	latency := time.Since(start).Milliseconds()

	if err != nil {
		a.mu.Lock()
		proxy.Status = StatusError
		proxy.Latency = 0
		proxy.CurrentIP = ""
		proxy.Location = nil
		proxy.LastCheck = time.Now().Format(time.RFC3339)
		a.mu.Unlock()
		a.addLog("connect-failed", id, proxy.Host+":"+proxy.Port, err.Error())
		runtime.EventsEmit(a.ctx, "proxy-updated", proxy)
		return proxy, fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ipResult struct {
		IP string `json:"ip"`
	}
	json.Unmarshal(body, &ipResult)

	geo := a.fetchGeoInfo(client, ipResult.IP)

	a.mu.Lock()
	proxy.Status = StatusConnected
	proxy.CurrentIP = ipResult.IP
	proxy.Latency = latency
	proxy.Location = geo
	proxy.Active = true
	proxy.LastCheck = time.Now().Format(time.RFC3339)
	proxy.ConnectedAt = time.Now().Format(time.RFC3339)
	a.mu.Unlock()

	a.addLog("connected", id, proxy.Host+":"+proxy.Port, fmt.Sprintf("IP: %s, %dms", ipResult.IP, latency))
	a.setLocalProxy(proxy, id)

	runtime.EventsEmit(a.ctx, "proxy-updated", proxy)
	return proxy, nil
}

func (a *App) DisconnectProxy(id string) (*ProxyEntry, error) {
	a.mu.Lock()
	proxy, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	if proxy.ConnectedAt != "" {
		if t, err := time.Parse(time.RFC3339, proxy.ConnectedAt); err == nil {
			proxy.TotalUptime += int64(time.Since(t).Seconds())
		}
	}
	proxy.Status = StatusOffline
	proxy.Active = false
	proxy.CurrentIP = ""
	proxy.Location = nil
	proxy.Latency = 0
	proxy.ConnectedAt = ""
	a.closeClient(id + ":" + proxy.SessionID)
	a.mu.Unlock()

	a.clearLocalProxyIf(id)

	a.addLog("disconnected", id, proxy.Host+":"+proxy.Port, "")
	runtime.EventsEmit(a.ctx, "proxy-updated", proxy)
	return proxy, nil
}

func (a *App) RotateProxy(id string) (*ProxyEntry, error) {
	a.mu.Lock()
	proxy, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	oldKey := id + ":" + proxy.SessionID
	a.closeClient(oldKey)
	proxy.SessionID = generateSessionID()
	proxy.Status = StatusConnecting
	client := a.getOrCreateClient(proxy)
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "proxy-updated", proxy)

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	latency := time.Since(start).Milliseconds()

	if err != nil {
		a.mu.Lock()
		proxy.Status = StatusError
		a.mu.Unlock()
		a.addLog("rotate-failed", id, proxy.Host+":"+proxy.Port, err.Error())
		runtime.EventsEmit(a.ctx, "proxy-updated", proxy)
		return proxy, fmt.Errorf("rotation failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ipResult struct {
		IP string `json:"ip"`
	}
	json.Unmarshal(body, &ipResult)

	geo := a.fetchGeoInfo(client, ipResult.IP)

	a.mu.Lock()
	proxy.Status = StatusConnected
	proxy.CurrentIP = ipResult.IP
	proxy.Latency = latency
	proxy.Location = geo
	proxy.Active = true
	proxy.LastCheck = time.Now().Format(time.RFC3339)
	a.mu.Unlock()

	a.addLog("rotated", id, proxy.Host+":"+proxy.Port, fmt.Sprintf("IP: %s, %dms", ipResult.IP, latency))
	a.setLocalProxy(proxy, id)

	runtime.EventsEmit(a.ctx, "proxy-updated", proxy)
	return proxy, nil
}

func (a *App) RestartProxy(id string) (*ProxyEntry, error) {
	a.DisconnectProxy(id)
	return a.ConnectProxy(id)
}

func (a *App) CheckHealth(id string) (*ProxyEntry, error) {
	a.mu.RLock()
	proxy, ok := a.proxies[id]
	if !ok || !proxy.Active {
		a.mu.RUnlock()
		return nil, fmt.Errorf("proxy not found or inactive")
	}
	client := a.getOrCreateClient(proxy)
	prevStatus := proxy.Status
	a.mu.RUnlock()

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	latency := time.Since(start).Milliseconds()

	if err != nil {
		a.mu.Lock()
		proxy.Status = StatusError
		proxy.LastCheck = time.Now().Format(time.RFC3339)
		a.mu.Unlock()
		if prevStatus == StatusConnected {
			a.addLog("dropped", id, proxy.Host+":"+proxy.Port, err.Error())
			runtime.EventsEmit(a.ctx, "proxy-dropped", map[string]string{
				"id":   id,
				"host": proxy.Host + ":" + proxy.Port,
			})
		}
		runtime.EventsEmit(a.ctx, "proxy-updated", proxy)
		return proxy, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ipResult struct {
		IP string `json:"ip"`
	}
	json.Unmarshal(body, &ipResult)

	a.mu.Lock()
	proxy.Latency = latency
	proxy.LastCheck = time.Now().Format(time.RFC3339)
	if ipResult.IP != "" && ipResult.IP != proxy.CurrentIP {
		proxy.CurrentIP = ipResult.IP
		geo := a.fetchGeoInfo(client, ipResult.IP)
		proxy.Location = geo
	}
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "proxy-updated", proxy)
	return proxy, nil
}

func (a *App) fetchGeoInfo(client *http.Client, ip string) *GeoInfo {
	resp, err := client.Get(fmt.Sprintf("https://ipwho.is/%s", ip))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var raw struct {
		Country     string  `json:"country"`
		CountryCode string  `json:"country_code"`
		Region      string  `json:"region"`
		City        string  `json:"city"`
		Lat         float64 `json:"latitude"`
		Lon         float64 `json:"longitude"`
		Connection  struct {
			ISP    string `json:"isp"`
			Org    string `json:"org"`
			ASN    int    `json:"asn"`
			Domain string `json:"domain"`
		} `json:"connection"`
		Timezone struct {
			ID string `json:"id"`
		} `json:"timezone"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}

	asStr := ""
	if raw.Connection.ASN > 0 {
		asStr = fmt.Sprintf("AS%d", raw.Connection.ASN)
	}

	return &GeoInfo{
		Country:     raw.Country,
		CountryCode: raw.CountryCode,
		Region:      raw.Region,
		City:        raw.City,
		ISP:         raw.Connection.ISP,
		Org:         raw.Connection.Org,
		Lat:         raw.Lat,
		Lon:         raw.Lon,
		Timezone:    raw.Timezone.ID,
		AS:          asStr,
	}
}

func (a *App) healthCheckLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	lastCheck := time.Now()

	for {
		select {
		case <-a.stopCh:
			return
		case <-ticker.C:
			cfg := a.loadConfig()
			interval := time.Duration(cfg.HealthCheckInterval) * time.Second
			if time.Since(lastCheck) < interval {
				continue
			}
			lastCheck = time.Now()

			a.mu.RLock()
			var activeIDs []string
			for id, p := range a.proxies {
				if p.Active {
					activeIDs = append(activeIDs, id)
				}
			}
			a.mu.RUnlock()

			for i, id := range activeIDs {
				select {
				case <-a.stopCh:
					return
				default:
				}
				a.CheckHealth(id)
				if i < len(activeIDs)-1 {
					time.Sleep(1500 * time.Millisecond)
				}
			}
		}
	}
}

func (a *App) autoRotateLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	lastRotation := time.Now()

	for {
		select {
		case <-a.stopCh:
			return
		case <-ticker.C:
			cfg := a.loadConfig()
			if cfg.AutoRotateInterval <= 0 {
				lastRotation = time.Now()
				continue
			}
			if time.Since(lastRotation) < time.Duration(cfg.AutoRotateInterval)*time.Second {
				continue
			}
			lastRotation = time.Now()

			a.mu.RLock()
			var activeIDs []string
			for id, p := range a.proxies {
				if p.Active {
					activeIDs = append(activeIDs, id)
				}
			}
			a.mu.RUnlock()

			for _, id := range activeIDs {
				select {
				case <-a.stopCh:
					return
				default:
				}
				a.RotateProxy(id)
			}

			if len(activeIDs) > 0 {
				a.addLog("auto-rotate", "", "", fmt.Sprintf("rotated %d proxies", len(activeIDs)))
			}
		}
	}
}

func (a *App) BenchmarkProxy(id string) (map[string]interface{}, error) {
	a.mu.RLock()
	proxy, ok := a.proxies[id]
	if !ok {
		a.mu.RUnlock()
		return nil, fmt.Errorf("proxy not found")
	}
	client := a.getOrCreateClient(proxy)
	a.mu.RUnlock()

	results := map[string]interface{}{}

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	connectMs := time.Since(start).Milliseconds()
	if err != nil {
		return nil, fmt.Errorf("connect failed: %v", err)
	}
	resp.Body.Close()
	results["connectMs"] = connectMs

	start = time.Now()
	resp, err = client.Get("https://speed.cloudflare.com/__down?bytes=102400")
	if err != nil {
		results["downloadMs"] = int64(-1)
		results["downloadMbps"] = float64(0)
	} else {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		dlMs := time.Since(start).Milliseconds()
		results["downloadMs"] = dlMs
		if dlMs > 0 {
			results["downloadMbps"] = float64(len(body)) * 8 / float64(dlMs) / 1000
		}
	}

	a.addLog("benchmark", id, proxy.Host+":"+proxy.Port, fmt.Sprintf("connect: %dms", connectMs))
	return results, nil
}

func (a *App) GetGroups() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	seen := map[string]bool{}
	var groups []string
	for _, p := range a.proxies {
		if p.Group != "" && !seen[p.Group] {
			seen[p.Group] = true
			groups = append(groups, p.Group)
		}
	}
	sort.Strings(groups)
	return groups
}

func (a *App) SetProxyGroup(id, group string) (*ProxyEntry, error) {
	a.mu.Lock()
	proxy, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	proxy.Group = group
	a.mu.Unlock()
	a.saveProxies()
	runtime.EventsEmit(a.ctx, "proxy-updated", proxy)
	return proxy, nil
}

func (a *App) ConnectGroup(group string) error {
	a.mu.RLock()
	var ids []string
	for id, p := range a.proxies {
		if p.Group == group && !p.Active {
			ids = append(ids, id)
		}
	}
	a.mu.RUnlock()

	for _, id := range ids {
		a.ConnectProxy(id)
	}
	return nil
}

func (a *App) DisconnectGroup(group string) error {
	a.mu.RLock()
	var ids []string
	for id, p := range a.proxies {
		if p.Group == group && p.Active {
			ids = append(ids, id)
		}
	}
	a.mu.RUnlock()

	for _, id := range ids {
		a.DisconnectProxy(id)
	}
	return nil
}


func (a *App) GetMyIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}
