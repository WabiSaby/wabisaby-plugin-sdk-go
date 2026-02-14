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

// NotificationType represents the type of notification.
type NotificationType string

const (
	NotificationTypeInfo    NotificationType = "info"
	NotificationTypeSuccess NotificationType = "success"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeError   NotificationType = "error"
)

// NotificationClient provides access to notification operations.
type NotificationClient struct {
	tenantID uuid.UUID
	pluginID uuid.UUID
	client   pluginpb.PluginCapabilitiesServiceClient
}

// NewNotificationClient creates a new notification client.
func NewNotificationClient(tenantID, pluginID uuid.UUID, client pluginpb.PluginCapabilitiesServiceClient) *NotificationClient {
	return &NotificationClient{
		tenantID: tenantID,
		pluginID: pluginID,
		client:   client,
	}
}

// Send sends a notification to a user.
// userID is optional - if empty, notification is sent to all users in the tenant.
func (c *NotificationClient) Send(
	ctx context.Context,
	userID string,
	title, message string,
	notifType NotificationType,
) (string, error) {
	req := &pluginpb.NotificationSendRequest{
		TenantId:         c.tenantID.String(),
		PluginId:         c.pluginID.String(),
		UserId:           userID,
		Title:            title,
		Message:          message,
		NotificationType: string(notifType),
	}

	resp, err := c.client.NotificationSend(ctx, req)
	if err != nil {
		return "", fmt.Errorf("notification send failed: %w", err)
	}

	if resp.Error != nil {
		return "", fmt.Errorf("notification error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.NotificationId, nil
}
