package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	"proxy_ui/internal/fsutil"
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

var (
	ipCheckURL   = "https://api.ipify.org?format=json"
	myIPURL      = "https://api.ipify.org?format=text"
	geoURLFormat = "https://ipwho.is/%s"
	speedTestURL = "https://speed.cloudflare.com/__down?bytes=102400"
	myIPTimeout  = 10 * time.Second

	newSOCKS5 = netproxy.SOCKS5
)

type systemProxy interface {
	Set(port int) error
	Clear() error
	Recover() error
}

type App struct {
	ctx        context.Context
	proxies    map[string]*ProxyEntry
	clients    map[string]*http.Client
	clientsMu  sync.Mutex
	mu         sync.RWMutex
	lpMu       sync.Mutex
	localProxy *proxy.LocalServer
	stopCh     chan struct{}
	logs       []LogEntry
	logMu      sync.RWMutex
	crypt      *crypto.Store
	storeErr   error // non-nil: proxies.json could not be loaded, never save over it

	dataFile       string
	emit           func(ctx context.Context, name string, data ...interface{})
	sysproxy       systemProxy
	newLocalServer func(proxy.UpstreamConfig) (*proxy.LocalServer, error)
}

func NewApp() *App {
	return &App{
		proxies:        make(map[string]*ProxyEntry),
		clients:        make(map[string]*http.Client),
		stopCh:         make(chan struct{}),
		logs:           make([]LogEntry, 0, 200),
		crypt:          crypto.NewStore(config.DataDir()),
		dataFile:       config.DataFile(),
		emit:           runtime.EventsEmit,
		sysproxy:       proxy.NewSystemProxy(config.SysProxyBackupFile()),
		newLocalServer: proxy.NewLocalServer,
	}
}

func (a *App) emitEvent(name string, data ...interface{}) {
	a.emit(a.ctx, name, data...)
}

// view copies p under the lock for bindings and events, so Wails marshals a
// value no goroutine is writing to. Callers must not hold a.mu.
func (a *App) view(p *ProxyEntry) *ProxyEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return snapshot(p)
}

// snapshot is the UI's copy of p. The password never leaves the backend.
func snapshot(p *ProxyEntry) *ProxyEntry {
	c := *p
	c.Password = ""
	if p.Location != nil {
		loc := *p.Location
		c.Location = &loc
	}
	return &c
}

func (a *App) emitProxy(p *ProxyEntry) {
	a.emitEvent("proxy-updated", a.view(p))
}

// --- Lifecycle ---

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.sysproxy.Recover()
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
		a.sysproxy.Clear()
	}
	a.lpMu.Unlock()
	close(a.stopCh)
}

// --- Window controls ---

func (a *App) MinimizeWindow() { runtime.WindowMinimise(a.ctx) }
func (a *App) CloseWindow()    { runtime.WindowHide(a.ctx) }
func (a *App) ToggleMaximize() { runtime.WindowToggleMaximise(a.ctx) }

// --- Config ---

func (a *App) GetLanguage() string        { return config.Load().Language }
func (a *App) GetTheme() string           { return config.Load().Theme }
func (a *App) GetAutoRotateInterval() int { return config.Load().AutoRotateInterval }

func (a *App) SetLanguage(lang string) {
	cfg := config.Load()
	cfg.Language = lang
	config.Save(cfg)
	a.emitEvent("language-changed", lang)
}

func (a *App) SetTheme(theme string) {
	cfg := config.Load()
	cfg.Theme = theme
	config.Save(cfg)
	a.emitEvent("theme-changed", theme)
}

func (a *App) SetAutoRotateInterval(seconds int) {
	cfg := config.Load()
	cfg.AutoRotateInterval = seconds
	config.Save(cfg)
	a.emitEvent("config-changed", cfg)
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
	a.emitEvent("log-entry", entry)
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

// loadProxies reads the store. If the file exists but cannot be read
// (invalid JSON, undecryptable credentials) the store goes read-only so a
// later save never overwrites data we could not load.
func (a *App) loadProxies() {
	data, err := os.ReadFile(a.dataFile)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		a.markStoreReadOnly(err)
		return
	}
	var entries []*ProxyEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		a.markStoreReadOnly(fmt.Errorf("parse %s: %w", a.dataFile, err))
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, p := range entries {
		user, errU := a.crypt.Decrypt(p.Username)
		pass, errP := a.crypt.Decrypt(p.Password)
		if err := errors.Join(errU, errP); err != nil {
			a.storeErr = fmt.Errorf("decrypt credentials of %s:%s: %w", p.Host, p.Port, err)
		} else {
			p.Username, p.Password = user, pass
		}
		p.Status = StatusOffline
		p.Active = false
		p.CurrentIP = ""
		p.Latency = 0
		a.proxies[p.ID] = p
	}
	if a.storeErr != nil {
		a.emitEvent("store-error", a.storeErr.Error())
	}
}

func (a *App) markStoreReadOnly(err error) {
	a.mu.Lock()
	a.storeErr = err
	a.mu.Unlock()
	a.emitEvent("store-error", err.Error())
}

func (a *App) checkStoreWritable() error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.storeErr != nil {
		return fmt.Errorf("proxy store is read-only: %w", a.storeErr)
	}
	return nil
}

// GetStoreError lets the UI learn about a read-only store on mount; the
// store-error event fires during startup, before the webview listens.
func (a *App) GetStoreError() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.storeErr == nil {
		return ""
	}
	return a.storeErr.Error()
}

func (a *App) saveProxies() error {
	if err := a.checkStoreWritable(); err != nil {
		return err
	}
	a.mu.RLock()
	entries := make([]*ProxyEntry, 0, len(a.proxies))
	for _, p := range a.proxies {
		clone := *p
		user, errU := a.crypt.Encrypt(clone.Username)
		pass, errP := a.crypt.Encrypt(clone.Password)
		if err := errors.Join(errU, errP); err != nil {
			a.mu.RUnlock()
			return fmt.Errorf("encrypt credentials: %w", err)
		}
		clone.Username, clone.Password = user, pass
		clone.Status = StatusOffline
		clone.Active = false
		clone.CurrentIP = ""
		clone.Latency = 0
		entries = append(entries, &clone)
	}
	a.mu.RUnlock()

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	if err := fsutil.WriteFileAtomic(a.dataFile, data, 0o600); err != nil {
		a.emitEvent("store-error", err.Error())
		return err
	}
	return nil
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
	if err := a.checkStoreWritable(); err != nil {
		return nil, err
	}
	p, err := parseProxy(raw)
	if err != nil {
		return nil, err
	}

	a.mu.Lock()
	a.proxies[p.ID] = p
	a.mu.Unlock()

	a.addLog("added", p.ID, p.Host+":"+p.Port, "")
	if err := a.saveProxies(); err != nil {
		return a.view(p), err
	}
	a.emitProxy(p)
	return a.view(p), nil
}

func (a *App) AddProxiesBatch(raw string) ([]*ProxyEntry, error) {
	if err := a.checkStoreWritable(); err != nil {
		return nil, err
	}
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
	err := a.saveProxies()
	a.emitEvent("proxies-changed")
	views := make([]*ProxyEntry, len(added))
	for i, p := range added {
		views[i] = a.view(p)
	}
	return views, err
}

func (a *App) RemoveProxy(id string) error {
	if err := a.checkStoreWritable(); err != nil {
		return err
	}
	a.mu.RLock()
	p := a.proxies[id]
	active := p != nil && p.Active
	a.mu.RUnlock()
	if active {
		a.DisconnectProxy(id)
	}
	a.clearLocalProxyIf(id)

	a.mu.Lock()
	host := ""
	if p != nil {
		host = p.Host + ":" + p.Port
		a.closeClient(id + ":" + p.SessionID)
	}
	delete(a.proxies, id)
	a.mu.Unlock()
	a.addLog("removed", id, host, "")
	err := a.saveProxies()
	a.emitEvent("proxies-changed")
	return err
}

func (a *App) GetProxies() []*ProxyEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make([]*ProxyEntry, 0, len(a.proxies))
	for _, p := range a.proxies {
		result = append(result, snapshot(p))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Order < result[j].Order
	})
	return result
}

func (a *App) ReorderProxies(ids []string) error {
	if err := a.checkStoreWritable(); err != nil {
		return err
	}
	a.mu.Lock()
	for i, id := range ids {
		if p, ok := a.proxies[id]; ok {
			p.Order = i
		}
	}
	a.mu.Unlock()
	return a.saveProxies()
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
	if err := a.checkStoreWritable(); err != nil {
		return nil, err
	}
	a.mu.Lock()
	p, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	p.Name = name
	a.mu.Unlock()
	if err := a.saveProxies(); err != nil {
		return a.view(p), err
	}
	a.emitProxy(p)
	return a.view(p), nil
}

// copiesOf returns unlocked copies of the proxies matching keep, safe to use
// for network checks without holding a.mu.
func (a *App) copiesOf(keep func(*ProxyEntry) bool) []*ProxyEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var out []*ProxyEntry
	for _, p := range a.proxies {
		if keep(p) {
			c := *p
			out = append(out, &c)
		}
	}
	return out
}

func (a *App) SetProxyGroup(id, group string) (*ProxyEntry, error) {
	if err := a.checkStoreWritable(); err != nil {
		return nil, err
	}
	a.mu.Lock()
	p, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	p.Group = group
	a.mu.Unlock()
	if err := a.saveProxies(); err != nil {
		return a.view(p), err
	}
	a.emitProxy(p)
	return a.view(p), nil
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

// getOrCreateClient reads p's fields, so callers hold a.mu (read or write);
// the clients map has its own lock because readers of a.mu call this concurrently.
func (a *App) getOrCreateClient(p *ProxyEntry) *http.Client {
	clientKey := p.ID + ":" + p.SessionID

	a.clientsMu.Lock()
	defer a.clientsMu.Unlock()
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
		var dialContext func(ctx context.Context, network, addr string) (net.Conn, error)
		dialer, err := newSOCKS5("tcp", net.JoinHostPort(p.Host, p.Port), auth, &net.Dialer{Timeout: 15 * time.Second})
		if err != nil {
			// Never fall back to a direct connection: it would expose the real IP.
			dialContext = func(context.Context, string, string) (net.Conn, error) {
				return nil, fmt.Errorf("socks5 setup: %w", err)
			}
		} else if cd, ok := dialer.(netproxy.ContextDialer); ok {
			dialContext = cd.DialContext
		} else {
			dialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		}
		transport = &http.Transport{
			DialContext:         dialContext,
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
	a.clientsMu.Lock()
	defer a.clientsMu.Unlock()
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

func (a *App) setLocalProxy(p *ProxyEntry, id string) error {
	a.mu.RLock()
	cfg := a.toUpstreamConfig(p)
	a.mu.RUnlock()

	a.lpMu.Lock()
	defer a.lpMu.Unlock()

	if a.localProxy != nil {
		a.localProxy.Stop()
		a.localProxy = nil
		a.sysproxy.Clear()
	}

	ls, err := a.newLocalServer(cfg)
	if err != nil {
		return fmt.Errorf("start local proxy: %w", err)
	}
	if err := a.sysproxy.Set(ls.Port()); err != nil {
		ls.Stop()
		return fmt.Errorf("set system proxy: %w", err)
	}
	a.localProxy = ls
	a.addLog("system-proxy", id, net.JoinHostPort(cfg.Host, cfg.Port), fmt.Sprintf("local proxy on port %d", ls.Port()))
	return nil
}

// updateLocalProxy points the running local server at p's current session
// without restarting it, so the port, the registry and open connections survive.
func (a *App) updateLocalProxy(p *ProxyEntry, id string) error {
	a.mu.RLock()
	cfg := a.toUpstreamConfig(p)
	a.mu.RUnlock()

	a.lpMu.Lock()
	if a.localProxy != nil && a.localProxy.UpstreamID() == id {
		a.localProxy.UpdateUpstream(cfg)
		a.lpMu.Unlock()
		return nil
	}
	a.lpMu.Unlock()
	return a.setLocalProxy(p, id)
}

// verifyExitIP fetches the exit IP through client. It succeeds only with
// HTTP 200 and a body whose "ip" field parses as an IP address.
func verifyExitIP(client *http.Client) (string, int64, error) {
	start := time.Now()
	resp, err := client.Get(ipCheckURL)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	latency := time.Since(start).Milliseconds()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("ip check returned HTTP %d", resp.StatusCode)
	}
	var result struct {
		IP string `json:"ip"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4096)).Decode(&result); err != nil {
		return "", 0, fmt.Errorf("ip check: invalid response: %w", err)
	}
	if net.ParseIP(result.IP) == nil {
		return "", 0, fmt.Errorf("ip check: invalid IP %q", result.IP)
	}
	return result.IP, latency, nil
}

// failLocalProxy marks p as failed after its local proxy could not be (re)started.
func (a *App) failLocalProxy(p *ProxyEntry, id, action string, err error) (*ProxyEntry, error) {
	a.mu.Lock()
	p.Status = StatusError
	p.Active = false
	p.CurrentIP = ""
	p.Location = nil
	p.ConnectedAt = ""
	p.LastCheck = time.Now().Format(time.RFC3339)
	a.mu.Unlock()
	a.addLog(action+"-failed", id, p.Host+":"+p.Port, err.Error())
	a.emitProxy(p)
	return a.view(p), err
}

func (a *App) clearLocalProxyIf(id string) {
	a.lpMu.Lock()
	defer a.lpMu.Unlock()

	if a.localProxy != nil && a.localProxy.UpstreamID() == id {
		a.localProxy.Stop()
		a.localProxy = nil
		a.sysproxy.Clear()
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

	a.emitProxy(p)

	ip, latency, err := verifyExitIP(client)
	if err != nil {
		a.mu.Lock()
		p.Status = StatusError
		p.Latency = 0
		p.CurrentIP = ""
		p.Location = nil
		p.LastCheck = time.Now().Format(time.RFC3339)
		a.mu.Unlock()
		a.addLog("connect-failed", id, p.Host+":"+p.Port, err.Error())
		a.emitProxy(p)
		return a.view(p), fmt.Errorf("connection failed: %w", err)
	}

	geo := a.fetchGeoInfo(client, ip)

	if err := a.setLocalProxy(p, id); err != nil {
		return a.failLocalProxy(p, id, "connect", err)
	}

	a.mu.Lock()
	p.Status = StatusConnected
	p.CurrentIP = ip
	p.Latency = latency
	p.Location = geo
	p.Active = true
	p.LastCheck = time.Now().Format(time.RFC3339)
	p.ConnectedAt = time.Now().Format(time.RFC3339)
	a.mu.Unlock()

	a.addLog("connected", id, p.Host+":"+p.Port, fmt.Sprintf("IP: %s, %dms", ip, latency))
	a.emitProxy(p)
	return a.view(p), nil
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
	a.emitProxy(p)
	return a.view(p), nil
}

func (a *App) RotateProxy(id string) (*ProxyEntry, error) {
	a.mu.Lock()
	p, ok := a.proxies[id]
	if !ok {
		a.mu.Unlock()
		return nil, fmt.Errorf("proxy not found")
	}
	oldSession := p.SessionID
	prevStatus := p.Status
	candidate := *p
	candidate.SessionID = generateSessionID()
	client := a.getOrCreateClient(&candidate)
	p.Status = StatusConnecting
	a.mu.Unlock()

	a.emitProxy(p)

	ip, latency, err := verifyExitIP(client)
	if err != nil {
		a.closeClient(id + ":" + candidate.SessionID)
		a.mu.Lock()
		p.Status = prevStatus
		a.mu.Unlock()
		a.addLog("rotate-failed", id, p.Host+":"+p.Port, err.Error())
		a.emitProxy(p)
		return a.view(p), fmt.Errorf("rotation failed: %w", err)
	}

	geo := a.fetchGeoInfo(client, ip)

	a.closeClient(id + ":" + oldSession)
	a.mu.Lock()
	p.SessionID = candidate.SessionID
	p.Status = StatusConnected
	p.CurrentIP = ip
	p.Latency = latency
	p.Location = geo
	p.Active = true
	p.LastCheck = time.Now().Format(time.RFC3339)
	a.mu.Unlock()

	if err := a.updateLocalProxy(p, id); err != nil {
		return a.failLocalProxy(p, id, "rotate", err)
	}

	a.addLog("rotated", id, p.Host+":"+p.Port, fmt.Sprintf("IP: %s, %dms", ip, latency))
	a.emitProxy(p)
	return a.view(p), nil
}

// --- Health & rotation ---

func (a *App) checkHealth(id string) (*ProxyEntry, error) {
	a.mu.RLock()
	p, ok := a.proxies[id]
	if !ok || !p.Active {
		a.mu.RUnlock()
		return nil, fmt.Errorf("proxy not found or inactive")
	}
	client := a.getOrCreateClient(p)
	prevStatus := p.Status
	a.mu.RUnlock()

	ip, latency, err := verifyExitIP(client)
	if err != nil {
		a.mu.Lock()
		p.Status = StatusError
		p.LastCheck = time.Now().Format(time.RFC3339)
		a.mu.Unlock()
		if prevStatus == StatusConnected {
			a.addLog("dropped", id, p.Host+":"+p.Port, err.Error())
			a.emitEvent("proxy-dropped", map[string]string{
				"id":   id,
				"host": p.Host + ":" + p.Port,
			})
			a.tryFailover(id)
		}
		a.emitProxy(p)
		return a.view(p), nil
	}

	a.mu.RLock()
	ipChanged := ip != p.CurrentIP
	a.mu.RUnlock()
	var geo *GeoInfo
	if ipChanged {
		geo = a.fetchGeoInfo(client, ip)
	}

	a.mu.Lock()
	if p.Active {
		p.Status = StatusConnected
	}
	p.Latency = latency
	p.LastCheck = time.Now().Format(time.RFC3339)
	if ipChanged {
		p.CurrentIP = ip
		p.Location = geo
	}
	a.mu.Unlock()

	a.emitProxy(p)
	return a.view(p), nil
}

func (a *App) fetchGeoInfo(client *http.Client, ip string) *GeoInfo {
	resp, err := client.Get(fmt.Sprintf(geoURLFormat, ip))
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
				a.checkHealth(id)
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

	_, connectMs, err := verifyExitIP(client)
	if err != nil {
		return nil, fmt.Errorf("connect failed: %w", err)
	}
	results["connectMs"] = connectMs

	start := time.Now()
	resp, err := client.Get(speedTestURL)
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
	resp, err := (&http.Client{Timeout: myIPTimeout}).Get(myIPURL)
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
	proxies := a.copiesOf(func(p *ProxyEntry) bool { return !p.Active })

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
			_, lat, err := verifyExitIP(a.getOrCreateClient(p))
			if err != nil {
				results <- result{p.ID, 999999}
				return
			}
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
	proxies := a.copiesOf(func(*ProxyEntry) bool { return true })

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

			ip, latency, err := verifyExitIP(client)
			r.Latency = latency
			if err != nil {
				r.Error = err.Error()
				resultsCh <- r
				return
			}

			r.Alive = true
			r.ExitIP = ip

			if geo := a.fetchGeoInfo(client, ip); geo != nil {
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
			if live, ok := a.proxies[p.ID]; ok {
				live.Latency = r.Latency
				live.LastCheck = time.Now().Format(time.RFC3339)
			}
			a.mu.Unlock()

			resultsCh <- r
		}(p)
	}

	go func() { wg.Wait(); close(resultsCh) }()

	var results []BulkCheckResult
	for r := range resultsCh {
		results = append(results, r)
		a.emitEvent("bulk-check-progress", r)
	}

	a.addLog("bulk-check", "", "", fmt.Sprintf("checked %d proxies", len(results)))
	return results
}

// --- DNS Leak Test ---

func (a *App) CheckDNSLeak() DNSLeakResult {
	result := DNSLeakResult{}

	active := a.copiesOf(func(p *ProxyEntry) bool { return p.Active })
	if len(active) == 0 {
		result.Leaked = true
		return result
	}

	if ip, _, err := verifyExitIP(a.getOrCreateClient(active[0])); err == nil {
		result.ProxyIP = ip
	}
	if ip, err := a.GetMyIP(); err == nil {
		result.DNSIP = strings.TrimSpace(ip)
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
	if err := a.checkStoreWritable(); err != nil {
		return nil, err
	}
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

	var err error
	if len(added) > 0 {
		a.addLog("batch-added", "", "", fmt.Sprintf("%d proxies (multi-format)", len(added)))
		err = a.saveProxies()
		a.emitEvent("proxies-changed")
	}
	views := make([]*ProxyEntry, len(added))
	for i, p := range added {
		views[i] = a.view(p)
	}
	return views, err
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
