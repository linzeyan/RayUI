package service

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/RayUI/RayUI/internal/model"
	"golang.org/x/net/proxy"
)

// SpeedTestService provides latency and speed testing.
type SpeedTestService struct{}

// TCPPing measures TCP connection round-trip time in milliseconds.
// Returns -1 on timeout.
func (s *SpeedTestService) TCPPing(address string, port int, timeout time.Duration) int {
	addr := net.JoinHostPort(address, fmt.Sprintf("%d", port))
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return -1
	}
	conn.Close()
	return int(time.Since(start).Milliseconds())
}

// RealPing performs an HTTP HEAD request through a SOCKS5 proxy to measure latency.
// Returns -1 on timeout or error.
func (s *SpeedTestService) RealPing(socksAddr string, testURL string, timeout time.Duration) int {
	dialer, err := proxy.SOCKS5("tcp", socksAddr, nil, proxy.Direct)
	if err != nil {
		return -1
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	start := time.Now()
	req, err := http.NewRequest(http.MethodHead, testURL, nil)
	if err != nil {
		return -1
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1
	}
	resp.Body.Close()

	return int(time.Since(start).Milliseconds())
}

// BandwidthTest downloads data through a SOCKS5 proxy and measures speed.
// Returns bytes per second, or 0 on error.
func (s *SpeedTestService) BandwidthTest(socksAddr string, downloadURL string, timeout time.Duration) int64 {
	dialer, err := proxy.SOCKS5("tcp", socksAddr, nil, proxy.Direct)
	if err != nil {
		return 0
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	start := time.Now()
	resp, err := client.Get(downloadURL)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	// Read up to 10MB or until timeout
	const maxBytes = 10 * 1024 * 1024
	reader := io.LimitReader(resp.Body, maxBytes)

	n, err := io.Copy(io.Discard, reader)
	if err != nil && n == 0 {
		return 0
	}

	elapsed := time.Since(start).Seconds()
	if elapsed == 0 {
		return 0
	}

	return int64(float64(n) / elapsed)
}

// TestProfile runs a TCP ping test on a single profile.
func (s *SpeedTestService) TestProfile(profile model.ProfileItem, timeout time.Duration) model.SpeedTestResult {
	latency := s.TCPPing(profile.Address, profile.Port, timeout)
	return model.SpeedTestResult{
		ProfileID: profile.ID,
		Latency:   latency,
	}
}

// TestProfileThroughProxy tests a profile using real HTTP requests through a local proxy.
// socksPort is the local SOCKS5 proxy port (e.g., 10808).
// pingURL is the URL for latency test (HTTP HEAD).
// downloadURL is the URL for bandwidth test (optional, empty to skip).
func (s *SpeedTestService) TestProfileThroughProxy(profile model.ProfileItem, socksPort int, pingURL string, downloadURL string, timeout time.Duration) model.SpeedTestResult {
	socksAddr := fmt.Sprintf("127.0.0.1:%d", socksPort)

	result := model.SpeedTestResult{
		ProfileID: profile.ID,
		Latency:   -1,
		Speed:     0,
	}

	// Test real latency through proxy
	if pingURL != "" {
		result.Latency = s.RealPing(socksAddr, pingURL, timeout)
	} else {
		// Fallback to TCP ping
		result.Latency = s.TCPPing(profile.Address, profile.Port, timeout)
	}

	// Test bandwidth if downloadURL is provided and latency test succeeded
	if downloadURL != "" && result.Latency > 0 {
		result.Speed = s.BandwidthTest(socksAddr, downloadURL, timeout)
	}

	return result
}

// TestProfiles runs TCP ping tests on multiple profiles concurrently.
func (s *SpeedTestService) TestProfiles(profiles []model.ProfileItem, concurrent int, timeout time.Duration) []model.SpeedTestResult {
	if concurrent <= 0 {
		concurrent = 4
	}

	results := make([]model.SpeedTestResult, len(profiles))
	sem := make(chan struct{}, concurrent)
	var wg sync.WaitGroup

	for i, p := range profiles {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, profile model.ProfileItem) {
			defer wg.Done()
			defer func() { <-sem }()
			results[idx] = s.TestProfile(profile, timeout)
		}(i, p)
	}

	wg.Wait()
	return results
}

// DefaultPingURL is the default URL for real ping tests.
const DefaultPingURL = "https://www.gstatic.com/generate_204"

// DefaultDownloadURL is the default URL for bandwidth tests (Cloudflare speed test).
const DefaultDownloadURL = "https://speed.cloudflare.com/__down?bytes=10000000"

// ProxyTestConfig holds configuration for proxy-based speed tests.
type ProxyTestConfig struct {
	SocksPort   int
	PingURL     string
	DownloadURL string
	Timeout     time.Duration
}

// DefaultProxyTestConfig returns a default configuration for proxy testing.
func DefaultProxyTestConfig() ProxyTestConfig {
	return ProxyTestConfig{
		SocksPort:   10808,
		PingURL:     DefaultPingURL,
		DownloadURL: "", // Empty means skip bandwidth test
		Timeout:     10 * time.Second,
	}
}

// ParseProxyURL parses a SOCKS5 proxy URL and returns the address.
func ParseProxyURL(proxyURL string) (string, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}
