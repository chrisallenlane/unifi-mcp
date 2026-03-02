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

// --- list_devices ---

func TestListDevices_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":                "bbb00000-0000-0000-0000-000000000001",
						"name":              "US-8-60W",
						"model":             "US860W",
						"macAddress":        "aa:bb:cc:dd:ee:ff",
						"ipAddress":         "192.168.1.2",
						"state":             "ONLINE",
						"firmwareVersion":   "6.6.57",
						"firmwareUpdatable": false,
						"supported":         true,
						"features":          []string{"switching"},
						"interfaces":        []string{"ports"},
					},
					{
						"id":                "bbb00000-0000-0000-0000-000000000002",
						"name":              "UAP-AC-PRO",
						"model":             "UAP-AC-PRO",
						"macAddress":        "11:22:33:44:55:66",
						"ipAddress":         "192.168.1.3",
						"state":             "ONLINE",
						"firmwareUpdatable": true,
						"supported":         true,
						"features":          []string{"accessPoint"},
						"interfaces":        []string{"radios"},
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

	tool := NewListDevices(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "US-8-60W") {
		t.Errorf("result should contain device name: %s", result)
	}
	if !strings.Contains(result, "aa:bb:cc:dd:ee:ff") {
		t.Errorf("result should contain MAC address: %s", result)
	}
	if !strings.Contains(result, "192.168.1.2") {
		t.Errorf("result should contain IP address: %s", result)
	}
	if !strings.Contains(result, "ONLINE") {
		t.Errorf("result should contain state: %s", result)
	}
}

func TestListDevices_Execute_Empty(t *testing.T) {
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

	tool := NewListDevices(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No adopted devices found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListDevices_Execute_NoSiteID(t *testing.T) {
	tool := &ListDevices{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID provided")
	}
}

func TestListDevices_Description(t *testing.T) {
	tool := &ListDevices{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListDevices_InputSchema(t *testing.T) {
	tool := &ListDevices{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}

// --- get_device ---

func TestGetDevice_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":                "bbb00000-0000-0000-0000-000000000001",
				"name":              "US-8-60W",
				"model":             "US860W",
				"macAddress":        "aa:bb:cc:dd:ee:ff",
				"ipAddress":         "192.168.1.2",
				"state":             "ONLINE",
				"firmwareVersion":   "6.6.57",
				"firmwareUpdatable": false,
				"supported":         true,
				"configurationId":   "cfg-001",
				"features":          map[string]interface{}{},
				"interfaces":        map[string]interface{}{},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewGetDevice(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "US-8-60W") {
		t.Errorf("result should contain device name: %s", result)
	}
	if !strings.Contains(result, "aa:bb:cc:dd:ee:ff") {
		t.Errorf("result should contain MAC address: %s", result)
	}
	if !strings.Contains(result, "6.6.57") {
		t.Errorf("result should contain firmware version: %s", result)
	}
}

func TestGetDevice_Execute_InvalidUUID(t *testing.T) {
	tool := &GetDevice{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"deviceId": "not-a-uuid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestGetDevice_Description(t *testing.T) {
	tool := &GetDevice{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestGetDevice_InputSchema(t *testing.T) {
	tool := &GetDevice{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "deviceId" {
			found = true
		}
	}
	if !found {
		t.Error("deviceId should be required")
	}
}

// --- adopt_device ---

func TestAdoptDevice_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":                "bbb00000-0000-0000-0000-000000000001",
				"name":              "Newly Adopted",
				"model":             "US860W",
				"macAddress":        "aa:bb:cc:dd:ee:ff",
				"ipAddress":         "192.168.1.5",
				"state":             "ADOPTING",
				"firmwareUpdatable": false,
				"supported":         true,
				"configurationId":   "cfg-002",
				"features":          map[string]interface{}{},
				"interfaces":        map[string]interface{}{},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewAdoptDevice(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"macAddress": "aa:bb:cc:dd:ee:ff"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Newly Adopted") {
		t.Errorf("result should contain device name: %s", result)
	}
	if !strings.Contains(result, "ADOPTING") {
		t.Errorf("result should contain state: %s", result)
	}
}

func TestAdoptDevice_Execute_MissingMac(t *testing.T) {
	tool := &AdoptDevice{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when macAddress missing")
	}
}

func TestAdoptDevice_Description(t *testing.T) {
	tool := &AdoptDevice{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestAdoptDevice_InputSchema(t *testing.T) {
	tool := &AdoptDevice{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "macAddress" {
			found = true
		}
	}
	if !found {
		t.Error("macAddress should be required")
	}
}

// --- remove_device ---

func TestRemoveDevice_Execute(t *testing.T) {
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

	tool := NewRemoveDevice(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "removed successfully") {
		t.Errorf("result should indicate success: %s", result)
	}
}

func TestRemoveDevice_Execute_InvalidUUID(t *testing.T) {
	tool := &RemoveDevice{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"deviceId": "bad-uuid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestRemoveDevice_Description(t *testing.T) {
	tool := &RemoveDevice{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestRemoveDevice_InputSchema(t *testing.T) {
	tool := &RemoveDevice{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "deviceId" {
			found = true
		}
	}
	if !found {
		t.Error("deviceId should be required")
	}
}

// --- execute_device_action ---

func TestExecuteDeviceAction_Execute(t *testing.T) {
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

	tool := NewExecuteDeviceAction(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001", "action": "restart"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "restart") {
		t.Errorf("result should contain action name: %s", result)
	}
}

func TestExecuteDeviceAction_Execute_MissingAction(t *testing.T) {
	tool := &ExecuteDeviceAction{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error when action missing")
	}
}

func TestExecuteDeviceAction_Execute_InvalidUUID(t *testing.T) {
	tool := &ExecuteDeviceAction{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "not-a-uuid", "action": "restart"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestExecuteDeviceAction_Description(t *testing.T) {
	tool := &ExecuteDeviceAction{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestExecuteDeviceAction_InputSchema(t *testing.T) {
	tool := &ExecuteDeviceAction{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	foundDevice := false
	foundAction := false
	for _, r := range required {
		if r == "deviceId" {
			foundDevice = true
		}
		if r == "action" {
			foundAction = true
		}
	}
	if !foundDevice {
		t.Error("deviceId should be required")
	}
	if !foundAction {
		t.Error("action should be required")
	}
}

// --- execute_port_action ---

func TestExecutePortAction_Execute(t *testing.T) {
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

	tool := NewExecutePortAction(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001", "portIdx": 2, "action": "cycle_poe"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "cycle_poe") {
		t.Errorf("result should contain action name: %s", result)
	}
	if !strings.Contains(result, "port 2") {
		t.Errorf("result should contain port index: %s", result)
	}
}

func TestExecutePortAction_Execute_MissingPortIdx(t *testing.T) {
	tool := &ExecutePortAction{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001", "action": "cycle_poe"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error when portIdx missing")
	}
}

func TestExecutePortAction_Execute_MissingAction(t *testing.T) {
	tool := &ExecutePortAction{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001", "portIdx": 2}`,
		),
	)
	if err == nil {
		t.Fatal("expected error when action missing")
	}
}

func TestExecutePortAction_Execute_InvalidUUID(t *testing.T) {
	tool := &ExecutePortAction{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bad-uuid", "portIdx": 1, "action": "cycle_poe"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestExecutePortAction_Description(t *testing.T) {
	tool := &ExecutePortAction{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestExecutePortAction_InputSchema(t *testing.T) {
	tool := &ExecutePortAction{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	foundDevice := false
	foundPort := false
	foundAction := false
	for _, r := range required {
		switch r {
		case "deviceId":
			foundDevice = true
		case "portIdx":
			foundPort = true
		case "action":
			foundAction = true
		}
	}
	if !foundDevice {
		t.Error("deviceId should be required")
	}
	if !foundPort {
		t.Error("portIdx should be required")
	}
	if !foundAction {
		t.Error("action should be required")
	}
}

// --- get_device_statistics ---

func TestGetDeviceStatistics_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"cpuUtilizationPct":    12.5,
				"memoryUtilizationPct": 45.0,
				"uptimeSec":            86400,
				"loadAverage1Min":      0.5,
				"loadAverage5Min":      0.4,
				"loadAverage15Min":     0.3,
				"lastHeartbeatAt":      "2026-03-01T10:00:00Z",
				"interfaces":           map[string]interface{}{},
				"uplink": map[string]interface{}{
					"rxRateBps": 1000000,
					"txRateBps": 500000,
				},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewGetDeviceStatistics(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "12.5") {
		t.Errorf("result should contain CPU utilization: %s", result)
	}
	if !strings.Contains(result, "86400") {
		t.Errorf("result should contain uptime: %s", result)
	}
}

func TestGetDeviceStatistics_Execute_InvalidUUID(t *testing.T) {
	tool := &GetDeviceStatistics{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"deviceId": "not-a-uuid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestGetDeviceStatistics_Description(t *testing.T) {
	tool := &GetDeviceStatistics{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestGetDeviceStatistics_InputSchema(t *testing.T) {
	tool := &GetDeviceStatistics{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "deviceId" {
			found = true
		}
	}
	if !found {
		t.Error("deviceId should be required")
	}
}

// --- list_pending_devices ---

func TestListPendingDevices_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"model":                 "USW-8",
						"macAddress":            "ff:ee:dd:cc:bb:aa",
						"ipAddress":             "10.0.0.50",
						"state":                 "PENDING_ADOPTION",
						"firmwareVersion":       "6.5.0",
						"firmwareUpdatable":     false,
						"supported":             true,
						"adoptionTargetSiteIds": []string{},
						"features":              []string{"switching"},
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

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewListPendingDevices(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "USW-8") {
		t.Errorf("result should contain model: %s", result)
	}
	if !strings.Contains(result, "ff:ee:dd:cc:bb:aa") {
		t.Errorf("result should contain MAC address: %s", result)
	}
	if !strings.Contains(result, "PENDING_ADOPTION") {
		t.Errorf("result should contain state: %s", result)
	}
}

func TestListPendingDevices_Execute_Empty(t *testing.T) {
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

	tool := NewListPendingDevices(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No pending devices found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListPendingDevices_Description(t *testing.T) {
	tool := &ListPendingDevices{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListPendingDevices_InputSchema(t *testing.T) {
	tool := &ListPendingDevices{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}
