package proxycore

import (
	"net"
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid https", "https://example.com", false},
		{"valid http", "http://example.com/path?q=1", false},
		{"valid with port", "https://example.com:8080/path", false},
		{"valid IP", "http://8.8.8.8", false},

		// Blocked schemes
		{"ftp", "ftp://example.com", true},
		{"file", "file:///etc/passwd", true},
		{"javascript", "javascript:alert(1)", true},
		{"no scheme", "example.com", true},
		{"empty", "", true},

		// Blocked hosts
		{"localhost", "http://localhost", true},
		{"localhost port", "http://localhost:8080", true},
		{".local", "http://myhost.local", true},
		{".localhost", "http://sub.localhost", true},
		{".internal", "http://metadata.internal", true},
		{"metadata hostname", "http://metadata", true},
		{"metadata.google.internal", "http://metadata.google.internal", true},

		// Blocked IPs - private ranges
		{"127.0.0.1", "http://127.0.0.1", true},
		{"127.0.0.2", "http://127.0.0.2", true},
		{"10.0.0.1", "http://10.0.0.1", true},
		{"172.16.0.1", "http://172.16.0.1", true},
		{"192.168.1.1", "http://192.168.1.1", true},
		{"0.0.0.0", "http://0.0.0.0", true},

		// Link-local
		{"169.254.1.1", "http://169.254.1.1", true},

		// Cloud metadata
		{"AWS metadata", "http://169.254.169.254", true},
		{"AWS ECS", "http://169.254.170.2", true},
		{"Alibaba", "http://100.100.100.200", true},

		// IPv6
		{"IPv6 loopback", "http://[::1]", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestValidateIP(t *testing.T) {
	tests := []struct {
		ip      string
		wantErr bool
	}{
		{"8.8.8.8", false},
		{"1.1.1.1", false},
		{"142.250.185.206", false},
		{"127.0.0.1", true},
		{"10.0.0.1", true},
		{"192.168.0.1", true},
		{"169.254.169.254", true},
		{"0.0.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if ip == nil {
				t.Fatalf("failed to parse IP %s", tt.ip)
			}
			err := validateIP(ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateIP(%s) error = %v, wantErr %v", tt.ip, err, tt.wantErr)
			}
		})
	}
}
