# CLAUDE.md

This file provides guidance to Claude Code when working in this repository.

## Project Overview

**unifi-mcp-server** is an MCP (Model Context Protocol) server for the UniFi
Network Integration API. It wraps UniFi firewall management endpoints as MCP
tools that Claude and other AI assistants can call.

**Tech Stack:**
- **Language**: Go 1.24+
- **Protocol**: MCP via JSON-RPC 2.0 over stdio
- **HTTP Client**: Generated from UniFi OpenAPI 3.1.0 spec via `oapi-codegen`
- **Module**: `github.com/chrisallenlane/unifi-mcp-server`

## Project Structure

```
unifi-mcp-server/
├── api/
│   └── unifi-network.json       # UniFi OpenAPI 3.1.0 spec (v10.1.84)
├── adr/
│   └── 001-generate-models-from-openapi-spec.md
├── cmd/
│   └── unifi-mcp-server/
│       └── main.go              # Entry point; reads env vars, wires client
├── internal/
│   ├── server/
│   │   ├── server.go            # JSON-RPC server, request routing, tool registry
│   │   ├── server_test.go
│   │   └── types.go             # JSON-RPC request/response types
│   ├── tools/
│   │   ├── tool.go              # Tool interface
│   │   ├── helpers.go           # UUID helpers, siteId resolution, error helpers
│   │   ├── helpers_test.go
│   │   ├── info.go              # get_info tool
│   │   ├── info_test.go
│   │   ├── sites.go             # list_sites tool
│   │   ├── sites_test.go
│   │   ├── firewall_zones.go    # 5 firewall zone tools
│   │   ├── firewall_zones_test.go
│   │   ├── firewall_policies.go # 8 firewall policy tools
│   │   └── firewall_policies_test.go
│   └── unifi/
│       ├── types.gen.go         # Generated model structs (DO NOT EDIT)
│       └── client.gen.go        # Generated HTTP client (DO NOT EDIT)
├── dist/                        # Build output
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

### Entry Point (`cmd/unifi-mcp-server/main.go`)

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

- `resolveSiteID(explicit, default)` - resolves site UUID from param or env default
- `resolveUUID(name, value)` - parses and validates a UUID string
- `resolveUUIDs(name, values)` - parses a slice of UUID strings
- `unexpectedStatusError(statusCode, body)` - formats an error for bad status codes
- `siteIDSchema()` - standard JSON schema snippet for the `siteId` parameter

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `UNIFI_API_URL` | Yes | Base URL of the UniFi controller |
| `UNIFI_API_KEY` | Yes | API key (X-API-Key header) |
| `UNIFI_SITE_ID` | No | Default site UUID; tools accept `siteId` to override |
| `UNIFI_INSECURE` | No | Non-empty value skips TLS verification |

## Available Tools (15)

**Info / Sites:**
- `get_info` - controller application version
- `list_sites` - list all sites

**Firewall Zones (5):**
- `list_firewall_zones`, `get_firewall_zone`, `create_firewall_zone`,
  `update_firewall_zone`, `delete_firewall_zone`

**Firewall Policies (8):**
- `list_firewall_policies`, `get_firewall_policy`, `create_firewall_policy`,
  `update_firewall_policy`, `delete_firewall_policy`, `patch_firewall_policy`,
  `get_firewall_policy_ordering`, `update_firewall_policy_ordering`

## Development Workflow

```bash
make fmt        # format with golines + gofumpt (80-col wrapping)
make lint       # lint with revive (excludes internal/unifi/)
make vet        # go vet
make test       # go test ./...
make check      # fmt + lint + vet + test
make build      # build to dist/unifi-mcp-server
make generate   # regenerate internal/unifi/ from api/unifi-network.json
make coverage   # test coverage report
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

- Accept `siteId` as an optional parameter; call `resolveSiteID` to resolve it
- Use `resolveUUID` / `resolveUUIDs` for UUID parameters
- Call the generated client method (e.g., `t.client.GetXyzWithResponse(ctx, ...)`)
- Check `resp.JSON200 == nil` and call `unexpectedStatusError` on failure
- Return human-readable strings, not JSON

## Dependencies

- `github.com/oapi-codegen/runtime` - runtime support for generated client
- `github.com/google/uuid` - UUID parsing and formatting
- `github.com/stretchr/testify` - test assertions

Build-time only: `oapi-codegen` (for `make generate`), `golines`, `gofumpt`,
`revive` (invoked via `go run` in the Makefile, no manual install needed).

## Version Information

- MCP Protocol Version: `2024-11-05`
- Server Name/Version: `unifi-mcp-server` / `0.1.0`
- OpenAPI Spec: UniFi Network v10.1.84

## Resources

- [MCP Specification](https://modelcontextprotocol.io/)
- [UniFi Developer Portal](https://developer.ui.com/)
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
