package tools

import (
	"encoding/json"
	"testing"
)

func FuzzParseArgs(f *testing.F) {
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"siteId":"550e8400-e29b-41d4-a716-446655440000"}`))
	f.Add([]byte(`null`))
	f.Add([]byte(`{invalid}`))
	f.Add([]byte(`"string"`))
	f.Add([]byte(`[]`))
	f.Add([]byte(`0`))
	f.Add([]byte(``))

	f.Fuzz(func(_ *testing.T, data []byte) {
		var params struct {
			SiteID string `json:"siteId"`
			ID     string `json:"id"`
			Name   string `json:"name"`
		}

		// Must not panic. Errors are acceptable.
		_ = parseArgs(json.RawMessage(data), &params)
	})
}

func FuzzResolveUUID(f *testing.F) {
	f.Add("550e8400-e29b-41d4-a716-446655440000")
	f.Add("")
	f.Add("not-a-uuid")
	f.Add("550e8400e29b41d4a716446655440000")
	f.Add("\x00\x00\x00\x00")

	f.Fuzz(func(_ *testing.T, data string) {
		// Must not panic. Errors are acceptable.
		_, _ = resolveUUID("testField", data)
	})
}
