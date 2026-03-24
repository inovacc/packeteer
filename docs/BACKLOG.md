# Backlog

## Priority Levels

| Priority | Timeline |
|----------|----------|
| P1 | This sprint |
| P2 | This quarter |
| P3 | Future |

## Items

### P1 — This Sprint

- **Test coverage for tools resources/prompts**
  - Description: RegisterResources (0%) and RegisterPrompts (0%) need direct unit tests
  - Effort: Medium
  - Category: Tech Debt

- **Test coverage for setup package**
  - Description: Install functions at 0%, PrintInstructions at 55%, overall 26.5%
  - Effort: Medium
  - Category: Tech Debt

### P2 — This Quarter

- **Npcap auto-install on Windows**
  - Description: Setup command should detect missing Npcap and offer to install for live capture
  - Effort: Medium
  - Category: Feature

- **Structured output for all tools**
  - Description: Extend ParseTSharkJSON to statistics, decode_packet, extract_fields
  - Effort: Large
  - Category: Feature

### P3 — Future

- **WebSocket transport**
  - Description: Support WebSocket-based MCP connections for browser integrations
  - Effort: Large
  - Category: Feature

## Resolved

| Item | Resolution | Date |
|------|------------|------|
| Sample pcap + quickstart demo | Implemented: testdata/sample.pcap + E2E tests | 2026-03-24 |
| Structured JSON output parsing | Implemented: ParseTSharkJSON + summarize flag | 2026-03-24 |
| Concurrent capture management | Implemented: CaptureLimiter with semaphore | 2026-03-24 |
| Resource/prompt test coverage | Covered via server integration tests (91.3%) | 2026-03-24 |
| Streamable HTTP transport | Implemented: --transport http --port 8080 | 2026-03-24 |
| CI integration tests | Implemented: build tag + tshark in GitHub Actions | 2026-03-24 |
| GoReleaser pipeline | Validated: 6 platform builds pass | 2026-03-24 |
