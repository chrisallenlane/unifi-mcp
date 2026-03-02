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
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":       "ddd00000-0000-0000-0000-000000000001",
						"name":     "Corp RADIUS",
						"metadata": map[string]string{"origin": "USER_DEFINED"},
					},
				},
				"count":      1,
				"limit":      25,
				"offset":     0,
				"totalCount": 1,
			})
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

func TestListRadiusProfiles_Description(t *testing.T) {
	tool := &ListRadiusProfiles{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListRadiusProfiles_InputSchema(t *testing.T) {
	tool := &ListRadiusProfiles{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}

// --- list_device_tags tests ---

func TestListDeviceTags_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":   "eee00000-0000-0000-0000-000000000001",
						"name": "Floor 1 APs",
						"deviceIds": []string{
							"fff00000-0000-0000-0000-000000000001",
							"fff00000-0000-0000-0000-000000000002",
						},
						"metadata": map[string]string{"origin": "USER_DEFINED"},
					},
				},
				"count":      1,
				"limit":      25,
				"offset":     0,
				"totalCount": 1,
			})
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

func TestListDeviceTags_Description(t *testing.T) {
	tool := &ListDeviceTags{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListDeviceTags_InputSchema(t *testing.T) {
	tool := &ListDeviceTags{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}

// --- list_dpi_categories tests ---

func TestListDpiCategories_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"id": 1, "name": "Streaming Media"},
					{"id": 2, "name": "Social Networking"},
				},
				"count":      2,
				"limit":      25,
				"offset":     0,
				"totalCount": 2,
			})
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

func TestListDpiCategories_Description(t *testing.T) {
	tool := &ListDpiCategories{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListDpiCategories_InputSchema(t *testing.T) {
	tool := &ListDpiCategories{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}

// --- list_dpi_applications tests ---

func TestListDpiApplications_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"id": 100, "name": "Netflix"},
					{"id": 101, "name": "YouTube"},
				},
				"count":      2,
				"limit":      25,
				"offset":     0,
				"totalCount": 2,
			})
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

func TestListDpiApplications_Description(t *testing.T) {
	tool := &ListDpiApplications{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListDpiApplications_InputSchema(t *testing.T) {
	tool := &ListDpiApplications{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}

// --- list_countries tests ---

func TestListCountries_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"code": "US", "name": "United States"},
					{"code": "DE", "name": "Germany"},
					{"code": "JP", "name": "Japan"},
				},
				"count":      3,
				"limit":      25,
				"offset":     0,
				"totalCount": 3,
			})
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

func TestListCountries_Description(t *testing.T) {
	tool := &ListCountries{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListCountries_InputSchema(t *testing.T) {
	tool := &ListCountries{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}
