// Package tools implements MCP tool definitions for the UniFi MCP server.
package tools

import (
	"fmt"

	"github.com/google/uuid"
)

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
