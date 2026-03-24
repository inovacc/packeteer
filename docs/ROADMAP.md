# Roadmap

## Current Status
**Overall Progress:** 75% - Core MCP server with all 10 tools implemented

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

### Phase 4: Polish & Release [IN PROGRESS]
- [ ] Increase test coverage to 80%+
- [ ] Integration tests with real tshark
- [ ] MCP resources and prompts
- [ ] HTTP transport option
- [ ] v1.0.0 release

## Test Coverage

**Current:** 27.0%  |  **Target:** 80%

| Package | Coverage | Status |
|---------|----------|--------|
| internal/safety | 74.5% | Good |
| internal/tools | 68.4% | Needs improvement |
| internal/executor | 0.0% | No tests (mock used indirectly) |
| internal/output | 0.0% | No tests |
| internal/server | 0.0% | No tests |
| cmd | 0.0% | Scaffold code |
