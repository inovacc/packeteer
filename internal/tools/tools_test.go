package tools

import (
	"context"
	"testing"

	"github.com/inovacc/sharkline/internal/executor"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestListInterfaces(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["tshark"] = executor.MockResponse{
		Stdout: []byte("1. eth0\n2. lo (Loopback)\n3. wlan0\n"),
	}

	handler := NewListInterfacesHandler(mock)
	result, _, err := handler(context.Background(), nil, ListInterfacesInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got error")
	}

	text := result.Content[0].(*mcp.TextContent).Text
	if text == "" {
		t.Fatal("expected non-empty result")
	}

	if len(mock.Calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(mock.Calls))
	}
	if mock.Calls[0].Binary != "tshark" {
		t.Fatalf("expected tshark, got %s", mock.Calls[0].Binary)
	}
}

func TestReadPcap(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["tshark"] = executor.MockResponse{
		Stdout: []byte(`[{"_index":"packets-1","_source":{"layers":{"frame":{"number":"1"}}}}]`),
	}

	handler := NewReadPcapHandler(mock)

	t.Run("valid pcap", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, ReadPcapInput{
			FilePath:   "/tmp/test.pcap",
			MaxPackets: 10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.IsError {
			t.Fatal("expected success")
		}
	})

	t.Run("invalid extension", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, ReadPcapInput{
			FilePath: "/tmp/test.txt",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected error for invalid extension")
		}
	})

	t.Run("path traversal", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, ReadPcapInput{
			FilePath: "/tmp/../etc/test.pcap",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected error for path traversal")
		}
	})
}

func TestCapturePackets(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["tshark"] = executor.MockResponse{
		Stdout: []byte(`[{"_source":{"layers":{"frame":{"number":"1"}}}}]`),
	}

	handler := NewCaptureHandler(mock, nil)

	t.Run("valid capture", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, CaptureInput{
			Interface:   "eth0",
			PacketCount: 10,
			Duration:    5,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.IsError {
			t.Fatal("expected success")
		}
	})

	t.Run("dangerous interface name", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, CaptureInput{
			Interface: "eth0;whoami",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected error for dangerous interface name")
		}
	})
}

func TestExtractFields(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["tshark"] = executor.MockResponse{
		Stdout: []byte("192.168.1.1\t192.168.1.2\t80\n"),
	}

	handler := NewExtractFieldsHandler(mock)

	t.Run("valid extraction", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, ExtractFieldsInput{
			FilePath: "/tmp/test.pcap",
			Fields:   []string{"ip.src", "ip.dst", "tcp.port"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.IsError {
			t.Fatal("expected success")
		}
	})

	t.Run("no fields", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, ExtractFieldsInput{
			FilePath: "/tmp/test.pcap",
			Fields:   []string{},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected error for empty fields")
		}
	})

	t.Run("dangerous field name", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, ExtractFieldsInput{
			FilePath: "/tmp/test.pcap",
			Fields:   []string{"ip.src; rm -rf /"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected error for dangerous field name")
		}
	})
}

func TestStatistics(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["tshark"] = executor.MockResponse{
		Stdout: []byte("Protocol Hierarchy Statistics\neth  100.00%\n  ip  95.00%\n"),
	}

	handler := NewStatisticsHandler(mock)

	result, _, err := handler(context.Background(), nil, StatisticsInput{
		FilePath: "/tmp/test.pcap",
		StatType: "io,phs",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestCaptureInfo(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["capinfos"] = executor.MockResponse{
		Stdout: []byte("File name:           test.pcap\nPackets:             1234\nCapture duration:    10.5 seconds\n"),
	}

	handler := NewCaptureInfoHandler(mock)

	result, _, err := handler(context.Background(), nil, CaptureInfoInput{
		FilePath: "/tmp/test.pcap",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestMergePcaps(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["mergecap"] = executor.MockResponse{}

	handler := NewMergePcapsHandler(mock)

	t.Run("valid merge", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, MergePcapsInput{
			InputFiles: []string{"/tmp/a.pcap", "/tmp/b.pcap"},
			OutputFile: "/tmp/merged.pcap",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.IsError {
			t.Fatal("expected success")
		}
	})

	t.Run("too few files", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, MergePcapsInput{
			InputFiles: []string{"/tmp/a.pcap"},
			OutputFile: "/tmp/merged.pcap",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected error for single file")
		}
	})
}

func TestListProtocols(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["tshark"] = executor.MockResponse{
		Stdout: []byte("TCP\ttcp\tTransmission Control Protocol\nUDP\tudp\tUser Datagram Protocol\nHTTP\thttp\tHypertext Transfer Protocol\n"),
	}

	handler := NewListProtocolsHandler(mock)

	t.Run("no filter", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, ListProtocolsInput{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.IsError {
			t.Fatal("expected success")
		}
	})

	t.Run("with filter", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, ListProtocolsInput{
			Filter: "tcp",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.IsError {
			t.Fatal("expected success")
		}
	})
}

func TestFilterPcap(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["editcap"] = executor.MockResponse{}

	handler := NewFilterPcapHandler(mock)

	t.Run("valid filter", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, FilterPcapInput{
			InputFile:  "/tmp/input.pcap",
			OutputFile: "/tmp/output.pcap",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.IsError {
			t.Fatal("expected success")
		}
	})

	t.Run("with time range", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, FilterPcapInput{
			InputFile:  "/tmp/input.pcap",
			OutputFile: "/tmp/output.pcap",
			StartTime:  "2024-01-01 00:00:00",
			EndTime:    "2024-01-01 01:00:00",
			MaxPackets: 500,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.IsError {
			t.Fatal("expected success")
		}
	})

	t.Run("missing output file", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, FilterPcapInput{
			InputFile: "/tmp/input.pcap",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected error for missing output file")
		}
	})

	t.Run("invalid input extension", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, FilterPcapInput{
			InputFile:  "/tmp/input.txt",
			OutputFile: "/tmp/output.pcap",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected error for invalid input extension")
		}
	})

	t.Run("path traversal in output", func(t *testing.T) {
		result, _, err := handler(context.Background(), nil, FilterPcapInput{
			InputFile:  "/tmp/input.pcap",
			OutputFile: "/tmp/../etc/output.pcap",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsError {
			t.Fatal("expected error for path traversal")
		}
	})
}

func TestDecodePacket(t *testing.T) {
	mock := executor.NewMockExecutor()
	mock.Responses["tshark"] = executor.MockResponse{
		Stdout: []byte("Frame 1: 66 bytes on wire\n  Ethernet II\n    Source: 00:11:22:33:44:55\n"),
	}

	handler := NewDecodePacketHandler(mock)

	result, _, err := handler(context.Background(), nil, DecodePacketInput{
		FilePath:     "/tmp/test.pcap",
		PacketNumber: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}
