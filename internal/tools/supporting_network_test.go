package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// --- list_wans tests ---

func TestListWans_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":   "aaa00000-0000-0000-0000-000000000001",
					"name": "WAN1",
				},
				map[string]interface{}{
					"id":   "aaa00000-0000-0000-0000-000000000002",
					"name": "WAN2",
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListWans(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "WAN1") {
		t.Errorf("result should contain 'WAN1': %s", result)
	}
	if !strings.Contains(result, "WAN2") {
		t.Errorf("result should contain 'WAN2': %s", result)
	}
}

func TestListWans_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListWans(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No WANs found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListWans_Execute_NoSiteID(t *testing.T) {
	tool := &ListWans{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID provided")
	}
	if !strings.Contains(err.Error(), "siteId") {
		t.Errorf("error should mention siteId: %v", err)
	}
}

// --- list_vpn_tunnels tests ---

func TestListVpnTunnels_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":       "bbb00000-0000-0000-0000-000000000001",
					"name":     "Office-to-DC",
					"type":     "IPSEC",
					"metadata": map[string]string{"origin": "USER_DEFINED"},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListVpnTunnels(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Office-to-DC") {
		t.Errorf(
			"result should contain 'Office-to-DC': %s",
			result,
		)
	}
	if !strings.Contains(result, "IPSEC") {
		t.Errorf("result should contain 'IPSEC': %s", result)
	}
}

func TestListVpnTunnels_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListVpnTunnels(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No VPN tunnels found." {
		t.Errorf("unexpected result: %s", result)
	}
}

// --- list_vpn_servers tests ---

func TestListVpnServers_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":       "ccc00000-0000-0000-0000-000000000001",
					"name":     "WireGuard Server",
					"type":     "WIREGUARD",
					"enabled":  true,
					"metadata": map[string]string{"origin": "USER_DEFINED"},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListVpnServers(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "WireGuard Server") {
		t.Errorf(
			"result should contain 'WireGuard Server': %s",
			result,
		)
	}
	if !strings.Contains(result, "WIREGUARD") {
		t.Errorf(
			"result should contain 'WIREGUARD': %s",
			result,
		)
	}
	if !strings.Contains(result, "Enabled: true") {
		t.Errorf(
			"result should contain 'Enabled: true': %s",
			result,
		)
	}
}

func TestListVpnServers_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListVpnServers(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No VPN servers found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListWans_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListWans(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for API error response")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestListVpnTunnels_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListVpnTunnels(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for API error response")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestListVpnServers_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListVpnServers(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for API error response")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

// --- description tests ---

func TestListWans_Description(t *testing.T) {
	tool := NewListWans(nil, "")
	desc := tool.Description()
	if desc == "" {
		t.Fatal("Description() should not be empty")
	}
	if !strings.Contains(desc, "WAN") {
		t.Errorf("description should mention WAN: %s", desc)
	}
}

// --- invalid JSON tests ---

func TestListWans_Execute_InvalidJSON(t *testing.T) {
	tool := &ListWans{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListVpnTunnels_Execute_InvalidJSON(t *testing.T) {
	tool := &ListVpnTunnels{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListVpnServers_Execute_InvalidJSON(t *testing.T) {
	tool := &ListVpnServers{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// --- network error tests ---

func TestListWans_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewListWans(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list WANs") {
		t.Errorf(
			"error should contain 'failed to list WANs': %v",
			err,
		)
	}
}

func TestListVpnTunnels_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewListVpnTunnels(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list VPN tunnels") {
		t.Errorf(
			"error should contain 'failed to list VPN tunnels': %v",
			err,
		)
	}
}

func TestListVpnServers_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewListVpnServers(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list VPN servers") {
		t.Errorf(
			"error should contain 'failed to list VPN servers': %v",
			err,
		)
	}
}

// --- formatting tests ---

func TestListWans_Execute_Formatting(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":   "aaa00000-0000-0000-0000-000000000001",
					"name": "Primary WAN",
				},
				map[string]interface{}{
					"id":   "aaa00000-0000-0000-0000-000000000002",
					"name": "Failover WAN",
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListWans(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "WANs (2 of 2):") {
		t.Errorf(
			"result should contain header 'WANs (2 of 2):': %s",
			result,
		)
	}
	if !strings.Contains(result, "1. Primary WAN") {
		t.Errorf(
			"result should contain '1. Primary WAN': %s",
			result,
		)
	}
	if !strings.Contains(result, "2. Failover WAN") {
		t.Errorf(
			"result should contain '2. Failover WAN': %s",
			result,
		)
	}
}
