package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/safety"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MergePcapsInput struct {
	InputFiles []string `json:"input_files" jsonschema:"list of pcap/pcapng file paths to merge (min 2 required)"`
	OutputFile string   `json:"output_file" jsonschema:"destination merged pcap file path (required)"`
}

func NewMergePcapsHandler(exec executor.CommandExecutor) func(context.Context, *mcp.CallToolRequest, MergePcapsInput) (*mcp.CallToolResult, struct{}, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input MergePcapsInput) (*mcp.CallToolResult, struct{}, error) {
		if len(input.InputFiles) < 2 {
			return errorResult("at least 2 input files are required"), struct{}{}, nil
		}
		if input.OutputFile == "" {
			return errorResult("output_file is required"), struct{}{}, nil
		}

		for _, f := range input.InputFiles {
			if err := safety.ValidateFilePath(f); err != nil {
				return errorResult(fmt.Sprintf("input file %q: %s", f, err.Error())), struct{}{}, nil
			}
		}
		if err := safety.ValidateOutputPath(input.OutputFile); err != nil {
			return errorResult("output: " + err.Error()), struct{}{}, nil
		}

		args := []string{"-w", input.OutputFile}
		args = append(args, input.InputFiles...)

		_, _, err := exec.Execute(ctx, "mergecap", args)
		if err != nil {
			return errorResult("Failed to merge pcaps: " + err.Error()), struct{}{}, nil
		}

		return textResult(fmt.Sprintf("Merged %d files into: %s\nInput files: %s",
			len(input.InputFiles), input.OutputFile, strings.Join(input.InputFiles, ", "))), struct{}{}, nil
	}
}
