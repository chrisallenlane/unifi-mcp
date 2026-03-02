package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

func TestListNetworks_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":         "aaa00000-0000-0000-0000-000000000001",
						"name":       "Default",
						"vlanId":     1,
						"management": "GATEWAY",
						"enabled":    true,
						"default":    true,
						"metadata": map[string]string{
							"origin": "SYSTEM_DEFINED",
						},
					},
					{
						"id":         "aaa00000-0000-0000-0000-000000000002",
						"name":       "IoT",
						"vlanId":     100,
						"management": "GATEWAY",
						"enabled":    true,
						"default":    false,
						"metadata": map[string]string{
							"origin": "USER_DEFINED",
						},
					},
				},
				"count":      2,
				"limit":      25,
				"offset":     0,
				"totalCount": 2,
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewListNetworks(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Default") {
		t.Errorf(
			"result should contain 'Default': %s",
			result,
		)
	}
	if !strings.Contains(result, "IoT") {
		t.Errorf(
			"result should contain 'IoT': %s",
			result,
		)
	}
	if !strings.Contains(result, "GATEWAY") {
		t.Errorf(
			"result should contain 'GATEWAY': %s",
			result,
		)
	}
}

func TestListNetworks_Execute_Empty(t *testing.T) {
	srv := httptest.NewServer(
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

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewListNetworks(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No networks found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListNetworks_Execute_NoSiteID(t *testing.T) {
	tool := &ListNetworks{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID provided")
	}
}

func TestListNetworks_Description(t *testing.T) {
	tool := &ListNetworks{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestListNetworks_InputSchema(t *testing.T) {
	tool := &ListNetworks{}
	schema := tool.InputSchema()
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}
}

func TestGetNetwork_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         "aaa00000-0000-0000-0000-000000000001",
				"name":       "IoT",
				"vlanId":     100,
				"management": "GATEWAY",
				"enabled":    true,
				"default":    false,
				"metadata":   map[string]string{"origin": "USER_DEFINED"},
				"dhcpGuarding": map[string]interface{}{
					"trustedDhcpServerIpAddresses": []string{
						"192.168.1.1",
					},
				},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewGetNetwork(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "IoT") {
		t.Errorf(
			"result should contain 'IoT': %s",
			result,
		)
	}
	if !strings.Contains(result, "VLAN ID: 100") {
		t.Errorf(
			"result should contain 'VLAN ID: 100': %s",
			result,
		)
	}
	if !strings.Contains(result, "192.168.1.1") {
		t.Errorf(
			"result should contain DHCP server IP: %s",
			result,
		)
	}
}

func TestGetNetwork_Execute_InvalidUUID(t *testing.T) {
	tool := &GetNetwork{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"networkId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestGetNetwork_Description(t *testing.T) {
	tool := &GetNetwork{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestGetNetwork_InputSchema(t *testing.T) {
	tool := &GetNetwork{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "networkId" {
			found = true
		}
	}
	if !found {
		t.Error("networkId should be required")
	}
}

func TestCreateNetwork_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         "aaa00000-0000-0000-0000-000000000001",
				"name":       "Guest",
				"vlanId":     200,
				"management": "GATEWAY",
				"enabled":    true,
				"default":    false,
				"metadata":   map[string]string{"origin": "USER_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewCreateNetwork(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Guest", "enabled": true, "management": "GATEWAY", "vlanId": 200}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Network created") {
		t.Errorf(
			"result should contain 'Network created': %s",
			result,
		)
	}
	if !strings.Contains(result, "Guest") {
		t.Errorf(
			"result should contain 'Guest': %s",
			result,
		)
	}
}

func TestCreateNetwork_Description(t *testing.T) {
	tool := &CreateNetwork{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestCreateNetwork_InputSchema(t *testing.T) {
	tool := &CreateNetwork{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	foundName := false
	foundMgmt := false
	for _, r := range required {
		if r == "name" {
			foundName = true
		}
		if r == "management" {
			foundMgmt = true
		}
	}
	if !foundName {
		t.Error("name should be required")
	}
	if !foundMgmt {
		t.Error("management should be required")
	}
}

func TestUpdateNetwork_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         "aaa00000-0000-0000-0000-000000000001",
				"name":       "Guest Updated",
				"vlanId":     200,
				"management": "GATEWAY",
				"enabled":    false,
				"default":    false,
				"metadata":   map[string]string{"origin": "USER_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewUpdateNetwork(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001", "name": "Guest Updated", "enabled": false, "management": "GATEWAY", "vlanId": 200}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "Network updated") {
		t.Errorf(
			"result should contain 'Network updated': %s",
			result,
		)
	}
}

func TestUpdateNetwork_Execute_InvalidUUID(t *testing.T) {
	tool := &UpdateNetwork{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "not-valid", "name": "x", "enabled": true, "management": "GATEWAY", "vlanId": 2}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestUpdateNetwork_Description(t *testing.T) {
	tool := &UpdateNetwork{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestUpdateNetwork_InputSchema(t *testing.T) {
	tool := &UpdateNetwork{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "networkId" {
			found = true
		}
	}
	if !found {
		t.Error("networkId should be required")
	}
}

func TestDeleteNetwork_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewDeleteNetwork(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Network deleted." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestDeleteNetwork_Execute_InvalidUUID(t *testing.T) {
	tool := &DeleteNetwork{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"networkId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestDeleteNetwork_Description(t *testing.T) {
	tool := &DeleteNetwork{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestDeleteNetwork_InputSchema(t *testing.T) {
	tool := &DeleteNetwork{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "networkId" {
			found = true
		}
	}
	if !found {
		t.Error("networkId should be required")
	}
}

func TestGetNetworkReferences_Execute(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"referenceResources": []map[string]interface{}{
					{
						"resourceType":   "CLIENT",
						"referenceCount": 5,
					},
					{
						"resourceType":   "WIFI_BROADCAST",
						"referenceCount": 2,
						"references": []map[string]interface{}{
							{
								"referenceId": "bbb00000-0000-0000-0000-000000000001",
							},
							{
								"referenceId": "bbb00000-0000-0000-0000-000000000002",
							},
						},
					},
				},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewGetNetworkReferences(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "CLIENT") {
		t.Errorf(
			"result should contain 'CLIENT': %s",
			result,
		)
	}
	if !strings.Contains(result, "WIFI_BROADCAST") {
		t.Errorf(
			"result should contain 'WIFI_BROADCAST': %s",
			result,
		)
	}
	if !strings.Contains(result, "5 reference(s)") {
		t.Errorf(
			"result should contain reference count: %s",
			result,
		)
	}
}

func TestGetNetworkReferences_Execute_Empty(t *testing.T) {
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"referenceResources": []interface{}{},
			})
		}),
	)
	defer srv.Close()

	client, err := unifi.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	tool := NewGetNetworkReferences(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No resources reference this network." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestGetNetworkReferences_Execute_InvalidUUID(t *testing.T) {
	tool := &GetNetworkReferences{defaultSiteID: testSiteID}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"networkId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestGetNetworkReferences_Description(t *testing.T) {
	tool := &GetNetworkReferences{}
	if tool.Description() == "" {
		t.Error("description should not be empty")
	}
}

func TestGetNetworkReferences_InputSchema(t *testing.T) {
	tool := &GetNetworkReferences{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "networkId" {
			found = true
		}
	}
	if !found {
		t.Error("networkId should be required")
	}
}
