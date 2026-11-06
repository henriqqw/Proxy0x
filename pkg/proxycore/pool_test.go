package proxycore

import "testing"

func newTestPool(t *testing.T) *Pool {
	t.Helper()
	dir := t.TempDir()
	store := NewStore(dir)
	return NewPool(store)
}

func TestParseProxy(t *testing.T) {
	tests := []struct {
		raw       string
		wantHost  string
		wantPort  string
		wantUser  string
		wantProto string
		wantErr   bool
	}{
		{"1.2.3.4:8080", "1.2.3.4", "8080", "", "http", false},
		{"host.com:3128:user:pass", "host.com", "3128", "user", "http", false},
		{"socks5://1.2.3.4:1080", "1.2.3.4", "1080", "", "socks5", false},
		{"socks5://host:1080:u:p", "host", "1080", "u", "socks5", false},
		{"http://1.2.3.4:8080", "1.2.3.4", "8080", "", "http", false},
		{"", "", "", "", "", true},
		{"noport", "", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			entry, err := ParseProxy(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseProxy(%q) error = %v, wantErr %v", tt.raw, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if entry.Host != tt.wantHost {
				t.Errorf("host = %q, want %q", entry.Host, tt.wantHost)
			}
			if entry.Port != tt.wantPort {
				t.Errorf("port = %q, want %q", entry.Port, tt.wantPort)
			}
			if entry.Username != tt.wantUser {
				t.Errorf("username = %q, want %q", entry.Username, tt.wantUser)
			}
			if entry.Protocol != tt.wantProto {
				t.Errorf("protocol = %q, want %q", entry.Protocol, tt.wantProto)
			}
			if entry.ID == "" {
				t.Error("ID should not be empty")
			}
			if entry.SessionID == "" {
				t.Error("SessionID should not be empty")
			}
		})
	}
}

func TestPoolAddGetRemove(t *testing.T) {
	pool := newTestPool(t)

	entry, err := pool.Add("1.2.3.4:8080:user:pass")
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if entry.Host != "1.2.3.4" {
		t.Errorf("Host = %q, want 1.2.3.4", entry.Host)
	}

	list := pool.List()
	if len(list) != 1 {
		t.Fatalf("List len = %d, want 1", len(list))
	}

	got, ok := pool.Get(entry.ID)
	if !ok || got.ID != entry.ID {
		t.Error("Get failed")
	}

	if err := pool.Remove(entry.ID); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	if len(pool.List()) != 0 {
		t.Error("List should be empty after remove")
	}

	if err := pool.Remove("nonexistent"); err == nil {
		t.Error("Remove nonexistent should error")
	}
}

func TestPoolStats(t *testing.T) {
	pool := newTestPool(t)

	pool.Add("1.2.3.4:8080")
	pool.Add("5.6.7.8:3128")

	stats := pool.Stats()
	if stats.Total != 2 {
		t.Errorf("Total = %d, want 2", stats.Total)
	}
	if stats.Active != 0 {
		t.Errorf("Active = %d, want 0", stats.Active)
	}
}

func TestPoolBestProxyNoActive(t *testing.T) {
	pool := newTestPool(t)
	pool.Add("1.2.3.4:8080")

	_, err := pool.BestProxy()
	if err == nil {
		t.Error("BestProxy should error when no proxy is active")
	}
}

func TestPoolScore(t *testing.T) {
	e := &ProxyEntry{Successes: 10, Failures: 3}
	if e.Score() != 17 {
		t.Errorf("Score() = %d, want 17", e.Score())
	}

	e2 := &ProxyEntry{Successes: 0, Failures: 5}
	if e2.Score() != -5 {
		t.Errorf("Score() = %d, want -5", e2.Score())
	}
}

func TestPoolLogs(t *testing.T) {
	pool := newTestPool(t)

	pool.AddLog("test", "id1", "host1", "d1")
	pool.AddLog("test", "id2", "host2", "d2")

	logs := pool.Logs()
	if len(logs) != 2 {
		t.Fatalf("Logs len = %d, want 2", len(logs))
	}
	if logs[0].ProxyID != "id2" {
		t.Errorf("Logs[0] = %q, want id2 (newest first)", logs[0].ProxyID)
	}
}

func TestPoolPersistence(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	p1 := NewPool(store)
	p1.Add("1.2.3.4:8080:user:pass")
	p1.Add("5.6.7.8:3128")

	p2 := NewPool(store)
	if err := p2.LoadFromStore(); err != nil {
		t.Fatalf("LoadFromStore: %v", err)
	}

	list := p2.List()
	if len(list) != 2 {
		t.Fatalf("Persisted pool len = %d, want 2", len(list))
	}

	var found bool
	for _, e := range list {
		if e.Username == "user" && e.Password == "pass" {
			found = true
		}
	}
	if !found {
		t.Error("Credentials should survive encrypt/decrypt cycle")
	}
}

func TestPoolAddress(t *testing.T) {
	e := &ProxyEntry{Host: "1.2.3.4", Port: "8080"}
	if e.Address() != "1.2.3.4:8080" {
		t.Errorf("Address() = %q, want 1.2.3.4:8080", e.Address())
	}
}

func TestPoolLogOverflow(t *testing.T) {
	pool := newTestPool(t)
	for i := 0; i < 250; i++ {
		pool.AddLog("test", "id", "host", "detail")
	}
	logs := pool.Logs()
	if len(logs) != maxLogs {
		t.Errorf("Logs len = %d, want %d (capped)", len(logs), maxLogs)
	}
}
