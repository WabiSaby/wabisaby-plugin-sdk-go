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

// QueueItem represents a queue item with song information.
type QueueItem struct {
	ID            string                 `json:"id"`
	TenantID      string                 `json:"tenant_id"`
	SongID        string                 `json:"song_id"`
	RequesterName string                 `json:"requester_name"`
	Position      int                    `json:"position"`
	IsPriority    bool                   `json:"is_priority"`
	Status        string                 `json:"status"`
	Song          map[string]interface{} `json:"song,omitempty"`
}

// QueueClient provides access to queue operations.
type QueueClient struct {
	tenantID uuid.UUID
	pluginID uuid.UUID
	client   pluginpb.PluginCapabilitiesServiceClient
}

// NewQueueClient creates a new queue client.
func NewQueueClient(tenantID, pluginID uuid.UUID, client pluginpb.PluginCapabilitiesServiceClient) *QueueClient {
	return &QueueClient{
		tenantID: tenantID,
		pluginID: pluginID,
		client:   client,
	}
}

// Get retrieves the current queue.
func (c *QueueClient) Get(ctx context.Context) ([]*QueueItem, error) {
	req := &pluginpb.QueueGetRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
	}

	resp, err := c.client.QueueGet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("queue get failed: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("queue error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	var items []*QueueItem
	if err := json.Unmarshal(resp.QueueData, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal queue data: %w", err)
	}

	return items, nil
}

// Add adds a song to the queue.
// songData can be either a song ID (string) or a full song object (map[string]interface{}).
// position is the position to insert at (-1 for end of queue).
func (c *QueueClient) Add(ctx context.Context, songData interface{}, position int) error {
	songJSON, err := json.Marshal(songData)
	if err != nil {
		return fmt.Errorf("failed to marshal song data: %w", err)
	}

	req := &pluginpb.QueueAddRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		SongData: songJSON,
		Position: int32(position),
	}

	resp, err := c.client.QueueAdd(ctx, req)
	if err != nil {
		return fmt.Errorf("queue add failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("queue error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return nil
}

// Remove removes a queue item by position.
func (c *QueueClient) Remove(ctx context.Context, position int) error {
	req := &pluginpb.QueueRemoveRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		Position: int32(position),
	}

	resp, err := c.client.QueueRemove(ctx, req)
	if err != nil {
		return fmt.Errorf("queue remove failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("queue error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return nil
}

// Reorder moves a queue item from one position to another.
func (c *QueueClient) Reorder(ctx context.Context, fromPosition, toPosition int) error {
	req := &pluginpb.QueueReorderRequest{
		TenantId:     c.tenantID.String(),
		PluginId:     c.pluginID.String(),
		FromPosition: int32(fromPosition),
		ToPosition:   int32(toPosition),
	}

	resp, err := c.client.QueueReorder(ctx, req)
	if err != nil {
		return fmt.Errorf("queue reorder failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("queue error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return nil
}
