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
	pluginpb "github.com/WabiSaby/WabiSaby-Protos/go/plugin"
)

// LogLevel represents a log level.
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// Logger provides structured logging for plugins.
type Logger struct {
	tenantID uuid.UUID
	pluginID uuid.UUID
	client   pluginpb.PluginCapabilitiesServiceClient
}

// NewLogger creates a new logger.
func NewLogger(tenantID, pluginID uuid.UUID, client pluginpb.PluginCapabilitiesServiceClient) *Logger {
	return &Logger{
		tenantID: tenantID,
		pluginID: pluginID,
		client:   client,
	}
}

// log sends a log message to the core.
func (l *Logger) log(ctx context.Context, level LogLevel, message string, fields map[string]string) error {
	req := &pluginpb.LogRequest{
		TenantId: l.tenantID.String(),
		PluginId: l.pluginID.String(),
		Level:    string(level),
		Message:  message,
		Fields:   fields,
	}

	resp, err := l.client.Log(ctx, req)
	if err != nil {
		return fmt.Errorf("log failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("log error: %s - %s", resp.Error.Code, resp.Error.Message)
	}

	return nil
}

// Debug logs a debug message.
func (l *Logger) Debug(ctx context.Context, message string, fields ...map[string]string) error {
	mergedFields := mergeFields(fields...)
	return l.log(ctx, LogLevelDebug, message, mergedFields)
}

// Info logs an info message.
func (l *Logger) Info(ctx context.Context, message string, fields ...map[string]string) error {
	mergedFields := mergeFields(fields...)
	return l.log(ctx, LogLevelInfo, message, mergedFields)
}

// Warn logs a warning message.
func (l *Logger) Warn(ctx context.Context, message string, fields ...map[string]string) error {
	mergedFields := mergeFields(fields...)
	return l.log(ctx, LogLevelWarn, message, mergedFields)
}

// Error logs an error message.
func (l *Logger) Error(ctx context.Context, message string, fields ...map[string]string) error {
	mergedFields := mergeFields(fields...)
	return l.log(ctx, LogLevelError, message, mergedFields)
}

// mergeFields merges multiple field maps into one.
func mergeFields(fields ...map[string]string) map[string]string {
	if len(fields) == 0 {
		return nil
	}

	merged := make(map[string]string)
	for _, f := range fields {
		for k, v := range f {
			merged[k] = v
		}
	}

	if len(merged) == 0 {
		return nil
	}

	return merged
}
