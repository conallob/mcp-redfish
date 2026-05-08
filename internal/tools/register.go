// Package tools registers Redfish MCP tools with an MCPServer.
// Read-only tools are always registered. Read-write tools are only registered
// when readOnly is false.
package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/conallob/mcp-redfish/internal/redfish"
)

// Register adds all Redfish tools to s.
// sol may be nil; console tools will report "not configured" in that case.
// When readOnly is true, tools that modify state are omitted.
func Register(s *server.MCPServer, c *redfish.Client, sol *redfish.SOLClient, readOnly bool) {
	// --- Read-only tools ---

	s.AddTool(mcp.NewTool("redfish_get_service_root",
		mcp.WithDescription("Get the Redfish service root, showing API version and available resource collections."),
	), handleGetServiceRoot(c))

	s.AddTool(mcp.NewTool("redfish_list_systems",
		mcp.WithDescription("List all computer systems managed by this BMC."),
	), handleListSystems(c))

	s.AddTool(mcp.NewTool("redfish_get_system",
		mcp.WithDescription("Get detailed information about a computer system: manufacturer, model, serial number, power state, BIOS version, processor and memory summary."),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
	), handleGetSystem(c))

	s.AddTool(mcp.NewTool("redfish_get_processors",
		mcp.WithDescription("List all processors (CPUs) installed in a system with socket, core/thread count, and speed."),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
	), handleGetProcessors(c))

	s.AddTool(mcp.NewTool("redfish_get_memory",
		mcp.WithDescription("List all memory modules (DIMMs) installed in a system with type, capacity, and slot location."),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
	), handleGetMemory(c))

	s.AddTool(mcp.NewTool("redfish_get_storage",
		mcp.WithDescription("List storage controllers and drives for a system, including media type, protocol, and capacity."),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
	), handleGetStorage(c))

	s.AddTool(mcp.NewTool("redfish_get_network_interfaces",
		mcp.WithDescription("List network interfaces (NICs) installed in a system."),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
	), handleGetNetworkInterfaces(c))

	s.AddTool(mcp.NewTool("redfish_get_thermal",
		mcp.WithDescription("Get temperature sensor readings and fan speeds from the chassis."),
	), handleGetThermal(c))

	s.AddTool(mcp.NewTool("redfish_get_power",
		mcp.WithDescription("Get power consumption metrics, voltage readings, and power supply status."),
	), handleGetPower(c))

	s.AddTool(mcp.NewTool("redfish_get_event_log",
		mcp.WithDescription("Retrieve entries from the system event log (SEL), ordered most-recent first."),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of log entries to return (default: 50, 0 = all)."),
		),
	), handleGetEventLog(c))

	s.AddTool(mcp.NewTool("redfish_list_managers",
		mcp.WithDescription("List all BMC/management controllers associated with this Redfish endpoint."),
	), handleListManagers(c))

	s.AddTool(mcp.NewTool("redfish_get_manager",
		mcp.WithDescription("Get detailed information about a BMC/management controller: firmware version, model, date/time."),
		mcp.WithString("manager_id",
			mcp.Description("Manager ID or OData path. Omit to use the first manager."),
		),
	), handleGetManager(c))

	s.AddTool(mcp.NewTool("redfish_get_bios",
		mcp.WithDescription("Retrieve the current BIOS/UEFI attribute settings for a system."),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
	), handleGetBios(c))

	// --- Serial console tools (always registered; return "not configured" when sol is nil) ---

	s.AddTool(mcp.NewTool("redfish_get_console_status",
		mcp.WithDescription("Show the current Serial-over-LAN connection state and how many lines are buffered."),
	), handleGetConsoleStatus(sol))

	s.AddTool(mcp.NewTool("redfish_get_console_output",
		mcp.WithDescription("Return recent Serial-over-LAN console output buffered since the server started."),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of lines to return (default: 100)."),
		),
	), handleGetConsoleOutput(sol))

	if readOnly {
		return
	}

	// --- Read-write tools (omitted in read-only mode) ---

	s.AddTool(mcp.NewTool("redfish_reset_system",
		mcp.WithDescription("Reset or change the power state of a system. Common reset types: On, ForceOff, GracefulShutdown, GracefulRestart, ForceRestart, PushPowerButton."),
		mcp.WithString("reset_type",
			mcp.Required(),
			mcp.Description("Reset type to perform."),
			mcp.Enum("On", "ForceOff", "GracefulShutdown", "GracefulRestart", "ForceRestart",
				"Nmi", "ForceOn", "PushPowerButton", "PowerCycle"),
		),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
	), handleResetSystem(c))

	s.AddTool(mcp.NewTool("redfish_set_indicator_led",
		mcp.WithDescription("Control the chassis indicator LED to identify a physical server."),
		mcp.WithString("state",
			mcp.Required(),
			mcp.Description("LED state to set."),
			mcp.Enum("Lit", "Blinking", "Off"),
		),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
	), handleSetIndicatorLED(c))

	s.AddTool(mcp.NewTool("redfish_clear_event_log",
		mcp.WithDescription("Clear (delete) all entries from the system event log. This action is irreversible."),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
	), handleClearEventLog(c))

	s.AddTool(mcp.NewTool("redfish_set_bios_attribute",
		mcp.WithDescription("Set a single BIOS/UEFI attribute. Changes are staged and applied on the next system reboot."),
		mcp.WithString("attribute",
			mcp.Required(),
			mcp.Description("BIOS attribute name (e.g. ProcHyperthreading, BootMode)."),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("New value for the attribute."),
		),
		mcp.WithString("system_id",
			mcp.Description("System ID or OData path. Omit to use the first system."),
		),
	), handleSetBiosAttribute(c))
}
