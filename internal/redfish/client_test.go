package redfish_test

import (
	"context"
	"testing"
)

var ctx = context.Background()

func TestGetServiceRoot(t *testing.T) {
	client := newTestServer(t)
	sr, err := client.GetServiceRoot(ctx)
	if err != nil {
		t.Fatalf("GetServiceRoot() error: %v", err)
	}
	if sr.RedfishVersion != "1.16.0" {
		t.Errorf("RedfishVersion = %q; want 1.16.0", sr.RedfishVersion)
	}
	if sr.Systems.ODataID != "/redfish/v1/Systems" {
		t.Errorf("Systems.ODataID = %q; want /redfish/v1/Systems", sr.Systems.ODataID)
	}
}

func TestListSystems(t *testing.T) {
	client := newTestServer(t)
	systems, err := client.ListSystems(ctx)
	if err != nil {
		t.Fatalf("ListSystems() error: %v", err)
	}
	if len(systems) != 1 {
		t.Fatalf("ListSystems() returned %d systems; want 1", len(systems))
	}
	if systems[0].Model != "TestServer 9000" {
		t.Errorf("Model = %q; want TestServer 9000", systems[0].Model)
	}
}

func TestGetSystem_DefaultsToFirst(t *testing.T) {
	client := newTestServer(t)
	sys, err := client.GetSystem(ctx, "")
	if err != nil {
		t.Fatalf("GetSystem(\"\") error: %v", err)
	}
	if sys.SerialNumber != "SN-TEST-001" {
		t.Errorf("SerialNumber = %q; want SN-TEST-001", sys.SerialNumber)
	}
}

func TestGetSystem_ByID(t *testing.T) {
	client := newTestServer(t)
	sys, err := client.GetSystem(ctx, "1")
	if err != nil {
		t.Fatalf("GetSystem(\"1\") error: %v", err)
	}
	if sys.ID != "1" {
		t.Errorf("ID = %q; want 1", sys.ID)
	}
}

func TestGetProcessors(t *testing.T) {
	client := newTestServer(t)
	procs, err := client.GetProcessors(ctx, "")
	if err != nil {
		t.Fatalf("GetProcessors() error: %v", err)
	}
	if len(procs) != 2 {
		t.Fatalf("GetProcessors() returned %d; want 2", len(procs))
	}
	if procs[0].TotalCores != 32 {
		t.Errorf("TotalCores = %d; want 32", procs[0].TotalCores)
	}
	if procs[0].TotalThreads != 64 {
		t.Errorf("TotalThreads = %d; want 64", procs[0].TotalThreads)
	}
}

func TestGetMemory(t *testing.T) {
	client := newTestServer(t)
	mods, err := client.GetMemory(ctx, "")
	if err != nil {
		t.Fatalf("GetMemory() error: %v", err)
	}
	if len(mods) != 2 {
		t.Fatalf("GetMemory() returned %d modules; want 2", len(mods))
	}
	if mods[0].CapacityMiB != 32768 {
		t.Errorf("CapacityMiB = %d; want 32768", mods[0].CapacityMiB)
	}
}

func TestGetStorage(t *testing.T) {
	client := newTestServer(t)
	controllers, err := client.GetStorage(ctx, "")
	if err != nil {
		t.Fatalf("GetStorage() error: %v", err)
	}
	if len(controllers) != 1 {
		t.Fatalf("GetStorage() returned %d controllers; want 1", len(controllers))
	}
	if len(controllers[0].Drives) != 1 {
		t.Errorf("Drives count = %d; want 1", len(controllers[0].Drives))
	}
}

func TestGetNetworkInterfaces(t *testing.T) {
	client := newTestServer(t)
	nics, err := client.GetNetworkInterfaces(ctx, "")
	if err != nil {
		t.Fatalf("GetNetworkInterfaces() error: %v", err)
	}
	if len(nics) != 1 {
		t.Fatalf("GetNetworkInterfaces() returned %d; want 1", len(nics))
	}
	if nics[0].ID != "NIC.Embedded.1" {
		t.Errorf("ID = %q; want NIC.Embedded.1", nics[0].ID)
	}
}

func TestGetThermal(t *testing.T) {
	client := newTestServer(t)
	thermal, err := client.GetThermal(ctx)
	if err != nil {
		t.Fatalf("GetThermal() error: %v", err)
	}
	if len(thermal.Temperatures) != 2 {
		t.Errorf("Temperatures count = %d; want 2", len(thermal.Temperatures))
	}
	if len(thermal.Fans) != 2 {
		t.Errorf("Fans count = %d; want 2", len(thermal.Fans))
	}
	if thermal.Temperatures[0].ReadingCelsius != 24.0 {
		t.Errorf("ReadingCelsius = %f; want 24.0", thermal.Temperatures[0].ReadingCelsius)
	}
}

func TestGetPower(t *testing.T) {
	client := newTestServer(t)
	power, err := client.GetPower(ctx)
	if err != nil {
		t.Fatalf("GetPower() error: %v", err)
	}
	if len(power.PowerControl) != 1 {
		t.Errorf("PowerControl count = %d; want 1", len(power.PowerControl))
	}
	if power.PowerControl[0].PowerConsumedWatts != 320 {
		t.Errorf("PowerConsumedWatts = %f; want 320", power.PowerControl[0].PowerConsumedWatts)
	}
}

func TestGetEventLog(t *testing.T) {
	client := newTestServer(t)
	entries, err := client.GetEventLog(ctx, "", 0)
	if err != nil {
		t.Fatalf("GetEventLog() error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("entries count = %d; want 2", len(entries))
	}
}

func TestGetEventLog_Limit(t *testing.T) {
	client := newTestServer(t)
	entries, err := client.GetEventLog(ctx, "", 1)
	if err != nil {
		t.Fatalf("GetEventLog() error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("entries count = %d; want 1", len(entries))
	}
}

func TestListManagers(t *testing.T) {
	client := newTestServer(t)
	mgrs, err := client.ListManagers(ctx)
	if err != nil {
		t.Fatalf("ListManagers() error: %v", err)
	}
	if len(mgrs) != 1 {
		t.Fatalf("ListManagers() returned %d; want 1", len(mgrs))
	}
	if mgrs[0].FirmwareVersion != "6.40.30.00" {
		t.Errorf("FirmwareVersion = %q; want 6.40.30.00", mgrs[0].FirmwareVersion)
	}
}

func TestGetManager_DefaultsToFirst(t *testing.T) {
	client := newTestServer(t)
	mgr, err := client.GetManager(ctx, "")
	if err != nil {
		t.Fatalf("GetManager(\"\") error: %v", err)
	}
	if mgr.ManagerType != "BMC" {
		t.Errorf("ManagerType = %q; want BMC", mgr.ManagerType)
	}
}

func TestGetBios(t *testing.T) {
	client := newTestServer(t)
	bios, err := client.GetBios(ctx, "")
	if err != nil {
		t.Fatalf("GetBios() error: %v", err)
	}
	if bios.Attributes["BootMode"] != "Uefi" {
		t.Errorf("BootMode = %v; want Uefi", bios.Attributes["BootMode"])
	}
}

func TestResetSystem(t *testing.T) {
	client := newTestServer(t)
	if err := client.ResetSystem(ctx, "", "GracefulShutdown"); err != nil {
		t.Fatalf("ResetSystem() error: %v", err)
	}
}

func TestSetIndicatorLED(t *testing.T) {
	client := newTestServer(t)
	if err := client.SetIndicatorLED(ctx, "", "Blinking"); err != nil {
		t.Fatalf("SetIndicatorLED() error: %v", err)
	}
}

func TestClearEventLog(t *testing.T) {
	client := newTestServer(t)
	if err := client.ClearEventLog(ctx, ""); err != nil {
		t.Fatalf("ClearEventLog() error: %v", err)
	}
}

func TestSetBiosAttribute(t *testing.T) {
	client := newTestServer(t)
	if err := client.SetBiosAttribute(ctx, "", "BootMode", "Bios"); err != nil {
		t.Fatalf("SetBiosAttribute() error: %v", err)
	}
}

func TestNewClient_HTTPPrefix(t *testing.T) {
	client := newTestServer(t)
	sr, err := client.GetServiceRoot(ctx)
	if err != nil {
		t.Fatalf("GetServiceRoot() with http:// prefix error: %v", err)
	}
	if sr.ID != "RootService" {
		t.Errorf("ID = %q; want RootService", sr.ID)
	}
}
