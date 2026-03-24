# Feature Requests

## Completed Features

### Live Packet Capture
- **Status:** Completed
- **Description:** Capture live network packets with BPF and display filters, timeout and count limits, concurrent capture limiting
- **Tool:** `capture_packets`

### Pcap File Analysis
- **Status:** Completed
- **Description:** Read and analyze pcap/pcapng files with Wireshark display filters, JSON output, optional structured summaries
- **Tool:** `read_pcap`

### Protocol Field Extraction
- **Status:** Completed
- **Description:** Extract specific protocol fields (ip.src, tcp.port, etc.) from capture files
- **Tool:** `extract_fields`

### Network Statistics
- **Status:** Completed
- **Description:** Protocol hierarchy, TCP/UDP/IP conversations, endpoints, I/O statistics
- **Tool:** `get_statistics`

### Capture File Metadata
- **Status:** Completed
- **Description:** Get file info via capinfos — packet count, duration, encapsulation type
- **Tool:** `get_capture_info`

### Pcap Filtering
- **Status:** Completed
- **Description:** Filter/extract packets from pcap by time range and count via editcap
- **Tool:** `filter_pcap`

### Pcap Merging
- **Status:** Completed
- **Description:** Merge multiple capture files chronologically via mergecap
- **Tool:** `merge_pcaps`

### Protocol Listing
- **Status:** Completed
- **Description:** List all available Wireshark protocol dissectors with optional name filter
- **Tool:** `list_protocols`

### Verbose Packet Decode
- **Status:** Completed
- **Description:** Full protocol layer decode of specific packets from a capture file
- **Tool:** `decode_packet`

### Interface Listing
- **Status:** Completed
- **Description:** List available network interfaces for packet capture
- **Tool:** `list_interfaces`

### MCP Resources for Capture Files
- **Status:** Completed (v1.0.0)
- **Description:** Browse pcap files as MCP resources via `sharkline://captures` and `sharkline://captures/{filename}`

### Analysis Workflow Prompts
- **Status:** Completed (v1.0.0)
- **Description:** 3 guided MCP prompts: analyze-traffic, investigate-connection, security-scan

### HTTP Transport
- **Status:** Completed (v1.0.0)
- **Description:** Streamable HTTP transport for remote MCP connections (`--transport http --port 8080`)

### Structured JSON Packet Parsing
- **Status:** Completed (v1.1.0)
- **Description:** Parse tshark JSON into typed packet summaries (source, dest, protocol, info) via `summarize=true`

### Concurrent Capture Limiting
- **Status:** Completed (v1.1.0)
- **Description:** Semaphore-based CaptureLimiter (default max 3) prevents resource exhaustion

### Auto-Install Setup Command
- **Status:** Completed (v1.1.0)
- **Description:** `sharkline setup --install` auto-installs Wireshark via winget/choco/brew/apt/dnf/pacman

## Proposed Features

### Npcap Auto-Install
- **Priority:** P2
- **Status:** Proposed
- **Description:** Detect missing Npcap on Windows and offer to install for live capture support
- **Motivation:** Live capture fails without Npcap; setup command should handle this

### Structured Output for All Tools
- **Priority:** P2
- **Status:** Proposed
- **Description:** Extend structured JSON parsing to statistics, decode_packet, and extract_fields
- **Motivation:** Consistent typed output across all tools improves AI consumption
