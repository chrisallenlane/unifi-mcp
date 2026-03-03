package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListPendingDevices_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"model":                 "USW-8",
					"macAddress":            "ff:ee:dd:cc:bb:aa",
					"ipAddress":             "10.0.0.50",
					"state":                 "PENDING_ADOPTION",
					"firmwareVersion":       "6.5.0",
					"firmwareUpdatable":     false,
					"supported":             true,
					"adoptionTargetSiteIds": []string{},
					"features":              []string{"switching"},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListPendingDevices(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "USW-8") {
		t.Errorf("result should contain model: %s", result)
	}
	if !strings.Contains(result, "ff:ee:dd:cc:bb:aa") {
		t.Errorf("result should contain MAC address: %s", result)
	}
	if !strings.Contains(result, "PENDING_ADOPTION") {
		t.Errorf("result should contain state: %s", result)
	}
}

func TestListPendingDevices_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListPendingDevices(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No pending devices found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListPendingDevices_Description(t *testing.T) {
	tool := &ListPendingDevices{}
	d := tool.Description()
	if d == "" {
		t.Fatal("description should not be empty")
	}
	if !strings.Contains(d, "pending") {
		t.Errorf(
			"description should mention pending: %s",
			d,
		)
	}
}

func TestListPendingDevices_Execute_InvalidJSON(t *testing.T) {
	tool := &ListPendingDevices{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListPendingDevices_Execute_NetworkError(
	t *testing.T,
) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewListPendingDevices(client)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to list pending devices",
	) {
		t.Errorf(
			"error should contain 'failed to list pending devices': %v",
			err,
		)
	}
}

func TestListPendingDevices_Execute_Formatting(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"model":             "USW-8",
						"macAddress":        "aa:bb:cc:dd:ee:ff",
						"ipAddress":         "10.0.0.1",
						"state":             "PENDING_ADOPTION",
						"firmwareVersion":   "6.5.0",
						"firmwareUpdatable": false,
						"supported":         true,
						"features":          []string{},
					},
					{
						"model":             "UAP-AC",
						"macAddress":        "11:22:33:44:55:66",
						"ipAddress":         "10.0.0.2",
						"state":             "PENDING_ADOPTION",
						"firmwareUpdatable": false,
						"supported":         true,
						"features":          []string{},
					},
				},
				"count":      2,
				"limit":      25,
				"offset":     0,
				"totalCount": 5,
			})
		}),
	)
	defer srv.Close()

	tool := NewListPendingDevices(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(
		result,
		"Pending Devices (2 of 5):",
	) {
		t.Errorf(
			"result should contain header: %s",
			result,
		)
	}
	if !strings.Contains(result, "1. Model: USW-8") {
		t.Errorf(
			"result should contain '1. Model: USW-8': %s",
			result,
		)
	}
	if !strings.Contains(result, "2. Model: UAP-AC") {
		t.Errorf(
			"result should contain '2. Model: UAP-AC': %s",
			result,
		)
	}
}

func TestListPendingDevices_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewListPendingDevices(client)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error on API failure")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf(
			"error should contain '500': %v",
			err,
		)
	}
}
