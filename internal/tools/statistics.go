package tools

import (
	"context"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/output"
	"github.com/inovacc/packeteer/internal/safety"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type StatisticsInput struct {
	FilePath      string `json:"file_path" jsonschema:"path to pcap/pcapng file (required)"`
	StatType      string `json:"stat_type" jsonschema:"statistics type: 'io,phs' (protocol hierarchy), 'conv,tcp' (TCP conversations), 'conv,udp' (UDP conversations), 'conv,ip' (IP conversations), 'endpoints,tcp', 'endpoints,ip', 'io,stat,1' (1-second intervals)"`
	DisplayFilter string `json:"display_filter,omitempty" jsonschema:"Wireshark display filter to apply"`
}

func NewStatisticsHandler(exec executor.CommandExecutor) func(context.Context, *mcp.CallToolRequest, StatisticsInput) (*mcp.CallToolResult, struct{}, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input StatisticsInput) (*mcp.CallToolResult, struct{}, error) {
		if err := safety.ValidateFilePath(input.FilePath); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}
		if err := safety.SanitizeStatType(input.StatType); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}
		if err := safety.SanitizeDisplayFilter(input.DisplayFilter); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}

		args := []string{
			"-r", input.FilePath,
			"-q",
			"-z", input.StatType,
		}

		if input.DisplayFilter != "" {
			args = append(args, "-Y", input.DisplayFilter)
		}

		stdout, _, err := exec.Execute(ctx, "tshark", args)
		if err != nil {
			return errorResult("Failed to get statistics: " + err.Error()), struct{}{}, nil
		}

		metadata := map[string]string{
			"File":      input.FilePath,
			"Stat Type": input.StatType,
		}

		return textResult(output.FormatResult(string(stdout), metadata)), struct{}{}, nil
	}
}
