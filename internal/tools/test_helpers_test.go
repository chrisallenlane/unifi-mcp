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
