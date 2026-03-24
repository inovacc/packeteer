package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/inovacc/sharkline/internal/executor"
	"github.com/inovacc/sharkline/internal/output"
	"github.com/inovacc/sharkline/internal/safety"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type DecodePacketInput struct {
	FilePath      string `json:"file_path" jsonschema:"path to pcap/pcapng file (required)"`
	PacketNumber  int    `json:"packet_number,omitempty" jsonschema:"specific packet number to decode (1-based, default: first 5 packets)"`
	DisplayFilter string `json:"display_filter,omitempty" jsonschema:"Wireshark display filter"`
	MaxPackets    int    `json:"max_packets,omitempty" jsonschema:"max packets to decode verbosely (max 10, default 5)"`
}

func NewDecodePacketHandler(exec executor.CommandExecutor) func(context.Context, *mcp.CallToolRequest, DecodePacketInput) (*mcp.CallToolResult, struct{}, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input DecodePacketInput) (*mcp.CallToolResult, struct{}, error) {
		if err := safety.ValidateFilePath(input.FilePath); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}
		if err := safety.SanitizeDisplayFilter(input.DisplayFilter); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}

		maxPkts := input.MaxPackets
		if maxPkts <= 0 {
			maxPkts = 5
		}
		if maxPkts > 10 {
			maxPkts = 10
		}

		args := []string{
			"-r", input.FilePath,
			"-V",
		}

		if input.PacketNumber > 0 {
			// Use display filter to select specific packet by frame number
			frameFilter := fmt.Sprintf("frame.number == %d", input.PacketNumber)
			if input.DisplayFilter != "" {
				frameFilter = fmt.Sprintf("(%s) && frame.number == %d", input.DisplayFilter, input.PacketNumber)
			}
			args = append(args, "-Y", frameFilter, "-c", "1")
		} else {
			args = append(args, "-c", strconv.Itoa(maxPkts))
			if input.DisplayFilter != "" {
				args = append(args, "-Y", input.DisplayFilter)
			}
		}

		stdout, _, err := exec.Execute(ctx, "tshark", args)
		if err != nil {
			return errorResult("Failed to decode packet: " + err.Error()), struct{}{}, nil
		}

		truncated, wasTruncated := output.Truncate(stdout, output.DefaultMaxBytes)
		metadata := map[string]string{
			"File":      input.FilePath,
			"Truncated": strconv.FormatBool(wasTruncated),
		}

		if input.PacketNumber > 0 {
			metadata["Packet"] = strconv.Itoa(input.PacketNumber)
		} else {
			metadata["Max Packets"] = strconv.Itoa(maxPkts)
		}

		return textResult(output.FormatResult(string(truncated), metadata)), struct{}{}, nil
	}
}
