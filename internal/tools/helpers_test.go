package tools

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResolveSiteID(t *testing.T) {
	tests := []struct {
		name      string
		explicit  string
		defaultID string
		wantErr   bool
		wantUUID  string
	}{
		{
			name:      "explicit provided",
			explicit:  "550e8400-e29b-41d4-a716-446655440000",
			defaultID: "",
			wantErr:   false,
			wantUUID:  "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:      "falls back to default",
			explicit:  "",
			defaultID: "550e8400-e29b-41d4-a716-446655440000",
			wantErr:   false,
			wantUUID:  "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:      "explicit takes precedence",
			explicit:  "550e8400-e29b-41d4-a716-446655440000",
			defaultID: "660e8400-e29b-41d4-a716-446655440000",
			wantErr:   false,
			wantUUID:  "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:      "neither provided",
			explicit:  "",
			defaultID: "",
			wantErr:   true,
		},
		{
			name:      "invalid UUID",
			explicit:  "not-a-uuid",
			defaultID: "",
			wantErr:   true,
		},
		{
			name:      "invalid default UUID",
			explicit:  "",
			defaultID: "not-a-uuid",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveSiteID(
				tt.explicit,
				tt.defaultID,
			)
			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantUUID != "" &&
				result.String() != tt.wantUUID {
				t.Errorf(
					"UUID = %s, want %s",
					result.String(),
					tt.wantUUID,
				)
			}
		})
	}
}

func TestResolveUUID(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			field:   "testId",
			value:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "empty value",
			field:   "testId",
			value:   "",
			wantErr: true,
		},
		{
			name:    "invalid UUID",
			field:   "testId",
			value:   "not-valid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveUUID(tt.field, tt.value)
			if tt.wantErr && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.wantErr && result.String() == "" {
				t.Error("Expected non-empty UUID")
			}
		})
	}
}

func TestResolveUUIDs(t *testing.T) {
	t.Run("valid UUIDs", func(t *testing.T) {
		ids, err := resolveUUIDs("testIds", []string{
			"550e8400-e29b-41d4-a716-446655440000",
			"660e8400-e29b-41d4-a716-446655440001",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ids) != 2 {
			t.Fatalf("expected 2 UUIDs, got %d", len(ids))
		}
		if ids[0].String() != "550e8400-e29b-41d4-a716-446655440000" {
			t.Errorf("first UUID = %s", ids[0].String())
		}
		if ids[1].String() != "660e8400-e29b-41d4-a716-446655440001" {
			t.Errorf("second UUID = %s", ids[1].String())
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		ids, err := resolveUUIDs("testIds", []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ids) != 0 {
			t.Errorf("expected 0 UUIDs, got %d", len(ids))
		}
	})

	t.Run("invalid UUID in slice", func(t *testing.T) {
		_, err := resolveUUIDs("testIds", []string{
			"550e8400-e29b-41d4-a716-446655440000",
			"not-valid",
		})
		if err == nil {
			t.Error("expected error for invalid UUID")
		}
	})
}

func TestParseArgs(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		var dst struct {
			Name string `json:"name"`
		}
		err := parseArgs(
			json.RawMessage(`{"name": "test"}`),
			&dst,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dst.Name != "test" {
			t.Errorf("Name = %s, want test", dst.Name)
		}
	})

	t.Run("empty args", func(t *testing.T) {
		var dst struct {
			Name string `json:"name"`
		}
		err := parseArgs(json.RawMessage{}, &dst)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dst.Name != "" {
			t.Errorf("Name should be empty, got %s", dst.Name)
		}
	})

	t.Run("nil args", func(t *testing.T) {
		var dst struct {
			Name string `json:"name"`
		}
		err := parseArgs(nil, &dst)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		var dst struct {
			Name string `json:"name"`
		}
		err := parseArgs(
			json.RawMessage(`{invalid}`),
			&dst,
		)
		if err == nil {
			t.Error("expected error for malformed JSON")
		}
	})
}

func TestStripKeys(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		keys    []string
		want    map[string]json.RawMessage
		wantErr bool
	}{
		{
			name:  "strips single key",
			input: `{"siteId":"abc","name":"test"}`,
			keys:  []string{"siteId"},
			want: map[string]json.RawMessage{
				"name": json.RawMessage(`"test"`),
			},
		},
		{
			name:  "strips multiple keys",
			input: `{"siteId":"abc","resourceId":"def","name":"test"}`,
			keys:  []string{"siteId", "resourceId"},
			want: map[string]json.RawMessage{
				"name": json.RawMessage(`"test"`),
			},
		},
		{
			name:  "absent key is a no-op",
			input: `{"name":"test"}`,
			keys:  []string{"siteId"},
			want: map[string]json.RawMessage{
				"name": json.RawMessage(`"test"`),
			},
		},
		{
			name:  "no keys strips nothing",
			input: `{"siteId":"abc","name":"test"}`,
			keys:  nil,
			want: map[string]json.RawMessage{
				"siteId": json.RawMessage(`"abc"`),
				"name":   json.RawMessage(`"test"`),
			},
		},
		{
			name:    "invalid JSON returns error",
			input:   `{invalid}`,
			keys:    []string{"siteId"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stripKeys(
				json.RawMessage(tt.input),
				tt.keys...,
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
			if tt.wantErr {
				return
			}

			var gotMap map[string]json.RawMessage
			if err := json.Unmarshal(got, &gotMap); err != nil {
				t.Fatalf(
					"failed to unmarshal result: %v",
					err,
				)
			}

			if len(gotMap) != len(tt.want) {
				t.Errorf(
					"got %d keys, want %d",
					len(gotMap),
					len(tt.want),
				)
			}
			for k, wantVal := range tt.want {
				gotVal, ok := gotMap[k]
				if !ok {
					t.Errorf("missing key %q in result", k)
					continue
				}
				if string(gotVal) != string(wantVal) {
					t.Errorf(
						"key %q: got %s, want %s",
						k,
						gotVal,
						wantVal,
					)
				}
			}
			for k := range gotMap {
				if _, ok := tt.want[k]; !ok {
					t.Errorf(
						"unexpected key %q in result",
						k,
					)
				}
			}
		})
	}
}

func TestListSchema(t *testing.T) {
	schema := listSchema()

	if schema["type"] != "object" {
		t.Errorf("type = %v, want object", schema["type"])
	}

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("properties should be a map")
	}

	for _, key := range []string{"siteId", "limit", "offset"} {
		if _, ok := props[key]; !ok {
			t.Errorf("properties missing key %q", key)
		}
	}
}

func TestUnexpectedStatusError(t *testing.T) {
	err := unexpectedStatusError(401, []byte("unauthorized"))
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "401") {
		t.Errorf("error should contain status code: %s", msg)
	}
	if !strings.Contains(msg, "unauthorized") {
		t.Errorf("error should contain body: %s", msg)
	}
}
