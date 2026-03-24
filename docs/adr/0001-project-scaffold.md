# ADR-0001: Project Scaffold and Tooling Choices

## Status
Accepted

## Context
Setting up Packeteer — a Go MCP server wrapping Wireshark CLI tools for AI-driven packet capture and network analysis. Requires choosing standard structure, tooling, and conventions.

## Decision
- **Structure:** Hexagonal/Clean Architecture (cmd/, internal/)
- **CLI Framework:** Cobra via omni scaffold
- **MCP SDK:** github.com/modelcontextprotocol/go-sdk/mcp (stdio transport)
- **Task Runner:** Taskfile (over Makefile) for cross-platform support
- **Linting:** golangci-lint v2 with curated ruleset
- **Releases:** GoReleaser for automated cross-platform builds
- **Module Path:** github.com/inovacc/sharkline

## Consequences

### Positive
- Consistent project structure across all projects
- Cross-platform build and task support
- Automated release pipeline from day one
- Strict code quality from the start
- CommandExecutor interface enables testing without tshark installed

### Negative
- Requires installing external tools (golangci-lint, goreleaser, task)
- Requires Wireshark/tshark installed for runtime use
