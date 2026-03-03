package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListNetworks_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
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
				map[string]interface{}{
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
			))
		}),
	)
	defer srv.Close()

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
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

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
	if !strings.Contains(err.Error(), "siteId") {
		t.Errorf("error should mention siteId: %v", err)
	}
}

func TestGetNetwork_Execute(t *testing.T) {
	client, srv := testClient(t,
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
	tool := &GetNetwork{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"networkId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "networkId") {
		t.Errorf("error should mention networkId: %v", err)
	}
}

func TestGetNetwork_InputSchema(t *testing.T) {
	tool := &GetNetwork{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "networkId")
}

func TestCreateNetwork_Execute(t *testing.T) {
	client, srv := testClient(t,
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

func TestCreateNetwork_InputSchema(t *testing.T) {
	tool := &CreateNetwork{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "name")
	requireContains(t, required, "management")
}

func TestUpdateNetwork_Execute(t *testing.T) {
	client, srv := testClient(t,
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
	tool := &UpdateNetwork{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "not-valid", "name": "x", "enabled": true, "management": "GATEWAY", "vlanId": 2}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "networkId") {
		t.Errorf("error should mention networkId: %v", err)
	}
}

func TestUpdateNetwork_InputSchema(t *testing.T) {
	tool := &UpdateNetwork{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "networkId")
}

func TestDeleteNetwork_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

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
	tool := &DeleteNetwork{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"networkId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "networkId") {
		t.Errorf("error should mention networkId: %v", err)
	}
}

func TestDeleteNetwork_InputSchema(t *testing.T) {
	tool := &DeleteNetwork{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "networkId")
}

func TestGetNetworkReferences_Execute(t *testing.T) {
	client, srv := testClient(t,
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
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"referenceResources": []interface{}{},
			})
		}),
	)
	defer srv.Close()

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
	tool := &GetNetworkReferences{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"networkId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
	if !strings.Contains(err.Error(), "networkId") {
		t.Errorf("error should mention networkId: %v", err)
	}
}

func TestGetNetworkReferences_InputSchema(t *testing.T) {
	tool := &GetNetworkReferences{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	requireContains(t, required, "networkId")
}

// --- optional field branches ---

func TestGetNetwork_Execute_WithDhcpGuarding(t *testing.T) {
	client, srv := testClient(t,
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

	if !strings.Contains(result, "DHCP Guarding:") {
		t.Errorf(
			"result should contain 'DHCP Guarding:': %s",
			result,
		)
	}
	if !strings.Contains(result, "Trusted Servers: 192.168.1.1") {
		t.Errorf(
			"result should contain trusted server IP: %s",
			result,
		)
	}
}

func TestGetNetwork_Execute_WithEmptyDhcpGuarding(t *testing.T) {
	client, srv := testClient(t,
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
					"trustedDhcpServerIpAddresses": []string{},
				},
			})
		}),
	)
	defer srv.Close()

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

	if !strings.Contains(result, "DHCP Guarding:") {
		t.Errorf(
			"result should contain 'DHCP Guarding:': %s",
			result,
		)
	}
	if !strings.Contains(result, "(none)") {
		t.Errorf(
			"result should contain '(none)' for empty trusted servers: %s",
			result,
		)
	}
}

// --- API error tests ---

func TestListNetworks_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewListNetworks(client, testSiteID)
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

func TestGetNetwork_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewGetNetwork(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
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

func TestCreateNetwork_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewCreateNetwork(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Guest", "enabled": true, "management": "GATEWAY", "vlanId": 200}`,
		),
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

func TestUpdateNetwork_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewUpdateNetwork(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001", "name": "x", "enabled": true, "management": "GATEWAY", "vlanId": 2}`,
		),
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

func TestDeleteNetwork_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewDeleteNetwork(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
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

// --- formatting and coverage tests ---

func TestListNetworks_Description(t *testing.T) {
	tool := &ListNetworks{}
	d := tool.Description()
	if d == "" {
		t.Fatal("description should not be empty")
	}
	if !strings.Contains(d, "etwork") {
		t.Errorf(
			"description should mention networks: %s",
			d,
		)
	}
}

func TestListNetworks_Execute_InvalidJSON(t *testing.T) {
	tool := &ListNetworks{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{invalid`),
	)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestListNetworks_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewListNetworks(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to list networks") {
		t.Errorf(
			"error should contain 'failed to list networks': %v",
			err,
		)
	}
}

func TestListNetworks_Execute_Formatting(t *testing.T) {
	client, srv := testClient(t,
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
				"totalCount": 5,
			})
		}),
	)
	defer srv.Close()

	tool := NewListNetworks(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// verify header uses totalCount, not len(data)
	if !strings.Contains(result, "Networks (2 of 5):") {
		t.Errorf(
			"result should contain 'Networks (2 of 5):': %s",
			result,
		)
	}

	// verify numbering
	if !strings.Contains(result, "1. Default") {
		t.Errorf(
			"result should contain '1. Default': %s",
			result,
		)
	}
	if !strings.Contains(result, "2. IoT") {
		t.Errorf(
			"result should contain '2. IoT': %s",
			result,
		)
	}
}

func TestGetNetwork_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewGetNetwork(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "failed to get network") {
		t.Errorf(
			"error should contain 'failed to get network': %v",
			err,
		)
	}
}

func TestCreateNetwork_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewCreateNetwork(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"name": "Guest", "enabled": true, "management": "GATEWAY", "vlanId": 200}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to create network",
	) {
		t.Errorf(
			"error should contain 'failed to create network': %v",
			err,
		)
	}
}

func TestUpdateNetwork_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewUpdateNetwork(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001", "name": "x", "enabled": true, "management": "GATEWAY", "vlanId": 2}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to update network",
	) {
		t.Errorf(
			"error should contain 'failed to update network': %v",
			err,
		)
	}
}

func TestDeleteNetwork_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewDeleteNetwork(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to delete network",
	) {
		t.Errorf(
			"error should contain 'failed to delete network': %v",
			err,
		)
	}
}

func TestGetNetworkReferences_Execute_NetworkError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(
			func(_ http.ResponseWriter, _ *http.Request) {},
		),
	)
	srv.Close()

	tool := NewGetNetworkReferences(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for network failure")
	}
	if !strings.Contains(
		err.Error(),
		"failed to get network references",
	) {
		t.Errorf(
			"error should contain 'failed to get network references': %v",
			err,
		)
	}
}

func TestGetNetworkReferences_Execute_Formatting(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"referenceResources": []map[string]interface{}{
					{
						"resourceType":   "CLIENT",
						"referenceCount": 3,
					},
					{
						"resourceType":   "WIFI_BROADCAST",
						"referenceCount": 1,
						"references": []map[string]interface{}{
							{
								"referenceId": "bbb00000-0000-0000-0000-000000000001",
							},
						},
					},
				},
			})
		}),
	)
	defer srv.Close()

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

	// verify count header
	if !strings.Contains(
		result,
		"Network references (2 resource types):",
	) {
		t.Errorf(
			"result should contain header with count: %s",
			result,
		)
	}

	// verify reference count format
	if !strings.Contains(result, "CLIENT: 3 reference(s)") {
		t.Errorf(
			"result should contain 'CLIENT: 3 reference(s)': %s",
			result,
		)
	}

	// verify reference IDs listed
	if !strings.Contains(
		result,
		"bbb00000-0000-0000-0000-000000000001",
	) {
		t.Errorf(
			"result should contain reference ID: %s",
			result,
		)
	}
}

func TestGetNetworkReferences_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer srv.Close()

	tool := NewGetNetworkReferences(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"networkId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
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
