package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

// --- list_acl_rules ---

// ListACLRules implements the list_acl_rules MCP tool.
type ListACLRules struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListACLRules creates a new ListACLRules tool.
func NewListACLRules(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListACLRules {
	return &ListACLRules{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListACLRules) Description() string {
	return "List ACL rules for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListACLRules) InputSchema() map[string]interface{} {
	return listSchema()
}

// Execute runs the tool.
func (t *ListACLRules) Execute(
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

	resp, err := t.client.GetAclRulePageWithResponse(
		ctx,
		siteID,
		&unifi.GetAclRulePageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list ACL rules: %w",
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
		return "No ACL rules found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"ACL Rules (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, r := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Type: %s\n   Action: %s\n   Enabled: %t\n   Index: %d\n",
			i+1,
			r.Name,
			r.Id.String(),
			r.Type,
			r.Action,
			r.Enabled,
			r.Index,
		)
	}

	return b.String(), nil
}

// --- get_acl_rule ---

// GetACLRule implements the get_acl_rule MCP tool.
type GetACLRule struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewGetACLRule creates a new GetACLRule tool.
func NewGetACLRule(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetACLRule {
	return &GetACLRule{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *GetACLRule) Description() string {
	return "Get details of a specific ACL rule"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetACLRule) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"aclRuleId": map[string]interface{}{
				"type":        "string",
				"description": "ACL rule UUID",
			},
		},
		"required": []string{"aclRuleId"},
	}
}

// Execute runs the tool.
func (t *GetACLRule) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID    string `json:"siteId"`
		ACLRuleID string `json:"aclRuleId"`
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

	ruleID, err := resolveUUID(
		"aclRuleId",
		params.ACLRuleID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetAclRuleWithResponse(
		ctx,
		siteID,
		ruleID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get ACL rule: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatACLRule(resp.JSON200), nil
}

func formatACLRule(r *unifi.ACLRule) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\n", r.Name)
	fmt.Fprintf(&b, "ID: %s\n", r.Id.String())
	fmt.Fprintf(&b, "Type: %s\n", r.Type)
	fmt.Fprintf(&b, "Action: %s\n", string(r.Action))
	fmt.Fprintf(&b, "Enabled: %t\n", r.Enabled)
	fmt.Fprintf(&b, "Index: %d\n", r.Index)
	if r.Description != nil {
		fmt.Fprintf(&b, "Description: %s\n", *r.Description)
	}
	if r.SourceFilter != nil {
		src, _ := json.MarshalIndent(
			r.SourceFilter,
			"",
			"  ",
		)
		fmt.Fprintf(&b, "Source Filter:\n%s\n", src)
	}
	if r.DestinationFilter != nil {
		dst, _ := json.MarshalIndent(
			r.DestinationFilter,
			"",
			"  ",
		)
		fmt.Fprintf(&b, "Destination Filter:\n%s\n", dst)
	}
	if r.EnforcingDeviceFilter != nil {
		fmt.Fprintf(
			&b,
			"Enforcing Device Filter Type: %s\n",
			r.EnforcingDeviceFilter.Type,
		)
	}
	return b.String()
}

// aclRuleInputSchema returns the common JSON schema properties for
// create/update ACL rule tools.
func aclRuleInputSchema() map[string]interface{} {
	return map[string]interface{}{
		"siteId": siteIDSchema(),
		"type": map[string]interface{}{
			"type":        "string",
			"description": "ACL rule type",
			"enum":        []string{"IPV4", "MAC"},
		},
		"name": map[string]interface{}{
			"type":        "string",
			"description": "ACL rule name",
		},
		"enabled": map[string]interface{}{
			"type":        "boolean",
			"description": "Whether the rule is enabled",
		},
		"action": map[string]interface{}{
			"type":        "string",
			"description": "ACL rule action",
			"enum":        []string{"ALLOW", "BLOCK"},
		},
		"description": map[string]interface{}{
			"type":        "string",
			"description": "ACL rule description (optional)",
		},
		"sourceFilter": map[string]interface{}{
			"type":        "object",
			"description": "Traffic source filter (type-specific)",
		},
		"destinationFilter": map[string]interface{}{
			"type":        "object",
			"description": "Traffic destination filter (type-specific)",
		},
		"enforcingDeviceFilter": map[string]interface{}{
			"type":        "object",
			"description": "Device filter for enforcement (optional)",
		},
	}
}

// --- create_acl_rule ---

// CreateACLRule implements the create_acl_rule MCP tool.
type CreateACLRule struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewCreateACLRule creates a new CreateACLRule tool.
func NewCreateACLRule(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *CreateACLRule {
	return &CreateACLRule{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *CreateACLRule) Description() string {
	return "Create a new ACL rule (IPv4 or MAC-based access control)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *CreateACLRule) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": aclRuleInputSchema(),
		"required": []string{
			"type",
			"name",
			"enabled",
			"action",
		},
	}
}

// Execute runs the tool.
func (t *CreateACLRule) Execute(
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

	resp, err := t.client.CreateAclRuleWithBodyWithResponse(
		ctx,
		siteID,
		"application/json",
		bytes.NewReader(args),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create ACL rule: %w",
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
		"ACL rule created:\n%s",
		formatACLRule(resp.JSON201),
	), nil
}

// --- update_acl_rule ---

// UpdateACLRule implements the update_acl_rule MCP tool.
type UpdateACLRule struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewUpdateACLRule creates a new UpdateACLRule tool.
func NewUpdateACLRule(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *UpdateACLRule {
	return &UpdateACLRule{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *UpdateACLRule) Description() string {
	return "Update an existing ACL rule"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *UpdateACLRule) InputSchema() map[string]interface{} {
	props := aclRuleInputSchema()
	props["aclRuleId"] = map[string]interface{}{
		"type":        "string",
		"description": "ACL rule UUID",
	}
	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required": []string{
			"aclRuleId",
			"type",
			"name",
			"enabled",
			"action",
		},
	}
}

// Execute runs the tool.
func (t *UpdateACLRule) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID    string `json:"siteId"`
		ACLRuleID string `json:"aclRuleId"`
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

	ruleID, err := resolveUUID(
		"aclRuleId",
		params.ACLRuleID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.UpdateAclRuleWithBodyWithResponse(
		ctx,
		siteID,
		ruleID,
		"application/json",
		bytes.NewReader(args),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update ACL rule: %w",
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
		"ACL rule updated:\n%s",
		formatACLRule(resp.JSON200),
	), nil
}

// --- delete_acl_rule ---

// DeleteACLRule implements the delete_acl_rule MCP tool.
type DeleteACLRule struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewDeleteACLRule creates a new DeleteACLRule tool.
func NewDeleteACLRule(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *DeleteACLRule {
	return &DeleteACLRule{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *DeleteACLRule) Description() string {
	return "Delete an ACL rule"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *DeleteACLRule) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"aclRuleId": map[string]interface{}{
				"type":        "string",
				"description": "ACL rule UUID",
			},
		},
		"required": []string{"aclRuleId"},
	}
}

// Execute runs the tool.
func (t *DeleteACLRule) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID    string `json:"siteId"`
		ACLRuleID string `json:"aclRuleId"`
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

	ruleID, err := resolveUUID(
		"aclRuleId",
		params.ACLRuleID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.DeleteAclRuleWithResponse(
		ctx,
		siteID,
		ruleID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete ACL rule: %w",
			err,
		)
	}

	if resp.StatusCode() != 200 {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return "ACL rule deleted.", nil
}

// --- get_acl_rule_ordering ---

// GetACLRuleOrdering implements the get_acl_rule_ordering
// MCP tool.
type GetACLRuleOrdering struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewGetACLRuleOrdering creates a new GetACLRuleOrdering tool.
func NewGetACLRuleOrdering(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetACLRuleOrdering {
	return &GetACLRuleOrdering{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *GetACLRuleOrdering) Description() string {
	return "Get the current ACL rule ordering"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetACLRuleOrdering) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
		},
	}
}

// Execute runs the tool.
func (t *GetACLRuleOrdering) Execute(
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

	resp, err := t.client.GetAclRuleOrderingWithResponse(
		ctx,
		siteID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get ACL rule ordering: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	ids := resp.JSON200.OrderedAclRuleIds
	if len(ids) == 0 {
		return "No ACL rules in ordering.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"ACL Rule Ordering (%d rules):\n",
		len(ids),
	)
	for i, id := range ids {
		fmt.Fprintf(&b, "%d. %s\n", i+1, id.String())
	}

	return b.String(), nil
}

// --- update_acl_rule_ordering ---

// UpdateACLRuleOrdering implements the
// update_acl_rule_ordering MCP tool.
type UpdateACLRuleOrdering struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewUpdateACLRuleOrdering creates a new
// UpdateACLRuleOrdering tool.
func NewUpdateACLRuleOrdering(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *UpdateACLRuleOrdering {
	return &UpdateACLRuleOrdering{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *UpdateACLRuleOrdering) Description() string {
	return "Update the ACL rule ordering"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *UpdateACLRuleOrdering) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"orderedAclRuleIds": map[string]interface{}{
				"type":        "array",
				"description": "Ordered list of ACL rule UUIDs",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"required": []string{"orderedAclRuleIds"},
	}
}

// Execute runs the tool.
func (t *UpdateACLRuleOrdering) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID            string   `json:"siteId"`
		OrderedACLRuleIDs []string `json:"orderedAclRuleIds"`
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

	ruleIDs, err := resolveUUIDs(
		"orderedAclRuleIds",
		params.OrderedACLRuleIDs,
	)
	if err != nil {
		return "", err
	}

	body := unifi.ACLRuleOrdering{
		OrderedAclRuleIds: ruleIDs,
	}

	resp, err := t.client.UpdateAclRuleOrderingWithResponse(
		ctx,
		siteID,
		body,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update ACL rule ordering: %w",
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
		"ACL rule ordering updated (%d rules).",
		len(resp.JSON200.OrderedAclRuleIds),
	), nil
}
