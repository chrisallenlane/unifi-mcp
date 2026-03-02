package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListDNSPolicies_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":       "aaa00000-0000-0000-0000-000000000001",
					"type":     "A_RECORD",
					"domain":   "nas.local",
					"enabled":  true,
					"metadata": map[string]string{"origin": "USER_DEFINED"},
				},
				map[string]interface{}{
					"id":       "aaa00000-0000-0000-0000-000000000002",
					"type":     "CNAME_RECORD",
					"domain":   "wiki.local",
					"enabled":  false,
					"metadata": map[string]string{"origin": "USER_DEFINED"},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListDNSPolicies(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "A_RECORD") {
		t.Errorf(
			"result should contain 'A_RECORD': %s",
			result,
		)
	}
	if !strings.Contains(result, "nas.local") {
		t.Errorf(
			"result should contain 'nas.local': %s",
			result,
		)
	}
}

func TestListDNSPolicies_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListDNSPolicies(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No DNS policies found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListDNSPolicies_Execute_NoSiteID(t *testing.T) {
	tool := &ListDNSPolicies{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID provided")
	}
}

func TestGetDNSPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       "aaa00000-0000-0000-0000-000000000001",
				"type":     "A_RECORD",
				"domain":   "nas.local",
				"enabled":  true,
				"metadata": map[string]string{"origin": "USER_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	tool := NewGetDNSPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "A_RECORD") {
		t.Errorf(
			"result should contain 'A_RECORD': %s",
			result,
		)
	}
	if !strings.Contains(result, "nas.local") {
		t.Errorf(
			"result should contain 'nas.local': %s",
			result,
		)
	}
}

func TestGetDNSPolicy_Execute_InvalidUUID(t *testing.T) {
	tool := &GetDNSPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"dnsPolicyId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestGetDNSPolicy_InputSchema(t *testing.T) {
	tool := &GetDNSPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "dnsPolicyId" {
			found = true
		}
	}
	if !found {
		t.Error("dnsPolicyId should be required")
	}
}

func TestCreateDNSPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       "aaa00000-0000-0000-0000-000000000001",
				"type":     "A_RECORD",
				"domain":   "nas.local",
				"enabled":  true,
				"metadata": map[string]string{"origin": "USER_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateDNSPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"type": "A_RECORD", "enabled": true, "domain": "nas.local", "address": "192.168.1.50"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "DNS policy created") {
		t.Errorf(
			"result should contain 'DNS policy created': %s",
			result,
		)
	}
	if !strings.Contains(result, "A_RECORD") {
		t.Errorf(
			"result should contain 'A_RECORD': %s",
			result,
		)
	}
}

func TestCreateDNSPolicy_InputSchema(t *testing.T) {
	tool := &CreateDNSPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	foundType := false
	foundEnabled := false
	for _, r := range required {
		if r == "type" {
			foundType = true
		}
		if r == "enabled" {
			foundEnabled = true
		}
	}
	if !foundType {
		t.Error("type should be required")
	}
	if !foundEnabled {
		t.Error("enabled should be required")
	}
}

func TestUpdateDNSPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       "aaa00000-0000-0000-0000-000000000001",
				"type":     "A_RECORD",
				"domain":   "nas.local",
				"enabled":  false,
				"metadata": map[string]string{"origin": "USER_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	tool := NewUpdateDNSPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "aaa00000-0000-0000-0000-000000000001", "type": "A_RECORD", "enabled": false}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "DNS policy updated") {
		t.Errorf(
			"result should contain 'DNS policy updated': %s",
			result,
		)
	}
}

func TestUpdateDNSPolicy_Execute_InvalidUUID(t *testing.T) {
	tool := &UpdateDNSPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "not-valid", "type": "A_RECORD", "enabled": true}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestUpdateDNSPolicy_InputSchema(t *testing.T) {
	tool := &UpdateDNSPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "dnsPolicyId" {
			found = true
		}
	}
	if !found {
		t.Error("dnsPolicyId should be required")
	}
}

func TestDeleteDNSPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	tool := NewDeleteDNSPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "DNS policy deleted." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestDeleteDNSPolicy_Execute_InvalidUUID(t *testing.T) {
	tool := &DeleteDNSPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"dnsPolicyId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestDeleteDNSPolicy_InputSchema(t *testing.T) {
	tool := &DeleteDNSPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "dnsPolicyId" {
			found = true
		}
	}
	if !found {
		t.Error("dnsPolicyId should be required")
	}
}
