package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestListFirewallZones_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":   "aaa00000-0000-0000-0000-000000000001",
					"name": "LAN",
					"networkIds": []string{
						"bbb00000-0000-0000-0000-000000000001",
					},
					"metadata": map[string]string{
						"origin": "SYSTEM_DEFINED",
					},
				},
				map[string]interface{}{
					"id":         "aaa00000-0000-0000-0000-000000000002",
					"name":       "DMZ",
					"networkIds": []string{},
					"metadata": map[string]string{
						"origin": "USER_DEFINED",
					},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListFirewallZones(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "LAN") {
		t.Errorf("result should contain 'LAN': %s", result)
	}
	if !strings.Contains(result, "DMZ") {
		t.Errorf("result should contain 'DMZ': %s", result)
	}
	if !strings.Contains(result, "SYSTEM_DEFINED") {
		t.Errorf(
			"result should contain 'SYSTEM_DEFINED': %s",
			result,
		)
	}
}

func TestListFirewallZones_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListFirewallZones(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "No firewall zones found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListFirewallZones_Execute_NoSiteID(t *testing.T) {
	tool := &ListFirewallZones{}
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

func TestListFirewallZones_DefaultSiteFallback(t *testing.T) {
	var gotPath string
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListFirewallZones(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotPath, testSiteID) {
		t.Errorf(
			"expected path to contain default site ID, got: %s",
			gotPath,
		)
	}
}

func TestGetFirewallZone_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         "aaa00000-0000-0000-0000-000000000001",
				"name":       "LAN",
				"networkIds": []string{"bbb00000-0000-0000-0000-000000000001"},
				"metadata":   map[string]string{"origin": "SYSTEM_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	tool := NewGetFirewallZone(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallZoneId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "LAN") {
		t.Errorf("result should contain 'LAN': %s", result)
	}
}

func TestGetFirewallZone_Execute_MissingZoneID(t *testing.T) {
	tool := &GetFirewallZone{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for missing zone ID")
	}
	if !strings.Contains(err.Error(), "firewallZoneId") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetFirewallZone_Execute_InvalidUUID(t *testing.T) {
	tool := &GetFirewallZone{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallZoneId": "not-valid"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "firewallZoneId") {
		t.Errorf("error should mention firewallZoneId: %v", err)
	}
}

func TestGetFirewallZone_InputSchema(t *testing.T) {
	tool := &GetFirewallZone{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "firewallZoneId")
}

func TestCreateFirewallZone_Execute(t *testing.T) {
	var gotBody map[string]interface{}
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         "aaa00000-0000-0000-0000-000000000003",
				"name":       "TestZone",
				"networkIds": []string{"bbb00000-0000-0000-0000-000000000001"},
				"metadata":   map[string]string{"origin": "USER_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateFirewallZone(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "TestZone", "networkIds": ["bbb00000-0000-0000-0000-000000000001"]}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "TestZone") {
		t.Errorf("result should contain 'TestZone': %s", result)
	}
	if !strings.Contains(result, "created") {
		t.Errorf("result should mention 'created': %s", result)
	}

	if gotBody["name"] != "TestZone" {
		t.Errorf(
			"request body name = %v, want TestZone",
			gotBody["name"],
		)
	}
}

func TestCreateFirewallZone_Execute_MissingName(t *testing.T) {
	tool := &CreateFirewallZone{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"networkIds": []}`),
	)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCreateFirewallZone_InputSchema(t *testing.T) {
	tool := &CreateFirewallZone{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	if len(required) < 2 {
		t.Error("expected at least 2 required fields")
	}
}

func TestUpdateFirewallZone_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         "aaa00000-0000-0000-0000-000000000001",
				"name":       "UpdatedZone",
				"networkIds": []string{},
				"metadata":   map[string]string{"origin": "USER_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	tool := NewUpdateFirewallZone(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallZoneId": "aaa00000-0000-0000-0000-000000000001", "name": "UpdatedZone", "networkIds": []}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "UpdatedZone") {
		t.Errorf(
			"result should contain 'UpdatedZone': %s",
			result,
		)
	}
	if !strings.Contains(result, "updated") {
		t.Errorf("result should mention 'updated': %s", result)
	}
}

func TestUpdateFirewallZone_Execute_MissingZoneID(t *testing.T) {
	tool := &UpdateFirewallZone{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "test", "networkIds": []}`),
	)
	if err == nil {
		t.Fatal("expected error for missing zone ID")
	}
	if !strings.Contains(err.Error(), "firewallZoneId") {
		t.Errorf("error should mention firewallZoneId: %v", err)
	}
}

func TestDeleteFirewallZone_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	tool := NewDeleteFirewallZone(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallZoneId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "deleted successfully") {
		t.Errorf("result should confirm deletion: %s", result)
	}
}

func TestDeleteFirewallZone_Execute_MissingZoneID(t *testing.T) {
	tool := &DeleteFirewallZone{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for missing zone ID")
	}
	if !strings.Contains(err.Error(), "firewallZoneId") {
		t.Errorf("error should mention firewallZoneId: %v", err)
	}
}

func TestDeleteFirewallZone_InputSchema(t *testing.T) {
	tool := &DeleteFirewallZone{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "firewallZoneId")
}

func TestListFirewallZones_Execute_Formatting(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":   "aaa00000-0000-0000-0000-000000000001",
					"name": "LAN",
					"networkIds": []string{
						"bbb00000-0000-0000-0000-000000000001",
					},
					"metadata": map[string]string{
						"origin": "SYSTEM_DEFINED",
					},
				},
				map[string]interface{}{
					"id":         "aaa00000-0000-0000-0000-000000000002",
					"name":       "DMZ",
					"networkIds": []string{},
					"metadata": map[string]string{
						"origin": "USER_DEFINED",
					},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListFirewallZones(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify 1-based numbering
	if !strings.Contains(result, "1. Name: LAN") {
		t.Errorf(
			"result should contain '1. Name: LAN': %s",
			result,
		)
	}
	if !strings.Contains(result, "2. Name: DMZ") {
		t.Errorf(
			"result should contain '2. Name: DMZ': %s",
			result,
		)
	}

	// Verify blank line separator between zones
	if !strings.Contains(result, "\n\n2. ") {
		t.Errorf(
			"result should have blank line between zones: %s",
			result,
		)
	}

	// Verify network IDs shown for LAN (has IDs)
	if !strings.Contains(
		result,
		"Network IDs: bbb00000-0000-0000-0000-000000000001",
	) {
		t.Errorf(
			"result should show network ID for LAN: %s",
			result,
		)
	}

	// Verify "(none)" shown for DMZ (no IDs)
	if !strings.Contains(result, "Network IDs: (none)") {
		t.Errorf(
			"result should show '(none)' for DMZ: %s",
			result,
		)
	}
}

func TestListFirewallZones_Execute_InvalidJSON(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			t.Fatal("API should not be called for invalid JSON")
		}),
	)
	defer srv.Close()

	tool := NewListFirewallZones(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListFirewallZones_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewListFirewallZones(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to list firewall zones",
	) {
		t.Errorf(
			"error should contain 'failed to list firewall zones': %v",
			err,
		)
	}
}

func TestListFirewallZones_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListFirewallZones(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestGetFirewallZone_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewGetFirewallZone(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallZoneId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestCreateFirewallZone_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewCreateFirewallZone(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "TestZone", "networkIds": ["bbb00000-0000-0000-0000-000000000001"]}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestUpdateFirewallZone_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewUpdateFirewallZone(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallZoneId": "aaa00000-0000-0000-0000-000000000001", "name": "TestZone", "networkIds": []}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestDeleteFirewallZone_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewDeleteFirewallZone(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallZoneId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}
