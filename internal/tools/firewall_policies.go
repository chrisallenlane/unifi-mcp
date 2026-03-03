package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

// formatPolicy formats a single firewall policy for display.
func formatPolicy(p *unifi.FirewallPolicy) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\n", p.Name)
	fmt.Fprintf(&b, "ID: %s\n", p.Id.String())
	fmt.Fprintf(&b, "Enabled: %t\n", p.Enabled)
	fmt.Fprintf(&b, "Action: %s\n", p.Action.Type)
	fmt.Fprintf(
		&b,
		"Source Zone: %s\n",
		p.Source.ZoneId.String(),
	)
	fmt.Fprintf(
		&b,
		"Destination Zone: %s\n",
		p.Destination.ZoneId.String(),
	)
	fmt.Fprintf(
		&b,
		"IP Protocol Scope: %s\n",
		p.IpProtocolScope.IpVersion,
	)
	fmt.Fprintf(&b, "Index: %d\n", p.Index)
	fmt.Fprintf(&b, "Logging: %t\n", p.LoggingEnabled)
	fmt.Fprintf(&b, "Origin: %s\n", p.Metadata.Origin)

	if p.Description != nil && *p.Description != "" {
		fmt.Fprintf(&b, "Description: %s\n", *p.Description)
	}
	if p.ConnectionStateFilter != nil &&
		len(*p.ConnectionStateFilter) > 0 {
		states := make(
			[]string,
			len(*p.ConnectionStateFilter),
		)
		for i, s := range *p.ConnectionStateFilter {
			states[i] = string(s)
		}
		fmt.Fprintf(
			&b,
			"Connection State Filter: %s\n",
			strings.Join(states, ", "),
		)
	}
	if p.IpsecFilter != nil {
		fmt.Fprintf(
			&b,
			"IPsec Filter: %s\n",
			string(*p.IpsecFilter),
		)
	}
	if p.Schedule != nil {
		fmt.Fprintf(&b, "Schedule Mode: %s\n", p.Schedule.Mode)
	}
	if p.Source.TrafficFilter != nil {
		fmt.Fprintf(
			&b,
			"Source Traffic Filter: %s\n",
			p.Source.TrafficFilter.Type,
		)
	}
	if p.Destination.TrafficFilter != nil {
		fmt.Fprintf(
			&b,
			"Destination Traffic Filter: %s\n",
			p.Destination.TrafficFilter.Type,
		)
	}

	return b.String()
}

// trafficFilterSchema returns the JSON schema for a traffic filter object on
// a firewall policy source or destination endpoint. The allowedTypes slice
// lists the valid values for the discriminating "type" field; the valid nested
// sub-objects depend on the type chosen:
//
//   - IP_ADDRESS   → ipAddressFilter  (type, matchOpposite, items)
//   - DOMAIN       → domainFilter     (type, domains)
//   - MAC_ADDRESS  → macAddressFilter (macAddresses)          [source only]
//   - NETWORK      → networkFilter    (matchOpposite, networkIds)
//   - PORT         → portFilter       (portRanges)
func trafficFilterSchema(
	direction string,
	allowedTypes []string,
) map[string]any {
	return map[string]any{
		"type": "object",
		"description": fmt.Sprintf(
			"Optional %s traffic filter. Set \"type\" to one of the"+
				" allowed values, then include the corresponding"+
				" nested filter object (e.g. ipAddressFilter when"+
				" type is IP_ADDRESS).",
			direction,
		),
		"properties": map[string]any{
			"type": map[string]any{
				"type":        "string",
				"description": "Traffic filter discriminator",
				"enum":        allowedTypes,
			},
			"ipAddressFilter": map[string]any{
				"type":        "object",
				"description": "Required when type is IP_ADDRESS",
				"properties": map[string]any{
					"type": map[string]any{
						"type": "string",
						"enum": []string{
							"IP_ADDRESSES",
							"TRAFFIC_MATCHING_LIST",
						},
					},
					"matchOpposite": map[string]any{
						"type": "boolean",
						"description": "Match all addresses except" +
							" the specified ones",
					},
					"items": map[string]any{
						"type":        "array",
						"description": "IP addresses, ranges, or subnets",
						"items": map[string]any{
							"type": "object",
						},
					},
					"trafficMatchingListId": map[string]any{
						"type":        "string",
						"description": "UUID of a Traffic Matching List",
					},
				},
				"required": []string{"type", "matchOpposite"},
			},
			"domainFilter": map[string]any{
				"type":        "object",
				"description": "Required when type is DOMAIN",
				"properties": map[string]any{
					"type": map[string]any{
						"type": "string",
						"enum": []string{"DOMAINS"},
					},
					"domains": map[string]any{
						"type":        "array",
						"description": "Domain names to match",
						"items": map[string]any{
							"type": "string",
						},
					},
				},
				"required": []string{"type", "domains"},
			},
			"macAddressFilter": map[string]any{
				"type":        "object",
				"description": "Required when type is MAC_ADDRESS (source only)",
				"properties": map[string]any{
					"macAddresses": map[string]any{
						"type":        "array",
						"description": "MAC addresses to match",
						"items": map[string]any{
							"type": "string",
						},
					},
				},
				"required": []string{"macAddresses"},
			},
			"networkFilter": map[string]any{
				"type":        "object",
				"description": "Required when type is NETWORK",
				"properties": map[string]any{
					"matchOpposite": map[string]any{
						"type": "boolean",
						"description": "Match all networks except" +
							" the selected ones",
					},
					"networkIds": map[string]any{
						"type":        "array",
						"description": "Network UUIDs to match",
						"items": map[string]any{
							"type": "string",
						},
					},
				},
				"required": []string{"matchOpposite", "networkIds"},
			},
		},
		"required": []string{"type"},
	}
}

// policyInputSchema returns the common JSON schema properties for
// create/update firewall policy tools.
func policyInputSchema() map[string]interface{} {
	return map[string]interface{}{
		"name": map[string]interface{}{
			"type":        "string",
			"description": "Name of the firewall policy",
		},
		"enabled": map[string]interface{}{
			"type":        "boolean",
			"description": "Whether the policy is enabled",
		},
		"action": map[string]interface{}{
			"type":        "object",
			"description": "Policy action",
			"properties": map[string]interface{}{
				"type": map[string]interface{}{
					"type":        "string",
					"description": "Action type: ALLOW, BLOCK, or REJECT",
					"enum": []string{
						"ALLOW",
						"BLOCK",
						"REJECT",
					},
				},
			},
			"required": []string{"type"},
		},
		"source": map[string]any{
			"type":        "object",
			"description": "Traffic source",
			"properties": map[string]any{
				"zoneId": map[string]any{
					"type":        "string",
					"description": "Source firewall zone UUID",
				},
				"trafficFilter": trafficFilterSchema(
					"source",
					[]string{
						"PORT",
						"NETWORK",
						"MAC_ADDRESS",
						"IP_ADDRESS",
						"IPV6_IID",
						"REGION",
						"VPN_SERVER",
						"SITE_TO_SITE_VPN_TUNNEL",
					},
				),
			},
			"required": []string{"zoneId"},
		},
		"destination": map[string]any{
			"type":        "object",
			"description": "Traffic destination",
			"properties": map[string]any{
				"zoneId": map[string]any{
					"type":        "string",
					"description": "Destination firewall zone UUID",
				},
				"trafficFilter": trafficFilterSchema(
					"destination",
					[]string{
						"PORT",
						"NETWORK",
						"IP_ADDRESS",
						"IPV6_IID",
						"REGION",
						"VPN_SERVER",
						"SITE_TO_SITE_VPN_TUNNEL",
						"DOMAIN",
						"APPLICATION",
						"APPLICATION_CATEGORY",
					},
				),
			},
			"required": []string{"zoneId"},
		},
		"ipProtocolScope": map[string]interface{}{
			"type":        "object",
			"description": "IP protocol scope",
			"properties": map[string]interface{}{
				"ipVersion": map[string]interface{}{
					"type":        "string",
					"description": "IP version: IPV4, IPV6, or IPV4_AND_IPV6",
					"enum": []string{
						"IPV4",
						"IPV6",
						"IPV4_AND_IPV6",
					},
				},
			},
			"required": []string{"ipVersion"},
		},
		"loggingEnabled": map[string]interface{}{
			"type":        "boolean",
			"description": "Whether logging is enabled for this policy",
		},
		"description": map[string]interface{}{
			"type":        "string",
			"description": "Optional description of the policy",
		},
		"connectionStateFilter": map[string]interface{}{
			"type":        "array",
			"description": "Optional connection state filter",
			"items": map[string]interface{}{
				"type": "string",
				"enum": []string{
					"NEW",
					"INVALID",
					"ESTABLISHED",
					"RELATED",
				},
			},
		},
		"ipsecFilter": map[string]interface{}{
			"type":        "string",
			"description": "Optional IPsec filter: MATCH_ENCRYPTED or MATCH_NOT_ENCRYPTED",
			"enum": []string{
				"MATCH_ENCRYPTED",
				"MATCH_NOT_ENCRYPTED",
			},
		},
		"schedule": map[string]interface{}{
			"type":        "object",
			"description": "Optional schedule",
			"properties": map[string]interface{}{
				"mode": map[string]interface{}{
					"type":        "string",
					"description": "Schedule mode",
				},
			},
		},
	}
}

// policyParams holds common parameters for create/update.
type policyParams struct {
	SiteID                string                `json:"siteId"`
	FirewallPolicyID      string                `json:"firewallPolicyId"`
	Name                  string                `json:"name"`
	Enabled               bool                  `json:"enabled"`
	Action                policyActionParams    `json:"action"`
	Source                policyEndpointParams  `json:"source"`
	Destination           policyEndpointParams  `json:"destination"`
	IPProtocolScope       policyIPScopeParams   `json:"ipProtocolScope"`
	LoggingEnabled        bool                  `json:"loggingEnabled"`
	Description           *string               `json:"description,omitempty"`
	ConnectionStateFilter []string              `json:"connectionStateFilter,omitempty"`
	IpsecFilter           *string               `json:"ipsecFilter,omitempty"`
	Schedule              *policyScheduleParams `json:"schedule,omitempty"`
}

type policyActionParams struct {
	Type string `json:"type"`
}

type policyEndpointParams struct {
	ZoneID        string           `json:"zoneId"`
	TrafficFilter *json.RawMessage `json:"trafficFilter,omitempty"`
}

type policyIPScopeParams struct {
	IPVersion string `json:"ipVersion"`
}

type policyScheduleParams struct {
	Mode string `json:"mode"`
}

// buildRequestBody converts parsed params to JSON bytes suitable for
// forwarding directly to the UniFi API. It uses map[string]any so that
// traffic filter objects — whose nested sub-objects are discriminated unions
// not captured by the generated types — are forwarded verbatim.
func buildRequestBody(params *policyParams) ([]byte, error) {
	sourceZoneID, err := resolveUUID(
		"source.zoneId",
		params.Source.ZoneID,
	)
	if err != nil {
		return nil, err
	}

	destZoneID, err := resolveUUID(
		"destination.zoneId",
		params.Destination.ZoneID,
	)
	if err != nil {
		return nil, err
	}

	source := map[string]any{
		"zoneId": sourceZoneID.String(),
	}
	if params.Source.TrafficFilter != nil {
		var tf any
		if err := json.Unmarshal(
			*params.Source.TrafficFilter,
			&tf,
		); err != nil {
			return nil, fmt.Errorf(
				"invalid source.trafficFilter: %w",
				err,
			)
		}
		source["trafficFilter"] = tf
	}

	destination := map[string]any{
		"zoneId": destZoneID.String(),
	}
	if params.Destination.TrafficFilter != nil {
		var tf any
		if err := json.Unmarshal(
			*params.Destination.TrafficFilter,
			&tf,
		); err != nil {
			return nil, fmt.Errorf(
				"invalid destination.trafficFilter: %w",
				err,
			)
		}
		destination["trafficFilter"] = tf
	}

	body := map[string]any{
		"name":    params.Name,
		"enabled": params.Enabled,
		"action": map[string]any{
			"type": params.Action.Type,
		},
		"source":      source,
		"destination": destination,
		"ipProtocolScope": map[string]any{
			"ipVersion": params.IPProtocolScope.IPVersion,
		},
		"loggingEnabled": params.LoggingEnabled,
	}

	if params.Description != nil {
		body["description"] = *params.Description
	}

	if len(params.ConnectionStateFilter) > 0 {
		body["connectionStateFilter"] = params.ConnectionStateFilter
	}

	if params.IpsecFilter != nil {
		body["ipsecFilter"] = *params.IpsecFilter
	}

	if params.Schedule != nil {
		body["schedule"] = map[string]any{
			"mode": params.Schedule.Mode,
		}
	}

	return json.Marshal(body)
}

// ListFirewallPolicies implements the list_firewall_policies MCP tool.
type ListFirewallPolicies struct {
	baseTool
}

// NewListFirewallPolicies creates a new ListFirewallPolicies tool.
func NewListFirewallPolicies(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListFirewallPolicies {
	return &ListFirewallPolicies{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *ListFirewallPolicies) Description() string {
	return "List firewall policies for a UniFi site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListFirewallPolicies) InputSchema() map[string]interface{} {
	return listSchema()
}

// Execute runs the tool.
func (t *ListFirewallPolicies) Execute(
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

	resp, err := t.client.GetFirewallPoliciesWithResponse(
		ctx,
		siteID,
		&unifi.GetFirewallPoliciesParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list firewall policies: %w",
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
		return "No firewall policies found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Firewall Policies (%d of %d):\n\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, policy := range page.Data {
		fmt.Fprintf(&b, "%d. %s", i+1, formatPolicy(&policy))
		if i < len(page.Data)-1 {
			fmt.Fprintln(&b)
		}
	}

	return b.String(), nil
}

// GetFirewallPolicy implements the get_firewall_policy MCP tool.
type GetFirewallPolicy struct {
	baseTool
}

// NewGetFirewallPolicy creates a new GetFirewallPolicy tool.
func NewGetFirewallPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetFirewallPolicy {
	return &GetFirewallPolicy{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *GetFirewallPolicy) Description() string {
	return "Get details of a specific firewall policy"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetFirewallPolicy) InputSchema() map[string]interface{} {
	return siteAndIDSchema(
		"firewallPolicyId",
		"Firewall policy UUID",
	)
}

// Execute runs the tool.
func (t *GetFirewallPolicy) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID           string `json:"siteId"`
		FirewallPolicyID string `json:"firewallPolicyId"`
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

	policyID, err := resolveUUID(
		"firewallPolicyId",
		params.FirewallPolicyID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetFirewallPolicyWithResponse(
		ctx,
		siteID,
		policyID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get firewall policy: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatPolicy(resp.JSON200), nil
}

// CreateFirewallPolicy implements the create_firewall_policy MCP tool.
type CreateFirewallPolicy struct {
	baseTool
}

// NewCreateFirewallPolicy creates a new CreateFirewallPolicy tool.
func NewCreateFirewallPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *CreateFirewallPolicy {
	return &CreateFirewallPolicy{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *CreateFirewallPolicy) Description() string {
	return "Create a new firewall policy"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *CreateFirewallPolicy) InputSchema() map[string]interface{} {
	props := policyInputSchema()
	props["siteId"] = siteIDSchema()

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required": []string{
			"name",
			"enabled",
			"action",
			"source",
			"destination",
			"ipProtocolScope",
			"loggingEnabled",
		},
	}
}

// Execute runs the tool.
func (t *CreateFirewallPolicy) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params policyParams
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	if params.Name == "" {
		return "", fmt.Errorf("name is required")
	}
	if params.Action.Type == "" {
		return "", fmt.Errorf("action.type is required")
	}

	siteID, err := resolveSiteID(
		params.SiteID,
		t.defaultSiteID,
	)
	if err != nil {
		return "", err
	}

	bodyJSON, err := buildRequestBody(&params)
	if err != nil {
		return "", err
	}

	resp, err := t.client.CreateFirewallPolicyWithBodyWithResponse(
		ctx,
		siteID,
		"application/json",
		bytes.NewReader(bodyJSON),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create firewall policy: %w",
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
		"Firewall policy created:\n%s",
		formatPolicy(resp.JSON201),
	), nil
}

// UpdateFirewallPolicy implements the update_firewall_policy tool.
type UpdateFirewallPolicy struct {
	baseTool
}

// NewUpdateFirewallPolicy creates a new UpdateFirewallPolicy tool.
func NewUpdateFirewallPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *UpdateFirewallPolicy {
	return &UpdateFirewallPolicy{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *UpdateFirewallPolicy) Description() string {
	return "Update an existing firewall policy"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *UpdateFirewallPolicy) InputSchema() map[string]interface{} {
	props := policyInputSchema()
	props["siteId"] = siteIDSchema()
	props["firewallPolicyId"] = map[string]interface{}{
		"type":        "string",
		"description": "Firewall policy UUID to update",
	}

	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required": []string{
			"firewallPolicyId",
			"name",
			"enabled",
			"action",
			"source",
			"destination",
			"ipProtocolScope",
			"loggingEnabled",
		},
	}
}

// Execute runs the tool.
func (t *UpdateFirewallPolicy) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params policyParams
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

	policyID, err := resolveUUID(
		"firewallPolicyId",
		params.FirewallPolicyID,
	)
	if err != nil {
		return "", err
	}

	bodyJSON, err := buildRequestBody(&params)
	if err != nil {
		return "", err
	}

	resp, err := t.client.UpdateFirewallPolicyWithBodyWithResponse(
		ctx,
		siteID,
		policyID,
		"application/json",
		bytes.NewReader(bodyJSON),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update firewall policy: %w",
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
		"Firewall policy updated:\n%s",
		formatPolicy(resp.JSON200),
	), nil
}

// DeleteFirewallPolicy implements the delete_firewall_policy tool.
type DeleteFirewallPolicy struct {
	baseTool
}

// NewDeleteFirewallPolicy creates a new DeleteFirewallPolicy tool.
func NewDeleteFirewallPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *DeleteFirewallPolicy {
	return &DeleteFirewallPolicy{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *DeleteFirewallPolicy) Description() string {
	return "Delete a firewall policy"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *DeleteFirewallPolicy) InputSchema() map[string]interface{} {
	return siteAndIDSchema(
		"firewallPolicyId",
		"Firewall policy UUID to delete",
	)
}

// Execute runs the tool.
func (t *DeleteFirewallPolicy) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID           string `json:"siteId"`
		FirewallPolicyID string `json:"firewallPolicyId"`
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

	policyID, err := resolveUUID(
		"firewallPolicyId",
		params.FirewallPolicyID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.DeleteFirewallPolicyWithResponse(
		ctx,
		siteID,
		policyID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete firewall policy: %w",
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
		"Firewall policy %s deleted successfully.",
		policyID.String(),
	), nil
}

// PatchFirewallPolicy implements the patch_firewall_policy tool.
type PatchFirewallPolicy struct {
	baseTool
}

// NewPatchFirewallPolicy creates a new PatchFirewallPolicy tool.
func NewPatchFirewallPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *PatchFirewallPolicy {
	return &PatchFirewallPolicy{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *PatchFirewallPolicy) Description() string {
	return "Partially update a firewall policy (toggle logging)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *PatchFirewallPolicy) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"firewallPolicyId": map[string]interface{}{
				"type":        "string",
				"description": "Firewall policy UUID to patch",
			},
			"loggingEnabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether logging is enabled",
			},
		},
		"required": []string{"firewallPolicyId"},
	}
}

// Execute runs the tool.
func (t *PatchFirewallPolicy) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID           string `json:"siteId"`
		FirewallPolicyID string `json:"firewallPolicyId"`
		LoggingEnabled   *bool  `json:"loggingEnabled"`
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

	policyID, err := resolveUUID(
		"firewallPolicyId",
		params.FirewallPolicyID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.PatchFirewallPolicyWithResponse(
		ctx,
		siteID,
		policyID,
		unifi.PatchFirewallPolicyJSONRequestBody{
			LoggingEnabled: params.LoggingEnabled,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to patch firewall policy: %w",
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
		"Firewall policy patched:\n%s",
		formatPolicy(resp.JSON200),
	), nil
}
