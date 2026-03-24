package cmd

import (
	"log/slog"
	"os"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Packeteer MCP server on stdio transport",
	Long: `Start the MCP server that exposes Wireshark CLI tools for packet capture
and network analysis. Communicates over stdin/stdout using JSON-RPC.

Add to your Claude Desktop config:
  {
    "mcpServers": {
      "packeteer": {
        "command": "packeteer",
        "args": ["serve"]
      }
    }
  }`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		wiresharkDir, _ := cmd.Flags().GetString("wireshark-dir")
		captureDir, _ := cmd.Flags().GetString("capture-dir")

		logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		exec := executor.NewRealExecutor(logger, wiresharkDir)

		var opts []server.Option
		if captureDir != "" {
			opts = append(opts, server.WithCaptureDir(captureDir))
		}

		srv := server.New(exec, logger, opts...)

		logger.Info("starting packeteer MCP server", "transport", "stdio")

		return srv.Run(cmd.Context(), &mcp.StdioTransport{})
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("wireshark-dir", "", "path to Wireshark installation directory (auto-detected if not set)")
	serveCmd.Flags().String("capture-dir", "", "directory containing pcap files for MCP resource browsing")
}
