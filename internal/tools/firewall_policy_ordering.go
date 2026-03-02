package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

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
