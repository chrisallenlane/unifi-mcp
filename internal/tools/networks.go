package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp-server/internal/unifi"
)

// --- list_networks ---

// ListNetworks implements the list_networks MCP tool.
type ListNetworks struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListNetworks creates a new ListNetworks tool.
func NewListNetworks(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListNetworks {
	return &ListNetworks{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListNetworks) Description() string {
	return "List networks for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListNetworks) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"limit": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum number of networks to return",
			},
			"offset": map[string]interface{}{
				"type":        "integer",
				"description": "Number of networks to skip",
			},
		},
	}
}

// Execute runs the tool.
func (t *ListNetworks) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID string `json:"siteId"`
		Limit  *int32 `json:"limit"`
		Offset *int32 `json:"offset"`
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

	resp, err := t.client.GetNetworksOverviewPageWithResponse(
		ctx,
		siteID,
		&unifi.GetNetworksOverviewPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list networks: %w",
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
		return "No networks found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Networks (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, n := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   VLAN: %d\n   Management: %s\n   Enabled: %t\n   Default: %t\n",
			i+1,
			n.Name,
			n.Id.String(),
			n.VlanId,
			n.Management,
			n.Enabled,
			n.Default,
		)
	}

	return b.String(), nil
}

// --- get_network ---

// GetNetwork implements the get_network MCP tool.
type GetNetwork struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewGetNetwork creates a new GetNetwork tool.
func NewGetNetwork(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetNetwork {
	return &GetNetwork{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *GetNetwork) Description() string {
	return "Get details of a specific network"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetNetwork) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"networkId": map[string]interface{}{
				"type":        "string",
				"description": "Network UUID",
			},
		},
		"required": []string{"networkId"},
	}
}

// Execute runs the tool.
func (t *GetNetwork) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID    string `json:"siteId"`
		NetworkID string `json:"networkId"`
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

	networkID, err := resolveUUID(
		"networkId",
		params.NetworkID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetNetworkDetailsWithResponse(
		ctx,
		siteID,
		networkID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get network: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatNetworkDetails(resp.JSON200), nil
}

func formatNetworkDetails(n *unifi.NetworkDetails) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\n", n.Name)
	fmt.Fprintf(&b, "ID: %s\n", n.Id.String())
	fmt.Fprintf(&b, "VLAN ID: %d\n", n.VlanId)
	fmt.Fprintf(&b, "Management: %s\n", n.Management)
	fmt.Fprintf(&b, "Enabled: %t\n", n.Enabled)
	fmt.Fprintf(&b, "Default: %t\n", n.Default)
	if n.DhcpGuarding != nil {
		fmt.Fprintf(&b, "DHCP Guarding:\n")
		if len(n.DhcpGuarding.TrustedDhcpServerIpAddresses) > 0 {
			fmt.Fprintf(
				&b,
				"  Trusted Servers: %s\n",
				strings.Join(
					n.DhcpGuarding.TrustedDhcpServerIpAddresses,
					", ",
				),
			)
		} else {
			fmt.Fprintf(
				&b,
				"  Trusted Servers: (none)\n",
			)
		}
	}
	return b.String()
}

// --- create_network ---

// CreateNetwork implements the create_network MCP tool.
type CreateNetwork struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewCreateNetwork creates a new CreateNetwork tool.
func NewCreateNetwork(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *CreateNetwork {
	return &CreateNetwork{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *CreateNetwork) Description() string {
	return "Create a new network (VLAN, subnet, DHCP settings)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *CreateNetwork) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Network name",
			},
			"enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the network is enabled",
			},
			"management": map[string]interface{}{
				"type":        "string",
				"description": "Network management type",
				"enum": []string{
					"GATEWAY",
					"SWITCH",
					"UNMANAGED",
				},
			},
			"vlanId": map[string]interface{}{
				"type":        "integer",
				"description": "VLAN ID (1 for default, >= 2 for additional)",
			},
			"dhcpGuarding": map[string]interface{}{
				"type":        "object",
				"description": "DHCP guarding settings (optional)",
				"properties": map[string]interface{}{
					"trustedDhcpServerIpAddresses": map[string]interface{}{
						"type":        "array",
						"description": "List of trusted DHCP server IP addresses",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
		"required": []string{
			"name",
			"enabled",
			"management",
			"vlanId",
		},
	}
}

// Execute runs the tool.
func (t *CreateNetwork) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID string `json:"siteId"`
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

	resp, err := t.client.CreateNetworkWithBodyWithResponse(
		ctx,
		siteID,
		"application/json",
		bytes.NewReader(args),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create network: %w",
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
		"Network created:\n%s",
		formatNetworkDetails(resp.JSON201),
	), nil
}

// --- update_network ---

// UpdateNetwork implements the update_network MCP tool.
type UpdateNetwork struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewUpdateNetwork creates a new UpdateNetwork tool.
func NewUpdateNetwork(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *UpdateNetwork {
	return &UpdateNetwork{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *UpdateNetwork) Description() string {
	return "Update an existing network"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *UpdateNetwork) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"networkId": map[string]interface{}{
				"type":        "string",
				"description": "Network UUID",
			},
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Network name",
			},
			"enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the network is enabled",
			},
			"management": map[string]interface{}{
				"type":        "string",
				"description": "Network management type",
				"enum": []string{
					"GATEWAY",
					"SWITCH",
					"UNMANAGED",
				},
			},
			"vlanId": map[string]interface{}{
				"type":        "integer",
				"description": "VLAN ID",
			},
			"dhcpGuarding": map[string]interface{}{
				"type":        "object",
				"description": "DHCP guarding settings (optional)",
				"properties": map[string]interface{}{
					"trustedDhcpServerIpAddresses": map[string]interface{}{
						"type":        "array",
						"description": "List of trusted DHCP server IP addresses",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
		"required": []string{
			"networkId",
			"name",
			"enabled",
			"management",
			"vlanId",
		},
	}
}

// Execute runs the tool.
func (t *UpdateNetwork) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID    string `json:"siteId"`
		NetworkID string `json:"networkId"`
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

	networkID, err := resolveUUID(
		"networkId",
		params.NetworkID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.UpdateNetworkWithBodyWithResponse(
		ctx,
		siteID,
		networkID,
		"application/json",
		bytes.NewReader(args),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update network: %w",
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
		"Network updated:\n%s",
		formatNetworkDetails(resp.JSON200),
	), nil
}

// --- delete_network ---

// DeleteNetwork implements the delete_network MCP tool.
type DeleteNetwork struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewDeleteNetwork creates a new DeleteNetwork tool.
func NewDeleteNetwork(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *DeleteNetwork {
	return &DeleteNetwork{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *DeleteNetwork) Description() string {
	return "Delete a network"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *DeleteNetwork) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"networkId": map[string]interface{}{
				"type":        "string",
				"description": "Network UUID",
			},
			"force": map[string]interface{}{
				"type":        "boolean",
				"description": "Force deletion even if resources reference this network",
			},
		},
		"required": []string{"networkId"},
	}
}

// Execute runs the tool.
func (t *DeleteNetwork) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID    string `json:"siteId"`
		NetworkID string `json:"networkId"`
		Force     *bool  `json:"force"`
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

	networkID, err := resolveUUID(
		"networkId",
		params.NetworkID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.DeleteNetworkWithResponse(
		ctx,
		siteID,
		networkID,
		&unifi.DeleteNetworkParams{
			Force: params.Force,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete network: %w",
			err,
		)
	}

	if resp.StatusCode() != 200 {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return "Network deleted.", nil
}

// --- get_network_references ---

// GetNetworkReferences implements the get_network_references
// MCP tool.
type GetNetworkReferences struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewGetNetworkReferences creates a new GetNetworkReferences
// tool.
func NewGetNetworkReferences(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetNetworkReferences {
	return &GetNetworkReferences{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *GetNetworkReferences) Description() string {
	return "Get resources that reference a network (useful before deleting)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetNetworkReferences) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"networkId": map[string]interface{}{
				"type":        "string",
				"description": "Network UUID",
			},
		},
		"required": []string{"networkId"},
	}
}

// Execute runs the tool.
func (t *GetNetworkReferences) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID    string `json:"siteId"`
		NetworkID string `json:"networkId"`
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

	networkID, err := resolveUUID(
		"networkId",
		params.NetworkID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetNetworkReferencesWithResponse(
		ctx,
		siteID,
		networkID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get network references: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	refs := resp.JSON200
	if len(refs.ReferenceResources) == 0 {
		return "No resources reference this network.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Network references (%d resource types):\n",
		len(refs.ReferenceResources),
	)
	for _, res := range refs.ReferenceResources {
		fmt.Fprintf(
			&b,
			"- %s: %d reference(s)\n",
			res.ResourceType,
			res.ReferenceCount,
		)
		if res.References != nil {
			for _, ref := range *res.References {
				fmt.Fprintf(
					&b,
					"  - %s\n",
					ref.ReferenceId.String(),
				)
			}
		}
	}

	return b.String(), nil
}
