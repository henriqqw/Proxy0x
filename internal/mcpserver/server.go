package mcpserver

import (
	"proxy_ui/pkg/proxycore"

	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	pool    *proxycore.Pool
	fetcher *proxycore.Fetcher
	mcp     *server.MCPServer
}

func New(pool *proxycore.Pool) *MCPServer {
	s := &MCPServer{
		pool:    pool,
		fetcher: proxycore.NewFetcher(pool),
	}

	s.mcp = server.NewMCPServer(
		"Proxy0x",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	s.registerTools()
	s.registerResources()
	s.registerPrompts()

	return s
}

func (s *MCPServer) Server() *server.MCPServer {
	return s.mcp
}
