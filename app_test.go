package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"proxy_ui/internal/crypto"
)

type fakeSysProxy struct {
	mu      sync.Mutex
	enabled bool
	port    int
}

func (f *fakeSysProxy) Set(port int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.enabled, f.port = true, port
	return nil
}

func (f *fakeSysProxy) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.enabled = false
}

func (f *fakeSysProxy) isEnabled() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.enabled
}

// newTestApp builds an App that never touches the Wails runtime, the real
// registry or the user's %APPDATA%.
func newTestApp(t *testing.T) (*App, *fakeSysProxy) {
	t.Helper()
	dir := t.TempDir()
	sys := &fakeSysProxy{}
	a := NewApp()
	a.ctx = context.Background()
	a.dataFile = filepath.Join(dir, "proxies.json")
	a.crypt = crypto.NewStore(dir)
	a.emit = func(context.Context, string, ...interface{}) {}
	a.sysproxy = sys
	t.Cleanup(func() {
		a.lpMu.Lock()
		if a.localProxy != nil {
			a.localProxy.Stop()
			a.localProxy = nil
		}
		a.lpMu.Unlock()
	})
	return a, sys
}

// fakeUpstream is an HTTP proxy that answers the IP check and geo lookups.
type fakeUpstream struct {
	mu       sync.Mutex
	status   int
	body     string
	requests int
}

func (f *fakeUpstream) set(status int, body string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.status, f.body = status, body
}

func startFakeUpstream(t *testing.T) (*fakeUpstream, string, string) {
	t.Helper()
	f := &fakeUpstream{status: 200, body: `{"ip":"203.0.113.7"}`}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.requests++
		status, body := f.status, f.body
		f.mu.Unlock()
		if strings.HasPrefix(r.URL.Host, "geo.test") {
			w.Write([]byte(`{"country":"Testland","country_code":"TL","city":"Testville"}`))
			return
		}
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	oldIP, oldGeo := ipCheckURL, geoURLFormat
	ipCheckURL, geoURLFormat = "http://ipcheck.test/", "http://geo.test/%s"
	t.Cleanup(func() { ipCheckURL, geoURLFormat = oldIP, oldGeo })

	host, port, _ := net.SplitHostPort(strings.TrimPrefix(srv.URL, "http://"))
	return f, host, port
}

func addTestProxy(t *testing.T, a *App, host, port string) *ProxyEntry {
	t.Helper()
	p, err := a.AddProxy(host + ":" + port + ":user:pass")
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestGetOrCreateClientIsSafeForConcurrentReaders(t *testing.T) {
	for round := 0; round < 20; round++ {
		a, _ := newTestApp(t)
		var wg sync.WaitGroup
		for i := 0; i < 64; i++ {
			p := &ProxyEntry{ID: "p", Host: "127.0.0.1", Port: "1", Protocol: "http", SessionID: strconv.Itoa(i)}
			wg.Add(1)
			go func() {
				defer wg.Done()
				a.mu.RLock()
				a.getOrCreateClient(p)
				a.mu.RUnlock()
			}()
		}
		wg.Wait()
		if len(a.clients) != 64 {
			t.Fatalf("round %d: %d clients, want 64", round, len(a.clients))
		}
	}
}
