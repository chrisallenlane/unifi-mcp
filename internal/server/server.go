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

	"github.com/chrisallenlane/unifi-mcp/internal/tools"
	"github.com/chrisallenlane/unifi-mcp/internal/unifi"
)

// Constants for server configuration
const (
	MCPProtocolVersion   = "2024-11-05"
	ServerName           = "unifi-mcp"
	ServerVersion        = "0.1.1"
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
	}

	s.registerTools()

	return s
}

// registerTools registers all available tools
func (s *Server) registerTools() {
	c, sid := s.client, s.defaultSiteID

	s.tools = map[string]tools.Tool{
		// Info / Sites
		"get_info":   tools.NewGetInfo(c),
		"list_sites": tools.NewListSites(c, sid),

		// Firewall Zones
		"list_firewall_zones":  tools.NewListFirewallZones(c, sid),
		"get_firewall_zone":    tools.NewGetFirewallZone(c, sid),
		"create_firewall_zone": tools.NewCreateFirewallZone(c, sid),
		"update_firewall_zone": tools.NewUpdateFirewallZone(c, sid),
		"delete_firewall_zone": tools.NewDeleteFirewallZone(c, sid),

		// Firewall Policies
		"list_firewall_policies": tools.NewListFirewallPolicies(
			c,
			sid,
		),
		"get_firewall_policy": tools.NewGetFirewallPolicy(c, sid),
		"create_firewall_policy": tools.NewCreateFirewallPolicy(
			c,
			sid,
		),
		"update_firewall_policy": tools.NewUpdateFirewallPolicy(
			c,
			sid,
		),
		"delete_firewall_policy": tools.NewDeleteFirewallPolicy(
			c,
			sid,
		),
		"patch_firewall_policy": tools.NewPatchFirewallPolicy(c, sid),
		"get_firewall_policy_ordering": tools.NewGetFirewallPolicyOrdering(
			c,
			sid,
		),
		"update_firewall_policy_ordering": tools.NewUpdateFirewallPolicyOrdering(
			c,
			sid,
		),

		// Networks
		"list_networks":          tools.NewListNetworks(c, sid),
		"get_network":            tools.NewGetNetwork(c, sid),
		"create_network":         tools.NewCreateNetwork(c, sid),
		"update_network":         tools.NewUpdateNetwork(c, sid),
		"delete_network":         tools.NewDeleteNetwork(c, sid),
		"get_network_references": tools.NewGetNetworkReferences(c, sid),

		// Clients
		"list_clients":          tools.NewListClients(c, sid),
		"get_client":            tools.NewGetClient(c, sid),
		"execute_client_action": tools.NewExecuteClientAction(c, sid),

		// Devices
		"list_devices":          tools.NewListDevices(c, sid),
		"get_device":            tools.NewGetDevice(c, sid),
		"adopt_device":          tools.NewAdoptDevice(c, sid),
		"remove_device":         tools.NewRemoveDevice(c, sid),
		"execute_device_action": tools.NewExecuteDeviceAction(c, sid),
		"execute_port_action":   tools.NewExecutePortAction(c, sid),
		"get_device_statistics": tools.NewGetDeviceStatistics(c, sid),
		"list_pending_devices":  tools.NewListPendingDevices(c),

		// ACL Rules
		"list_acl_rules":           tools.NewListACLRules(c, sid),
		"get_acl_rule":             tools.NewGetACLRule(c, sid),
		"create_acl_rule":          tools.NewCreateACLRule(c, sid),
		"update_acl_rule":          tools.NewUpdateACLRule(c, sid),
		"delete_acl_rule":          tools.NewDeleteACLRule(c, sid),
		"get_acl_rule_ordering":    tools.NewGetACLRuleOrdering(c, sid),
		"update_acl_rule_ordering": tools.NewUpdateACLRuleOrdering(c, sid),

		// DNS Policies
		"list_dns_policies": tools.NewListDNSPolicies(c, sid),
		"get_dns_policy":    tools.NewGetDNSPolicy(c, sid),
		"create_dns_policy": tools.NewCreateDNSPolicy(c, sid),
		"update_dns_policy": tools.NewUpdateDNSPolicy(c, sid),
		"delete_dns_policy": tools.NewDeleteDNSPolicy(c, sid),

		// Traffic Matching Lists
		"list_traffic_matching_lists": tools.NewListTrafficMatchingLists(
			c,
			sid,
		),
		"get_traffic_matching_list": tools.NewGetTrafficMatchingList(c, sid),
		"create_traffic_matching_list": tools.NewCreateTrafficMatchingList(
			c,
			sid,
		),
		"update_traffic_matching_list": tools.NewUpdateTrafficMatchingList(
			c,
			sid,
		),
		"delete_traffic_matching_list": tools.NewDeleteTrafficMatchingList(
			c,
			sid,
		),

		// WiFi Broadcasts
		"list_wifi_broadcasts":  tools.NewListWiFiBroadcasts(c, sid),
		"get_wifi_broadcast":    tools.NewGetWiFiBroadcast(c, sid),
		"create_wifi_broadcast": tools.NewCreateWiFiBroadcast(c, sid),
		"update_wifi_broadcast": tools.NewUpdateWiFiBroadcast(c, sid),
		"delete_wifi_broadcast": tools.NewDeleteWiFiBroadcast(c, sid),

		// Hotspot Vouchers
		"list_vouchers":   tools.NewListVouchers(c, sid),
		"get_voucher":     tools.NewGetVoucher(c, sid),
		"create_vouchers": tools.NewCreateVouchers(c, sid),
		"delete_vouchers": tools.NewDeleteVouchers(c, sid),
		"delete_voucher":  tools.NewDeleteVoucher(c, sid),

		// Supporting Read-Only
		"list_wans":             tools.NewListWans(c, sid),
		"list_vpn_tunnels":      tools.NewListVpnTunnels(c, sid),
		"list_vpn_servers":      tools.NewListVpnServers(c, sid),
		"list_radius_profiles":  tools.NewListRadiusProfiles(c, sid),
		"list_device_tags":      tools.NewListDeviceTags(c, sid),
		"list_dpi_categories":   tools.NewListDpiCategories(c),
		"list_dpi_applications": tools.NewListDpiApplications(c),
		"list_countries":        tools.NewListCountries(c),
	}
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
		if resp == nil {
			continue
		}
		if err := encoder.Encode(resp); err != nil {
			log.Printf("Failed to encode response: %v", err)
			return err
		}
	}

	return scanner.Err()
}

// handleRequest processes a JSON-RPC request. Returns nil for
// notifications (JSON-RPC 2.0 requests without an id), which MUST
// NOT receive a response per the spec. The MCP handshake sends
// "notifications/initialized" after initialize; responding to it
// causes strict clients to abort the handshake.
func (s *Server) handleRequest(
	ctx context.Context,
	req *JSONRPCRequest,
) *JSONRPCResponse {
	if req.ID == nil {
		return nil
	}

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
