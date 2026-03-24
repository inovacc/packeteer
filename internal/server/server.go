package server

import (
	"log/slog"

	"github.com/inovacc/sharkline/internal/executor"
	"github.com/inovacc/sharkline/internal/safety"
	"github.com/inovacc/sharkline/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "1.1.0-dev"

// New creates and configures the Sharkline MCP server with all tools, resources, and prompts registered.
func New(exec executor.CommandExecutor, logger *slog.Logger, opts ...Option) *mcp.Server {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "sharkline",
			Version: version,
		},
		&mcp.ServerOptions{
			Instructions: `Sharkline provides Wireshark CLI tools for packet capture and network analysis.

Available tools cover the full Wireshark suite:
- tshark: capture, read, filter, extract fields, decode packets, list interfaces/protocols, statistics
- capinfos: capture file metadata
- editcap: filter/extract packets from pcap files
- mergecap: combine multiple capture files

Safety: captures are limited to 30s/1000 packets, max 3 concurrent. File paths are validated. Filters are sanitized.
Use summarize=true on read_pcap/capture_packets for structured packet summaries instead of raw JSON.`,
		},
	)

	captureLimiter := safety.NewCaptureLimiter(cfg.maxConcurrentCaptures)

	// Tier 1 — Essential tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_interfaces",
		Description: "List available network interfaces for packet capture (tshark -D). Returns interface names and descriptions.",
	}, tools.NewListInterfacesHandler(exec))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "capture_packets",
		Description: "Capture live network packets with optional BPF and display filters. Returns JSON-formatted packet data. Captures are limited to 30 seconds, 1000 packets, and 3 concurrent. Set summarize=true for structured summaries.",
	}, tools.NewCaptureHandler(exec, captureLimiter))

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

	logger.Info("sharkline MCP server initialized",
		"version", version,
		"tools", 10,
		"prompts", 3,
		"capture_dir", cfg.captureDir,
	)

	return server
}

type config struct {
	captureDir            string
	maxConcurrentCaptures int
}

// Option configures the MCP server.
type Option func(*config)

// WithCaptureDir sets the directory for browsing capture files via MCP resources.
func WithCaptureDir(dir string) Option {
	return func(c *config) {
		c.captureDir = dir
	}
}

// WithMaxConcurrentCaptures sets the maximum number of simultaneous live captures.
// Defaults to 3 if not set or set to 0.
func WithMaxConcurrentCaptures(max int) Option {
	return func(c *config) {
		c.maxConcurrentCaptures = max
	}
}
