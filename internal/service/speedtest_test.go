package service

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RayUI/RayUI/internal/model"
)

func TestParseProxyURL(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"socks5://127.0.0.1:1080", "127.0.0.1:1080", false},
		{"socks5://localhost:10808", "localhost:10808", false},
		{"http://proxy.example.com:8080", "proxy.example.com:8080", false},
		{"://invalid", "", true},
	}

	for _, tt := range tests {
		got, err := ParseProxyURL(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseProxyURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseProxyURL(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestDefaultProxyTestConfig(t *testing.T) {
	cfg := DefaultProxyTestConfig()

	if cfg.SocksPort != 10808 {
		t.Errorf("SocksPort = %d, want 10808", cfg.SocksPort)
	}
	if cfg.PingURL != DefaultPingURL {
		t.Errorf("PingURL = %q, want %q", cfg.PingURL, DefaultPingURL)
	}
	if cfg.DownloadURL != "" {
		t.Errorf("DownloadURL = %q, want empty", cfg.DownloadURL)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s", cfg.Timeout)
	}
}

func TestTCPPing(t *testing.T) {
	// Start a local TCP listener
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	// Accept connections in background
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	svc := &SpeedTestService{}

	// Successful ping
	latency := svc.TCPPing("127.0.0.1", addr.Port, 5*time.Second)
	if latency < 0 {
		t.Errorf("TCPPing() = %d, want >= 0", latency)
	}

	// Timeout on unreachable port
	latency = svc.TCPPing("127.0.0.1", 1, 100*time.Millisecond)
	if latency != -1 {
		t.Errorf("TCPPing(unreachable) = %d, want -1", latency)
	}
}

func TestTestProfile(t *testing.T) {
	// Start a TCP listener to ping
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	svc := &SpeedTestService{}

	profile := model.ProfileItem{
		ID:      "test-1",
		Address: "127.0.0.1",
		Port:    addr.Port,
	}

	result := svc.TestProfile(profile, 5*time.Second)
	if result.ProfileID != "test-1" {
		t.Errorf("ProfileID = %q, want test-1", result.ProfileID)
	}
	if result.Latency < 0 {
		t.Errorf("Latency = %d, want >= 0", result.Latency)
	}
}

func TestTestProfiles(t *testing.T) {
	// Start multiple TCP listeners
	listeners := make([]net.Listener, 3)
	for i := range listeners {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen[%d]: %v", i, err)
		}
		listeners[i] = ln
		go func(l net.Listener) {
			for {
				conn, err := l.Accept()
				if err != nil {
					return
				}
				conn.Close()
			}
		}(ln)
	}
	defer func() {
		for _, ln := range listeners {
			ln.Close()
		}
	}()

	profiles := make([]model.ProfileItem, len(listeners))
	for i, ln := range listeners {
		addr := ln.Addr().(*net.TCPAddr)
		profiles[i] = model.ProfileItem{
			ID:      fmt.Sprintf("profile-%d", i),
			Address: "127.0.0.1",
			Port:    addr.Port,
		}
	}

	svc := &SpeedTestService{}
	results := svc.TestProfiles(profiles, 2, 5*time.Second)

	if len(results) != len(profiles) {
		t.Fatalf("len(results) = %d, want %d", len(results), len(profiles))
	}

	for i, r := range results {
		if r.ProfileID != profiles[i].ID {
			t.Errorf("results[%d].ProfileID = %q, want %q", i, r.ProfileID, profiles[i].ID)
		}
		if r.Latency < 0 {
			t.Errorf("results[%d].Latency = %d, want >= 0", i, r.Latency)
		}
	}
}

func TestTestProfilesConcurrencyDefault(t *testing.T) {
	svc := &SpeedTestService{}
	// Empty profiles with concurrent=0 should use default 4
	results := svc.TestProfiles(nil, 0, time.Second)
	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}

func TestTestProfileThroughProxy(t *testing.T) {
	// Start a TCP listener as fallback
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	svc := &SpeedTestService{}

	profile := model.ProfileItem{
		ID:      "test-proxy",
		Address: "127.0.0.1",
		Port:    addr.Port,
	}

	// With empty pingURL, falls back to TCP ping
	result := svc.TestProfileThroughProxy(profile, 10808, "", "", time.Second)
	if result.ProfileID != "test-proxy" {
		t.Errorf("ProfileID = %q, want test-proxy", result.ProfileID)
	}
	// Latency should be >= 0 from TCP ping (direct connection to our listener)
	if result.Latency < 0 {
		t.Errorf("Latency = %d, want >= 0 (TCP ping fallback)", result.Latency)
	}
}

func TestConstants(t *testing.T) {
	if DefaultPingURL == "" {
		t.Error("DefaultPingURL should not be empty")
	}
	if DefaultDownloadURL == "" {
		t.Error("DefaultDownloadURL should not be empty")
	}
}

func TestRealPingNoProxy(t *testing.T) {
	// Create a test HTTP server
	_ = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	svc := &SpeedTestService{}

	// RealPing to a non-existent SOCKS proxy should return -1
	latency := svc.RealPing("127.0.0.1:1", "http://example.com", 500*time.Millisecond)
	if latency != -1 {
		t.Errorf("RealPing(invalid proxy) = %d, want -1", latency)
	}
}
