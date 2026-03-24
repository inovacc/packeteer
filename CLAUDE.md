# CLAUDE.md — Sharkline

## Overview

Sharkline is a Go MCP server wrapping Wireshark CLI tools (tshark, capinfos, editcap, mergecap) for AI-driven packet capture and network analysis.

## Build & Test

```bash
task build          # Build binary
task test           # Run tests
task lint           # Run golangci-lint
go run . serve      # Start MCP server (stdio)
go run . setup      # Check/install Wireshark dependencies
```

## Architecture

- **Hexagonal/Clean:** `cmd/` → `internal/server/` → `internal/tools/` → `internal/executor/`
- **Key port:** `CommandExecutor` interface in `internal/executor/executor.go`
- **Transports:** stdio (default) and Streamable HTTP (`--transport http --port 8080`)
- **Logger:** `log/slog` to stderr (stdout reserved for MCP protocol)

### Packages

| Package | Purpose |
|---------|---------|
| `cmd/` | Cobra CLI commands: serve, setup, version, aicontext, cmdtree |
| `internal/server/` | MCP server factory, tool/resource/prompt registration |
| `internal/tools/` | 10 tool handlers, resource browser, 3 prompt workflows |
| `internal/executor/` | CommandExecutor interface, RealExecutor, MockExecutor |
| `internal/safety/` | Input validation, filter sanitization, CaptureLimiter |
| `internal/output/` | Output truncation, structured JSON packet parser |
| `internal/setup/` | Cross-platform Wireshark detection and auto-install |

## Conventions

- MCP SDK: `github.com/modelcontextprotocol/go-sdk/mcp`
- Typed tool inputs with `json` + `jsonschema` struct tags
- All CLI commands go through `CommandExecutor` — never call `exec.Command` directly in tool handlers
- Safety: validate paths, sanitize filters, clamp timeouts (30s max), limit packet count (1000 max), max 3 concurrent captures
- Tests: table-driven with mock executor + in-memory MCP transport
