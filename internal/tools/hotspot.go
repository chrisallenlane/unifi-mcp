package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

// --- list_vouchers ---

// ListVouchers implements the list_vouchers MCP tool.
type ListVouchers struct {
	baseTool
}

// NewListVouchers creates a new ListVouchers tool.
func NewListVouchers(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListVouchers {
	return &ListVouchers{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *ListVouchers) Description() string {
	return "List hotspot vouchers for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListVouchers) InputSchema() map[string]interface{} {
	return listSchema()
}

// Execute runs the tool.
func (t *ListVouchers) Execute(
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

	resp, err := t.client.GetVouchersWithResponse(
		ctx,
		siteID,
		&unifi.GetVouchersParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list vouchers: %w",
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
		return "No vouchers found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Vouchers (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, v := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. Code: %s\n   ID: %s\n   Name: %s\n   Duration: %d min\n   Expired: %t\n",
			i+1,
			v.Code,
			v.Id.String(),
			v.Name,
			v.TimeLimitMinutes,
			v.Expired,
		)
	}

	return b.String(), nil
}

// --- get_voucher ---

// GetVoucher implements the get_voucher MCP tool.
type GetVoucher struct {
	baseTool
}

// NewGetVoucher creates a new GetVoucher tool.
func NewGetVoucher(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *GetVoucher {
	return &GetVoucher{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *GetVoucher) Description() string {
	return "Get details of a specific hotspot voucher"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *GetVoucher) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"voucherId": map[string]interface{}{
				"type":        "string",
				"description": "Voucher UUID",
			},
		},
		"required": []string{"voucherId"},
	}
}

// Execute runs the tool.
func (t *GetVoucher) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID    string `json:"siteId"`
		VoucherID string `json:"voucherId"`
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

	voucherID, err := resolveUUID(
		"voucherId",
		params.VoucherID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.GetVoucherWithResponse(
		ctx,
		siteID,
		voucherID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to get voucher: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return formatVoucher(resp.JSON200), nil
}

func formatVoucher(v *unifi.HotspotVoucherDetails) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Code: %s\n", v.Code)
	fmt.Fprintf(&b, "ID: %s\n", v.Id.String())
	fmt.Fprintf(&b, "Name: %s\n", v.Name)
	fmt.Fprintf(
		&b,
		"Duration: %d minutes\n",
		v.TimeLimitMinutes,
	)
	fmt.Fprintf(&b, "Expired: %t\n", v.Expired)
	fmt.Fprintf(
		&b,
		"Authorized Guests: %d\n",
		v.AuthorizedGuestCount,
	)
	if v.AuthorizedGuestLimit != nil {
		fmt.Fprintf(
			&b,
			"Guest Limit: %d\n",
			*v.AuthorizedGuestLimit,
		)
	}
	fmt.Fprintf(
		&b,
		"Created At: %s\n",
		v.CreatedAt.String(),
	)
	if v.ActivatedAt != nil {
		fmt.Fprintf(
			&b,
			"Activated At: %s\n",
			v.ActivatedAt.String(),
		)
	}
	if v.ExpiresAt != nil {
		fmt.Fprintf(
			&b,
			"Expires At: %s\n",
			v.ExpiresAt.String(),
		)
	}
	if v.DataUsageLimitMBytes != nil {
		fmt.Fprintf(
			&b,
			"Data Limit: %d MB\n",
			*v.DataUsageLimitMBytes,
		)
	}
	if v.RxRateLimitKbps != nil {
		fmt.Fprintf(
			&b,
			"Download Limit: %d kbps\n",
			*v.RxRateLimitKbps,
		)
	}
	if v.TxRateLimitKbps != nil {
		fmt.Fprintf(
			&b,
			"Upload Limit: %d kbps\n",
			*v.TxRateLimitKbps,
		)
	}
	return b.String()
}

// --- create_vouchers ---

// CreateVouchers implements the create_vouchers MCP tool.
type CreateVouchers struct {
	baseTool
}

// NewCreateVouchers creates a new CreateVouchers tool.
func NewCreateVouchers(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *CreateVouchers {
	return &CreateVouchers{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *CreateVouchers) Description() string {
	return "Generate hotspot vouchers (creates multiple vouchers at once)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *CreateVouchers) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Voucher note (applied to all generated vouchers)",
			},
			"timeLimitMinutes": map[string]interface{}{
				"type":        "integer",
				"description": "Duration in minutes for network access",
			},
			"count": map[string]interface{}{
				"type":        "integer",
				"description": "Number of vouchers to generate",
			},
			"authorizedGuestLimit": map[string]interface{}{
				"type":        "integer",
				"description": "Max guests per voucher (optional)",
			},
			"dataUsageLimitMBytes": map[string]interface{}{
				"type":        "integer",
				"description": "Data usage limit in MB (optional)",
			},
			"rxRateLimitKbps": map[string]interface{}{
				"type":        "integer",
				"description": "Download rate limit in kbps (optional)",
			},
			"txRateLimitKbps": map[string]interface{}{
				"type":        "integer",
				"description": "Upload rate limit in kbps (optional)",
			},
		},
		"required": []string{
			"name",
			"timeLimitMinutes",
		},
	}
}

// Execute runs the tool.
func (t *CreateVouchers) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID               string `json:"siteId"`
		Name                 string `json:"name"`
		TimeLimitMinutes     int64  `json:"timeLimitMinutes"`
		Count                *int32 `json:"count"`
		AuthorizedGuestLimit *int64 `json:"authorizedGuestLimit"`
		DataUsageLimitMBytes *int64 `json:"dataUsageLimitMBytes"`
		RxRateLimitKbps      *int64 `json:"rxRateLimitKbps"`
		TxRateLimitKbps      *int64 `json:"txRateLimitKbps"`
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

	body := unifi.HotspotVoucherCreationRequest{
		Name:                 params.Name,
		TimeLimitMinutes:     params.TimeLimitMinutes,
		Count:                params.Count,
		AuthorizedGuestLimit: params.AuthorizedGuestLimit,
		DataUsageLimitMBytes: params.DataUsageLimitMBytes,
		RxRateLimitKbps:      params.RxRateLimitKbps,
		TxRateLimitKbps:      params.TxRateLimitKbps,
	}

	resp, err := t.client.CreateVouchersWithResponse(
		ctx,
		siteID,
		body,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to create vouchers: %w",
			err,
		)
	}

	if resp.JSON201 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	result := resp.JSON201
	if result.Vouchers == nil || len(*result.Vouchers) == 0 {
		return "No vouchers generated.", nil
	}

	vouchers := *result.Vouchers
	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Generated %d voucher(s):\n",
		len(vouchers),
	)
	for i, v := range vouchers {
		fmt.Fprintf(
			&b,
			"%d. Code: %s  ID: %s\n",
			i+1,
			v.Code,
			v.Id.String(),
		)
	}

	return b.String(), nil
}

// --- delete_vouchers ---

// DeleteVouchers implements the delete_vouchers MCP tool.
type DeleteVouchers struct {
	baseTool
}

// NewDeleteVouchers creates a new DeleteVouchers tool.
func NewDeleteVouchers(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *DeleteVouchers {
	return &DeleteVouchers{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *DeleteVouchers) Description() string {
	return "Bulk delete hotspot vouchers matching a filter"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *DeleteVouchers) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"filter": map[string]interface{}{
				"type":        "string",
				"description": "Filter expression for bulk deletion",
			},
		},
		"required": []string{"filter"},
	}
}

// Execute runs the tool.
func (t *DeleteVouchers) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID string `json:"siteId"`
		Filter string `json:"filter"`
	}
	if err := parseArgs(args, &params); err != nil {
		return "", err
	}

	if params.Filter == "" {
		return "", fmt.Errorf("filter is required")
	}

	siteID, err := resolveSiteID(
		params.SiteID,
		t.defaultSiteID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.DeleteVouchersWithResponse(
		ctx,
		siteID,
		&unifi.DeleteVouchersParams{
			Filter: params.Filter,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete vouchers: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	deleted := int64(0)
	if resp.JSON200.VouchersDeleted != nil {
		deleted = *resp.JSON200.VouchersDeleted
	}

	return fmt.Sprintf(
		"Deleted %d voucher(s).",
		deleted,
	), nil
}

// --- delete_voucher ---

// DeleteVoucher implements the delete_voucher MCP tool.
type DeleteVoucher struct {
	baseTool
}

// NewDeleteVoucher creates a new DeleteVoucher tool.
func NewDeleteVoucher(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *DeleteVoucher {
	return &DeleteVoucher{baseTool{c, defaultSiteID}}
}

// Description returns a description of the tool.
func (t *DeleteVoucher) Description() string {
	return "Delete a specific hotspot voucher"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *DeleteVoucher) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			"voucherId": map[string]interface{}{
				"type":        "string",
				"description": "Voucher UUID",
			},
		},
		"required": []string{"voucherId"},
	}
}

// Execute runs the tool.
func (t *DeleteVoucher) Execute(
	ctx context.Context,
	args json.RawMessage,
) (string, error) {
	var params struct {
		SiteID    string `json:"siteId"`
		VoucherID string `json:"voucherId"`
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

	voucherID, err := resolveUUID(
		"voucherId",
		params.VoucherID,
	)
	if err != nil {
		return "", err
	}

	resp, err := t.client.DeleteVoucherWithResponse(
		ctx,
		siteID,
		voucherID,
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to delete voucher: %w",
			err,
		)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	return "Voucher deleted.", nil
}
