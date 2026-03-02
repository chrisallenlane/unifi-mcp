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
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":   "aaa00000-0000-0000-0000-000000000001",
						"name": "WAN1",
					},
					{
						"id":   "aaa00000-0000-0000-0000-000000000002",
						"name": "WAN2",
					},
				},
				"count":      2,
				"limit":      25,
				"offset":     0,
				"totalCount": 2,
			})
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
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data":       []interface{}{},
				"count":      0,
				"limit":      25,
				"offset":     0,
				"totalCount": 0,
			})
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
}

func TestListWans_Description(t *testing.T) {
	tool := &ListWans{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListWans_InputSchema(t *testing.T) {
	tool := &ListWans{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}

// --- list_vpn_tunnels tests ---

func TestListVpnTunnels_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":       "bbb00000-0000-0000-0000-000000000001",
						"name":     "Office-to-DC",
						"type":     "IPSEC",
						"metadata": map[string]string{"origin": "USER_DEFINED"},
					},
				},
				"count":      1,
				"limit":      25,
				"offset":     0,
				"totalCount": 1,
			})
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
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data":       []interface{}{},
				"count":      0,
				"limit":      25,
				"offset":     0,
				"totalCount": 0,
			})
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

func TestListVpnTunnels_Description(t *testing.T) {
	tool := &ListVpnTunnels{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListVpnTunnels_InputSchema(t *testing.T) {
	tool := &ListVpnTunnels{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}

// --- list_vpn_servers tests ---

func TestListVpnServers_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":       "ccc00000-0000-0000-0000-000000000001",
						"name":     "WireGuard Server",
						"type":     "WIREGUARD",
						"enabled":  true,
						"metadata": map[string]string{"origin": "USER_DEFINED"},
					},
				},
				"count":      1,
				"limit":      25,
				"offset":     0,
				"totalCount": 1,
			})
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
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data":       []interface{}{},
				"count":      0,
				"limit":      25,
				"offset":     0,
				"totalCount": 0,
			})
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

func TestListVpnServers_Description(t *testing.T) {
	tool := &ListVpnServers{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListVpnServers_InputSchema(t *testing.T) {
	tool := &ListVpnServers{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}
