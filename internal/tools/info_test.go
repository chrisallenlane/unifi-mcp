package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestGetInfo_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"applicationVersion": "10.1.84",
			})
		}),
	)
	defer srv.Close()

	tool := NewGetInfo(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "10.1.84") {
		t.Errorf(
			"result should contain version '10.1.84': %s",
			result,
		)
	}
}

func TestGetInfo_Execute_Error(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
		}),
	)
	defer srv.Close()

	tool := NewGetInfo(client)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for unauthorized response")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestGetInfo_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close() // close immediately to force a network error

	tool := NewGetInfo(client)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to get info") {
		t.Errorf(
			"error should contain 'failed to get info': %v",
			err,
		)
	}
}

func TestGetInfo_Description(t *testing.T) {
	tool := NewGetInfo(nil)
	desc := tool.Description()
	if desc == "" {
		t.Fatal("Description() should not be empty")
	}
	if !strings.Contains(desc, "UniFi") {
		t.Error("Description() should mention UniFi")
	}
}

func TestGetInfo_InputSchema(t *testing.T) {
	tool := NewGetInfo(nil)
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf(
			"schema type should be 'object', got %v",
			schema["type"],
		)
	}
	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("schema should have properties map")
	}
	if len(props) != 0 {
		t.Errorf(
			"get_info should have no properties, got %d",
			len(props),
		)
	}
}

func TestGetInfo_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewGetInfo(client)
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
