package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListWiFiBroadcasts_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":      "aaa00000-0000-0000-0000-000000000001",
					"name":    "Home WiFi",
					"type":    "STANDARD",
					"enabled": true,
					"metadata": map[string]string{
						"origin": "USER_DEFINED",
					},
					"securityConfiguration": map[string]string{
						"type": "WPA2_WPA3",
					},
				},
				map[string]interface{}{
					"id":      "aaa00000-0000-0000-0000-000000000002",
					"name":    "IoT Network",
					"type":    "IOT_OPTIMIZED",
					"enabled": true,
					"metadata": map[string]string{
						"origin": "USER_DEFINED",
					},
					"securityConfiguration": map[string]string{
						"type": "WPA2",
					},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListWiFiBroadcasts(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Home WiFi") {
		t.Errorf(
			"result should contain 'Home WiFi': %s",
			result,
		)
	}
	if !strings.Contains(result, "WPA2_WPA3") {
		t.Errorf(
			"result should contain 'WPA2_WPA3': %s",
			result,
		)
	}
}

func TestListWiFiBroadcasts_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListWiFiBroadcasts(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No WiFi broadcasts found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListWiFiBroadcasts_Execute_NoSiteID(t *testing.T) {
	tool := &ListWiFiBroadcasts{}
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

func TestGetWiFiBroadcast_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       "aaa00000-0000-0000-0000-000000000001",
				"name":     "Home WiFi",
				"type":     "STANDARD",
				"enabled":  true,
				"hideName": false,
				"metadata": map[string]string{
					"origin": "USER_DEFINED",
				},
				"securityConfiguration": map[string]string{
					"type": "WPA2_WPA3",
				},
				"clientIsolationEnabled":              false,
				"multicastToUnicastConversionEnabled": true,
				"uapsdEnabled":                        true,
				"network": map[string]string{
					"type": "DEFAULT",
				},
			})
		}),
	)
	defer srv.Close()

	tool := NewGetWiFiBroadcast(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Home WiFi") {
		t.Errorf(
			"result should contain 'Home WiFi': %s",
			result,
		)
	}
	if !strings.Contains(result, "WPA2_WPA3") {
		t.Errorf(
			"result should contain 'WPA2_WPA3': %s",
			result,
		)
	}
	if !strings.Contains(result, "U-APSD: true") {
		t.Errorf(
			"result should contain U-APSD setting: %s",
			result,
		)
	}
}

func TestGetWiFiBroadcast_Execute_InvalidUUID(t *testing.T) {
	tool := &GetWiFiBroadcast{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "not-valid"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "wifiBroadcastId") {
		t.Errorf("error should mention wifiBroadcastId: %v", err)
	}
}

func TestGetWiFiBroadcast_InputSchema(t *testing.T) {
	tool := &GetWiFiBroadcast{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "wifiBroadcastId")
}

func TestCreateWiFiBroadcast_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       "aaa00000-0000-0000-0000-000000000001",
				"name":     "Guest WiFi",
				"type":     "STANDARD",
				"enabled":  true,
				"hideName": false,
				"metadata": map[string]string{
					"origin": "USER_DEFINED",
				},
				"securityConfiguration": map[string]string{
					"type": "OPEN",
				},
				"clientIsolationEnabled":              true,
				"multicastToUnicastConversionEnabled": false,
				"uapsdEnabled":                        false,
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateWiFiBroadcast(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Guest WiFi", "enabled": true, "type": "STANDARD", "securityConfiguration": {"type": "OPEN"}}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "WiFi broadcast created") {
		t.Errorf(
			"result should contain 'WiFi broadcast created': %s",
			result,
		)
	}
	if !strings.Contains(result, "Guest WiFi") {
		t.Errorf(
			"result should contain 'Guest WiFi': %s",
			result,
		)
	}
}

func TestCreateWiFiBroadcast_InputSchema(t *testing.T) {
	tool := &CreateWiFiBroadcast{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "name")
	requireContains(t, required, "securityConfiguration")
}

func TestUpdateWiFiBroadcast_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       "aaa00000-0000-0000-0000-000000000001",
				"name":     "Updated WiFi",
				"type":     "STANDARD",
				"enabled":  false,
				"hideName": true,
				"metadata": map[string]string{
					"origin": "USER_DEFINED",
				},
				"securityConfiguration": map[string]string{
					"type": "WPA3",
				},
				"clientIsolationEnabled":              false,
				"multicastToUnicastConversionEnabled": false,
				"uapsdEnabled":                        false,
			})
		}),
	)
	defer srv.Close()

	tool := NewUpdateWiFiBroadcast(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "aaa00000-0000-0000-0000-000000000001", "name": "Updated WiFi", "enabled": false, "type": "STANDARD", "securityConfiguration": {"type": "WPA3"}}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "WiFi broadcast updated") {
		t.Errorf(
			"result should contain 'WiFi broadcast updated': %s",
			result,
		)
	}
}

func TestUpdateWiFiBroadcast_Execute_InvalidUUID(t *testing.T) {
	tool := &UpdateWiFiBroadcast{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "not-valid", "name": "x", "enabled": true, "type": "STANDARD", "securityConfiguration": {"type": "OPEN"}}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "wifiBroadcastId") {
		t.Errorf("error should mention wifiBroadcastId: %v", err)
	}
}

func TestUpdateWiFiBroadcast_InputSchema(t *testing.T) {
	tool := &UpdateWiFiBroadcast{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "wifiBroadcastId")
}

func TestDeleteWiFiBroadcast_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	tool := NewDeleteWiFiBroadcast(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "WiFi broadcast deleted." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestDeleteWiFiBroadcast_Execute_InvalidUUID(t *testing.T) {
	tool := &DeleteWiFiBroadcast{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "not-valid"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "wifiBroadcastId") {
		t.Errorf("error should mention wifiBroadcastId: %v", err)
	}
}

func TestDeleteWiFiBroadcast_InputSchema(t *testing.T) {
	tool := &DeleteWiFiBroadcast{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "wifiBroadcastId")
}

func TestGetWiFiBroadcast_Execute_WithOptionalFields(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       "aaa00000-0000-0000-0000-000000000001",
				"name":     "Enterprise WiFi",
				"type":     "STANDARD",
				"enabled":  true,
				"hideName": false,
				"metadata": map[string]string{
					"origin": "USER_DEFINED",
				},
				"securityConfiguration": map[string]interface{}{
					"type": "WPA2",
					"radiusConfiguration": map[string]interface{}{
						"authServer": "radius.local",
						"authPort":   1812,
					},
				},
				"clientIsolationEnabled":              false,
				"multicastToUnicastConversionEnabled": false,
				"uapsdEnabled":                        false,
				"broadcastingDeviceFilter": map[string]interface{}{
					"type": "ALL",
				},
				"clientFilteringPolicy": map[string]interface{}{
					"action": "ALLOW",
					"macAddressFilter": []string{
						"AA:BB:CC:DD:EE:FF",
					},
				},
			})
		}),
	)
	defer srv.Close()

	tool := NewGetWiFiBroadcast(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checks := []string{
		"RADIUS Configuration",
		"Broadcasting Device Filter Type:",
		"Client Filtering:",
	}
	for _, s := range checks {
		if !strings.Contains(result, s) {
			t.Errorf(
				"result should contain %q: %s",
				s,
				result,
			)
		}
	}
}

func TestListWiFiBroadcasts_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListWiFiBroadcasts(client, testSiteID)
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

func TestGetWiFiBroadcast_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewGetWiFiBroadcast(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestCreateWiFiBroadcast_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewCreateWiFiBroadcast(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Test", "enabled": true, "type": "STANDARD", "securityConfiguration": {"type": "OPEN"}}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestUpdateWiFiBroadcast_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewUpdateWiFiBroadcast(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "aaa00000-0000-0000-0000-000000000001", "name": "Test", "enabled": true, "type": "STANDARD", "securityConfiguration": {"type": "OPEN"}}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestDeleteWiFiBroadcast_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewDeleteWiFiBroadcast(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}
