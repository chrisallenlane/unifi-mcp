package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

// formatTrafficMatchingList formats a single traffic matching list for
// display.
func formatTrafficMatchingList(l *unifi.TrafficMatchingList) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\n", l.Name)
	fmt.Fprintf(&b, "ID: %s\n", l.Id.String())
	fmt.Fprintf(&b, "Type: %s\n", l.Type)
	return b.String()
}

// --- list_traffic_matching_lists ---

// ListTrafficMatchingLists implements the list_traffic_matching_lists
// MCP tool.
type ListTrafficMatchingLists struct {
	baseTool
}

// NewListTrafficMatchingLists creates a new ListTrafficMatchingLists
// tool.
func NewListTrafficMatchingLists(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListTrafficMatchingLists {
	return &ListTrafficMatchingLists{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *ListTrafficMatchingLists) Description() string {
	return "List traffic matching lists for a UniFi site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListTrafficMatchingLists) InputSchema() map[string]interface{} {
	return listSchema()
}

// Execute runs the tool.
func (t *ListTrafficMatchingLists) Execute(
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

	resp, err := t.client.GetTrafficMatchingListsWithResponse(
		ctx,
		siteID,
		&unifi.GetTrafficMatchingListsParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list traffic matching lists: %w",
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
		return "No traffic matching lists found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Traffic Matching Lists (%d of %d):\n\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, item := range page.Data {
		fmt.Fprintf(&b, "%d. %s", i+1, formatTrafficMatchingList(&item))
		if i < len(page.Data)-1 {
			fmt.Fprintln(&b)
		}
	}

	return b.String(), nil
}

// --- get_traffic_matching_list ---

// GetTrafficMatchingList implements the get_traffic_matching_list
// MCP tool.
type GetTrafficMatchingList struct {
	baseTool
}

// NewGetTrafficMatchingList creates a new GetTrafficMatchingList tool.
func NewGetTrafficMatchingList(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetTrafficMatchingList {
	return &GetTrafficMatchingList{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *GetTrafficMatchingList) Description() string {
	return "Get details of a specific traffic matching list"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetTrafficMatchingList) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"trafficMatchingListId": map[string]interface{}{
				"type":        "string",
				"description": "Traffic matching list UUID",
			},
		},
		"required": []string{"trafficMatchingListId"},
	}
}

// Execute runs the tool.
func (t *GetTrafficMatchingList) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID                string `json:"siteId"`
		TrafficMatchingListID string `json:"trafficMatchingListId"`
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

	listID, err := resolveUUID(
		"trafficMatchingListId",
		params.TrafficMatchingListID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetTrafficMatchingListWithResponse(
		ctx,
		siteID,
		listID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get traffic matching list: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatTrafficMatchingList(resp.JSON200), nil
}

// trafficMatchingListInputSchema returns the common JSON schema
// properties for create/update traffic matching list tools.
func trafficMatchingListInputSchema() map[string]interface{} {
	return map[string]interface{}{
		"siteId": siteIDSchema(),
		"name": map[string]interface{}{
			"type":        "string",
			"description": "Name of the traffic matching list",
		},
		"type": map[string]interface{}{
			"type":        "string",
			"description": "Type of traffic matching list",
			"enum": []string{
				"IPV4_ADDRESSES",
				"IPV6_ADDRESSES",
				"PORTS",
			},
		},
	}
}

// --- create_traffic_matching_list ---

// CreateTrafficMatchingList implements the create_traffic_matching_list
// MCP tool.
type CreateTrafficMatchingList struct {
	baseTool
}

// NewCreateTrafficMatchingList creates a new CreateTrafficMatchingList
// tool.
func NewCreateTrafficMatchingList(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *CreateTrafficMatchingList {
	return &CreateTrafficMatchingList{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *CreateTrafficMatchingList) Description() string {
	return "Create a new traffic matching list"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *CreateTrafficMatchingList) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": trafficMatchingListInputSchema(),
		"required":   []string{"name", "type"},
	}
}

// Execute runs the tool.
func (t *CreateTrafficMatchingList) Execute(
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

	resp, err := t.client.CreateTrafficMatchingListWithBodyWithResponse(
		ctx,
		siteID,
		"application/json",
		bytes.NewReader(args),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create traffic matching list: %w",
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
		"Traffic matching list created:\n%s",
		formatTrafficMatchingList(resp.JSON201),
	), nil
}

// --- update_traffic_matching_list ---

// UpdateTrafficMatchingList implements the update_traffic_matching_list
// MCP tool.
type UpdateTrafficMatchingList struct {
	baseTool
}

// NewUpdateTrafficMatchingList creates a new UpdateTrafficMatchingList
// tool.
func NewUpdateTrafficMatchingList(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *UpdateTrafficMatchingList {
	return &UpdateTrafficMatchingList{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *UpdateTrafficMatchingList) Description() string {
	return "Update an existing traffic matching list"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *UpdateTrafficMatchingList) InputSchema() map[string]interface{} {
	props := trafficMatchingListInputSchema()
	props["trafficMatchingListId"] = map[string]interface{}{
		"type":        "string",
		"description": "Traffic matching list UUID to update",
	}
	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required": []string{
			"trafficMatchingListId",
			"name",
			"type",
		},
	}
}

// Execute runs the tool.
func (t *UpdateTrafficMatchingList) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID                string `json:"siteId"`
		TrafficMatchingListID string `json:"trafficMatchingListId"`
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

	listID, err := resolveUUID(
		"trafficMatchingListId",
		params.TrafficMatchingListID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.UpdateTrafficMatchingListWithBodyWithResponse(
		ctx,
		siteID,
		listID,
		"application/json",
		bytes.NewReader(args),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to update traffic matching list: %w",
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
		"Traffic matching list updated:\n%s",
		formatTrafficMatchingList(resp.JSON200),
	), nil
}

// --- delete_traffic_matching_list ---

// DeleteTrafficMatchingList implements the delete_traffic_matching_list
// MCP tool.
type DeleteTrafficMatchingList struct {
	baseTool
}

// NewDeleteTrafficMatchingList creates a new DeleteTrafficMatchingList
// tool.
func NewDeleteTrafficMatchingList(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *DeleteTrafficMatchingList {
	return &DeleteTrafficMatchingList{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *DeleteTrafficMatchingList) Description() string {
	return "Delete a traffic matching list"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *DeleteTrafficMatchingList) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"trafficMatchingListId": map[string]interface{}{
				"type":        "string",
				"description": "Traffic matching list UUID to delete",
			},
		},
		"required": []string{"trafficMatchingListId"},
	}
}

// Execute runs the tool.
func (t *DeleteTrafficMatchingList) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID                string `json:"siteId"`
		TrafficMatchingListID string `json:"trafficMatchingListId"`
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

	listID, err := resolveUUID(
		"trafficMatchingListId",
		params.TrafficMatchingListID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.DeleteTrafficMatchingListWithResponse(
		ctx,
		siteID,
		listID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete traffic matching list: %w",
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
		"Traffic matching list %s deleted successfully.",
		listID.String(),
	), nil
}
