# Roadmap

## Current Status
**Overall Progress:** 90% - All tools, resources, prompts implemented. Polish remaining.

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
- [x] MCP resources for capture file browsing
- [x] MCP prompts (analyze-traffic, investigate-connection, security-scan)
- [x] Unit tests for executor, output, safety, server, tools
- [ ] CI integration tests with real tshark
- [ ] Streamable HTTP transport option
- [ ] GoReleaser pipeline validation
- [ ] v1.0.0 release

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
