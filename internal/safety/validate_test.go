package safety

import (
	"testing"
	"time"
)

func TestValidateFilePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid pcap", "/tmp/capture.pcap", false},
		{"valid pcapng", "/tmp/capture.pcapng", false},
		{"valid cap", "/tmp/capture.cap", false},
		{"valid gz", "/tmp/capture.pcap.gz", false},
		{"empty path", "", true},
		{"path traversal", "/tmp/../etc/passwd.pcap", true},
		{"invalid extension", "/tmp/capture.txt", true},
		{"no extension", "/tmp/capture", true},
		{"windows path", `C:\captures\test.pcap`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeDisplayFilter(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		wantErr bool
	}{
		{"empty filter", "", false},
		{"simple protocol", "tcp", false},
		{"ip comparison", "ip.addr == 192.168.1.1", false},
		{"port filter", "tcp.port == 80", false},
		{"complex filter", "tcp.port == 80 && http.request.method == \"GET\"", false},
		{"dangerous semicolon", "tcp; rm -rf /", true},
		{"dangerous backtick", "tcp`whoami`", true},
		{"dangerous dollar", "tcp$HOME", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SanitizeDisplayFilter(tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeDisplayFilter(%q) error = %v, wantErr %v", tt.filter, err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeCaptureFilter(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		wantErr bool
	}{
		{"empty", "", false},
		{"simple", "tcp port 80", false},
		{"host filter", "host 192.168.1.1", false},
		{"dangerous semicolon", "tcp; rm -rf /", true},
		{"dangerous backtick", "tcp`id`", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SanitizeCaptureFilter(tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeCaptureFilter(%q) error = %v, wantErr %v", tt.filter, err, tt.wantErr)
			}
		})
	}
}

func TestClampTimeout(t *testing.T) {
	tests := []struct {
		name      string
		requested time.Duration
		max       time.Duration
		want      time.Duration
	}{
		{"zero returns max", 0, 30 * time.Second, 30 * time.Second},
		{"negative returns max", -1, 30 * time.Second, 30 * time.Second},
		{"within limit", 10 * time.Second, 30 * time.Second, 10 * time.Second},
		{"exceeds limit", 60 * time.Second, 30 * time.Second, 30 * time.Second},
		{"exact limit", 30 * time.Second, 30 * time.Second, 30 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClampTimeout(tt.requested, tt.max)
			if got != tt.want {
				t.Errorf("ClampTimeout(%v, %v) = %v, want %v", tt.requested, tt.max, got, tt.want)
			}
		})
	}
}

func TestClampPacketCount(t *testing.T) {
	tests := []struct {
		name      string
		requested int
		want      int
	}{
		{"zero returns max", 0, MaxPacketCount},
		{"negative returns max", -1, MaxPacketCount},
		{"within limit", 100, 100},
		{"exceeds limit", 5000, MaxPacketCount},
		{"exact limit", MaxPacketCount, MaxPacketCount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClampPacketCount(tt.requested)
			if got != tt.want {
				t.Errorf("ClampPacketCount(%d) = %d, want %d", tt.requested, got, tt.want)
			}
		})
	}
}

func TestSanitizeInterfaceName(t *testing.T) {
	tests := []struct {
		name    string
		iface   string
		wantErr bool
	}{
		{"valid eth0", "eth0", false},
		{"valid loopback", "lo", false},
		{"empty", "", true},
		{"dangerous", "eth0;whoami", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SanitizeInterfaceName(tt.iface)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeInterfaceName(%q) error = %v, wantErr %v", tt.iface, err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeFieldName(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		wantErr bool
	}{
		{"ip.src", "ip.src", false},
		{"tcp.port", "tcp.port", false},
		{"http.request.method", "http.request.method", false},
		{"empty", "", true},
		{"injection attempt", "ip.src; rm", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SanitizeFieldName(tt.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeFieldName(%q) error = %v, wantErr %v", tt.field, err, tt.wantErr)
			}
		})
	}
}

func TestValidateOutputPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"empty is ok", "", false},
		{"valid pcap", "/tmp/output.pcap", false},
		{"valid pcapng", "/tmp/output.pcapng", false},
		{"path traversal", "/tmp/../etc/output.pcap", true},
		{"invalid extension", "/tmp/output.csv", true},
		{"windows path", `C:\captures\out.pcap`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOutputPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutputPath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeStatType(t *testing.T) {
	tests := []struct {
		name    string
		stat    string
		wantErr bool
	}{
		{"protocol hierarchy", "io,phs", false},
		{"tcp conversations", "conv,tcp", false},
		{"udp conversations", "conv,udp", false},
		{"ip endpoints", "endpoints,ip", false},
		{"io stat interval", "io,stat,1", false},
		{"empty", "", true},
		{"invalid chars", "io,phs\x00", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SanitizeStatType(tt.stat)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeStatType(%q) error = %v, wantErr %v", tt.stat, err, tt.wantErr)
			}
		})
	}
}
