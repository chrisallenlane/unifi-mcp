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
			json.NewEncoder(w).Encode(paginatedResponse(
				mockPolicyJSON(),
			))
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
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
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
	if !strings.Contains(err.Error(), "siteId") {
		t.Errorf("error should mention siteId: %v", err)
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
	if !strings.Contains(err.Error(), "firewallPolicyId") {
		t.Errorf("error should mention firewallPolicyId: %v", err)
	}
}

func TestGetFirewallPolicy_InputSchema(t *testing.T) {
	tool := &GetFirewallPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "firewallPolicyId")
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

func TestCreateFirewallPolicy_Execute_MissingActionType(
	t *testing.T,
) {
	tool := &CreateFirewallPolicy{
		baseTool{defaultSiteID: testSiteID},
	}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Test", "enabled": true, "action": {}, "source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"}, "destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"}, "ipProtocolScope": {"ipVersion": "IPV4"}, "loggingEnabled": false}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for missing action.type")
	}
	if !strings.Contains(err.Error(), "action.type is required") {
		t.Errorf("unexpected error: %v", err)
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
	if !strings.Contains(err.Error(), "firewallPolicyId") {
		t.Errorf("error should mention firewallPolicyId: %v", err)
	}
}

func TestUpdateFirewallPolicy_Execute_MissingName(t *testing.T) {
	tool := &UpdateFirewallPolicy{
		baseTool{defaultSiteID: testSiteID},
	}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallPolicyId": "ccc00000-0000-0000-0000-000000000001", "enabled": true, "action": {"type": "ALLOW"}, "source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"}, "destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"}, "ipProtocolScope": {"ipVersion": "IPV4"}, "loggingEnabled": false}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("unexpected error: %v", err)
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
	if !strings.Contains(err.Error(), "firewallPolicyId") {
		t.Errorf("error should mention firewallPolicyId: %v", err)
	}
}

func TestDeleteFirewallPolicy_InputSchema(t *testing.T) {
	tool := &DeleteFirewallPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "firewallPolicyId")
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
	if !strings.Contains(err.Error(), "firewallPolicyId") {
		t.Errorf("error should mention firewallPolicyId: %v", err)
	}
}

func TestPatchFirewallPolicy_InputSchema(t *testing.T) {
	tool := &PatchFirewallPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "firewallPolicyId")
}

func TestCreateFirewallPolicy_Execute_InvalidSourceZoneID(
	t *testing.T,
) {
	tool := &CreateFirewallPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{
			"name": "Test",
			"enabled": true,
			"action": {"type": "ALLOW"},
			"source": {"zoneId": "not-a-uuid"},
			"destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"},
			"ipProtocolScope": {"ipVersion": "IPV4"},
			"loggingEnabled": false
		}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid source zone UUID")
	}
	if !strings.Contains(err.Error(), "source.zoneId") {
		t.Errorf(
			"error should mention source.zoneId: %v",
			err,
		)
	}
}

func TestCreateFirewallPolicy_Execute_InvalidDestinationZoneID(
	t *testing.T,
) {
	tool := &CreateFirewallPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{
			"name": "Test",
			"enabled": true,
			"action": {"type": "ALLOW"},
			"source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"},
			"destination": {"zoneId": "not-a-uuid"},
			"ipProtocolScope": {"ipVersion": "IPV4"},
			"loggingEnabled": false
		}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid destination zone UUID")
	}
	if !strings.Contains(err.Error(), "destination.zoneId") {
		t.Errorf(
			"error should mention destination.zoneId: %v",
			err,
		)
	}
}

func mockPolicyWithOptionalFieldsJSON() map[string]interface{} {
	desc := "Test policy description"
	return map[string]interface{}{
		"id":      "ccc00000-0000-0000-0000-000000000001",
		"name":    "Allow LAN to WAN",
		"enabled": true,
		"action":  map[string]string{"type": "ALLOW"},
		"source": map[string]interface{}{
			"zoneId": "aaa00000-0000-0000-0000-000000000001",
			"trafficFilter": map[string]interface{}{
				"type": "NETWORK",
			},
		},
		"destination": map[string]interface{}{
			"zoneId": "aaa00000-0000-0000-0000-000000000002",
			"trafficFilter": map[string]interface{}{
				"type": "IP_ADDRESS",
			},
		},
		"ipProtocolScope": map[string]string{
			"ipVersion": "IPV4_AND_IPV6",
		},
		"loggingEnabled": false,
		"index":          0,
		"metadata": map[string]string{
			"origin": "USER_DEFINED",
		},
		"description":           desc,
		"connectionStateFilter": []string{"NEW", "ESTABLISHED"},
		"ipsecFilter":           "MATCH_ENCRYPTED",
		"schedule":              map[string]string{"mode": "ALWAYS"},
	}
}

func TestGetFirewallPolicy_Execute_WithOptionalFields(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(
				mockPolicyWithOptionalFieldsJSON(),
			)
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
	if !strings.Contains(result, "Test policy description") {
		t.Errorf(
			"result should contain description value: %s",
			result,
		)
	}
	if !strings.Contains(result, "NEW") {
		t.Errorf(
			"result should contain connection state filter: %s",
			result,
		)
	}
	if !strings.Contains(result, "MATCH_ENCRYPTED") {
		t.Errorf(
			"result should contain ipsec filter value: %s",
			result,
		)
	}
	if !strings.Contains(result, "ALWAYS") {
		t.Errorf(
			"result should contain schedule mode value: %s",
			result,
		)
	}
	if !strings.Contains(result, "NETWORK") {
		t.Errorf(
			"result should contain source traffic filter: %s",
			result,
		)
	}
	if !strings.Contains(result, "IP_ADDRESS") {
		t.Errorf(
			"result should contain destination traffic filter: %s",
			result,
		)
	}
}

func TestCreateFirewallPolicy_Execute_WithAllOptionalFields(
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
		"name": "Test Policy",
		"enabled": true,
		"action": {"type": "ALLOW"},
		"source": {
			"zoneId": "aaa00000-0000-0000-0000-000000000001",
			"trafficFilter": {"type": "NETWORK"}
		},
		"destination": {
			"zoneId": "aaa00000-0000-0000-0000-000000000002",
			"trafficFilter": {"type": "IP_ADDRESS"}
		},
		"ipProtocolScope": {"ipVersion": "IPV4"},
		"loggingEnabled": false,
		"schedule": {"mode": "ALWAYS"}
	}`
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src, ok := gotBody["source"].(map[string]interface{})
	if !ok {
		t.Fatal("request body should have source object")
	}
	if src["trafficFilter"] == nil {
		t.Error("request body source.trafficFilter should be set")
	}

	dst, ok := gotBody["destination"].(map[string]interface{})
	if !ok {
		t.Fatal("request body should have destination object")
	}
	if dst["trafficFilter"] == nil {
		t.Error(
			"request body destination.trafficFilter should be set",
		)
	}

	if gotBody["schedule"] == nil {
		t.Error("request body schedule should be set")
	}
}

// TestBuildRequestBody_SourceIPAddressFilter verifies that when
// source.trafficFilter.type is IP_ADDRESS and an ipAddressFilter is provided,
// the built request body includes the ipAddressFilter with the correct
// addresses. This test FAILS against the current code because
// policyTrafficFilterParams does not have an IpAddressFilter field.
func TestBuildRequestBody_SourceIPAddressFilter(t *testing.T) {
	var gotBody map[string]any
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
		"name": "Test IP Address Filter",
		"enabled": true,
		"action": {"type": "ALLOW"},
		"source": {
			"zoneId": "aaa00000-0000-0000-0000-000000000001",
			"trafficFilter": {
				"type": "IP_ADDRESS",
				"ipAddressFilter": {
					"type": "IP_ADDRESSES",
					"matchOpposite": false,
					"items": [{"cidr": "192.168.1.0/24"}]
				}
			}
		},
		"destination": {
			"zoneId": "aaa00000-0000-0000-0000-000000000002"
		},
		"ipProtocolScope": {"ipVersion": "IPV4"},
		"loggingEnabled": false
	}`
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src, ok := gotBody["source"].(map[string]any)
	if !ok {
		t.Fatal("request body should have source object")
	}
	tf, ok := src["trafficFilter"].(map[string]any)
	if !ok {
		t.Fatal("request body source.trafficFilter should be set")
	}
	if tf["type"] != "IP_ADDRESS" {
		t.Errorf(
			"source.trafficFilter.type = %v, want IP_ADDRESS",
			tf["type"],
		)
	}
	// This assertion FAILS with the current code: ipAddressFilter is dropped
	// because policyTrafficFilterParams has no field for it.
	if tf["ipAddressFilter"] == nil {
		t.Error(
			"source.trafficFilter.ipAddressFilter should be set" +
				" when type is IP_ADDRESS",
		)
	}
}

// TestBuildRequestBody_DestinationDomainFilter verifies that when
// destination.trafficFilter.type is DOMAIN and a domainFilter is provided,
// the built request body includes the domainFilter with the correct domains.
// This test FAILS against the current code because policyTrafficFilterParams
// does not have a DomainFilter field.
func TestBuildRequestBody_DestinationDomainFilter(t *testing.T) {
	var gotBody map[string]any
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
		"name": "Test Domain Filter",
		"enabled": true,
		"action": {"type": "BLOCK"},
		"source": {
			"zoneId": "aaa00000-0000-0000-0000-000000000001"
		},
		"destination": {
			"zoneId": "aaa00000-0000-0000-0000-000000000002",
			"trafficFilter": {
				"type": "DOMAIN",
				"domainFilter": {
					"type": "DOMAINS",
					"domains": ["example.com", "evil.test"]
				}
			}
		},
		"ipProtocolScope": {"ipVersion": "IPV4"},
		"loggingEnabled": false
	}`
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dst, ok := gotBody["destination"].(map[string]any)
	if !ok {
		t.Fatal("request body should have destination object")
	}
	tf, ok := dst["trafficFilter"].(map[string]any)
	if !ok {
		t.Fatal("request body destination.trafficFilter should be set")
	}
	if tf["type"] != "DOMAIN" {
		t.Errorf(
			"destination.trafficFilter.type = %v, want DOMAIN",
			tf["type"],
		)
	}
	// This assertion FAILS with the current code: domainFilter is dropped
	// because policyTrafficFilterParams has no field for it.
	if tf["domainFilter"] == nil {
		t.Error(
			"destination.trafficFilter.domainFilter should be set" +
				" when type is DOMAIN",
		)
	}
	df, ok := tf["domainFilter"].(map[string]any)
	if ok {
		domains, ok := df["domains"].([]any)
		if !ok || len(domains) != 2 {
			t.Errorf(
				"domainFilter.domains = %v, want 2 entries",
				df["domains"],
			)
		}
	}
}

// TestBuildRequestBody_SourceMACAddressFilter verifies that when
// source.trafficFilter.type is MAC_ADDRESS and a macAddressFilter is provided,
// the built request body includes the macAddressFilter with the correct MAC
// addresses. This test FAILS against the current code because
// policyTrafficFilterParams does not have a MacAddressFilter field.
func TestBuildRequestBody_SourceMACAddressFilter(t *testing.T) {
	var gotBody map[string]any
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
		"name": "Test MAC Address Filter",
		"enabled": true,
		"action": {"type": "BLOCK"},
		"source": {
			"zoneId": "aaa00000-0000-0000-0000-000000000001",
			"trafficFilter": {
				"type": "MAC_ADDRESS",
				"macAddressFilter": {
					"macAddresses": ["aa:bb:cc:dd:ee:ff", "11:22:33:44:55:66"]
				}
			}
		},
		"destination": {
			"zoneId": "aaa00000-0000-0000-0000-000000000002"
		},
		"ipProtocolScope": {"ipVersion": "IPV4"},
		"loggingEnabled": false
	}`
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src, ok := gotBody["source"].(map[string]any)
	if !ok {
		t.Fatal("request body should have source object")
	}
	tf, ok := src["trafficFilter"].(map[string]any)
	if !ok {
		t.Fatal("request body source.trafficFilter should be set")
	}
	if tf["type"] != "MAC_ADDRESS" {
		t.Errorf(
			"source.trafficFilter.type = %v, want MAC_ADDRESS",
			tf["type"],
		)
	}
	// This assertion FAILS with the current code: macAddressFilter is dropped
	// because policyTrafficFilterParams has no field for it.
	if tf["macAddressFilter"] == nil {
		t.Error(
			"source.trafficFilter.macAddressFilter should be set" +
				" when type is MAC_ADDRESS",
		)
	}
	maf, ok := tf["macAddressFilter"].(map[string]any)
	if ok {
		macs, ok := maf["macAddresses"].([]any)
		if !ok || len(macs) != 2 {
			t.Errorf(
				"macAddressFilter.macAddresses = %v, want 2 entries",
				maf["macAddresses"],
			)
		}
	}
}

// TestBuildRequestBody_SourceNetworkFilter verifies that when
// source.trafficFilter.type is NETWORK and a networkFilter is provided, the
// built request body includes the networkFilter with the correct network IDs.
// This test FAILS against the current code because policyTrafficFilterParams
// does not have a NetworkFilter field.
func TestBuildRequestBody_SourceNetworkFilter(t *testing.T) {
	var gotBody map[string]any
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
		"name": "Test Network Filter",
		"enabled": true,
		"action": {"type": "ALLOW"},
		"source": {
			"zoneId": "aaa00000-0000-0000-0000-000000000001",
			"trafficFilter": {
				"type": "NETWORK",
				"networkFilter": {
					"matchOpposite": false,
					"networkIds": [
						"bbb00000-0000-0000-0000-000000000001",
						"bbb00000-0000-0000-0000-000000000002"
					]
				}
			}
		},
		"destination": {
			"zoneId": "aaa00000-0000-0000-0000-000000000002"
		},
		"ipProtocolScope": {"ipVersion": "IPV4"},
		"loggingEnabled": false
	}`
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	src, ok := gotBody["source"].(map[string]any)
	if !ok {
		t.Fatal("request body should have source object")
	}
	tf, ok := src["trafficFilter"].(map[string]any)
	if !ok {
		t.Fatal("request body source.trafficFilter should be set")
	}
	if tf["type"] != "NETWORK" {
		t.Errorf(
			"source.trafficFilter.type = %v, want NETWORK",
			tf["type"],
		)
	}
	// This assertion FAILS with the current code: networkFilter is dropped
	// because policyTrafficFilterParams has no field for it.
	if tf["networkFilter"] == nil {
		t.Error(
			"source.trafficFilter.networkFilter should be set" +
				" when type is NETWORK",
		)
	}
	nf, ok := tf["networkFilter"].(map[string]any)
	if ok {
		ids, ok := nf["networkIds"].([]any)
		if !ok || len(ids) != 2 {
			t.Errorf(
				"networkFilter.networkIds = %v, want 2 entries",
				nf["networkIds"],
			)
		}
		if nf["matchOpposite"] != false {
			t.Errorf(
				"networkFilter.matchOpposite = %v, want false",
				nf["matchOpposite"],
			)
		}
	}
}

func TestListFirewallPolicies_Description(t *testing.T) {
	tool := &ListFirewallPolicies{}
	desc := tool.Description()
	if desc == "" {
		t.Fatal("Description() should not be empty")
	}
	if !strings.Contains(desc, "firewall") {
		t.Error("Description() should mention firewall")
	}
}

func TestCreateFirewallPolicy_Description(t *testing.T) {
	tool := &CreateFirewallPolicy{}
	desc := tool.Description()
	if desc == "" {
		t.Fatal("Description() should not be empty")
	}
	if !strings.Contains(desc, "firewall") {
		t.Error("Description() should mention firewall")
	}
}

func TestListFirewallPolicies_Execute_Formatting(t *testing.T) {
	p1 := mockPolicyJSON()
	p2 := mockPolicyJSON()
	p2["id"] = "ccc00000-0000-0000-0000-000000000002"
	p2["name"] = "Block WAN to LAN"

	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(
				paginatedResponse(p1, p2),
			)
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

	// Verify 1-based numbering
	if !strings.Contains(result, "1. Name: Allow LAN to WAN") {
		t.Errorf(
			"result should contain '1. Name: Allow LAN to WAN': %s",
			result,
		)
	}
	if !strings.Contains(result, "2. Name: Block WAN to LAN") {
		t.Errorf(
			"result should contain '2. Name: Block WAN to LAN': %s",
			result,
		)
	}

	// Verify blank line separator between policies
	if !strings.Contains(result, "\n\n2. ") {
		t.Errorf(
			"result should have blank line between policies: %s",
			result,
		)
	}
}

func TestDeleteFirewallPolicy_Execute_Message(t *testing.T) {
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
	if !strings.Contains(
		result,
		"ccc00000-0000-0000-0000-000000000001",
	) {
		t.Errorf(
			"result should contain policy ID: %s",
			result,
		)
	}
	if !strings.Contains(result, "Firewall policy") {
		t.Errorf(
			"result should mention 'Firewall policy': %s",
			result,
		)
	}
}

func TestFormatPolicy_EmptyConnectionStateFilter(t *testing.T) {
	// Policy with non-nil but empty ConnectionStateFilter
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			p := mockPolicyJSON()
			p["connectionStateFilter"] = []string{}
			json.NewEncoder(w).Encode(p)
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
	// Empty filter should not appear in output
	if strings.Contains(result, "Connection State Filter") {
		t.Errorf(
			"result should not show empty connection state filter: %s",
			result,
		)
	}
}

func TestCreateFirewallPolicy_EmptyConnectionStateFilter(
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
		"name": "Test",
		"enabled": true,
		"action": {"type": "ALLOW"},
		"source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"},
		"destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"},
		"ipProtocolScope": {"ipVersion": "IPV4"},
		"loggingEnabled": false,
		"connectionStateFilter": []
	}`
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Empty connectionStateFilter should NOT be sent in body
	if gotBody["connectionStateFilter"] != nil {
		t.Errorf(
			"empty connectionStateFilter should not be sent: %v",
			gotBody["connectionStateFilter"],
		)
	}
}

func TestListFirewallPolicies_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListFirewallPolicies(client, testSiteID)
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

func TestGetFirewallPolicy_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewGetFirewallPolicy(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallPolicyId": "ccc00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestCreateFirewallPolicy_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewCreateFirewallPolicy(client, testSiteID)
	args := `{
		"name": "Test",
		"enabled": true,
		"action": {"type": "ALLOW"},
		"source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"},
		"destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"},
		"ipProtocolScope": {"ipVersion": "IPV4"},
		"loggingEnabled": false
	}`
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestUpdateFirewallPolicy_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewUpdateFirewallPolicy(client, testSiteID)
	args := `{
		"firewallPolicyId": "ccc00000-0000-0000-0000-000000000001",
		"name": "Test",
		"enabled": true,
		"action": {"type": "ALLOW"},
		"source": {"zoneId": "aaa00000-0000-0000-0000-000000000001"},
		"destination": {"zoneId": "aaa00000-0000-0000-0000-000000000002"},
		"ipProtocolScope": {"ipVersion": "IPV4"},
		"loggingEnabled": false
	}`
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestDeleteFirewallPolicy_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewDeleteFirewallPolicy(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallPolicyId": "ccc00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestPatchFirewallPolicy_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewPatchFirewallPolicy(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"firewallPolicyId": "ccc00000-0000-0000-0000-000000000001", "loggingEnabled": true}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}
