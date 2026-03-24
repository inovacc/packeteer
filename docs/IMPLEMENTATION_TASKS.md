# Implementation Tasks

## Domain 1: Test Coverage Improvement

### 1.1 Add executor package tests
- **What:** Test MockExecutor response matching and RealExecutor binary path resolution
- **Files:** `internal/executor/executor_test.go`
- **Environment:** Go code
- **Dependencies:** None
- **Effort:** Small

### 1.2 Add output package tests
- **What:** Test Truncate() with various sizes and FormatResult() with metadata
- **Files:** `internal/output/truncate_test.go`
- **Environment:** Go code
- **Dependencies:** None
- **Effort:** Small

### 1.3 Add filter_pcap tool tests
- **What:** Test FilterPcapHandler with valid inputs, missing output, time range filtering
- **Files:** `internal/tools/tools_test.go`
- **Environment:** Go code
- **Dependencies:** None
- **Effort:** Small

### 1.4 Add server integration test
- **What:** Test server.New() registers all 10 tools via in-memory MCP transport
- **Files:** `internal/server/server_test.go`
- **Environment:** Go code
- **Dependencies:** None
- **Effort:** Medium

### 1.5 Add safety edge case tests
- **What:** Test ValidateOutputPath and SanitizeStatType (currently 0% coverage)
- **Files:** `internal/safety/validate_test.go`
- **Environment:** Go code
- **Dependencies:** None
- **Effort:** Small

## Domain 2: MCP Server Enhancements

### 2.1 Add MCP resources for capture files
- **What:** Expose pcap files as MCP resources with URI template `sharkline://captures/{filename}`
- **Files:** `internal/server/server.go`, `internal/tools/resources.go`
- **Environment:** Go code
- **Dependencies:** None
- **Effort:** Medium

### 2.2 Add MCP prompts for common analysis workflows
- **What:** Add prompts like "analyze-traffic" and "investigate-connection" with guided parameters
- **Files:** `internal/server/server.go`, `internal/tools/prompts.go`
- **Environment:** Go code
- **Dependencies:** None
- **Effort:** Medium

### 2.3 Add Streamable HTTP transport option
- **What:** Support `--transport http --port 8080` flag for remote MCP connections
- **Files:** `cmd/serve.go`, `internal/server/server.go`
- **Environment:** Go code
- **Dependencies:** None
- **Effort:** Medium

## Domain 3: Safety & Robustness

### 3.1 Add concurrent capture management
- **What:** Track active captures, enforce max concurrent limit, provide cancel mechanism
- **Files:** `internal/tools/capture.go`, `internal/safety/limits.go`
- **Environment:** Go code
- **Dependencies:** None
- **Effort:** Large

### 3.2 Add structured JSON output parsing
- **What:** Parse tshark JSON output into typed Go structs instead of raw text passthrough
- **Files:** `internal/tools/types.go`, update all handlers
- **Environment:** Go code
- **Dependencies:** None
- **Effort:** Large

## Domain 4: CI/CD

### 4.1 Add GitHub Actions integration test
- **What:** Install tshark in CI, run integration tests against real tshark binary
- **Files:** `.github/workflows/test.yml`
- **Environment:** GitHub Actions
- **Dependencies:** 1.4
- **Effort:** Medium

## Suggested Implementation Order

1. **1.1, 1.2, 1.3, 1.5** (parallel) — Quick coverage wins
2. **1.4** — Server integration test
3. **2.1** — MCP resources
4. **2.2** — MCP prompts
5. **3.1** — Concurrent capture management
6. **2.3** — HTTP transport
7. **4.1** — CI integration tests
8. **3.2** — Structured output parsing
