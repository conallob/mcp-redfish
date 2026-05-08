package redfish

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Config holds the connection parameters for a Redfish BMC.
type Config struct {
	Host     string
	Username string
	Password string
	Insecure bool
}

// Client makes authenticated requests to a Redfish BMC.
type Client struct {
	cfg        Config
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a Client from the given config.
func NewClient(cfg Config) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.Insecure, //nolint:gosec
		},
	}

	host := cfg.Host
	var baseURL string
	switch {
	case strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://"):
		baseURL = strings.TrimRight(host, "/")
	default:
		baseURL = "https://" + strings.TrimRight(host, "/")
	}

	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// get performs a GET request and decodes the JSON body into v.
func (c *Client) get(ctx context.Context, path string, v interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if c.cfg.Username != "" {
		req.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request GET %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GET %s returned HTTP %d: %s", path, resp.StatusCode, trimBody(body))
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

// post performs a POST request, encoding body as JSON, and decodes the response into v (may be nil).
func (c *Client) post(ctx context.Context, path string, body interface{}, v interface{}) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encoding request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.cfg.Username != "" {
		req.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request POST %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("POST %s returned HTTP %d: %s", path, resp.StatusCode, trimBody(body))
	}
	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

// patch performs a PATCH request encoding body as JSON.
func (c *Client) patch(ctx context.Context, path string, body interface{}) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encoding request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.baseURL+path, bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.cfg.Username != "" {
		req.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request PATCH %s: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("PATCH %s returned HTTP %d: %s", path, resp.StatusCode, trimBody(body))
	}
	return nil
}

// --- Service root ---

// GetServiceRoot retrieves the Redfish service root.
func (c *Client) GetServiceRoot(ctx context.Context) (*ServiceRoot, error) {
	var sr ServiceRoot
	return &sr, c.get(ctx, "/redfish/v1/", &sr)
}

// --- Systems ---

// ListSystems returns all computer systems in the Systems collection.
func (c *Client) ListSystems(ctx context.Context) ([]ComputerSystem, error) {
	var coll Collection
	if err := c.get(ctx, "/redfish/v1/Systems", &coll); err != nil {
		return nil, err
	}
	systems := make([]ComputerSystem, 0, len(coll.Members))
	for _, m := range coll.Members {
		var sys ComputerSystem
		if err := c.get(ctx, m.ODataID, &sys); err != nil {
			return nil, fmt.Errorf("fetching system %s: %w", m.ODataID, err)
		}
		systems = append(systems, sys)
	}
	return systems, nil
}

// GetSystem retrieves the computer system at the given OData path.
// If id is empty, it returns the first system in the collection.
func (c *Client) GetSystem(ctx context.Context, id string) (*ComputerSystem, error) {
	path, err := c.resolveSystemPath(ctx, id)
	if err != nil {
		return nil, err
	}
	var sys ComputerSystem
	return &sys, c.get(ctx, path, &sys)
}

// resolveSystemPath returns the OData path for a system ID.
// If id is empty, the first system in the collection is used.
func (c *Client) resolveSystemPath(ctx context.Context, id string) (string, error) {
	if id != "" {
		if strings.HasPrefix(id, "/") {
			return id, nil
		}
		return "/redfish/v1/Systems/" + id, nil
	}
	var coll Collection
	if err := c.get(ctx, "/redfish/v1/Systems", &coll); err != nil {
		return "", err
	}
	if len(coll.Members) == 0 {
		return "", fmt.Errorf("no systems found")
	}
	return coll.Members[0].ODataID, nil
}

// GetProcessors returns all processors for the given system.
func (c *Client) GetProcessors(ctx context.Context, systemID string) ([]Processor, error) {
	sysPath, err := c.resolveSystemPath(ctx, systemID)
	if err != nil {
		return nil, err
	}
	var sys ComputerSystem
	if err := c.get(ctx, sysPath, &sys); err != nil {
		return nil, err
	}
	var coll Collection
	if err := c.get(ctx, sys.Processors.ODataID, &coll); err != nil {
		return nil, err
	}
	procs := make([]Processor, 0, len(coll.Members))
	for _, m := range coll.Members {
		var p Processor
		if err := c.get(ctx, m.ODataID, &p); err != nil {
			return nil, fmt.Errorf("fetching processor %s: %w", m.ODataID, err)
		}
		procs = append(procs, p)
	}
	return procs, nil
}

// GetMemory returns all memory modules for the given system.
func (c *Client) GetMemory(ctx context.Context, systemID string) ([]MemoryModule, error) {
	sysPath, err := c.resolveSystemPath(ctx, systemID)
	if err != nil {
		return nil, err
	}
	var sys ComputerSystem
	if err := c.get(ctx, sysPath, &sys); err != nil {
		return nil, err
	}
	var coll Collection
	if err := c.get(ctx, sys.Memory.ODataID, &coll); err != nil {
		return nil, err
	}
	mods := make([]MemoryModule, 0, len(coll.Members))
	for _, m := range coll.Members {
		var mod MemoryModule
		if err := c.get(ctx, m.ODataID, &mod); err != nil {
			return nil, fmt.Errorf("fetching memory module %s: %w", m.ODataID, err)
		}
		mods = append(mods, mod)
	}
	return mods, nil
}

// GetStorage returns all storage controllers (with their drives) for the given system.
func (c *Client) GetStorage(ctx context.Context, systemID string) ([]StorageController, error) {
	sysPath, err := c.resolveSystemPath(ctx, systemID)
	if err != nil {
		return nil, err
	}
	var sys ComputerSystem
	if err := c.get(ctx, sysPath, &sys); err != nil {
		return nil, err
	}
	var coll Collection
	if err := c.get(ctx, sys.Storage.ODataID, &coll); err != nil {
		return nil, err
	}
	controllers := make([]StorageController, 0, len(coll.Members))
	for _, m := range coll.Members {
		var sc StorageController
		if err := c.get(ctx, m.ODataID, &sc); err != nil {
			return nil, fmt.Errorf("fetching storage controller %s: %w", m.ODataID, err)
		}
		controllers = append(controllers, sc)
	}
	return controllers, nil
}

// GetDrive retrieves a single drive by its OData path.
func (c *Client) GetDrive(ctx context.Context, drivePath string) (*Drive, error) {
	var d Drive
	return &d, c.get(ctx, drivePath, &d)
}

// GetNetworkInterfaces returns all network interfaces for the given system.
func (c *Client) GetNetworkInterfaces(ctx context.Context, systemID string) ([]NetworkInterface, error) {
	sysPath, err := c.resolveSystemPath(ctx, systemID)
	if err != nil {
		return nil, err
	}
	var sys ComputerSystem
	if err := c.get(ctx, sysPath, &sys); err != nil {
		return nil, err
	}
	var coll Collection
	if err := c.get(ctx, sys.NetworkInterfaces.ODataID, &coll); err != nil {
		return nil, err
	}
	nics := make([]NetworkInterface, 0, len(coll.Members))
	for _, m := range coll.Members {
		var nic NetworkInterface
		if err := c.get(ctx, m.ODataID, &nic); err != nil {
			return nil, fmt.Errorf("fetching network interface %s: %w", m.ODataID, err)
		}
		nics = append(nics, nic)
	}
	return nics, nil
}

// GetThermal retrieves thermal sensor data from the first chassis.
func (c *Client) GetThermal(ctx context.Context) (*Thermal, error) {
	chassisPath, err := c.firstChassisPath(ctx)
	if err != nil {
		return nil, err
	}
	var thermal Thermal
	return &thermal, c.get(ctx, chassisPath+"/Thermal", &thermal)
}

// GetPower retrieves power data from the first chassis.
func (c *Client) GetPower(ctx context.Context) (*Power, error) {
	chassisPath, err := c.firstChassisPath(ctx)
	if err != nil {
		return nil, err
	}
	var power Power
	return &power, c.get(ctx, chassisPath+"/Power", &power)
}

// firstChassisPath returns the OData path of the first chassis.
func (c *Client) firstChassisPath(ctx context.Context) (string, error) {
	var coll Collection
	if err := c.get(ctx, "/redfish/v1/Chassis", &coll); err != nil {
		return "", err
	}
	if len(coll.Members) == 0 {
		return "", fmt.Errorf("no chassis found")
	}
	return coll.Members[0].ODataID, nil
}

// GetEventLog returns log entries from the System log service.
func (c *Client) GetEventLog(ctx context.Context, systemID string, limit int) ([]LogEntry, error) {
	sysPath, err := c.resolveSystemPath(ctx, systemID)
	if err != nil {
		return nil, err
	}
	var sys ComputerSystem
	if err := c.get(ctx, sysPath, &sys); err != nil {
		return nil, err
	}
	var logServicesColl Collection
	if err := c.get(ctx, sys.LogServices.ODataID, &logServicesColl); err != nil {
		return nil, err
	}

	// Prefer a log service named "System" or "Log"; fall back to the first one.
	logPath := ""
	for _, m := range logServicesColl.Members {
		if strings.Contains(m.ODataID, "System") || strings.Contains(m.ODataID, "Log") {
			logPath = m.ODataID
			break
		}
	}
	if logPath == "" && len(logServicesColl.Members) > 0 {
		logPath = logServicesColl.Members[0].ODataID
	}
	if logPath == "" {
		return nil, fmt.Errorf("no log service found")
	}

	var entriesColl struct {
		Members      []LogEntry `json:"Members"`
		MembersCount int        `json:"Members@odata.count"`
	}
	if err := c.get(ctx, logPath+"/Entries", &entriesColl); err != nil {
		return nil, err
	}

	entries := entriesColl.Members
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

// --- Managers ---

// ListManagers returns all managers (BMCs) in the system.
func (c *Client) ListManagers(ctx context.Context) ([]Manager, error) {
	var coll Collection
	if err := c.get(ctx, "/redfish/v1/Managers", &coll); err != nil {
		return nil, err
	}
	managers := make([]Manager, 0, len(coll.Members))
	for _, m := range coll.Members {
		var mgr Manager
		if err := c.get(ctx, m.ODataID, &mgr); err != nil {
			return nil, fmt.Errorf("fetching manager %s: %w", m.ODataID, err)
		}
		managers = append(managers, mgr)
	}
	return managers, nil
}

// GetManager retrieves a manager by ID. If id is empty, returns the first manager.
func (c *Client) GetManager(ctx context.Context, id string) (*Manager, error) {
	path, err := c.resolveManagerPath(ctx, id)
	if err != nil {
		return nil, err
	}
	var mgr Manager
	return &mgr, c.get(ctx, path, &mgr)
}

func (c *Client) resolveManagerPath(ctx context.Context, id string) (string, error) {
	if id != "" {
		if strings.HasPrefix(id, "/") {
			return id, nil
		}
		return "/redfish/v1/Managers/" + id, nil
	}
	var coll Collection
	if err := c.get(ctx, "/redfish/v1/Managers", &coll); err != nil {
		return "", err
	}
	if len(coll.Members) == 0 {
		return "", fmt.Errorf("no managers found")
	}
	return coll.Members[0].ODataID, nil
}

// GetBios retrieves the BIOS resource for the given system.
func (c *Client) GetBios(ctx context.Context, systemID string) (*Bios, error) {
	sysPath, err := c.resolveSystemPath(ctx, systemID)
	if err != nil {
		return nil, err
	}
	var sys ComputerSystem
	if err := c.get(ctx, sysPath, &sys); err != nil {
		return nil, err
	}
	var bios Bios
	return &bios, c.get(ctx, sys.Bios.ODataID, &bios)
}

// --- Write operations ---

// ResetSystem posts a reset action to the given system.
func (c *Client) ResetSystem(ctx context.Context, systemID, resetType string) error {
	sysPath, err := c.resolveSystemPath(ctx, systemID)
	if err != nil {
		return err
	}
	var sys ComputerSystem
	if err := c.get(ctx, sysPath, &sys); err != nil {
		return err
	}
	target := sys.Actions.Reset.Target
	if target == "" {
		target = sysPath + "/Actions/ComputerSystem.Reset"
	}
	return c.post(ctx, target, ResetRequest{ResetType: resetType}, nil)
}

// SetIndicatorLED patches the IndicatorLED field on the given system.
func (c *Client) SetIndicatorLED(ctx context.Context, systemID, state string) error {
	sysPath, err := c.resolveSystemPath(ctx, systemID)
	if err != nil {
		return err
	}
	return c.patch(ctx, sysPath, IndicatorLEDRequest{IndicatorLED: state})
}

// ClearEventLog posts a ClearLog action to the System log service.
func (c *Client) ClearEventLog(ctx context.Context, systemID string) error {
	sysPath, err := c.resolveSystemPath(ctx, systemID)
	if err != nil {
		return err
	}
	var sys ComputerSystem
	if err := c.get(ctx, sysPath, &sys); err != nil {
		return err
	}
	var logServicesColl Collection
	if err := c.get(ctx, sys.LogServices.ODataID, &logServicesColl); err != nil {
		return err
	}
	logPath := ""
	for _, m := range logServicesColl.Members {
		if strings.Contains(m.ODataID, "System") || strings.Contains(m.ODataID, "Log") {
			logPath = m.ODataID
			break
		}
	}
	if logPath == "" && len(logServicesColl.Members) > 0 {
		logPath = logServicesColl.Members[0].ODataID
	}
	if logPath == "" {
		return fmt.Errorf("no log service found")
	}
	return c.post(ctx, logPath+"/Actions/LogService.ClearLog", struct{}{}, nil)
}

// SetBiosAttribute patches a single BIOS attribute on the given system.
func (c *Client) SetBiosAttribute(ctx context.Context, systemID, key string, value interface{}) error {
	sysPath, err := c.resolveSystemPath(ctx, systemID)
	if err != nil {
		return err
	}
	var sys ComputerSystem
	if err := c.get(ctx, sysPath, &sys); err != nil {
		return err
	}
	biosSettingsPath := sys.Bios.ODataID + "/Settings"
	return c.patch(ctx, biosSettingsPath, SetAttributesRequest{
		Attributes: map[string]interface{}{key: value},
	})
}

// trimBody truncates a response body for display in error messages.
func trimBody(b []byte) string {
	const max = 256
	if len(b) > max {
		return string(b[:max]) + "..."
	}
	return string(b)
}
