package tools

// Tests for bug #17 (create/update_traffic_matching_list missing items field)
// and bug #18 (create_dns_policy schema field names don't match API).
//
// These tests are written FIRST and are expected to FAIL against the current
// implementation. Once the bugs are fixed the tests should pass.

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// --- bug #17: traffic_matching_list items field ---

// TestCreateTrafficMatchingList_Execute_ItemsForwardedInBody verifies that
// the items array supplied by the caller is included in the request body
// sent to the UniFi API.
func TestCreateTrafficMatchingList_Execute_ItemsForwardedInBody(
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
		"name": "Port Block",
		"type": "PORTS",
		"items": [
			{"type": "PORT_NUMBER", "portNumber": 80},
			{"type": "PORT_NUMBER_RANGE", "from": 8000, "to": 9000}
		]
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	items, present := gotBody["items"]
	if !present {
		t.Fatal(
			"request body must contain 'items', but it was missing",
		)
	}

	itemSlice, ok := items.([]any)
	if !ok {
		t.Fatalf(
			"items must be an array, got %T",
			items,
		)
	}
	if len(itemSlice) != 2 {
		t.Errorf(
			"items should have 2 elements, got %d",
			len(itemSlice),
		)
	}
}

// TestCreateTrafficMatchingList_InputSchema_HasItemsField verifies that the
// input schema exposes the items field.
func TestCreateTrafficMatchingList_InputSchema_HasItemsField(t *testing.T) {
	tool := &CreateTrafficMatchingList{}
	schema := tool.InputSchema()

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("properties should be a map")
	}

	if _, present := props["items"]; !present {
		t.Error(
			"InputSchema properties must include 'items', but it is missing",
		)
	}
}

// TestCreateTrafficMatchingList_InputSchema_ItemsRequired verifies that items
// is included in the required list for create.
func TestCreateTrafficMatchingList_InputSchema_ItemsRequired(t *testing.T) {
	tool := &CreateTrafficMatchingList{}
	schema := tool.InputSchema()

	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "items")
}

// TestUpdateTrafficMatchingList_Execute_ItemsForwardedInBody verifies that
// the items array supplied by the caller is included in the request body for
// update operations.
func TestUpdateTrafficMatchingList_Execute_ItemsForwardedInBody(
	t *testing.T,
) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			updated := mockTrafficMatchingListJSON()
			updated["name"] = "Updated Port Block"
			json.NewEncoder(w).Encode(updated)
		}),
	)
	defer srv.Close()

	tool := NewUpdateTrafficMatchingList(client, testSiteID)
	args := `{
		"trafficMatchingListId": "ddd00000-0000-0000-0000-000000000001",
		"name": "Updated Port Block",
		"type": "PORTS",
		"items": [
			{"type": "PORT_NUMBER", "portNumber": 443}
		]
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	items, present := gotBody["items"]
	if !present {
		t.Fatal(
			"request body must contain 'items', but it was missing",
		)
	}

	itemSlice, ok := items.([]any)
	if !ok {
		t.Fatalf(
			"items must be an array, got %T",
			items,
		)
	}
	if len(itemSlice) != 1 {
		t.Errorf(
			"items should have 1 element, got %d",
			len(itemSlice),
		)
	}
}

// TestUpdateTrafficMatchingList_InputSchema_ItemsRequired verifies that items
// is included in the required list for update.
func TestUpdateTrafficMatchingList_InputSchema_ItemsRequired(t *testing.T) {
	tool := &UpdateTrafficMatchingList{}
	schema := tool.InputSchema()

	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "items")
}

// --- bug #18: dns_policy schema field names ---

// TestCreateDNSPolicy_Execute_ARecord_CorrectFieldNames verifies that an
// A_RECORD creation request sends ipv4Address (not address) to the API.
func TestCreateDNSPolicy_Execute_ARecord_CorrectFieldNames(t *testing.T) {
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
				"domain":  "nas.local",
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateDNSPolicy(client, testSiteID)
	args := `{
		"type": "A_RECORD",
		"enabled": true,
		"domain": "nas.local",
		"ipv4Address": "192.168.1.50",
		"ttlSeconds": 300
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The API field name is ipv4Address, not address.
	if _, present := gotBody["ipv4Address"]; !present {
		t.Error(
			"request body must contain 'ipv4Address', but it was missing",
		)
	}
	if _, present := gotBody["address"]; present {
		t.Error(
			"request body must not contain old field 'address', " +
				"but it was found",
		)
	}
	if _, present := gotBody["ttlSeconds"]; !present {
		t.Error(
			"request body must contain 'ttlSeconds', but it was missing",
		)
	}
}

// TestCreateDNSPolicy_Execute_CNameRecord_CorrectFieldNames verifies that a
// CNAME_RECORD creation request sends targetDomain (not target) to the API.
func TestCreateDNSPolicy_Execute_CNameRecord_CorrectFieldNames(
	t *testing.T,
) {
	var gotBody map[string]any
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"id":      "ccc00000-0000-0000-0000-000000000002",
				"type":    "CNAME_RECORD",
				"enabled": true,
				"domain":  "alias.local",
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateDNSPolicy(client, testSiteID)
	args := `{
		"type": "CNAME_RECORD",
		"enabled": true,
		"domain": "alias.local",
		"targetDomain": "real.local"
	}`
	_, err := tool.Execute(context.Background(), json.RawMessage(args))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The API field name is targetDomain, not target.
	if _, present := gotBody["targetDomain"]; !present {
		t.Error(
			"request body must contain 'targetDomain', but it was missing",
		)
	}
	if _, present := gotBody["target"]; present {
		t.Error(
			"request body must not contain old field 'target', " +
				"but it was found",
		)
	}
}

// TestCreateDNSPolicy_InputSchema_CorrectFieldNames verifies that the input
// schema exposes the correct API field names and not the old incorrect ones.
func TestCreateDNSPolicy_InputSchema_CorrectFieldNames(t *testing.T) {
	tool := &CreateDNSPolicy{}
	schema := tool.InputSchema()

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("properties should be a map")
	}

	// Fields that MUST be present (correct API names).
	mustHave := []string{
		"ipv4Address",
		"ipv6Address",
		"targetDomain",
		"mailServerDomain",
		"text",
		"ipAddress",
		"serverDomain",
		"service",
		"protocol",
		"ttlSeconds",
	}
	for _, f := range mustHave {
		if _, present := props[f]; !present {
			t.Errorf(
				"InputSchema properties must include %q, but it is missing",
				f,
			)
		}
	}

	// Fields that must NOT appear (old incorrect names).
	mustNotHave := []string{
		"address",
		"target",
		"server",
		"value",
		"forwardTo",
	}
	for _, f := range mustNotHave {
		if _, present := props[f]; present {
			t.Errorf(
				"InputSchema properties must not include old field %q",
				f,
			)
		}
	}
}
