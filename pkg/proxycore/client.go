package proxycore

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	netproxy "golang.org/x/net/proxy"
)

func BuildProxyURL(p *ProxyEntry) string {
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

func BuildHTTPClient(p *ProxyEntry) *http.Client {
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
		dialer, err := netproxy.SOCKS5("tcp", net.JoinHostPort(p.Host, p.Port), auth, netproxy.Direct)
		if err != nil {
			dialer = netproxy.Direct
		}
		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
			TLSHandshakeTimeout: 10 * time.Second,
			MaxIdleConns:        2,
			IdleConnTimeout:     90 * time.Second,
		}
	} else {
		proxyURL, _ := url.Parse(BuildProxyURL(p))
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

	return &http.Client{
		Transport: transport,
		Timeout:   20 * time.Second,
	}
}

func CloseClient(c *http.Client) {
	if t, ok := c.Transport.(*http.Transport); ok {
		t.CloseIdleConnections()
	}
}
