package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListPendingDevices_Execute(t *testing.T) {
	client, srv := testClient(t,
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
