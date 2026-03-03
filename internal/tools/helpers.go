// Package tools implements MCP tool definitions for the UniFi MCP server.
package tools

import (
	"encoding/json"
	"fmt"

	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
	"github.com/google/uuid"
)

// baseTool holds the common fields shared by all tool implementations.
type baseTool struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
}

// resolveSiteID resolves a site ID from an explicit parameter or the
// default. Returns an error if neither is provided.
func resolveSiteID(
	explicit string,
	defaultID string,
) (uuid.UUID, error) {
	id := explicit
	if id == "" {
		id = defaultID
	}
	if id == "" {
		return uuid.UUID{}, fmt.Errorf(
			"siteId is required (provide it as a parameter or set UNIFI_SITE_ID)",
		)
	}
	return resolveUUID("siteId", id)
}

// resolveUUID parses and validates a UUID string.
func resolveUUID(
	name string,
	value string,
) (uuid.UUID, error) {
	if value == "" {
		return uuid.UUID{}, fmt.Errorf(
			"%s is required",
			name,
		)
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf(
			"invalid %s: %w",
			name,
			err,
		)
	}

	return parsed, nil
}

// resolveUUIDs parses a slice of UUID strings.
func resolveUUIDs(
	name string,
	values []string,
) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, len(values))
	for i, v := range values {
		parsed, err := resolveUUID(name, v)
		if err != nil {
			return nil, err
		}
		ids[i] = parsed
	}
	return ids, nil
}

// unexpectedStatusError returns a formatted error for unexpected
// HTTP status codes.
func unexpectedStatusError(
	statusCode int,
	body []byte,
) error {
	return fmt.Errorf(
		"unexpected status %d: %s",
		statusCode,
		string(body),
	)
}

// stripKeys removes the named keys from a JSON object and returns the
// re-encoded bytes. It is used by create/update tools to scrub
// MCP-only parameters (siteId, resource ID fields) from the request
// body before forwarding it to the UniFi API.
func stripKeys(
	args json.RawMessage,
	keys ...string,
) ([]byte, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(args, &m); err != nil {
		return nil, fmt.Errorf(
			"failed to parse arguments: %w",
			err,
		)
	}
	for _, k := range keys {
		delete(m, k)
	}
	return json.Marshal(m)
}

// parseArgs unmarshals JSON arguments into a destination struct.
func parseArgs(
	args json.RawMessage,
	dst interface{},
) error {
	if len(args) > 0 {
		if err := json.Unmarshal(args, dst); err != nil {
			return fmt.Errorf(
				"failed to parse arguments: %w",
				err,
			)
		}
	}
	return nil
}

// siteIDSchema returns the standard JSON schema for the siteId
// parameter.
func siteIDSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": "Site UUID (uses UNIFI_SITE_ID if not provided)",
	}
}

// paginationSchema returns the standard JSON schema properties
// for limit and offset parameters.
func paginationSchema() map[string]interface{} {
	return map[string]interface{}{
		"limit": map[string]interface{}{
			"type":        "integer",
			"description": "Maximum number of items to return",
		},
		"offset": map[string]interface{}{
			"type":        "integer",
			"description": "Number of items to skip",
		},
	}
}

// siteAndIDSchema returns the standard JSON schema for operations
// that take a siteId and a single resource ID.
func siteAndIDSchema(
	idName string,
	idDesc string,
) map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"siteId": siteIDSchema(),
			idName: map[string]interface{}{
				"type":        "string",
				"description": idDesc,
			},
		},
		"required": []string{idName},
	}
}

// listSchema returns the standard JSON schema for list operations
// with siteId + pagination parameters.
func listSchema() map[string]interface{} {
	props := map[string]interface{}{
		"siteId": siteIDSchema(),
	}
	for k, v := range paginationSchema() {
		props[k] = v
	}
	return map[string]interface{}{
		"type":       "object",
		"properties": props,
	}
}
