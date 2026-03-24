# Milestones

## v0.1.0 - Foundation
- **Target:** Complete
- **Status:** Complete
- **Test Coverage:** 27.0% (safety 74.5%, tools 68.4%, executor 0%, output 0%)
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
- **Target:** TBD
- **Status:** Not Started
- **Coverage Target:** 80%+
- **Goals:**
  - [ ] Executor package tests
  - [ ] Output package tests
  - [ ] Server integration tests (in-memory MCP transport)
  - [ ] Safety edge case tests (ValidateOutputPath, SanitizeStatType)
  - [ ] Filter pcap handler tests

## v0.3.0 - MCP Enhancements
- **Target:** TBD
- **Status:** Not Started
- **Goals:**
  - [ ] MCP resources for capture files
  - [ ] MCP prompts for common analysis workflows
  - [ ] Streamable HTTP transport option

## v1.0.0 - First Stable Release
- **Target:** TBD
- **Status:** Not Started
- **Goals:**
  - [ ] Full test coverage (80%+)
  - [ ] Integration tests with real tshark in CI
  - [ ] Documentation complete
  - [ ] GoReleaser pipeline validated
