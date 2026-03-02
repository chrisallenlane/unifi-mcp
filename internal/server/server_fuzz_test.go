package server

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

func FuzzRun(f *testing.F) {
	f.Add([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize"}`))
	f.Add([]byte(`{invalid json}`))
	f.Add(
		[]byte(
			`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"x","arguments":{}}}`,
		),
	)
	f.Add([]byte(`null`))
	f.Add([]byte(``))
	f.Add([]byte(`{}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		client, err := unifi.NewClientWithResponses(
			"http://localhost",
		)
		if err != nil {
			t.Skip("failed to create client")
		}

		s := New(client, "default-site")

		stdin := strings.NewReader(string(data) + "\n")
		var stdout bytes.Buffer

		// Must not panic. Errors are acceptable.
		_ = s.Run(context.Background(), stdin, &stdout)
	})
}

func FuzzHandleCallTool(f *testing.F) {
	f.Add([]byte(`{"name":"test_tool","arguments":{}}`))
	f.Add([]byte(`{"name":"nonexistent","arguments":{}}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`{"name":"test_tool","arguments":null}`))
	f.Add([]byte(`{"name":"test_tool"}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		client, err := unifi.NewClientWithResponses(
			"http://localhost",
		)
		if err != nil {
			t.Skip("failed to create client")
		}

		s := New(client, "default-site")
		s.tools["test_tool"] = &stubTool{
			result: "ok",
		}

		req := &JSONRPCRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "tools/call",
			Params:  json.RawMessage(data),
		}

		// Must not panic. Error responses are acceptable.
		resp := s.handleRequest(context.Background(), req)
		if resp == nil {
			t.Fatal("response should never be nil")
		}
	})
}
