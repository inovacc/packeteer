package tools

import (
	"context"
	"strings"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/output"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListProtocolsInput struct {
	Filter string `json:"filter,omitempty" jsonschema:"filter protocols by name (case-insensitive substring match)"`
}

func NewListProtocolsHandler(exec executor.CommandExecutor) func(context.Context, *mcp.CallToolRequest, ListProtocolsInput) (*mcp.CallToolResult, struct{}, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input ListProtocolsInput) (*mcp.CallToolResult, struct{}, error) {
		stdout, _, err := exec.Execute(ctx, "tshark", []string{"-G", "protocols"})
		if err != nil {
			return errorResult("Failed to list protocols: " + err.Error()), struct{}{}, nil
		}

		result := string(stdout)

		if input.Filter != "" {
			filterLower := strings.ToLower(input.Filter)
			var filtered []string
			for _, line := range strings.Split(result, "\n") {
				if strings.Contains(strings.ToLower(line), filterLower) {
					filtered = append(filtered, line)
				}
			}
			result = strings.Join(filtered, "\n")
			if result == "" {
				return textResult("No protocols matching: " + input.Filter), struct{}{}, nil
			}
		}

		truncated, _ := output.Truncate([]byte(result), output.DefaultMaxBytes)

		metadata := map[string]string{
			"Command": "tshark -G protocols",
		}
		if input.Filter != "" {
			metadata["Filter"] = input.Filter
		}

		return textResult(output.FormatResult(string(truncated), metadata)), struct{}{}, nil
	}
}
