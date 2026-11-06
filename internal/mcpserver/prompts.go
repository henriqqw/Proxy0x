package mcpserver

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerPrompts() {
	s.mcp.AddPrompt(
		mcp.NewPrompt("scrape_workflow",
			mcp.WithPromptDescription("Guided workflow for web scraping through the proxy pool"),
			mcp.WithArgument("target_url", mcp.ArgumentDescription("URL to scrape"), mcp.RequiredArgument()),
			mcp.WithArgument("data_type", mcp.ArgumentDescription("Type of data to extract (prices, text, links)")),
		),
		s.handleScrapeWorkflow,
	)

	s.mcp.AddPrompt(
		mcp.NewPrompt("geo_check",
			mcp.WithPromptDescription("Verify URL accessibility from a geographic location"),
			mcp.WithArgument("url", mcp.ArgumentDescription("URL to check"), mcp.RequiredArgument()),
			mcp.WithArgument("country", mcp.ArgumentDescription("Target country code (US, BR, DE)")),
		),
		s.handleGeoCheck,
	)
}

func (s *MCPServer) handleScrapeWorkflow(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	targetURL := req.Params.Arguments["target_url"]
	dataType := req.Params.Arguments["data_type"]
	if dataType == "" {
		dataType = "general content"
	}

	prompt := fmt.Sprintf(`You are using Proxy0x's rotating proxy pool to scrape web content.

Target: %s
Extract: %s

Workflow:
1. proxy_status — check pool has active proxies
2. proxy_connect — activate one if needed
3. proxy_fetch url=%s — fetch the page
4. On 403/429: proxy_rotate, then retry
5. Parse the response body for the requested data
6. For pagination: rotate between requests to avoid rate-limits
7. proxy_health — verify proxies stay healthy

The proxy_fetch tool auto-selects the best proxy by score.`, targetURL, dataType, targetURL)

	return mcp.NewGetPromptResult(
		"Web scraping with Proxy0x rotating proxies",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(prompt)),
		},
	), nil
}

func (s *MCPServer) handleGeoCheck(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	targetURL := req.Params.Arguments["url"]
	country := req.Params.Arguments["country"]
	if country == "" {
		country = "any available"
	}

	prompt := fmt.Sprintf(`Verify geographic access to a URL using Proxy0x.

URL: %s
Country: %s

Workflow:
1. proxy_list — find a proxy in the target country
2. proxy_connect — activate it if needed
3. proxy_fetch url=%s proxy_id=<id> — fetch through that proxy
4. Check response:
   - 200 = accessible
   - 403 = geo-blocked
   - 301/302 = geo-redirect (check Location header)
   - Look for country-specific content or block messages
5. Report: accessible/blocked/redirected, proxy location, response details`, targetURL, country, targetURL)

	return mcp.NewGetPromptResult(
		"Geographic access check with Proxy0x",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(prompt)),
		},
	), nil
}
