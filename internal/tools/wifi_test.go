package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

func TestListWiFiBroadcasts_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
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
					{
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
				},
				"count":      2,
				"limit":      25,
				"offset":     0,
				"totalCount": 2,
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

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
	srv := httptest.NewServer(
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

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

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
}

func TestListWiFiBroadcasts_Description(t *testing.T) {
	tool := &ListWiFiBroadcasts{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListWiFiBroadcasts_InputSchema(t *testing.T) {
	tool := &ListWiFiBroadcasts{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf(
			"schema type = %v, want object",
			schema["type"],
		)
	}
}

func TestGetWiFiBroadcast_Execute(t *testing.T) {
	srv := httptest.NewServer(
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

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

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
	tool := &GetWiFiBroadcast{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "not-valid"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestGetWiFiBroadcast_Description(t *testing.T) {
	tool := &GetWiFiBroadcast{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestGetWiFiBroadcast_InputSchema(t *testing.T) {
	tool := &GetWiFiBroadcast{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "wifiBroadcastId" {
			found = true
		}
	}
	if !found {
		t.Error("wifiBroadcastId should be required")
	}
}

func TestCreateWiFiBroadcast_Execute(t *testing.T) {
	srv := httptest.NewServer(
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

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

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

func TestCreateWiFiBroadcast_Description(t *testing.T) {
	tool := &CreateWiFiBroadcast{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestCreateWiFiBroadcast_InputSchema(t *testing.T) {
	tool := &CreateWiFiBroadcast{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	foundName := false
	foundSec := false
	for _, r := range required {
		if r == "name" {
			foundName = true
		}
		if r == "securityConfiguration" {
			foundSec = true
		}
	}
	if !foundName {
		t.Error("name should be required")
	}
	if !foundSec {
		t.Error("securityConfiguration should be required")
	}
}

func TestUpdateWiFiBroadcast_Execute(t *testing.T) {
	srv := httptest.NewServer(
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

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

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
	tool := &UpdateWiFiBroadcast{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "not-valid", "name": "x", "enabled": true, "type": "STANDARD", "securityConfiguration": {"type": "OPEN"}}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestUpdateWiFiBroadcast_Description(t *testing.T) {
	tool := &UpdateWiFiBroadcast{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestUpdateWiFiBroadcast_InputSchema(t *testing.T) {
	tool := &UpdateWiFiBroadcast{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "wifiBroadcastId" {
			found = true
		}
	}
	if !found {
		t.Error("wifiBroadcastId should be required")
	}
}

func TestDeleteWiFiBroadcast_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

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
	tool := &DeleteWiFiBroadcast{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"wifiBroadcastId": "not-valid"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestDeleteWiFiBroadcast_Description(t *testing.T) {
	tool := &DeleteWiFiBroadcast{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestDeleteWiFiBroadcast_InputSchema(t *testing.T) {
	tool := &DeleteWiFiBroadcast{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "wifiBroadcastId" {
			found = true
		}
	}
	if !found {
		t.Error("wifiBroadcastId should be required")
	}
}
