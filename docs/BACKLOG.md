# Backlog

## Priority Levels

| Priority | Timeline |
|----------|----------|
| P1 | This sprint |
| P2 | This quarter |
| P3 | Future |

## Items

### P1 — This Sprint

- **CI integration tests with tshark**
  - Description: Install tshark in GitHub Actions, write integration tests with build tag
  - Effort: Medium
  - Category: Infrastructure

- **GoReleaser pipeline validation**
  - Description: Run snapshot release, verify cross-platform binaries
  - Effort: Small
  - Category: Infrastructure

### P2 — This Quarter

- **Streamable HTTP transport**
  - Description: Support `--transport http --port 8080` for remote connections
  - Effort: Medium
  - Category: Feature

- **Resource/prompt test coverage**
  - Description: RegisterResources and RegisterPrompts have 0% coverage
  - Effort: Medium
  - Category: Tech Debt

### P3 — Future

- **Structured JSON output parsing**
  - Description: Parse tshark JSON into typed Go structs instead of raw text
  - Effort: Large
  - Category: Feature

- **Concurrent capture management**
  - Description: Track active captures, enforce max concurrent limit
  - Effort: Large
  - Category: Feature
