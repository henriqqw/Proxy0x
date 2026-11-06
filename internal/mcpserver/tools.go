package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"proxy_ui/pkg/proxycore"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerTools() {
	s.mcp.AddTool(
		mcp.NewTool("proxy_fetch",
			mcp.WithDescription("Fetch a URL through the rotating proxy pool. Selects the best proxy by score, validates against SSRF, and returns the response with status, headers and body."),
			mcp.WithString("url", mcp.Required(), mcp.Description("Target URL (http/https only)")),
			mcp.WithString("method", mcp.Description("HTTP method"), mcp.Enum("GET", "POST", "PUT", "DELETE", "HEAD", "PATCH")),
			mcp.WithString("body", mcp.Description("Request body for POST/PUT/PATCH")),
			mcp.WithString("proxy_id", mcp.Description("Specific proxy ID to use instead of auto-select")),
		),
		s.handleProxyFetch,
	)

	s.mcp.AddTool(
		mcp.NewTool("proxy_list",
			mcp.WithDescription("List all proxies in the pool with status, IP, latency, location and score. Credentials are never exposed."),
		),
		s.handleProxyList,
	)

	s.mcp.AddTool(
		mcp.NewTool("proxy_status",
			mcp.WithDescription("Pool overview: total, active, healthy, errored counts and average latency."),
		),
		s.handleProxyStatus,
	)

	s.mcp.AddTool(
		mcp.NewTool("proxy_connect",
			mcp.WithDescription("Connect a proxy by ID. Verifies connectivity and resolves the exit IP."),
			mcp.WithString("id", mcp.Required(), mcp.Description("Proxy ID")),
		),
		s.handleProxyConnect,
	)

	s.mcp.AddTool(
		mcp.NewTool("proxy_disconnect",
			mcp.WithDescription("Disconnect an active proxy by ID."),
			mcp.WithString("id", mcp.Required(), mcp.Description("Proxy ID")),
		),
		s.handleProxyDisconnect,
	)

	s.mcp.AddTool(
		mcp.NewTool("proxy_rotate",
			mcp.WithDescription("Rotate a proxy's session to get a new exit IP."),
			mcp.WithString("id", mcp.Required(), mcp.Description("Proxy ID")),
		),
		s.handleProxyRotate,
	)

	s.mcp.AddTool(
		mcp.NewTool("proxy_add",
			mcp.WithDescription("Add a proxy. Format: host:port or host:port:user:pass. Prefix socks5:// for SOCKS5."),
			mcp.WithString("proxy", mcp.Required(), mcp.Description("Proxy string")),
		),
		s.handleProxyAdd,
	)

	s.mcp.AddTool(
		mcp.NewTool("proxy_remove",
			mcp.WithDescription("Remove a proxy from the pool by ID."),
			mcp.WithString("id", mcp.Required(), mcp.Description("Proxy ID")),
		),
		s.handleProxyRemove,
	)

	s.mcp.AddTool(
		mcp.NewTool("proxy_health",
			mcp.WithDescription("Run health check on all active proxies."),
		),
		s.handleProxyHealth,
	)
}

func (s *MCPServer) handleProxyFetch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	urlStr, err := req.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError("url is required"), nil
	}
	method := req.GetString("method", "GET")
	body := req.GetString("body", "")
	proxyID := req.GetString("proxy_id", "")

	result, err := s.fetcher.Fetch(ctx, urlStr, method, nil, body, proxyID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

func (s *MCPServer) handleProxyList(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	proxies := s.pool.List()

	type safeProxy struct {
		ID        string                `json:"id"`
		Host      string                `json:"host"`
		Port      string                `json:"port"`
		Protocol  string                `json:"protocol"`
		Group     string                `json:"group,omitempty"`
		Status    proxycore.ProxyStatus `json:"status"`
		CurrentIP string                `json:"currentIp,omitempty"`
		Location  *proxycore.GeoInfo    `json:"location,omitempty"`
		Latency   int64                 `json:"latencyMs"`
		Active    bool                  `json:"active"`
		Score     int                   `json:"score"`
		LastCheck string                `json:"lastCheck,omitempty"`
	}

	result := make([]safeProxy, len(proxies))
	for i, p := range proxies {
		result[i] = safeProxy{
			ID:        p.ID,
			Host:      p.Host,
			Port:      p.Port,
			Protocol:  p.Protocol,
			Group:     p.Group,
			Status:    p.Status,
			CurrentIP: p.CurrentIP,
			Location:  p.Location,
			Latency:   p.Latency,
			Active:    p.Active,
			Score:     p.Score(),
			LastCheck: p.LastCheck,
		}
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

func (s *MCPServer) handleProxyStatus(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stats := s.pool.Stats()
	data, _ := json.MarshalIndent(stats, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

func (s *MCPServer) handleProxyConnect(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := req.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}

	proxy, err := s.pool.Connect(id)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, _ := json.MarshalIndent(map[string]interface{}{
		"id":       proxy.ID,
		"status":   proxy.Status,
		"ip":       proxy.CurrentIP,
		"latencyMs": proxy.Latency,
		"location": proxy.Location,
	}, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

func (s *MCPServer) handleProxyDisconnect(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := req.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}

	_, err = s.pool.Disconnect(id)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(`{"id":%q,"status":"disconnected"}`, id)), nil
}

func (s *MCPServer) handleProxyRotate(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := req.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}

	proxy, err := s.pool.Rotate(id)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, _ := json.MarshalIndent(map[string]interface{}{
		"id":       proxy.ID,
		"status":   proxy.Status,
		"ip":       proxy.CurrentIP,
		"latencyMs": proxy.Latency,
		"location": proxy.Location,
	}, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

func (s *MCPServer) handleProxyAdd(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	proxyStr, err := req.RequireString("proxy")
	if err != nil {
		return mcp.NewToolResultError("proxy string is required"), nil
	}

	proxy, err := s.pool.Add(proxyStr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	data, _ := json.MarshalIndent(map[string]interface{}{
		"id":       proxy.ID,
		"host":     proxy.Host,
		"port":     proxy.Port,
		"protocol": proxy.Protocol,
		"status":   "added",
	}, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

func (s *MCPServer) handleProxyRemove(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := req.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("id is required"), nil
	}

	if err := s.pool.Remove(id); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf(`{"id":%q,"status":"removed"}`, id)), nil
}

func (s *MCPServer) handleProxyHealth(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	proxies := s.pool.List()

	type healthResult struct {
		ID      string                `json:"id"`
		Host    string                `json:"host"`
		Status  proxycore.ProxyStatus `json:"status"`
		Latency int64                 `json:"latencyMs"`
		IP      string                `json:"ip,omitempty"`
	}

	var results []healthResult
	for _, p := range proxies {
		if !p.Active {
			continue
		}
		checked, _ := s.pool.CheckHealth(p.ID)
		if checked != nil {
			results = append(results, healthResult{
				ID:      checked.ID,
				Host:    checked.Address(),
				Status:  checked.Status,
				Latency: checked.Latency,
				IP:      checked.CurrentIP,
			})
		}
	}

	if results == nil {
		results = []healthResult{}
	}

	data, _ := json.MarshalIndent(results, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}
