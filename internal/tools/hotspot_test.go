package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListVouchers_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":                   "aaa00000-0000-0000-0000-000000000001",
					"code":                 "ABCD-1234",
					"name":                 "Guest WiFi",
					"timeLimitMinutes":     60,
					"expired":              false,
					"authorizedGuestCount": 0,
					"createdAt":            "2026-03-01T10:00:00Z",
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListVouchers(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "ABCD-1234") {
		t.Errorf(
			"result should contain voucher code: %s",
			result,
		)
	}
	if !strings.Contains(result, "Guest WiFi") {
		t.Errorf(
			"result should contain voucher name: %s",
			result,
		)
	}
}

func TestListVouchers_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListVouchers(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No vouchers found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListVouchers_Execute_NoSiteID(t *testing.T) {
	tool := &ListVouchers{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID provided")
	}
}

func TestGetVoucher_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":                   "aaa00000-0000-0000-0000-000000000001",
				"code":                 "ABCD-1234",
				"name":                 "Guest WiFi",
				"timeLimitMinutes":     60,
				"expired":              false,
				"authorizedGuestCount": 2,
				"authorizedGuestLimit": 5,
				"createdAt":            "2026-03-01T10:00:00Z",
			})
		}),
	)
	defer srv.Close()

	tool := NewGetVoucher(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"voucherId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "ABCD-1234") {
		t.Errorf(
			"result should contain voucher code: %s",
			result,
		)
	}
	if !strings.Contains(result, "Authorized Guests: 2") {
		t.Errorf(
			"result should contain guest count: %s",
			result,
		)
	}
	if !strings.Contains(result, "Guest Limit: 5") {
		t.Errorf(
			"result should contain guest limit: %s",
			result,
		)
	}
}

func TestGetVoucher_Execute_InvalidUUID(t *testing.T) {
	tool := &GetVoucher{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"voucherId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestGetVoucher_InputSchema(t *testing.T) {
	tool := &GetVoucher{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "voucherId" {
			found = true
		}
	}
	if !found {
		t.Error("voucherId should be required")
	}
}

func TestCreateVouchers_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"vouchers": []map[string]interface{}{
					{
						"id":                   "bbb00000-0000-0000-0000-000000000001",
						"code":                 "WXYZ-5678",
						"name":                 "Event Pass",
						"timeLimitMinutes":     120,
						"expired":              false,
						"authorizedGuestCount": 0,
						"createdAt":            "2026-03-01T12:00:00Z",
					},
				},
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateVouchers(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Event Pass", "timeLimitMinutes": 120, "count": 1}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Generated 1 voucher(s)") {
		t.Errorf(
			"result should contain generation count: %s",
			result,
		)
	}
	if !strings.Contains(result, "WXYZ-5678") {
		t.Errorf(
			"result should contain voucher code: %s",
			result,
		)
	}
}

func TestCreateVouchers_InputSchema(t *testing.T) {
	tool := &CreateVouchers{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	foundName := false
	foundTime := false
	for _, r := range required {
		if r == "name" {
			foundName = true
		}
		if r == "timeLimitMinutes" {
			foundTime = true
		}
	}
	if !foundName {
		t.Error("name should be required")
	}
	if !foundTime {
		t.Error("timeLimitMinutes should be required")
	}
}

func TestDeleteVouchers_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"vouchersDeleted": 3,
			})
		}),
	)
	defer srv.Close()

	tool := NewDeleteVouchers(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"filter": "expired eq true"}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Deleted 3 voucher(s)") {
		t.Errorf(
			"result should contain deletion count: %s",
			result,
		)
	}
}

func TestDeleteVouchers_Execute_MissingFilter(t *testing.T) {
	tool := &DeleteVouchers{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when filter missing")
	}
}

func TestDeleteVouchers_InputSchema(t *testing.T) {
	tool := &DeleteVouchers{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "filter" {
			found = true
		}
	}
	if !found {
		t.Error("filter should be required")
	}
}

func TestDeleteVoucher_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"vouchersDeleted": 1,
			})
		}),
	)
	defer srv.Close()

	tool := NewDeleteVoucher(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"voucherId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Voucher deleted." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestDeleteVoucher_Execute_InvalidUUID(t *testing.T) {
	tool := &DeleteVoucher{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"voucherId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestDeleteVoucher_InputSchema(t *testing.T) {
	tool := &DeleteVoucher{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "voucherId" {
			found = true
		}
	}
	if !found {
		t.Error("voucherId should be required")
	}
}
