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
	"fmt"

	"github.com/google/uuid"
	pluginpb "github.com/wabisaby/wabisaby-protos-go/go/plugin"
)

// SecretsClient provides access to encrypted secret storage.
// Note: Secrets are encrypted at rest. The value returned by Get is encrypted.
// Plugins should handle encryption/decryption themselves or use a secure key management service.
type SecretsClient struct {
	tenantID uuid.UUID
	pluginID uuid.UUID
	client   pluginpb.PluginCapabilitiesServiceClient
}

// NewSecretsClient creates a new secrets client.
func NewSecretsClient(tenantID, pluginID uuid.UUID, client pluginpb.PluginCapabilitiesServiceClient) *SecretsClient {
	return &SecretsClient{
		tenantID: tenantID,
		pluginID: pluginID,
		client:   client,
	}
}

// Get retrieves an encrypted secret value.
// Returns the encrypted value as a string.
func (c *SecretsClient) Get(ctx context.Context, key string) (string, error) {
	req := &pluginpb.SecretGetRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		Key:      key,
	}

	resp, err := c.client.SecretGet(ctx, req)
	if err != nil {
		return "", fmt.Errorf("secret get failed: %w", err)
	}

	if resp.Error != nil {
		return "", fmt.Errorf("secret error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Value, nil
}

// Set stores an encrypted secret value.
// Note: The value should be encrypted before calling this method.
func (c *SecretsClient) Set(ctx context.Context, key, encryptedValue string) error {
	req := &pluginpb.SecretSetRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		Key:      key,
		Value:    encryptedValue,
	}

	resp, err := c.client.SecretSet(ctx, req)
	if err != nil {
		return fmt.Errorf("secret set failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("secret error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return nil
}
