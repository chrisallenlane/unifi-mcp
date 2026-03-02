package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chrisallenlane/unifi-mcp-server/internal/unifi"
)

// --- list_wans ---

// ListWans implements the list_wans MCP tool.
type ListWans struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListWans creates a new ListWans tool.
func NewListWans(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListWans {
	return &ListWans{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListWans) Description() string {
	return "List WAN interfaces for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListWans) InputSchema() map[string]interface{} {
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
func (t *ListWans) Execute(
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

	resp, err := t.client.GetWansOverviewPageWithResponse(
		ctx,
		siteID,
		&unifi.GetWansOverviewPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to list WANs: %w", err)
	}

	if resp.JSON200 == nil {
		return "", unexpectedStatusError(
			resp.StatusCode(),
			resp.Body,
		)
	}

	page := resp.JSON200
	if len(page.Data) == 0 {
		return "No WANs found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"WANs (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, wan := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n",
			i+1,
			wan.Name,
			wan.Id.String(),
		)
	}

	return b.String(), nil
}

// --- list_vpn_tunnels ---

// ListVpnTunnels implements the list_vpn_tunnels MCP tool.
type ListVpnTunnels struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListVpnTunnels creates a new ListVpnTunnels tool.
func NewListVpnTunnels(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListVpnTunnels {
	return &ListVpnTunnels{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListVpnTunnels) Description() string {
	return "List site-to-site VPN tunnels for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListVpnTunnels) InputSchema() map[string]interface{} {
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
func (t *ListVpnTunnels) Execute(
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

	resp, err := t.client.GetSiteToSiteVpnTunnelPageWithResponse(
		ctx,
		siteID,
		&unifi.GetSiteToSiteVpnTunnelPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list VPN tunnels: %w",
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
		return "No VPN tunnels found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"VPN Tunnels (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, tunnel := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Type: %s\n",
			i+1,
			tunnel.Name,
			tunnel.Id.String(),
			tunnel.Type,
		)
	}

	return b.String(), nil
}

// --- list_vpn_servers ---

// ListVpnServers implements the list_vpn_servers MCP tool.
type ListVpnServers struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListVpnServers creates a new ListVpnServers tool.
func NewListVpnServers(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListVpnServers {
	return &ListVpnServers{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListVpnServers) Description() string {
	return "List VPN servers for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListVpnServers) InputSchema() map[string]interface{} {
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
func (t *ListVpnServers) Execute(
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

	resp, err := t.client.GetVpnServerPageWithResponse(
		ctx,
		siteID,
		&unifi.GetVpnServerPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list VPN servers: %w",
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
		return "No VPN servers found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"VPN Servers (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, server := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Type: %s\n   Enabled: %t\n",
			i+1,
			server.Name,
			server.Id.String(),
			server.Type,
			server.Enabled,
		)
	}

	return b.String(), nil
}

// --- list_radius_profiles ---

// ListRadiusProfiles implements the list_radius_profiles MCP tool.
type ListRadiusProfiles struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListRadiusProfiles creates a new ListRadiusProfiles tool.
func NewListRadiusProfiles(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListRadiusProfiles {
	return &ListRadiusProfiles{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListRadiusProfiles) Description() string {
	return "List RADIUS profiles for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListRadiusProfiles) InputSchema() map[string]interface{} {
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
func (t *ListRadiusProfiles) Execute(
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

	resp, err := t.client.GetRadiusProfileOverviewPageWithResponse(
		ctx,
		siteID,
		&unifi.GetRadiusProfileOverviewPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list RADIUS profiles: %w",
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
		return "No RADIUS profiles found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"RADIUS Profiles (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, profile := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n",
			i+1,
			profile.Name,
			profile.Id.String(),
		)
	}

	return b.String(), nil
}

// --- list_device_tags ---

// ListDeviceTags implements the list_device_tags MCP tool.
type ListDeviceTags struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// NewListDeviceTags creates a new ListDeviceTags tool.
func NewListDeviceTags(
	c *unifi.ClientWithResponses,
	defaultSiteID string,
) *ListDeviceTags {
	return &ListDeviceTags{
		client:        c,
		defaultSiteID: defaultSiteID,
	}
}

// Description returns a description of the tool.
func (t *ListDeviceTags) Description() string {
	return "List device tags for a site"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListDeviceTags) InputSchema() map[string]interface{} {
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
func (t *ListDeviceTags) Execute(
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

	resp, err := t.client.GetDeviceTagPageWithResponse(
		ctx,
		siteID,
		&unifi.GetDeviceTagPageParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list device tags: %w",
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
		return "No device tags found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Device Tags (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, tag := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s\n   ID: %s\n   Devices: %d\n",
			i+1,
			tag.Name,
			tag.Id.String(),
			len(tag.DeviceIds),
		)
	}

	return b.String(), nil
}

// --- list_dpi_categories ---

// ListDpiCategories implements the list_dpi_categories MCP tool.
type ListDpiCategories struct {
	client *unifi.ClientWithResponses
}

// NewListDpiCategories creates a new ListDpiCategories tool.
func NewListDpiCategories(
	c *unifi.ClientWithResponses,
) *ListDpiCategories {
	return &ListDpiCategories{client: c}
}

// Description returns a description of the tool.
func (t *ListDpiCategories) Description() string {
	return "List DPI application categories (for use in firewall policy traffic filters)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListDpiCategories) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": paginationSchema(),
	}
}

// Execute runs the tool.
func (t *ListDpiCategories) Execute(
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

	resp, err := t.client.GetDpiApplicationCategoriesWithResponse(
		ctx,
		&unifi.GetDpiApplicationCategoriesParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list DPI categories: %w",
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
		return "No DPI categories found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"DPI Categories (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, cat := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s (ID: %d)\n",
			i+1,
			cat.Name,
			cat.Id,
		)
	}

	return b.String(), nil
}

// --- list_dpi_applications ---

// ListDpiApplications implements the list_dpi_applications MCP tool.
type ListDpiApplications struct {
	client *unifi.ClientWithResponses
}

// NewListDpiApplications creates a new ListDpiApplications tool.
func NewListDpiApplications(
	c *unifi.ClientWithResponses,
) *ListDpiApplications {
	return &ListDpiApplications{client: c}
}

// Description returns a description of the tool.
func (t *ListDpiApplications) Description() string {
	return "List DPI applications (for use in firewall policy traffic filters)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListDpiApplications) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": paginationSchema(),
	}
}

// Execute runs the tool.
func (t *ListDpiApplications) Execute(
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

	resp, err := t.client.GetDpiApplicationsWithResponse(
		ctx,
		&unifi.GetDpiApplicationsParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list DPI applications: %w",
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
		return "No DPI applications found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"DPI Applications (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, app := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s (ID: %d)\n",
			i+1,
			app.Name,
			app.Id,
		)
	}

	return b.String(), nil
}

// --- list_countries ---

// ListCountries implements the list_countries MCP tool.
type ListCountries struct {
	client *unifi.ClientWithResponses
}

// NewListCountries creates a new ListCountries tool.
func NewListCountries(
	c *unifi.ClientWithResponses,
) *ListCountries {
	return &ListCountries{client: c}
}

// Description returns a description of the tool.
func (t *ListCountries) Description() string {
	return "List country codes (for use in region-based firewall filtering)"
}

// InputSchema returns the JSON schema for the tool's input.
func (t *ListCountries) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": paginationSchema(),
	}
}

// Execute runs the tool.
func (t *ListCountries) Execute(
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

	resp, err := t.client.GetCountriesWithResponse(
		ctx,
		&unifi.GetCountriesParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to list countries: %w",
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
		return "No countries found.", nil
	}

	var b strings.Builder
	fmt.Fprintf(
		&b,
		"Countries (%d of %d):\n",
		len(page.Data),
		page.TotalCount,
	)
	for i, country := range page.Data {
		fmt.Fprintf(
			&b,
			"%d. %s (%s)\n",
			i+1,
			country.Name,
			country.Code,
		)
	}

	return b.String(), nil
}
