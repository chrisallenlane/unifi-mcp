# CLAUDE.md

This file provides guidance to Claude Code when working in this repository.

## Project Overview

**unifi-mcp** is an MCP (Model Context Protocol) server for the UniFi
Network Integration API. It exposes 67 UniFi controller operations as MCP
tools that Claude and other AI assistants can call, covering firewall
management, networks, clients, devices, ACL rules, DNS policies, traffic
matching lists, WiFi broadcasts, hotspot vouchers, and supporting read-only
endpoints.

**Tech Stack:**
- **Language**: Go 1.24+
- **Protocol**: MCP via JSON-RPC 2.0 over stdio
- **HTTP Client**: Generated from UniFi OpenAPI 3.1.0 spec via `oapi-codegen`
- **Module**: `github.com/chrisallenlane/unifi-mcp`

## Project Structure

```
unifi-mcp/
├── api/
│   └── unifi-network.json           # UniFi OpenAPI 3.1.0 spec (v10.1.84)
├── adr/
│   └── 001-generate-models-from-openapi-spec.md
├── cmd/
│   └── unifi-mcp/
│       └── main.go                  # Entry point; reads env vars, wires client
├── internal/
│   ├── server/
│   │   ├── server.go                # JSON-RPC server, request routing, tool registry
│   │   ├── server_test.go
│   │   ├── server_fuzz_test.go      # Fuzz tests for Run, handleCallTool
│   │   └── types.go                 # JSON-RPC request/response types
│   ├── tools/
│   │   ├── tool.go                  # Tool interface
│   │   ├── helpers.go               # baseTool struct, UUID helpers, siteId resolution, error helpers, shared schema builders
│   │   ├── helpers_test.go
│   │   ├── helpers_fuzz_test.go     # Fuzz tests for parseArgs, resolveUUID
│   │   ├── test_helpers_test.go     # Shared test helpers (testClient, testSiteID)
│   │   ├── info.go                  # get_info tool
│   │   ├── info_test.go
│   │   ├── sites.go                 # list_sites tool
│   │   ├── sites_test.go
│   │   ├── firewall_zones.go        # 5 firewall zone tools
│   │   ├── firewall_zones_test.go
│   │   ├── firewall_policies.go     # 6 firewall policy CRUD tools
│   │   ├── firewall_policies_test.go
│   │   ├── firewall_policy_ordering.go     # 2 firewall policy ordering tools
│   │   ├── firewall_policy_ordering_test.go
│   │   ├── networks.go              # 6 network tools
│   │   ├── networks_test.go
│   │   ├── clients.go               # 3 client tools
│   │   ├── clients_test.go
│   │   ├── devices.go               # 7 device tools (adopt, remove, get, list, actions, port, stats)
│   │   ├── devices_test.go
│   │   ├── devices_pending.go       # list_pending_devices tool
│   │   ├── devices_pending_test.go
│   │   ├── acl_rules.go             # 7 ACL rule tools
│   │   ├── acl_rules_test.go
│   │   ├── dns_policies.go          # 5 DNS policy tools
│   │   ├── dns_policies_test.go
│   │   ├── traffic_matching_lists.go # 5 traffic matching list tools
│   │   ├── traffic_matching_lists_test.go
│   │   ├── wifi.go                  # 5 WiFi broadcast tools
│   │   ├── wifi_test.go
│   │   ├── hotspot.go               # 5 hotspot voucher tools
│   │   ├── hotspot_test.go
│   │   ├── supporting_network.go    # 3 network supporting tools (WANs, VPN tunnels, VPN servers)
│   │   ├── supporting_network_test.go
│   │   ├── supporting_reference.go  # 5 reference supporting tools (RADIUS, device tags, DPI, countries)
│   │   ├── supporting_reference_test.go
│   │   ├── body_fields_test.go      # Tests for siteId/ID stripping from API request bodies
│   │   └── bug17_18_test.go         # Tests for traffic matching list items and DNS policy fields
│   └── unifi/
│       ├── types.gen.go             # Generated model structs (DO NOT EDIT)
│       └── client.gen.go            # Generated HTTP client (DO NOT EDIT)
├── dist/                            # Build output
├── Makefile
├── go.mod
└── go.sum
```

## Architecture

### MCP Protocol

The server implements MCP via JSON-RPC 2.0 over stdio:

1. **Stdin** → JSON-RPC requests from Claude
2. **Server** → routes to handlers, executes tools
3. **Stdout** → JSON-RPC responses back to Claude

Key methods: `initialize`, `tools/list`, `tools/call`

### Entry Point (`cmd/unifi-mcp/main.go`)

Reads environment variables, constructs a `unifi.ClientWithResponses` with the
API key injected as an `X-API-Key` request editor, then starts the server.
`UNIFI_INSECURE` skips TLS verification for controllers with self-signed certs.

### Generated Client (`internal/unifi/`)

`oapi-codegen` generates `ClientWithResponses` from `api/unifi-network.json`.
Each endpoint becomes a typed method (e.g., `GetFirewallZonesWithResponse`).
Responses use typed fields like `resp.JSON200`; unexpected status codes use
`resp.StatusCode()` and `resp.Body`.

**Do not edit generated files.** Regenerate with `make generate`.

### Tool Interface (`internal/tools/tool.go`)

```go
type Tool interface {
    Execute(ctx context.Context, args json.RawMessage) (string, error)
    Description() string
    InputSchema() map[string]interface{}
}
```

Tools return human-readable formatted strings, not raw JSON.

### Tool Registration (`internal/server/server.go`)

Tools are registered in `registerTools()`:

```go
s.tools["tool_name"] = tools.NewToolName(s.client, s.defaultSiteID)
```

The server exposes all registered tools via `tools/list`.

### Helpers (`internal/tools/helpers.go`)

`helpers.go` also defines `baseTool`, a struct embedded by all tool implementations:

```go
type baseTool struct {
    client        *unifi.ClientWithResponses
    defaultSiteID string
}
```

Helper functions:

- `resolveSiteID(explicit, default)` - resolves site UUID from param or env default
- `resolveUUID(name, value)` - parses and validates a UUID string
- `resolveUUIDs(name, values)` - parses a slice of UUID strings
- `unexpectedStatusError(statusCode, body)` - formats an error for bad status codes
- `parseArgs(args, dst)` - unmarshals JSON-RPC args into a typed struct
- `stripKeys(args, keys...)` - removes MCP-only keys (e.g. `siteId`) from raw JSON before forwarding to API
- `siteIDSchema()` - standard JSON schema snippet for the `siteId` parameter
- `paginationSchema()` - standard JSON schema snippet for `limit` and `offset` parameters
- `siteAndIDSchema(idName, idDesc)` - schema for operations taking a siteId + one resource ID
- `listSchema()` - schema for list operations with siteId + pagination parameters

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `UNIFI_API_URL` | Yes | Base URL of the UniFi controller |
| `UNIFI_API_KEY` | Yes | API key (X-API-Key header) |
| `UNIFI_SITE_ID` | No | Default site UUID; tools accept `siteId` to override |
| `UNIFI_INSECURE` | No | Non-empty value skips TLS verification |

## Available Tools (67)

**Info / Sites (2):**
- `get_info` - controller application version
- `list_sites` - list all sites

**Firewall Zones (5):**
- `list_firewall_zones`, `get_firewall_zone`, `create_firewall_zone`,
  `update_firewall_zone`, `delete_firewall_zone`

**Firewall Policies (8):**
- `list_firewall_policies`, `get_firewall_policy`, `create_firewall_policy`,
  `update_firewall_policy`, `delete_firewall_policy`, `patch_firewall_policy`,
  `get_firewall_policy_ordering`, `update_firewall_policy_ordering`

**Networks (6):**
- `list_networks`, `get_network`, `create_network`, `update_network`,
  `delete_network`, `get_network_references`

**Clients (3):**
- `list_clients`, `get_client`, `execute_client_action`

**Devices (8):**
- `list_devices`, `get_device`, `adopt_device`, `remove_device`,
  `execute_device_action`, `execute_port_action`, `get_device_statistics`,
  `list_pending_devices`

**ACL Rules (7):**
- `list_acl_rules`, `get_acl_rule`, `create_acl_rule`, `update_acl_rule`,
  `delete_acl_rule`, `get_acl_rule_ordering`, `update_acl_rule_ordering`

**DNS Policies (5):**
- `list_dns_policies`, `get_dns_policy`, `create_dns_policy`,
  `update_dns_policy`, `delete_dns_policy`

**Traffic Matching Lists (5):**
- `list_traffic_matching_lists`, `get_traffic_matching_list`,
  `create_traffic_matching_list`, `update_traffic_matching_list`,
  `delete_traffic_matching_list`

**WiFi Broadcasts (5):**
- `list_wifi_broadcasts`, `get_wifi_broadcast`, `create_wifi_broadcast`,
  `update_wifi_broadcast`, `delete_wifi_broadcast`

**Hotspot Vouchers (5):**
- `list_vouchers`, `get_voucher`, `create_vouchers`, `delete_vouchers`,
  `delete_voucher`

**Supporting Read-Only (8):**
- `list_wans`, `list_vpn_tunnels`, `list_vpn_servers`, `list_radius_profiles`,
  `list_device_tags`, `list_dpi_categories`, `list_dpi_applications`,
  `list_countries`

## Development Workflow

```bash
make fmt           # format with golines + gofumpt (80-col wrapping)
make lint          # lint with revive (excludes internal/unifi/)
make vet           # go vet
make test          # go test ./...
make check         # fmt + lint + vet + test
make build         # fmt + lint + vet, then build to dist/unifi-mcp
make build-release # cross-compile release binaries for all platforms
make install       # build and install to $GOPATH/bin
make generate      # regenerate internal/unifi/ from api/unifi-network.json
make coverage      # test coverage report (coverage.out + coverage.html)
make fuzz          # run fuzz tests (FUZZ_TIME=10s by default)
make sloc          # count source lines of code (requires scc)
make clean         # remove compiled executables from dist/
```

## Adding a New Tool

1. **Create `internal/tools/my_tool.go`** implementing the `Tool` interface.
   Use an existing file (e.g., `info.go`, `sites.go`) as a reference.

2. **Register in `internal/server/server.go`** inside `registerTools()`:
   ```go
   s.tools["my_tool"] = tools.NewMyTool(s.client, s.defaultSiteID)
   ```

3. **Write tests** in `internal/tools/my_tool_test.go`.

4. **Run `make check`** to verify.

### Tool Conventions

- Embed `baseTool` in each tool struct (provides `client` and `defaultSiteID` via promotion)
- Accept `siteId` as an optional parameter; call `resolveSiteID` to resolve it
- Use `resolveUUID` / `resolveUUIDs` for UUID parameters
- Use `listSchema()` for list operations and `siteAndIDSchema()` for single-resource operations
- Call the generated client method (e.g., `t.client.GetXyzWithResponse(ctx, ...)`)
- Check `resp.JSON200 == nil` and call `unexpectedStatusError` on failure
- Return human-readable strings, not JSON

## Dependencies

Runtime:
- `github.com/oapi-codegen/runtime` - runtime support for generated client
- `github.com/google/uuid` - UUID parsing and formatting

Test (indirect):
- `github.com/stretchr/testify` - test assertions

Build-time only:
- `golines`, `gofumpt`, `revive` - invoked via `go run` in the Makefile, no manual install needed
- `oapi-codegen` - called directly by `make generate`; must be installed separately if you need to regenerate the client

## Version Information

- MCP Protocol Version: `2024-11-05`
- Server Name/Version: `unifi-mcp` / `0.1.1`
- OpenAPI Spec: UniFi Network v10.1.84

## Resources

- [MCP Specification](https://modelcontextprotocol.io/)
- [UniFi Developer Portal](https://developer.ui.com/)
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
