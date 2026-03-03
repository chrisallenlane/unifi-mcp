package tools

// Tests that create/update tools do NOT forward MCP-only parameters
// (siteId, resource ID fields) in the HTTP request body sent to the
// UniFi API. The current code forwards raw args directly via
// bytes.NewReader(args), which leaks these fields and causes the API
// to respond with 400 BAD_REQUEST.
//
// These tests reproduce bug #15 and are expected to FAIL against the
// current implementation.

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// --- acl_rules ---

func TestCreateACLRule_Execute_BodyExcludesSiteID(t *testing.T) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"id":      "aaa00000-0000-0000-0000-000000000001",
				"name":    "Block IoT",
				"type":    "IPV4",
				"action":  "BLOCK",
				"enabled": true,
				"index":   0,
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateACLRule(client, testSiteID)
	args := `{
		"siteId": "` + testSiteID + `",
		"type": "IPV4",
		"name": "Block IoT",
		"enabled": true,
		"action": "BLOCK"
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotBody["siteId"]; present {
		t.Error(
			"request body must not contain 'siteId', " +
				"but it was found",
		)
	}
}

func TestUpdateACLRule_Execute_BodyExcludesMCPFields(t *testing.T) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":      "aaa00000-0000-0000-0000-000000000001",
				"name":    "Block IoT",
				"type":    "IPV4",
				"action":  "BLOCK",
				"enabled": true,
				"index":   0,
			})
		}),
	)
	defer srv.Close()

	tool := NewUpdateACLRule(client, testSiteID)
	args := `{
		"siteId": "` + testSiteID + `",
		"aclRuleId": "aaa00000-0000-0000-0000-000000000001",
		"type": "IPV4",
		"name": "Block IoT",
		"enabled": true,
		"action": "BLOCK"
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotBody["siteId"]; present {
		t.Error(
			"request body must not contain 'siteId', " +
				"but it was found",
		)
	}
	if _, present := gotBody["aclRuleId"]; present {
		t.Error(
			"request body must not contain 'aclRuleId', " +
				"but it was found",
		)
	}
}

// --- wifi ---

func TestCreateWiFiBroadcast_Execute_BodyExcludesSiteID(t *testing.T) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"id":       "bbb00000-0000-0000-0000-000000000001",
				"name":     "Guest WiFi",
				"type":     "STANDARD",
				"enabled":  true,
				"hideName": false,
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
	args := `{
		"siteId": "` + testSiteID + `",
		"name": "Guest WiFi",
		"enabled": true,
		"type": "STANDARD",
		"securityConfiguration": {"type": "OPEN"}
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotBody["siteId"]; present {
		t.Error(
			"request body must not contain 'siteId', " +
				"but it was found",
		)
	}
}

func TestUpdateWiFiBroadcast_Execute_BodyExcludesMCPFields(t *testing.T) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":       "bbb00000-0000-0000-0000-000000000001",
				"name":     "Updated WiFi",
				"type":     "STANDARD",
				"enabled":  false,
				"hideName": false,
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
	args := `{
		"siteId": "` + testSiteID + `",
		"wifiBroadcastId": "bbb00000-0000-0000-0000-000000000001",
		"name": "Updated WiFi",
		"enabled": false,
		"type": "STANDARD",
		"securityConfiguration": {"type": "WPA3"}
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotBody["siteId"]; present {
		t.Error(
			"request body must not contain 'siteId', " +
				"but it was found",
		)
	}
	if _, present := gotBody["wifiBroadcastId"]; present {
		t.Error(
			"request body must not contain 'wifiBroadcastId', " +
				"but it was found",
		)
	}
}

// --- traffic_matching_lists ---

func TestCreateTrafficMatchingList_Execute_BodyExcludesSiteID(
	t *testing.T,
) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(mockTrafficMatchingListJSON())
		}),
	)
	defer srv.Close()

	tool := NewCreateTrafficMatchingList(client, testSiteID)
	args := `{
		"siteId": "` + testSiteID + `",
		"name": "Block List",
		"type": "IPV4_ADDRESSES"
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotBody["siteId"]; present {
		t.Error(
			"request body must not contain 'siteId', " +
				"but it was found",
		)
	}
}

func TestUpdateTrafficMatchingList_Execute_BodyExcludesMCPFields(
	t *testing.T,
) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			updated := mockTrafficMatchingListJSON()
			updated["name"] = "Updated List"
			json.NewEncoder(w).Encode(updated)
		}),
	)
	defer srv.Close()

	tool := NewUpdateTrafficMatchingList(client, testSiteID)
	args := `{
		"siteId": "` + testSiteID + `",
		"trafficMatchingListId": "ddd00000-0000-0000-0000-000000000001",
		"name": "Updated List",
		"type": "IPV4_ADDRESSES"
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotBody["siteId"]; present {
		t.Error(
			"request body must not contain 'siteId', " +
				"but it was found",
		)
	}
	if _, present := gotBody["trafficMatchingListId"]; present {
		t.Error(
			"request body must not contain 'trafficMatchingListId', " +
				"but it was found",
		)
	}
}

// --- dns_policies ---

func TestCreateDNSPolicy_Execute_BodyExcludesSiteID(t *testing.T) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"id":      "ccc00000-0000-0000-0000-000000000001",
				"type":    "A_RECORD",
				"enabled": true,
				"domain":  "example.local",
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateDNSPolicy(client, testSiteID)
	args := `{
		"siteId": "` + testSiteID + `",
		"type": "A_RECORD",
		"enabled": true,
		"domain": "example.local",
		"ipv4Address": "192.168.1.1"
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotBody["siteId"]; present {
		t.Error(
			"request body must not contain 'siteId', " +
				"but it was found",
		)
	}
}

func TestUpdateDNSPolicy_Execute_BodyExcludesMCPFields(t *testing.T) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":      "ccc00000-0000-0000-0000-000000000001",
				"type":    "A_RECORD",
				"enabled": false,
				"domain":  "example.local",
			})
		}),
	)
	defer srv.Close()

	tool := NewUpdateDNSPolicy(client, testSiteID)
	args := `{
		"siteId": "` + testSiteID + `",
		"dnsPolicyId": "ccc00000-0000-0000-0000-000000000001",
		"type": "A_RECORD",
		"enabled": false
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotBody["siteId"]; present {
		t.Error(
			"request body must not contain 'siteId', " +
				"but it was found",
		)
	}
	if _, present := gotBody["dnsPolicyId"]; present {
		t.Error(
			"request body must not contain 'dnsPolicyId', " +
				"but it was found",
		)
	}
}

// --- networks ---

func TestCreateNetwork_Execute_BodyExcludesSiteID(t *testing.T) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"id":         "eee00000-0000-0000-0000-000000000001",
				"name":       "IoT",
				"vlanId":     100,
				"management": "GATEWAY",
				"enabled":    true,
				"default":    false,
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateNetwork(client, testSiteID)
	args := `{
		"siteId": "` + testSiteID + `",
		"name": "IoT",
		"enabled": true,
		"management": "GATEWAY",
		"vlanId": 100
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotBody["siteId"]; present {
		t.Error(
			"request body must not contain 'siteId', " +
				"but it was found",
		)
	}
}

func TestUpdateNetwork_Execute_BodyExcludesMCPFields(t *testing.T) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":         "eee00000-0000-0000-0000-000000000001",
				"name":       "IoT Updated",
				"vlanId":     100,
				"management": "GATEWAY",
				"enabled":    true,
				"default":    false,
			})
		}),
	)
	defer srv.Close()

	tool := NewUpdateNetwork(client, testSiteID)
	args := `{
		"siteId": "` + testSiteID + `",
		"networkId": "eee00000-0000-0000-0000-000000000001",
		"name": "IoT Updated",
		"enabled": true,
		"management": "GATEWAY",
		"vlanId": 100
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotBody["siteId"]; present {
		t.Error(
			"request body must not contain 'siteId', " +
				"but it was found",
		)
	}
	if _, present := gotBody["networkId"]; present {
		t.Error(
			"request body must not contain 'networkId', " +
				"but it was found",
		)
	}
}
