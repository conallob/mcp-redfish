# mcp-redfish

A local [MCP](https://modelcontextprotocol.io/) server for interacting with server BMCs via the [Redfish API](https://www.dmtf.org/standards/redfish). Exposes read-only and read-write capabilities as MCP tools that can be used with Claude Code, Claude Desktop, or any MCP-compatible client.

## Features

- **13 read-only tools** – service root, system info, processors, memory, storage, networking, thermals, power, event log, managers, BIOS
- **4 read-write tools** – system reset, indicator LED control, event log clearing, BIOS attribute updates
- **Read-only mode** – optionally restrict to safe, non-destructive operations
- **TLS support** with optional certificate verification bypass for lab/dev BMCs
- **Basic Auth** for BMC authentication
- Released as both a **container image** and **Homebrew formula**

## Installation

### Homebrew

```bash
brew tap conallob/tap
brew install mcp-redfish
```

### Container

```bash
docker pull ghcr.io/conallob/mcp-redfish:latest
```

### From source

```bash
go install github.com/conallob/mcp-redfish@latest
```

## Usage

```
mcp-redfish --host <bmc-host> [options]

Options:
  --host       BMC hostname or IP (env: REDFISH_HOST)
  --username   BMC username (env: REDFISH_USERNAME)
  --password   BMC password (env: REDFISH_PASSWORD)
  --insecure   Skip TLS certificate verification (env: REDFISH_INSECURE)
  --read-only  Restrict to read-only tools only (env: REDFISH_READ_ONLY)
  --version    Print version and exit
```

## MCP Client Configuration

### Claude Code / Claude Desktop

Add to your MCP server config (`~/.claude.json` or `claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "redfish": {
      "command": "mcp-redfish",
      "args": ["--host", "192.168.1.100", "--username", "admin", "--insecure"],
      "env": {
        "REDFISH_PASSWORD": "your-bmc-password"
      }
    }
  }
}
```

For read-only access (safe for production use):

```json
{
  "mcpServers": {
    "redfish-ro": {
      "command": "mcp-redfish",
      "args": ["--host", "192.168.1.100", "--read-only"],
      "env": {
        "REDFISH_USERNAME": "readonly-user",
        "REDFISH_PASSWORD": "password"
      }
    }
  }
}
```

### Container usage

```json
{
  "mcpServers": {
    "redfish": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-e", "REDFISH_HOST=192.168.1.100",
        "-e", "REDFISH_USERNAME=admin",
        "-e", "REDFISH_PASSWORD=password",
        "ghcr.io/conallob/mcp-redfish:latest"
      ]
    }
  }
}
```

## Available Tools

### Read-only

| Tool | Description |
|------|-------------|
| `redfish_get_service_root` | Redfish API version and available collections |
| `redfish_list_systems` | List all managed computer systems |
| `redfish_get_system` | System details: model, serial, power state, BIOS version |
| `redfish_get_processors` | CPU socket, model, core/thread count |
| `redfish_get_memory` | DIMM type, capacity, speed, slot location |
| `redfish_get_storage` | Storage controllers and drives |
| `redfish_get_network_interfaces` | NIC inventory |
| `redfish_get_thermal` | Temperature sensors and fan speeds |
| `redfish_get_power` | Power consumption, voltages, PSU status |
| `redfish_get_event_log` | System event log (SEL) entries |
| `redfish_list_managers` | BMC/management controller inventory |
| `redfish_get_manager` | BMC firmware version, model, date/time |
| `redfish_get_bios` | BIOS/UEFI attribute settings |

### Read-write (disabled with `--read-only`)

| Tool | Description |
|------|-------------|
| `redfish_reset_system` | Power on/off/restart the system |
| `redfish_set_indicator_led` | Control the physical ID LED (Lit/Blinking/Off) |
| `redfish_clear_event_log` | Clear the system event log |
| `redfish_set_bios_attribute` | Stage a BIOS attribute change (applied on next reboot) |

## Development

```bash
# Run tests
go test ./...

# Run tests with race detection
go test -race ./...

# Build
go build -o mcp-redfish .

# Lint
golangci-lint run
```

## Release

Releases are built automatically by [GoReleaser](https://goreleaser.com/) when a `v*` tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

This produces:
- GitHub release with cross-compiled binaries (linux/darwin/windows, amd64/arm64)
- Multi-arch container images published to `ghcr.io/conallob/mcp-redfish`
- Homebrew formula updated in `conallob/homebrew-tap`

## License

BSD 3-Clause — see [LICENSE](LICENSE).
