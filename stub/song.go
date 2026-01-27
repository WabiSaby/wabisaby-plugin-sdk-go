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

// Song represents a song with metadata.
type Song map[string]interface{}

// SongClient provides access to song operations.
type SongClient struct {
	tenantID uuid.UUID
	pluginID uuid.UUID
	client   pluginpb.PluginCapabilitiesServiceClient
}

// NewSongClient creates a new song client.
func NewSongClient(tenantID, pluginID uuid.UUID, client pluginpb.PluginCapabilitiesServiceClient) *SongClient {
	return &SongClient{
		tenantID: tenantID,
		pluginID: pluginID,
		client:   client,
	}
}

// Search searches for songs by query.
// Searches by title, artist, or channel.
func (c *SongClient) Search(ctx context.Context, query string, limit int) ([]Song, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	req := &pluginpb.SongSearchRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		Query:    query,
		Limit:    int32(limit),
	}

	resp, err := c.client.SongSearch(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("song search failed: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("song error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	if len(resp.Songs) == 0 {
		return []Song{}, nil
	}

	var songs []Song
	if err := json.Unmarshal(resp.Songs, &songs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal songs: %w", err)
	}

	return songs, nil
}

// Get retrieves a song by ID.
// Returns nil if the song doesn't exist.
func (c *SongClient) Get(ctx context.Context, songID string) (Song, error) {
	req := &pluginpb.SongGetRequest{
		TenantId: c.tenantID.String(),
		PluginId: c.pluginID.String(),
		SongId:   songID,
	}

	resp, err := c.client.SongGet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("song get failed: %w", err)
	}

	if resp.Error != nil {
		if resp.Error.Code == "NOT_FOUND" {
			return nil, nil
		}
		return nil, fmt.Errorf("song error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	if len(resp.Song) == 0 {
		return nil, nil
	}

	var song Song
	if err := json.Unmarshal(resp.Song, &song); err != nil {
		return nil, fmt.Errorf("failed to unmarshal song: %w", err)
	}

	return song, nil
}
