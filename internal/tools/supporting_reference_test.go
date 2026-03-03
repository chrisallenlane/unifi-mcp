package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// --- list_radius_profiles tests ---

func TestListRadiusProfiles_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":       "ddd00000-0000-0000-0000-000000000001",
					"name":     "Corp RADIUS",
					"metadata": map[string]string{"origin": "USER_DEFINED"},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListRadiusProfiles(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Corp RADIUS") {
		t.Errorf(
			"result should contain 'Corp RADIUS': %s",
			result,
		)
	}
}

func TestListRadiusProfiles_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListRadiusProfiles(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No RADIUS profiles found." {
		t.Errorf("unexpected result: %s", result)
	}
}

// --- list_device_tags tests ---

func TestListDeviceTags_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":   "eee00000-0000-0000-0000-000000000001",
					"name": "Floor 1 APs",
					"deviceIds": []string{
						"fff00000-0000-0000-0000-000000000001",
						"fff00000-0000-0000-0000-000000000002",
					},
					"metadata": map[string]string{"origin": "USER_DEFINED"},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListDeviceTags(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Floor 1 APs") {
		t.Errorf(
			"result should contain 'Floor 1 APs': %s",
			result,
		)
	}
	if !strings.Contains(result, "Devices: 2") {
		t.Errorf(
			"result should contain 'Devices: 2': %s",
			result,
		)
	}
}

func TestListDeviceTags_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListDeviceTags(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No device tags found." {
		t.Errorf("unexpected result: %s", result)
	}
}

// --- list_dpi_categories tests ---

func TestListDpiCategories_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{"id": 1, "name": "Streaming Media"},
				map[string]interface{}{"id": 2, "name": "Social Networking"},
			))
		}),
	)
	defer srv.Close()

	tool := NewListDpiCategories(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Streaming Media") {
		t.Errorf(
			"result should contain 'Streaming Media': %s",
			result,
		)
	}
	if !strings.Contains(result, "Social Networking") {
		t.Errorf(
			"result should contain 'Social Networking': %s",
			result,
		)
	}
}

func TestListDpiCategories_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListDpiCategories(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No DPI categories found." {
		t.Errorf("unexpected result: %s", result)
	}
}

// --- list_dpi_applications tests ---

func TestListDpiApplications_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{"id": 100, "name": "Netflix"},
				map[string]interface{}{"id": 101, "name": "YouTube"},
			))
		}),
	)
	defer srv.Close()

	tool := NewListDpiApplications(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Netflix") {
		t.Errorf(
			"result should contain 'Netflix': %s",
			result,
		)
	}
	if !strings.Contains(result, "YouTube") {
		t.Errorf(
			"result should contain 'YouTube': %s",
			result,
		)
	}
}

func TestListDpiApplications_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListDpiApplications(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No DPI applications found." {
		t.Errorf("unexpected result: %s", result)
	}
}

// --- list_countries tests ---

func TestListCountries_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{"code": "US", "name": "United States"},
				map[string]interface{}{"code": "DE", "name": "Germany"},
				map[string]interface{}{"code": "JP", "name": "Japan"},
			))
		}),
	)
	defer srv.Close()

	tool := NewListCountries(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "United States") {
		t.Errorf(
			"result should contain 'United States': %s",
			result,
		)
	}
	if !strings.Contains(result, "US") {
		t.Errorf("result should contain 'US': %s", result)
	}
	if !strings.Contains(result, "Germany") {
		t.Errorf(
			"result should contain 'Germany': %s",
			result,
		)
	}
}

func TestListCountries_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListCountries(client)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No countries found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListRadiusProfiles_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListRadiusProfiles(client, testSiteID)
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

func TestListDeviceTags_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListDeviceTags(client, testSiteID)
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

func TestListDpiCategories_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListDpiCategories(client)
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

func TestListDpiApplications_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListDpiApplications(client)
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

func TestListCountries_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListCountries(client)
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

// --- description tests ---

func TestListRadiusProfiles_Description(t *testing.T) {
	tool := NewListRadiusProfiles(nil, "")
	desc := tool.Description()
	if desc == "" {
		t.Fatal("Description() should not be empty")
	}
	if !strings.Contains(desc, "RADIUS") {
		t.Errorf("description should mention RADIUS: %s", desc)
	}
}

// --- invalid JSON tests ---

func TestListRadiusProfiles_Execute_InvalidJSON(t *testing.T) {
	tool := &ListRadiusProfiles{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListDeviceTags_Execute_InvalidJSON(t *testing.T) {
	tool := &ListDeviceTags{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListDpiCategories_Execute_InvalidJSON(t *testing.T) {
	tool := &ListDpiCategories{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListDpiApplications_Execute_InvalidJSON(t *testing.T) {
	tool := &ListDpiApplications{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListCountries_Execute_InvalidJSON(t *testing.T) {
	tool := &ListCountries{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// --- network error tests ---

func TestListRadiusProfiles_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewListRadiusProfiles(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list RADIUS profiles") {
		t.Errorf(
			"error should contain 'failed to list RADIUS profiles': %v",
			err,
		)
	}
}

func TestListDeviceTags_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewListDeviceTags(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list device tags") {
		t.Errorf(
			"error should contain 'failed to list device tags': %v",
			err,
		)
	}
}

func TestListDpiCategories_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewListDpiCategories(client)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list DPI categories") {
		t.Errorf(
			"error should contain 'failed to list DPI categories': %v",
			err,
		)
	}
}

func TestListDpiApplications_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewListDpiApplications(client)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list DPI applications") {
		t.Errorf(
			"error should contain 'failed to list DPI applications': %v",
			err,
		)
	}
}

func TestListCountries_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)
	srv.Close()

	tool := NewListCountries(client)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list countries") {
		t.Errorf(
			"error should contain 'failed to list countries': %v",
			err,
		)
	}
}

// --- formatting tests ---

func TestListRadiusProfiles_Execute_Formatting(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":       "ddd00000-0000-0000-0000-000000000001",
					"name":     "Corp RADIUS",
					"metadata": map[string]string{"origin": "USER_DEFINED"},
				},
				map[string]interface{}{
					"id":       "ddd00000-0000-0000-0000-000000000002",
					"name":     "Guest RADIUS",
					"metadata": map[string]string{"origin": "USER_DEFINED"},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListRadiusProfiles(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "RADIUS Profiles (2 of 2):") {
		t.Errorf(
			"result should contain header 'RADIUS Profiles (2 of 2):': %s",
			result,
		)
	}
	if !strings.Contains(result, "1. Corp RADIUS") {
		t.Errorf(
			"result should contain '1. Corp RADIUS': %s",
			result,
		)
	}
	if !strings.Contains(result, "2. Guest RADIUS") {
		t.Errorf(
			"result should contain '2. Guest RADIUS': %s",
			result,
		)
	}
}
