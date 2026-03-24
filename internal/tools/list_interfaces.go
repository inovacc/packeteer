package tools

import (
	"context"
	"strings"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/output"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListInterfacesInput struct{}

func NewListInterfacesHandler(exec executor.CommandExecutor) func(context.Context, *mcp.CallToolRequest, ListInterfacesInput) (*mcp.CallToolResult, struct{}, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ ListInterfacesInput) (*mcp.CallToolResult, struct{}, error) {
		stdout, _, err := exec.Execute(ctx, "tshark", []string{"-D"})
		if err != nil {
			return errorResult("Failed to list interfaces: " + err.Error()), struct{}{}, nil
		}

		result := strings.TrimSpace(string(stdout))
		if result == "" {
			return textResult("No network interfaces found. Ensure you have appropriate permissions."), struct{}{}, nil
		}

		formatted := output.FormatResult(result, map[string]string{
			"Command": "tshark -D",
		})

		return textResult(formatted), struct{}{}, nil
	}
}
