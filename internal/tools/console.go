package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/conallob/mcp-redfish/internal/redfish"
)

func handleGetConsoleStatus(sol *redfish.SOLClient) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if sol == nil {
			return mcp.NewToolResultText("Serial console not configured. Set --sol-url or REDFISH_SOL_URL."), nil
		}
		state, lastErr := sol.State()
		var sb strings.Builder
		fmt.Fprintf(&sb, "Serial Console Status\n")
		fmt.Fprintf(&sb, "  State:          %s\n", state)
		if lastErr != "" {
			fmt.Fprintf(&sb, "  Last error:     %s\n", lastErr)
		}
		fmt.Fprintf(&sb, "  Buffered lines: %d\n", len(sol.RecentLines(0)))
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetConsoleOutput(sol *redfish.SOLClient) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if sol == nil {
			return mcp.NewToolResultErrorf("Serial console not configured. Set --sol-url or REDFISH_SOL_URL."), nil
		}
		limit := req.GetInt("limit", 100)
		if limit <= 0 {
			limit = 100
		}
		state, _ := sol.State()
		lines := sol.RecentLines(limit)
		if len(lines) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf(
				"No console output buffered (connection state: %s).", state)), nil
		}
		return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
	}
}
