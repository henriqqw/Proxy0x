package proxycore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	MaxResponseBody = 100 * 1024 // 100KB
	MaxRedirects    = 10
)

type Fetcher struct {
	pool *Pool
}

func NewFetcher(pool *Pool) *Fetcher {
	return &Fetcher{pool: pool}
}

func (f *Fetcher) Fetch(ctx context.Context, targetURL, method string, headers map[string]string, body string, proxyID string) (*FetchResult, error) {
	if err := ValidateURL(targetURL); err != nil {
		return nil, fmt.Errorf("SSRF blocked: %w", err)
	}

	if method == "" {
		method = "GET"
	}
	method = strings.ToUpper(method)

	var proxy *ProxyEntry
	var err error
	if proxyID != "" {
		proxy, _ = f.pool.Get(proxyID)
		if proxy == nil || !proxy.Active {
			return nil, fmt.Errorf("proxy %s not found or not active", proxyID)
		}
	} else {
		proxy, err = f.pool.ActiveProxy()
		if err != nil {
			return nil, err
		}
	}

	baseClient := f.pool.ClientFor(proxy)

	fetchClient := &http.Client{
		Transport: baseClient.Transport,
		Timeout:   30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= MaxRedirects {
				return fmt.Errorf("too many redirects (%d)", MaxRedirects)
			}
			return ValidateURL(req.URL.String())
		},
	}

	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, targetURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Proxy0x-MCP/1.0")
	}

	start := time.Now()
	resp, err := fetchClient.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		f.pool.AddLog("fetch-failed", proxy.ID, proxy.Address(), err.Error())
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	limited := io.LimitReader(resp.Body, MaxResponseBody+1)
	respBody, _ := io.ReadAll(limited)
	truncated := len(respBody) > MaxResponseBody
	if truncated {
		respBody = respBody[:MaxResponseBody]
	}

	respHeaders := make(map[string]string)
	for k := range resp.Header {
		respHeaders[k] = resp.Header.Get(k)
	}

	result := &FetchResult{
		StatusCode: resp.StatusCode,
		Headers:    respHeaders,
		Body:       string(respBody),
		ProxyID:    proxy.ID,
		ProxyIP:    proxy.CurrentIP,
		LatencyMs:  latency,
		Truncated:  truncated,
	}

	f.pool.AddLog("fetched", proxy.ID, proxy.Address(),
		fmt.Sprintf("%s %s -> %d (%dms)", method, targetURL, resp.StatusCode, latency))

	return result, nil
}
