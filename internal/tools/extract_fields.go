package tools

import (
	"context"
	"strconv"
	"strings"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/output"
	"github.com/inovacc/packeteer/internal/safety"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ExtractFieldsInput struct {
	FilePath      string   `json:"file_path" jsonschema:"path to pcap/pcapng file (required)"`
	Fields        []string `json:"fields" jsonschema:"protocol fields to extract (e.g. ['ip.src', 'ip.dst', 'tcp.port'])"`
	DisplayFilter string   `json:"display_filter,omitempty" jsonschema:"Wireshark display filter to apply"`
	MaxPackets    int      `json:"max_packets,omitempty" jsonschema:"maximum packets to process (max 1000, default 100)"`
	Separator     string   `json:"separator,omitempty" jsonschema:"field separator character (default tab)"`
	ShowHeader    bool     `json:"show_header,omitempty" jsonschema:"include field names as header row"`
}

func NewExtractFieldsHandler(exec executor.CommandExecutor) func(context.Context, *mcp.CallToolRequest, ExtractFieldsInput) (*mcp.CallToolResult, struct{}, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input ExtractFieldsInput) (*mcp.CallToolResult, struct{}, error) {
		if err := safety.ValidateFilePath(input.FilePath); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}
		if err := safety.SanitizeDisplayFilter(input.DisplayFilter); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}
		if len(input.Fields) == 0 {
			return errorResult("at least one field is required"), struct{}{}, nil
		}
		for _, f := range input.Fields {
			if err := safety.SanitizeFieldName(f); err != nil {
				return errorResult(err.Error()), struct{}{}, nil
			}
		}

		count := input.MaxPackets
		if count <= 0 {
			count = 100
		}
		count = safety.ClampPacketCount(count)

		separator := input.Separator
		if separator == "" {
			separator = "\t"
		}

		args := []string{
			"-r", input.FilePath,
			"-T", "fields",
			"-c", strconv.Itoa(count),
			"-E", "separator=" + separator,
		}

		if input.ShowHeader {
			args = append(args, "-E", "header=y")
		}

		for _, field := range input.Fields {
			args = append(args, "-e", field)
		}

		if input.DisplayFilter != "" {
			args = append(args, "-Y", input.DisplayFilter)
		}

		stdout, _, err := exec.Execute(ctx, "tshark", args)
		if err != nil {
			return errorResult("Failed to extract fields: " + err.Error()), struct{}{}, nil
		}

		truncated, wasTruncated := output.Truncate(stdout, output.DefaultMaxBytes)
		metadata := map[string]string{
			"File":      input.FilePath,
			"Fields":    strings.Join(input.Fields, ", "),
			"Truncated": strconv.FormatBool(wasTruncated),
		}

		return textResult(output.FormatResult(string(truncated), metadata)), struct{}{}, nil
	}
}
