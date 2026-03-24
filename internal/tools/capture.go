package tools

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/inovacc/sharkline/internal/executor"
	"github.com/inovacc/sharkline/internal/output"
	"github.com/inovacc/sharkline/internal/safety"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CaptureInput struct {
	Interface     string `json:"interface" jsonschema:"network interface name or index (required)"`
	CaptureFilter string `json:"capture_filter,omitempty" jsonschema:"BPF capture filter expression (e.g. 'tcp port 80')"`
	DisplayFilter string `json:"display_filter,omitempty" jsonschema:"Wireshark display filter (e.g. 'http.request')"`
	Duration      int    `json:"duration,omitempty" jsonschema:"capture duration in seconds (max 30, default 10)"`
	PacketCount   int    `json:"packet_count,omitempty" jsonschema:"max packets to capture (max 1000, default 100)"`
	OutputFile    string `json:"output_file,omitempty" jsonschema:"path to save pcap file (.pcap or .pcapng)"`
	Summarize     bool   `json:"summarize,omitempty" jsonschema:"parse JSON into structured packet summaries (default false)"`
}

// NewCaptureHandler creates a capture tool handler. If limiter is non-nil,
// concurrent captures are limited.
func NewCaptureHandler(exec executor.CommandExecutor, limiter *safety.CaptureLimiter) func(context.Context, *mcp.CallToolRequest, CaptureInput) (*mcp.CallToolResult, struct{}, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input CaptureInput) (*mcp.CallToolResult, struct{}, error) {
		if limiter != nil {
			release, err := limiter.Acquire(ctx)
			if err != nil {
				return errorResult(fmt.Sprintf("Capture limit reached: %v (active: %d/%d)", err, limiter.Active(), limiter.Max())), struct{}{}, nil
			}
			defer release()
		}

		if err := safety.SanitizeInterfaceName(input.Interface); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}
		if err := safety.SanitizeCaptureFilter(input.CaptureFilter); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}
		if err := safety.SanitizeDisplayFilter(input.DisplayFilter); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}
		if err := safety.ValidateOutputPath(input.OutputFile); err != nil {
			return errorResult(err.Error()), struct{}{}, nil
		}

		duration := safety.ClampTimeout(time.Duration(input.Duration)*time.Second, safety.MaxCaptureTimeout)
		count := input.PacketCount
		if count <= 0 {
			count = 100
		}
		count = safety.ClampPacketCount(count)

		args := []string{
			"-i", input.Interface,
			"-c", strconv.Itoa(count),
			"-a", fmt.Sprintf("duration:%d", int(duration.Seconds())),
			"-T", "json",
		}

		if input.CaptureFilter != "" {
			args = append(args, "-f", input.CaptureFilter)
		}
		if input.DisplayFilter != "" {
			args = append(args, "-Y", input.DisplayFilter)
		}
		if input.OutputFile != "" {
			args = append(args, "-w", input.OutputFile)
		}

		captureCtx, cancel := context.WithTimeout(ctx, duration+5*time.Second)
		defer cancel()

		stdout, _, err := exec.Execute(captureCtx, "tshark", args)
		if err != nil {
			return errorResult("Capture failed: " + err.Error()), struct{}{}, nil
		}

		var resultBytes []byte
		if input.Summarize {
			parsed, parseErr := output.ParseTSharkJSON(stdout, count)
			if parseErr != nil {
				resultBytes = stdout
			} else {
				resultBytes = parsed
			}
		} else {
			resultBytes = stdout
		}

		truncated, wasTruncated := output.Truncate(resultBytes, output.DefaultMaxBytes)
		metadata := map[string]string{
			"Interface":    input.Interface,
			"Max Packets":  strconv.Itoa(count),
			"Max Duration": duration.String(),
			"Truncated":    strconv.FormatBool(wasTruncated),
		}

		return textResult(output.FormatResult(string(truncated), metadata)), struct{}{}, nil
	}
}
