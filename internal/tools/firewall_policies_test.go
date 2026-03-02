package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func mockPolicyJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":      "ccc00000-0000-0000-0000-000000000001",
		"name":    "Allow LAN to WAN",
		"enabled": true,
		"action":  map[string]string{"type": "ALLOW"},
		"source": map[string]interface{}{
			"zoneId": "aaa00000-0000-0000-0000-000000000001",
		},
		"destination": map[string]interface{}{
			"zoneId": "aaa00000-0000-0000-0000-000000000002",
		},
		"ipProtocolScope": map[string]string{
			"ipVersion": "IPV4_AND_IPV6",
		},
		"loggingEnabled": false,
		"index":          0,
		"metadata": map[string]string{
			"origin": "USER_DEFINED",
		},
	}
}

func TestListFirewallPolicies_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data":       []interface{}{mockPolicyJSON()},
				"count":      1,
				"limit":      25,
				"offset":     0,
				"totalCount": 1,
			})
		}),
	)
	defer srv.Close()

	tool := NewListFirewallPolicies(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Allow LAN to WAN") {
		t.Errorf(
			"result should contain policy name: %s",
			result,
		)
	}
	if !strings.Contains(result, "ALLOW") {
		t.Errorf("result should contain action: %s", result)
	}
}

func TestListFirewallPolicies_Execute_Empty(t *testing.T) {
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

	tool := NewListFirewallPolicies(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "No firewall policies found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListFirewallPolicies_Execute_NoSiteID(t *testing.T) {
	tool := &ListFirewallPolicies{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID")
	}
}

func TestListFirewallPolicies_Description(t *testing.T) {
	tool := &ListFirewallPolicies{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListFirewallPolicies_InputSchema(t *testing.T) {
	tool := &ListFirewallPolicies{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}

func TestGetFirewallPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockPolicyJSON())
		}),
	)
	defer srv.Close()

	tool := NewGetFirewallPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallPolicyId": "ccc00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Allow LAN to WAN") {
		t.Errorf("result should contain policy name: %s", result)
	}
}

func TestGetFirewallPolicy_Execute_MissingPolicyID(t *testing.T) {
	tool := &GetFirewallPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for missing policy ID")
	}
}

func TestGetFirewallPolicy_Description(t *testing.T) {
	tool := &GetFirewallPolicy{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestGetFirewallPolicy_InputSchema(t *testing.T) {
	tool := &GetFirewallPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "firewallPolicyId" {
			found = true
		}
	}
	if !found {
		t.Error("firewallPolicyId should be required")
	}
}

func TestCreateFirewallPolicy_Execute(t *testing.T) {
	var gotBody map[string]interface{}
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(mockPolicyJSON())
		}),
	)
	defer srv.Close()

	tool := NewCreateFirewallPolicy(client, testSiteID)
	args := `{
		"name": "Allow LAN to WAN",
		"enabled": true,
		"action": {"type": "ALLOW"},
		"source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"},
		"destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"},
		"ipProtocolScope": {"ipVersion": "IPV4_AND_IPV6"},
		"loggingEnabled": false
	}`
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "created") {
		t.Errorf("result should mention created: %s", result)
	}
	if !strings.Contains(result, "Allow LAN to WAN") {
		t.Errorf("result should contain policy name: %s", result)
	}

	if gotBody["name"] != "Allow LAN to WAN" {
		t.Errorf(
			"request body name = %v, want Allow LAN to WAN",
			gotBody["name"],
		)
	}
	action, ok := gotBody["action"].(map[string]interface{})
	if !ok || action["type"] != "ALLOW" {
		t.Errorf(
			"request body action.type = %v, want ALLOW",
			gotBody["action"],
		)
	}
}

func TestCreateFirewallPolicy_Execute_WithOptionalFields(
	t *testing.T,
) {
	var gotBody map[string]interface{}
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(mockPolicyJSON())
		}),
	)
	defer srv.Close()

	tool := NewCreateFirewallPolicy(client, testSiteID)
	args := `{
		"name": "Block Policy",
		"enabled": true,
		"action": {"type": "BLOCK"},
		"source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"},
		"destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"},
		"ipProtocolScope": {"ipVersion": "IPV4"},
		"loggingEnabled": true,
		"description": "Test block policy",
		"connectionStateFilter": ["NEW", "ESTABLISHED"],
		"ipsecFilter": "MATCH_NOT_ENCRYPTED"
	}`
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotBody["description"] != "Test block policy" {
		t.Errorf(
			"request body description = %v",
			gotBody["description"],
		)
	}
	csf, ok := gotBody["connectionStateFilter"].([]interface{})
	if !ok || len(csf) != 2 {
		t.Errorf(
			"request body connectionStateFilter = %v",
			gotBody["connectionStateFilter"],
		)
	}
	if gotBody["ipsecFilter"] != "MATCH_NOT_ENCRYPTED" {
		t.Errorf(
			"request body ipsecFilter = %v",
			gotBody["ipsecFilter"],
		)
	}
}

func TestCreateFirewallPolicy_Execute_MissingName(t *testing.T) {
	tool := &CreateFirewallPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"enabled": true, "action": {"type": "ALLOW"}, "source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"}, "destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"}, "ipProtocolScope": {"ipVersion": "IPV4"}, "loggingEnabled": false}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCreateFirewallPolicy_Description(t *testing.T) {
	tool := &CreateFirewallPolicy{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestCreateFirewallPolicy_InputSchema(t *testing.T) {
	tool := &CreateFirewallPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	if len(required) < 7 {
		t.Errorf("expected at least 7 required fields, got %d", len(required))
	}
}

func TestUpdateFirewallPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockPolicyJSON())
		}),
	)
	defer srv.Close()

	tool := NewUpdateFirewallPolicy(client, testSiteID)
	args := `{
		"firewallPolicyId": "ccc00000-0000-0000-0000-000000000001",
		"name": "Allow LAN to WAN",
		"enabled": true,
		"action": {"type": "ALLOW"},
		"source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"},
		"destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"},
		"ipProtocolScope": {"ipVersion": "IPV4_AND_IPV6"},
		"loggingEnabled": false
	}`
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "updated") {
		t.Errorf("result should mention updated: %s", result)
	}
}

func TestUpdateFirewallPolicy_Execute_MissingPolicyID(t *testing.T) {
	tool := &UpdateFirewallPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "test", "enabled": true, "action": {"type": "ALLOW"}, "source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"}, "destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"}, "ipProtocolScope": {"ipVersion": "IPV4"}, "loggingEnabled": false}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for missing policy ID")
	}
}

func TestUpdateFirewallPolicy_Description(t *testing.T) {
	tool := &UpdateFirewallPolicy{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestDeleteFirewallPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	tool := NewDeleteFirewallPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallPolicyId": "ccc00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "deleted successfully") {
		t.Errorf("result should confirm deletion: %s", result)
	}
}

func TestDeleteFirewallPolicy_Execute_MissingPolicyID(
	t *testing.T,
) {
	tool := &DeleteFirewallPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for missing policy ID")
	}
}

func TestDeleteFirewallPolicy_Description(t *testing.T) {
	tool := &DeleteFirewallPolicy{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestDeleteFirewallPolicy_InputSchema(t *testing.T) {
	tool := &DeleteFirewallPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "firewallPolicyId" {
			found = true
		}
	}
	if !found {
		t.Error("firewallPolicyId should be required")
	}
}

func TestPatchFirewallPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			p := mockPolicyJSON()
			p["loggingEnabled"] = true
			json.NewEncoder(w).Encode(p)
		}),
	)
	defer srv.Close()

	tool := NewPatchFirewallPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallPolicyId": "ccc00000-0000-0000-0000-000000000001", "loggingEnabled": true}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "patched") {
		t.Errorf("result should mention patched: %s", result)
	}
	if !strings.Contains(result, "Logging: true") {
		t.Errorf("result should show logging enabled: %s", result)
	}
}

func TestPatchFirewallPolicy_Execute_MissingPolicyID(
	t *testing.T,
) {
	tool := &PatchFirewallPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for missing policy ID")
	}
}

func TestPatchFirewallPolicy_Description(t *testing.T) {
	tool := &PatchFirewallPolicy{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestPatchFirewallPolicy_InputSchema(t *testing.T) {
	tool := &PatchFirewallPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "firewallPolicyId" {
			found = true
		}
	}
	if !found {
		t.Error("firewallPolicyId should be required")
	}
}
