# Setup Guide

## Build

```bash
make build
# Binary: dist/unifi-mcp
```

## Quick connectivity test

```bash
UNIFI_API_URL=https://192.168.1.1 \
UNIFI_API_KEY=your-api-key \
UNIFI_INSECURE=1 \
echo '{"jsonrpc":"2.0","id":1,"method":"initialize"}' | ./dist/unifi-mcp
```

## Configure with Claude Code

```bash
claude mcp add unifi-mcp /path/to/dist/unifi-mcp \
  -s user \
  -e UNIFI_API_URL=https://192.168.1.1 \
  -e UNIFI_API_KEY=your-api-key \
  -e UNIFI_SITE_ID=your-site-uuid \
  -e UNIFI_INSECURE=1
```

Omit `UNIFI_INSECURE` if your controller has a valid TLS certificate.
Omit `UNIFI_SITE_ID` if you prefer to pass `siteId` explicitly to each tool.

**Scope options:**
- `-s user` - Available across all projects (recommended)
- `-s local` - Private to the current project only
- `-s project` - Saved to `.mcp.json` for team sharing

**Verify:**
```bash
claude mcp list
claude mcp get unifi-mcp
```

## Configure with Claude Desktop

Add to your Claude Desktop config file:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "unifi-mcp": {
      "command": "/path/to/dist/unifi-mcp",
      "env": {
        "UNIFI_API_URL": "https://192.168.1.1",
        "UNIFI_API_KEY": "your-api-key",
        "UNIFI_SITE_ID": "your-site-uuid",
        "UNIFI_INSECURE": "1"
      }
    }
  }
}
```

Restart Claude Desktop after saving.

## Updating

After rebuilding:

```bash
make build

# Claude Code: remove and re-add
claude mcp remove unifi-mcp
claude mcp add unifi-mcp /path/to/dist/unifi-mcp \
  -s user \
  -e UNIFI_API_URL=https://192.168.1.1 \
  -e UNIFI_API_KEY=your-api-key

# Claude Desktop: restart the app
```

## Troubleshooting

**Server not appearing:**
```bash
claude mcp list
claude mcp get unifi-mcp
```

**TLS errors:** Set `UNIFI_INSECURE=1` if using a self-signed certificate.

**Authentication errors:** Verify `UNIFI_API_KEY` is correct. The Integration
API uses `X-API-Key` header authentication, not session cookies.

**Site ID errors:** Use `list_sites` to find the UUID for your site, then set
`UNIFI_SITE_ID` or pass `siteId` directly to each tool call.

**Binary not found:** Use an absolute path:
```bash
claude mcp add unifi-mcp /home/user/unifi-mcp/dist/unifi-mcp
```

## Security Notes

- Store credentials in environment variables, not in code or config files
- Use `claude mcp add` with `-e` flags rather than hardcoding secrets
- The server communicates via stdio only - no network port is opened
