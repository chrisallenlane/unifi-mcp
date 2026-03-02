package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp-server/internal/unifi"
)

// --- list_clients ---

// ListClients implements the list_clients MCP tool.
type ListClients struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListClients creates a new ListClients tool.
func NewListClients(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListClients {
	return &ListClients{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListClients) Description() string {
	return "List connected clients (wired, wireless, VPN, Teleport) for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListClients) InputSchema() map[string]interface{} {
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
func (t *ListClients) Execute(
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

	resp, err := t.client.GetConnectedClientOverviewPageWithResponse(
		ctx,
		siteID,
		&unifi.GetConnectedClientOverviewPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list clients: %w",
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
		return "No connected clients found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Connected Clients (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, c := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Type: %s\n",
			i+1,
			c.Name,
			c.Id.String(),
			c.Type,
		)
		if c.IpAddress != nil {
			fmt.Fprintf(&b, "   IP: %s\n", *c.IpAddress)
		}
		if c.ConnectedAt != nil {
			fmt.Fprintf(
				&b,
				"   Connected At: %s\n",
				c.ConnectedAt.String(),
			)
		}
	}

	return b.String(), nil
}

// --- get_client ---

// GetClient implements the get_client MCP tool.
type GetClient struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewGetClient creates a new GetClient tool.
func NewGetClient(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetClient {
	return &GetClient{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *GetClient) Description() string {
	return "Get detailed information about a connected client"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetClient) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"clientId": map[string]interface{}{
				"type":        "string",
				"description": "Client UUID",
			},
		},
		"required": []string{"clientId"},
	}
}

// Execute runs the tool.
func (t *GetClient) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID   string `json:"siteId"`
		ClientID string `json:"clientId"`
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

	clientID, err := resolveUUID("clientId", params.ClientID)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetConnectedClientDetailsWithResponse(
		ctx,
		siteID,
		clientID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get client: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatClientDetails(resp.JSON200), nil
}

func formatClientDetails(
	c *unifi.ClientDetails,
) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Name: %s\n", c.Name)
	fmt.Fprintf(&b, "ID: %s\n", c.Id.String())
	fmt.Fprintf(&b, "Type: %s\n", c.Type)
	if c.IpAddress != nil {
		fmt.Fprintf(&b, "IP Address: %s\n", *c.IpAddress)
	}
	if c.ConnectedAt != nil {
		fmt.Fprintf(
			&b,
			"Connected At: %s\n",
			c.ConnectedAt.String(),
		)
	}
	if c.Access != nil {
		accessJSON, err := json.MarshalIndent(
			c.Access,
			"",
			"  ",
		)
		if err == nil {
			fmt.Fprintf(
				&b,
				"Access: %s\n",
				string(accessJSON),
			)
		}
	}
	return b.String()
}

// --- execute_client_action ---

// ExecuteClientAction implements the execute_client_action MCP tool.
type ExecuteClientAction struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewExecuteClientAction creates a new ExecuteClientAction tool.
func NewExecuteClientAction(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ExecuteClientAction {
	return &ExecuteClientAction{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ExecuteClientAction) Description() string {
	return "Execute an action on a connected client (authorize/unauthorize guest access)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ExecuteClientAction) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"clientId": map[string]interface{}{
				"type":        "string",
				"description": "Client UUID",
			},
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action to perform (AUTHORIZE_GUEST_ACCESS or UNAUTHORIZE_GUEST_ACCESS)",
				"enum": []string{
					"AUTHORIZE_GUEST_ACCESS",
					"UNAUTHORIZE_GUEST_ACCESS",
				},
			},
		},
		"required": []string{"clientId", "action"},
	}
}

// Execute runs the tool.
func (t *ExecuteClientAction) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID   string `json:"siteId"`
		ClientID string `json:"clientId"`
		Action   string `json:"action"`
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

	clientID, err := resolveUUID("clientId", params.ClientID)
	if err != nil {
		return "", err
	}

	if params.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	resp, err := t.client.ExecuteConnectedClientActionWithResponse(
		ctx,
		siteID,
		clientID,
		unifi.ClientActionRequest{
			Action: params.Action,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute client action: %w",
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
		"Action executed: %s",
		resp.JSON200.Action,
	), nil
}
