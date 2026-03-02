package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetFirewallPolicyOrdering_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("sourceFirewallZoneId") == "" {
				t.Error("expected sourceFirewallZoneId param")
			}
			if q.Get("destinationFirewallZoneId") == "" {
				t.Error(
					"expected destinationFirewallZoneId param",
				)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"orderedFirewallPolicyIds": map[string]interface{}{
					"beforeSystemDefined": []string{
						"ccc00000-0000-0000-0000-000000000001",
					},
					"afterSystemDefined": []string{},
				},
			})
		}),
	)
	defer srv.Close()

	tool := NewGetFirewallPolicyOrdering(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Before System-Defined") {
		t.Errorf(
			"result should contain 'Before System-Defined': %s",
			result,
		)
	}
	if !strings.Contains(
		result,
		"ccc00000-0000-0000-0000-000000000001",
	) {
		t.Errorf("result should contain policy ID: %s", result)
	}
}

func TestGetFirewallPolicyOrdering_Execute_MissingParams(
	t *testing.T,
) {
	tool := &GetFirewallPolicyOrdering{
		baseTool{defaultSiteID: testSiteID},
	}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for missing params")
	}
	if !strings.Contains(err.Error(), "sourceZoneId") {
		t.Errorf("error should mention sourceZoneId: %v", err)
	}
}

func TestGetFirewallPolicyOrdering_InputSchema(t *testing.T) {
	tool := &GetFirewallPolicyOrdering{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	if len(required) != 2 {
		t.Errorf("expected 2 required fields, got %d", len(required))
	}
}

func TestUpdateFirewallPolicyOrdering_Execute(t *testing.T) {
	var gotBody map[string]interface{}
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &gotBody)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"orderedFirewallPolicyIds": map[string]interface{}{
					"beforeSystemDefined": []string{
						"ccc00000-0000-0000-0000-000000000001",
					},
					"afterSystemDefined": []string{},
				},
			})
		}),
	)
	defer srv.Close()

	tool := NewUpdateFirewallPolicyOrdering(client, testSiteID)
	args := `{
		"sourceZoneId": "aaa00000-0000-0000-0000-000000000001",
		"destinationZoneId": "aaa00000-0000-0000-0000-000000000002",
		"beforeSystemDefined": ["ccc00000-0000-0000-0000-000000000001"],
		"afterSystemDefined": []
	}`
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(args),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "updated successfully") {
		t.Errorf("result should confirm update: %s", result)
	}

	ordering, ok := gotBody["orderedFirewallPolicyIds"].(map[string]interface{})
	if !ok {
		t.Fatal("request body should have orderedFirewallPolicyIds")
	}
	before, ok := ordering["beforeSystemDefined"].([]interface{})
	if !ok || len(before) != 1 {
		t.Errorf(
			"request body beforeSystemDefined = %v",
			ordering["beforeSystemDefined"],
		)
	}
}

func TestUpdateFirewallPolicyOrdering_Execute_MissingParams(
	t *testing.T,
) {
	tool := &UpdateFirewallPolicyOrdering{
		baseTool{defaultSiteID: testSiteID},
	}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for missing params")
	}
	if !strings.Contains(err.Error(), "sourceZoneId") {
		t.Errorf("error should mention sourceZoneId: %v", err)
	}
}

func TestUpdateFirewallPolicyOrdering_InputSchema(t *testing.T) {
	tool := &UpdateFirewallPolicyOrdering{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	if len(required) != 4 {
		t.Errorf(
			"expected 4 required fields, got %d",
			len(required),
		)
	}
}

func TestGetFirewallPolicyOrdering_Execute_Formatting(
	t *testing.T,
) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"orderedFirewallPolicyIds": map[string]interface{}{
					"beforeSystemDefined": []string{
						"ccc00000-0000-0000-0000-000000000001",
					},
					"afterSystemDefined": []string{
						"ccc00000-0000-0000-0000-000000000002",
					},
				},
			})
		}),
	)
	defer srv.Close()

	tool := NewGetFirewallPolicyOrdering(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Policy Ordering:") {
		t.Errorf(
			"result should contain 'Policy Ordering:': %s",
			result,
		)
	}
	if !strings.Contains(
		result,
		"1. ccc00000-0000-0000-0000-000000000001",
	) {
		t.Errorf(
			"result should contain 1-based numbering for before: %s",
			result,
		)
	}
	if !strings.Contains(
		result,
		"1. ccc00000-0000-0000-0000-000000000002",
	) {
		t.Errorf(
			"result should contain 1-based numbering for after: %s",
			result,
		)
	}
}

func TestGetFirewallPolicyOrdering_Execute_NoSiteID(t *testing.T) {
	tool := &GetFirewallPolicyOrdering{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error when no site ID")
	}
	if !strings.Contains(err.Error(), "siteId") {
		t.Errorf("error should mention siteId: %v", err)
	}
}

func TestGetFirewallPolicyOrdering_Execute_InvalidDestZone(
	t *testing.T,
) {
	tool := &GetFirewallPolicyOrdering{
		baseTool{defaultSiteID: testSiteID},
	}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "not-valid"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid destination zone ID")
	}
	if !strings.Contains(err.Error(), "destinationZoneId") {
		t.Errorf(
			"error should mention destinationZoneId: %v",
			err,
		)
	}
}

func TestGetFirewallPolicyOrdering_Execute_NetworkError(
	t *testing.T,
) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewGetFirewallPolicyOrdering(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to get policy ordering",
	) {
		t.Errorf("error should contain expected message: %v", err)
	}
}

func TestUpdateFirewallPolicyOrdering_Execute_NoSiteID(
	t *testing.T,
) {
	tool := &UpdateFirewallPolicyOrdering{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002", "beforeSystemDefined": [], "afterSystemDefined": []}`,
		),
	)
	if err == nil {
		t.Fatal("expected error when no site ID")
	}
	if !strings.Contains(err.Error(), "siteId") {
		t.Errorf("error should mention siteId: %v", err)
	}
}

func TestUpdateFirewallPolicyOrdering_Execute_InvalidDestZone(
	t *testing.T,
) {
	tool := &UpdateFirewallPolicyOrdering{
		baseTool{defaultSiteID: testSiteID},
	}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "not-valid", "beforeSystemDefined": [], "afterSystemDefined": []}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid destination zone ID")
	}
	if !strings.Contains(err.Error(), "destinationZoneId") {
		t.Errorf(
			"error should mention destinationZoneId: %v",
			err,
		)
	}
}

func TestUpdateFirewallPolicyOrdering_Execute_InvalidBeforeUUID(
	t *testing.T,
) {
	tool := &UpdateFirewallPolicyOrdering{
		baseTool{defaultSiteID: testSiteID},
	}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002", "beforeSystemDefined": ["not-valid"], "afterSystemDefined": []}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid before UUID")
	}
	if !strings.Contains(err.Error(), "beforeSystemDefined") {
		t.Errorf(
			"error should mention beforeSystemDefined: %v",
			err,
		)
	}
}

func TestUpdateFirewallPolicyOrdering_Execute_InvalidAfterUUID(
	t *testing.T,
) {
	tool := &UpdateFirewallPolicyOrdering{
		baseTool{defaultSiteID: testSiteID},
	}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002", "beforeSystemDefined": [], "afterSystemDefined": ["not-valid"]}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid after UUID")
	}
	if !strings.Contains(err.Error(), "afterSystemDefined") {
		t.Errorf(
			"error should mention afterSystemDefined: %v",
			err,
		)
	}
}

func TestUpdateFirewallPolicyOrdering_Execute_NetworkError(
	t *testing.T,
) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewUpdateFirewallPolicyOrdering(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002", "beforeSystemDefined": [], "afterSystemDefined": []}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to update policy ordering",
	) {
		t.Errorf("error should contain expected message: %v", err)
	}
}

func TestGetFirewallPolicyOrdering_Execute_EmptyBefore(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"orderedFirewallPolicyIds": map[string]interface{}{
					"beforeSystemDefined": []string{},
					"afterSystemDefined": []string{
						"ccc00000-0000-0000-0000-000000000001",
					},
				},
			})
		}),
	)
	defer srv.Close()

	tool := NewGetFirewallPolicyOrdering(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Before System-Defined") {
		t.Errorf(
			"result should contain 'Before System-Defined': %s",
			result,
		)
	}
	if !strings.Contains(result, "(none)") {
		t.Errorf(
			"result should contain '(none)' for empty before section: %s",
			result,
		)
	}
	if !strings.Contains(
		result,
		"ccc00000-0000-0000-0000-000000000001",
	) {
		t.Errorf(
			"result should contain after-system-defined policy ID: %s",
			result,
		)
	}
}

func TestGetFirewallPolicyOrdering_Execute_EmptyAfter(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"orderedFirewallPolicyIds": map[string]interface{}{
					"beforeSystemDefined": []string{
						"ccc00000-0000-0000-0000-000000000001",
					},
					"afterSystemDefined": []string{},
				},
			})
		}),
	)
	defer srv.Close()

	tool := NewGetFirewallPolicyOrdering(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "After System-Defined") {
		t.Errorf(
			"result should contain 'After System-Defined': %s",
			result,
		)
	}
	if !strings.Contains(result, "(none)") {
		t.Errorf(
			"result should contain '(none)' for empty after section: %s",
			result,
		)
	}
	if !strings.Contains(
		result,
		"ccc00000-0000-0000-0000-000000000001",
	) {
		t.Errorf(
			"result should contain before-system-defined policy ID: %s",
			result,
		)
	}
}

func TestGetFirewallPolicyOrdering_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewGetFirewallPolicyOrdering(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"sourceZoneId": "aaa00000-0000-0000-0000-000000000001", "destinationZoneId": "aaa00000-0000-0000-0000-000000000002"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestUpdateFirewallPolicyOrdering_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewUpdateFirewallPolicyOrdering(client, testSiteID)
	args := `{
		"sourceZoneId": "aaa00000-0000-0000-0000-000000000001",
		"destinationZoneId": "aaa00000-0000-0000-0000-000000000002",
		"beforeSystemDefined": ["ccc00000-0000-0000-0000-000000000001"],
		"afterSystemDefined": []
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
