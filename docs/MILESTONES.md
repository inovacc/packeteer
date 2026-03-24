# Milestones

## v1.0.0 - First Stable Release
- **Released:** 2026-03-24
- **Status:** Complete
- **Test Coverage:** 37.7%
- **Goals:**
  - [x] 10 MCP tools (tshark, capinfos, editcap, mergecap)
  - [x] 3 MCP prompts (analyze-traffic, investigate-connection, security-scan)
  - [x] MCP resources for capture file browsing
  - [x] Stdio and Streamable HTTP transports
  - [x] Safety guardrails (path validation, filter sanitization, timeout/count clamping)
  - [x] GoReleaser pipeline (linux/windows/darwin x amd64/arm64)
  - [x] CI with tshark integration tests

## v1.1.0 - Hardening
- **Released:** 2026-03-24
- **Status:** Complete
- **Test Coverage:** 40.4%
- **Goals:**
  - [x] Structured JSON output parsing (summarize=true on read_pcap/capture_packets)
  - [x] Concurrent capture management (CaptureLimiter, max 3)
  - [x] Sample pcap with end-to-end integration tests
  - [x] Setup command with cross-platform auto-install (winget/choco/brew/apt/dnf/pacman)
  - [x] Project rename: Packeteer → Sharkline

## v1.2.0 - Coverage & Polish
- **Target:** TBD
- **Status:** Not Started
- **Goals:**
  - [ ] Test coverage 80%+ (tools resources/prompts, setup package)
  - [ ] Npcap auto-install on Windows for live capture
  - [ ] Structured output for all tools (not just read_pcap/capture_packets)
