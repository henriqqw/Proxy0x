package mcpserver

import (
	"context"
	"encoding/json"
	"testing"

	"proxy_ui/pkg/proxycore"

	"github.com/mark3labs/mcp-go/mcp"
)

func newTestServer(t *testing.T) *MCPServer {
	t.Helper()
	store := proxycore.NewStore(t.TempDir())
	pool := proxycore.NewPool(store)
	return New(pool)
}

func TestServerCreation(t *testing.T) {
	srv := newTestServer(t)
	if srv.Server() == nil {
		t.Fatal("MCP server should not be nil")
	}
	if srv.pool == nil {
		t.Fatal("pool should not be nil")
	}
	if srv.fetcher == nil {
		t.Fatal("fetcher should not be nil")
	}
}

func TestHandleProxyAddIsReadOnly(t *testing.T) {
	srv := newTestServer(t)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{
		"proxy": "1.2.3.4:8080:user:pass",
	}

	result, err := srv.handleProxyAdd(context.Background(), req)
	if err != nil {
		t.Fatalf("handleProxyAdd protocol error: %v", err)
	}
	if !result.IsError {
		t.Fatal("proxy_add must refuse: the MCP server is read-only")
	}
	if len(srv.pool.List()) != 0 {
		t.Fatal("pool must not change")
	}
}

func TestHandleProxyList(t *testing.T) {
	srv := newTestServer(t)

	srv.pool.Add("1.2.3.4:8080")
	srv.pool.Add("5.6.7.8:3128")

	req := mcp.CallToolRequest{}
	result, err := srv.handleProxyList(context.Background(), req)
	if err != nil {
		t.Fatalf("handleProxyList error: %v", err)
	}
	if result.IsError {
		t.Fatal("handleProxyList returned tool error")
	}

	// Verify JSON contains both proxies
	text := extractText(t, result)
	var proxies []json.RawMessage
	if err := json.Unmarshal([]byte(text), &proxies); err != nil {
		t.Fatalf("result is not valid JSON array: %v", err)
	}
	if len(proxies) != 2 {
		t.Errorf("got %d proxies, want 2", len(proxies))
	}
}

func TestHandleProxyStatus(t *testing.T) {
	srv := newTestServer(t)

	srv.pool.Add("1.2.3.4:8080")

	req := mcp.CallToolRequest{}
	result, err := srv.handleProxyStatus(context.Background(), req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result.IsError {
		t.Fatal("returned tool error")
	}

	text := extractText(t, result)
	var stats proxycore.PoolStats
	if err := json.Unmarshal([]byte(text), &stats); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if stats.Total != 1 {
		t.Errorf("total = %d, want 1", stats.Total)
	}
}

func TestHandleProxyRemoveIsReadOnly(t *testing.T) {
	srv := newTestServer(t)

	entry, _ := srv.pool.Add("1.2.3.4:8080")

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{
		"id": entry.ID,
	}

	result, err := srv.handleProxyRemove(context.Background(), req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !result.IsError {
		t.Fatal("proxy_remove must refuse: the MCP server is read-only")
	}
	if len(srv.pool.List()) != 1 {
		t.Error("pool must not change")
	}
}

func TestHandleProxyRemoveNotFound(t *testing.T) {
	srv := newTestServer(t)

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{
		"id": "nonexistent",
	}

	result, err := srv.handleProxyRemove(context.Background(), req)
	if err != nil {
		t.Fatalf("protocol error: %v", err)
	}
	if !result.IsError {
		t.Error("removing non-existent should return tool error")
	}
}

func TestHandleProxyFetchSSRF(t *testing.T) {
	srv := newTestServer(t)

	// Add and "connect" a proxy manually for the test
	entry, _ := srv.pool.Add("1.2.3.4:8080")
	// Manually mark as active+connected
	e, _ := srv.pool.Get(entry.ID)
	e.Active = true
	e.Status = proxycore.StatusConnected

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{
		"url": "http://169.254.169.254/latest/meta-data/",
	}

	result, err := srv.handleProxyFetch(context.Background(), req)
	if err != nil {
		t.Fatalf("protocol error: %v", err)
	}
	if !result.IsError {
		t.Error("SSRF URL should be blocked")
	}
}

func TestHandleProxyHealth(t *testing.T) {
	srv := newTestServer(t)

	// No active proxies
	req := mcp.CallToolRequest{}
	result, err := srv.handleProxyHealth(context.Background(), req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result.IsError {
		t.Fatal("returned tool error")
	}

	text := extractText(t, result)
	var results []json.RawMessage
	if err := json.Unmarshal([]byte(text), &results); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("got %d results, want 0 (no active proxies)", len(results))
	}
}

func TestNoCredentialsInOutput(t *testing.T) {
	srv := newTestServer(t)
	srv.pool.Add("host.com:8080:secretuser:secretpass")

	req := mcp.CallToolRequest{}
	result, err := srv.handleProxyList(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	text := extractText(t, result)
	if contains(text, "secretuser") || contains(text, "secretpass") {
		t.Error("credentials should not appear in proxy_list output")
	}
}

func extractText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if result == nil || len(result.Content) == 0 {
		t.Fatal("result has no content")
	}
	tc, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("content is not TextContent, got %T", result.Content[0])
	}
	return tc.Text
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
