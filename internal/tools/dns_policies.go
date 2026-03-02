package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

// --- list_dns_policies ---

// ListDNSPolicies implements the list_dns_policies MCP tool.
type ListDNSPolicies struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListDNSPolicies creates a new ListDNSPolicies tool.
func NewListDNSPolicies(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListDNSPolicies {
	return &ListDNSPolicies{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListDNSPolicies) Description() string {
	return "List DNS policies for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListDNSPolicies) InputSchema() map[string]interface{} {
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
func (t *ListDNSPolicies) Execute(
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

	resp, err := t.client.GetDnsPolicyPageWithResponse(
		ctx,
		siteID,
		&unifi.GetDnsPolicyPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list DNS policies: %w",
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
		return "No DNS policies found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"DNS Policies (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, policy := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Enabled: %t\n",
			i+1,
			policy.Type,
			policy.Id.String(),
			policy.Enabled,
		)
		if policy.Domain != nil {
			fmt.Fprintf(
				&b,
				"   Domain: %s\n",
				*policy.Domain,
			)
		}
	}

	return b.String(), nil
}

// --- get_dns_policy ---

// GetDNSPolicy implements the get_dns_policy MCP tool.
type GetDNSPolicy struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewGetDNSPolicy creates a new GetDNSPolicy tool.
func NewGetDNSPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetDNSPolicy {
	return &GetDNSPolicy{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *GetDNSPolicy) Description() string {
	return "Get a specific DNS policy"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetDNSPolicy) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"dnsPolicyId": map[string]interface{}{
				"type":        "string",
				"description": "DNS policy UUID",
			},
		},
		"required": []string{"dnsPolicyId"},
	}
}

// Execute runs the tool.
func (t *GetDNSPolicy) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID      string `json:"siteId"`
		DNSPolicyID string `json:"dnsPolicyId"`
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
		"dnsPolicyId",
		params.DNSPolicyID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetDnsPolicyWithResponse(
		ctx,
		siteID,
		policyID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get DNS policy: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatDNSPolicy(resp.JSON200), nil
}

func formatDNSPolicy(p *unifi.DNSPolicy) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Type: %s\n", p.Type)
	fmt.Fprintf(&b, "ID: %s\n", p.Id.String())
	fmt.Fprintf(&b, "Enabled: %t\n", p.Enabled)
	if p.Domain != nil {
		fmt.Fprintf(&b, "Domain: %s\n", *p.Domain)
	}
	return b.String()
}

// --- create_dns_policy ---

// CreateDNSPolicy implements the create_dns_policy MCP tool.
type CreateDNSPolicy struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewCreateDNSPolicy creates a new CreateDNSPolicy tool.
func NewCreateDNSPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *CreateDNSPolicy {
	return &CreateDNSPolicy{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *CreateDNSPolicy) Description() string {
	return "Create a new DNS policy (A, AAAA, CNAME, MX, SRV, TXT record or domain forwarding)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *CreateDNSPolicy) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"type": map[string]interface{}{
				"type":        "string",
				"description": "DNS record type",
				"enum": []string{
					"A_RECORD",
					"AAAA_RECORD",
					"CNAME_RECORD",
					"FORWARD_DOMAIN",
					"MX_RECORD",
					"SRV_RECORD",
					"TXT_RECORD",
				},
			},
			"enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the policy is enabled",
			},
			"domain": map[string]interface{}{
				"type":        "string",
				"description": "Domain name",
			},
			"address": map[string]interface{}{
				"type":        "string",
				"description": "IP address (for A/AAAA records)",
			},
			"target": map[string]interface{}{
				"type":        "string",
				"description": "Target hostname (for CNAME/SRV records)",
			},
			"server": map[string]interface{}{
				"type":        "string",
				"description": "Mail server (for MX records)",
			},
			"priority": map[string]interface{}{
				"type":        "integer",
				"description": "Priority (for MX/SRV records)",
			},
			"weight": map[string]interface{}{
				"type":        "integer",
				"description": "Weight (for SRV records)",
			},
			"port": map[string]interface{}{
				"type":        "integer",
				"description": "Port (for SRV records)",
			},
			"value": map[string]interface{}{
				"type":        "string",
				"description": "Text value (for TXT records)",
			},
			"forwardTo": map[string]interface{}{
				"type":        "string",
				"description": "DNS server to forward to (for FORWARD_DOMAIN)",
			},
		},
		"required": []string{"type", "enabled"},
	}
}

// Execute runs the tool.
func (t *CreateDNSPolicy) Execute(
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

	resp, err := t.client.CreateDnsPolicyWithBodyWithResponse(
		ctx,
		siteID,
		"application/json",
		bytes.NewReader(args),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create DNS policy: %w",
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
		"DNS policy created:\n%s",
		formatDNSPolicy(resp.JSON201),
	), nil
}

// --- update_dns_policy ---

// UpdateDNSPolicy implements the update_dns_policy MCP tool.
type UpdateDNSPolicy struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewUpdateDNSPolicy creates a new UpdateDNSPolicy tool.
func NewUpdateDNSPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *UpdateDNSPolicy {
	return &UpdateDNSPolicy{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *UpdateDNSPolicy) Description() string {
	return "Update an existing DNS policy"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *UpdateDNSPolicy) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"dnsPolicyId": map[string]interface{}{
				"type":        "string",
				"description": "DNS policy UUID",
			},
			"type": map[string]interface{}{
				"type":        "string",
				"description": "DNS record type",
				"enum": []string{
					"A_RECORD",
					"AAAA_RECORD",
					"CNAME_RECORD",
					"FORWARD_DOMAIN",
					"MX_RECORD",
					"SRV_RECORD",
					"TXT_RECORD",
				},
			},
			"enabled": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the policy is enabled",
			},
			"domain": map[string]interface{}{
				"type":        "string",
				"description": "Domain name",
			},
			"address": map[string]interface{}{
				"type":        "string",
				"description": "IP address (for A/AAAA records)",
			},
			"target": map[string]interface{}{
				"type":        "string",
				"description": "Target hostname (for CNAME/SRV records)",
			},
			"server": map[string]interface{}{
				"type":        "string",
				"description": "Mail server (for MX records)",
			},
			"priority": map[string]interface{}{
				"type":        "integer",
				"description": "Priority (for MX/SRV records)",
			},
			"weight": map[string]interface{}{
				"type":        "integer",
				"description": "Weight (for SRV records)",
			},
			"port": map[string]interface{}{
				"type":        "integer",
				"description": "Port (for SRV records)",
			},
			"value": map[string]interface{}{
				"type":        "string",
				"description": "Text value (for TXT records)",
			},
			"forwardTo": map[string]interface{}{
				"type":        "string",
				"description": "DNS server to forward to (for FORWARD_DOMAIN)",
			},
		},
		"required": []string{
			"dnsPolicyId",
			"type",
			"enabled",
		},
	}
}

// Execute runs the tool.
func (t *UpdateDNSPolicy) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID      string `json:"siteId"`
		DNSPolicyID string `json:"dnsPolicyId"`
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
		"dnsPolicyId",
		params.DNSPolicyID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.UpdateDnsPolicyWithBodyWithResponse(
		ctx,
		siteID,
		policyID,
		"application/json",
		bytes.NewReader(args),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update DNS policy: %w",
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
		"DNS policy updated:\n%s",
		formatDNSPolicy(resp.JSON200),
	), nil
}

// --- delete_dns_policy ---

// DeleteDNSPolicy implements the delete_dns_policy MCP tool.
type DeleteDNSPolicy struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewDeleteDNSPolicy creates a new DeleteDNSPolicy tool.
func NewDeleteDNSPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *DeleteDNSPolicy {
	return &DeleteDNSPolicy{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *DeleteDNSPolicy) Description() string {
	return "Delete a DNS policy"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *DeleteDNSPolicy) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"dnsPolicyId": map[string]interface{}{
				"type":        "string",
				"description": "DNS policy UUID",
			},
		},
		"required": []string{"dnsPolicyId"},
	}
}

// Execute runs the tool.
func (t *DeleteDNSPolicy) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID      string `json:"siteId"`
		DNSPolicyID string `json:"dnsPolicyId"`
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
		"dnsPolicyId",
		params.DNSPolicyID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.DeleteDnsPolicyWithResponse(
		ctx,
		siteID,
		policyID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete DNS policy: %w",
			err,
		)
	}

	if resp.StatusCode() != 200 {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return "DNS policy deleted.", nil
}
