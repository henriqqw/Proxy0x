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
	Raw         string      `json:"raw"`
	Host        string      `json:"host"`
	Port        string      `json:"port"`
	Username    string      `json:"username"`
	Password    string      `json:"password"`
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
	stopCh     chan struct{}
	logs       []LogEntry
	logMu      sync.RWMutex
	localProxy *LocalProxy
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
	Language string `json:"language"`
	Theme    string `json:"theme"`
}

func (a *App) loadConfig() AppConfig {
	data, err := os.ReadFile(a.configFile())
	if err != nil {
		return AppConfig{Language: "en", Theme: "dark"}
	}
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{Language: "en", Theme: "dark"}
	}
	if cfg.Language == "" {
		cfg.Language = "en"
	}
	if cfg.Theme == "" {
		cfg.Theme = "dark"
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
	a.loadProxies()
	go a.setupTray()
	go a.healthCheckLoop()
}

func (a *App) shutdown(ctx context.Context) {
	if a.localProxy != nil {
		a.localProxy.stop()
		a.localProxy = nil
		clearSystemProxy()
	}
	close(a.stopCh)
}

func (a *App) MinimizeWindow() {
	runtime.WindowMinimise(a.ctx)
}

func (a *App) CloseWindow() {
	runtime.WindowHide(a.ctx)
}

func (a *App) QuitApp() {
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

	parts := strings.SplitN(raw, ":", 4)
	if len(parts) < 4 {
		if len(parts) == 2 {
			return &ProxyEntry{
				ID:        generateID(),
				Raw:       raw,
				Host:      parts[0],
				Port:      parts[1],
				Status:    StatusOffline,
				SessionID: generateSessionID(),
			}, nil
		}
		return nil, fmt.Errorf("invalid proxy format, expected host:port:user:pass")
	}

	return &ProxyEntry{
		ID:        generateID(),
		Raw:       raw,
		Host:      parts[0],
		Port:      parts[1],
		Username:  parts[2],
		Password:  parts[3],
		Status:    StatusOffline,
		SessionID: generateSessionID(),
	}, nil
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

func (a *App) buildProxyURL(proxy *ProxyEntry) string {
	username := proxy.Username
	if username != "" && proxy.SessionID != "" {
		username = username + "_session-" + proxy.SessionID
	}
	if username != "" {
		return fmt.Sprintf("http://%s:%s@%s:%s",
			url.QueryEscape(username),
			url.QueryEscape(proxy.Password),
			proxy.Host, proxy.Port)
	}
	return fmt.Sprintf("http://%s:%s", proxy.Host, proxy.Port)
}

func (a *App) getOrCreateClient(proxy *ProxyEntry) *http.Client {
	clientKey := proxy.ID + ":" + proxy.SessionID
	if client, ok := a.clients[clientKey]; ok {
		return client
	}
	proxyURL, _ := url.Parse(a.buildProxyURL(proxy))
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
			MaxIdleConns:        2,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 20 * time.Second,
	}
	a.clients[clientKey] = client
	return client
}

func (a *App) ConnectProxy(id string) (*ProxyEntry, error) {
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

	if a.localProxy != nil {
		a.localProxy.stop()
		clearSystemProxy()
	}
	lp, lpErr := newLocalProxy(proxy)
	if lpErr == nil {
		a.localProxy = lp
		setSystemProxy(lp.port)
		a.addLog("system-proxy", id, proxy.Host+":"+proxy.Port, fmt.Sprintf("local proxy on port %d", lp.port))
	}

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
	proxy.Latency = 0
	proxy.ConnectedAt = ""
	delete(a.clients, id+":"+proxy.SessionID)
	a.mu.Unlock()

	if a.localProxy != nil && a.localProxy.upstream.ID == id {
		a.localProxy.stop()
		a.localProxy = nil
		clearSystemProxy()
	}

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
	oldSession := proxy.SessionID
	delete(a.clients, id+":"+oldSession)
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

	if a.localProxy != nil && a.localProxy.upstream.ID == id {
		a.localProxy.stop()
		lp, lpErr := newLocalProxy(proxy)
		if lpErr == nil {
			a.localProxy = lp
			setSystemProxy(lp.port)
		}
	}

	runtime.EventsEmit(a.ctx, "proxy-updated", proxy)
	return proxy, nil
}

func (a *App) RestartProxy(id string) (*ProxyEntry, error) {
	a.DisconnectProxy(id)
	time.Sleep(500 * time.Millisecond)
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
	resp, err := client.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=66846719", ip))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var geo GeoInfo
	if err := json.Unmarshal(body, &geo); err != nil {
		return nil
	}
	return &geo
}

func (a *App) healthCheckLoop() {
	ticker := time.NewTicker(90 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.stopCh:
			return
		case <-ticker.C:
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

func (a *App) GetMyIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}
