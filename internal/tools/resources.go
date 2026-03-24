package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/inovacc/sharkline/internal/executor"
	"github.com/inovacc/sharkline/internal/output"
	"github.com/inovacc/sharkline/internal/safety"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterResources adds MCP resources and resource templates to the server.
func RegisterResources(server *mcp.Server, exec executor.CommandExecutor, captureDir string) {
	// Resource template for pcap files.
	server.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: "sharkline://captures/{filename}",
			Name:        "Capture File",
			Description: "Read metadata and summary of a pcap/pcapng capture file from the captures directory",
			MIMEType:    "text/plain",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			uri := req.Params.URI
			filename := strings.TrimPrefix(uri, "sharkline://captures/")

			if filename == "" || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
				return nil, mcp.ResourceNotFoundError(uri)
			}

			ext := strings.ToLower(filepath.Ext(filename))
			allowed := map[string]bool{".pcap": true, ".pcapng": true, ".cap": true}
			if !allowed[ext] {
				return nil, mcp.ResourceNotFoundError(uri)
			}

			filePath := filepath.Join(captureDir, filename)
			if err := safety.ValidateFilePath(filePath); err != nil {
				return nil, mcp.ResourceNotFoundError(uri)
			}

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return nil, mcp.ResourceNotFoundError(uri)
			}

			// Get capinfos summary.
			stdout, _, err := exec.Execute(ctx, "capinfos", []string{filePath})
			if err != nil {
				return &mcp.ReadResourceResult{
					Contents: []*mcp.ResourceContents{
						{URI: uri, MIMEType: "text/plain", Text: "Error reading capture info: " + err.Error()},
					},
				}, nil
			}

			truncated, _ := output.Truncate(stdout, output.DefaultMaxBytes)

			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{
					{URI: uri, MIMEType: "text/plain", Text: string(truncated)},
				},
			}, nil
		},
	)

	// Static resource listing available captures.
	server.AddResource(
		&mcp.Resource{
			URI:         "sharkline://captures",
			Name:        "Available Captures",
			Description: "List all pcap/pcapng files in the captures directory",
			MIMEType:    "text/plain",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			if captureDir == "" {
				return &mcp.ReadResourceResult{
					Contents: []*mcp.ResourceContents{
						{URI: req.Params.URI, MIMEType: "text/plain", Text: "No captures directory configured. Use --capture-dir flag."},
					},
				}, nil
			}

			entries, err := os.ReadDir(captureDir)
			if err != nil {
				return &mcp.ReadResourceResult{
					Contents: []*mcp.ResourceContents{
						{URI: req.Params.URI, MIMEType: "text/plain", Text: "Error reading captures directory: " + err.Error()},
					},
				}, nil
			}

			var files []string
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				ext := strings.ToLower(filepath.Ext(entry.Name()))
				if ext == ".pcap" || ext == ".pcapng" || ext == ".cap" {
					info, _ := entry.Info()
					size := "unknown"
					if info != nil {
						size = formatSize(info.Size())
					}
					files = append(files, entry.Name()+" ("+size+")")
				}
			}

			if len(files) == 0 {
				return &mcp.ReadResourceResult{
					Contents: []*mcp.ResourceContents{
						{URI: req.Params.URI, MIMEType: "text/plain", Text: "No capture files found in: " + captureDir},
					},
				}, nil
			}

			result := "Capture files in " + captureDir + ":\n\n"
			for _, f := range files {
				result += "  - " + f + "\n"
			}

			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{
					{URI: req.Params.URI, MIMEType: "text/plain", Text: result},
				},
			}, nil
		},
	)
}

func formatSize(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
