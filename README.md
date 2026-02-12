# WabiSaby Plugin SDK

The WabiSaby Plugin SDK provides a Go library for developing plugins that extend the WabiSaby platform. Plugins can provide storage providers, content resolvers, downloaders, and more.

## Installation

```bash
go get github.com/wabisaby/wabisaby-plugin-sdk
```

For local development, use a replace directive in your `go.mod`:

```go
replace github.com/wabisaby/wabisaby-plugin-sdk => ../WabiSaby-Plugin-SDK
```

## Quick Start

```go
package main

import (
    sdk "github.com/wabisaby/wabisaby-plugin-sdk"
)

type MyPlugin struct {
    *sdk.BasePlugin
}

func (p *MyPlugin) Initialize(ctx *sdk.Context) error {
    ctx.Logger.Info("Plugin initialized")
    return nil
}

func (p *MyPlugin) ExecuteCommand(ctx *sdk.Context, command string, args []interface{}) (interface{}, error) {
    switch command {
    case "hello":
        return map[string]interface{}{"message": "Hello, World!"}, nil
    default:
        return nil, fmt.Errorf("unknown command: %s", command)
    }
}

func main() {
    plugin := &MyPlugin{
        BasePlugin: sdk.NewBasePlugin(),
    }
    sdk.Serve(plugin)
}
```

## Core Interfaces

### Plugin
The base interface that all plugins must implement:

```go
type Plugin interface {
    Initialize(ctx *Context) error
    Shutdown(ctx *Context) error
}
```

### CommandPlugin
For plugins that handle commands:

```go
type CommandPlugin interface {
    Plugin
    ExecuteCommand(ctx *Context, command string, args []interface{}) (interface{}, error)
}
```

### StorageProvider
For storage provider plugins:

```go
type StorageProvider interface {
    Plugin
    UploadAudio(ctx *Context, req *UploadAudioRequest) (string, error)
    GetFileSizeMB(ctx *Context, cdnURL string) (float64, error)
    DeleteAudio(ctx *Context, cdnURL string) error
}
```

## Context

The `Context` provides access to all plugin capabilities:

```go
type Context struct {
    context.Context
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
}
```

## Examples

### Using Storage

```go
func (p *MyPlugin) UploadFile(ctx *sdk.Context, filePath string) error {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }
    
    return ctx.Storage.Set(ctx, "my-key", data)
}
```

### Making HTTP Requests

```go
func (p *MyPlugin) FetchData(ctx *sdk.Context, url string) ([]byte, error) {
    resp, err := ctx.HTTP.Fetch(ctx, &sdk.HTTPRequest{
        URL:    url,
        Method: "GET",
    })
    if err != nil {
        return nil, err
    }
    
    return resp.Body, nil
}
```

### Sending Notifications

```go
func (p *MyPlugin) NotifyUser(ctx *sdk.Context, message string) error {
    return ctx.Notification.Send(ctx, &sdk.NotificationRequest{
        Title:   "Plugin Notification",
        Message: message,
        Type:    "info",
    })
}
```

## Dependencies

The SDK depends on:
- `github.com/wabisaby/wabisaby/api/generated/proto/plugin` - Protobuf definitions (via replace directive for local dev)
- `google.golang.org/grpc` - gRPC library
- `github.com/hashicorp/go-plugin` - Plugin framework

## License

Copyright (c) 2026 WabiSaby. All rights reserved.

This source code is proprietary and confidential. Unauthorized copying, modification, distribution, or use of this software, via any medium is strictly prohibited without the express written permission of WabiSaby.
