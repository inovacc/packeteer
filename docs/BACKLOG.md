# Backlog

## Priority Levels

| Priority | Timeline |
|----------|----------|
| P1 | This sprint |
| P2 | This quarter |
| P3 | Future |

## Items

### P1 — This Sprint

- **Resource/prompt test coverage**
  - Description: RegisterResources and RegisterPrompts have 0% coverage, dragging tools to 58%
  - Effort: Medium
  - Category: Tech Debt

- **Sample pcap + quickstart demo**
  - Description: Bundle a small pcap, end-to-end MCP test, README quickstart section
  - Effort: Medium
  - Category: Onboarding

### P2 — This Quarter

- **Structured JSON output parsing**
  - Description: Parse tshark JSON into typed Go structs for read_pcap/capture_packets
  - Effort: Large
  - Category: Feature

- **Concurrent capture management**
  - Description: Semaphore limiting active captures, cancel mechanism
  - Effort: Large
  - Category: Safety
