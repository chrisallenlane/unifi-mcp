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
	if !strings.Contains(err.Error(), "siteId") {
		t.Errorf("error should mention siteId: %v", err)
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
	if !strings.Contains(err.Error(), "voucherId") {
		t.Errorf("error should mention voucherId: %v", err)
	}
}

func TestGetVoucher_InputSchema(t *testing.T) {
	tool := &GetVoucher{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "voucherId")
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
	requireContains(t, required, "name")
	requireContains(t, required, "timeLimitMinutes")
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
	if !strings.Contains(err.Error(), "filter") {
		t.Errorf("error should mention filter: %v", err)
	}
}

func TestDeleteVouchers_InputSchema(t *testing.T) {
	tool := &DeleteVouchers{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "filter")
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
	if !strings.Contains(err.Error(), "voucherId") {
		t.Errorf("error should mention voucherId: %v", err)
	}
}

func TestDeleteVoucher_InputSchema(t *testing.T) {
	tool := &DeleteVoucher{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "voucherId")
}

func TestGetVoucher_Execute_WithOptionalFields(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":                   "aaa00000-0000-0000-0000-000000000001",
				"code":                 "ABCD-1234",
				"name":                 "Guest WiFi",
				"timeLimitMinutes":     60,
				"expired":              false,
				"authorizedGuestCount": 0,
				"authorizedGuestLimit": 5,
				"createdAt":            "2025-01-01T00:00:00Z",
				"activatedAt":          "2025-01-01T00:00:00Z",
				"expiresAt":            "2025-02-01T00:00:00Z",
				"dataUsageLimitMBytes": 1024,
				"rxRateLimitKbps":      5000,
				"txRateLimitKbps":      1000,
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

	checks := []string{
		"Guest Limit:",
		"Activated At:",
		"Expires At:",
		"Data Limit:",
		"Download Limit:",
		"Upload Limit:",
	}
	for _, s := range checks {
		if !strings.Contains(result, s) {
			t.Errorf(
				"result should contain %q: %s",
				s,
				result,
			)
		}
	}
}

func TestCreateVouchers_Execute_EmptyResult(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"vouchers": []interface{}{},
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateVouchers(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Event Pass", "timeLimitMinutes": 60}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No vouchers generated." {
		t.Errorf("unexpected result: %s", result)
	}
}

// --- formatting and coverage tests ---

func TestListVouchers_Description(t *testing.T) {
	tool := &ListVouchers{}
	d := tool.Description()
	if d == "" {
		t.Fatal("description should not be empty")
	}
	if !strings.Contains(d, "oucher") {
		t.Errorf(
			"description should mention vouchers: %s",
			d,
		)
	}
}

func TestListVouchers_Execute_InvalidJSON(t *testing.T) {
	tool := &ListVouchers{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListVouchers_Execute_Formatting(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":                   "aaa00000-0000-0000-0000-000000000001",
						"code":                 "ABCD-1234",
						"name":                 "Guest WiFi",
						"timeLimitMinutes":     60,
						"expired":              false,
						"authorizedGuestCount": 0,
						"createdAt":            "2026-03-01T10:00:00Z",
					},
					{
						"id":                   "aaa00000-0000-0000-0000-000000000002",
						"code":                 "WXYZ-5678",
						"name":                 "Event Pass",
						"timeLimitMinutes":     120,
						"expired":              false,
						"authorizedGuestCount": 0,
						"createdAt":            "2026-03-01T11:00:00Z",
					},
				},
				"count":      2,
				"limit":      25,
				"offset":     0,
				"totalCount": 11,
			})
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

	// verify header uses totalCount, not len(data)
	if !strings.Contains(result, "Vouchers (2 of 11):") {
		t.Errorf(
			"result should contain 'Vouchers (2 of 11):': %s",
			result,
		)
	}

	// verify numbering
	if !strings.Contains(result, "1. Code: ABCD-1234") {
		t.Errorf(
			"result should contain '1. Code: ABCD-1234': %s",
			result,
		)
	}
	if !strings.Contains(result, "2. Code: WXYZ-5678") {
		t.Errorf(
			"result should contain '2. Code: WXYZ-5678': %s",
			result,
		)
	}
}

func TestListVouchers_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewListVouchers(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list vouchers") {
		t.Errorf(
			"error should contain 'failed to list vouchers': %v",
			err,
		)
	}
}

func TestGetVoucher_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewGetVoucher(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"voucherId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to get voucher") {
		t.Errorf(
			"error should contain 'failed to get voucher': %v",
			err,
		)
	}
}

func TestCreateVouchers_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewCreateVouchers(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Event Pass", "timeLimitMinutes": 60}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to create vouchers") {
		t.Errorf(
			"error should contain 'failed to create vouchers': %v",
			err,
		)
	}
}

func TestDeleteVouchers_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewDeleteVouchers(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"filter": "expired eq true"}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to delete vouchers") {
		t.Errorf(
			"error should contain 'failed to delete vouchers': %v",
			err,
		)
	}
}

func TestDeleteVoucher_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewDeleteVoucher(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"voucherId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to delete voucher") {
		t.Errorf(
			"error should contain 'failed to delete voucher': %v",
			err,
		)
	}
}

func TestListVouchers_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListVouchers(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestGetVoucher_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewGetVoucher(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"voucherId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestCreateVouchers_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewCreateVouchers(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Event Pass", "timeLimitMinutes": 60}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestDeleteVouchers_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewDeleteVouchers(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"filter": "expired eq true"}`),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestDeleteVoucher_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewDeleteVoucher(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"voucherId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}
