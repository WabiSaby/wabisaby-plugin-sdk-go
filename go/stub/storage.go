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

package stub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	pluginpb "github.com/wabisaby/wabisaby/api/generated/proto/plugin"
)

// StorageClient provides access to plugin storage (key-value store).
type StorageClient struct {
	tenantID uuid.UUID
	pluginID uuid.UUID
	client   pluginpb.PluginCapabilitiesServiceClient
}

// NewStorageClient creates a new storage client.
func NewStorageClient(tenantID, pluginID uuid.UUID, client pluginpb.PluginCapabilitiesServiceClient) *StorageClient {
	return &StorageClient{
		tenantID: tenantID,
		pluginID: pluginID,
		client:   client,
	}
}

// Get retrieves a value from storage.
// Returns nil if the key doesn't exist.
func (c *StorageClient) Get(ctx context.Context, key string) (interface{}, error) {
	req := &pluginpb.StorageGetRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		Key:      key,
	}

	resp, err := c.client.StorageGet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("storage get failed: %w", err)
	}

	if resp.Error != nil {
		if resp.Error.Code == "NOT_FOUND" {
			return nil, nil
		}
		return nil, fmt.Errorf("storage error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	if len(resp.Value) == 0 {
		return nil, nil
	}

	var value interface{}
	if err := json.Unmarshal(resp.Value, &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal storage value: %w", err)
	}

	return value, nil
}

// Set stores a value in storage.
// The value must be JSON-serializable.
func (c *StorageClient) Set(ctx context.Context, key string, value interface{}) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	req := &pluginpb.StorageSetRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		Key:      key,
		Value:    valueJSON,
	}

	resp, err := c.client.StorageSet(ctx, req)
	if err != nil {
		return fmt.Errorf("storage set failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("storage error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return nil
}

// Delete removes a value from storage.
func (c *StorageClient) Delete(ctx context.Context, key string) error {
	req := &pluginpb.StorageDeleteRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		Key:      key,
	}

	resp, err := c.client.StorageDelete(ctx, req)
	if err != nil {
		return fmt.Errorf("storage delete failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("storage error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return nil
}

// Keys lists all keys in storage, optionally filtered by prefix.
func (c *StorageClient) Keys(ctx context.Context, prefix string) ([]string, error) {
	req := &pluginpb.StorageKeysRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		Prefix:   prefix,
	}

	resp, err := c.client.StorageKeys(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("storage keys failed: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("storage error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Keys, nil
}
