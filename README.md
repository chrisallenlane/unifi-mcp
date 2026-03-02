# unifi-mcp-server

An [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) server for
the [UniFi Network Integration API](https://developer.ui.com/). Exposes 67
UniFi controller operations as MCP tools usable by Claude and other AI
assistants.

## Features

- Query controller info and sites
- Full CRUD for firewall zones, firewall policies, networks, ACL rules, DNS
  policies, traffic matching lists, and WiFi broadcasts
- Client and device management (listing, status, actions)
- Hotspot voucher management
- Supporting read-only tools (WANs, VPN tunnels, VPN servers, RADIUS profiles,
  device tags, DPI categories, DPI applications, countries)
- Firewall policy and ACL rule ordering
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

**Info / Sites**

| Tool | Description |
|---|---|
| `get_info` | Get controller application version |
| `list_sites` | List all sites |

**Firewall Zones**

| Tool | Description |
|---|---|
| `list_firewall_zones` | List firewall zones for a site |
| `get_firewall_zone` | Get a specific firewall zone |
| `create_firewall_zone` | Create a new firewall zone |
| `update_firewall_zone` | Update an existing firewall zone |
| `delete_firewall_zone` | Delete a firewall zone |

**Firewall Policies**

| Tool | Description |
|---|---|
| `list_firewall_policies` | List firewall policies for a site |
| `get_firewall_policy` | Get a specific firewall policy |
| `create_firewall_policy` | Create a new firewall policy |
| `update_firewall_policy` | Update an existing firewall policy |
| `delete_firewall_policy` | Delete a firewall policy |
| `patch_firewall_policy` | Partially update a firewall policy |
| `get_firewall_policy_ordering` | Get firewall policy ordering for a site |
| `update_firewall_policy_ordering` | Update firewall policy ordering for a site |

**Networks**

| Tool | Description |
|---|---|
| `list_networks` | List networks for a site |
| `get_network` | Get a specific network |
| `create_network` | Create a new network |
| `update_network` | Update an existing network |
| `delete_network` | Delete a network |
| `get_network_references` | Get references to a network |

**Clients**

| Tool | Description |
|---|---|
| `list_clients` | List clients for a site |
| `get_client` | Get a specific client |
| `execute_client_action` | Execute an action on a client |

**Devices**

| Tool | Description |
|---|---|
| `list_devices` | List devices for a site |
| `get_device` | Get a specific device |
| `adopt_device` | Adopt a device |
| `remove_device` | Remove a device |
| `execute_device_action` | Execute an action on a device |
| `execute_port_action` | Execute an action on a device port |
| `get_device_statistics` | Get statistics for a device |
| `list_pending_devices` | List devices pending adoption |

**ACL Rules**

| Tool | Description |
|---|---|
| `list_acl_rules` | List ACL rules for a site |
| `get_acl_rule` | Get a specific ACL rule |
| `create_acl_rule` | Create a new ACL rule |
| `update_acl_rule` | Update an existing ACL rule |
| `delete_acl_rule` | Delete an ACL rule |
| `get_acl_rule_ordering` | Get ACL rule ordering for a site |
| `update_acl_rule_ordering` | Update ACL rule ordering for a site |

**DNS Policies**

| Tool | Description |
|---|---|
| `list_dns_policies` | List DNS policies for a site |
| `get_dns_policy` | Get a specific DNS policy |
| `create_dns_policy` | Create a new DNS policy |
| `update_dns_policy` | Update an existing DNS policy |
| `delete_dns_policy` | Delete a DNS policy |

**Traffic Matching Lists**

| Tool | Description |
|---|---|
| `list_traffic_matching_lists` | List traffic matching lists for a site |
| `get_traffic_matching_list` | Get a specific traffic matching list |
| `create_traffic_matching_list` | Create a new traffic matching list |
| `update_traffic_matching_list` | Update an existing traffic matching list |
| `delete_traffic_matching_list` | Delete a traffic matching list |

**WiFi Broadcasts**

| Tool | Description |
|---|---|
| `list_wifi_broadcasts` | List WiFi broadcasts for a site |
| `get_wifi_broadcast` | Get a specific WiFi broadcast |
| `create_wifi_broadcast` | Create a new WiFi broadcast |
| `update_wifi_broadcast` | Update an existing WiFi broadcast |
| `delete_wifi_broadcast` | Delete a WiFi broadcast |

**Hotspot Vouchers**

| Tool | Description |
|---|---|
| `list_vouchers` | List hotspot vouchers for a site |
| `get_voucher` | Get a specific hotspot voucher |
| `create_vouchers` | Create hotspot vouchers |
| `delete_vouchers` | Delete multiple hotspot vouchers |
| `delete_voucher` | Delete a specific hotspot voucher |

**Supporting Read-Only**

| Tool | Description |
|---|---|
| `list_wans` | List WAN configurations for a site |
| `list_vpn_tunnels` | List VPN tunnels for a site |
| `list_vpn_servers` | List VPN servers for a site |
| `list_radius_profiles` | List RADIUS profiles for a site |
| `list_device_tags` | List device tags for a site |
| `list_dpi_categories` | List DPI categories |
| `list_dpi_applications` | List DPI applications |
| `list_countries` | List available countries |

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
