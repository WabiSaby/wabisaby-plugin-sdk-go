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
	"fmt"

	"github.com/wabisaby/wabisaby-plugin-sdk/go/stub"
)

// NotificationType represents the type of notification.
// Re-exported from stub for convenience.
type NotificationType = stub.NotificationType

// Notification type constants.
const (
	NotificationTypeInfo    NotificationType = stub.NotificationTypeInfo
	NotificationTypeSuccess NotificationType = stub.NotificationTypeSuccess
	NotificationTypeWarning NotificationType = stub.NotificationTypeWarning
	NotificationTypeError   NotificationType = stub.NotificationTypeError
)

// HTTPRequest represents an HTTP request configuration.
type HTTPRequest struct {
	URL       string
	Method    string
	Headers   map[string]string
	Body      []byte
	TimeoutMs int32
}

// HTTPResponse is an alias to stub.HTTPResponse for convenience.
type HTTPResponse = stub.HTTPResponse

// ConfigAccessor provides typed access to plugin configuration.
type ConfigAccessor struct {
	data map[string]interface{}
}

// NewConfigAccessor creates a new config accessor.
func NewConfigAccessor(data map[string]interface{}) *ConfigAccessor {
	return &ConfigAccessor{data: data}
}

// Get returns a raw config value.
func (c *ConfigAccessor) Get(key string) interface{} {
	if c == nil || c.data == nil {
		return nil
	}
	return c.data[key]
}

// GetString returns a string config value.
func (c *ConfigAccessor) GetString(key string, defaultVal ...string) string {
	val := c.Get(key)
	if s, ok := val.(string); ok {
		return s
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return ""
}

// GetInt returns an integer config value.
func (c *ConfigAccessor) GetInt(key string, defaultVal ...int) int {
	val := c.Get(key)
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0
}

// GetBool returns a boolean config value.
func (c *ConfigAccessor) GetBool(key string, defaultVal ...bool) bool {
	val := c.Get(key)
	if b, ok := val.(bool); ok {
		return b
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return false
}

// GetFloat returns a float64 config value.
func (c *ConfigAccessor) GetFloat(key string, defaultVal ...float64) float64 {
	val := c.Get(key)
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0.0
}

// Has checks if a config key exists.
func (c *ConfigAccessor) Has(key string) bool {
	if c == nil || c.data == nil {
		return false
	}
	_, exists := c.data[key]
	return exists
}

// ContextLogger provides slog-style logging for plugins.
type ContextLogger struct {
	logger *stub.Logger
	ctx    *Context
}

// NewContextLogger creates a new context logger.
func NewContextLogger(logger *stub.Logger, ctx *Context) *ContextLogger {
	return &ContextLogger{logger: logger, ctx: ctx}
}

// Info logs an info message with key-value pairs.
func (l *ContextLogger) Info(msg string, keysAndValues ...interface{}) {
	fields := toStringMap(keysAndValues...)
	if err := l.logger.Info(l.ctx, msg, fields); err != nil {
		// Logging errors are non-critical, ignore
		_ = err
	}
}

// Debug logs a debug message with key-value pairs.
func (l *ContextLogger) Debug(msg string, keysAndValues ...interface{}) {
	fields := toStringMap(keysAndValues...)
	if err := l.logger.Debug(l.ctx, msg, fields); err != nil {
		// Logging errors are non-critical, ignore
		_ = err
	}
}

// Warn logs a warning message with key-value pairs.
func (l *ContextLogger) Warn(msg string, keysAndValues ...interface{}) {
	fields := toStringMap(keysAndValues...)
	if err := l.logger.Warn(l.ctx, msg, fields); err != nil {
		// Logging errors are non-critical, ignore
		_ = err
	}
}

// Error logs an error message with key-value pairs.
func (l *ContextLogger) Error(msg string, keysAndValues ...interface{}) {
	fields := toStringMap(keysAndValues...)
	if err := l.logger.Error(l.ctx, msg, fields); err != nil {
		// Logging errors are non-critical, ignore
		_ = err
	}
}

// toStringMap converts key-value pairs to map[string]string.
func toStringMap(keysAndValues ...interface{}) map[string]string {
	if len(keysAndValues) == 0 {
		return nil
	}

	result := make(map[string]string)
	for i := 0; i < len(keysAndValues)-1; i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keysAndValues[i])
		}
		result[key] = fmt.Sprintf("%v", keysAndValues[i+1])
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

// HTTPService wraps HTTPClient with a more ergonomic API.
type HTTPService struct {
	client *stub.HTTPClient
}

// NewHTTPService creates a new HTTP service wrapper.
func NewHTTPService(client *stub.HTTPClient) HTTPService {
	return HTTPService{client: client}
}

// Fetch makes an HTTP request using HTTPRequest struct.
func (h HTTPService) Fetch(ctx *Context, req *HTTPRequest) (*HTTPResponse, error) {
	method := req.Method
	if method == "" {
		method = "GET"
	}
	return h.client.Fetch(ctx, method, req.URL, req.Headers, req.Body, req.TimeoutMs)
}

// Get makes a GET request.
func (h HTTPService) Get(ctx *Context, url string) (*HTTPResponse, error) {
	return h.client.Get(ctx, url, nil)
}

// GetWithHeaders makes a GET request with custom headers.
func (h HTTPService) GetWithHeaders(ctx *Context, url string, headers map[string]string) (*HTTPResponse, error) {
	return h.client.Get(ctx, url, headers)
}

// Post makes a POST request with JSON body.
func (h HTTPService) Post(ctx *Context, url string, body []byte) (*HTTPResponse, error) {
	return h.client.Post(ctx, url, map[string]string{"Content-Type": "application/json"}, body)
}

// PostWithHeaders makes a POST request with custom headers.
func (h HTTPService) PostWithHeaders(ctx *Context, url string, headers map[string]string, body []byte) (*HTTPResponse, error) {
	return h.client.Post(ctx, url, headers, body)
}

// NotificationService wraps NotificationClient with a simpler API.
type NotificationService struct {
	client *stub.NotificationClient
}

// NewNotificationService creates a new notification service wrapper.
func NewNotificationService(client *stub.NotificationClient) NotificationService {
	return NotificationService{client: client}
}

// Send sends a notification to all users in the tenant.
func (n NotificationService) Send(ctx *Context, message string, notifType NotificationType) error {
	_, err := n.client.Send(ctx, "", "Plugin Notification", message, notifType)
	return err
}

// SendWithTitle sends a notification with a custom title.
func (n NotificationService) SendWithTitle(ctx *Context, title, message string, notifType NotificationType) error {
	_, err := n.client.Send(ctx, "", title, message, notifType)
	return err
}

// SendToUser sends a notification to a specific user.
func (n NotificationService) SendToUser(ctx *Context, userID, title, message string, notifType NotificationType) error {
	_, err := n.client.Send(ctx, userID, title, message, notifType)
	return err
}
