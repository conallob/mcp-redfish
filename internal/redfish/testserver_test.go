package redfish_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conallob/mcp-redfish/internal/redfish"
)

// redfishFixtures holds the mock server response data.
var redfishFixtures = map[string]interface{}{
	"/redfish/v1/": redfish.ServiceRoot{
		ODataID:        "/redfish/v1/",
		ID:             "RootService",
		Name:           "Root Service",
		RedfishVersion: "1.16.0",
		UUID:           "92384634-2938-2342-8820-489239905423",
		Product:        "Test Server",
		Vendor:         "Test Vendor",
		Systems:        redfish.ODataLink{ODataID: "/redfish/v1/Systems"},
		Chassis:        redfish.ODataLink{ODataID: "/redfish/v1/Chassis"},
		Managers:       redfish.ODataLink{ODataID: "/redfish/v1/Managers"},
	},
	"/redfish/v1/Systems": redfish.Collection{
		ODataID:      "/redfish/v1/Systems",
		Name:         "Computer System Collection",
		MembersCount: 1,
		Members:      []redfish.ODataLink{{ODataID: "/redfish/v1/Systems/1"}},
	},
	"/redfish/v1/Systems/1": redfish.ComputerSystem{
		ODataID:      "/redfish/v1/Systems/1",
		ID:           "1",
		Name:         "System",
		Manufacturer: "Acme Corp",
		Model:        "TestServer 9000",
		SerialNumber: "SN-TEST-001",
		HostName:     "testserver.example.com",
		PowerState:   "On",
		BIOSVersion:  "2.14.0",
		IndicatorLED: "Off",
		Status:       redfish.Status{State: "Enabled", Health: "OK"},
		ProcessorSummary: redfish.ProcessorSummary{
			Count:  2,
			Model:  "Intel Xeon Gold 6338",
			Status: redfish.Status{State: "Enabled", Health: "OK"},
		},
		MemorySummary: redfish.MemorySummary{
			TotalSystemMemoryGiB: 256,
			Status:               redfish.Status{State: "Enabled", Health: "OK"},
		},
		Processors:        redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/Processors"},
		Memory:            redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/Memory"},
		Storage:           redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/Storage"},
		NetworkInterfaces: redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/NetworkInterfaces"},
		LogServices:       redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/LogServices"},
		Bios:              redfish.ODataLink{ODataID: "/redfish/v1/Systems/1/Bios"},
		Actions: redfish.SystemActions{
			Reset: redfish.ActionTarget{
				Target: "/redfish/v1/Systems/1/Actions/ComputerSystem.Reset",
				AllowableValues: []string{
					"On", "ForceOff", "GracefulShutdown", "GracefulRestart", "ForceRestart",
				},
			},
		},
	},
	"/redfish/v1/Systems/1/Processors": redfish.Collection{
		ODataID:      "/redfish/v1/Systems/1/Processors",
		Name:         "Processors Collection",
		MembersCount: 2,
		Members: []redfish.ODataLink{
			{ODataID: "/redfish/v1/Systems/1/Processors/CPU1"},
			{ODataID: "/redfish/v1/Systems/1/Processors/CPU2"},
		},
	},
	"/redfish/v1/Systems/1/Processors/CPU1": redfish.Processor{
		ODataID:               "/redfish/v1/Systems/1/Processors/CPU1",
		ID:                    "CPU1",
		Name:                  "Processor 1",
		ProcessorType:         "CPU",
		ProcessorArchitecture: "x86",
		InstructionSet:        "x86-64",
		Manufacturer:          "Intel",
		Model:                 "Intel Xeon Gold 6338",
		MaxSpeedMHz:           2000,
		TotalCores:            32,
		TotalThreads:          64,
		Socket:                "CPU1",
		Status:                redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/Processors/CPU2": redfish.Processor{
		ODataID:               "/redfish/v1/Systems/1/Processors/CPU2",
		ID:                    "CPU2",
		Name:                  "Processor 2",
		ProcessorType:         "CPU",
		ProcessorArchitecture: "x86",
		InstructionSet:        "x86-64",
		Manufacturer:          "Intel",
		Model:                 "Intel Xeon Gold 6338",
		MaxSpeedMHz:           2000,
		TotalCores:            32,
		TotalThreads:          64,
		Socket:                "CPU2",
		Status:                redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/Memory": redfish.Collection{
		ODataID:      "/redfish/v1/Systems/1/Memory",
		Name:         "Memory Collection",
		MembersCount: 2,
		Members: []redfish.ODataLink{
			{ODataID: "/redfish/v1/Systems/1/Memory/DIMM1"},
			{ODataID: "/redfish/v1/Systems/1/Memory/DIMM2"},
		},
	},
	"/redfish/v1/Systems/1/Memory/DIMM1": redfish.MemoryModule{
		ODataID:          "/redfish/v1/Systems/1/Memory/DIMM1",
		ID:               "DIMM1",
		Name:             "DIMM Slot 1",
		MemoryType:       "DRAM",
		MemoryDeviceType: "DDR4",
		CapacityMiB:      32768,
		OperatingSpeedMhz: 3200,
		Manufacturer:     "Samsung",
		SerialNumber:     "MEM-001",
		DeviceLocator:    "DIMM.Socket.A1",
		Status:           redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/Memory/DIMM2": redfish.MemoryModule{
		ODataID:          "/redfish/v1/Systems/1/Memory/DIMM2",
		ID:               "DIMM2",
		Name:             "DIMM Slot 2",
		MemoryType:       "DRAM",
		MemoryDeviceType: "DDR4",
		CapacityMiB:      32768,
		OperatingSpeedMhz: 3200,
		Manufacturer:     "Samsung",
		SerialNumber:     "MEM-002",
		DeviceLocator:    "DIMM.Socket.B1",
		Status:           redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/Storage": redfish.Collection{
		ODataID:      "/redfish/v1/Systems/1/Storage",
		Name:         "Storage Collection",
		MembersCount: 1,
		Members: []redfish.ODataLink{
			{ODataID: "/redfish/v1/Systems/1/Storage/RAID.Slot.1-1"},
		},
	},
	"/redfish/v1/Systems/1/Storage/RAID.Slot.1-1": redfish.StorageController{
		ODataID: "/redfish/v1/Systems/1/Storage/RAID.Slot.1-1",
		ID:      "RAID.Slot.1-1",
		Name:    "PERC H750",
		StorageControllers: []redfish.StorageControllerSummary{
			{
				Name:               "PERC H750 Front",
				Manufacturer:       "Dell",
				Model:              "H750",
				FirmwareVersion:    "52.24.3-4839",
				SupportedRAIDTypes: []string{"RAID0", "RAID1", "RAID5", "RAID6", "RAID10"},
				Status:             redfish.Status{State: "Enabled", Health: "OK"},
			},
		},
		Drives: []redfish.ODataLink{
			{ODataID: "/redfish/v1/Systems/1/Storage/RAID.Slot.1-1/Drives/Disk.Bay.0:Enclosure.Internal.0-1:RAID.Slot.1-1"},
		},
		Status: redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/Storage/RAID.Slot.1-1/Drives/Disk.Bay.0:Enclosure.Internal.0-1:RAID.Slot.1-1": redfish.Drive{
		ODataID:       "/redfish/v1/Systems/1/Storage/RAID.Slot.1-1/Drives/Disk.Bay.0:Enclosure.Internal.0-1:RAID.Slot.1-1",
		ID:            "Disk.Bay.0",
		Name:          "Physical Disk 0",
		Manufacturer:  "Seagate",
		Model:         "ST18000NM000J",
		SerialNumber:  "DISK-001",
		MediaType:     "HDD",
		Protocol:      "SAS",
		CapacityBytes: 18000000000000,
		Status:        redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/NetworkInterfaces": redfish.Collection{
		ODataID:      "/redfish/v1/Systems/1/NetworkInterfaces",
		Name:         "Network Interface Collection",
		MembersCount: 1,
		Members: []redfish.ODataLink{
			{ODataID: "/redfish/v1/Systems/1/NetworkInterfaces/NIC.Embedded.1"},
		},
	},
	"/redfish/v1/Systems/1/NetworkInterfaces/NIC.Embedded.1": redfish.NetworkInterface{
		ODataID:     "/redfish/v1/Systems/1/NetworkInterfaces/NIC.Embedded.1",
		ID:          "NIC.Embedded.1",
		Name:        "Embedded NIC 1",
		Description: "Integrated Dual-Port NIC",
		Status:      redfish.Status{State: "Enabled", Health: "OK"},
	},
	"/redfish/v1/Systems/1/LogServices": redfish.Collection{
		ODataID:      "/redfish/v1/Systems/1/LogServices",
		Name:         "Log Service Collection",
		MembersCount: 1,
		Members: []redfish.ODataLink{
			{ODataID: "/redfish/v1/Systems/1/LogServices/System"},
		},
	},
	"/redfish/v1/Systems/1/LogServices/System/Entries": map[string]interface{}{
		"@odata.id":           "/redfish/v1/Systems/1/LogServices/System/Entries",
		"Members@odata.count": 2,
		"Members": []redfish.LogEntry{
			{
				ODataID:   "/redfish/v1/Systems/1/LogServices/System/Entries/1",
				ID:        "1",
				Created:   "2026-05-01T10:00:00Z",
				Severity:  "OK",
				Message:   "System powered on.",
				EntryType: "Event",
			},
			{
				ODataID:   "/redfish/v1/Systems/1/LogServices/System/Entries/2",
				ID:        "2",
				Created:   "2026-05-01T09:00:00Z",
				Severity:  "Warning",
				Message:   "Fan1 reading above threshold.",
				EntryType: "SEL",
			},
		},
	},
	"/redfish/v1/Systems/1/Bios": redfish.Bios{
		ODataID:     "/redfish/v1/Systems/1/Bios",
		ID:          "Bios",
		Name:        "BIOS Configuration Current Settings",
		Description: "BIOS Configuration",
		Attributes: map[string]interface{}{
			"BootMode":              "Uefi",
			"ProcHyperthreading":    "Enabled",
			"NumaNodesPerSocket":    2,
			"ProcVirtualization":    "Enabled",
		},
	},
	"/redfish/v1/Chassis": redfish.Collection{
		ODataID:      "/redfish/v1/Chassis",
		Name:         "Chassis Collection",
		MembersCount: 1,
		Members: []redfish.ODataLink{
			{ODataID: "/redfish/v1/Chassis/System.Embedded.1"},
		},
	},
	"/redfish/v1/Chassis/System.Embedded.1/Thermal": redfish.Thermal{
		ODataID: "/redfish/v1/Chassis/System.Embedded.1/Thermal",
		ID:      "Thermal",
		Name:    "Thermal",
		Temperatures: []redfish.Temperature{
			{
				Name:                      "Inlet Temp",
				ReadingCelsius:            24.0,
				UpperThresholdCritical:    47.0,
				UpperThresholdNonCritical: 42.0,
				Status:                    redfish.Status{State: "Enabled", Health: "OK"},
			},
			{
				Name:           "CPU1 Temp",
				ReadingCelsius: 55.0,
				UpperThresholdCritical: 90.0,
				Status:         redfish.Status{State: "Enabled", Health: "OK"},
			},
		},
		Fans: []redfish.Fan{
			{Name: "Fan.Slot.1", Reading: 5880, ReadingUnits: "RPM", Status: redfish.Status{State: "Enabled", Health: "OK"}},
			{Name: "Fan.Slot.2", Reading: 5760, ReadingUnits: "RPM", Status: redfish.Status{State: "Enabled", Health: "OK"}},
		},
	},
	"/redfish/v1/Chassis/System.Embedded.1/Power": redfish.Power{
		ODataID: "/redfish/v1/Chassis/System.Embedded.1/Power",
		ID:      "Power",
		Name:    "Power",
		PowerControl: []redfish.PowerControl{
			{
				Name:               "System Power Control",
				PowerConsumedWatts: 320,
				PowerCapacityWatts: 750,
				PowerMetrics: redfish.PowerMetrics{
					AverageConsumedWatts: 310,
					MinConsumedWatts:     280,
					MaxConsumedWatts:     380,
				},
				Status: redfish.Status{State: "Enabled", Health: "OK"},
			},
		},
		PowerSupplies: []redfish.PowerSupply{
			{
				Name:                 "PSU.Slot.1",
				Manufacturer:         "Dell",
				Model:                "PWR SPLY 750W",
				PowerCapacityWatts:   750,
				LastPowerOutputWatts: 320,
				PowerSupplyType:      "AC",
				LineInputVoltage:     230,
				Status:               redfish.Status{State: "Enabled", Health: "OK"},
			},
		},
	},
	"/redfish/v1/Managers": redfish.Collection{
		ODataID:      "/redfish/v1/Managers",
		Name:         "Manager Collection",
		MembersCount: 1,
		Members: []redfish.ODataLink{
			{ODataID: "/redfish/v1/Managers/iDRAC.Embedded.1"},
		},
	},
	"/redfish/v1/Managers/iDRAC.Embedded.1": redfish.Manager{
		ODataID:             "/redfish/v1/Managers/iDRAC.Embedded.1",
		ID:                  "iDRAC.Embedded.1",
		Name:                "Manager",
		ManagerType:         "BMC",
		Manufacturer:        "Dell",
		Model:               "iDRAC9",
		FirmwareVersion:     "6.40.30.00",
		UUID:                "3432b9b5-4cdc-8135-6380-5c0719ef6c53",
		PowerState:          "On",
		DateTime:            "2026-05-08T00:00:00+00:00",
		DateTimeLocalOffset: "+00:00",
		Status:              redfish.Status{State: "Enabled", Health: "OK"},
	},
}

// newTestServer creates an httptest server that responds to Redfish paths
// from redfishFixtures. POST/PATCH requests return 204.
func newTestServer(t *testing.T) (*httptest.Server, *redfish.Client) {
	t.Helper()
	mux := http.NewServeMux()
	for path, fixture := range redfishFixtures {
		p, f := path, fixture
		mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(f); err != nil {
					t.Errorf("encoding fixture %s: %v", p, err)
				}
			case http.MethodPost, http.MethodPatch:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.NotFound(w, r)
			}
		})
	}

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	client := redfish.NewClient(redfish.Config{
		Host:     srv.URL,
		Username: "admin",
		Password: "password",
	})
	return srv, client
}
