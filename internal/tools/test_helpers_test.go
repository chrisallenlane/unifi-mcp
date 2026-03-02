package tools

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

const testSiteID = "550e8400-e29b-41d4-a716-446655440000"

func testClient(
	t *testing.T,
	handler http.Handler,
) (*unifi.ClientWithResponses, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client, srv
}

// paginatedResponse wraps data items in the standard UniFi
// paginated response envelope used by list endpoints.
func paginatedResponse(
	data ...map[string]interface{},
) map[string]interface{} {
	if data == nil {
		data = []map[string]interface{}{}
	}
	return map[string]interface{}{
		"data":       data,
		"count":      len(data),
		"limit":      25,
		"offset":     0,
		"totalCount": len(data),
	}
}

// emptyPaginatedResponse returns an empty paginated response
// with an untyped empty slice matching the pattern used by the
// UniFi API for empty results.
func emptyPaginatedResponse() map[string]interface{} {
	return map[string]interface{}{
		"data":       []interface{}{},
		"count":      0,
		"limit":      25,
		"offset":     0,
		"totalCount": 0,
	}
}

// requireContains asserts that a string slice contains the
// given value.
func requireContains(
	t *testing.T,
	slice []string,
	value string,
) {
	t.Helper()
	for _, s := range slice {
		if s == value {
			return
		}
	}
	t.Errorf("%q should be in required fields", value)
}
