package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ConnectionsService provides access to active connections via Clash API.
type ConnectionsService struct {
	baseURL string
	client  *http.Client
}

// NewConnectionsService creates a new ConnectionsService.
func NewConnectionsService() *ConnectionsService {
	return &ConnectionsService{
		baseURL: "http://127.0.0.1:9090",
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Connection represents a single active connection.
type Connection struct {
	ID          string   `json:"id"`
	Metadata    Metadata `json:"metadata"`
	Upload      int64    `json:"upload"`
	Download    int64    `json:"download"`
	Start       string   `json:"start"`
	Chains      []string `json:"chains"`
	Rule        string   `json:"rule"`
	RulePayload string   `json:"rulePayload"`
}

// Metadata contains connection metadata.
type Metadata struct {
	Network         string `json:"network"`
	Type            string `json:"type"`
	SourceIP        string `json:"sourceIP"`
	DestinationIP   string `json:"destinationIP"`
	SourcePort      string `json:"sourcePort"`
	DestinationPort string `json:"destinationPort"`
	Host            string `json:"host"`
	DNSMode         string `json:"dnsMode"`
	ProcessPath     string `json:"processPath"`
}

// ConnectionsResponse is the response from GET /connections.
type ConnectionsResponse struct {
	DownloadTotal int64        `json:"downloadTotal"`
	UploadTotal   int64        `json:"uploadTotal"`
	Connections   []Connection `json:"connections"`
}

// GetConnections fetches all active connections.
func (s *ConnectionsService) GetConnections() (*ConnectionsResponse, error) {
	resp, err := s.client.Get(s.baseURL + "/connections")
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result ConnectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// CloseConnection closes a specific connection by ID.
func (s *ConnectionsService) CloseConnection(id string) error {
	req, err := http.NewRequest(http.MethodDelete, s.baseURL+"/connections/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// CloseAllConnections closes all active connections.
func (s *ConnectionsService) CloseAllConnections() error {
	req, err := http.NewRequest(http.MethodDelete, s.baseURL+"/connections", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to close all connections: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}
