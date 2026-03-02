// Package main is the entry point for the UniFi MCP server.
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/chrisallenlane/unifi-mcp-server/internal/server"
	"github.com/chrisallenlane/unifi-mcp-server/internal/unifi"
)

func main() {
	apiURL := os.Getenv("UNIFI_API_URL")
	if apiURL == "" {
		fmt.Fprintln(os.Stderr, "UNIFI_API_URL is required")
		os.Exit(1)
	}

	apiKey := os.Getenv("UNIFI_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "UNIFI_API_KEY is required")
		os.Exit(1)
	}

	siteID := os.Getenv("UNIFI_SITE_ID")

	httpClient := &http.Client{}
	if os.Getenv("UNIFI_INSECURE") != "" {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // #nosec G402
			},
		}
	}

	addAPIKey := func(
		_ context.Context,
		req *http.Request,
	) error {
		req.Header.Set("X-API-Key", apiKey)
		return nil
	}

	client, err := unifi.NewClientWithResponses(
		apiURL,
		unifi.WithHTTPClient(httpClient),
		unifi.WithRequestEditorFn(addAPIKey),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create client: %v\n", err)
		os.Exit(1)
	}

	s := server.New(client, siteID)
	if err := s.Run(context.Background(), os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
