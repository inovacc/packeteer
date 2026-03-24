package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterPrompts adds MCP prompts for common network analysis workflows.
func RegisterPrompts(server *mcp.Server) {
	server.AddPrompt(
		&mcp.Prompt{
			Name:        "analyze-traffic",
			Description: "Analyze network traffic from a pcap file — protocol breakdown, top talkers, and anomalies",
			Arguments: []*mcp.PromptArgument{
				{Name: "file_path", Description: "Path to the pcap/pcapng file", Required: true},
				{Name: "focus", Description: "Area to focus on: 'overview', 'dns', 'http', 'tls', 'tcp-issues'", Required: false},
			},
		},
		func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			filePath := req.Params.Arguments["file_path"]
			focus := req.Params.Arguments["focus"]
			if focus == "" {
				focus = "overview"
			}

			instructions := fmt.Sprintf(`Analyze the network traffic in %q with a focus on %q.

Follow these steps:
1. First, use get_capture_info to understand the file's scope (duration, packet count, protocols)
2. Use get_statistics with stat_type "io,phs" to see the protocol hierarchy
3. Use get_statistics with stat_type "conv,ip" to identify top talkers
4. Based on the focus area:
   - overview: summarize protocols, top conversations, any anomalies
   - dns: use extract_fields with fields ["dns.qry.name", "dns.a", "dns.resp.type"] filtered by "dns"
   - http: use extract_fields with fields ["http.host", "http.request.uri", "http.response.code"] filtered by "http"
   - tls: use extract_fields with fields ["tls.handshake.extensions_server_name", "tls.handshake.version"] filtered by "tls"
   - tcp-issues: use get_statistics with "conv,tcp" and look for retransmissions, resets, zero windows
5. Summarize findings with key observations and any security concerns`, filePath, focus)

			return &mcp.GetPromptResult{
				Description: "Network traffic analysis workflow",
				Messages: []*mcp.PromptMessage{
					{
						Role:    "user",
						Content: &mcp.TextContent{Text: instructions},
					},
				},
			}, nil
		},
	)

	server.AddPrompt(
		&mcp.Prompt{
			Name:        "investigate-connection",
			Description: "Deep-dive into a specific network connection between two hosts",
			Arguments: []*mcp.PromptArgument{
				{Name: "file_path", Description: "Path to the pcap/pcapng file", Required: true},
				{Name: "source_ip", Description: "Source IP address", Required: true},
				{Name: "dest_ip", Description: "Destination IP address", Required: true},
				{Name: "port", Description: "Port number to filter on (optional)", Required: false},
			},
		},
		func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			filePath := req.Params.Arguments["file_path"]
			srcIP := req.Params.Arguments["source_ip"]
			dstIP := req.Params.Arguments["dest_ip"]
			port := req.Params.Arguments["port"]

			filter := fmt.Sprintf("ip.addr == %s && ip.addr == %s", srcIP, dstIP)
			if port != "" {
				filter += fmt.Sprintf(" && tcp.port == %s", port)
			}

			instructions := fmt.Sprintf(`Investigate the network connection between %s and %s in %q.

Follow these steps:
1. Use read_pcap with display_filter %q to see the conversation packets
2. Use extract_fields with fields ["frame.time", "ip.src", "ip.dst", "tcp.srcport", "tcp.dstport", "tcp.flags", "frame.len"] and the same filter
3. Use get_statistics with stat_type "conv,tcp" to see the conversation metrics
4. If HTTP traffic is present, use extract_fields with fields ["http.request.method", "http.host", "http.request.uri", "http.response.code"]
5. If TLS traffic is present, use extract_fields with fields ["tls.handshake.extensions_server_name", "tls.handshake.version"]
6. Use decode_packet to verbosely decode the first few packets of the connection for handshake analysis

Provide a timeline of the connection, identify the application protocol, note any errors or anomalies.`, srcIP, dstIP, filePath, filter)

			return &mcp.GetPromptResult{
				Description: "Connection investigation workflow",
				Messages: []*mcp.PromptMessage{
					{
						Role:    "user",
						Content: &mcp.TextContent{Text: instructions},
					},
				},
			}, nil
		},
	)

	server.AddPrompt(
		&mcp.Prompt{
			Name:        "security-scan",
			Description: "Scan a pcap file for potential security issues and suspicious activity",
			Arguments: []*mcp.PromptArgument{
				{Name: "file_path", Description: "Path to the pcap/pcapng file", Required: true},
			},
		},
		func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			filePath := req.Params.Arguments["file_path"]

			instructions := fmt.Sprintf(`Perform a security-focused analysis of the traffic in %q.

Follow these steps:
1. Use get_capture_info to understand the capture scope
2. Use get_statistics with "io,phs" for protocol overview
3. Check for DNS exfiltration: extract_fields with fields ["dns.qry.name", "dns.qry.type"] filtered by "dns"
4. Check for unusual ports: extract_fields with fields ["ip.src", "ip.dst", "tcp.dstport"] filtered by "tcp"
5. Check for cleartext credentials: read_pcap filtered by "http.authbasic" or "ftp.request.command"
6. Check for TLS issues: extract_fields with fields ["tls.handshake.version", "tls.handshake.ciphersuite"] filtered by "tls.handshake"
7. Check for port scanning: get_statistics with "conv,tcp" and look for many connections from one source to many ports
8. Check for ARP spoofing: read_pcap filtered by "arp.duplicate-address-detected"

Report findings organized by severity: Critical, Warning, Informational.`, filePath)

			return &mcp.GetPromptResult{
				Description: "Security scan workflow",
				Messages: []*mcp.PromptMessage{
					{
						Role:    "user",
						Content: &mcp.TextContent{Text: instructions},
					},
				},
			}, nil
		},
	)
}
