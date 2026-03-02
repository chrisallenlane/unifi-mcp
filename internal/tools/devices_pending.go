package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

// ListPendingDevices implements the list_pending_devices MCP tool.
// This tool is not site-scoped.
type ListPendingDevices struct {
	baseTool
}

// NewListPendingDevices creates a new ListPendingDevices tool.
func NewListPendingDevices(
	c *unifi.ClientWithResponses,
) *ListPendingDevices {
	return &ListPendingDevices{baseTool{client: c}}
}

// Description returns a description of the tool.
func (t *ListPendingDevices) Description() string {
	return "List devices pending adoption (not site-scoped)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListPendingDevices) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": paginationSchema(),
	}
}

// Execute runs the tool.
func (t *ListPendingDevices) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		Limit  *int32 `json:"limit"`
		Offset *int32 `json:"offset"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	resp, err := t.client.GetPendingDevicePageWithResponse(
		ctx,
		&unifi.GetPendingDevicePageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list pending devices: %w",
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
		return "No pending devices found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Pending Devices (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, d := range page.Data {
		fw := ""
		if d.FirmwareVersion != nil {
			fw = *d.FirmwareVersion
		}
		fmt.Fprintf(
			&b,
			"%d. Model: %s\n   MAC: %s\n   IP: %s\n"+
				"   State: %s\n   Firmware: %s\n",
			i+1,
			d.Model,
			d.MacAddress,
			d.IpAddress,
			d.State,
			fw,
		)
	}

	return b.String(), nil
}
