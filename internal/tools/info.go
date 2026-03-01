package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/unifi-mcp-server/internal/unifi"
)

// GetInfo implements the get_info MCP tool.
type GetInfo struct {
	client *unifi.ClientWithResponses
}

// NewGetInfo creates a new GetInfo tool.
func NewGetInfo(c *unifi.ClientWithResponses) *GetInfo {
	return &GetInfo{client: c}
}

// Description returns a description of the tool.
func (t *GetInfo) Description() string {
	return "Get UniFi controller application info (version, connectivity check)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetInfo) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

// Execute runs the tool.
func (t *GetInfo) Execute(
	ctx context.Context,
	_ json.RawMessage,
) (string, error) {
	resp, err := t.client.GetInfoWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get info: %w", err)
	}

	if resp.JSON200 == nil {
		return "", fmt.Errorf(
			"unexpected status %d: %s",
			resp.StatusCode(),
			string(resp.Body),
		)
	}

	return fmt.Sprintf(
		"Application Version: %s",
		resp.JSON200.ApplicationVersion,
	), nil
}
