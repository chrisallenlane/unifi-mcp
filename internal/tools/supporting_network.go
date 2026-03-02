package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

// --- list_wans ---

// ListWans implements the list_wans MCP tool.
type ListWans struct {
	baseTool
}

// NewListWans creates a new ListWans tool.
func NewListWans(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListWans {
	return &ListWans{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *ListWans) Description() string {
	return "List WAN interfaces for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListWans) InputSchema() map[string]interface{} {
	return listSchema()
}

// Execute runs the tool.
func (t *ListWans) Execute(
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

	resp, err := t.client.GetWansOverviewPageWithResponse(
		ctx,
		siteID,
		&unifi.GetWansOverviewPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to list WANs: %w", err)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	page := resp.JSON200
	if len(page.Data) == 0 {
		return "No WANs found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"WANs (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, wan := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n",
			i+1,
			wan.Name,
			wan.Id.String(),
		)
	}

	return b.String(), nil
}

// --- list_vpn_tunnels ---

// ListVpnTunnels implements the list_vpn_tunnels MCP tool.
type ListVpnTunnels struct {
	baseTool
}

// NewListVpnTunnels creates a new ListVpnTunnels tool.
func NewListVpnTunnels(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListVpnTunnels {
	return &ListVpnTunnels{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *ListVpnTunnels) Description() string {
	return "List site-to-site VPN tunnels for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListVpnTunnels) InputSchema() map[string]interface{} {
	return listSchema()
}

// Execute runs the tool.
func (t *ListVpnTunnels) Execute(
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

	resp, err := t.client.GetSiteToSiteVpnTunnelPageWithResponse(
		ctx,
		siteID,
		&unifi.GetSiteToSiteVpnTunnelPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list VPN tunnels: %w",
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
		return "No VPN tunnels found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"VPN Tunnels (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, tunnel := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Type: %s\n",
			i+1,
			tunnel.Name,
			tunnel.Id.String(),
			tunnel.Type,
		)
	}

	return b.String(), nil
}

// --- list_vpn_servers ---

// ListVpnServers implements the list_vpn_servers MCP tool.
type ListVpnServers struct {
	baseTool
}

// NewListVpnServers creates a new ListVpnServers tool.
func NewListVpnServers(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListVpnServers {
	return &ListVpnServers{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *ListVpnServers) Description() string {
	return "List VPN servers for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListVpnServers) InputSchema() map[string]interface{} {
	return listSchema()
}

// Execute runs the tool.
func (t *ListVpnServers) Execute(
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

	resp, err := t.client.GetVpnServerPageWithResponse(
		ctx,
		siteID,
		&unifi.GetVpnServerPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list VPN servers: %w",
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
		return "No VPN servers found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"VPN Servers (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, server := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Type: %s\n   Enabled: %t\n",
			i+1,
			server.Name,
			server.Id.String(),
			server.Type,
			server.Enabled,
		)
	}

	return b.String(), nil
}
