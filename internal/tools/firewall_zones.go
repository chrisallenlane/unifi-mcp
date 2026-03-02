package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp-server/internal/unifi"
)

// formatZone formats a single firewall zone for display.
func formatZone(zone *unifi.FirewallZone) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\n", zone.Name)
	fmt.Fprintf(&b, "ID: %s\n", zone.Id.String())
	fmt.Fprintf(&b, "Origin: %s\n", zone.Metadata.Origin)

	if len(zone.NetworkIds) > 0 {
		ids := make([]string, len(zone.NetworkIds))
		for i, id := range zone.NetworkIds {
			ids[i] = id.String()
		}
		fmt.Fprintf(
			&b,
			"Network IDs: %s\n",
			strings.Join(ids, ", "),
		)
	} else {
		fmt.Fprintln(&b, "Network IDs: (none)")
	}

	return b.String()
}

// ListFirewallZones implements the list_firewall_zones MCP tool.
type ListFirewallZones struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListFirewallZones creates a new ListFirewallZones tool.
func NewListFirewallZones(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListFirewallZones {
	return &ListFirewallZones{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListFirewallZones) Description() string {
	return "List firewall zones for a UniFi site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListFirewallZones) InputSchema() map[string]interface{} {
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
func (t *ListFirewallZones) Execute(
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

	resp, err := t.client.GetFirewallZonesWithResponse(
		ctx,
		siteID,
		&unifi.GetFirewallZonesParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list firewall zones: %w",
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
		return "No firewall zones found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Firewall Zones (%d of %d):\n\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, zone := range page.Data {
		fmt.Fprintf(&b, "%d. %s", i+1, formatZone(&zone))
		if i < len(page.Data)-1 {
			fmt.Fprintln(&b)
		}
	}

	return b.String(), nil
}

// GetFirewallZone implements the get_firewall_zone MCP tool.
type GetFirewallZone struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewGetFirewallZone creates a new GetFirewallZone tool.
func NewGetFirewallZone(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetFirewallZone {
	return &GetFirewallZone{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *GetFirewallZone) Description() string {
	return "Get details of a specific firewall zone"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetFirewallZone) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"firewallZoneId": map[string]interface{}{
				"type":        "string",
				"description": "Firewall zone UUID",
			},
		},
		"required": []string{"firewallZoneId"},
	}
}

// Execute runs the tool.
func (t *GetFirewallZone) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID         string `json:"siteId"`
		FirewallZoneID string `json:"firewallZoneId"`
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

	zoneID, err := resolveUUID(
		"firewallZoneId",
		params.FirewallZoneID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetFirewallZoneWithResponse(
		ctx,
		siteID,
		zoneID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get firewall zone: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatZone(resp.JSON200), nil
}

// CreateFirewallZone implements the create_firewall_zone MCP tool.
type CreateFirewallZone struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewCreateFirewallZone creates a new CreateFirewallZone tool.
func NewCreateFirewallZone(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *CreateFirewallZone {
	return &CreateFirewallZone{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *CreateFirewallZone) Description() string {
	return "Create a new firewall zone"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *CreateFirewallZone) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the firewall zone",
			},
			"networkIds": map[string]interface{}{
				"type":        "array",
				"description": "List of network UUIDs to include in the zone",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"required": []string{"name", "networkIds"},
	}
}

// Execute runs the tool.
func (t *CreateFirewallZone) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID     string   `json:"siteId"`
		Name       string   `json:"name"`
		NetworkIDs []string `json:"networkIds"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	if params.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	siteID, err := resolveSiteID(
		params.SiteID,
		t.defaultSiteID,
	)
	if err != nil {
		return "", err
	}

	networkIDs, err := resolveUUIDs(
		"networkIds",
		params.NetworkIDs,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.CreateFirewallZoneWithResponse(
		ctx,
		siteID,
		unifi.CreateFirewallZoneJSONRequestBody{
			Name:       params.Name,
			NetworkIds: networkIDs,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create firewall zone: %w",
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
		"Firewall zone created:\n%s",
		formatZone(resp.JSON201),
	), nil
}

// UpdateFirewallZone implements the update_firewall_zone MCP tool.
type UpdateFirewallZone struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewUpdateFirewallZone creates a new UpdateFirewallZone tool.
func NewUpdateFirewallZone(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *UpdateFirewallZone {
	return &UpdateFirewallZone{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *UpdateFirewallZone) Description() string {
	return "Update an existing firewall zone"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *UpdateFirewallZone) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"firewallZoneId": map[string]interface{}{
				"type":        "string",
				"description": "Firewall zone UUID to update",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "New name for the firewall zone",
			},
			"networkIds": map[string]interface{}{
				"type":        "array",
				"description": "New list of network UUIDs",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"required": []string{
			"firewallZoneId",
			"name",
			"networkIds",
		},
	}
}

// Execute runs the tool.
func (t *UpdateFirewallZone) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID         string   `json:"siteId"`
		FirewallZoneID string   `json:"firewallZoneId"`
		Name           string   `json:"name"`
		NetworkIDs     []string `json:"networkIds"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	if params.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	siteID, err := resolveSiteID(
		params.SiteID,
		t.defaultSiteID,
	)
	if err != nil {
		return "", err
	}

	zoneID, err := resolveUUID(
		"firewallZoneId",
		params.FirewallZoneID,
	)
	if err != nil {
		return "", err
	}

	networkIDs, err := resolveUUIDs(
		"networkIds",
		params.NetworkIDs,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.UpdateFirewallZoneWithResponse(
		ctx,
		siteID,
		zoneID,
		unifi.UpdateFirewallZoneJSONRequestBody{
			Name:       params.Name,
			NetworkIds: networkIDs,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update firewall zone: %w",
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
		"Firewall zone updated:\n%s",
		formatZone(resp.JSON200),
	), nil
}

// DeleteFirewallZone implements the delete_firewall_zone MCP tool.
type DeleteFirewallZone struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewDeleteFirewallZone creates a new DeleteFirewallZone tool.
func NewDeleteFirewallZone(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *DeleteFirewallZone {
	return &DeleteFirewallZone{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *DeleteFirewallZone) Description() string {
	return "Delete a user-defined firewall zone"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *DeleteFirewallZone) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"firewallZoneId": map[string]interface{}{
				"type":        "string",
				"description": "Firewall zone UUID to delete",
			},
		},
		"required": []string{"firewallZoneId"},
	}
}

// Execute runs the tool.
func (t *DeleteFirewallZone) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID         string `json:"siteId"`
		FirewallZoneID string `json:"firewallZoneId"`
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

	zoneID, err := resolveUUID(
		"firewallZoneId",
		params.FirewallZoneID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.DeleteFirewallZoneWithResponse(
		ctx,
		siteID,
		zoneID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete firewall zone: %w",
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
		"Firewall zone %s deleted successfully.",
		zoneID.String(),
	), nil
}
