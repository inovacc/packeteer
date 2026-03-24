# Packeteer

**Packet intelligence, on demand.**

Packeteer is an MCP (Model Context Protocol) server that wraps Wireshark's CLI tools — tshark, capinfos, editcap, and mergecap — giving AI assistants the ability to capture, analyze, and dissect network traffic.

## Prerequisites

- [Wireshark](https://www.wireshark.org/download.html) installed (provides tshark, capinfos, editcap, mergecap)
- Go 1.23+

## Installation

```bash
go install github.com/inovacc/packeteer@latest
```

## Claude Desktop Configuration

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "packeteer": {
      "command": "packeteer",
      "args": ["serve"]
    }
  }
}
```

## MCP Tools

### Tier 1 — Essential

| Tool | Description |
|------|-------------|
| `list_interfaces` | Show available network interfaces |
| `capture_packets` | Live capture with BPF/display filters, timeout, count limits |
| `read_pcap` | Read and analyze pcap files with display filters |
| `extract_fields` | Extract specific protocol fields from captures |
| `get_statistics` | Protocol hierarchy, conversations, endpoints |

### Tier 2 — Extended

| Tool | Description |
|------|-------------|
| `get_capture_info` | Capture file metadata (capinfos) |
| `filter_pcap` | Filter/extract packets from pcap (editcap) |
| `merge_pcaps` | Combine multiple capture files (mergecap) |
| `list_protocols` | Available protocol dissectors |
| `decode_packet` | Verbose decode of specific packets |

## MCP Prompts

| Prompt | Description |
|--------|-------------|
| `analyze-traffic` | Guided workflow for protocol breakdown, top talkers, anomaly detection |
| `investigate-connection` | Deep-dive into a specific connection between two hosts |
| `security-scan` | Scan for DNS exfiltration, cleartext credentials, port scanning, TLS issues |

## MCP Resources

| Resource | Description |
|----------|-------------|
| `packeteer://captures` | List available pcap files in the captures directory |
| `packeteer://captures/{filename}` | Get metadata for a specific capture file |

## Usage

### Stdio Transport (default)

```bash
packeteer serve
```

### HTTP Transport (remote)

```bash
packeteer serve --transport http --port 8080
```

### With Capture Directory

```bash
packeteer serve --capture-dir /path/to/pcaps
```

## Development

```bash
# Build
task build

# Run
task run

# Test
task test

# Lint
task lint
```

## Release

```bash
# Create a snapshot release
task release:snapshot

# Create a production release (requires git tag)
git tag v1.0.0
task release
```

## License

MIT
