# unifi-mcp-server

An [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) server for
the [UniFi Network Integration API](https://developer.ui.com/). Exposes UniFi
firewall management as MCP tools usable by Claude and other AI assistants.

## Features

- Query controller info and sites
- Full CRUD for firewall zones and firewall policies
- Firewall policy ordering
- Type-safe client generated from the official UniFi OpenAPI 3.1.0 spec

## Requirements

- Go 1.24+
- UniFi controller running firmware ~7.x/8.x or later (Integration API)
- An API key from the UniFi controller

## Installation

```bash
git clone https://github.com/chrisallenlane/unifi-mcp-server.git
cd unifi-mcp-server
make build
# Binary: dist/unifi-mcp-server
```

## Configuration

| Variable | Required | Description |
|---|---|---|
| `UNIFI_API_URL` | Yes | Base URL of your UniFi controller (e.g. `https://192.168.1.1`) |
| `UNIFI_API_KEY` | Yes | API key from the UniFi controller |
| `UNIFI_SITE_ID` | No | Default site UUID; tools accept `siteId` to override |
| `UNIFI_INSECURE` | No | Set to any non-empty value to skip TLS verification |

## Tools

| Tool | Description |
|---|---|
| `get_info` | Get controller application version |
| `list_sites` | List all sites |
| `list_firewall_zones` | List firewall zones for a site |
| `get_firewall_zone` | Get a specific firewall zone |
| `create_firewall_zone` | Create a new firewall zone |
| `update_firewall_zone` | Update an existing firewall zone |
| `delete_firewall_zone` | Delete a firewall zone |
| `list_firewall_policies` | List firewall policies for a site |
| `get_firewall_policy` | Get a specific firewall policy |
| `create_firewall_policy` | Create a new firewall policy |
| `update_firewall_policy` | Update an existing firewall policy |
| `delete_firewall_policy` | Delete a firewall policy |
| `patch_firewall_policy` | Partially update a firewall policy |
| `get_firewall_policy_ordering` | Get firewall policy ordering for a site |
| `update_firewall_policy_ordering` | Update firewall policy ordering for a site |

## Usage

See [SETUP.md](SETUP.md) for configuration with Claude Code or Claude Desktop.

## Development

```bash
# Format, lint, vet, and test
make check

# Build binary
make build

# Regenerate types/client from OpenAPI spec
make generate
```

For development guidance, see [CLAUDE.md](CLAUDE.md).

## Architecture

Types and HTTP client are generated from the official UniFi Network OpenAPI
3.1.0 spec (`api/unifi-network.json`) using `oapi-codegen`. See
[adr/001-generate-models-from-openapi-spec.md](adr/001-generate-models-from-openapi-spec.md)
for rationale.

## License

MIT License - see LICENSE file for details
