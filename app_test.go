package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"proxy_ui/internal/crypto"
	"proxy_ui/internal/proxy"
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
	geoDelay time.Duration
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
		status, body, geoDelay := f.status, f.body, f.geoDelay
		f.mu.Unlock()
		if strings.HasPrefix(r.URL.Host, "geo.test") {
			time.Sleep(geoDelay)
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

func failingLocalServer(proxy.UpstreamConfig) (*proxy.LocalServer, error) {
	return nil, errors.New("listen failed")
}

func TestLocalServerStartFailureFailsConnect(t *testing.T) {
	a, sys := newTestApp(t)
	_, host, port := startFakeUpstream(t)
	p := addTestProxy(t, a, host, port)
	a.newLocalServer = failingLocalServer

	if _, err := a.ConnectProxy(p.ID); err == nil {
		t.Fatal("connect should fail when the local proxy cannot start")
	}
	if p.Status != StatusError || p.Active {
		t.Fatalf("status=%s active=%v, want error/false", p.Status, p.Active)
	}
	if a.localProxy != nil || sys.isEnabled() {
		t.Fatal("no local proxy or system proxy should be left behind")
	}
}

func TestRotateWithLocalServerFailureDoesNotPanicOnDisconnect(t *testing.T) {
	a, sys := newTestApp(t)
	_, host, port := startFakeUpstream(t)
	p := addTestProxy(t, a, host, port)
	if _, err := a.ConnectProxy(p.ID); err != nil {
		t.Fatal(err)
	}
	a.newLocalServer = failingLocalServer

	a.RotateProxy(p.ID)
	a.DisconnectProxy(p.ID) // used to panic: close of closed channel

	if a.localProxy != nil || sys.isEnabled() {
		t.Fatal("no local proxy or system proxy should be left behind")
	}
}

func TestRemoveActiveProxyDisconnectsFirst(t *testing.T) {
	a, sys := newTestApp(t)
	_, host, port := startFakeUpstream(t)
	p := addTestProxy(t, a, host, port)
	if _, err := a.ConnectProxy(p.ID); err != nil {
		t.Fatal(err)
	}

	a.RemoveProxy(p.ID)

	if a.localProxy != nil || sys.isEnabled() {
		t.Fatal("removed proxy must not keep routing traffic")
	}
	if n := len(a.clients); n != 0 {
		t.Fatalf("%d HTTP clients leaked", n)
	}
	if len(a.GetProxies()) != 0 {
		t.Fatal("proxy not removed")
	}
}

func TestConnectRequiresHTTP200AndValidIP(t *testing.T) {
	cases := map[string]struct {
		status int
		body   string
	}{
		"rate limited": {429, `{"ip":"203.0.113.7"}`},
		"server error": {500, ``},
		"garbage body": {200, `<html>blocked</html>`},
		"empty ip":     {200, `{"ip":""}`},
		"invalid ip":   {200, `{"ip":"not-an-ip"}`},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			a, sys := newTestApp(t)
			up, host, port := startFakeUpstream(t)
			up.set(tc.status, tc.body)
			p := addTestProxy(t, a, host, port)

			if _, err := a.ConnectProxy(p.ID); err == nil {
				t.Fatal("connect should fail")
			}
			if p.Status != StatusError || p.Active || sys.isEnabled() {
				t.Fatalf("status=%s active=%v sysproxy=%v, want error/false/false", p.Status, p.Active, sys.isEnabled())
			}
		})
	}
}

func TestHealthCheckFailsOnNon200(t *testing.T) {
	a, _ := newTestApp(t)
	up, host, port := startFakeUpstream(t)
	p := addTestProxy(t, a, host, port)
	if _, err := a.ConnectProxy(p.ID); err != nil {
		t.Fatal(err)
	}
	up.set(503, ``)

	a.CheckHealth(p.ID)

	if p.Status != StatusError {
		t.Fatalf("status=%s, want error", p.Status)
	}
}

func TestHealthCheckDoesNotHoldLockDuringGeoLookup(t *testing.T) {
	a, _ := newTestApp(t)
	up, host, port := startFakeUpstream(t)
	p := addTestProxy(t, a, host, port)
	if _, err := a.ConnectProxy(p.ID); err != nil {
		t.Fatal(err)
	}
	up.mu.Lock()
	up.body = `{"ip":"198.51.100.9"}` // IP changed -> geo lookup
	up.geoDelay = 500 * time.Millisecond
	up.mu.Unlock()

	go a.CheckHealth(p.ID)
	time.Sleep(150 * time.Millisecond) // health check is now inside the geo lookup

	start := time.Now()
	a.GetProxies()
	if d := time.Since(start); d > 100*time.Millisecond {
		t.Fatalf("GetProxies blocked %v while the health check did network I/O", d)
	}
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
