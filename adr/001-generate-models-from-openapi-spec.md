# ADR 001: Generate models and client from official OpenAPI spec

## Status

Accepted

## Date

2026-03-01

## Context

We are building an MCP server that wraps the UniFi Network API. We need Go
types for API request/response bodies and an HTTP client to call the API
endpoints. There are several approaches:

1. **Hand-write models and client code.** Define Go structs and HTTP methods
   manually based on observed API behavior and documentation.

2. **Generate models and client from the official OpenAPI spec.** UniFi
   controller firmware embeds an OpenAPI 3.1.0 specification. Extract it and
   use `oapi-codegen` to produce typed Go code.

3. **Use an existing third-party Go UniFi library.** Several community
   libraries exist (e.g., `unpoller/unifi`, `paultyng/go-unifi`), but these
   target the older unofficial API, not the newer Integration API with API key
   authentication.

### Key considerations

- The UniFi Integration API (introduced in firmware ~7.x/8.x) is distinct
  from the older session-cookie-based controller API. Existing community
  libraries do not cover it.

- The official OpenAPI spec (v10.1.84) documents 67 operations and 360
  schemas. Hand-writing this would be error-prone and tedious to maintain.

- `oapi-codegen` produces a typed `ClientWithResponses` that gives us
  compile-time safety, method-per-endpoint coverage, and automatic
  request/response marshaling.

- The generated code has 12 `interface{}` fields from unresolved
  discriminated unions in the spec. None of these affect the priority
  endpoints (firewall zones, firewall policies). They are acceptable for
  best-effort formatting elsewhere.

- Checking the spec into the repo (`api/unifi-network.json`) makes generation
  reproducible and pins us to a known API version. Updating the spec is an
  explicit, reviewable action.

## Decision

Generate all Go types and the HTTP client from the official UniFi Network
OpenAPI 3.1.0 specification using `oapi-codegen`.

Specifically:

- The spec is checked into the repo at `api/unifi-network.json`.
- `oapi-codegen` generates two files into `internal/unifi/`:
  - `types.gen.go` — all model structs
  - `client.gen.go` — typed HTTP client with `ClientWithResponses`
- A `make generate` target reproduces the generated code.
- Generated files carry `// Code generated ... DO NOT EDIT.` headers and are
  excluded from linting/formatting tools.
- The 12 `interface{}` fields are accepted as-is.

## Consequences

### Positive

- **Correctness**: Types match the official spec exactly, reducing the risk of
  serialization bugs or missing fields.
- **Coverage**: All 67 operations and 360 schemas are available from day one,
  even if we only wrap a subset as MCP tools initially.
- **Maintainability**: Updating to a new firmware version means replacing the
  spec file and re-running `make generate`. Structural changes surface as
  compile errors.
- **No external runtime dependencies**: `oapi-codegen` is a build-time tool.
  The generated code depends only on `github.com/oapi-codegen/runtime` for
  UUID types and response parsing.

### Negative

- **Generated code is large**: ~1,971 lines of types and ~9,518 lines of
  client code. This adds to the repo size but is excluded from linting.
- **12 `interface{}` fields**: Some API responses will require runtime type
  assertions for fields with discriminated unions. This affects a small
  minority of endpoints and can be addressed later if needed.
- **Spec extraction is manual**: The OpenAPI spec must be pulled from
  controller firmware. There is no official distribution channel for it
  (yet). Updating requires access to the firmware image.
- **Build-time dependency on `oapi-codegen`**: Developers need this tool to
  regenerate code, though the generated files are checked in so casual
  contributors don't need it.
