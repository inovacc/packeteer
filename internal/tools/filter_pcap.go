package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/safety"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type FilterPcapInput struct {
	InputFile  string `json:"input_file" jsonschema:"source pcap/pcapng file path (required)"`
	OutputFile string `json:"output_file" jsonschema:"destination pcap file path (required)"`
	StartTime  string `json:"start_time,omitempty" jsonschema:"start time filter (e.g. '2024-01-01 00:00:00')"`
	EndTime    string `json:"end_time,omitempty" jsonschema:"end time filter"`
	MaxPackets int    `json:"max_packets,omitempty" jsonschema:"maximum packets to keep (max 1000)"`
}

func NewFilterPcapHandler(exec executor.CommandExecutor) func(context.Context, *mcp.CallToolRequest, FilterPcapInput) (*mcp.CallToolResult, struct{}, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input FilterPcapInput) (*mcp.CallToolResult, struct{}, error) {
		if err := safety.ValidateFilePath(input.InputFile); err != nil {
			return errorResult("input: " + err.Error()), struct{}{}, nil
		}
		if err := safety.ValidateOutputPath(input.OutputFile); err != nil {
			return errorResult("output: " + err.Error()), struct{}{}, nil
		}
		if input.OutputFile == "" {
			return errorResult("output_file is required"), struct{}{}, nil
		}

		args := []string{}

		if input.StartTime != "" {
			args = append(args, "-A", input.StartTime)
		}
		if input.EndTime != "" {
			args = append(args, "-B", input.EndTime)
		}
		if input.MaxPackets > 0 {
			count := safety.ClampPacketCount(input.MaxPackets)
			args = append(args, "-c", strconv.Itoa(count))
		}

		args = append(args, input.InputFile, input.OutputFile)

		_, _, err := exec.Execute(ctx, "editcap", args)
		if err != nil {
			return errorResult("Failed to filter pcap: " + err.Error()), struct{}{}, nil
		}

		return textResult(fmt.Sprintf("Filtered pcap written to: %s", input.OutputFile)), struct{}{}, nil
	}
}
