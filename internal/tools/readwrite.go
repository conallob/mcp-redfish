package tools

import (
	"context"
	"fmt"

	"github.com/conallob/mcp-redfish/internal/redfish"
	"github.com/mark3labs/mcp-go/mcp"
)

func handleResetSystem(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resetType, err := req.RequireString("reset_type")
		if err != nil {
			return mcp.NewToolResultErrorf("reset_type is required: %v", err), nil
		}
		systemID := req.GetString("system_id", "")
		if err := c.ResetSystem(systemID, resetType); err != nil {
			return mcp.NewToolResultErrorf("failed to reset system: %v", err), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("System reset initiated: %s", resetType)), nil
	}
}

func handleSetIndicatorLED(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		state, err := req.RequireString("state")
		if err != nil {
			return mcp.NewToolResultErrorf("state is required: %v", err), nil
		}
		systemID := req.GetString("system_id", "")
		if err := c.SetIndicatorLED(systemID, state); err != nil {
			return mcp.NewToolResultErrorf("failed to set indicator LED: %v", err), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Indicator LED set to: %s", state)), nil
	}
}

func handleClearEventLog(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		systemID := req.GetString("system_id", "")
		if err := c.ClearEventLog(systemID); err != nil {
			return mcp.NewToolResultErrorf("failed to clear event log: %v", err), nil
		}
		return mcp.NewToolResultText("Event log cleared successfully."), nil
	}
}

func handleSetBiosAttribute(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		key, err := req.RequireString("attribute")
		if err != nil {
			return mcp.NewToolResultErrorf("attribute is required: %v", err), nil
		}
		value, err := req.RequireString("value")
		if err != nil {
			return mcp.NewToolResultErrorf("value is required: %v", err), nil
		}
		systemID := req.GetString("system_id", "")
		if err := c.SetBiosAttribute(systemID, key, value); err != nil {
			return mcp.NewToolResultErrorf("failed to set BIOS attribute: %v", err), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("BIOS attribute %q set to %q (pending reboot).", key, value)), nil
	}
}
