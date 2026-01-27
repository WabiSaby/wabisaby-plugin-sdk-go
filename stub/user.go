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
	pluginpb "github.com/wabisaby/wabisaby-protos/go/plugin"
)

// UserInfo represents user information.
type UserInfo map[string]interface{}

// UserClient provides access to user operations.
type UserClient struct {
	tenantID uuid.UUID
	pluginID uuid.UUID
	client   pluginpb.PluginCapabilitiesServiceClient
}

// NewUserClient creates a new user client.
func NewUserClient(tenantID, pluginID uuid.UUID, client pluginpb.PluginCapabilitiesServiceClient) *UserClient {
	return &UserClient{
		tenantID: tenantID,
		pluginID: pluginID,
		client:   client,
	}
}

// Get retrieves user information by ID.
// Returns nil if the user doesn't exist.
func (c *UserClient) Get(ctx context.Context, userID string) (UserInfo, error) {
	req := &pluginpb.UserGetRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		UserId:   userID,
	}

	resp, err := c.client.UserGet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("user get failed: %w", err)
	}

	if resp.Error != nil {
		if resp.Error.Code == "NOT_FOUND" {
			return nil, nil
		}
		return nil, fmt.Errorf("user error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	if len(resp.User) == 0 {
		return nil, nil
	}

	var user UserInfo
	if err := json.Unmarshal(resp.User, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return user, nil
}
