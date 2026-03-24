package server

import (
	"log/slog"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "0.1.0"

// New creates and configures the Packeteer MCP server with all tools, resources, and prompts registered.
func New(exec executor.CommandExecutor, logger *slog.Logger, opts ...Option) *mcp.Server {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "packeteer",
			Version: version,
		},
		&mcp.ServerOptions{
			Instructions: `Packeteer provides Wireshark CLI tools for packet capture and network analysis.

Available tools cover the full Wireshark suite:
- tshark: capture, read, filter, extract fields, decode packets, list interfaces/protocols, statistics
- capinfos: capture file metadata
- editcap: filter/extract packets from pcap files
- mergecap: combine multiple capture files

Safety: captures are limited to 30s/1000 packets. File paths are validated. Filters are sanitized.`,
		},
	)

	// Tier 1 — Essential tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_interfaces",
		Description: "List available network interfaces for packet capture (tshark -D). Returns interface names and descriptions.",
	}, tools.NewListInterfacesHandler(exec))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "capture_packets",
		Description: "Capture live network packets with optional BPF and display filters. Returns JSON-formatted packet data. Captures are limited to 30 seconds and 1000 packets for safety.",
	}, tools.NewCaptureHandler(exec))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "read_pcap",
		Description: "Read and analyze an existing pcap/pcapng capture file. Supports Wireshark display filters. Returns JSON-formatted packet data.",
	}, tools.NewReadPcapHandler(exec))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "extract_fields",
		Description: "Extract specific protocol fields from a pcap file (e.g., ip.src, tcp.port, http.host). Returns tab-separated field values.",
	}, tools.NewExtractFieldsHandler(exec))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_statistics",
		Description: "Generate network statistics from a pcap file. Supports protocol hierarchy (io,phs), TCP/UDP/IP conversations (conv,tcp), endpoints, and I/O statistics.",
	}, tools.NewStatisticsHandler(exec))

	// Tier 2 — Extended tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_capture_info",
		Description: "Get metadata about a capture file using capinfos: packet count, duration, file size, encapsulation type, and more.",
	}, tools.NewCaptureInfoHandler(exec))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "filter_pcap",
		Description: "Filter and extract packets from a pcap file into a new file using editcap. Supports time range filtering and packet count limits.",
	}, tools.NewFilterPcapHandler(exec))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "merge_pcaps",
		Description: "Merge multiple pcap/pcapng files into a single file using mergecap. Files are merged chronologically.",
	}, tools.NewMergePcapsHandler(exec))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_protocols",
		Description: "List all available protocol dissectors in the installed Wireshark version. Optionally filter by name.",
	}, tools.NewListProtocolsHandler(exec))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "decode_packet",
		Description: "Verbose decode of packets from a pcap file showing all protocol layers and field values. Limited to 10 packets max to prevent excessive output.",
	}, tools.NewDecodePacketHandler(exec))

	// Resources
	tools.RegisterResources(server, exec, cfg.captureDir)

	// Prompts
	tools.RegisterPrompts(server)

	logger.Info("packeteer MCP server initialized",
		"version", version,
		"tools", 10,
		"prompts", 3,
		"capture_dir", cfg.captureDir,
	)

	return server
}

type config struct {
	captureDir string
}

// Option configures the MCP server.
type Option func(*config)

// WithCaptureDir sets the directory for browsing capture files via MCP resources.
func WithCaptureDir(dir string) Option {
	return func(c *config) {
		c.captureDir = dir
	}
}
