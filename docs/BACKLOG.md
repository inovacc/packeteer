# Backlog

## Priority Levels

| Priority | Timeline |
|----------|----------|
| P1 | This sprint |
| P2 | This quarter |
| P3 | Future |

## Items

### P1 — This Sprint

- **Test coverage for executor package**
  - Description: Add tests for MockExecutor response matching and RealExecutor.BinaryPath()
  - Effort: Small
  - Category: Tech Debt

- **Test coverage for output package**
  - Description: Add tests for Truncate() boundary conditions and FormatResult()
  - Effort: Small
  - Category: Tech Debt

- **Missing filter_pcap tests**
  - Description: FilterPcapHandler has 0% coverage — add valid/invalid input tests
  - Effort: Small
  - Category: Tech Debt

- **Safety edge cases**
  - Description: ValidateOutputPath and SanitizeStatType have 0% coverage
  - Effort: Small
  - Category: Tech Debt

### P2 — This Quarter

- **Server integration test**
  - Description: Test server.New() with in-memory MCP transport, verify all 10 tools registered
  - Effort: Medium
  - Category: Tech Debt

- **MCP resources for capture files**
  - Description: Expose pcap files as browsable MCP resources
  - Effort: Medium
  - Category: Feature

- **MCP prompts for analysis workflows**
  - Description: Add guided prompts like "analyze-traffic" and "investigate-connection"
  - Effort: Medium
  - Category: Feature

- **CI integration tests with tshark**
  - Description: Install tshark in GitHub Actions, run integration tests
  - Effort: Medium
  - Category: Infrastructure

### P3 — Future

- **Streamable HTTP transport**
  - Description: Support `--transport http --port 8080` for remote connections
  - Effort: Medium
  - Category: Feature

- **Structured JSON output parsing**
  - Description: Parse tshark JSON into typed Go structs instead of raw text
  - Effort: Large
  - Category: Feature

- **Concurrent capture management**
  - Description: Track active captures, enforce max concurrent limit
  - Effort: Large
  - Category: Feature
