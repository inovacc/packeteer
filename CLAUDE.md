# CLAUDE.md — Packeteer

## Overview

Packeteer is a Go MCP server wrapping Wireshark CLI tools (tshark, capinfos, editcap, mergecap) for AI-driven packet capture and network analysis.

## Build & Test

```bash
task build          # Build binary
task test           # Run tests
task lint           # Run golangci-lint
go run . serve      # Start MCP server (stdio)
```

## Architecture

- **Hexagonal/Clean:** `cmd/` → `internal/server/` → `internal/tools/` → `internal/executor/`
- **Key port:** `CommandExecutor` interface in `internal/executor/executor.go`
- **Transport:** stdio (MCP JSON-RPC over stdin/stdout)
- **Logger:** `log/slog` to stderr (stdout reserved for MCP protocol)

## Conventions

- MCP SDK: `github.com/modelcontextprotocol/go-sdk/mcp`
- Typed tool inputs with `json` + `jsonschema` struct tags
- All CLI commands go through `CommandExecutor` — never call `exec.Command` directly in tool handlers
- Safety: validate paths, sanitize filters, clamp timeouts (30s max), limit packet count (1000 max)
- Tests: table-driven with mock executor, 80%+ coverage target
