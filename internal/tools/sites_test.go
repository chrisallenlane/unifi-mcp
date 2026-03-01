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

func TestListSites_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/sites" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
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
						"name":              "Branch Office",
						"internalReference": "branch",
					},
				},
				"count":      2,
				"limit":      25,
				"offset":     0,
				"totalCount": 2,
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

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
	srv := httptest.NewServer(
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

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewListSites(client, "")
	_, err = tool.Execute(
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

func TestListSites_Description(t *testing.T) {
	tool := &ListSites{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
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
