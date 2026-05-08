package redfish

// ODataLink is a reference to another Redfish resource.
type ODataLink struct {
	ODataID string `json:"@odata.id"`
}

// Collection is a generic Redfish resource collection.
type Collection struct {
	ODataID      string      `json:"@odata.id"`
	Name         string      `json:"Name"`
	MembersCount int         `json:"Members@odata.count"`
	Members      []ODataLink `json:"Members"`
}

// ServiceRoot is the Redfish API entry point.
type ServiceRoot struct {
	ODataID        string    `json:"@odata.id"`
	ID             string    `json:"Id"`
	Name           string    `json:"Name"`
	RedfishVersion string    `json:"RedfishVersion"`
	UUID           string    `json:"UUID"`
	Product        string    `json:"Product"`
	Vendor         string    `json:"Vendor"`
	Systems        ODataLink `json:"Systems"`
	Chassis        ODataLink `json:"Chassis"`
	Managers       ODataLink `json:"Managers"`
}

// Status represents health and state of a Redfish resource.
type Status struct {
	State        string `json:"State"`
	Health       string `json:"Health"`
	HealthRollup string `json:"HealthRollup"`
}

// ProcessorSummary is a summary of system processor(s).
type ProcessorSummary struct {
	Count  int    `json:"Count"`
	Model  string `json:"Model"`
	Status Status `json:"Status"`
}

// MemorySummary is a summary of system memory.
type MemorySummary struct {
	TotalSystemMemoryGiB float64 `json:"TotalSystemMemoryGiB"`
	Status               Status  `json:"Status"`
}

// Boot contains boot configuration.
type Boot struct {
	BootSourceOverrideEnabled string `json:"BootSourceOverrideEnabled"`
	BootSourceOverrideTarget  string `json:"BootSourceOverrideTarget"`
	BootSourceOverrideMode    string `json:"BootSourceOverrideMode"`
}

// ActionTarget describes a Redfish action endpoint.
type ActionTarget struct {
	Target          string   `json:"target"`
	AllowableValues []string `json:"ResetType@Redfish.AllowableValues"`
}

// SystemActions lists actions available on a ComputerSystem.
type SystemActions struct {
	Reset ActionTarget `json:"#ComputerSystem.Reset"`
}

// SystemLinks lists related resources for a ComputerSystem.
type SystemLinks struct {
	Chassis   []ODataLink `json:"Chassis"`
	ManagedBy []ODataLink `json:"ManagedBy"`
}

// ComputerSystem represents a Redfish managed server.
type ComputerSystem struct {
	ODataID            string           `json:"@odata.id"`
	ID                 string           `json:"Id"`
	Name               string           `json:"Name"`
	Description        string           `json:"Description"`
	SystemType         string           `json:"SystemType"`
	Manufacturer       string           `json:"Manufacturer"`
	Model              string           `json:"Model"`
	SKU                string           `json:"SKU"`
	SerialNumber       string           `json:"SerialNumber"`
	PartNumber         string           `json:"PartNumber"`
	HostName           string           `json:"HostName"`
	Status             Status           `json:"Status"`
	IndicatorLED       string           `json:"IndicatorLED"`
	PowerState         string           `json:"PowerState"`
	BIOSVersion        string           `json:"BiosVersion"`
	ProcessorSummary   ProcessorSummary `json:"ProcessorSummary"`
	MemorySummary      MemorySummary    `json:"MemorySummary"`
	Boot               Boot             `json:"Boot"`
	Links              SystemLinks      `json:"Links"`
	Actions            SystemActions    `json:"Actions"`
	Processors         ODataLink        `json:"Processors"`
	Memory             ODataLink        `json:"Memory"`
	Storage            ODataLink        `json:"Storage"`
	NetworkInterfaces  ODataLink        `json:"NetworkInterfaces"`
	EthernetInterfaces ODataLink        `json:"EthernetInterfaces"`
	LogServices        ODataLink        `json:"LogServices"`
	Bios               ODataLink        `json:"Bios"`
}

// Processor represents a CPU installed in the system.
type Processor struct {
	ODataID               string `json:"@odata.id"`
	ID                    string `json:"Id"`
	Name                  string `json:"Name"`
	ProcessorType         string `json:"ProcessorType"`
	ProcessorArchitecture string `json:"ProcessorArchitecture"`
	InstructionSet        string `json:"InstructionSet"`
	Manufacturer          string `json:"Manufacturer"`
	Model                 string `json:"Model"`
	MaxSpeedMHz           int    `json:"MaxSpeedMHz"`
	TotalCores            int    `json:"TotalCores"`
	TotalThreads          int    `json:"TotalThreads"`
	Socket                string `json:"Socket"`
	Status                Status `json:"Status"`
}

// MemoryModule represents a single DIMM or memory device.
type MemoryModule struct {
	ODataID           string `json:"@odata.id"`
	ID                string `json:"Id"`
	Name              string `json:"Name"`
	MemoryType        string `json:"MemoryType"`
	MemoryDeviceType  string `json:"MemoryDeviceType"`
	CapacityMiB       int    `json:"CapacityMiB"`
	OperatingSpeedMhz int    `json:"OperatingSpeedMhz"`
	Manufacturer      string `json:"Manufacturer"`
	SerialNumber      string `json:"SerialNumber"`
	PartNumber        string `json:"PartNumber"`
	DeviceLocator     string `json:"DeviceLocator"`
	Status            Status `json:"Status"`
}

// StorageControllerSummary describes a storage controller.
type StorageControllerSummary struct {
	Name               string   `json:"Name"`
	Manufacturer       string   `json:"Manufacturer"`
	Model              string   `json:"Model"`
	FirmwareVersion    string   `json:"FirmwareVersion"`
	SupportedRAIDTypes []string `json:"SupportedRAIDTypes"`
	Status             Status   `json:"Status"`
}

// StorageController represents a storage subsystem.
type StorageController struct {
	ODataID            string                     `json:"@odata.id"`
	ID                 string                     `json:"Id"`
	Name               string                     `json:"Name"`
	StorageControllers []StorageControllerSummary `json:"StorageControllers"`
	Drives             []ODataLink                `json:"Drives"`
	Status             Status                     `json:"Status"`
}

// Drive represents a physical storage drive.
type Drive struct {
	ODataID          string  `json:"@odata.id"`
	ID               string  `json:"Id"`
	Name             string  `json:"Name"`
	Manufacturer     string  `json:"Manufacturer"`
	Model            string  `json:"Model"`
	SerialNumber     string  `json:"SerialNumber"`
	MediaType        string  `json:"MediaType"`
	Protocol         string  `json:"Protocol"`
	CapacityBytes    int64   `json:"CapacityBytes"`
	RotationSpeedRPM float64 `json:"RotationSpeedRPM"`
	Status           Status  `json:"Status"`
}

// NetworkInterface represents a NIC or network adapter.
type NetworkInterface struct {
	ODataID     string `json:"@odata.id"`
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Status      Status `json:"Status"`
}

// Temperature is a single thermal sensor reading.
type Temperature struct {
	Name                      string  `json:"Name"`
	ReadingCelsius            float64 `json:"ReadingCelsius"`
	UpperThresholdFatal       float64 `json:"UpperThresholdFatal"`
	UpperThresholdCritical    float64 `json:"UpperThresholdCritical"`
	UpperThresholdNonCritical float64 `json:"UpperThresholdNonCritical"`
	LowerThresholdFatal       float64 `json:"LowerThresholdFatal"`
	LowerThresholdCritical    float64 `json:"LowerThresholdCritical"`
	Status                    Status  `json:"Status"`
}

// Fan is a single fan sensor reading.
type Fan struct {
	Name                string `json:"Name"`
	Reading             int    `json:"Reading"`
	ReadingUnits        string `json:"ReadingUnits"`
	LowerThresholdFatal int    `json:"LowerThresholdFatal"`
	Status              Status `json:"Status"`
}

// Thermal aggregates temperature and fan readings for a chassis.
type Thermal struct {
	ODataID      string        `json:"@odata.id"`
	ID           string        `json:"Id"`
	Name         string        `json:"Name"`
	Temperatures []Temperature `json:"Temperatures"`
	Fans         []Fan         `json:"Fans"`
}

// PowerMetrics contains power usage statistics.
type PowerMetrics struct {
	IntervalInMin        int     `json:"IntervalInMin"`
	MinConsumedWatts     float64 `json:"MinConsumedWatts"`
	MaxConsumedWatts     float64 `json:"MaxConsumedWatts"`
	AverageConsumedWatts float64 `json:"AverageConsumedWatts"`
}

// PowerControl represents a power domain with consumption data.
type PowerControl struct {
	Name                string       `json:"Name"`
	PowerConsumedWatts  float64      `json:"PowerConsumedWatts"`
	PowerCapacityWatts  float64      `json:"PowerCapacityWatts"`
	PowerAllocatedWatts float64      `json:"PowerAllocatedWatts"`
	PowerMetrics        PowerMetrics `json:"PowerMetrics"`
	Status              Status       `json:"Status"`
}

// Voltage is a single voltage sensor reading.
type Voltage struct {
	Name                   string  `json:"Name"`
	ReadingVolts           float64 `json:"ReadingVolts"`
	UpperThresholdFatal    float64 `json:"UpperThresholdFatal"`
	UpperThresholdCritical float64 `json:"UpperThresholdCritical"`
	LowerThresholdFatal    float64 `json:"LowerThresholdFatal"`
	LowerThresholdCritical float64 `json:"LowerThresholdCritical"`
	Status                 Status  `json:"Status"`
}

// PowerSupply represents a PSU installed in the chassis.
type PowerSupply struct {
	Name                 string  `json:"Name"`
	Manufacturer         string  `json:"Manufacturer"`
	Model                string  `json:"Model"`
	SerialNumber         string  `json:"SerialNumber"`
	PowerCapacityWatts   float64 `json:"PowerCapacityWatts"`
	LastPowerOutputWatts float64 `json:"LastPowerOutputWatts"`
	PowerSupplyType      string  `json:"PowerSupplyType"`
	LineInputVoltage     float64 `json:"LineInputVoltage"`
	Status               Status  `json:"Status"`
}

// Power aggregates power control, voltage, and PSU data for a chassis.
type Power struct {
	ODataID       string         `json:"@odata.id"`
	ID            string         `json:"Id"`
	Name          string         `json:"Name"`
	PowerControl  []PowerControl `json:"PowerControl"`
	Voltages      []Voltage      `json:"Voltages"`
	PowerSupplies []PowerSupply  `json:"PowerSupplies"`
}

// LogEntry represents a single entry in a system event log.
type LogEntry struct {
	ODataID      string `json:"@odata.id"`
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	Created      string `json:"Created"`
	Severity     string `json:"Severity"`
	Message      string `json:"Message"`
	MessageID    string `json:"MessageId"`
	EntryType    string `json:"EntryType"`
	SensorType   string `json:"SensorType"`
	SensorNumber int    `json:"SensorNumber"`
}

// ManagerActions lists actions available on a Manager.
type ManagerActions struct {
	Reset ActionTarget `json:"#Manager.Reset"`
}

// Manager represents a BMC or management controller.
type Manager struct {
	ODataID             string         `json:"@odata.id"`
	ID                  string         `json:"Id"`
	Name                string         `json:"Name"`
	Description         string         `json:"Description"`
	ManagerType         string         `json:"ManagerType"`
	Manufacturer        string         `json:"Manufacturer"`
	Model               string         `json:"Model"`
	FirmwareVersion     string         `json:"FirmwareVersion"`
	UUID                string         `json:"UUID"`
	Status              Status         `json:"Status"`
	PowerState          string         `json:"PowerState"`
	DateTime            string         `json:"DateTime"`
	DateTimeLocalOffset string         `json:"DateTimeLocalOffset"`
	EthernetInterfaces  ODataLink      `json:"EthernetInterfaces"`
	NetworkProtocol     ODataLink      `json:"NetworkProtocol"`
	LogServices         ODataLink      `json:"LogServices"`
	Actions             ManagerActions `json:"Actions"`
}

// Bios represents BIOS/UEFI firmware settings.
type Bios struct {
	ODataID     string                 `json:"@odata.id"`
	ID          string                 `json:"Id"`
	Name        string                 `json:"Name"`
	Description string                 `json:"Description"`
	Attributes  map[string]interface{} `json:"Attributes"`
}

// ResetRequest is the payload for a system or manager reset action.
type ResetRequest struct {
	ResetType string `json:"ResetType"`
}

// SetAttributesRequest is the payload for updating BIOS attributes.
type SetAttributesRequest struct {
	Attributes map[string]interface{} `json:"Attributes"`
}

// IndicatorLEDRequest is the payload for setting the indicator LED.
type IndicatorLEDRequest struct {
	IndicatorLED string `json:"IndicatorLED"`
}
