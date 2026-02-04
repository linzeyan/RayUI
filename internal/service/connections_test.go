package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestConnectionsService(serverURL string) *ConnectionsService {
	return &ConnectionsService{
		baseURL: serverURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func TestGetConnections(t *testing.T) {
	expected := ConnectionsResponse{
		DownloadTotal: 1024,
		UploadTotal:   512,
		Connections: []Connection{
			{
				ID:       "abc-123",
				Upload:   100,
				Download: 200,
				Chains:   []string{"DIRECT"},
				Rule:     "MATCH",
				Metadata: Metadata{
					Network:         "tcp",
					Host:            "example.com",
					DestinationPort: "443",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/connections" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	svc := newTestConnectionsService(server.URL)
	resp, err := svc.GetConnections()
	if err != nil {
		t.Fatalf("GetConnections() error: %v", err)
	}

	if resp.DownloadTotal != 1024 {
		t.Errorf("DownloadTotal = %d, want 1024", resp.DownloadTotal)
	}
	if resp.UploadTotal != 512 {
		t.Errorf("UploadTotal = %d, want 512", resp.UploadTotal)
	}
	if len(resp.Connections) != 1 {
		t.Fatalf("len(Connections) = %d, want 1", len(resp.Connections))
	}
	conn := resp.Connections[0]
	if conn.ID != "abc-123" {
		t.Errorf("conn.ID = %q, want abc-123", conn.ID)
	}
	if conn.Metadata.Host != "example.com" {
		t.Errorf("conn.Metadata.Host = %q, want example.com", conn.Metadata.Host)
	}
}

func TestGetConnectionsServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := newTestConnectionsService(server.URL)
	_, err := svc.GetConnections()
	if err == nil {
		t.Error("GetConnections() should return error on 500")
	}
}

func TestCloseConnection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/connections/abc-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := newTestConnectionsService(server.URL)
	err := svc.CloseConnection("abc-123")
	if err != nil {
		t.Fatalf("CloseConnection() error: %v", err)
	}
}

func TestCloseConnectionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	svc := newTestConnectionsService(server.URL)
	err := svc.CloseConnection("bad-id")
	if err == nil {
		t.Error("CloseConnection() should return error on 400")
	}
}

func TestCloseAllConnections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/connections" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := newTestConnectionsService(server.URL)
	err := svc.CloseAllConnections()
	if err != nil {
		t.Fatalf("CloseAllConnections() error: %v", err)
	}
}

func TestCloseAllConnectionsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	svc := newTestConnectionsService(server.URL)
	err := svc.CloseAllConnections()
	if err == nil {
		t.Error("CloseAllConnections() should return error on 403")
	}
}

func TestNewConnectionsService(t *testing.T) {
	svc := NewConnectionsService()
	if svc.baseURL != "http://127.0.0.1:9090" {
		t.Errorf("baseURL = %q, want http://127.0.0.1:9090", svc.baseURL)
	}
	if svc.client == nil {
		t.Error("client should not be nil")
	}
}
