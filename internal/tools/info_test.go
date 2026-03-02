package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetInfo_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/info" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
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

	if result != "Application Version: 10.1.84" {
		t.Errorf("unexpected result: %s", result)
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
}

func TestGetInfo_Description(t *testing.T) {
	tool := &GetInfo{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestGetInfo_InputSchema(t *testing.T) {
	tool := &GetInfo{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}
