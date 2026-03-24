package tools

import (
	"context"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/output"
	"github.com/inovacc/packeteer/internal/safety"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CaptureInfoInput struct {
	FilePath string `json:"file_path" jsonschema:"path to pcap/pcapng file (required)"`
}

func NewCaptureInfoHandler(exec executor.CommandExecutor) func(context.Context, *mcp.CallToolRequest, CaptureInfoInput) (*mcp.CallToolResult, struct{}, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input CaptureInfoInput) (*mcp.CallToolResult, struct{}, error) {
		if err := safety.ValidateFilePath(input.FilePath); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}

		stdout, _, err := exec.Execute(ctx, "capinfos", []string{input.FilePath})
		if err != nil {
			return errorResult("Failed to get capture info: " + err.Error()), struct{}{}, nil
		}

		metadata := map[string]string{
			"File":    input.FilePath,
			"Command": "capinfos",
		}

		return textResult(output.FormatResult(string(stdout), metadata)), struct{}{}, nil
	}
}
