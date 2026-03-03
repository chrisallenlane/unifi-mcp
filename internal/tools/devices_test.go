package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// --- list_devices ---

func TestListDevices_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
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
				map[string]interface{}{
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
			))
		}),
	)
	defer srv.Close()

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
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

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
	if !strings.Contains(err.Error(), "siteId") {
		t.Errorf("error should mention siteId: %v", err)
	}
}

// --- get_device ---

func TestGetDevice_Execute(t *testing.T) {
	client, srv := testClient(t,
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
	tool := &GetDevice{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"deviceId": "not-a-uuid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "deviceId") {
		t.Errorf("error should mention deviceId: %v", err)
	}
}

func TestGetDevice_InputSchema(t *testing.T) {
	tool := &GetDevice{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "deviceId")
}

// --- adopt_device ---

func TestAdoptDevice_Execute(t *testing.T) {
	client, srv := testClient(t,
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
	tool := &AdoptDevice{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when macAddress missing")
	}
	if !strings.Contains(err.Error(), "macAddress") {
		t.Errorf("error should mention macAddress: %v", err)
	}
}

func TestAdoptDevice_InputSchema(t *testing.T) {
	tool := &AdoptDevice{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "macAddress")
}

// --- remove_device ---

func TestRemoveDevice_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

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
	tool := &RemoveDevice{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"deviceId": "bad-uuid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "deviceId") {
		t.Errorf("error should mention deviceId: %v", err)
	}
}

func TestRemoveDevice_InputSchema(t *testing.T) {
	tool := &RemoveDevice{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "deviceId")
}

// --- execute_device_action ---

func TestExecuteDeviceAction_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

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
	tool := &ExecuteDeviceAction{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error when action missing")
	}
	if !strings.Contains(err.Error(), "action") {
		t.Errorf("error should mention action: %v", err)
	}
}

func TestExecuteDeviceAction_Execute_InvalidUUID(t *testing.T) {
	tool := &ExecuteDeviceAction{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "not-a-uuid", "action": "restart"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "deviceId") {
		t.Errorf("error should mention deviceId: %v", err)
	}
}

func TestExecuteDeviceAction_InputSchema(t *testing.T) {
	tool := &ExecuteDeviceAction{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "deviceId")
	requireContains(t, required, "action")
}

// --- execute_port_action ---

func TestExecutePortAction_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

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
	tool := &ExecutePortAction{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001", "action": "cycle_poe"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error when portIdx missing")
	}
	if !strings.Contains(err.Error(), "portIdx") {
		t.Errorf("error should mention portIdx: %v", err)
	}
}

func TestExecutePortAction_Execute_MissingAction(t *testing.T) {
	tool := &ExecutePortAction{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001", "portIdx": 2}`,
		),
	)
	if err == nil {
		t.Fatal("expected error when action missing")
	}
	if !strings.Contains(err.Error(), "action") {
		t.Errorf("error should mention action: %v", err)
	}
}

func TestExecutePortAction_Execute_InvalidUUID(t *testing.T) {
	tool := &ExecutePortAction{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bad-uuid", "portIdx": 1, "action": "cycle_poe"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "deviceId") {
		t.Errorf("error should mention deviceId: %v", err)
	}
}

func TestExecutePortAction_InputSchema(t *testing.T) {
	tool := &ExecutePortAction{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "deviceId")
	requireContains(t, required, "portIdx")
	requireContains(t, required, "action")
}

// --- get_device_statistics ---

func TestGetDeviceStatistics_Execute(t *testing.T) {
	client, srv := testClient(t,
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
	tool := &GetDeviceStatistics{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"deviceId": "not-a-uuid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "deviceId") {
		t.Errorf("error should mention deviceId: %v", err)
	}
}

func TestGetDeviceStatistics_InputSchema(t *testing.T) {
	tool := &GetDeviceStatistics{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "deviceId")
}

// --- optional field branches ---

func TestGetDevice_Execute_WithOptionalFields(t *testing.T) {
	client, srv := testClient(t,
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
				"features":          map[string]interface{}{},
				"interfaces":        map[string]interface{}{},
				"adoptedAt":         "2025-01-01T00:00:00Z",
				"provisionedAt":     "2025-01-02T00:00:00Z",
			})
		}),
	)
	defer srv.Close()

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

	if !strings.Contains(result, "Adopted At:") {
		t.Errorf(
			"result should contain 'Adopted At:': %s",
			result,
		)
	}
	if !strings.Contains(result, "2025-01-01") {
		t.Errorf(
			"result should contain adopted-at date: %s",
			result,
		)
	}
	if !strings.Contains(result, "Provisioned At:") {
		t.Errorf(
			"result should contain 'Provisioned At:': %s",
			result,
		)
	}
	if !strings.Contains(result, "2025-01-02") {
		t.Errorf(
			"result should contain provisioned-at date: %s",
			result,
		)
	}
}

func TestGetDeviceStatistics_Execute_WithNextHeartbeat(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"lastHeartbeatAt": "2026-03-01T10:00:00Z",
				"nextHeartbeatAt": "2026-03-01T10:00:30Z",
				"interfaces":      map[string]interface{}{},
			})
		}),
	)
	defer srv.Close()

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

	if !strings.Contains(result, "Next Heartbeat:") {
		t.Errorf(
			"result should contain 'Next Heartbeat:': %s",
			result,
		)
	}
	if !strings.Contains(result, "2026-03-01") {
		t.Errorf(
			"result should contain next-heartbeat date: %s",
			result,
		)
	}
}

// --- formatting and coverage tests ---

func TestListDevices_Description(t *testing.T) {
	tool := &ListDevices{}
	d := tool.Description()
	if d == "" {
		t.Fatal("description should not be empty")
	}
	if !strings.Contains(d, "device") {
		t.Errorf(
			"description should mention devices: %s",
			d,
		)
	}
}

func TestListDevices_Execute_InvalidJSON(t *testing.T) {
	tool := &ListDevices{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListDevices_Execute_Formatting(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":                "bbb00000-0000-0000-0000-000000000001",
						"name":              "Switch",
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
						"name":              "AP",
						"model":             "UAP",
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
				"totalCount": 10,
			})
		}),
	)
	defer srv.Close()

	tool := NewListDevices(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// verify header uses totalCount
	if !strings.Contains(
		result,
		"Adopted Devices (2 of 10):",
	) {
		t.Errorf(
			"result should contain 'Adopted Devices (2 of 10):': %s",
			result,
		)
	}

	// verify numbering
	if !strings.Contains(result, "1. Switch") {
		t.Errorf(
			"result should contain '1. Switch': %s",
			result,
		)
	}
	if !strings.Contains(result, "2. AP") {
		t.Errorf(
			"result should contain '2. AP': %s",
			result,
		)
	}
}

func TestListDevices_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewListDevices(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to list devices",
	) {
		t.Errorf(
			"error should contain 'failed to list devices': %v",
			err,
		)
	}
}

func TestGetDevice_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewGetDevice(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to get device",
	) {
		t.Errorf(
			"error should contain 'failed to get device': %v",
			err,
		)
	}
}

func TestGetDevice_Execute_FeaturesJSON(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":                "bbb00000-0000-0000-0000-000000000001",
				"name":              "Switch",
				"model":             "US860W",
				"macAddress":        "aa:bb:cc:dd:ee:ff",
				"ipAddress":         "192.168.1.2",
				"state":             "ONLINE",
				"firmwareUpdatable": false,
				"supported":         true,
				"features": map[string]interface{}{
					"switching": true,
				},
				"interfaces": map[string]interface{}{},
			})
		}),
	)
	defer srv.Close()

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

	// verify Features JSON is included
	if !strings.Contains(result, "Features:") {
		t.Errorf(
			"result should contain 'Features:': %s",
			result,
		)
	}
	if !strings.Contains(result, "switching") {
		t.Errorf(
			"result should contain feature key: %s",
			result,
		)
	}
}

func TestGetDeviceStatistics_Execute_NetworkError(
	t *testing.T,
) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewGetDeviceStatistics(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to get device statistics",
	) {
		t.Errorf(
			"error should contain 'failed to get device statistics': %v",
			err,
		)
	}
}

// --- transport error tests ---

func TestAdoptDevice_Execute_TransportError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
	)
	srv.Close()

	tool := NewAdoptDevice(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"macAddress": "aa:bb:cc:dd:ee:ff"}`),
	)
	if err == nil {
		t.Fatal("expected error on transport failure")
	}
	if !strings.Contains(err.Error(), "failed to adopt device") {
		t.Errorf("error should wrap with context: %v", err)
	}
}

func TestRemoveDevice_Execute_TransportError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
	)
	srv.Close()

	tool := NewRemoveDevice(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error on transport failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to remove device",
	) {
		t.Errorf("error should wrap with context: %v", err)
	}
}

func TestExecuteDeviceAction_Execute_TransportError(
	t *testing.T,
) {
	client, srv := testClient(t,
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
	)
	srv.Close()

	tool := NewExecuteDeviceAction(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001", "action": "restart"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error on transport failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to execute device action",
	) {
		t.Errorf("error should wrap with context: %v", err)
	}
}

func TestExecutePortAction_Execute_TransportError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
	)
	srv.Close()

	tool := NewExecutePortAction(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001", "portIdx": 2, "action": "cycle_poe"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error on transport failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to execute port action",
	) {
		t.Errorf("error should wrap with context: %v", err)
	}
}

// --- API error tests ---

func TestListDevices_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewListDevices(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error on API failure")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"error should contain '500': %v",
			err,
		)
	}
}

func TestGetDevice_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewGetDevice(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error on API failure")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"error should contain '500': %v",
			err,
		)
	}
}

func TestAdoptDevice_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewAdoptDevice(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"macAddress": "aa:bb:cc:dd:ee:ff"}`),
	)
	if err == nil {
		t.Fatal("expected error on API failure")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"error should contain '500': %v",
			err,
		)
	}
}

func TestRemoveDevice_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewRemoveDevice(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error on API failure")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"error should contain '500': %v",
			err,
		)
	}
}

func TestExecuteDeviceAction_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewExecuteDeviceAction(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001", "action": "restart"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error on API failure")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"error should contain '500': %v",
			err,
		)
	}
}

func TestExecutePortAction_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewExecutePortAction(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001", "portIdx": 2, "action": "cycle_poe"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error on API failure")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"error should contain '500': %v",
			err,
		)
	}
}

func TestGetDeviceStatistics_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewGetDeviceStatistics(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"deviceId": "bbb00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error on API failure")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"error should contain '500': %v",
			err,
		)
	}
}
