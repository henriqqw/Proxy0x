package mcpserver

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerResources() {
	s.mcp.AddResource(
		mcp.NewResource("proxy://status", "Pool Status",
			mcp.WithResourceDescription("Current proxy pool statistics"),
			mcp.WithMIMEType("application/json"),
		),
		s.handleResourceStatus,
	)

	s.mcp.AddResource(
		mcp.NewResource("proxy://pool", "Proxy Pool",
			mcp.WithResourceDescription("All proxies with current state (no credentials)"),
			mcp.WithMIMEType("application/json"),
		),
		s.handleResourcePool,
	)

	s.mcp.AddResource(
		mcp.NewResource("proxy://logs", "Activity Logs",
			mcp.WithResourceDescription("Recent proxy activity (last 200 entries)"),
			mcp.WithMIMEType("application/json"),
		),
		s.handleResourceLogs,
	)
}

func (s *MCPServer) handleResourceStatus(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	stats := s.pool.Stats()
	data, _ := json.MarshalIndent(stats, "", "  ")
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "proxy://status",
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}

func (s *MCPServer) handleResourcePool(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	proxies := s.pool.List()

	type safeEntry struct {
		ID        string `json:"id"`
		Host      string `json:"host"`
		Port      string `json:"port"`
		Protocol  string `json:"protocol"`
		Group     string `json:"group,omitempty"`
		Status    string `json:"status"`
		CurrentIP string `json:"currentIp,omitempty"`
		Latency   int64  `json:"latencyMs"`
		Active    bool   `json:"active"`
	}

	safe := make([]safeEntry, len(proxies))
	for i, p := range proxies {
		safe[i] = safeEntry{
			ID:        p.ID,
			Host:      p.Host,
			Port:      p.Port,
			Protocol:  p.Protocol,
			Group:     p.Group,
			Status:    string(p.Status),
			CurrentIP: p.CurrentIP,
			Latency:   p.Latency,
			Active:    p.Active,
		}
	}

	data, _ := json.MarshalIndent(safe, "", "  ")
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "proxy://pool",
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}

func (s *MCPServer) handleResourceLogs(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	logs := s.pool.Logs()
	data, _ := json.MarshalIndent(logs, "", "  ")
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "proxy://logs",
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}
