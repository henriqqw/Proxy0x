package proxy

import (
	"bufio"
	"net"
	"net/http"
	"strconv"
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

func TestStopReturnsQuicklyWithPlainHTTPInFlight(t *testing.T) {
	host, port := fakeHTTPUpstream(t, okResponse)
	ls := startServer(t, UpstreamConfig{ID: "x", Host: host, Port: port, Protocol: "http"})
	c := dialServer(t, ls)
	sendGet(t, c, "/").Body.Close()

	stopWithin(t, ls, time.Second)
}
