# Milestones

## v0.1.0 - Foundation
- **Status:** Complete
- **Test Coverage:** 37.7% (safety 100%, output 100%, server 85%, executor 68%, tools 58%)
- **Goals:**
  - [x] Project scaffolding with Cobra CLI
  - [x] Branding (Packeteer)
  - [x] CommandExecutor interface with RealExecutor and MockExecutor
  - [x] Safety module (path validation, filter sanitization, timeout/count clamping)
  - [x] Output truncation module
  - [x] MCP server with stdio transport
  - [x] All 10 MCP tools implemented (5 tier-1, 5 tier-2)
  - [x] Unit tests for safety and tool handlers

## v0.2.0 - Test Coverage & Robustness
- **Status:** Complete
- **Goals:**
  - [x] Executor package tests (68.1%)
  - [x] Output package tests (100%)
  - [x] Server integration tests via in-memory MCP transport (85%)
  - [x] Safety edge case tests — ValidateOutputPath, SanitizeStatType (100%)
  - [x] Filter pcap handler tests (95%)

## v0.3.0 - MCP Enhancements
- **Status:** In Progress
- **Goals:**
  - [x] MCP resources for capture files (packeteer://captures, packeteer://captures/{filename})
  - [x] MCP prompts: analyze-traffic, investigate-connection, security-scan
  - [ ] Streamable HTTP transport option

## v1.0.0 - First Stable Release
- **Target:** TBD
- **Status:** Not Started
- **Goals:**
  - [ ] CI integration tests with real tshark
  - [ ] GoReleaser pipeline validated
  - [ ] Documentation complete
  - [ ] Tag and publish release
