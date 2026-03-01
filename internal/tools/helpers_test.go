package tools

import (
	"testing"
)

func TestResolveSiteID(t *testing.T) {
	tests := []struct {
		name      string
		explicit  string
		defaultID string
		wantErr   bool
	}{
		{
			name:      "explicit provided",
			explicit:  "550e8400-e29b-41d4-a716-446655440000",
			defaultID: "",
			wantErr:   false,
		},
		{
			name:      "falls back to default",
			explicit:  "",
			defaultID: "550e8400-e29b-41d4-a716-446655440000",
			wantErr:   false,
		},
		{
			name:      "explicit takes precedence",
			explicit:  "550e8400-e29b-41d4-a716-446655440000",
			defaultID: "660e8400-e29b-41d4-a716-446655440000",
			wantErr:   false,
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
			if !tt.wantErr && result.String() == "" {
				t.Error("Expected non-empty UUID")
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
