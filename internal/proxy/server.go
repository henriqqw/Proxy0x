package proxy

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	netproxy "golang.org/x/net/proxy"
)

type UpstreamConfig struct {
	ID        string
	Host      string
	Port      string
	Username  string
	Password  string
	SessionID string
	Protocol  string
}

type LocalServer struct {
	listener net.Listener
	upstream UpstreamConfig
	port     int
	stopCh   chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
	connsMu  sync.Mutex
	conns    []net.Conn
}

func NewLocalServer(upstream UpstreamConfig) (*LocalServer, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	return newLocalServer(listener, upstream), nil
}

func newLocalServer(listener net.Listener, upstream UpstreamConfig) *LocalServer {
	ls := &LocalServer{
		listener: listener,
		upstream: upstream,
		port:     listener.Addr().(*net.TCPAddr).Port,
		stopCh:   make(chan struct{}),
	}
	go ls.serve()
	return ls
}

func (ls *LocalServer) Port() int {
	return ls.port
}

// UpdateUpstream switches new connections to u; open tunnels and the port are kept.
func (ls *LocalServer) UpdateUpstream(u UpstreamConfig) {
	ls.mu.Lock()
	ls.upstream = u
	ls.mu.Unlock()
}

func (ls *LocalServer) UpstreamID() string {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return ls.upstream.ID
}

func (ls *LocalServer) Stop() {
	close(ls.stopCh)
	ls.listener.Close()

	ls.connsMu.Lock()
	for _, c := range ls.conns {
		c.Close()
	}
	ls.conns = nil
	ls.connsMu.Unlock()

	ls.wg.Wait()
}

const (
	acceptBackoffMin = 5 * time.Millisecond
	acceptBackoffMax = time.Second
)

// upstreamTimeout bounds dialing the upstream proxy and its handshake.
var upstreamTimeout = 15 * time.Second

func (ls *LocalServer) serve() {
	var delay time.Duration
	for {
		conn, err := ls.listener.Accept()
		if err != nil {
			if delay == 0 {
				delay = acceptBackoffMin
			} else {
				delay = min(delay*2, acceptBackoffMax)
			}
			select {
			case <-ls.stopCh:
				return
			case <-time.After(delay):
				continue
			}
		}
		delay = 0
		ls.wg.Add(1)
		go func() {
			defer ls.wg.Done()
			ls.handleConn(conn)
		}()
	}
}

func (ls *LocalServer) trackConn(c net.Conn) {
	ls.connsMu.Lock()
	ls.conns = append(ls.conns, c)
	ls.connsMu.Unlock()
}

func (ls *LocalServer) untrackConn(c net.Conn) {
	ls.connsMu.Lock()
	for i, conn := range ls.conns {
		if conn == c {
			ls.conns = append(ls.conns[:i], ls.conns[i+1:]...)
			break
		}
	}
	ls.connsMu.Unlock()
}

func (ls *LocalServer) proxyAuth() string {
	ls.mu.RLock()
	u := ls.upstream
	ls.mu.RUnlock()

	username := u.Username
	if username != "" && u.SessionID != "" {
		username = username + "_session-" + u.SessionID
	}
	if username == "" {
		return ""
	}
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+u.Password))
}

func (ls *LocalServer) upstreamAddr() string {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return net.JoinHostPort(ls.upstream.Host, ls.upstream.Port)
}

func (ls *LocalServer) dialViaSocks5(target string) (net.Conn, error) {
	ls.mu.RLock()
	u := ls.upstream
	ls.mu.RUnlock()

	username := u.Username
	if username != "" && u.SessionID != "" {
		username = username + "_session-" + u.SessionID
	}
	var auth *netproxy.Auth
	if username != "" {
		auth = &netproxy.Auth{User: username, Password: u.Password}
	}
	dialer, err := netproxy.SOCKS5("tcp", net.JoinHostPort(u.Host, u.Port), auth, &net.Dialer{Timeout: upstreamTimeout})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), upstreamTimeout)
	defer cancel()
	return dialer.(netproxy.ContextDialer).DialContext(ctx, "tcp", target)
}

func (ls *LocalServer) handleConn(clientConn net.Conn) {
	ls.trackConn(clientConn)
	defer func() {
		ls.untrackConn(clientConn)
		clientConn.Close()
	}()
	br := bufio.NewReader(clientConn)
	req, err := http.ReadRequest(br)
	if err != nil {
		return
	}
	if req.Method == "CONNECT" {
		ls.handleConnect(clientConn, req)
	} else {
		ls.handleHTTP(clientConn, req)
	}
}

func (ls *LocalServer) handleConnect(clientConn net.Conn, req *http.Request) {
	ls.mu.RLock()
	proto := ls.upstream.Protocol
	ls.mu.RUnlock()

	var upConn net.Conn
	var err error

	if proto == "socks5" {
		upConn, err = ls.dialViaSocks5(req.Host)
		if err != nil {
			fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			return
		}
		ls.trackConn(upConn)
		defer func() { ls.untrackConn(upConn); upConn.Close() }()
		clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	} else {
		upConn, err = net.DialTimeout("tcp", ls.upstreamAddr(), upstreamTimeout)
		if err != nil {
			fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			return
		}
		ls.trackConn(upConn)
		defer func() { ls.untrackConn(upConn); upConn.Close() }()

		connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n", req.Host, req.Host)
		if auth := ls.proxyAuth(); auth != "" {
			connectReq += "Proxy-Authorization: " + auth + "\r\n"
		}
		connectReq += "\r\n"
		upConn.SetDeadline(time.Now().Add(upstreamTimeout))
		upConn.Write([]byte(connectReq))

		upReader := bufio.NewReader(upConn)
		resp, err := http.ReadResponse(upReader, req)
		if err != nil {
			fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			return
		}
		resp.Body.Close()
		upConn.SetDeadline(time.Time{})

		if resp.StatusCode != 200 {
			fmt.Fprintf(clientConn, "HTTP/1.1 %s\r\n\r\n", resp.Status)
			return
		}

		clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	}

	const idleTimeout = 5 * time.Minute
	done := make(chan struct{}, 2)
	go func() {
		copyWithDeadline(upConn, clientConn, idleTimeout)
		done <- struct{}{}
	}()
	go func() {
		copyWithDeadline(clientConn, upConn, idleTimeout)
		done <- struct{}{}
	}()
	<-done
}

func copyWithDeadline(dst, src net.Conn, idle time.Duration) {
	buf := make([]byte, 32*1024)
	for {
		src.SetReadDeadline(time.Now().Add(idle))
		n, err := src.Read(buf)
		if n > 0 {
			dst.SetWriteDeadline(time.Now().Add(30 * time.Second))
			if _, wErr := dst.Write(buf[:n]); wErr != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}

func (ls *LocalServer) handleHTTP(clientConn net.Conn, req *http.Request) {
	ls.mu.RLock()
	proto := ls.upstream.Protocol
	ls.mu.RUnlock()

	var upConn net.Conn
	var err error

	// One request per client connection: no keep-alive through the local proxy.
	req.Close = true
	req.Header.Del("Proxy-Connection")

	if proto == "socks5" {
		host := req.URL.Host
		if _, _, splitErr := net.SplitHostPort(host); splitErr != nil {
			host = net.JoinHostPort(host, "80")
		}
		upConn, err = ls.dialViaSocks5(host)
		if err != nil {
			fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			return
		}
		ls.trackConn(upConn)
		defer func() { ls.untrackConn(upConn); upConn.Close() }()
		upConn.SetWriteDeadline(time.Now().Add(30 * time.Second))
		req.Write(upConn)
	} else {
		upConn, err = net.DialTimeout("tcp", ls.upstreamAddr(), upstreamTimeout)
		if err != nil {
			fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			return
		}
		ls.trackConn(upConn)
		defer func() { ls.untrackConn(upConn); upConn.Close() }()

		if auth := ls.proxyAuth(); auth != "" {
			req.Header.Set("Proxy-Authorization", auth)
		}
		if !req.URL.IsAbs() {
			req.URL.Scheme = "http"
			req.URL.Host = req.Host
		}

		upConn.SetWriteDeadline(time.Now().Add(30 * time.Second))
		req.WriteProxy(upConn)
	}

	upConn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	resp, err := http.ReadResponse(bufio.NewReader(upConn), req)
	if err != nil {
		fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
		return
	}
	defer resp.Body.Close()
	resp.Close = true
	resp.Write(clientConn)
}
