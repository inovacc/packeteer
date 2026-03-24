# Roadmap

## Current Status
**Overall Progress:** 100% - v1.0.0 released

## Phases

### Phase 1: Foundation [COMPLETE]
- [x] Project scaffold (structure, tooling, CI config)
- [x] Branding and identity (Packeteer)
- [x] CommandExecutor interface and RealExecutor implementation
- [x] MockExecutor for testing
- [x] Safety/validation module (path, filter, timeout, count)
- [x] MCP server skeleton with stdio transport
- [x] Output truncation module

### Phase 2: Core Tools [COMPLETE]
- [x] list_interfaces — tshark -D
- [x] read_pcap — read/analyze pcap files with display filters
- [x] capture_packets — live capture with filters, timeout, count limits
- [x] extract_fields — field extraction from pcap files
- [x] get_statistics — protocol hierarchy, conversations, endpoints

### Phase 3: Extended Tools [COMPLETE]
- [x] get_capture_info — capinfos wrapper
- [x] filter_pcap — editcap wrapper
- [x] merge_pcaps — mergecap wrapper
- [x] list_protocols — protocol dissector listing
- [x] decode_packet — verbose packet decode

### Phase 4: Polish & Release [COMPLETE]
- [x] MCP resources for capture file browsing
- [x] MCP prompts (analyze-traffic, investigate-connection, security-scan)
- [x] Unit tests for executor, output, safety, server, tools
- [x] CI integration tests with real tshark
- [x] Streamable HTTP transport option
- [x] GoReleaser pipeline validation (6 platforms)
- [x] v1.0.0 released — 2026-03-24

### Phase 5: Hardening & v1.1.0 [NOT STARTED]
- [ ] Increase test coverage to 80%+ (resource/prompt tests)
- [ ] Sample pcap + end-to-end quickstart demo
- [ ] Structured JSON output parsing for typed packet data
- [ ] Concurrent capture management with semaphore limiting

## Test Coverage

**Current:** 37.7%  |  **Target:** 80%

| Package | Coverage | Status |
|---------|----------|--------|
| internal/safety | 100.0% | Complete |
| internal/output | 100.0% | Complete |
| internal/server | 85.0% | Good |
| internal/executor | 68.1% | Good |
| internal/tools | 58.0% | Needs improvement (resources/prompts untested) |
| cmd | 0.0% | Scaffold code |
