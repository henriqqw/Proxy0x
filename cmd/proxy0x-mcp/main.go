package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"proxy_ui/internal/mcpserver"
	"proxy_ui/pkg/proxycore"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("failed to get config dir: %v", err)
	}
	dataDir := filepath.Join(configDir, "Proxy0x")

	store := proxycore.NewStore(dataDir)
	pool := proxycore.NewPool(store)

	if err := pool.LoadFromStore(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load proxies: %v\n", err)
	}

	mcpSrv := mcpserver.New(pool)

	fmt.Fprintf(os.Stderr, "Proxy0x MCP server starting (stdio)...\n")
	if err := server.ServeStdio(mcpSrv.Server()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
