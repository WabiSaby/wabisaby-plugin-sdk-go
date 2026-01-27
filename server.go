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
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	hashicorp_plugin "github.com/hashicorp/go-plugin"
	pluginpb "github.com/wabisaby/wabisaby-protos/go/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server wraps a plugin implementation and provides the gRPC server.
type Server struct {
	pluginpb.UnimplementedPluginExecutionServiceServer

	plugin             Plugin
	capabilitiesClient pluginpb.PluginCapabilitiesServiceClient
	capabilitiesConn   *grpc.ClientConn

	// Initialization state tracking
	initOnce     sync.Once
	initErr      error
	shutdownOnce sync.Once
}

// NewServer creates a new plugin server.
// It connects to the capabilities service using WABISABY_CAPABILITIES_ADDR environment variable.
func NewServer(plugin Plugin) (*Server, error) {
	// Connect to capabilities service
	capabilitiesAddr := os.Getenv("WABISABY_CAPABILITIES_ADDR")
	if capabilitiesAddr == "" {
		return nil, fmt.Errorf("WABISABY_CAPABILITIES_ADDR environment variable not set")
	}

	conn, err := grpc.NewClient(
		capabilitiesAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to capabilities service: %w", err)
	}

	capabilitiesClient := pluginpb.NewPluginCapabilitiesServiceClient(conn)

	return &Server{
		plugin:             plugin,
		capabilitiesClient: capabilitiesClient,
		capabilitiesConn:   conn,
	}, nil
}

// Close closes the capabilities connection.
func (s *Server) Close() error {
	if s.capabilitiesConn != nil {
		return s.capabilitiesConn.Close()
	}
	return nil
}

// ExecuteCommand implements PluginExecutionServiceServer.ExecuteCommand.
func (s *Server) ExecuteCommand(ctx context.Context, req *pluginpb.ExecuteCommandRequest) (*pluginpb.ExecuteCommandResponse, error) {
	tenantID, err := uuid.Parse(req.TenantId)
	if err != nil {
		return &pluginpb.ExecuteCommandResponse{
			Result: &pluginpb.ExecuteCommandResponse_Error{
				Error: &pluginpb.PluginError{
					Code:    "INVALID_ARGUMENT",
					Message: fmt.Sprintf("invalid tenant ID: %v", err),
				},
			},
		}, nil
	}

	pluginID, err := uuid.Parse(req.PluginId)
	if err != nil {
		return &pluginpb.ExecuteCommandResponse{
			Result: &pluginpb.ExecuteCommandResponse_Error{
				Error: &pluginpb.PluginError{
					Code:    "INVALID_ARGUMENT",
					Message: fmt.Sprintf("invalid plugin ID: %v", err),
				},
			},
		}, nil
	}

	// Check if plugin implements CommandPlugin
	commandPlugin, ok := s.plugin.(CommandPlugin)
	if !ok {
		return &pluginpb.ExecuteCommandResponse{
			Result: &pluginpb.ExecuteCommandResponse_Error{
				Error: &pluginpb.PluginError{
					Code:    "NOT_SUPPORTED",
					Message: "plugin does not support command execution",
				},
			},
		}, nil
	}

	// Unmarshal arguments
	args := make([]interface{}, 0, len(req.Args))
	for _, argBytes := range req.Args {
		var arg interface{}
		if err := json.Unmarshal(argBytes, &arg); err != nil {
			return &pluginpb.ExecuteCommandResponse{
				Result: &pluginpb.ExecuteCommandResponse_Error{
					Error: &pluginpb.PluginError{
						Code:    "INVALID_ARGUMENT",
						Message: fmt.Sprintf("failed to unmarshal argument: %v", err),
					},
				},
			}, nil
		}
		args = append(args, arg)
	}

	// Create execution context with timeout if specified
	execCtx := ctx
	var cancel context.CancelFunc
	if req.TimeoutMs > 0 {
		execCtx, cancel = context.WithTimeout(ctx, time.Duration(req.TimeoutMs)*time.Millisecond)
		defer cancel()
	}

	// Create plugin context
	pluginCtx := NewContext(execCtx, tenantID, pluginID, s.capabilitiesClient, nil)

	startTime := time.Now()
	result, err := commandPlugin.ExecuteCommand(pluginCtx, req.Command, args)
	executionTime := time.Since(startTime)

	if err != nil {
		return &pluginpb.ExecuteCommandResponse{
			Result: &pluginpb.ExecuteCommandResponse_Error{
				Error: &pluginpb.PluginError{
					Code:    "EXECUTION_ERROR",
					Message: err.Error(),
				},
			},
			ExecutionTimeMs: executionTime.Milliseconds(),
		}, nil
	}

	// Marshal result
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return &pluginpb.ExecuteCommandResponse{
			Result: &pluginpb.ExecuteCommandResponse_Error{
				Error: &pluginpb.PluginError{
					Code:    "SERIALIZATION_ERROR",
					Message: fmt.Sprintf("failed to marshal result: %v", err),
				},
			},
			ExecutionTimeMs: executionTime.Milliseconds(),
		}, nil
	}

	return &pluginpb.ExecuteCommandResponse{
		Result: &pluginpb.ExecuteCommandResponse_Data{
			Data: resultJSON,
		},
		ExecutionTimeMs: executionTime.Milliseconds(),
	}, nil
}

// EnablePlugin implements PluginExecutionServiceServer.EnablePlugin.
func (s *Server) EnablePlugin(ctx context.Context, req *pluginpb.EnablePluginRequest) (*pluginpb.EnablePluginResponse, error) {
	tenantID, err := uuid.Parse(req.TenantId)
	if err != nil {
		return &pluginpb.EnablePluginResponse{
			Error: &pluginpb.PluginError{
				Code:    "INVALID_ARGUMENT",
				Message: fmt.Sprintf("invalid tenant ID: %v", err),
			},
		}, nil
	}

	pluginID, err := uuid.Parse(req.PluginId)
	if err != nil {
		return &pluginpb.EnablePluginResponse{
			Error: &pluginpb.PluginError{
				Code:    "INVALID_ARGUMENT",
				Message: fmt.Sprintf("invalid plugin ID: %v", err),
			},
		}, nil
	}

	// Decode config if provided
	var config map[string]interface{}
	if len(req.Config) > 0 {
		if err := json.Unmarshal(req.Config, &config); err != nil {
			return &pluginpb.EnablePluginResponse{
				Error: &pluginpb.PluginError{
					Code:    "INVALID_ARGUMENT",
					Message: fmt.Sprintf("failed to decode config: %v", err),
				},
			}, nil
		}
	}

	// Create context with capabilities and config
	// Note: For stateless plugins, this is a no-op, but we create the context
	// to validate config and ensure it's available if needed
	_ = NewContext(ctx, tenantID, pluginID, s.capabilitiesClient, config)

	// For stateless plugins, this is a no-op
	// Stateful plugins would maintain state here
	return &pluginpb.EnablePluginResponse{
		Success:    true,
		InstanceId: "default",
	}, nil
}

// DisablePlugin implements PluginExecutionServiceServer.DisablePlugin.
func (s *Server) DisablePlugin(ctx context.Context, req *pluginpb.DisablePluginRequest) (*pluginpb.DisablePluginResponse, error) {
	// For stateless plugins, this is a no-op
	return &pluginpb.DisablePluginResponse{
		Success: true,
	}, nil
}

// StreamEvents implements PluginExecutionServiceServer.StreamEvents.
// This is for stateful plugins that receive event streams.
func (s *Server) StreamEvents(stream pluginpb.PluginExecutionService_StreamEventsServer) error {
	// This would be implemented for stateful plugins
	// For now, return an error indicating it's not supported
	return fmt.Errorf("streaming events not yet implemented")
}

// InitializePlugin implements PluginExecutionServiceServer.InitializePlugin.
func (s *Server) InitializePlugin(ctx context.Context, req *pluginpb.InitializePluginRequest) (*pluginpb.InitializePluginResponse, error) {
	tenantID, err := uuid.Parse(req.TenantId)
	if err != nil {
		return &pluginpb.InitializePluginResponse{
			Error: &pluginpb.PluginError{
				Code:    "INVALID_ARGUMENT",
				Message: fmt.Sprintf("invalid tenant ID: %v", err),
			},
		}, nil
	}

	pluginID, err := uuid.Parse(req.PluginId)
	if err != nil {
		return &pluginpb.InitializePluginResponse{
			Error: &pluginpb.PluginError{
				Code:    "INVALID_ARGUMENT",
				Message: fmt.Sprintf("invalid plugin ID: %v", err),
			},
		}, nil
	}

	// Decode config if provided
	var config map[string]interface{}
	if len(req.Config) > 0 {
		if err := json.Unmarshal(req.Config, &config); err != nil {
			return &pluginpb.InitializePluginResponse{
				Error: &pluginpb.PluginError{
					Code:    "INVALID_ARGUMENT",
					Message: fmt.Sprintf("failed to decode config: %v", err),
				},
			}, nil
		}
	}

	// Use sync.Once to ensure Initialize is called only once per plugin process
	s.initOnce.Do(func() {
		pluginCtx := NewContext(ctx, tenantID, pluginID, s.capabilitiesClient, config)
		s.initErr = s.plugin.Initialize(pluginCtx)
	})

	if s.initErr != nil {
		return &pluginpb.InitializePluginResponse{
			Error: &pluginpb.PluginError{
				Code:    "INITIALIZATION_ERROR",
				Message: s.initErr.Error(),
			},
		}, nil
	}

	return &pluginpb.InitializePluginResponse{
		Success: true,
	}, nil
}

// ShutdownPlugin implements PluginExecutionServiceServer.ShutdownPlugin.
func (s *Server) ShutdownPlugin(ctx context.Context, req *pluginpb.ShutdownPluginRequest) (*pluginpb.ShutdownPluginResponse, error) {
	tenantID, err := uuid.Parse(req.TenantId)
	if err != nil {
		return &pluginpb.ShutdownPluginResponse{
			Error: &pluginpb.PluginError{
				Code:    "INVALID_ARGUMENT",
				Message: fmt.Sprintf("invalid tenant ID: %v", err),
			},
		}, nil
	}

	pluginID, err := uuid.Parse(req.PluginId)
	if err != nil {
		return &pluginpb.ShutdownPluginResponse{
			Error: &pluginpb.PluginError{
				Code:    "INVALID_ARGUMENT",
				Message: fmt.Sprintf("invalid plugin ID: %v", err),
			},
		}, nil
	}

	// Use sync.Once to ensure Shutdown is called only once
	var shutdownErr error
	s.shutdownOnce.Do(func() {
		pluginCtx := NewContext(ctx, tenantID, pluginID, s.capabilitiesClient, nil)
		shutdownErr = s.plugin.Shutdown(pluginCtx)
	})

	if shutdownErr != nil {
		return &pluginpb.ShutdownPluginResponse{
			Error: &pluginpb.PluginError{
				Code:    "SHUTDOWN_ERROR",
				Message: shutdownErr.Error(),
			},
		}, nil
	}

	return &pluginpb.ShutdownPluginResponse{
		Success: true,
	}, nil
}

// HealthCheck implements PluginExecutionServiceServer.HealthCheck.
func (s *Server) HealthCheck(ctx context.Context, req *pluginpb.HealthCheckRequest) (*pluginpb.HealthCheckResponse, error) {
	return &pluginpb.HealthCheckResponse{
		Status: pluginpb.HealthCheckResponse_SERVING,
	}, nil
}

// Serve starts the plugin server using HashiCorp go-plugin.
// This is the main entry point for plugin binaries.
func Serve(plugin Plugin) error {
	server, err := NewServer(plugin)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	defer server.Close()

	// Create plugin implementation
	pluginImpl := &PluginGRPC{
		Impl: server,
	}

	// Serve using go-plugin
	hashicorp_plugin.Serve(&hashicorp_plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig(),
		Plugins: map[string]hashicorp_plugin.Plugin{
			"plugin": pluginImpl,
		},
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			return grpc.NewServer(opts...)
		},
	})

	return nil
}
