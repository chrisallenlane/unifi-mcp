package tools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func mockTrafficMatchingListJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":   "ddd00000-0000-0000-0000-000000000001",
		"name": "Block List",
		"type": "IPV4_ADDRESSES",
	}
}

func TestListTrafficMatchingLists_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				mockTrafficMatchingListJSON(),
			))
		}),
	)
	defer srv.Close()

	tool := NewListTrafficMatchingLists(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Block List") {
		t.Errorf("result should contain list name: %s", result)
	}
	if !strings.Contains(result, "IPV4_ADDRESSES") {
		t.Errorf("result should contain type: %s", result)
	}
}

func TestListTrafficMatchingLists_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListTrafficMatchingLists(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "No traffic matching lists found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListTrafficMatchingLists_Execute_NoSiteID(t *testing.T) {
	tool := &ListTrafficMatchingLists{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID")
	}
}

func TestGetTrafficMatchingList_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockTrafficMatchingListJSON())
		}),
	)
	defer srv.Close()

	tool := NewGetTrafficMatchingList(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"trafficMatchingListId": "ddd00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Block List") {
		t.Errorf("result should contain list name: %s", result)
	}
	if !strings.Contains(result, "IPV4_ADDRESSES") {
		t.Errorf("result should contain type: %s", result)
	}
}

func TestGetTrafficMatchingList_Execute_MissingID(t *testing.T) {
	tool := &GetTrafficMatchingList{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for missing traffic matching list ID")
	}
}

func TestGetTrafficMatchingList_InputSchema(t *testing.T) {
	tool := &GetTrafficMatchingList{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "trafficMatchingListId" {
			found = true
		}
	}
	if !found {
		t.Error("trafficMatchingListId should be required")
	}
}

func TestCreateTrafficMatchingList_Execute(t *testing.T) {
	var gotBody map[string]interface{}
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
	args := `{"name": "Block List", "type": "IPV4_ADDRESSES"}`
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
	if !strings.Contains(result, "Block List") {
		t.Errorf("result should contain list name: %s", result)
	}

	if gotBody["name"] != "Block List" {
		t.Errorf(
			"request body name = %v, want Block List",
			gotBody["name"],
		)
	}
	if gotBody["type"] != "IPV4_ADDRESSES" {
		t.Errorf(
			"request body type = %v, want IPV4_ADDRESSES",
			gotBody["type"],
		)
	}
}

func TestCreateTrafficMatchingList_Execute_NoSiteID(t *testing.T) {
	tool := &CreateTrafficMatchingList{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "Test", "type": "PORTS"}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID")
	}
}

func TestCreateTrafficMatchingList_InputSchema(t *testing.T) {
	tool := &CreateTrafficMatchingList{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	if len(required) < 2 {
		t.Errorf(
			"expected at least 2 required fields, got %d",
			len(required),
		)
	}
}

func TestUpdateTrafficMatchingList_Execute(t *testing.T) {
	var gotBody map[string]interface{}
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
		"trafficMatchingListId": "ddd00000-0000-0000-0000-000000000001",
		"name": "Updated List",
		"type": "IPV4_ADDRESSES"
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
	if !strings.Contains(result, "Updated List") {
		t.Errorf("result should contain updated name: %s", result)
	}
}

func TestUpdateTrafficMatchingList_Execute_MissingID(t *testing.T) {
	tool := &UpdateTrafficMatchingList{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"name": "Test", "type": "PORTS"}`),
	)
	if err == nil {
		t.Fatal("expected error for missing traffic matching list ID")
	}
}

func TestUpdateTrafficMatchingList_InputSchema(t *testing.T) {
	tool := &UpdateTrafficMatchingList{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "trafficMatchingListId" {
			found = true
		}
	}
	if !found {
		t.Error("trafficMatchingListId should be required")
	}
}

func TestDeleteTrafficMatchingList_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	tool := NewDeleteTrafficMatchingList(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"trafficMatchingListId": "ddd00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "deleted successfully") {
		t.Errorf("result should confirm deletion: %s", result)
	}
}

func TestDeleteTrafficMatchingList_Execute_MissingID(t *testing.T) {
	tool := &DeleteTrafficMatchingList{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for missing traffic matching list ID")
	}
}

func TestDeleteTrafficMatchingList_InputSchema(t *testing.T) {
	tool := &DeleteTrafficMatchingList{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "trafficMatchingListId" {
			found = true
		}
	}
	if !found {
		t.Error("trafficMatchingListId should be required")
	}
}
