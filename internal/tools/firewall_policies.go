package tools

import (
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
		"source": map[string]interface{}{
			"type":        "object",
			"description": "Traffic source",
			"properties": map[string]interface{}{
				"zoneId": map[string]interface{}{
					"type":        "string",
					"description": "Source firewall zone UUID",
				},
				"trafficFilter": map[string]interface{}{
					"type":        "object",
					"description": "Optional source traffic filter",
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Traffic filter type",
						},
					},
				},
			},
			"required": []string{"zoneId"},
		},
		"destination": map[string]interface{}{
			"type":        "object",
			"description": "Traffic destination",
			"properties": map[string]interface{}{
				"zoneId": map[string]interface{}{
					"type":        "string",
					"description": "Destination firewall zone UUID",
				},
				"trafficFilter": map[string]interface{}{
					"type":        "object",
					"description": "Optional destination traffic filter",
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Traffic filter type",
						},
					},
				},
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
	ZoneID        string                     `json:"zoneId"`
	TrafficFilter *policyTrafficFilterParams `json:"trafficFilter,omitempty"`
}

type policyIPScopeParams struct {
	IPVersion string `json:"ipVersion"`
}

type policyScheduleParams struct {
	Mode string `json:"mode"`
}

type policyTrafficFilterParams struct {
	Type string `json:"type"`
}

// buildRequestBody converts parsed params to the generated request type.
func buildRequestBody(
	params *policyParams,
) (unifi.CreateOrUpdateFirewallPolicy, error) {
	sourceZoneID, err := resolveUUID(
		"source.zoneId",
		params.Source.ZoneID,
	)
	if err != nil {
		return unifi.CreateOrUpdateFirewallPolicy{}, err
	}

	destZoneID, err := resolveUUID(
		"destination.zoneId",
		params.Destination.ZoneID,
	)
	if err != nil {
		return unifi.CreateOrUpdateFirewallPolicy{}, err
	}

	body := unifi.CreateOrUpdateFirewallPolicy{
		Name:    params.Name,
		Enabled: params.Enabled,
		Action: unifi.FirewallPolicyAction{
			Type: params.Action.Type,
		},
		Source: unifi.FirewallPolicySource{
			ZoneId: sourceZoneID,
		},
		Destination: unifi.FirewallPolicyDestination{
			ZoneId: destZoneID,
		},
		IpProtocolScope: unifi.FirewallPolicyIPProtocolScope{
			IpVersion: params.IPProtocolScope.IPVersion,
		},
		LoggingEnabled: params.LoggingEnabled,
		Description:    params.Description,
	}

	if params.Source.TrafficFilter != nil {
		body.Source.TrafficFilter = &unifi.FirewallPolicySourceTrafficFilter{
			Type: params.Source.TrafficFilter.Type,
		}
	}

	if params.Destination.TrafficFilter != nil {
		body.Destination.TrafficFilter = &unifi.FirewallPolicyDestinationTrafficFilter{
			Type: params.Destination.TrafficFilter.Type,
		}
	}

	if len(params.ConnectionStateFilter) > 0 {
		filters := make(
			[]unifi.CreateOrUpdateFirewallPolicyConnectionStateFilter,
			len(params.ConnectionStateFilter),
		)
		for i, f := range params.ConnectionStateFilter {
			filters[i] = unifi.CreateOrUpdateFirewallPolicyConnectionStateFilter(
				f,
			)
		}
		body.ConnectionStateFilter = &filters
	}

	if params.IpsecFilter != nil {
		f := unifi.CreateOrUpdateFirewallPolicyIpsecFilter(
			*params.IpsecFilter,
		)
		body.IpsecFilter = &f
	}

	if params.Schedule != nil {
		body.Schedule = &unifi.FirewallSchedule{
			Mode: params.Schedule.Mode,
		}
	}

	return body, nil
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
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"firewallPolicyId": map[string]interface{}{
				"type":        "string",
				"description": "Firewall policy UUID",
			},
		},
		"required": []string{"firewallPolicyId"},
	}
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

	body, err := buildRequestBody(&params)
	if err != nil {
		return "", err
	}

	resp, err := t.client.CreateFirewallPolicyWithResponse(
		ctx,
		siteID,
		body,
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

	body, err := buildRequestBody(&params)
	if err != nil {
		return "", err
	}

	resp, err := t.client.UpdateFirewallPolicyWithResponse(
		ctx,
		siteID,
		policyID,
		body,
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
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"firewallPolicyId": map[string]interface{}{
				"type":        "string",
				"description": "Firewall policy UUID to delete",
			},
		},
		"required": []string{"firewallPolicyId"},
	}
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

// GetFirewallPolicyOrdering implements the
// get_firewall_policy_ordering MCP tool.
type GetFirewallPolicyOrdering struct {
	baseTool
}

// NewGetFirewallPolicyOrdering creates a new tool.
func NewGetFirewallPolicyOrdering(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetFirewallPolicyOrdering {
	return &GetFirewallPolicyOrdering{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *GetFirewallPolicyOrdering) Description() string {
	return "Get firewall policy ordering between two zones"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetFirewallPolicyOrdering) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"sourceZoneId": map[string]interface{}{
				"type":        "string",
				"description": "Source firewall zone UUID",
			},
			"destinationZoneId": map[string]interface{}{
				"type":        "string",
				"description": "Destination firewall zone UUID",
			},
		},
		"required": []string{
			"sourceZoneId",
			"destinationZoneId",
		},
	}
}

// Execute runs the tool.
func (t *GetFirewallPolicyOrdering) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID            string `json:"siteId"`
		SourceZoneID      string `json:"sourceZoneId"`
		DestinationZoneID string `json:"destinationZoneId"`
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

	srcZone, err := resolveUUID(
		"sourceZoneId",
		params.SourceZoneID,
	)
	if err != nil {
		return "", err
	}

	dstZone, err := resolveUUID(
		"destinationZoneId",
		params.DestinationZoneID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetFirewallPolicyOrderingWithResponse(
		ctx,
		siteID,
		&unifi.GetFirewallPolicyOrderingParams{
			SourceFirewallZoneId:      srcZone,
			DestinationFirewallZoneId: dstZone,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get policy ordering: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	ordering := resp.JSON200.OrderedFirewallPolicyIds
	var b strings.Builder
	fmt.Fprintln(&b, "Policy Ordering:")
	fmt.Fprintln(&b, "\nBefore System-Defined:")
	if len(ordering.BeforeSystemDefined) == 0 {
		fmt.Fprintln(&b, "  (none)")
	}
	for i, id := range ordering.BeforeSystemDefined {
		fmt.Fprintf(&b, "  %d. %s\n", i+1, id.String())
	}
	fmt.Fprintln(&b, "\nAfter System-Defined:")
	if len(ordering.AfterSystemDefined) == 0 {
		fmt.Fprintln(&b, "  (none)")
	}
	for i, id := range ordering.AfterSystemDefined {
		fmt.Fprintf(&b, "  %d. %s\n", i+1, id.String())
	}

	return b.String(), nil
}

// UpdateFirewallPolicyOrdering implements the
// update_firewall_policy_ordering MCP tool.
type UpdateFirewallPolicyOrdering struct {
	baseTool
}

// NewUpdateFirewallPolicyOrdering creates a new tool.
func NewUpdateFirewallPolicyOrdering(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *UpdateFirewallPolicyOrdering {
	return &UpdateFirewallPolicyOrdering{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *UpdateFirewallPolicyOrdering) Description() string {
	return "Update firewall policy ordering between two zones"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *UpdateFirewallPolicyOrdering) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"sourceZoneId": map[string]interface{}{
				"type":        "string",
				"description": "Source firewall zone UUID",
			},
			"destinationZoneId": map[string]interface{}{
				"type":        "string",
				"description": "Destination firewall zone UUID",
			},
			"beforeSystemDefined": map[string]interface{}{
				"type":        "array",
				"description": "Policy UUIDs ordered before system-defined policies",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"afterSystemDefined": map[string]interface{}{
				"type":        "array",
				"description": "Policy UUIDs ordered after system-defined policies",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"required": []string{
			"sourceZoneId",
			"destinationZoneId",
			"beforeSystemDefined",
			"afterSystemDefined",
		},
	}
}

// Execute runs the tool.
func (t *UpdateFirewallPolicyOrdering) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID              string   `json:"siteId"`
		SourceZoneID        string   `json:"sourceZoneId"`
		DestinationZoneID   string   `json:"destinationZoneId"`
		BeforeSystemDefined []string `json:"beforeSystemDefined"`
		AfterSystemDefined  []string `json:"afterSystemDefined"`
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

	srcZone, err := resolveUUID(
		"sourceZoneId",
		params.SourceZoneID,
	)
	if err != nil {
		return "", err
	}

	dstZone, err := resolveUUID(
		"destinationZoneId",
		params.DestinationZoneID,
	)
	if err != nil {
		return "", err
	}

	beforeIDs, err := resolveUUIDs(
		"beforeSystemDefined",
		params.BeforeSystemDefined,
	)
	if err != nil {
		return "", err
	}

	afterIDs, err := resolveUUIDs(
		"afterSystemDefined",
		params.AfterSystemDefined,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.UpdateFirewallPolicyOrderingWithResponse(
		ctx,
		siteID,
		&unifi.UpdateFirewallPolicyOrderingParams{
			SourceFirewallZoneId:      srcZone,
			DestinationFirewallZoneId: dstZone,
		},
		unifi.UpdateFirewallPolicyOrderingJSONRequestBody{
			OrderedFirewallPolicyIds: unifi.OrderedFirewallPolicyIDs{
				BeforeSystemDefined: beforeIDs,
				AfterSystemDefined:  afterIDs,
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update policy ordering: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return "Firewall policy ordering updated successfully.", nil
}
