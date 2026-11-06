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
)

type LocalProxy struct {
	listener net.Listener
	upstream *ProxyEntry
	port     int
	stopCh   chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
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

func (lp *LocalProxy) handleConn(clientConn net.Conn) {
	defer clientConn.Close()
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
	upConn, err := net.DialTimeout("tcp", lp.upstreamAddr(), 15*time.Second)
	if err != nil {
		fmt.Fprintf(clientConn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
		return
	}
	defer upConn.Close()

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

	done := make(chan struct{}, 2)
	go func() { io.Copy(upConn, clientConn); done <- struct{}{} }()
	go func() { io.Copy(clientConn, upConn); done <- struct{}{} }()
	<-done
}

func (lp *LocalProxy) handleHTTP(clientConn net.Conn, req *http.Request) {
	upConn, err := net.DialTimeout("tcp", lp.upstreamAddr(), 15*time.Second)
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
	req.WriteProxy(upConn)
	io.Copy(clientConn, upConn)
}

func (lp *LocalProxy) stop() {
	close(lp.stopCh)
	lp.listener.Close()
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
	ps := `Add-Type 'using System;using System.Runtime.InteropServices;public class WI{[DllImport("wininet.dll")]public static extern bool InternetSetOption(IntPtr h,int o,IntPtr b,int l);}';[WI]::InternetSetOption([IntPtr]::Zero,39,[IntPtr]::Zero,0);[WI]::InternetSetOption([IntPtr]::Zero,37,[IntPtr]::Zero,0)`
	runHidden("powershell", "-NoProfile", "-Command", ps)
}
