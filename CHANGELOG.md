# Changelog

## 0.1.1 — 2026-04-19

### Fixed

- **JSON-RPC notification handling**: The server no longer responds to notifications (JSON-RPC 2.0 requests without an `id`). Previously it returned a "Method not found" error for `notifications/initialized`, which caused strict MCP clients (including current Claude Code) to abort the handshake and expose zero tools.

## 0.1.0 — 2026-03-03

Initial release.

### Features

- MCP server exposing 67 UniFi Network API operations as tools
- JSON-RPC 2.0 over stdio
- HTTP client generated from UniFi OpenAPI 3.1.0 spec (v10.1.84)
- API key authentication via `X-API-Key` header
- Optional TLS verification skip for self-signed certificates
- Default site ID via environment variable, per-request override via `siteId`
  parameter

### Tools

- **Info / Sites (2):** `get_info`, `list_sites`
- **Firewall Zones (5):** list, get, create, update, delete
- **Firewall Policies (8):** list, get, create, update, delete, patch, get
  ordering, update ordering
- **Networks (6):** list, get, create, update, delete, get references
- **Clients (3):** list, get, execute action
- **Devices (8):** list, get, adopt, remove, execute device action, execute
  port action, get statistics, list pending
- **ACL Rules (7):** list, get, create, update, delete, get ordering, update
  ordering
- **DNS Policies (5):** list, get, create, update, delete
- **Traffic Matching Lists (5):** list, get, create, update, delete
- **WiFi Broadcasts (5):** list, get, create, update, delete
- **Hotspot Vouchers (5):** list, get, create, delete (bulk), delete (single)
- **Supporting Read-Only (8):** WANs, VPN tunnels, VPN servers, RADIUS
  profiles, device tags, DPI categories, DPI applications, countries
