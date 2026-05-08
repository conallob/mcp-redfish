package tools_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/conallob/mcp-redfish/internal/redfish"
	"github.com/conallob/mcp-redfish/internal/tools"
)

// fixtures is a shared set of mock Redfish responses for tool tests.
var fixtures = map[string]interface{}{
	"/redfish/v1/": redfish.ServiceRoot{
		ODataID:        "/redfish/v1/",
		ID:             "RootService",
		Name:           "Root Service",
		RedfishVersion: "1.16.0",
		Systems:        redfish.ODataLink{ODataID: "/redfish/v1/Systems"},
		Chassis:        redfish.ODataLink{ODataID: "/redfish/v1/Chassis"},
		Managers:       redfish.ODataLink{ODataID: "/redfish/v1/Managers"},
	},
	"/redfish/v1/Systems": redfish.Collection{
		MembersCount: 1,
		Members:      []redfish.ODataLink{{ODataID: "/redfish/v1/Systems/1"}},
	},
	"/redfish/v1/Systems/1": redfish.ComputerSystem{
		ODataID:      "/redfish/v1/Systems/1",
		ID:           "1",
		Name:         "System",
		Manufacturer: "Acme",
		Model:        "TestBox",
		SerialNumber: "SN-001",
		PowerState:   "On",
		BIOSVersion:  "1.0",
		Status:       redfish.Status{State: "Enabled", Health: "OK"},
		ProcessorSummary: redfish.ProcessorSummary{
			Count: 1,
			Model: "Test CPU",
		},
		MemorySummary:     redfish.MemorySummary{TotalSystemMemoryGiB: 64},
		Processors:        redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/Processors"},
		Memory:            redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/Memory"},
		Storage:           redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/Storage"},
		NetworkInterfaces: redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/NetworkInterfaces"},
		LogServices:       redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/LogServices"},
		Bios:              redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/Bios"},
		Actions: redfish.SystemActions{
			Reset: redfish.ActionTarget{
				Target: "/redfish/v1/Systems/1/Actions/ComputerSystem.Reset",
			},
		},
	},
	"/redfish/v1/Systems/1/Processors": redfish.Collection{
		MembersCount: 1,
		Members:      []redfish.ODataLink{{ODataID: "/redfish/v1/Systems/1/Processors/CPU1"}},
	},
	"/redfish/v1/Systems/1/Processors/CPU1": redfish.Processor{
		ID:           "CPU1",
		Name:         "CPU1",
		Manufacturer: "Intel",
		Model:        "Test CPU",
		TotalCores:   8,
		TotalThreads: 16,
		Socket:       "CPU1",
		Status:       redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/Memory": redfish.Collection{
		MembersCount: 1,
		Members:      []redfish.ODataLink{{ODataID: "/redfish/v1/Systems/1/Memory/DIMM1"}},
	},
	"/redfish/v1/Systems/1/Memory/DIMM1": redfish.MemoryModule{
		ID:               "DIMM1",
		Name:             "DIMM1",
		MemoryType:       "DRAM",
		MemoryDeviceType: "DDR4",
		CapacityMiB:      65536,
		Status:           redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/Storage": redfish.Collection{
		MembersCount: 1,
		Members:      []redfish.ODataLink{{ODataID: "/redfish/v1/Systems/1/Storage/CTRL1"}},
	},
	"/redfish/v1/Systems/1/Storage/CTRL1": redfish.StorageController{
		ID:   "CTRL1",
		Name: "Test Controller",
		StorageControllers: []redfish.StorageControllerSummary{
			{Name: "Test RAID", Manufacturer: "Test", Model: "CTRL1"},
		},
		Drives: []redfish.ODataLink{
			{ODataID: "/redfish/v1/Systems/1/Storage/CTRL1/Drives/Disk0"},
		},
		Status: redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/Storage/CTRL1/Drives/Disk0": redfish.Drive{
		ID:            "Disk0",
		Name:          "Disk 0",
		Manufacturer:  "WD",
		Model:         "WD4000",
		MediaType:     "SSD",
		Protocol:      "SATA",
		CapacityBytes: 4000000000000,
		Status:        redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/NetworkInterfaces": redfish.Collection{
		MembersCount: 1,
		Members:      []redfish.ODataLink{{ODataID: "/redfish/v1/Systems/1/NetworkInterfaces/NIC1"}},
	},
	"/redfish/v1/Systems/1/NetworkInterfaces/NIC1": redfish.NetworkInterface{
		ID:     "NIC1",
		Name:   "NIC1",
		Status: redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/LogServices": redfish.Collection{
		MembersCount: 1,
		Members:      []redfish.ODataLink{{ODataID: "/redfish/v1/Systems/1/LogServices/System"}},
	},
	"/redfish/v1/Systems/1/LogServices/System/Entries": map[string]interface{}{
		"Members@odata.count": 1,
		"Members": []redfish.LogEntry{
			{ID: "1", Created: "2026-01-01T00:00:00Z", Severity: "OK", Message: "Test log entry", EntryType: "Event"},
		},
	},
	"/redfish/v1/Systems/1/Bios": redfish.Bios{
		ID:   "Bios",
		Name: "BIOS",
		Attributes: map[string]interface{}{
			"BootMode": "Uefi",
		},
	},
	"/redfish/v1/Systems/1/Bios/Settings": redfish.Bios{}, // PATCH target
	"/redfish/v1/Chassis": redfish.Collection{
		MembersCount: 1,
		Members:      []redfish.ODataLink{{ODataID: "/redfish/v1/Chassis/1"}},
	},
	"/redfish/v1/Chassis/1/Thermal": redfish.Thermal{
		ID:   "Thermal",
		Name: "Thermal",
		Temperatures: []redfish.Temperature{
			{Name: "Inlet", ReadingCelsius: 22, Status: redfish.Status{State: "Enabled", Health: "OK"}},
		},
		Fans: []redfish.Fan{
			{Name: "Fan1", Reading: 5000, ReadingUnits: "RPM", Status: redfish.Status{State: "Enabled", Health: "OK"}},
		},
	},
	"/redfish/v1/Chassis/1/Power": redfish.Power{
		ID:   "Power",
		Name: "Power",
		PowerControl: []redfish.PowerControl{
			{Name: "Main", PowerConsumedWatts: 200, Status: redfish.Status{State: "Enabled", Health: "OK"}},
		},
	},
	"/redfish/v1/Managers": redfish.Collection{
		MembersCount: 1,
		Members:      []redfish.ODataLink{{ODataID: "/redfish/v1/Managers/BMC"}},
	},
	"/redfish/v1/Managers/BMC": redfish.Manager{
		ID:              "BMC",
		Name:            "BMC",
		ManagerType:     "BMC",
		FirmwareVersion: "1.0.0",
		Status:          redfish.Status{State: "Enabled", Health: "OK"},
	},
}

// newToolTestServer spins up a mock Redfish server and returns a configured client.
func newToolTestServer(t *testing.T) *redfish.Client {
	t.Helper()
	mux := http.NewServeMux()
	for path, fixture := range fixtures {
		p, f := path, fixture
		mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(f)
			case http.MethodPost, http.MethodPatch:
				w.WriteHeader(http.StatusNoContent)
			}
		})
	}
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return redfish.NewClient(redfish.Config{Host: srv.URL})
}

// callTool invokes a registered tool on a fresh MCPServer with the given args.
func callTool(t *testing.T, toolName string, args map[string]interface{}, readOnly bool) *mcp.CallToolResult {
	t.Helper()
	client := newToolTestServer(t)
	s := server.NewMCPServer("test", "0.0.0")
	tools.Register(s, client, readOnly)

	st := s.GetTool(toolName)
	if st == nil {
		t.Fatalf("tool %q not registered", toolName)
	}

	req := mcp.CallToolRequest{}
	req.Params.Name = toolName
	req.Params.Arguments = args

	result, err := st.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("tool %q handler returned unexpected error: %v", toolName, err)
	}
	return result
}

// --- Read-only tool tests ---

func TestTool_GetServiceRoot(t *testing.T) {
	result := callTool(t, "redfish_get_service_root", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "Redfish Service Root")
	assertContains(t, result, "1.16.0")
}

func TestTool_ListSystems(t *testing.T) {
	result := callTool(t, "redfish_list_systems", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "TestBox")
}

func TestTool_GetSystem(t *testing.T) {
	result := callTool(t, "redfish_get_system", map[string]interface{}{"system_id": "1"}, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "Acme")
	assertContains(t, result, "SN-001")
}

func TestTool_GetSystem_NoID(t *testing.T) {
	result := callTool(t, "redfish_get_system", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "TestBox")
}

func TestTool_GetProcessors(t *testing.T) {
	result := callTool(t, "redfish_get_processors", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "CPU1")
	assertContains(t, result, "16") // threads
}

func TestTool_GetMemory(t *testing.T) {
	result := callTool(t, "redfish_get_memory", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "DDR4")
}

func TestTool_GetStorage(t *testing.T) {
	result := callTool(t, "redfish_get_storage", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "SSD")
}

func TestTool_GetNetworkInterfaces(t *testing.T) {
	result := callTool(t, "redfish_get_network_interfaces", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "NIC1")
}

func TestTool_GetThermal(t *testing.T) {
	result := callTool(t, "redfish_get_thermal", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "Inlet")
	assertContains(t, result, "Fan1")
}

func TestTool_GetPower(t *testing.T) {
	result := callTool(t, "redfish_get_power", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "200")
}

func TestTool_GetEventLog(t *testing.T) {
	result := callTool(t, "redfish_get_event_log", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "Test log entry")
}

func TestTool_ListManagers(t *testing.T) {
	result := callTool(t, "redfish_list_managers", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "BMC")
}

func TestTool_GetManager(t *testing.T) {
	result := callTool(t, "redfish_get_manager", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "1.0.0")
}

func TestTool_GetBios(t *testing.T) {
	result := callTool(t, "redfish_get_bios", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "Uefi")
}

// --- Read-write tool tests ---

func TestTool_ResetSystem(t *testing.T) {
	result := callTool(t, "redfish_reset_system",
		map[string]interface{}{"reset_type": "GracefulShutdown"}, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "GracefulShutdown")
}

func TestTool_SetIndicatorLED(t *testing.T) {
	result := callTool(t, "redfish_set_indicator_led",
		map[string]interface{}{"state": "Blinking"}, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "Blinking")
}

func TestTool_ClearEventLog(t *testing.T) {
	result := callTool(t, "redfish_clear_event_log", nil, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "cleared")
}

func TestTool_SetBiosAttribute(t *testing.T) {
	result := callTool(t, "redfish_set_bios_attribute",
		map[string]interface{}{"attribute": "BootMode", "value": "Bios"}, false)
	assertToolSuccess(t, result)
	assertContains(t, result, "BootMode")
}

// --- Read-only mode tests ---

func TestReadOnlyMode_ExcludesWriteTools(t *testing.T) {
	writeTools := []string{
		"redfish_reset_system",
		"redfish_set_indicator_led",
		"redfish_clear_event_log",
		"redfish_set_bios_attribute",
	}
	for _, name := range writeTools {
		t.Run(name, func(t *testing.T) {
			client := newToolTestServer(t)
			s := server.NewMCPServer("test", "0.0.0")
			tools.Register(s, client, true /* readOnly */)

			if st := s.GetTool(name); st != nil {
				t.Errorf("write tool %q should not be registered in read-only mode", name)
			}
		})
	}
}

func TestReadOnlyMode_AllowsReadTools(t *testing.T) {
	result := callTool(t, "redfish_get_system", nil, true /* readOnly */)
	assertToolSuccess(t, result)
}

// --- Helpers ---

func assertToolSuccess(t *testing.T, result *mcp.CallToolResult) {
	t.Helper()
	if result == nil {
		t.Fatal("result is nil")
	}
	if result.IsError {
		t.Fatalf("tool returned error: %v", toolText(result))
	}
}

func assertContains(t *testing.T, result *mcp.CallToolResult, substr string) {
	t.Helper()
	text := toolText(result)
	if text == "" {
		t.Fatalf("result has no text content")
	}
	for i := range len(text) - len(substr) + 1 {
		if text[i:i+len(substr)] == substr {
			return
		}
	}
	t.Errorf("result text does not contain %q\nfull text:\n%s", substr, text)
}

func toolText(result *mcp.CallToolResult) string {
	if result == nil {
		return ""
	}
	for _, c := range result.Content {
		if tc, ok := c.(mcp.TextContent); ok {
			return tc.Text
		}
	}
	return ""
}
