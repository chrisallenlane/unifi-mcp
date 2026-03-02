package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisallenlane/unifi-mcp-server/internal/unifi"
)

func TestListACLRules_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":      "aaa00000-0000-0000-0000-000000000001",
						"name":    "Block IoT to LAN",
						"type":    "IPV4",
						"action":  "BLOCK",
						"enabled": true,
						"index":   0,
						"metadata": map[string]string{
							"origin": "USER_DEFINED",
						},
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

	tool := NewListACLRules(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Block IoT to LAN") {
		t.Errorf(
			"result should contain rule name: %s",
			result,
		)
	}
	if !strings.Contains(result, "IPV4") {
		t.Errorf(
			"result should contain 'IPV4': %s",
			result,
		)
	}
}

func TestListACLRules_Execute_Empty(t *testing.T) {
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

	tool := NewListACLRules(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No ACL rules found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListACLRules_Execute_NoSiteID(t *testing.T) {
	tool := &ListACLRules{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID provided")
	}
}

func TestListACLRules_Description(t *testing.T) {
	tool := &ListACLRules{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListACLRules_InputSchema(t *testing.T) {
	tool := &ListACLRules{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf(
			"schema type = %v, want object",
			schema["type"],
		)
	}
}

func TestGetACLRule_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":      "aaa00000-0000-0000-0000-000000000001",
				"name":    "Block IoT to LAN",
				"type":    "IPV4",
				"action":  "BLOCK",
				"enabled": true,
				"index":   0,
				"metadata": map[string]string{
					"origin": "USER_DEFINED",
				},
				"sourceFilter": map[string]interface{}{
					"type":      "NETWORK",
					"networkId": "bbb00000-0000-0000-0000-000000000001",
				},
				"description": "Block IoT devices from accessing LAN",
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewGetACLRule(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"aclRuleId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Block IoT to LAN") {
		t.Errorf(
			"result should contain rule name: %s",
			result,
		)
	}
	if !strings.Contains(result, "Source Filter") {
		t.Errorf(
			"result should contain source filter: %s",
			result,
		)
	}
	if !strings.Contains(
		result,
		"Block IoT devices from accessing LAN",
	) {
		t.Errorf(
			"result should contain description: %s",
			result,
		)
	}
}

func TestGetACLRule_Execute_InvalidUUID(t *testing.T) {
	tool := &GetACLRule{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"aclRuleId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestGetACLRule_Description(t *testing.T) {
	tool := &GetACLRule{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestGetACLRule_InputSchema(t *testing.T) {
	tool := &GetACLRule{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "aclRuleId" {
			found = true
		}
	}
	if !found {
		t.Error("aclRuleId should be required")
	}
}

func TestCreateACLRule_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":      "aaa00000-0000-0000-0000-000000000001",
				"name":    "Allow LAN to WAN",
				"type":    "IPV4",
				"action":  "ALLOW",
				"enabled": true,
				"index":   1,
				"metadata": map[string]string{
					"origin": "USER_DEFINED",
				},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewCreateACLRule(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"type": "IPV4", "name": "Allow LAN to WAN", "enabled": true, "action": "ALLOW"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "ACL rule created") {
		t.Errorf(
			"result should contain 'ACL rule created': %s",
			result,
		)
	}
	if !strings.Contains(result, "Allow LAN to WAN") {
		t.Errorf(
			"result should contain rule name: %s",
			result,
		)
	}
}

func TestCreateACLRule_Description(t *testing.T) {
	tool := &CreateACLRule{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestCreateACLRule_InputSchema(t *testing.T) {
	tool := &CreateACLRule{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	foundType := false
	foundName := false
	for _, r := range required {
		if r == "type" {
			foundType = true
		}
		if r == "name" {
			foundName = true
		}
	}
	if !foundType {
		t.Error("type should be required")
	}
	if !foundName {
		t.Error("name should be required")
	}
}

func TestUpdateACLRule_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":      "aaa00000-0000-0000-0000-000000000001",
				"name":    "Updated Rule",
				"type":    "IPV4",
				"action":  "BLOCK",
				"enabled": false,
				"index":   0,
				"metadata": map[string]string{
					"origin": "USER_DEFINED",
				},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewUpdateACLRule(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"aclRuleId": "aaa00000-0000-0000-0000-000000000001", "type": "IPV4", "name": "Updated Rule", "enabled": false, "action": "BLOCK"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "ACL rule updated") {
		t.Errorf(
			"result should contain 'ACL rule updated': %s",
			result,
		)
	}
}

func TestUpdateACLRule_Execute_InvalidUUID(t *testing.T) {
	tool := &UpdateACLRule{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"aclRuleId": "not-valid", "type": "IPV4", "name": "x", "enabled": true, "action": "ALLOW"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestUpdateACLRule_Description(t *testing.T) {
	tool := &UpdateACLRule{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestUpdateACLRule_InputSchema(t *testing.T) {
	tool := &UpdateACLRule{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "aclRuleId" {
			found = true
		}
	}
	if !found {
		t.Error("aclRuleId should be required")
	}
}

func TestDeleteACLRule_Execute(t *testing.T) {
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

	tool := NewDeleteACLRule(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"aclRuleId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "ACL rule deleted." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestDeleteACLRule_Execute_InvalidUUID(t *testing.T) {
	tool := &DeleteACLRule{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"aclRuleId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestDeleteACLRule_Description(t *testing.T) {
	tool := &DeleteACLRule{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestDeleteACLRule_InputSchema(t *testing.T) {
	tool := &DeleteACLRule{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "aclRuleId" {
			found = true
		}
	}
	if !found {
		t.Error("aclRuleId should be required")
	}
}

func TestGetACLRuleOrdering_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"orderedAclRuleIds": []string{
					"aaa00000-0000-0000-0000-000000000001",
					"aaa00000-0000-0000-0000-000000000002",
				},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewGetACLRuleOrdering(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "2 rules") {
		t.Errorf(
			"result should contain rule count: %s",
			result,
		)
	}
	if !strings.Contains(
		result,
		"aaa00000-0000-0000-0000-000000000001",
	) {
		t.Errorf(
			"result should contain first rule ID: %s",
			result,
		)
	}
}

func TestGetACLRuleOrdering_Execute_Empty(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"orderedAclRuleIds": []interface{}{},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewGetACLRuleOrdering(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No ACL rules in ordering." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestGetACLRuleOrdering_Description(t *testing.T) {
	tool := &GetACLRuleOrdering{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestGetACLRuleOrdering_InputSchema(t *testing.T) {
	tool := &GetACLRuleOrdering{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf(
			"schema type = %v, want object",
			schema["type"],
		)
	}
}

func TestUpdateACLRuleOrdering_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"orderedAclRuleIds": []string{
					"aaa00000-0000-0000-0000-000000000002",
					"aaa00000-0000-0000-0000-000000000001",
				},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewUpdateACLRuleOrdering(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"orderedAclRuleIds": ["aaa00000-0000-0000-0000-000000000002", "aaa00000-0000-0000-0000-000000000001"]}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "2 rules") {
		t.Errorf(
			"result should contain rule count: %s",
			result,
		)
	}
}

func TestUpdateACLRuleOrdering_Description(t *testing.T) {
	tool := &UpdateACLRuleOrdering{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestUpdateACLRuleOrdering_InputSchema(t *testing.T) {
	tool := &UpdateACLRuleOrdering{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "orderedAclRuleIds" {
			found = true
		}
	}
	if !found {
		t.Error("orderedAclRuleIds should be required")
	}
}
