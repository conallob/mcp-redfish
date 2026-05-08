package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/conallob/mcp-redfish/internal/redfish"
	"github.com/mark3labs/mcp-go/mcp"
)

func handleGetServiceRoot(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sr, err := c.GetServiceRoot()
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get service root: %v", err), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Redfish Service Root\n")
		fmt.Fprintf(&sb, "  Name:            %s\n", sr.Name)
		fmt.Fprintf(&sb, "  Redfish Version: %s\n", sr.RedfishVersion)
		fmt.Fprintf(&sb, "  UUID:            %s\n", sr.UUID)
		if sr.Product != "" {
			fmt.Fprintf(&sb, "  Product:         %s\n", sr.Product)
		}
		if sr.Vendor != "" {
			fmt.Fprintf(&sb, "  Vendor:          %s\n", sr.Vendor)
		}
		fmt.Fprintf(&sb, "  Systems:  %s\n", sr.Systems.ODataID)
		fmt.Fprintf(&sb, "  Chassis:  %s\n", sr.Chassis.ODataID)
		fmt.Fprintf(&sb, "  Managers: %s\n", sr.Managers.ODataID)
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleListSystems(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		systems, err := c.ListSystems()
		if err != nil {
			return mcp.NewToolResultErrorf("failed to list systems: %v", err), nil
		}
		if len(systems) == 0 {
			return mcp.NewToolResultText("No systems found."), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Systems (%d):\n", len(systems))
		for _, s := range systems {
			fmt.Fprintf(&sb, "  - ID: %s | Name: %s | Model: %s | Power: %s\n",
				s.ID, s.Name, s.Model, s.PowerState)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetSystem(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		systemID := req.GetString("system_id", "")
		sys, err := c.GetSystem(systemID)
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get system: %v", err), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "System: %s\n", sys.Name)
		fmt.Fprintf(&sb, "  ID:           %s\n", sys.ID)
		fmt.Fprintf(&sb, "  Manufacturer: %s\n", sys.Manufacturer)
		fmt.Fprintf(&sb, "  Model:        %s\n", sys.Model)
		if sys.SKU != "" {
			fmt.Fprintf(&sb, "  SKU:          %s\n", sys.SKU)
		}
		fmt.Fprintf(&sb, "  Serial:       %s\n", sys.SerialNumber)
		if sys.HostName != "" {
			fmt.Fprintf(&sb, "  Hostname:     %s\n", sys.HostName)
		}
		fmt.Fprintf(&sb, "  Power State:  %s\n", sys.PowerState)
		fmt.Fprintf(&sb, "  BIOS Version: %s\n", sys.BIOSVersion)
		fmt.Fprintf(&sb, "  Status:       %s / %s\n", sys.Status.Health, sys.Status.State)
		fmt.Fprintf(&sb, "  Processors:   %d x %s\n",
			sys.ProcessorSummary.Count, sys.ProcessorSummary.Model)
		fmt.Fprintf(&sb, "  Memory:       %.0f GiB\n", sys.MemorySummary.TotalSystemMemoryGiB)
		if sys.IndicatorLED != "" {
			fmt.Fprintf(&sb, "  Indicator LED: %s\n", sys.IndicatorLED)
		}
		if sys.Boot.BootSourceOverrideEnabled != "" && sys.Boot.BootSourceOverrideEnabled != "Disabled" {
			fmt.Fprintf(&sb, "  Boot Override: %s -> %s\n",
				sys.Boot.BootSourceOverrideEnabled, sys.Boot.BootSourceOverrideTarget)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetProcessors(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		systemID := req.GetString("system_id", "")
		procs, err := c.GetProcessors(systemID)
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get processors: %v", err), nil
		}
		if len(procs) == 0 {
			return mcp.NewToolResultText("No processors found."), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Processors (%d):\n", len(procs))
		for _, p := range procs {
			fmt.Fprintf(&sb, "\n  [%s] %s\n", p.ID, p.Name)
			fmt.Fprintf(&sb, "    Type:          %s (%s)\n", p.ProcessorType, p.ProcessorArchitecture)
			fmt.Fprintf(&sb, "    Manufacturer:  %s\n", p.Manufacturer)
			fmt.Fprintf(&sb, "    Model:         %s\n", p.Model)
			fmt.Fprintf(&sb, "    Socket:        %s\n", p.Socket)
			fmt.Fprintf(&sb, "    Cores/Threads: %d / %d\n", p.TotalCores, p.TotalThreads)
			if p.MaxSpeedMHz > 0 {
				fmt.Fprintf(&sb, "    Max Speed:     %d MHz\n", p.MaxSpeedMHz)
			}
			fmt.Fprintf(&sb, "    Status:        %s / %s\n", p.Status.Health, p.Status.State)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetMemory(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		systemID := req.GetString("system_id", "")
		modules, err := c.GetMemory(systemID)
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get memory: %v", err), nil
		}
		if len(modules) == 0 {
			return mcp.NewToolResultText("No memory modules found."), nil
		}
		var sb strings.Builder
		totalMiB := 0
		for _, m := range modules {
			totalMiB += m.CapacityMiB
		}
		fmt.Fprintf(&sb, "Memory Modules (%d, total %d GiB):\n", len(modules), totalMiB/1024)
		for _, m := range modules {
			if m.CapacityMiB == 0 {
				continue // Skip empty slots
			}
			fmt.Fprintf(&sb, "\n  [%s] %s\n", m.ID, m.Name)
			fmt.Fprintf(&sb, "    Type:         %s %s\n", m.MemoryType, m.MemoryDeviceType)
			fmt.Fprintf(&sb, "    Capacity:     %d MiB\n", m.CapacityMiB)
			if m.OperatingSpeedMhz > 0 {
				fmt.Fprintf(&sb, "    Speed:        %d MHz\n", m.OperatingSpeedMhz)
			}
			if m.Manufacturer != "" {
				fmt.Fprintf(&sb, "    Manufacturer: %s\n", m.Manufacturer)
			}
			if m.DeviceLocator != "" {
				fmt.Fprintf(&sb, "    Slot:         %s\n", m.DeviceLocator)
			}
			fmt.Fprintf(&sb, "    Status:       %s / %s\n", m.Status.Health, m.Status.State)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetStorage(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		systemID := req.GetString("system_id", "")
		controllers, err := c.GetStorage(systemID)
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get storage: %v", err), nil
		}
		if len(controllers) == 0 {
			return mcp.NewToolResultText("No storage controllers found."), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Storage Controllers (%d):\n", len(controllers))
		for _, sc := range controllers {
			fmt.Fprintf(&sb, "\n  [%s] %s\n", sc.ID, sc.Name)
			for _, ctrl := range sc.StorageControllers {
				fmt.Fprintf(&sb, "    Controller:  %s %s\n", ctrl.Manufacturer, ctrl.Model)
				if ctrl.FirmwareVersion != "" {
					fmt.Fprintf(&sb, "    Firmware:    %s\n", ctrl.FirmwareVersion)
				}
				if len(ctrl.SupportedRAIDTypes) > 0 {
					fmt.Fprintf(&sb, "    RAID Types:  %s\n", strings.Join(ctrl.SupportedRAIDTypes, ", "))
				}
			}
			fmt.Fprintf(&sb, "    Drives (%d):\n", len(sc.Drives))
			for _, d := range sc.Drives {
				drive, err := c.GetDrive(d.ODataID)
				if err != nil {
					fmt.Fprintf(&sb, "      - %s (error: %v)\n", d.ODataID, err)
					continue
				}
				sizeGB := drive.CapacityBytes / 1_000_000_000
				fmt.Fprintf(&sb, "      - [%s] %s %s | %s %s | %d GB | %s/%s\n",
					drive.ID, drive.Manufacturer, drive.Model,
					drive.MediaType, drive.Protocol,
					sizeGB, drive.Status.Health, drive.Status.State)
			}
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetNetworkInterfaces(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		systemID := req.GetString("system_id", "")
		nics, err := c.GetNetworkInterfaces(systemID)
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get network interfaces: %v", err), nil
		}
		if len(nics) == 0 {
			return mcp.NewToolResultText("No network interfaces found."), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Network Interfaces (%d):\n", len(nics))
		for _, n := range nics {
			fmt.Fprintf(&sb, "  - [%s] %s | Status: %s/%s\n",
				n.ID, n.Name, n.Status.Health, n.Status.State)
			if n.Description != "" {
				fmt.Fprintf(&sb, "      Description: %s\n", n.Description)
			}
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetThermal(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		thermal, err := c.GetThermal()
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get thermal data: %v", err), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Thermal Readings\n")
		fmt.Fprintf(&sb, "\nTemperatures (%d):\n", len(thermal.Temperatures))
		for _, t := range thermal.Temperatures {
			fmt.Fprintf(&sb, "  %-40s %6.1f°C  [%s/%s]",
				t.Name, t.ReadingCelsius, t.Status.Health, t.Status.State)
			if t.UpperThresholdCritical > 0 {
				fmt.Fprintf(&sb, "  warn@%.0f°C  crit@%.0f°C",
					t.UpperThresholdNonCritical, t.UpperThresholdCritical)
			}
			fmt.Fprintln(&sb)
		}
		fmt.Fprintf(&sb, "\nFans (%d):\n", len(thermal.Fans))
		for _, f := range thermal.Fans {
			units := f.ReadingUnits
			if units == "" {
				units = "RPM"
			}
			fmt.Fprintf(&sb, "  %-40s %6d %s  [%s/%s]\n",
				f.Name, f.Reading, units, f.Status.Health, f.Status.State)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetPower(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		power, err := c.GetPower()
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get power data: %v", err), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Power\n")
		if len(power.PowerControl) > 0 {
			fmt.Fprintf(&sb, "\nPower Consumption:\n")
			for _, pc := range power.PowerControl {
				fmt.Fprintf(&sb, "  %s\n", pc.Name)
				fmt.Fprintf(&sb, "    Consumed: %.0f W\n", pc.PowerConsumedWatts)
				if pc.PowerCapacityWatts > 0 {
					fmt.Fprintf(&sb, "    Capacity: %.0f W\n", pc.PowerCapacityWatts)
				}
				if pc.PowerMetrics.AverageConsumedWatts > 0 {
					fmt.Fprintf(&sb, "    Avg/Min/Max: %.0f / %.0f / %.0f W\n",
						pc.PowerMetrics.AverageConsumedWatts,
						pc.PowerMetrics.MinConsumedWatts,
						pc.PowerMetrics.MaxConsumedWatts)
				}
			}
		}
		if len(power.PowerSupplies) > 0 {
			fmt.Fprintf(&sb, "\nPower Supplies (%d):\n", len(power.PowerSupplies))
			for _, ps := range power.PowerSupplies {
				fmt.Fprintf(&sb, "  %s (%s %s)\n", ps.Name, ps.Manufacturer, ps.Model)
				fmt.Fprintf(&sb, "    Output: %.0f W  Capacity: %.0f W  Input: %.0f V  [%s/%s]\n",
					ps.LastPowerOutputWatts, ps.PowerCapacityWatts,
					ps.LineInputVoltage, ps.Status.Health, ps.Status.State)
			}
		}
		if len(power.Voltages) > 0 {
			fmt.Fprintf(&sb, "\nVoltages (%d):\n", len(power.Voltages))
			for _, v := range power.Voltages {
				fmt.Fprintf(&sb, "  %-40s %6.3f V  [%s/%s]\n",
					v.Name, v.ReadingVolts, v.Status.Health, v.Status.State)
			}
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetEventLog(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		systemID := req.GetString("system_id", "")
		limit := req.GetInt("limit", 50)
		entries, err := c.GetEventLog(systemID, limit)
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get event log: %v", err), nil
		}
		if len(entries) == 0 {
			return mcp.NewToolResultText("Event log is empty."), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Event Log (%d entries):\n\n", len(entries))
		for _, e := range entries {
			fmt.Fprintf(&sb, "[%s] %-10s %-12s %s\n",
				e.Created, e.Severity, e.EntryType, e.Message)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleListManagers(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		managers, err := c.ListManagers()
		if err != nil {
			return mcp.NewToolResultErrorf("failed to list managers: %v", err), nil
		}
		if len(managers) == 0 {
			return mcp.NewToolResultText("No managers found."), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Managers (%d):\n", len(managers))
		for _, m := range managers {
			fmt.Fprintf(&sb, "  - [%s] %s | Type: %s | FW: %s | Status: %s/%s\n",
				m.ID, m.Name, m.ManagerType, m.FirmwareVersion,
				m.Status.Health, m.Status.State)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetManager(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		managerID := req.GetString("manager_id", "")
		mgr, err := c.GetManager(managerID)
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get manager: %v", err), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Manager: %s\n", mgr.Name)
		fmt.Fprintf(&sb, "  ID:              %s\n", mgr.ID)
		fmt.Fprintf(&sb, "  Type:            %s\n", mgr.ManagerType)
		if mgr.Manufacturer != "" {
			fmt.Fprintf(&sb, "  Manufacturer:    %s\n", mgr.Manufacturer)
		}
		if mgr.Model != "" {
			fmt.Fprintf(&sb, "  Model:           %s\n", mgr.Model)
		}
		fmt.Fprintf(&sb, "  Firmware:        %s\n", mgr.FirmwareVersion)
		fmt.Fprintf(&sb, "  UUID:            %s\n", mgr.UUID)
		fmt.Fprintf(&sb, "  Power State:     %s\n", mgr.PowerState)
		fmt.Fprintf(&sb, "  Status:          %s / %s\n", mgr.Status.Health, mgr.Status.State)
		if mgr.DateTime != "" {
			fmt.Fprintf(&sb, "  Date/Time:       %s %s\n", mgr.DateTime, mgr.DateTimeLocalOffset)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

func handleGetBios(c *redfish.Client) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		systemID := req.GetString("system_id", "")
		bios, err := c.GetBios(systemID)
		if err != nil {
			return mcp.NewToolResultErrorf("failed to get BIOS: %v", err), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "BIOS: %s\n", bios.Name)
		if len(bios.Attributes) == 0 {
			fmt.Fprintf(&sb, "  No attributes available.\n")
			return mcp.NewToolResultText(sb.String()), nil
		}
		fmt.Fprintf(&sb, "  Attributes (%d):\n", len(bios.Attributes))
		for k, v := range bios.Attributes {
			fmt.Fprintf(&sb, "    %-50s = %v\n", k, v)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}
