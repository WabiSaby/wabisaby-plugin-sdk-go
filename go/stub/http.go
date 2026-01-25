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
	"github.com/wabisaby/wabisaby/pkg/utils/convertion"
	"time"

	"github.com/google/uuid"
	pluginpb "github.com/wabisaby/wabisaby/api/generated/proto/plugin"
)

// HTTPResponse represents an HTTP response.
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// HTTPClient provides access to HTTP operations.
type HTTPClient struct {
	tenantID uuid.UUID
	pluginID uuid.UUID
	client   pluginpb.PluginCapabilitiesServiceClient
}

// NewHTTPClient creates a new HTTP client.
func NewHTTPClient(tenantID, pluginID uuid.UUID, client pluginpb.PluginCapabilitiesServiceClient) *HTTPClient {
	return &HTTPClient{
		tenantID: tenantID,
		pluginID: pluginID,
		client:   client,
	}
}

// Fetch makes an HTTP request.
// method is the HTTP method (GET, POST, etc.).
// url is the request URL.
// headers is an optional map of HTTP headers.
// body is an optional request body.
// timeout is an optional timeout in milliseconds (0 = use default).
func (c *HTTPClient) Fetch(
	ctx context.Context,
	method, url string,
	headers map[string]string,
	body []byte,
	timeoutMs int32,
) (*HTTPResponse, error) {
	req := &pluginpb.HTTPFetchRequest{
		TenantId:  c.tenantID.String(),
		PluginId:  c.pluginID.String(),
		Url:       url,
		Method:    method,
		Headers:   headers,
		Body:      body,
		TimeoutMs: timeoutMs,
	}

	resp, err := c.client.HTTPFetch(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("HTTP fetch failed: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("HTTP error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return &HTTPResponse{
		StatusCode: int(resp.StatusCode),
		Headers:    resp.Headers,
		Body:       resp.Body,
	}, nil
}

// Get is a convenience method for GET requests.
func (c *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error) {
	return c.Fetch(ctx, "GET", url, headers, nil, 0)
}

// Post is a convenience method for POST requests.
func (c *HTTPClient) Post(ctx context.Context, url string, headers map[string]string, body []byte) (*HTTPResponse, error) {
	return c.Fetch(ctx, "POST", url, headers, body, 0)
}

// PostWithTimeout is a convenience method for POST requests with timeout.
func (c *HTTPClient) PostWithTimeout(ctx context.Context, url string, headers map[string]string, body []byte, timeout time.Duration) (*HTTPResponse, error) {
	timeoutMs := convertion.DurationToInt32Ms(timeout)
	return c.Fetch(ctx, "POST", url, headers, body, timeoutMs)
}
