package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

// --- list_devices ---

// ListDevices implements the list_devices MCP tool.
type ListDevices struct {
	baseTool
}

// NewListDevices creates a new ListDevices tool.
func NewListDevices(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListDevices {
	return &ListDevices{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *ListDevices) Description() string {
	return "List adopted devices for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListDevices) InputSchema() map[string]interface{} {
	return listSchema()
}

// Execute runs the tool.
func (t *ListDevices) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID string `json:"siteId"`
		Limit  *int32 `json:"limit"`
		Offset *int32 `json:"offset"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	siteID, err := resolveSiteID(params.SiteID, t.defaultSiteID)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetAdoptedDeviceOverviewPageWithResponse(
		ctx,
		siteID,
		&unifi.GetAdoptedDeviceOverviewPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to list devices: %w", err)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	page := resp.JSON200
	if len(page.Data) == 0 {
		return "No adopted devices found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Adopted Devices (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, d := range page.Data {
		fw := ""
		if d.FirmwareVersion != nil {
			fw = *d.FirmwareVersion
		}
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Model: %s\n   MAC: %s\n"+
				"   IP: %s\n   State: %s\n   Firmware: %s\n",
			i+1,
			d.Name,
			d.Id.String(),
			d.Model,
			d.MacAddress,
			d.IpAddress,
			d.State,
			fw,
		)
	}

	return b.String(), nil
}

// --- get_device ---

// GetDevice implements the get_device MCP tool.
type GetDevice struct {
	baseTool
}

// NewGetDevice creates a new GetDevice tool.
func NewGetDevice(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetDevice {
	return &GetDevice{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *GetDevice) Description() string {
	return "Get detailed information about an adopted device"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetDevice) InputSchema() map[string]interface{} {
	return siteAndIDSchema("deviceId", "Device UUID")
}

// Execute runs the tool.
func (t *GetDevice) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID   string `json:"siteId"`
		DeviceID string `json:"deviceId"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	siteID, err := resolveSiteID(params.SiteID, t.defaultSiteID)
	if err != nil {
		return "", err
	}

	deviceID, err := resolveUUID("deviceId", params.DeviceID)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetAdoptedDeviceDetailsWithResponse(
		ctx,
		siteID,
		deviceID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to get device: %w", err)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatDeviceDetails(resp.JSON200), nil
}

func formatDeviceDetails(d *unifi.AdoptedDeviceDetails) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\n", d.Name)
	fmt.Fprintf(&b, "ID: %s\n", d.Id.String())
	fmt.Fprintf(&b, "Model: %s\n", d.Model)
	fmt.Fprintf(&b, "MAC Address: %s\n", d.MacAddress)
	fmt.Fprintf(&b, "IP Address: %s\n", d.IpAddress)
	fmt.Fprintf(&b, "State: %s\n", d.State)
	fmt.Fprintf(&b, "Supported: %v\n", d.Supported)
	fmt.Fprintf(&b, "Firmware Updatable: %v\n", d.FirmwareUpdatable)
	if d.FirmwareVersion != nil {
		fmt.Fprintf(&b, "Firmware Version: %s\n", *d.FirmwareVersion)
	}
	if d.AdoptedAt != nil {
		fmt.Fprintf(&b, "Adopted At: %s\n", d.AdoptedAt.String())
	}
	if d.ProvisionedAt != nil {
		fmt.Fprintf(
			&b,
			"Provisioned At: %s\n",
			d.ProvisionedAt.String(),
		)
	}
	featuresJSON, err := json.MarshalIndent(d.Features, "", "  ")
	if err == nil {
		fmt.Fprintf(&b, "Features: %s\n", string(featuresJSON))
	}
	return b.String()
}

// --- adopt_device ---

// AdoptDevice implements the adopt_device MCP tool.
type AdoptDevice struct {
	baseTool
}

// NewAdoptDevice creates a new AdoptDevice tool.
func NewAdoptDevice(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *AdoptDevice {
	return &AdoptDevice{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *AdoptDevice) Description() string {
	return "Adopt a device into a site by MAC address"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *AdoptDevice) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"macAddress": map[string]interface{}{
				"type":        "string",
				"description": "MAC address of the device to adopt",
			},
		},
		"required": []string{"macAddress"},
	}
}

// Execute runs the tool.
func (t *AdoptDevice) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID     string `json:"siteId"`
		MacAddress string `json:"macAddress"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	siteID, err := resolveSiteID(params.SiteID, t.defaultSiteID)
	if err != nil {
		return "", err
	}

	if params.MacAddress == "" {
		return "", fmt.Errorf("macAddress is required")
	}

	resp, err := t.client.AdoptDeviceWithResponse(
		ctx,
		siteID,
		unifi.IntegrationDeviceAdoptionRequestDto{
			MacAddress:        params.MacAddress,
			IgnoreDeviceLimit: false,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to adopt device: %w", err)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatDeviceDetails(resp.JSON200), nil
}

// --- remove_device ---

// RemoveDevice implements the remove_device MCP tool.
type RemoveDevice struct {
	baseTool
}

// NewRemoveDevice creates a new RemoveDevice tool.
func NewRemoveDevice(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *RemoveDevice {
	return &RemoveDevice{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *RemoveDevice) Description() string {
	return "Remove (un-adopt) an adopted device from a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *RemoveDevice) InputSchema() map[string]interface{} {
	return siteAndIDSchema("deviceId", "Device UUID")
}

// Execute runs the tool.
func (t *RemoveDevice) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID   string `json:"siteId"`
		DeviceID string `json:"deviceId"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	siteID, err := resolveSiteID(params.SiteID, t.defaultSiteID)
	if err != nil {
		return "", err
	}

	deviceID, err := resolveUUID("deviceId", params.DeviceID)
	if err != nil {
		return "", err
	}

	resp, err := t.client.RemoveDeviceWithResponse(
		ctx,
		siteID,
		deviceID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to remove device: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return fmt.Sprintf(
		"Device %s removed successfully.",
		params.DeviceID,
	), nil
}

// --- execute_device_action ---

// ExecuteDeviceAction implements the execute_device_action MCP tool.
type ExecuteDeviceAction struct {
	baseTool
}

// NewExecuteDeviceAction creates a new ExecuteDeviceAction tool.
func NewExecuteDeviceAction(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ExecuteDeviceAction {
	return &ExecuteDeviceAction{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *ExecuteDeviceAction) Description() string {
	return "Execute an action on an adopted device (e.g. restart)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ExecuteDeviceAction) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"deviceId": map[string]interface{}{
				"type":        "string",
				"description": "Device UUID",
			},
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action to perform on the device",
			},
		},
		"required": []string{"deviceId", "action"},
	}
}

// Execute runs the tool.
func (t *ExecuteDeviceAction) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID   string `json:"siteId"`
		DeviceID string `json:"deviceId"`
		Action   string `json:"action"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	siteID, err := resolveSiteID(params.SiteID, t.defaultSiteID)
	if err != nil {
		return "", err
	}

	deviceID, err := resolveUUID("deviceId", params.DeviceID)
	if err != nil {
		return "", err
	}

	if params.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	resp, err := t.client.ExecuteAdoptedDeviceActionWithResponse(
		ctx,
		siteID,
		deviceID,
		unifi.DeviceActionRequest{
			Action: params.Action,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute device action: %w",
			err,
		)
	}

	if resp.StatusCode() != 200 {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return fmt.Sprintf(
		"Action %q executed on device %s.",
		params.Action,
		params.DeviceID,
	), nil
}

// --- execute_port_action ---

// ExecutePortAction implements the execute_port_action MCP tool.
type ExecutePortAction struct {
	baseTool
}

// NewExecutePortAction creates a new ExecutePortAction tool.
func NewExecutePortAction(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ExecutePortAction {
	return &ExecutePortAction{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *ExecutePortAction) Description() string {
	return "Execute an action on a port of an adopted device"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ExecutePortAction) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"deviceId": map[string]interface{}{
				"type":        "string",
				"description": "Device UUID",
			},
			"portIdx": map[string]interface{}{
				"type":        "integer",
				"description": "Port index (integer, not UUID)",
			},
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action to perform on the port",
			},
		},
		"required": []string{"deviceId", "portIdx", "action"},
	}
}

// Execute runs the tool.
func (t *ExecutePortAction) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID   string `json:"siteId"`
		DeviceID string `json:"deviceId"`
		PortIdx  *int32 `json:"portIdx"`
		Action   string `json:"action"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	siteID, err := resolveSiteID(params.SiteID, t.defaultSiteID)
	if err != nil {
		return "", err
	}

	deviceID, err := resolveUUID("deviceId", params.DeviceID)
	if err != nil {
		return "", err
	}

	if params.PortIdx == nil {
		return "", fmt.Errorf("portIdx is required")
	}

	if params.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	resp, err := t.client.ExecutePortActionWithResponse(
		ctx,
		siteID,
		deviceID,
		*params.PortIdx,
		unifi.PortActionRequest{
			Action: params.Action,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute port action: %w",
			err,
		)
	}

	if resp.StatusCode() != 200 {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return fmt.Sprintf(
		"Action %q executed on port %d of device %s.",
		params.Action,
		*params.PortIdx,
		params.DeviceID,
	), nil
}

// --- get_device_statistics ---

// GetDeviceStatistics implements the get_device_statistics MCP tool.
type GetDeviceStatistics struct {
	baseTool
}

// NewGetDeviceStatistics creates a new GetDeviceStatistics tool.
func NewGetDeviceStatistics(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetDeviceStatistics {
	return &GetDeviceStatistics{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *GetDeviceStatistics) Description() string {
	return "Get latest statistics for an adopted device"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetDeviceStatistics) InputSchema() map[string]interface{} {
	return siteAndIDSchema("deviceId", "Device UUID")
}

// Execute runs the tool.
func (t *GetDeviceStatistics) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID   string `json:"siteId"`
		DeviceID string `json:"deviceId"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	siteID, err := resolveSiteID(params.SiteID, t.defaultSiteID)
	if err != nil {
		return "", err
	}

	deviceID, err := resolveUUID("deviceId", params.DeviceID)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetAdoptedDeviceLatestStatisticsWithResponse(
		ctx,
		siteID,
		deviceID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get device statistics: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatDeviceStatistics(resp.JSON200), nil
}

func formatDeviceStatistics(
	s *unifi.LatestStatisticsForADevice,
) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Device Statistics:\n")
	if s.LastHeartbeatAt != nil {
		fmt.Fprintf(
			&b,
			"Last Heartbeat: %s\n",
			s.LastHeartbeatAt.String(),
		)
	}
	if s.NextHeartbeatAt != nil {
		fmt.Fprintf(
			&b,
			"Next Heartbeat: %s\n",
			s.NextHeartbeatAt.String(),
		)
	}
	if s.UptimeSec != nil {
		fmt.Fprintf(&b, "Uptime: %d seconds\n", *s.UptimeSec)
	}
	if s.CpuUtilizationPct != nil {
		fmt.Fprintf(
			&b,
			"CPU Utilization: %.1f%%\n",
			*s.CpuUtilizationPct,
		)
	}
	if s.MemoryUtilizationPct != nil {
		fmt.Fprintf(
			&b,
			"Memory Utilization: %.1f%%\n",
			*s.MemoryUtilizationPct,
		)
	}
	if s.LoadAverage1Min != nil {
		fmt.Fprintf(
			&b,
			"Load Average (1m/5m/15m): %.2f / %.2f / %.2f\n",
			*s.LoadAverage1Min,
			ptrFloat64OrZero(s.LoadAverage5Min),
			ptrFloat64OrZero(s.LoadAverage15Min),
		)
	}
	if s.Uplink != nil {
		if s.Uplink.RxRateBps != nil {
			fmt.Fprintf(
				&b,
				"Uplink RX: %d bps\n",
				*s.Uplink.RxRateBps,
			)
		}
		if s.Uplink.TxRateBps != nil {
			fmt.Fprintf(
				&b,
				"Uplink TX: %d bps\n",
				*s.Uplink.TxRateBps,
			)
		}
	}
	return b.String()
}

func ptrFloat64OrZero(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

// --- list_pending_devices ---

// ListPendingDevices implements the list_pending_devices MCP tool.
// This tool is not site-scoped.
type ListPendingDevices struct {
	baseTool
}

// NewListPendingDevices creates a new ListPendingDevices tool.
func NewListPendingDevices(
	c *unifi.ClientWithResponses,
) *ListPendingDevices {
	return &ListPendingDevices{baseTool{client: c}}
}

// Description returns a description of the tool.
func (t *ListPendingDevices) Description() string {
	return "List devices pending adoption (not site-scoped)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListPendingDevices) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": paginationSchema(),
	}
}

// Execute runs the tool.
func (t *ListPendingDevices) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		Limit  *int32 `json:"limit"`
		Offset *int32 `json:"offset"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	resp, err := t.client.GetPendingDevicePageWithResponse(
		ctx,
		&unifi.GetPendingDevicePageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list pending devices: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	page := resp.JSON200
	if len(page.Data) == 0 {
		return "No pending devices found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Pending Devices (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, d := range page.Data {
		fw := ""
		if d.FirmwareVersion != nil {
			fw = *d.FirmwareVersion
		}
		fmt.Fprintf(
			&b,
			"%d. Model: %s\n   MAC: %s\n   IP: %s\n"+
				"   State: %s\n   Firmware: %s\n",
			i+1,
			d.Model,
			d.MacAddress,
			d.IpAddress,
			d.State,
			fw,
		)
	}

	return b.String(), nil
}
