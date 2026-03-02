package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp-server/internal/unifi"
)

// ListSites implements the list_sites MCP tool.
type ListSites struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListSites creates a new ListSites tool.
func NewListSites(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListSites {
	return &ListSites{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListSites) Description() string {
	return "List UniFi sites (discover site IDs for use with other tools)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListSites) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"limit": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum number of sites to return",
			},
			"offset": map[string]interface{}{
				"type":        "integer",
				"description": "Number of sites to skip",
			},
		},
	}
}

// Execute runs the tool.
func (t *ListSites) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
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

	apiParams := &unifi.GetSiteOverviewPageParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	resp, err := t.client.GetSiteOverviewPageWithResponse(
		ctx,
		apiParams,
	)
	if err != nil {
		return "", fmt.Errorf("failed to list sites: %w", err)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	page := resp.JSON200
	if len(page.Data) == 0 {
		return "No sites found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Sites (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, site := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Internal Reference: %s\n",
			i+1,
			site.Name,
			site.Id.String(),
			site.InternalReference,
		)
	}

	return b.String(), nil
}
