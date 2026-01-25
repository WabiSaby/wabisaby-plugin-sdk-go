// Copyright (c) 2026 WabiSaby
// All rights reserved.
//
// This source code is proprietary and confidential. Unauthorized copying,
// modification, distribution, or use of this software, via any medium is
// strictly prohibited without the express written permission of WabiSaby.
//
// This software contains confidential and proprietary information of
// WabiSaby and its licensors. Use, disclosure, or reproduction
// is prohibited without the prior express written permission of WabiSaby.

package sdk

import (
	"context"

	"github.com/hashicorp/go-plugin"
	pluginpb "github.com/wabisaby/wabisaby/api/generated/proto/plugin"
	"google.golang.org/grpc"
)

// PluginGRPC implements the plugin.Plugin interface for gRPC
type PluginGRPC struct {
	plugin.Plugin
	// Impl will be set by the plugin binary
	Impl pluginpb.PluginExecutionServiceServer
}

// GRPCServer registers the gRPC server
func (p *PluginGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pluginpb.RegisterPluginExecutionServiceServer(s, p.Impl)
	return nil
}

// GRPCClient creates a gRPC client
func (p *PluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return pluginpb.NewPluginExecutionServiceClient(c), nil
}

// HandshakeConfig returns the handshake configuration
func HandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "WABISABY_PLUGIN",
		MagicCookieValue: "wabisaby-plugin-runtime",
	}
}
