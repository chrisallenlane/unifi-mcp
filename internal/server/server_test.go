package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

// stubTool implements tools.Tool for testing server dispatch
// without real HTTP calls.
type stubTool struct {
	result string
	err    error
}

func (s *stubTool) Execute(
	_ context.Context,
	_ json.RawMessage,
) (string, error) {
	return s.result, s.err
}

func (s *stubTool) Description() string {
	return "A test tool"
}

func (s *stubTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

func newTestServer(t *testing.T) *Server {
	t.Helper()
	client, err := unifi.NewClientWithResponses(
		"http://localhost",
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return New(client, "default-site")
}

func TestHandleInitialize(t *testing.T) {
	s := newTestServer(t)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.JSONRPC != "2.0" {
		t.Errorf("Response JSONRPC = %s, want 2.0", resp.JSONRPC)
	}

	if resp.ID != 1 {
		t.Errorf("Response ID = %v, want 1", resp.ID)
	}

	if resp.Error != nil {
		t.Errorf("Unexpected error: %+v", resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("Result should not be nil")
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Result should be a map")
	}

	if result["protocolVersion"] != MCPProtocolVersion {
		t.Errorf(
			"Protocol version = %v, want %s",
			result["protocolVersion"],
			MCPProtocolVersion,
		)
	}

	serverInfo, ok := result["serverInfo"].(map[string]string)
	if !ok {
		t.Fatal("serverInfo should be a map")
	}

	if serverInfo["name"] != ServerName {
		t.Errorf(
			"Server name = %s, want %s",
			serverInfo["name"],
			ServerName,
		)
	}

	if serverInfo["version"] != ServerVersion {
		t.Errorf(
			"Server version = %s, want %s",
			serverInfo["version"],
			ServerVersion,
		)
	}

	capabilities, ok := result["capabilities"].(map[string]interface{})
	if !ok {
		t.Fatal("capabilities should be a map")
	}

	if _, ok := capabilities["tools"]; !ok {
		t.Error("capabilities should contain 'tools'")
	}
}

func TestHandleListTools(t *testing.T) {
	s := newTestServer(t)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error != nil {
		t.Errorf("Unexpected error: %+v", resp.Error)
	}

	if resp.Result == nil {
		t.Fatal("Result should not be nil")
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Result should be a map")
	}

	toolList, ok := result["tools"].([]map[string]interface{})
	if !ok {
		t.Fatal("tools should be a slice of maps")
	}

	if len(toolList) == 0 {
		t.Fatal("tools list should not be empty")
	}

	for _, tool := range toolList {
		name, ok := tool["name"].(string)
		if !ok || name == "" {
			t.Errorf("tool missing or empty name: %v", tool)
			continue
		}

		desc, ok := tool["description"].(string)
		if !ok || desc == "" {
			t.Errorf(
				"tool %q missing or empty description",
				name,
			)
		}

		schema, ok := tool["inputSchema"].(map[string]interface{})
		if !ok {
			t.Errorf("tool %q missing inputSchema", name)
			continue
		}

		if schema["type"] != "object" {
			t.Errorf(
				"tool %q inputSchema type = %v, want object",
				name,
				schema["type"],
			)
		}
	}
}

func TestRegisterTools_Count(t *testing.T) {
	s := newTestServer(t)
	if len(s.tools) != 67 {
		t.Errorf(
			"expected 67 registered tools, got %d",
			len(s.tools),
		)
	}
}

func TestHandleCallTool_MissingArguments(t *testing.T) {
	s := newTestServer(t)
	s.tools["test_tool"] = &stubTool{result: "ok"}

	params := map[string]interface{}{
		"name": "test_tool",
	}
	paramsJSON, _ := json.Marshal(params)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      10,
		Method:  "tools/call",
		Params:  paramsJSON,
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("result should be a map")
	}

	content, ok := result["content"].([]map[string]interface{})
	if !ok {
		t.Fatal("content should be a slice of maps")
	}

	if len(content) != 1 {
		t.Fatalf("content length = %d, want 1", len(content))
	}

	if content[0]["text"] != "ok" {
		t.Errorf(
			"content text = %v, want 'ok'",
			content[0]["text"],
		)
	}
}

func TestHandleUnknownMethod(t *testing.T) {
	s := newTestServer(t)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      3,
		Method:  "unknown/method",
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error == nil {
		t.Fatal("Expected error for unknown method")
	}

	if resp.Error.Code != -32601 {
		t.Errorf(
			"Error code = %d, want -32601",
			resp.Error.Code,
		)
	}

	if !strings.Contains(
		resp.Error.Message,
		"unknown/method",
	) {
		t.Errorf(
			"Error message should contain method name, got: %s",
			resp.Error.Message,
		)
	}

	if resp.Result != nil {
		t.Error("Result should be nil for error response")
	}
}

// TestHandleNotification verifies that JSON-RPC notifications
// (requests without an id) receive no response, per JSON-RPC
// 2.0. The MCP handshake sends "notifications/initialized"
// after initialize; strict clients abort the handshake if the
// server replies to it.
func TestHandleNotification(t *testing.T) {
	s := newTestServer(t)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      nil,
		Method:  "notifications/initialized",
	}

	resp := s.handleRequest(context.Background(), req)

	if resp != nil {
		t.Fatalf(
			"Expected nil response for notification, got %+v",
			resp,
		)
	}
}

// TestRun_NotificationNoResponse verifies that Run() writes
// nothing to stdout for a notification.
func TestRun_NotificationNoResponse(t *testing.T) {
	s := newTestServer(t)

	input := `{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}` +
		"\n"
	stdin := strings.NewReader(input)
	var stdout bytes.Buffer

	if err := s.Run(context.Background(), stdin, &stdout); err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	if stdout.Len() != 0 {
		t.Errorf(
			"expected no output for notification, got: %s",
			stdout.String(),
		)
	}
}

func TestHandleCallTool_InvalidTool(t *testing.T) {
	s := newTestServer(t)

	params := map[string]interface{}{
		"name":      "nonexistent_tool",
		"arguments": json.RawMessage(`{}`),
	}
	paramsJSON, _ := json.Marshal(params)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      4,
		Method:  "tools/call",
		Params:  paramsJSON,
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error == nil {
		t.Fatal("Expected error for nonexistent tool")
	}

	if !strings.Contains(
		resp.Error.Message,
		"tool not found",
	) {
		t.Errorf(
			"Error message should mention 'tool not found', got: %s",
			resp.Error.Message,
		)
	}
}

func TestHandleCallTool_MalformedParams(t *testing.T) {
	s := newTestServer(t)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      5,
		Method:  "tools/call",
		Params:  json.RawMessage(`{invalid json}`),
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error == nil {
		t.Fatal("Expected error for malformed params")
	}

	if resp.Error.Code != -32603 {
		t.Errorf(
			"Error code = %d, want -32603",
			resp.Error.Code,
		)
	}
}

func TestHandleCallTool_Success(t *testing.T) {
	s := newTestServer(t)
	s.tools["test_tool"] = &stubTool{
		result: "success output",
	}

	params := map[string]interface{}{
		"name":      "test_tool",
		"arguments": json.RawMessage(`{}`),
	}
	paramsJSON, _ := json.Marshal(params)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      6,
		Method:  "tools/call",
		Params:  paramsJSON,
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error != nil {
		t.Fatalf("Unexpected error: %+v", resp.Error)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Result should be a map")
	}

	content, ok := result["content"].([]map[string]interface{})
	if !ok {
		t.Fatal("content should be a slice of maps")
	}

	if len(content) != 1 {
		t.Fatalf("content length = %d, want 1", len(content))
	}

	if content[0]["type"] != "text" {
		t.Errorf(
			"content type = %v, want text",
			content[0]["type"],
		)
	}

	if content[0]["text"] != "success output" {
		t.Errorf(
			"content text = %v, want 'success output'",
			content[0]["text"],
		)
	}
}

func TestHandleCallTool_ToolError(t *testing.T) {
	s := newTestServer(t)
	s.tools["failing_tool"] = &stubTool{
		err: fmt.Errorf("something broke"),
	}

	params := map[string]interface{}{
		"name":      "failing_tool",
		"arguments": json.RawMessage(`{}`),
	}
	paramsJSON, _ := json.Marshal(params)

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      7,
		Method:  "tools/call",
		Params:  paramsJSON,
	}

	resp := s.handleRequest(context.Background(), req)

	if resp.Error == nil {
		t.Fatal("Expected error for failing tool")
	}

	if resp.Error.Code != -32603 {
		t.Errorf(
			"Error code = %d, want -32603",
			resp.Error.Code,
		)
	}

	if !strings.Contains(
		resp.Error.Message,
		"something broke",
	) {
		t.Errorf(
			"Error should contain original error, got: %s",
			resp.Error.Message,
		)
	}
}

func TestRun(t *testing.T) {
	s := newTestServer(t)
	s.tools["test_tool"] = &stubTool{
		result: "hello world",
	}

	reqJSON := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"test_tool","arguments":{}}}`
	stdin := strings.NewReader(reqJSON + "\n")
	var stdout bytes.Buffer

	err := s.Run(context.Background(), stdin, &stdout)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	var resp JSONRPCResponse
	if err := json.NewDecoder(&stdout).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("result should be a map")
	}

	content, ok := result["content"].([]interface{})
	if !ok {
		t.Fatal("content should be an array")
	}

	if len(content) == 0 {
		t.Fatal("content should not be empty")
	}

	first, ok := content[0].(map[string]interface{})
	if !ok {
		t.Fatal("content item should be a map")
	}

	if first["text"] != "hello world" {
		t.Errorf(
			"text = %v, want 'hello world'",
			first["text"],
		)
	}
}

func TestRun_MalformedJSON(t *testing.T) {
	s := newTestServer(t)

	stdin := strings.NewReader("{invalid json}\n")
	var stdout bytes.Buffer

	err := s.Run(context.Background(), stdin, &stdout)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	var resp JSONRPCResponse
	if err := json.NewDecoder(&stdout).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error response")
	}

	if resp.Error.Code != -32700 {
		t.Errorf(
			"error code = %d, want -32700",
			resp.Error.Code,
		)
	}
}

func TestRun_MultipleRequests(t *testing.T) {
	s := newTestServer(t)
	s.tools["test_tool"] = &stubTool{
		result: "hello",
	}

	requests := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"test_tool","arguments":{}}}`,
	}, "\n") + "\n"

	stdin := strings.NewReader(requests)
	var stdout bytes.Buffer

	err := s.Run(context.Background(), stdin, &stdout)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	decoder := json.NewDecoder(&stdout)

	var resp1, resp2 JSONRPCResponse
	if err := decoder.Decode(&resp1); err != nil {
		t.Fatalf("failed to decode first response: %v", err)
	}
	if err := decoder.Decode(&resp2); err != nil {
		t.Fatalf("failed to decode second response: %v", err)
	}

	if resp1.Error != nil {
		t.Errorf(
			"first response has error: %+v",
			resp1.Error,
		)
	}
	if resp2.Error != nil {
		t.Errorf(
			"second response has error: %+v",
			resp2.Error,
		)
	}
}

type failWriter struct{}

func (f failWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("broken pipe")
}

func TestRun_EncodeFail(t *testing.T) {
	s := newTestServer(t)

	reqJSON := `{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"
	stdin := strings.NewReader(reqJSON)

	err := s.Run(context.Background(), stdin, failWriter{})
	if err == nil {
		t.Fatal("expected error when stdout write fails")
	}

	if !strings.Contains(err.Error(), "broken pipe") {
		t.Errorf("error should contain cause: %v", err)
	}
}

func TestRun_MalformedJSON_EncodeFail(t *testing.T) {
	s := newTestServer(t)

	stdin := strings.NewReader("{invalid json}\n")

	err := s.Run(context.Background(), stdin, failWriter{})
	if err != nil {
		t.Fatalf(
			"Run should not return error for malformed JSON "+
				"encode failure, got: %v",
			err,
		)
	}
}

func TestRun_EOF(t *testing.T) {
	s := newTestServer(t)

	stdin := strings.NewReader("")
	var stdout bytes.Buffer

	err := s.Run(context.Background(), stdin, &stdout)
	if err != nil {
		t.Fatalf(
			"Run should return nil on EOF, got: %v",
			err,
		)
	}

	if stdout.Len() != 0 {
		t.Errorf(
			"expected no output on empty input, got: %s",
			stdout.String(),
		)
	}
}
