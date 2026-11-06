package proxy

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

// fakeHTTPUpstream answers every proxied request with the given raw response
// and keeps the connection open, like a keep-alive upstream proxy.
func fakeHTTPUpstream(t *testing.T, rawResponse string) (host, port string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				br := bufio.NewReader(c)
				for {
					req, err := http.ReadRequest(br)
					if err != nil {
						return
					}
					req.Body.Close()
					c.Write([]byte(rawResponse))
				}
			}(c)
		}
	}()
	host, port, _ = net.SplitHostPort(ln.Addr().String())
	return host, port
}

const okResponse = "HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\nok"

func startServer(t *testing.T, upstream UpstreamConfig) *LocalServer {
	t.Helper()
	ls, err := NewLocalServer(upstream)
	if err != nil {
		t.Fatal(err)
	}
	return ls
}

func dialServer(t *testing.T, ls *LocalServer) net.Conn {
	t.Helper()
	c, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(ls.Port()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return c
}

func sendGet(t *testing.T, c net.Conn, path string) *http.Response {
	t.Helper()
	c.Write([]byte("GET http://example.com" + path + " HTTP/1.1\r\nHost: example.com\r\n\r\n"))
	c.SetReadDeadline(time.Now().Add(3 * time.Second))
	resp, err := http.ReadResponse(bufio.NewReader(c), nil)
	if err != nil {
		t.Fatalf("no response for %s: %v", path, err)
	}
	return resp
}

func stopWithin(t *testing.T, ls *LocalServer, limit time.Duration) {
	t.Helper()
	done := make(chan struct{})
	go func() { ls.Stop(); close(done) }()
	select {
	case <-done:
	case <-time.After(limit):
		t.Fatalf("Stop still blocked after %v", limit)
	}
}

func TestPlainHTTPClosesClientAfterResponse(t *testing.T) {
	cases := map[string]struct {
		raw  string
		body string
	}{
		"content-length": {okResponse, "ok"},
		"chunked":        {"HTTP/1.1 200 OK\r\nTransfer-Encoding: chunked\r\n\r\n5\r\nhello\r\n0\r\n\r\n", "hello"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			host, port := fakeHTTPUpstream(t, tc.raw)
			ls := startServer(t, UpstreamConfig{ID: "x", Host: host, Port: port, Protocol: "http"})
			defer ls.Stop()
			c := dialServer(t, ls)

			c.Write([]byte("GET http://example.com/ HTTP/1.1\r\nHost: example.com\r\nProxy-Connection: keep-alive\r\n\r\n"))
			c.SetReadDeadline(time.Now().Add(2 * time.Second))
			all, err := io.ReadAll(c)
			if err != nil {
				t.Fatalf("connection not closed after the response (keep-alive would hang the browser): %v", err)
			}
			resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(all)), nil)
			if err != nil {
				t.Fatalf("bad response: %v", err)
			}
			body, _ := io.ReadAll(resp.Body)
			if string(body) != tc.body {
				t.Fatalf("body = %q, want %q", body, tc.body)
			}
			if !resp.Close {
				t.Fatalf("response should announce Connection: close")
			}
		})
	}
}

// connectUpstream accepts CONNECT, records Proxy-Authorization and echoes the tunnel.
func connectUpstream(t *testing.T) (host, port string, auths chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	auths = make(chan string, 10)
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				br := bufio.NewReader(c)
				req, err := http.ReadRequest(br)
				if err != nil {
					return
				}
				auths <- req.Header.Get("Proxy-Authorization")
				c.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
				io.Copy(c, br)
			}(c)
		}
	}()
	host, port, _ = net.SplitHostPort(ln.Addr().String())
	return host, port, auths
}

func openTunnel(t *testing.T, ls *LocalServer) net.Conn {
	t.Helper()
	c := dialServer(t, ls)
	c.Write([]byte("CONNECT example.com:443 HTTP/1.1\r\nHost: example.com:443\r\n\r\n"))
	c.SetReadDeadline(time.Now().Add(2 * time.Second))
	br := bufio.NewReader(c)
	resp, err := http.ReadResponse(br, nil)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("tunnel not established: %v", err)
	}
	return c
}

func echo(t *testing.T, c net.Conn, msg string) {
	t.Helper()
	c.Write([]byte(msg))
	buf := make([]byte, len(msg))
	c.SetReadDeadline(time.Now().Add(2 * time.Second))
	if _, err := io.ReadFull(c, buf); err != nil || string(buf) != msg {
		t.Fatalf("tunnel broken: got %q, err %v", buf, err)
	}
}

func TestUpdateUpstreamKeepsOpenTunnelsAndPort(t *testing.T) {
	host, port, auths := connectUpstream(t)
	cfg := UpstreamConfig{ID: "x", Host: host, Port: port, Protocol: "http", Username: "u", Password: "p", SessionID: "s1"}
	ls := startServer(t, cfg)
	defer ls.Stop()
	before := ls.Port()

	c1 := openTunnel(t, ls)
	<-auths
	echo(t, c1, "ping")

	cfg.SessionID = "s2"
	ls.UpdateUpstream(cfg)

	echo(t, c1, "still-open")
	openTunnel(t, ls)
	want := "Basic " + base64.StdEncoding.EncodeToString([]byte("u_session-s2:p"))
	if got := <-auths; got != want {
		t.Fatalf("new connection used auth %q, want new session %q", got, want)
	}
	if ls.Port() != before {
		t.Fatal("port changed")
	}
}

func silentUpstream(t *testing.T) (host, port string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			t.Cleanup(func() { c.Close() })
		}
	}()
	host, port, _ = net.SplitHostPort(ln.Addr().String())
	return host, port
}

func TestSilentUpstreamTimesOut(t *testing.T) {
	old := upstreamTimeout
	upstreamTimeout = 200 * time.Millisecond
	t.Cleanup(func() { upstreamTimeout = old })

	for _, proto := range []string{"http", "socks5"} {
		t.Run(proto, func(t *testing.T) {
			host, port := silentUpstream(t)
			ls := startServer(t, UpstreamConfig{ID: "x", Host: host, Port: port, Protocol: proto})
			defer ls.Stop()
			c := dialServer(t, ls)

			c.Write([]byte("CONNECT example.com:443 HTTP/1.1\r\nHost: example.com:443\r\n\r\n"))
			c.SetReadDeadline(time.Now().Add(2 * time.Second))
			resp, err := http.ReadResponse(bufio.NewReader(c), nil)
			if err != nil {
				t.Fatalf("no answer within 2s from a silent upstream: %v", err)
			}
			if resp.StatusCode != http.StatusBadGateway {
				t.Fatalf("status = %d, want 502", resp.StatusCode)
			}
		})
	}
}

type failingListener struct {
	calls  atomic.Int64
	closed chan struct{}
}

func (l *failingListener) Accept() (net.Conn, error) {
	l.calls.Add(1)
	select {
	case <-l.closed:
		return nil, net.ErrClosed
	default:
		return nil, errors.New("accept: too many open files")
	}
}
func (l *failingListener) Close() error   { close(l.closed); return nil }
func (l *failingListener) Addr() net.Addr { return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1} }

func TestAcceptErrorsBackOff(t *testing.T) {
	ln := &failingListener{closed: make(chan struct{})}
	ls := newLocalServer(ln, UpstreamConfig{})
	time.Sleep(100 * time.Millisecond)
	stopWithin(t, ls, time.Second)

	if n := ln.calls.Load(); n > 20 {
		t.Fatalf("Accept called %d times in 100ms: busy loop instead of backoff", n)
	}
}

func TestStopReturnsQuicklyWithPlainHTTPInFlight(t *testing.T) {
	host, port := fakeHTTPUpstream(t, okResponse)
	ls := startServer(t, UpstreamConfig{ID: "x", Host: host, Port: port, Protocol: "http"})
	c := dialServer(t, ls)
	sendGet(t, c, "/").Body.Close()

	stopWithin(t, ls, time.Second)
}
