package proxycore

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const maxLogs = 200

type Pool struct {
	mu      sync.RWMutex
	proxies map[string]*ProxyEntry

	clientsMu sync.Mutex
	clients   map[string]*http.Client

	logMu sync.RWMutex
	logs  []LogEntry

	store *Store
}

func NewPool(store *Store) *Pool {
	return &Pool{
		proxies: make(map[string]*ProxyEntry),
		clients: make(map[string]*http.Client),
		logs:    make([]LogEntry, 0, maxLogs),
		store:   store,
	}
}

func (p *Pool) AddLog(action, proxyID, host, detail string) {
	entry := LogEntry{
		Time:    time.Now().Format(time.RFC3339),
		Action:  action,
		ProxyID: proxyID,
		Host:    host,
		Detail:  detail,
	}
	p.logMu.Lock()
	p.logs = append(p.logs, entry)
	if len(p.logs) > maxLogs {
		p.logs = p.logs[len(p.logs)-maxLogs:]
	}
	p.logMu.Unlock()
}

func (p *Pool) Logs() []LogEntry {
	p.logMu.RLock()
	defer p.logMu.RUnlock()
	result := make([]LogEntry, len(p.logs))
	copy(result, p.logs)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func GenerateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateSessionID() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func ParseProxy(raw string) (*ProxyEntry, error) {
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
		ID:        GenerateID(),
		Host:      parts[0],
		Port:      parts[1],
		Protocol:  protocol,
		Status:    StatusOffline,
		SessionID: GenerateSessionID(),
	}

	if len(parts) >= 4 {
		entry.Username = parts[2]
		entry.Password = parts[3]
	}

	return entry, nil
}

func (p *Pool) LoadFromStore() error {
	entries, err := p.store.Load()
	if err != nil {
		return err
	}
	p.mu.Lock()
	for _, e := range entries {
		p.proxies[e.ID] = e
	}
	p.mu.Unlock()
	return nil
}

func (p *Pool) List() []*ProxyEntry {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]*ProxyEntry, 0, len(p.proxies))
	for _, e := range p.proxies {
		result = append(result, e)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Order < result[j].Order
	})
	return result
}

func (p *Pool) Get(id string) (*ProxyEntry, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	e, ok := p.proxies[id]
	return e, ok
}

func (p *Pool) Add(raw string) (*ProxyEntry, error) {
	entry, err := ParseProxy(raw)
	if err != nil {
		return nil, err
	}
	p.mu.Lock()
	p.proxies[entry.ID] = entry
	p.mu.Unlock()
	p.AddLog("added", entry.ID, entry.Address(), "")
	return entry, nil
}

func (p *Pool) Remove(id string) error {
	p.mu.Lock()
	entry, ok := p.proxies[id]
	if !ok {
		p.mu.Unlock()
		return fmt.Errorf("proxy %s not found", id)
	}
	host := entry.Address()
	delete(p.proxies, id)
	p.closeClient(id + ":" + entry.SessionID)
	p.mu.Unlock()
	p.AddLog("removed", id, host, "")
	return nil
}

// getOrCreateClient is called under p.mu read locks, so the map has its own lock.
func (p *Pool) getOrCreateClient(proxy *ProxyEntry) *http.Client {
	clientKey := proxy.ID + ":" + proxy.SessionID
	p.clientsMu.Lock()
	defer p.clientsMu.Unlock()
	if c, ok := p.clients[clientKey]; ok {
		return c
	}
	c := BuildHTTPClient(proxy)
	p.clients[clientKey] = c
	return c
}

func (p *Pool) closeClient(key string) {
	p.clientsMu.Lock()
	defer p.clientsMu.Unlock()
	if c, ok := p.clients[key]; ok {
		CloseClient(c)
		delete(p.clients, key)
	}
}

func (p *Pool) Connect(id string) (*ProxyEntry, error) {
	p.mu.RLock()
	var activeID string
	for pid, pe := range p.proxies {
		if pe.Active && pid != id {
			activeID = pid
			break
		}
	}
	p.mu.RUnlock()
	if activeID != "" {
		p.Disconnect(activeID)
	}

	p.mu.Lock()
	proxy, ok := p.proxies[id]
	if !ok {
		p.mu.Unlock()
		return nil, fmt.Errorf("proxy %s not found", id)
	}
	proxy.Status = StatusConnecting
	client := p.getOrCreateClient(proxy)
	p.mu.Unlock()

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	latency := time.Since(start).Milliseconds()

	if err != nil {
		p.mu.Lock()
		proxy.Status = StatusError
		proxy.Latency = 0
		proxy.CurrentIP = ""
		proxy.Location = nil
		proxy.LastCheck = time.Now().Format(time.RFC3339)
		p.mu.Unlock()
		p.AddLog("connect-failed", id, proxy.Address(), err.Error())
		return proxy, fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ipResult struct {
		IP string `json:"ip"`
	}
	json.Unmarshal(body, &ipResult)

	geo := FetchGeoInfo(client, ipResult.IP)

	p.mu.Lock()
	proxy.Status = StatusConnected
	proxy.CurrentIP = ipResult.IP
	proxy.Latency = latency
	proxy.Location = geo
	proxy.Active = true
	proxy.LastCheck = time.Now().Format(time.RFC3339)
	proxy.ConnectedAt = time.Now().Format(time.RFC3339)
	p.mu.Unlock()

	p.AddLog("connected", id, proxy.Address(), fmt.Sprintf("IP: %s, %dms", ipResult.IP, latency))
	return proxy, nil
}

func (p *Pool) Disconnect(id string) (*ProxyEntry, error) {
	p.mu.Lock()
	proxy, ok := p.proxies[id]
	if !ok {
		p.mu.Unlock()
		return nil, fmt.Errorf("proxy %s not found", id)
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
	p.closeClient(id + ":" + proxy.SessionID)
	p.mu.Unlock()

	p.AddLog("disconnected", id, proxy.Address(), "")
	return proxy, nil
}

func (p *Pool) Rotate(id string) (*ProxyEntry, error) {
	p.mu.Lock()
	proxy, ok := p.proxies[id]
	if !ok {
		p.mu.Unlock()
		return nil, fmt.Errorf("proxy %s not found", id)
	}
	oldKey := id + ":" + proxy.SessionID
	p.closeClient(oldKey)
	proxy.SessionID = GenerateSessionID()
	proxy.Status = StatusConnecting
	client := p.getOrCreateClient(proxy)
	p.mu.Unlock()

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	latency := time.Since(start).Milliseconds()

	if err != nil {
		p.mu.Lock()
		proxy.Status = StatusError
		p.mu.Unlock()
		p.AddLog("rotate-failed", id, proxy.Address(), err.Error())
		return proxy, fmt.Errorf("rotation failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ipResult struct {
		IP string `json:"ip"`
	}
	json.Unmarshal(body, &ipResult)

	geo := FetchGeoInfo(client, ipResult.IP)

	p.mu.Lock()
	proxy.Status = StatusConnected
	proxy.CurrentIP = ipResult.IP
	proxy.Latency = latency
	proxy.Location = geo
	proxy.Active = true
	proxy.LastCheck = time.Now().Format(time.RFC3339)
	p.mu.Unlock()

	p.AddLog("rotated", id, proxy.Address(), fmt.Sprintf("IP: %s, %dms", ipResult.IP, latency))
	return proxy, nil
}

func (p *Pool) CheckHealth(id string) (*ProxyEntry, error) {
	p.mu.RLock()
	proxy, ok := p.proxies[id]
	if !ok || !proxy.Active {
		p.mu.RUnlock()
		return nil, fmt.Errorf("proxy %s not found or inactive", id)
	}
	client := p.getOrCreateClient(proxy)
	p.mu.RUnlock()

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	latency := time.Since(start).Milliseconds()

	if err != nil {
		p.mu.Lock()
		proxy.Status = StatusError
		proxy.LastCheck = time.Now().Format(time.RFC3339)
		p.mu.Unlock()
		p.AddLog("health-failed", id, proxy.Address(), err.Error())
		return proxy, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ipResult struct {
		IP string `json:"ip"`
	}
	json.Unmarshal(body, &ipResult)

	p.mu.Lock()
	proxy.Latency = latency
	proxy.LastCheck = time.Now().Format(time.RFC3339)
	if ipResult.IP != "" && ipResult.IP != proxy.CurrentIP {
		proxy.CurrentIP = ipResult.IP
		geo := FetchGeoInfo(client, ipResult.IP)
		proxy.Location = geo
	}
	p.mu.Unlock()

	return proxy, nil
}

func (p *Pool) Stats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := PoolStats{Total: len(p.proxies)}
	var totalLatency int64
	var latencyCount int

	for _, e := range p.proxies {
		if e.Active {
			stats.Active++
			if e.Status == StatusConnected {
				stats.Healthy++
			}
		}
		if e.Status == StatusError {
			stats.Errored++
		}
		if e.Latency > 0 {
			totalLatency += e.Latency
			latencyCount++
		}
	}

	if latencyCount > 0 {
		stats.AvgLatency = float64(totalLatency) / float64(latencyCount)
	}
	return stats
}

// ActiveProxy returns the connected proxy; only one can be active at a time.
func (p *Pool) ActiveProxy() (*ProxyEntry, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, e := range p.proxies {
		if e.Active && e.Status == StatusConnected {
			return e, nil
		}
	}
	return nil, fmt.Errorf("no active proxy available")
}

func (p *Pool) ClientFor(proxy *ProxyEntry) *http.Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.getOrCreateClient(proxy)
}

func FetchGeoInfo(client *http.Client, ip string) *GeoInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://ipwho.is/%s", ip), nil)
	if err != nil {
		return nil
	}
	resp, err := client.Do(req)
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
			ISP string `json:"isp"`
			Org string `json:"org"`
			ASN int    `json:"asn"`
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
