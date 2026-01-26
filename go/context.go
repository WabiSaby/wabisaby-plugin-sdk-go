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

	"github.com/google/uuid"
	"github.com/wabisaby/wabisaby-plugin-sdk/go/stub"
	pluginpb "github.com/wabisaby/wabisaby-protos/go/plugin"
)

// PluginSession provides execution context for a plugin.
type PluginSession struct {
	// Logger provides structured logging.
	Logger *stub.Logger

	// TenantID is the ID of the tenant executing this plugin.
	TenantID uuid.UUID

	// PluginID is the ID of this plugin.
	PluginID uuid.UUID

	// Config is the plugin configuration (JSON-decoded), nil if not provided.
	Config map[string]any
}

// Context provides the execution context for a plugin.
// It provides access to services (stub) and execution context (session).
type Context struct {
	// Context is the standard Go context for cancellation and timeouts.
	context.Context

	// Direct accessors for convenience (new)
	Storage      *stub.StorageClient
	HTTP         HTTPService
	Queue        *stub.QueueClient
	Notification NotificationService
	Secrets      *stub.SecretsClient
	Songs        *stub.SongClient
	Users        *stub.UserClient
	Logger       *ContextLogger
	TenantID     uuid.UUID
	PluginID     uuid.UUID
	Config       *ConfigAccessor

	// Backward compatibility - use GetStub() and GetSession() for access
	stub    *stub.PluginStub
	session *PluginSession
}

// GetStub returns the plugin stub with semantically grouped API services.
func (c *Context) GetStub() *stub.PluginStub {
	return c.stub
}

// GetSession returns the plugin session with execution context.
func (c *Context) GetSession() *PluginSession {
	return c.session
}

// NewContext creates a new plugin context.
func NewContext(
	ctx context.Context,
	tenantID, pluginID uuid.UUID,
	capabilitiesClient pluginpb.PluginCapabilitiesServiceClient,
	config map[string]interface{},
) *Context {
	// Create all clients
	storageClient := stub.NewStorageClient(tenantID, pluginID, capabilitiesClient)
	httpClient := stub.NewHTTPClient(tenantID, pluginID, capabilitiesClient)
	queueClient := stub.NewQueueClient(tenantID, pluginID, capabilitiesClient)
	notificationClient := stub.NewNotificationClient(tenantID, pluginID, capabilitiesClient)
	secretsClient := stub.NewSecretsClient(tenantID, pluginID, capabilitiesClient)
	songClient := stub.NewSongClient(tenantID, pluginID, capabilitiesClient)
	userClient := stub.NewUserClient(tenantID, pluginID, capabilitiesClient)
	logger := stub.NewLogger(tenantID, pluginID, capabilitiesClient)

	// Initialize PluginStub with grouped clients
	pluginStub := &stub.PluginStub{}
	pluginStub.Data.Storage = storageClient
	pluginStub.Data.Secrets = secretsClient
	pluginStub.Music.Queue = queueClient
	pluginStub.Music.Songs = songClient
	pluginStub.Users = userClient
	pluginStub.Communication.Notify = notificationClient
	pluginStub.Network.HTTP = httpClient

	// Initialize PluginSession
	session := &PluginSession{
		Logger:   logger,
		TenantID: tenantID,
		PluginID: pluginID,
		Config:   config,
	}

	// Create context with direct accessors
	pluginCtx := &Context{
		Context: ctx,
		stub:    pluginStub,
		session: session,
		// Direct accessors
		Storage:      storageClient,
		HTTP:         NewHTTPService(httpClient),
		Queue:        queueClient,
		Notification: NewNotificationService(notificationClient),
		Secrets:      secretsClient,
		Songs:        songClient,
		Users:        userClient,
		TenantID:     tenantID,
		PluginID:     pluginID,
		Config:       NewConfigAccessor(config),
	}

	// Create logger with context reference
	pluginCtx.Logger = NewContextLogger(logger, pluginCtx)

	return pluginCtx
}
