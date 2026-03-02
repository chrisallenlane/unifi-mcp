package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp-server/internal/unifi"
)

// --- list_wifi_broadcasts ---

// ListWiFiBroadcasts implements the list_wifi_broadcasts
// MCP tool.
type ListWiFiBroadcasts struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListWiFiBroadcasts creates a new ListWiFiBroadcasts
// tool.
func NewListWiFiBroadcasts(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListWiFiBroadcasts {
	return &ListWiFiBroadcasts{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListWiFiBroadcasts) Description() string {
	return "List WiFi broadcasts (SSIDs) for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListWiFiBroadcasts) InputSchema() map[string]interface{} {
	props := map[string]interface{}{
		"siteId": siteIDSchema(),
	}
	for k, v := range paginationSchema() {
		props[k] = v
	}
	return map[string]interface{}{
		"type":       "object",
		"properties": props,
	}
}

// Execute runs the tool.
func (t *ListWiFiBroadcasts) Execute(
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

	siteID, err := resolveSiteID(
		params.SiteID,
		t.defaultSiteID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetWifiBroadcastPageWithResponse(
		ctx,
		siteID,
		&unifi.GetWifiBroadcastPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list WiFi broadcasts: %w",
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
		return "No WiFi broadcasts found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"WiFi Broadcasts (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, w := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Type: %s\n   Enabled: %t\n   Security: %s\n",
			i+1,
			w.Name,
			w.Id.String(),
			w.Type,
			w.Enabled,
			w.SecurityConfiguration.Type,
		)
	}

	return b.String(), nil
}

// --- get_wifi_broadcast ---

// GetWiFiBroadcast implements the get_wifi_broadcast
// MCP tool.
type GetWiFiBroadcast struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewGetWiFiBroadcast creates a new GetWiFiBroadcast tool.
func NewGetWiFiBroadcast(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetWiFiBroadcast {
	return &GetWiFiBroadcast{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *GetWiFiBroadcast) Description() string {
	return "Get details of a specific WiFi broadcast (SSID)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetWiFiBroadcast) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"wifiBroadcastId": map[string]interface{}{
				"type":        "string",
				"description": "WiFi broadcast UUID",
			},
		},
		"required": []string{"wifiBroadcastId"},
	}
}

// Execute runs the tool.
func (t *GetWiFiBroadcast) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID          string `json:"siteId"`
		WiFiBroadcastID string `json:"wifiBroadcastId"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	siteID, err := resolveSiteID(
		params.SiteID,
		t.defaultSiteID,
	)
	if err != nil {
		return "", err
	}

	broadcastID, err := resolveUUID(
		"wifiBroadcastId",
		params.WiFiBroadcastID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetWifiBroadcastDetailsWithResponse(
		ctx,
		siteID,
		broadcastID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get WiFi broadcast: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatWiFiBroadcast(resp.JSON200), nil
}

func formatWiFiBroadcast(
	w *unifi.WifiBroadcastDetails,
) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\n", w.Name)
	fmt.Fprintf(&b, "ID: %s\n", w.Id.String())
	fmt.Fprintf(&b, "Type: %s\n", w.Type)
	fmt.Fprintf(&b, "Enabled: %t\n", w.Enabled)
	fmt.Fprintf(&b, "Hidden SSID: %t\n", w.HideName)
	fmt.Fprintf(
		&b,
		"Security: %s\n",
		w.SecurityConfiguration.Type,
	)
	fmt.Fprintf(
		&b,
		"Client Isolation: %t\n",
		w.ClientIsolationEnabled,
	)
	fmt.Fprintf(
		&b,
		"Multicast to Unicast: %t\n",
		w.MulticastToUnicastConversionEnabled,
	)
	fmt.Fprintf(
		&b,
		"U-APSD: %t\n",
		w.UapsdEnabled,
	)
	if w.Network != nil {
		fmt.Fprintf(
			&b,
			"Network Type: %s\n",
			w.Network.Type,
		)
	}
	if w.SecurityConfiguration.RadiusConfiguration != nil {
		radius, _ := json.MarshalIndent(
			w.SecurityConfiguration.RadiusConfiguration,
			"",
			"  ",
		)
		fmt.Fprintf(
			&b,
			"RADIUS Configuration:\n%s\n",
			radius,
		)
	}
	if w.BroadcastingDeviceFilter != nil {
		fmt.Fprintf(
			&b,
			"Broadcasting Device Filter Type: %s\n",
			w.BroadcastingDeviceFilter.Type,
		)
	}
	if w.ClientFilteringPolicy != nil {
		fmt.Fprintf(
			&b,
			"Client Filtering: %s (%d MACs)\n",
			w.ClientFilteringPolicy.Action,
			len(w.ClientFilteringPolicy.MacAddressFilter),
		)
	}
	return b.String()
}

// --- create_wifi_broadcast ---

// CreateWiFiBroadcast implements the create_wifi_broadcast
// MCP tool.
type CreateWiFiBroadcast struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewCreateWiFiBroadcast creates a new CreateWiFiBroadcast
// tool.
func NewCreateWiFiBroadcast(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *CreateWiFiBroadcast {
	return &CreateWiFiBroadcast{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *CreateWiFiBroadcast) Description() string {
	return "Create a new WiFi broadcast (SSID)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *CreateWiFiBroadcast) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"name": map[string]interface{}{
				"type":        "string",
				"description": "SSID name",
			},
			"enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the broadcast is enabled",
			},
			"type": map[string]interface{}{
				"type":        "string",
				"description": "Broadcast type",
				"enum": []string{
					"STANDARD",
					"IOT_OPTIMIZED",
				},
			},
			"securityConfiguration": map[string]interface{}{
				"type":        "object",
				"description": "Security settings (type: OPEN, WPA2, WPA3, WPA2_WPA3, WPA2_ENTERPRISE, WPA3_ENTERPRISE, WPA2_WPA3_ENTERPRISE)",
			},
			"hideName": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether to hide the SSID",
			},
			"clientIsolationEnabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether client isolation is enabled",
			},
			"network": map[string]interface{}{
				"type":        "object",
				"description": "Network reference (optional)",
			},
			"broadcastingDeviceFilter": map[string]interface{}{
				"type":        "object",
				"description": "Device filter for broadcasting (optional)",
			},
		},
		"required": []string{
			"name",
			"enabled",
			"type",
			"securityConfiguration",
		},
	}
}

// Execute runs the tool.
func (t *CreateWiFiBroadcast) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID string `json:"siteId"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	siteID, err := resolveSiteID(
		params.SiteID,
		t.defaultSiteID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.CreateWifiBroadcastWithBodyWithResponse(
		ctx,
		siteID,
		"application/json",
		bytes.NewReader(args),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create WiFi broadcast: %w",
			err,
		)
	}

	if resp.JSON201 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return fmt.Sprintf(
		"WiFi broadcast created:\n%s",
		formatWiFiBroadcast(resp.JSON201),
	), nil
}

// --- update_wifi_broadcast ---

// UpdateWiFiBroadcast implements the update_wifi_broadcast
// MCP tool.
type UpdateWiFiBroadcast struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewUpdateWiFiBroadcast creates a new UpdateWiFiBroadcast
// tool.
func NewUpdateWiFiBroadcast(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *UpdateWiFiBroadcast {
	return &UpdateWiFiBroadcast{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *UpdateWiFiBroadcast) Description() string {
	return "Update an existing WiFi broadcast (SSID)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *UpdateWiFiBroadcast) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"wifiBroadcastId": map[string]interface{}{
				"type":        "string",
				"description": "WiFi broadcast UUID",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "SSID name",
			},
			"enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the broadcast is enabled",
			},
			"type": map[string]interface{}{
				"type":        "string",
				"description": "Broadcast type",
				"enum": []string{
					"STANDARD",
					"IOT_OPTIMIZED",
				},
			},
			"securityConfiguration": map[string]interface{}{
				"type":        "object",
				"description": "Security settings",
			},
			"hideName": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether to hide the SSID",
			},
			"clientIsolationEnabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether client isolation is enabled",
			},
		},
		"required": []string{
			"wifiBroadcastId",
			"name",
			"enabled",
			"type",
			"securityConfiguration",
		},
	}
}

// Execute runs the tool.
func (t *UpdateWiFiBroadcast) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID          string `json:"siteId"`
		WiFiBroadcastID string `json:"wifiBroadcastId"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	siteID, err := resolveSiteID(
		params.SiteID,
		t.defaultSiteID,
	)
	if err != nil {
		return "", err
	}

	broadcastID, err := resolveUUID(
		"wifiBroadcastId",
		params.WiFiBroadcastID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.UpdateWifiBroadcastWithBodyWithResponse(
		ctx,
		siteID,
		broadcastID,
		"application/json",
		bytes.NewReader(args),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update WiFi broadcast: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return fmt.Sprintf(
		"WiFi broadcast updated:\n%s",
		formatWiFiBroadcast(resp.JSON200),
	), nil
}

// --- delete_wifi_broadcast ---

// DeleteWiFiBroadcast implements the delete_wifi_broadcast
// MCP tool.
type DeleteWiFiBroadcast struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewDeleteWiFiBroadcast creates a new DeleteWiFiBroadcast
// tool.
func NewDeleteWiFiBroadcast(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *DeleteWiFiBroadcast {
	return &DeleteWiFiBroadcast{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *DeleteWiFiBroadcast) Description() string {
	return "Delete a WiFi broadcast (SSID)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *DeleteWiFiBroadcast) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"wifiBroadcastId": map[string]interface{}{
				"type":        "string",
				"description": "WiFi broadcast UUID",
			},
			"force": map[string]interface{}{
				"type":        "boolean",
				"description": "Force deletion",
			},
		},
		"required": []string{"wifiBroadcastId"},
	}
}

// Execute runs the tool.
func (t *DeleteWiFiBroadcast) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID          string `json:"siteId"`
		WiFiBroadcastID string `json:"wifiBroadcastId"`
		Force           *bool  `json:"force"`
	}
	if len(args) > 0 {
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf(
				"failed to parse arguments: %w",
				err,
			)
		}
	}

	siteID, err := resolveSiteID(
		params.SiteID,
		t.defaultSiteID,
	)
	if err != nil {
		return "", err
	}

	broadcastID, err := resolveUUID(
		"wifiBroadcastId",
		params.WiFiBroadcastID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.DeleteWifiBroadcastWithResponse(
		ctx,
		siteID,
		broadcastID,
		&unifi.DeleteWifiBroadcastParams{
			Force: params.Force,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete WiFi broadcast: %w",
			err,
		)
	}

	if resp.StatusCode() != 200 {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return "WiFi broadcast deleted.", nil
}
