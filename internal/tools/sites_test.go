package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListSites_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":                "550e8400-e29b-41d4-a716-446655440000",
					"name":              "Default",
					"internalReference": "default",
				},
				map[string]interface{}{
					"id":                "660e8400-e29b-41d4-a716-446655440001",
					"name":              "Branch Office",
					"internalReference": "branch",
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListSites(client, "")
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Default") {
		t.Errorf("result should contain 'Default': %s", result)
	}
	if !strings.Contains(result, "Branch Office") {
		t.Errorf(
			"result should contain 'Branch Office': %s",
			result,
		)
	}
	if !strings.Contains(
		result,
		"550e8400-e29b-41d4-a716-446655440000",
	) {
		t.Errorf("result should contain site ID: %s", result)
	}
}

func TestListSites_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListSites(client, "")
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No sites found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListSites_Execute_WithParams(t *testing.T) {
	var gotLimit, gotOffset string
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotLimit = r.URL.Query().Get("limit")
			gotOffset = r.URL.Query().Get("offset")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data":       []interface{}{},
				"count":      0,
				"limit":      10,
				"offset":     5,
				"totalCount": 0,
			})
		}),
	)
	defer srv.Close()

	tool := NewListSites(client, "")
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"limit": 10, "offset": 5}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotLimit != "10" {
		t.Errorf("expected limit=10, got %s", gotLimit)
	}
	if gotOffset != "5" {
		t.Errorf("expected offset=5, got %s", gotOffset)
	}
}

func TestListSites_InputSchema(t *testing.T) {
	tool := &ListSites{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("properties should be a map")
	}
	if _, ok := props["limit"]; !ok {
		t.Error("schema should have limit property")
	}
	if _, ok := props["offset"]; !ok {
		t.Error("schema should have offset property")
	}
}

func TestListSites_Description(t *testing.T) {
	tool := &ListSites{}
	desc := tool.Description()
	if desc == "" {
		t.Fatal("Description() should not be empty")
	}
	if !strings.Contains(desc, "site") {
		t.Error("Description() should mention sites")
	}
}

func TestListSites_Execute_InvalidJSON(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			t.Fatal("API should not be called for invalid JSON")
		}),
	)
	defer srv.Close()

	tool := NewListSites(client, "")
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListSites_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewListSites(client, "")
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list sites") {
		t.Errorf(
			"error should contain 'failed to list sites': %v",
			err,
		)
	}
}

func TestListSites_Execute_Formatting(t *testing.T) {
	// Use totalCount=5 but only 2 items to distinguish
	// len(page.Data) from page.TotalCount in the header.
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":                "550e8400-e29b-41d4-a716-446655440000",
						"name":              "Default",
						"internalReference": "default",
					},
					{
						"id":                "660e8400-e29b-41d4-a716-446655440001",
						"name":              "Branch",
						"internalReference": "branch",
					},
				},
				"count":      2,
				"limit":      2,
				"offset":     0,
				"totalCount": 5,
			})
		}),
	)
	defer srv.Close()

	tool := NewListSites(client, "")
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify header shows page size vs total count
	if !strings.Contains(result, "Sites (2 of 5):") {
		t.Errorf(
			"result should contain 'Sites (2 of 5):': %s",
			result,
		)
	}

	// Verify 1-based numbering
	if !strings.Contains(result, "1. Default") {
		t.Errorf(
			"result should contain '1. Default': %s",
			result,
		)
	}
	if !strings.Contains(result, "2. Branch") {
		t.Errorf(
			"result should contain '2. Branch': %s",
			result,
		)
	}
}

func TestListSites_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListSites(client, "")
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for API error response")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}
