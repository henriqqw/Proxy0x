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
	"sort"
	"strings"
	"sync"
	"time"

	"proxy_ui/internal/config"
	"proxy_ui/internal/crypto"
	"proxy_ui/internal/killswitch"
	"proxy_ui/internal/proxy"

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
	Name        string      `json:"name"`
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

type BulkCheckResult struct {
	ID          string `json:"id"`
	Host        string `json:"host"`
	Alive       bool   `json:"alive"`
	Latency     int64  `json:"latency"`
	ExitIP      string `json:"exitIp"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Anonymity   string `json:"anonymity"`
	Error       string `json:"error,omitempty"`
}

type DNSLeakResult struct {
	ProxyIP   string `json:"proxyIp"`
	DNSIP     string `json:"dnsIp"`
	DNSServer string `json:"dnsServer"`
	Leaked    bool   `json:"leaked"`
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
	localProxy *proxy.LocalServer
	stopCh     chan struct{}
	logs       []LogEntry
	logMu      sync.RWMutex
	crypt      *crypto.Store
}

func NewApp() *App {
	return &App{
		proxies: make(map[string]*ProxyEntry),
		clients: make(map[string]*http.Client),
		stopCh:  make(chan struct{}),
		logs:    make([]LogEntry, 0, 200),
		crypt:   crypto.NewStore(config.DataDir()),
	}
}

// --- Lifecycle ---

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	proxy.ClearSystemProxy()
	a.loadProxies()
	go a.setupTray()
	go a.healthCheckLoop()
	go a.autoRotateLoop()
}

func (a *App) shutdown(ctx context.Context) {
	if config.Load().KillSwitch {
		killswitch.Disable()
		cfg := config.Load()
		cfg.KillSwitch = false
		config.Save(cfg)
	}
	a.lpMu.Lock()
	if a.localProxy != nil {
		a.localProxy.Stop()
		a.localProxy = nil
		proxy.ClearSystemProxy()
	}
	a.lpMu.Unlock()
	close(a.stopCh)
}

// --- Window controls ---

func (a *App) MinimizeWindow() { runtime.WindowMinimise(a.ctx) }
func (a *App) CloseWindow()    { runtime.WindowHide(a.ctx) }
func (a *App) ShowWindow()     { runtime.WindowShow(a.ctx) }
func (a *App) ToggleMaximize() { runtime.WindowToggleMaximise(a.ctx) }

func (a *App) QuitApp() {
	a.lpMu.Lock()
	if a.localProxy != nil {
		a.localProxy.Stop()
		a.localProxy = nil
	}
	a.lpMu.Unlock()
	proxy.ClearSystemProxy()
	runtime.Quit(a.ctx)
}

// --- Config ---

func (a *App) GetLanguage() string         { return config.Load().Language }
func (a *App) GetTheme() string            { return config.Load().Theme }
func (a *App) GetAutoRotateInterval() int  { return config.Load().AutoRotateInterval }
func (a *App) GetHealthCheckInterval() int { return config.Load().HealthCheckInterval }

func (a *App) SetLanguage(lang string) {
	cfg := config.Load()
	cfg.Language = lang
	config.Save(cfg)
	runtime.EventsEmit(a.ctx, "language-changed", lang)
}

func (a *App) SetTheme(theme string) {
	cfg := config.Load()
	cfg.Theme = theme
	config.Save(cfg)
	runtime.EventsEmit(a.ctx, "theme-changed", theme)
}

func (a *App) SetAutoRotateInterval(seconds int) {
	cfg := config.Load()
	cfg.AutoRotateInterval = seconds
	config.Save(cfg)
	runtime.EventsEmit(a.ctx, "config-changed", cfg)
}

func (a *App) SetHealthCheckInterval(seconds int) {
	if seconds < 60 {
		seconds = 60
	}
	cfg := config.Load()
	cfg.HealthCheckInterval = seconds
	config.Save(cfg)
	runtime.EventsEmit(a.ctx, "config-changed", cfg)
}

// --- Logging ---

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

// --- Persistence ---

func (a *App) loadProxies() {
	data, err := os.ReadFile(config.DataFile())
	if err != nil {
		return
	}
	var entries []*ProxyEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return
	}
	a.mu.Lock()
	for _, p := range entries {
		if dec, err := a.crypt.Decrypt(p.Username); err == nil {
			p.Username = dec
		}
		if dec, err := a.crypt.Decrypt(p.Password); err == nil {
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
		if enc, err := a.crypt.Encrypt(clone.Username); err == nil {
			clone.Username = enc
		}
		if enc, err := a.crypt.Encrypt(clone.Password); err == nil {
			clone.Password = enc
		}
		clone.Status = StatusOffline
		clone.Active = false
		clone.CurrentIP = ""
		clone.Latency = 0
		entries = append(entries, &clone)
	}
	a.mu.RUnlock()

	os.MkdirAll(config.DataDir(), 0o700)
	data, _ := json.MarshalIndent(entries, "", "  ")
	os.WriteFile(config.DataFile(), data, 0o600)
}

// --- Proxy CRUD ---

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
	p, err := parseProxy(raw)
	if err != nil {
		return nil, err
	}

	a.mu.Lock()
	a.proxies[p.ID] = p
	a.mu.Unlock()

	a.addLog("added", p.ID, p.Host+":"+p.Port, "")
	a.saveProxies()
	runtime.EventsEmit(a.ctx, "proxy-updated", p)
	return p, nil
}

func (a *App) AddProxiesBatch(raw string) ([]*ProxyEntry, error) {
	lines := strings.Split(raw, "\n")
	var added []*ProxyEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		p, err := parseProxy(line)
		if err != nil {
			continue
		}
		a.mu.Lock()
		a.proxies[p.ID] = p
		a.mu.Unlock()
		added = append(added, p)
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

// --- Groups ---

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

func (a *App) SetProxyName(id, name string) (*ProxyEntry, error) {
	a.mu.Lock()
	p, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	p.Name = name
	a.mu.Unlock()
	a.saveProxies()
	runtime.EventsEmit(a.ctx, "proxy-updated", p)
	return p, nil
}

func (a *App) SetProxyGroup(id, group string) (*ProxyEntry, error) {
	a.mu.Lock()
	p, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	p.Group = group
	a.mu.Unlock()
	a.saveProxies()
	runtime.EventsEmit(a.ctx, "proxy-updated", p)
	return p, nil
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

// --- Connection ---

func (a *App) buildProxyURL(p *ProxyEntry) string {
	scheme := "http"
	if p.Protocol == "socks5" {
		scheme = "socks5"
	}
	username := p.Username
	if username != "" && p.SessionID != "" {
		username = username + "_session-" + p.SessionID
	}
	if username != "" {
		return fmt.Sprintf("%s://%s:%s@%s:%s",
			scheme,
			url.QueryEscape(username),
			url.QueryEscape(p.Password),
			p.Host, p.Port)
	}
	return fmt.Sprintf("%s://%s:%s", scheme, p.Host, p.Port)
}

func (a *App) getOrCreateClient(p *ProxyEntry) *http.Client {
	clientKey := p.ID + ":" + p.SessionID

	if old, ok := a.clients[clientKey]; ok {
		return old
	}

	var transport *http.Transport

	if p.Protocol == "socks5" {
		username := p.Username
		if username != "" && p.SessionID != "" {
			username = username + "_session-" + p.SessionID
		}
		var auth *netproxy.Auth
		if username != "" {
			auth = &netproxy.Auth{User: username, Password: p.Password}
		}
		dialer, err := netproxy.SOCKS5("tcp", net.JoinHostPort(p.Host, p.Port), auth, netproxy.Direct)
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
		proxyURL, _ := url.Parse(a.buildProxyURL(p))
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

func (a *App) toUpstreamConfig(p *ProxyEntry) proxy.UpstreamConfig {
	return proxy.UpstreamConfig{
		ID:        p.ID,
		Host:      p.Host,
		Port:      p.Port,
		Username:  p.Username,
		Password:  p.Password,
		SessionID: p.SessionID,
		Protocol:  p.Protocol,
	}
}

func (a *App) setLocalProxy(p *ProxyEntry, id string) {
	a.lpMu.Lock()
	defer a.lpMu.Unlock()

	if a.localProxy != nil {
		a.localProxy.Stop()
		proxy.ClearSystemProxy()
	}

	ls, err := proxy.NewLocalServer(a.toUpstreamConfig(p))
	if err == nil {
		a.localProxy = ls
		proxy.SetSystemProxy(ls.Port())
		a.addLog("system-proxy", id, p.Host+":"+p.Port, fmt.Sprintf("local proxy on port %d", ls.Port()))
	}
}

func (a *App) clearLocalProxyIf(id string) {
	a.lpMu.Lock()
	defer a.lpMu.Unlock()

	if a.localProxy != nil && a.localProxy.UpstreamID() == id {
		a.localProxy.Stop()
		a.localProxy = nil
		proxy.ClearSystemProxy()
	}
}

func (a *App) ConnectProxy(id string) (*ProxyEntry, error) {
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
	p, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	p.Status = StatusConnecting
	client := a.getOrCreateClient(p)
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "proxy-updated", p)

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	latency := time.Since(start).Milliseconds()

	if err != nil {
		a.mu.Lock()
		p.Status = StatusError
		p.Latency = 0
		p.CurrentIP = ""
		p.Location = nil
		p.LastCheck = time.Now().Format(time.RFC3339)
		a.mu.Unlock()
		a.addLog("connect-failed", id, p.Host+":"+p.Port, err.Error())
		runtime.EventsEmit(a.ctx, "proxy-updated", p)
		return p, fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ipResult struct {
		IP string `json:"ip"`
	}
	json.Unmarshal(body, &ipResult)

	geo := a.fetchGeoInfo(client, ipResult.IP)

	a.mu.Lock()
	p.Status = StatusConnected
	p.CurrentIP = ipResult.IP
	p.Latency = latency
	p.Location = geo
	p.Active = true
	p.LastCheck = time.Now().Format(time.RFC3339)
	p.ConnectedAt = time.Now().Format(time.RFC3339)
	a.mu.Unlock()

	a.addLog("connected", id, p.Host+":"+p.Port, fmt.Sprintf("IP: %s, %dms", ipResult.IP, latency))
	a.setLocalProxy(p, id)

	runtime.EventsEmit(a.ctx, "proxy-updated", p)
	return p, nil
}

func (a *App) DisconnectProxy(id string) (*ProxyEntry, error) {
	a.mu.Lock()
	p, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	if p.ConnectedAt != "" {
		if t, err := time.Parse(time.RFC3339, p.ConnectedAt); err == nil {
			p.TotalUptime += int64(time.Since(t).Seconds())
		}
	}
	p.Status = StatusOffline
	p.Active = false
	p.CurrentIP = ""
	p.Location = nil
	p.Latency = 0
	p.ConnectedAt = ""
	a.closeClient(id + ":" + p.SessionID)
	a.mu.Unlock()

	a.clearLocalProxyIf(id)

	a.addLog("disconnected", id, p.Host+":"+p.Port, "")
	runtime.EventsEmit(a.ctx, "proxy-updated", p)
	return p, nil
}

func (a *App) RotateProxy(id string) (*ProxyEntry, error) {
	a.mu.Lock()
	p, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	oldKey := id + ":" + p.SessionID
	a.closeClient(oldKey)
	p.SessionID = generateSessionID()
	p.Status = StatusConnecting
	client := a.getOrCreateClient(p)
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "proxy-updated", p)

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	latency := time.Since(start).Milliseconds()

	if err != nil {
		a.mu.Lock()
		p.Status = StatusError
		a.mu.Unlock()
		a.addLog("rotate-failed", id, p.Host+":"+p.Port, err.Error())
		runtime.EventsEmit(a.ctx, "proxy-updated", p)
		return p, fmt.Errorf("rotation failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ipResult struct {
		IP string `json:"ip"`
	}
	json.Unmarshal(body, &ipResult)

	geo := a.fetchGeoInfo(client, ipResult.IP)

	a.mu.Lock()
	p.Status = StatusConnected
	p.CurrentIP = ipResult.IP
	p.Latency = latency
	p.Location = geo
	p.Active = true
	p.LastCheck = time.Now().Format(time.RFC3339)
	a.mu.Unlock()

	a.addLog("rotated", id, p.Host+":"+p.Port, fmt.Sprintf("IP: %s, %dms", ipResult.IP, latency))
	a.setLocalProxy(p, id)

	runtime.EventsEmit(a.ctx, "proxy-updated", p)
	return p, nil
}

func (a *App) RestartProxy(id string) (*ProxyEntry, error) {
	a.DisconnectProxy(id)
	return a.ConnectProxy(id)
}

// --- Health & rotation ---

func (a *App) CheckHealth(id string) (*ProxyEntry, error) {
	a.mu.RLock()
	p, ok := a.proxies[id]
	if !ok || !p.Active {
		a.mu.RUnlock()
		return nil, fmt.Errorf("proxy not found or inactive")
	}
	client := a.getOrCreateClient(p)
	prevStatus := p.Status
	a.mu.RUnlock()

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	latency := time.Since(start).Milliseconds()

	if err != nil {
		a.mu.Lock()
		p.Status = StatusError
		p.LastCheck = time.Now().Format(time.RFC3339)
		a.mu.Unlock()
		if prevStatus == StatusConnected {
			a.addLog("dropped", id, p.Host+":"+p.Port, err.Error())
			runtime.EventsEmit(a.ctx, "proxy-dropped", map[string]string{
				"id":   id,
				"host": p.Host + ":" + p.Port,
			})
			a.tryFailover(id)
		}
		runtime.EventsEmit(a.ctx, "proxy-updated", p)
		return p, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ipResult struct {
		IP string `json:"ip"`
	}
	json.Unmarshal(body, &ipResult)

	a.mu.Lock()
	p.Latency = latency
	p.LastCheck = time.Now().Format(time.RFC3339)
	if ipResult.IP != "" && ipResult.IP != p.CurrentIP {
		p.CurrentIP = ipResult.IP
		geo := a.fetchGeoInfo(client, ipResult.IP)
		p.Location = geo
	}
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "proxy-updated", p)
	return p, nil
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
			cfg := config.Load()
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
			cfg := config.Load()
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

// --- Benchmark ---

func (a *App) BenchmarkProxy(id string) (map[string]interface{}, error) {
	a.mu.RLock()
	p, ok := a.proxies[id]
	if !ok {
		a.mu.RUnlock()
		return nil, fmt.Errorf("proxy not found")
	}
	client := a.getOrCreateClient(p)
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

	a.addLog("benchmark", id, p.Host+":"+p.Port, fmt.Sprintf("connect: %dms", connectMs))
	return results, nil
}

// --- Public IP ---

func (a *App) GetMyIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

// --- Kill Switch ---

func (a *App) GetKillSwitch() bool {
	return config.Load().KillSwitch
}

func (a *App) SetKillSwitch(enabled bool) error {
	var err error
	if enabled {
		err = killswitch.Enable()
	} else {
		err = killswitch.Disable()
	}
	if err != nil {
		return err
	}
	cfg := config.Load()
	cfg.KillSwitch = enabled
	config.Save(cfg)
	return nil
}

// --- Auto Failover ---

func (a *App) GetAutoFailover() bool {
	return config.Load().AutoFailover
}

func (a *App) SetAutoFailover(enabled bool) {
	cfg := config.Load()
	cfg.AutoFailover = enabled
	config.Save(cfg)
}

func (a *App) tryFailover(failedID string) {
	if !config.Load().AutoFailover {
		return
	}
	a.mu.RLock()
	var best *ProxyEntry
	var bestLatency int64 = 999999
	for _, p := range a.proxies {
		if p.ID == failedID || p.Active {
			continue
		}
		if p.Latency > 0 && p.Latency < bestLatency {
			best = p
			bestLatency = p.Latency
		} else if best == nil {
			best = p
		}
	}
	a.mu.RUnlock()

	if best == nil {
		return
	}
	a.addLog("failover", best.ID, best.Host+":"+best.Port, fmt.Sprintf("auto-failover from %s", failedID))
	go a.ConnectProxy(best.ID)
}

// --- Auto-Selection by Latency ---

func (a *App) ConnectFastest() (*ProxyEntry, error) {
	a.mu.RLock()
	proxies := make([]*ProxyEntry, 0)
	for _, p := range a.proxies {
		if !p.Active {
			proxies = append(proxies, p)
		}
	}
	a.mu.RUnlock()

	if len(proxies) == 0 {
		return nil, fmt.Errorf("no proxies available")
	}

	type result struct {
		id      string
		latency int64
	}
	results := make(chan result, len(proxies))
	sem := make(chan struct{}, 10)

	for _, p := range proxies {
		go func(p *ProxyEntry) {
			sem <- struct{}{}
			defer func() { <-sem }()
			client := a.getOrCreateClient(p)
			start := time.Now()
			resp, err := client.Get("https://api.ipify.org?format=json")
			lat := time.Since(start).Milliseconds()
			if err != nil {
				results <- result{p.ID, 999999}
				return
			}
			resp.Body.Close()
			results <- result{p.ID, lat}
		}(p)
	}

	var bestID string
	var bestLat int64 = 999999
	for range proxies {
		r := <-results
		if r.latency < bestLat {
			bestLat = r.latency
			bestID = r.id
		}
	}

	if bestID == "" || bestLat >= 999999 {
		return nil, fmt.Errorf("no proxy responded")
	}

	a.addLog("auto-select", bestID, "", fmt.Sprintf("fastest: %dms", bestLat))
	return a.ConnectProxy(bestID)
}

// --- Bulk Proxy Checker ---

func (a *App) BulkCheckProxies() []BulkCheckResult {
	a.mu.RLock()
	proxies := make([]*ProxyEntry, 0, len(a.proxies))
	for _, p := range a.proxies {
		proxies = append(proxies, p)
	}
	a.mu.RUnlock()

	resultsCh := make(chan BulkCheckResult, len(proxies))
	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup

	realIP, _ := a.GetMyIP()

	for _, p := range proxies {
		wg.Add(1)
		go func(p *ProxyEntry) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			r := BulkCheckResult{ID: p.ID, Host: p.Host + ":" + p.Port}
			client := a.getOrCreateClient(p)

			start := time.Now()
			resp, err := client.Get("https://api.ipify.org?format=json")
			r.Latency = time.Since(start).Milliseconds()

			if err != nil {
				r.Error = err.Error()
				resultsCh <- r
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			var ipRes struct{ IP string `json:"ip"` }
			json.Unmarshal(body, &ipRes)

			r.Alive = true
			r.ExitIP = ipRes.IP

			if geo := a.fetchGeoInfo(client, ipRes.IP); geo != nil {
				r.Country = geo.Country
				r.CountryCode = geo.CountryCode
			}

			if realIP != "" && r.ExitIP != "" {
				if r.ExitIP == realIP {
					r.Anonymity = "transparent"
				} else {
					r.Anonymity = "elite"
				}
			} else {
				r.Anonymity = "unknown"
			}

			a.mu.Lock()
			p.Latency = r.Latency
			p.LastCheck = time.Now().Format(time.RFC3339)
			a.mu.Unlock()

			resultsCh <- r
		}(p)
	}

	go func() { wg.Wait(); close(resultsCh) }()

	var results []BulkCheckResult
	for r := range resultsCh {
		results = append(results, r)
		runtime.EventsEmit(a.ctx, "bulk-check-progress", r)
	}

	a.addLog("bulk-check", "", "", fmt.Sprintf("checked %d proxies", len(results)))
	return results
}

// --- DNS Leak Test ---

func (a *App) CheckDNSLeak() DNSLeakResult {
	result := DNSLeakResult{}

	a.mu.RLock()
	var activeProxy *ProxyEntry
	for _, p := range a.proxies {
		if p.Active {
			activeProxy = p
			break
		}
	}
	a.mu.RUnlock()

	if activeProxy == nil {
		result.Leaked = true
		return result
	}

	client := a.getOrCreateClient(activeProxy)
	resp, err := client.Get("https://api.ipify.org?format=text")
	if err == nil {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		result.ProxyIP = strings.TrimSpace(string(body))
	}

	directResp, err := http.Get("https://api.ipify.org?format=text")
	if err == nil {
		body, _ := io.ReadAll(directResp.Body)
		directResp.Body.Close()
		result.DNSIP = strings.TrimSpace(string(body))
	}

	if result.ProxyIP != "" && result.DNSIP != "" {
		result.Leaked = result.ProxyIP == result.DNSIP
	}

	return result
}

// --- UWP Loopback Exemption ---

func (a *App) EnableUWPSupport() error {
	script := `@echo off
CheckNetIsolation.exe LoopbackExempt -a -n="Microsoft.MicrosoftEdge_8wekyb3d8bbwe" >nul 2>&1
CheckNetIsolation.exe LoopbackExempt -a -n="Microsoft.Windows.Search_cw5n1h2txyewy" >nul 2>&1
CheckNetIsolation.exe LoopbackExempt -a -n="Microsoft.WindowsStore_8wekyb3d8bbwe" >nul 2>&1
CheckNetIsolation.exe LoopbackExempt -a -n="Microsoft.MSPaint_8wekyb3d8bbwe" >nul 2>&1
CheckNetIsolation.exe LoopbackExempt -a -n="Microsoft.MicrosoftEdge.Stable_8wekyb3d8bbwe" >nul 2>&1
`
	return killswitch.RunElevated(script)
}

// --- Import Formats ---

func (a *App) AddProxiesMultiFormat(raw string) ([]*ProxyEntry, error) {
	lines := strings.Split(raw, "\n")
	var added []*ProxyEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		var p *ProxyEntry
		var err error

		switch {
		case strings.HasPrefix(line, "socks5://") || strings.HasPrefix(line, "socks4://") || strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://"):
			p, err = parseProxyURI(line)
		case strings.Contains(line, ","):
			p, err = parseProxyCSV(line)
		default:
			p, err = parseProxy(line)
		}

		if err != nil || p == nil {
			continue
		}
		a.mu.Lock()
		a.proxies[p.ID] = p
		a.mu.Unlock()
		added = append(added, p)
	}

	if len(added) > 0 {
		a.addLog("batch-added", "", "", fmt.Sprintf("%d proxies (multi-format)", len(added)))
		a.saveProxies()
		runtime.EventsEmit(a.ctx, "proxies-changed")
	}
	return added, nil
}

func parseProxyURI(raw string) (*ProxyEntry, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	protocol := "http"
	if u.Scheme == "socks5" || u.Scheme == "socks4" {
		protocol = "socks5"
	}
	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "8080"
	}
	entry := &ProxyEntry{
		ID:        generateID(),
		Host:      host,
		Port:      port,
		Protocol:  protocol,
		Status:    StatusOffline,
		SessionID: generateSessionID(),
	}
	if u.User != nil {
		entry.Username = u.User.Username()
		entry.Password, _ = u.User.Password()
	}
	return entry, nil
}

func parseProxyCSV(line string) (*ProxyEntry, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid CSV format")
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	entry := &ProxyEntry{
		ID:        generateID(),
		Host:      parts[0],
		Port:      parts[1],
		Protocol:  "http",
		Status:    StatusOffline,
		SessionID: generateSessionID(),
	}
	if len(parts) >= 4 {
		entry.Username = parts[2]
		entry.Password = parts[3]
	}
	if len(parts) >= 5 {
		proto := strings.ToLower(parts[4])
		if proto == "socks5" || proto == "socks4" {
			entry.Protocol = "socks5"
		}
	}
	return entry, nil
}
