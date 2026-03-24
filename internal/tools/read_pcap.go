package tools

import (
	"context"
	"strconv"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/output"
	"github.com/inovacc/packeteer/internal/safety"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ReadPcapInput struct {
	FilePath      string `json:"file_path" jsonschema:"path to pcap/pcapng file to read (required)"`
	DisplayFilter string `json:"display_filter,omitempty" jsonschema:"Wireshark display filter (e.g. 'tcp.port == 80')"`
	MaxPackets    int    `json:"max_packets,omitempty" jsonschema:"maximum packets to return (max 1000, default 100)"`
}

func NewReadPcapHandler(exec executor.CommandExecutor) func(context.Context, *mcp.CallToolRequest, ReadPcapInput) (*mcp.CallToolResult, struct{}, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input ReadPcapInput) (*mcp.CallToolResult, struct{}, error) {
		if err := safety.ValidateFilePath(input.FilePath); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}
		if err := safety.SanitizeDisplayFilter(input.DisplayFilter); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}

		count := input.MaxPackets
		if count <= 0 {
			count = 100
		}
		count = safety.ClampPacketCount(count)

		args := []string{
			"-r", input.FilePath,
			"-T", "json",
			"-c", strconv.Itoa(count),
		}

		if input.DisplayFilter != "" {
			args = append(args, "-Y", input.DisplayFilter)
		}

		stdout, _, err := exec.Execute(ctx, "tshark", args)
		if err != nil {
			return errorResult("Failed to read pcap: " + err.Error()), struct{}{}, nil
		}

		truncated, wasTruncated := output.Truncate(stdout, output.DefaultMaxBytes)
		metadata := map[string]string{
			"File":       input.FilePath,
			"Max Packets": strconv.Itoa(count),
			"Truncated":  strconv.FormatBool(wasTruncated),
		}

		if input.DisplayFilter != "" {
			metadata["Filter"] = input.DisplayFilter
		}

		return textResult(output.FormatResult(string(truncated), metadata)), struct{}{}, nil
	}
}
