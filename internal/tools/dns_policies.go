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
	baseTool
}

// NewListDNSPolicies creates a new ListDNSPolicies tool.
func NewListDNSPolicies(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListDNSPolicies {
	return &ListDNSPolicies{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *ListDNSPolicies) Description() string {
	return "List DNS policies for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListDNSPolicies) InputSchema() map[string]any {
	return listSchema()
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

	// Re-parse raw body to access type-specific fields
	// that the generated DNSPolicy struct doesn't capture.
	var rawPage struct {
		Data []json.RawMessage `json:"data"`
	}
	_ = json.Unmarshal(resp.Body, &rawPage)

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
		if i < len(rawPage.Data) {
			b.WriteString(
				formatDNSRecordDetails(rawPage.Data[i]),
			)
		}
	}

	return b.String(), nil
}

// --- get_dns_policy ---

// GetDNSPolicy implements the get_dns_policy MCP tool.
type GetDNSPolicy struct {
	baseTool
}

// NewGetDNSPolicy creates a new GetDNSPolicy tool.
func NewGetDNSPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetDNSPolicy {
	return &GetDNSPolicy{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *GetDNSPolicy) Description() string {
	return "Get a specific DNS policy"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetDNSPolicy) InputSchema() map[string]any {
	return siteAndIDSchema("dnsPolicyId", "DNS policy UUID")
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

	return formatDNSPolicy(resp.JSON200, resp.Body), nil
}

func formatDNSPolicy(
	p *unifi.DNSPolicy,
	raw []byte,
) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Type: %s\n", p.Type)
	fmt.Fprintf(&b, "ID: %s\n", p.Id.String())
	fmt.Fprintf(&b, "Enabled: %t\n", p.Enabled)
	if p.Domain != nil {
		fmt.Fprintf(&b, "Domain: %s\n", *p.Domain)
	}
	if len(raw) > 0 {
		b.WriteString(formatDNSRecordDetails(raw))
	}
	return b.String()
}

// formatDNSRecordDetails extracts type-specific fields
// from a raw DNS policy JSON object. The generated
// DNSPolicy struct only captures common fields due to
// the discriminated union in the OpenAPI spec.
func formatDNSRecordDetails(
	raw json.RawMessage,
) string {
	var fields struct {
		IPv4Address      *string `json:"ipv4Address"`
		IPv6Address      *string `json:"ipv6Address"`
		TargetDomain     *string `json:"targetDomain"`
		MailServerDomain *string `json:"mailServerDomain"`
		ServerDomain     *string `json:"serverDomain"`
		Service          *string `json:"service"`
		Protocol         *string `json:"protocol"`
		Text             *string `json:"text"`
		IPAddress        *string `json:"ipAddress"`
		TTLSeconds       *int32  `json:"ttlSeconds"`
		Priority         *int32  `json:"priority"`
		Weight           *int32  `json:"weight"`
		Port             *int32  `json:"port"`
	}
	if err := json.Unmarshal(raw, &fields); err != nil {
		return ""
	}

	var b strings.Builder
	if fields.IPv4Address != nil {
		fmt.Fprintf(&b, "   Address: %s\n", *fields.IPv4Address)
	}
	if fields.IPv6Address != nil {
		fmt.Fprintf(&b, "   Address: %s\n", *fields.IPv6Address)
	}
	if fields.TargetDomain != nil {
		fmt.Fprintf(
			&b,
			"   Target: %s\n",
			*fields.TargetDomain,
		)
	}
	if fields.MailServerDomain != nil {
		fmt.Fprintf(
			&b,
			"   Mail Server: %s\n",
			*fields.MailServerDomain,
		)
	}
	if fields.ServerDomain != nil {
		fmt.Fprintf(
			&b,
			"   Server: %s\n",
			*fields.ServerDomain,
		)
	}
	if fields.Service != nil {
		fmt.Fprintf(
			&b,
			"   Service: %s\n",
			*fields.Service,
		)
	}
	if fields.Protocol != nil {
		fmt.Fprintf(
			&b,
			"   Protocol: %s\n",
			*fields.Protocol,
		)
	}
	if fields.Text != nil {
		fmt.Fprintf(&b, "   Text: %s\n", *fields.Text)
	}
	if fields.IPAddress != nil {
		fmt.Fprintf(
			&b,
			"   Forward To: %s\n",
			*fields.IPAddress,
		)
	}
	if fields.TTLSeconds != nil {
		fmt.Fprintf(
			&b,
			"   TTL: %d seconds\n",
			*fields.TTLSeconds,
		)
	}
	if fields.Priority != nil {
		fmt.Fprintf(
			&b, "   Priority: %d\n", *fields.Priority,
		)
	}
	if fields.Weight != nil {
		fmt.Fprintf(&b, "   Weight: %d\n", *fields.Weight)
	}
	if fields.Port != nil {
		fmt.Fprintf(&b, "   Port: %d\n", *fields.Port)
	}
	return b.String()
}

// dnsPolicyInputSchema returns the common JSON schema properties for
// create/update DNS policy tools.
func dnsPolicyInputSchema() map[string]any {
	return map[string]any{
		"siteId": siteIDSchema(),
		"type": map[string]any{
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
		"enabled": map[string]any{
			"type":        "boolean",
			"description": "Whether the policy is enabled",
		},
		"domain": map[string]any{
			"type":        "string",
			"description": "Domain name",
		},
		"ipv4Address": map[string]any{
			"type":        "string",
			"description": "IPv4 address (for A_RECORD)",
		},
		"ipv6Address": map[string]any{
			"type":        "string",
			"description": "IPv6 address (for AAAA_RECORD)",
		},
		"targetDomain": map[string]any{
			"type":        "string",
			"description": "Target domain (for CNAME_RECORD)",
		},
		"mailServerDomain": map[string]any{
			"type":        "string",
			"description": "Mail server domain (for MX_RECORD)",
		},
		"serverDomain": map[string]any{
			"type":        "string",
			"description": "Server domain (for SRV_RECORD)",
		},
		"service": map[string]any{
			"type":        "string",
			"description": "Service name, e.g. _http (for SRV_RECORD)",
		},
		"protocol": map[string]any{
			"type":        "string",
			"description": "Protocol, e.g. _tcp (for SRV_RECORD)",
		},
		"priority": map[string]any{
			"type":        "integer",
			"description": "Priority (for MX_RECORD and SRV_RECORD)",
		},
		"weight": map[string]any{
			"type":        "integer",
			"description": "Weight (for SRV_RECORD)",
		},
		"port": map[string]any{
			"type":        "integer",
			"description": "Port (for SRV_RECORD)",
		},
		"text": map[string]any{
			"type":        "string",
			"description": "Text value (for TXT_RECORD)",
		},
		"ipAddress": map[string]any{
			"type":        "string",
			"description": "DNS server IP to forward to (for FORWARD_DOMAIN)",
		},
		"ttlSeconds": map[string]any{
			"type":        "integer",
			"description": "TTL in seconds (for A_RECORD, AAAA_RECORD, CNAME_RECORD)",
		},
	}
}

// --- create_dns_policy ---

// CreateDNSPolicy implements the create_dns_policy MCP tool.
type CreateDNSPolicy struct {
	baseTool
}

// NewCreateDNSPolicy creates a new CreateDNSPolicy tool.
func NewCreateDNSPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *CreateDNSPolicy {
	return &CreateDNSPolicy{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *CreateDNSPolicy) Description() string {
	return "Create a new DNS policy (A, AAAA, CNAME, MX, SRV, TXT record or domain forwarding)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *CreateDNSPolicy) InputSchema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": dnsPolicyInputSchema(),
		"required":   []string{"type", "enabled"},
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

	body, err := stripKeys(args, "siteId")
	if err != nil {
		return "", err
	}

	resp, err := t.client.CreateDnsPolicyWithBodyWithResponse(
		ctx,
		siteID,
		"application/json",
		bytes.NewReader(body),
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
		formatDNSPolicy(resp.JSON201, resp.Body),
	), nil
}

// --- update_dns_policy ---

// UpdateDNSPolicy implements the update_dns_policy MCP tool.
type UpdateDNSPolicy struct {
	baseTool
}

// NewUpdateDNSPolicy creates a new UpdateDNSPolicy tool.
func NewUpdateDNSPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *UpdateDNSPolicy {
	return &UpdateDNSPolicy{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *UpdateDNSPolicy) Description() string {
	return "Update an existing DNS policy"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *UpdateDNSPolicy) InputSchema() map[string]any {
	props := dnsPolicyInputSchema()
	props["dnsPolicyId"] = map[string]any{
		"type":        "string",
		"description": "DNS policy UUID",
	}
	return map[string]any{
		"type":       "object",
		"properties": props,
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

	body, err := stripKeys(args, "siteId", "dnsPolicyId")
	if err != nil {
		return "", err
	}

	resp, err := t.client.UpdateDnsPolicyWithBodyWithResponse(
		ctx,
		siteID,
		policyID,
		"application/json",
		bytes.NewReader(body),
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
		formatDNSPolicy(resp.JSON200, resp.Body),
	), nil
}

// --- delete_dns_policy ---

// DeleteDNSPolicy implements the delete_dns_policy MCP tool.
type DeleteDNSPolicy struct {
	baseTool
}

// NewDeleteDNSPolicy creates a new DeleteDNSPolicy tool.
func NewDeleteDNSPolicy(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *DeleteDNSPolicy {
	return &DeleteDNSPolicy{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *DeleteDNSPolicy) Description() string {
	return "Delete a DNS policy"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *DeleteDNSPolicy) InputSchema() map[string]any {
	return siteAndIDSchema("dnsPolicyId", "DNS policy UUID")
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
