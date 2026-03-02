package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestFormatDNSRecordDetails(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{
			name:   "A record",
			input:  `{"ipv4Address": "192.168.1.1", "ttlSeconds": 300}`,
			expect: []string{"Address: 192.168.1.1", "TTL: 300"},
		},
		{
			name:   "AAAA record",
			input:  `{"ipv6Address": "::1"}`,
			expect: []string{"Address: ::1"},
		},
		{
			name:   "CNAME record",
			input:  `{"targetDomain": "target.local"}`,
			expect: []string{"Target: target.local"},
		},
		{
			name:  "MX record",
			input: `{"mailServerDomain": "mail.local", "priority": 10}`,
			expect: []string{
				"Mail Server: mail.local",
				"Priority: 10",
			},
		},
		{
			name:  "SRV record",
			input: `{"serverDomain": "srv.local", "service": "_http", "protocol": "_tcp", "priority": 10, "weight": 20, "port": 80}`,
			expect: []string{
				"Server: srv.local",
				"Service: _http",
				"Protocol: _tcp",
				"Priority: 10",
				"Weight: 20",
				"Port: 80",
			},
		},
		{
			name:   "TXT record",
			input:  `{"text": "v=spf1"}`,
			expect: []string{"Text: v=spf1"},
		},
		{
			name:   "FORWARD_DOMAIN",
			input:  `{"ipAddress": "8.8.8.8"}`,
			expect: []string{"Forward To: 8.8.8.8"},
		},
		{
			name:   "invalid JSON",
			input:  `{invalid}`,
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDNSRecordDetails(
				json.RawMessage(tt.input),
			)
			if tt.expect == nil {
				if result != "" {
					t.Errorf(
						"expected empty result, got: %s",
						result,
					)
				}
				return
			}
			for _, s := range tt.expect {
				if !strings.Contains(result, s) {
					t.Errorf(
						"expected %q in result: %s",
						s,
						result,
					)
				}
			}
		})
	}
}

func TestListDNSPolicies_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(paginatedResponse(
				map[string]interface{}{
					"id":       "aaa00000-0000-0000-0000-000000000001",
					"type":     "A_RECORD",
					"domain":   "nas.local",
					"enabled":  true,
					"metadata": map[string]string{"origin": "USER_DEFINED"},
				},
				map[string]interface{}{
					"id":       "aaa00000-0000-0000-0000-000000000002",
					"type":     "CNAME_RECORD",
					"domain":   "wiki.local",
					"enabled":  false,
					"metadata": map[string]string{"origin": "USER_DEFINED"},
				},
			))
		}),
	)
	defer srv.Close()

	tool := NewListDNSPolicies(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "A_RECORD") {
		t.Errorf(
			"result should contain 'A_RECORD': %s",
			result,
		)
	}
	if !strings.Contains(result, "nas.local") {
		t.Errorf(
			"result should contain 'nas.local': %s",
			result,
		)
	}
}

func TestListDNSPolicies_Execute_Empty(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emptyPaginatedResponse())
		}),
	)
	defer srv.Close()

	tool := NewListDNSPolicies(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "No DNS policies found." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestListDNSPolicies_Execute_NoSiteID(t *testing.T) {
	tool := &ListDNSPolicies{}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{}`),
	)
	if err == nil {
		t.Fatal("expected error when no site ID provided")
	}
}

func TestGetDNSPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       "aaa00000-0000-0000-0000-000000000001",
				"type":     "A_RECORD",
				"domain":   "nas.local",
				"enabled":  true,
				"metadata": map[string]string{"origin": "USER_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	tool := NewGetDNSPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "A_RECORD") {
		t.Errorf(
			"result should contain 'A_RECORD': %s",
			result,
		)
	}
	if !strings.Contains(result, "nas.local") {
		t.Errorf(
			"result should contain 'nas.local': %s",
			result,
		)
	}
}

func TestGetDNSPolicy_Execute_InvalidUUID(t *testing.T) {
	tool := &GetDNSPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"dnsPolicyId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestGetDNSPolicy_InputSchema(t *testing.T) {
	tool := &GetDNSPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "dnsPolicyId" {
			found = true
		}
	}
	if !found {
		t.Error("dnsPolicyId should be required")
	}
}

func TestCreateDNSPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       "aaa00000-0000-0000-0000-000000000001",
				"type":     "A_RECORD",
				"domain":   "nas.local",
				"enabled":  true,
				"metadata": map[string]string{"origin": "USER_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	tool := NewCreateDNSPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"type": "A_RECORD", "enabled": true, "domain": "nas.local", "address": "192.168.1.50"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "DNS policy created") {
		t.Errorf(
			"result should contain 'DNS policy created': %s",
			result,
		)
	}
	if !strings.Contains(result, "A_RECORD") {
		t.Errorf(
			"result should contain 'A_RECORD': %s",
			result,
		)
	}
}

func TestCreateDNSPolicy_InputSchema(t *testing.T) {
	tool := &CreateDNSPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	foundType := false
	foundEnabled := false
	for _, r := range required {
		if r == "type" {
			foundType = true
		}
		if r == "enabled" {
			foundEnabled = true
		}
	}
	if !foundType {
		t.Error("type should be required")
	}
	if !foundEnabled {
		t.Error("enabled should be required")
	}
}

func TestUpdateDNSPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       "aaa00000-0000-0000-0000-000000000001",
				"type":     "A_RECORD",
				"domain":   "nas.local",
				"enabled":  false,
				"metadata": map[string]string{"origin": "USER_DEFINED"},
			})
		}),
	)
	defer srv.Close()

	tool := NewUpdateDNSPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "aaa00000-0000-0000-0000-000000000001", "type": "A_RECORD", "enabled": false}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "DNS policy updated") {
		t.Errorf(
			"result should contain 'DNS policy updated': %s",
			result,
		)
	}
}

func TestUpdateDNSPolicy_Execute_InvalidUUID(t *testing.T) {
	tool := &UpdateDNSPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "not-valid", "type": "A_RECORD", "enabled": true}`,
		),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestUpdateDNSPolicy_InputSchema(t *testing.T) {
	tool := &UpdateDNSPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "dnsPolicyId" {
			found = true
		}
	}
	if !found {
		t.Error("dnsPolicyId should be required")
	}
}

func TestDeleteDNSPolicy_Execute(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	tool := NewDeleteDNSPolicy(client, testSiteID)
	result, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "DNS policy deleted." {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestDeleteDNSPolicy_Execute_InvalidUUID(t *testing.T) {
	tool := &DeleteDNSPolicy{baseTool{defaultSiteID: testSiteID}}
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(`{"dnsPolicyId": "not-valid"}`),
	)
	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestDeleteDNSPolicy_InputSchema(t *testing.T) {
	tool := &DeleteDNSPolicy{}
	schema := tool.InputSchema()
	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("required should be a string slice")
	}
	found := false
	for _, r := range required {
		if r == "dnsPolicyId" {
			found = true
		}
	}
	if !found {
		t.Error("dnsPolicyId should be required")
	}
}

func TestListDNSPolicies_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewListDNSPolicies(client, testSiteID)
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

func TestGetDNSPolicy_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewGetDNSPolicy(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestCreateDNSPolicy_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewCreateDNSPolicy(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"type": "A_RECORD", "enabled": true}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestUpdateDNSPolicy_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewUpdateDNSPolicy(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "aaa00000-0000-0000-0000-000000000001", "type": "A_RECORD", "enabled": true}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}

func TestDeleteDNSPolicy_Execute_APIError(t *testing.T) {
	client, srv := testClient(t,
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
		}),
	)
	defer srv.Close()

	tool := NewDeleteDNSPolicy(client, testSiteID)
	_, err := tool.Execute(
		context.Background(),
		json.RawMessage(
			`{"dnsPolicyId": "aaa00000-0000-0000-0000-000000000001"}`,
		),
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code: %v", err)
	}
}
