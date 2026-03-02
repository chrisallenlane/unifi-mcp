// Package server implements the MCP (Model Context Protocol) JSON-RPC server.
package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/chrisallenlane/unifi-mcp-server/internal/tools"
	"github.com/chrisallenlane/unifi-mcp-server/internal/unifi"
)

// Constants for server configuration
const (
	MCPProtocolVersion   = "2024-11-05"
	ServerName           = "unifi-mcp-server"
	ServerVersion        = "0.1.0"
	ToolExecutionTimeout = 30 * time.Second
)

// Server represents an MCP server
type Server struct {
	client        *unifi.ClientWithResponses
	defaultSiteID string
	tools         map[string]tools.Tool
}

// New creates a new MCP server
func New(
	client *unifi.ClientWithResponses,
	defaultSiteID string,
) *Server {
	s := &Server{
		client:        client,
		defaultSiteID: defaultSiteID,
		tools:         make(map[string]tools.Tool),
	}

	s.registerTools()

	return s
}

// registerTools registers all available tools
func (s *Server) registerTools() {
	s.tools["get_info"] = tools.NewGetInfo(s.client)
	s.tools["list_sites"] = tools.NewListSites(
		s.client,
		s.defaultSiteID,
	)
	s.tools["list_firewall_zones"] = tools.NewListFirewallZones(
		s.client,
		s.defaultSiteID,
	)
	s.tools["get_firewall_zone"] = tools.NewGetFirewallZone(
		s.client,
		s.defaultSiteID,
	)
	s.tools["create_firewall_zone"] = tools.NewCreateFirewallZone(
		s.client,
		s.defaultSiteID,
	)
	s.tools["update_firewall_zone"] = tools.NewUpdateFirewallZone(
		s.client,
		s.defaultSiteID,
	)
	s.tools["delete_firewall_zone"] = tools.NewDeleteFirewallZone(
		s.client,
		s.defaultSiteID,
	)
	s.tools["list_firewall_policies"] = tools.NewListFirewallPolicies(
		s.client,
		s.defaultSiteID,
	)
	s.tools["get_firewall_policy"] = tools.NewGetFirewallPolicy(
		s.client,
		s.defaultSiteID,
	)
	s.tools["create_firewall_policy"] = tools.NewCreateFirewallPolicy(
		s.client,
		s.defaultSiteID,
	)
	s.tools["update_firewall_policy"] = tools.NewUpdateFirewallPolicy(
		s.client,
		s.defaultSiteID,
	)
	s.tools["delete_firewall_policy"] = tools.NewDeleteFirewallPolicy(
		s.client,
		s.defaultSiteID,
	)
	s.tools["patch_firewall_policy"] = tools.NewPatchFirewallPolicy(
		s.client,
		s.defaultSiteID,
	)
	s.tools["get_firewall_policy_ordering"] = tools.NewGetFirewallPolicyOrdering(
		s.client,
		s.defaultSiteID,
	)
	s.tools["update_firewall_policy_ordering"] = tools.NewUpdateFirewallPolicyOrdering(
		s.client,
		s.defaultSiteID,
	)
	s.tools["list_wans"] = tools.NewListWans(
		s.client,
		s.defaultSiteID,
	)
	s.tools["list_vpn_tunnels"] = tools.NewListVpnTunnels(
		s.client,
		s.defaultSiteID,
	)
	s.tools["list_vpn_servers"] = tools.NewListVpnServers(
		s.client,
		s.defaultSiteID,
	)
	s.tools["list_radius_profiles"] = tools.NewListRadiusProfiles(
		s.client,
		s.defaultSiteID,
	)
	s.tools["list_device_tags"] = tools.NewListDeviceTags(
		s.client,
		s.defaultSiteID,
	)
	s.tools["list_dpi_categories"] = tools.NewListDpiCategories(
		s.client,
	)
	s.tools["list_dpi_applications"] = tools.NewListDpiApplications(
		s.client,
	)
	s.tools["list_countries"] = tools.NewListCountries(
		s.client,
	)
	s.tools["list_clients"] = tools.NewListClients(
		s.client,
		s.defaultSiteID,
	)
	s.tools["get_client"] = tools.NewGetClient(
		s.client,
		s.defaultSiteID,
	)
	s.tools["execute_client_action"] = tools.NewExecuteClientAction(
		s.client,
		s.defaultSiteID,
	)
	s.tools["list_dns_policies"] = tools.NewListDNSPolicies(
		s.client,
		s.defaultSiteID,
	)
	s.tools["get_dns_policy"] = tools.NewGetDNSPolicy(
		s.client,
		s.defaultSiteID,
	)
	s.tools["create_dns_policy"] = tools.NewCreateDNSPolicy(
		s.client,
		s.defaultSiteID,
	)
	s.tools["update_dns_policy"] = tools.NewUpdateDNSPolicy(
		s.client,
		s.defaultSiteID,
	)
	s.tools["delete_dns_policy"] = tools.NewDeleteDNSPolicy(
		s.client,
		s.defaultSiteID,
	)
}

// Run starts the MCP server and processes requests
func (s *Server) Run(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
) error {
	scanner := bufio.NewScanner(stdin)
	encoder := json.NewEncoder(stdout)

	for scanner.Scan() {
		line := scanner.Bytes()

		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			log.Printf("Failed to parse request: %v", err)
			// Send error response for malformed JSON-RPC request
			errResp := &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      nil,
				Error: &JSONRPCError{
					Code:    -32700,
					Message: fmt.Sprintf("Parse error: %v", err),
				},
			}
			if encErr := encoder.Encode(errResp); encErr != nil {
				log.Printf(
					"Failed to encode error response: %v",
					encErr,
				)
			}
			continue
		}

		resp := s.handleRequest(ctx, &req)
		if err := encoder.Encode(resp); err != nil {
			log.Printf("Failed to encode response: %v", err)
			return err
		}
	}

	return scanner.Err()
}

// handleRequest processes a JSON-RPC request
func (s *Server) handleRequest(
	ctx context.Context,
	req *JSONRPCRequest,
) *JSONRPCResponse {
	resp := &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "initialize":
		resp.Result = s.handleInitialize(ctx, req.Params)
	case "tools/list":
		resp.Result = s.handleListTools(ctx)
	case "tools/call":
		result, err := s.handleCallTool(ctx, req.Params)
		if err != nil {
			resp.Error = &JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			}
		} else {
			resp.Result = result
		}
	default:
		resp.Error = &JSONRPCError{
			Code: -32601,
			Message: fmt.Sprintf(
				"Method not found: %s",
				req.Method,
			),
		}
	}

	return resp
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(
	_ context.Context,
	_ json.RawMessage,
) interface{} {
	return map[string]interface{}{
		"protocolVersion": MCPProtocolVersion,
		"capabilities": map[string]interface{}{
			"tools": map[string]bool{},
		},
		"serverInfo": map[string]string{
			"name":    ServerName,
			"version": ServerVersion,
		},
	}
}

// handleListTools returns the list of available tools
func (s *Server) handleListTools(
	_ context.Context,
) interface{} {
	toolList := make(
		[]map[string]interface{},
		0,
		len(s.tools),
	)

	for name, tool := range s.tools {
		toolList = append(
			toolList,
			map[string]interface{}{
				"name":        name,
				"description": tool.Description(),
				"inputSchema": tool.InputSchema(),
			},
		)
	}

	return map[string]interface{}{
		"tools": toolList,
	}
}

// handleCallTool executes a tool
func (s *Server) handleCallTool(
	ctx context.Context,
	params json.RawMessage,
) (interface{}, error) {
	var callParams struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal(params, &callParams); err != nil {
		return nil, fmt.Errorf(
			"failed to parse tool call params: %w",
			err,
		)
	}

	tool, exists := s.tools[callParams.Name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", callParams.Name)
	}

	toolCtx, cancel := context.WithTimeout(
		ctx,
		ToolExecutionTimeout,
	)
	defer cancel()

	result, err := tool.Execute(toolCtx, callParams.Arguments)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": result,
			},
		},
	}, nil
}
