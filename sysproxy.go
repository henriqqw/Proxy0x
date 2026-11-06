package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"sync"
	"syscall"
	"time"
	"unsafe"

	netproxy "golang.org/x/net/proxy"
)

var wininet = syscall.NewLazyDLL("wininet.dll")
var internetSetOption = wininet.NewProc("InternetSetOptionW")

type LocalProxy struct {
	listener net.Listener
	upstream *ProxyEntry
	port     int
	stopCh   chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
	connsMu  sync.Mutex
	conns    []net.Conn
}

func newLocalProxy(upstream *ProxyEntry) (*LocalProxy, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	lp := &LocalProxy{
		listener: listener,
		upstream: upstream,
		port:     listener.Addr().(*net.TCPAddr).Port,
		stopCh:   make(chan struct{}),
	}
	go lp.serve()
	return lp, nil
}

func (lp *LocalProxy) serve() {
	for {
		conn, err := lp.listener.Accept()
		if err != nil {
			select {
			case <-lp.stopCh:
				return
			default:
				continue
			}
		}
		lp.wg.Add(1)
		go func() {
			defer lp.wg.Done()
			lp.handleConn(conn)
		}()
	}
}

func (lp *LocalProxy) proxyAuth() string {
	lp.mu.RLock()
	u := lp.upstream
	lp.mu.RUnlock()

	username := u.Username
	if username != "" && u.SessionID != "" {
		username = username + "_session-" + u.SessionID
	}
	if username == "" {
		return ""
	}
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+u.Password))
}

func (lp *LocalProxy) upstreamAddr() string {
	lp.mu.RLock()
	defer lp.mu.RUnlock()
	return net.JoinHostPort(lp.upstream.Host, lp.upstream.Port)
}

func (lp *LocalProxy) dialViaSocks5(target string) (net.Conn, error) {
	lp.mu.RLock()
	u := lp.upstream
	lp.mu.RUnlock()

	username := u.Username
	if username != "" && u.SessionID != "" {
		username = username + "_session-" + u.SessionID
	}
	var auth *netproxy.Auth
	if username != "" {
		auth = &netproxy.Auth{User: username, Password: u.Password}
	}
	dialer, err := netproxy.SOCKS5("tcp", net.JoinHostPort(u.Host, u.Port), auth, netproxy.Direct)
	if err != nil {
		return nil, err
	}
	return dialer.Dial("tcp", target)
}

func (lp *LocalProxy) trackConn(c net.Conn) {
	lp.connsMu.Lock()
	lp.conns = append(lp.conns, c)
	lp.connsMu.Unlock()
}

func (lp *LocalProxy) untrackConn(c net.Conn) {
	lp.connsMu.Lock()
	for i, conn := range lp.conns {
		if conn == c {
			lp.conns = append(lp.conns[:i], lp.conns[i+1:]...)
			break
		}
	}
	lp.connsMu.Unlock()
}

func (lp *LocalProxy) handleConn(clientConn net.Conn) {
	lp.trackConn(clientConn)
	defer func() {
		lp.untrackConn(clientConn)
		clientConn.Close()
	}()
	br := bufio.NewReader(clientConn)
	req, err := http.ReadRequest(br)
	if err != nil {
		return
	}
	if req.Method == "CONNECT" {
		lp.handleConnect(clientConn, req)
	} else {
		lp.handleHTTP(clientConn, req)
	}
}

func (lp *LocalProxy) handleConnect(clientConn net.Conn, req *http.Request) {
	lp.mu.RLock()
	proto := lp.upstream.Protocol
	lp.mu.RUnlock()

	var upConn net.Conn
	var err error

	if proto == "socks5" {
		upConn, err = lp.dialViaSocks5(req.Host)
		if err != nil {
			fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			return
		}
		lp.trackConn(upConn)
		defer func() { lp.untrackConn(upConn); upConn.Close() }()
		clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	} else {
		upConn, err = net.DialTimeout("tcp", lp.upstreamAddr(), 15*time.Second)
		if err != nil {
			fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			return
		}
		lp.trackConn(upConn)
		defer func() { lp.untrackConn(upConn); upConn.Close() }()

		connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n", req.Host, req.Host)
		if auth := lp.proxyAuth(); auth != "" {
			connectReq += "Proxy-Authorization: " + auth + "\r\n"
		}
		connectReq += "\r\n"
		upConn.Write([]byte(connectReq))

		upReader := bufio.NewReader(upConn)
		resp, err := http.ReadResponse(upReader, req)
		if err != nil {
			fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			return
		}
		resp.Body.Close()

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

func (lp *LocalProxy) handleHTTP(clientConn net.Conn, req *http.Request) {
	lp.mu.RLock()
	proto := lp.upstream.Protocol
	lp.mu.RUnlock()

	var upConn net.Conn
	var err error

	if proto == "socks5" {
		host := req.URL.Host
		if _, _, splitErr := net.SplitHostPort(host); splitErr != nil {
			host = net.JoinHostPort(host, "80")
		}
		upConn, err = lp.dialViaSocks5(host)
		if err != nil {
			fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			return
		}
		defer upConn.Close()
		upConn.SetWriteDeadline(time.Now().Add(30 * time.Second))
		req.Write(upConn)
	} else {
		upConn, err = net.DialTimeout("tcp", lp.upstreamAddr(), 15*time.Second)
		if err != nil {
			fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
			return
		}
		defer upConn.Close()

		if auth := lp.proxyAuth(); auth != "" {
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
	io.Copy(clientConn, upConn)
}

func (lp *LocalProxy) stop() {
	close(lp.stopCh)
	lp.listener.Close()

	lp.connsMu.Lock()
	for _, c := range lp.conns {
		c.Close()
	}
	lp.conns = nil
	lp.connsMu.Unlock()

	lp.wg.Wait()
}

func runHidden(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	return cmd.Run()
}

const regPath = `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`

func setSystemProxy(port int) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	bypass := "localhost;127.*;10.*;192.168.*;<local>"

	runHidden("reg", "add", regPath, "/v", "ProxyServer", "/t", "REG_SZ", "/d", addr, "/f")
	runHidden("reg", "add", regPath, "/v", "ProxyOverride", "/t", "REG_SZ", "/d", bypass, "/f")
	runHidden("reg", "add", regPath, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f")
	refreshProxy()
	return nil
}

func clearSystemProxy() {
	runHidden("reg", "add", regPath, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f")
	refreshProxy()
}

func refreshProxy() {
	internetSetOption.Call(0, 39, uintptr(unsafe.Pointer(nil)), 0)
	internetSetOption.Call(0, 37, uintptr(unsafe.Pointer(nil)), 0)
}
