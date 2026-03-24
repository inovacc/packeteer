package cmd

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/inovacc/packeteer/internal/executor"
	"github.com/inovacc/packeteer/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Packeteer MCP server",
	Long: `Start the MCP server that exposes Wireshark CLI tools for packet capture
and network analysis.

Transports:
  stdio (default) — communicates over stdin/stdout using JSON-RPC
  http            — serves Streamable HTTP on the specified port

Examples:
  # Stdio transport (Claude Desktop)
  packeteer serve

  # HTTP transport (remote connections)
  packeteer serve --transport http --port 8080

Claude Desktop config (stdio):
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
		transport, _ := cmd.Flags().GetString("transport")
		port, _ := cmd.Flags().GetInt("port")

		logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		exec := executor.NewRealExecutor(logger, wiresharkDir)

		var opts []server.Option
		if captureDir != "" {
			opts = append(opts, server.WithCaptureDir(captureDir))
		}

		srv := server.New(exec, logger, opts...)

		switch transport {
		case "stdio":
			logger.Info("starting packeteer MCP server", "transport", "stdio")
			return srv.Run(cmd.Context(), &mcp.StdioTransport{})

		case "http":
			addr := fmt.Sprintf(":%d", port)
			logger.Info("starting packeteer MCP server", "transport", "http", "addr", addr)

			handler := mcp.NewStreamableHTTPHandler(
				func(_ *http.Request) *mcp.Server { return srv },
				&mcp.StreamableHTTPOptions{
					Logger: logger,
				},
			)

			httpServer := &http.Server{
				Addr:    addr,
				Handler: handler,
			}

			go func() {
				<-cmd.Context().Done()
				_ = httpServer.Close()
			}()

			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return fmt.Errorf("HTTP server failed: %w", err)
			}
			return nil

		default:
			return fmt.Errorf("unknown transport %q: use 'stdio' or 'http'", transport)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("wireshark-dir", "", "path to Wireshark installation directory (auto-detected if not set)")
	serveCmd.Flags().String("capture-dir", "", "directory containing pcap files for MCP resource browsing")
	serveCmd.Flags().String("transport", "stdio", "transport type: 'stdio' or 'http'")
	serveCmd.Flags().Int("port", 8080, "HTTP server port (only used with --transport http)")
}
