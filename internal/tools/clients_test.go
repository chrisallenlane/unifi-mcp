package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListClients_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":          "aaa00000-0000-0000-0000-000000000001",
					"name":        "Chris's MacBook",
					"type":        "WIRELESS",
					"ipAddress":   "192.168.1.100",
					"connectedAt": "2026-03-01T10:00:00Z",
					"access":      map[string]interface{}{},
				},
				map[string]interface{}{
					"id":        "aaa00000-0000-0000-0000-000000000002",
					"name":      "NAS",
					"type":      "WIRED",
					"ipAddress": "192.168.1.50",
					"access":    map[string]interface{}{},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListClients(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Chris's MacBook") {
		t.Errorf(
			"result should contain client name: %s",
			result,
		)
	}
	if !strings.Contains(result, "WIRELESS") {
		t.Errorf(
			"result should contain 'WIRELESS': %s",
			result,
		)
	}
	if !strings.Contains(result, "192.168.1.100") {
		t.Errorf(
			"result should contain IP address: %s",
			result,
		)
	}
}

func TestListClients_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListClients(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No connected clients found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListClients_Execute_NoSiteID(t *testing.T) {
	tool := &ListClients{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID provided")
	}
}

func TestGetClient_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          "aaa00000-0000-0000-0000-000000000001",
				"name":        "Chris's MacBook",
				"type":        "WIRELESS",
				"ipAddress":   "192.168.1.100",
				"connectedAt": "2026-03-01T10:00:00Z",
				"access":      map[string]interface{}{"type": "STANDARD"},
			})
		}),
	)
	defer srv.Close()

	tool := NewGetClient(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"clientId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Chris's MacBook") {
		t.Errorf(
			"result should contain client name: %s",
			result,
		)
	}
	if !strings.Contains(result, "WIRELESS") {
		t.Errorf(
			"result should contain 'WIRELESS': %s",
			result,
		)
	}
	if !strings.Contains(result, "192.168.1.100") {
		t.Errorf(
			"result should contain IP address: %s",
			result,
		)
	}
}

func TestGetClient_Execute_InvalidUUID(t *testing.T) {
	tool := &GetClient{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"clientId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestGetClient_InputSchema(t *testing.T) {
	tool := &GetClient{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "clientId" {
			found = true
		}
	}
	if !found {
		t.Error("clientId should be required")
	}
}

func TestExecuteClientAction_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"action": "AUTHORIZE_GUEST_ACCESS",
			})
		}),
	)
	defer srv.Close()

	tool := NewExecuteClientAction(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"clientId": "aaa00000-0000-0000-0000-000000000001", "action": "AUTHORIZE_GUEST_ACCESS"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "AUTHORIZE_GUEST_ACCESS") {
		t.Errorf(
			"result should contain action: %s",
			result,
		)
	}
}

func TestExecuteClientAction_Execute_MissingAction(t *testing.T) {
	tool := &ExecuteClientAction{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"clientId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error when action missing")
	}
}

func TestExecuteClientAction_InputSchema(t *testing.T) {
	tool := &ExecuteClientAction{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	foundClient := false
	foundAction := false
	for _, r := range required {
		if r == "clientId" {
			foundClient = true
		}
		if r == "action" {
			foundAction = true
		}
	}
	if !foundClient {
		t.Error("clientId should be required")
	}
	if !foundAction {
		t.Error("action should be required")
	}
}

// --- optional field branches ---

func TestListClients_Execute_WithOptionalFields(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":          "aaa00000-0000-0000-0000-000000000001",
					"name":        "Chris's Laptop",
					"type":        "WIRED",
					"ipAddress":   "192.168.1.100",
					"connectedAt": "2025-01-01T00:00:00Z",
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListClients(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "IP:") {
		t.Errorf(
			"result should contain 'IP:': %s",
			result,
		)
	}
	if !strings.Contains(result, "Connected At:") {
		t.Errorf(
			"result should contain 'Connected At:': %s",
			result,
		)
	}
	if !strings.Contains(result, "192.168.1.100") {
		t.Errorf(
			"result should contain IP address: %s",
			result,
		)
	}
}

// --- API error tests ---

func TestListClients_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewListClients(client, testSiteID)
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

func TestGetClient_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewGetClient(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"clientId": "aaa00000-0000-0000-0000-000000000001"}`,
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

func TestExecuteClientAction_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewExecuteClientAction(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"clientId": "aaa00000-0000-0000-0000-000000000001", "action": "AUTHORIZE_GUEST_ACCESS"}`,
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
