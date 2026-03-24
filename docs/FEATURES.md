# Feature Requests

## Completed Features

### Live Packet Capture
- **Status:** Completed
- **Description:** Capture live network packets with BPF and display filters, timeout and count limits
- **Tool:** `capture_packets`

### Pcap File Analysis
- **Status:** Completed
- **Description:** Read and analyze pcap/pcapng files with Wireshark display filters, JSON output
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

## Proposed Features

### MCP Resources for Capture Files
- **Priority:** P2
- **Status:** Proposed
- **Description:** Expose pcap files as browsable MCP resources with URI templates
- **Motivation:** Allow AI to discover and browse available captures without needing file paths

### Analysis Workflow Prompts
- **Priority:** P2
- **Status:** Proposed
- **Description:** Guided MCP prompts for common workflows (traffic analysis, connection investigation)
- **Motivation:** Reduce the learning curve for AI assistants interacting with packet data

### HTTP Transport
- **Priority:** P3
- **Status:** Proposed
- **Description:** Support Streamable HTTP transport for remote MCP connections
- **Motivation:** Enable use from remote AI systems without local process execution
