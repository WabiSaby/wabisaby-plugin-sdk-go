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

// Plugin is the base interface that all plugins must implement.
type Plugin interface {
	// Initialize is called when the plugin is first loaded.
	// Use this to set up any initial state or resources.
	Initialize(ctx *Context) error

	// Shutdown is called when the plugin is being unloaded.
	// Use this to clean up any resources.
	Shutdown(ctx *Context) error
}

// CommandPlugin handles command execution (stateless plugins).
type CommandPlugin interface {
	Plugin

	// ExecuteCommand executes a command with the given arguments.
	// Returns the result as a JSON-serializable value.
	ExecuteCommand(ctx *Context, command string, args []interface{}) (interface{}, error)
}

// BasePlugin provides default implementations for all plugin interfaces.
// Plugins can embed BasePlugin to automatically satisfy all interfaces
// and only override the methods they need.
type BasePlugin struct {
	router *CommandRouter
}

// NewBasePlugin creates a new BasePlugin instance.
func NewBasePlugin() *BasePlugin {
	return &BasePlugin{
		router: NewCommandRouter(),
	}
}

// Initialize is called when the plugin is first loaded.
// Override this method to set up initial state or resources.
func (p *BasePlugin) Initialize(ctx *Context) error {
	return nil
}

// Shutdown is called when the plugin is being unloaded.
// Override this method to clean up any resources.
func (p *BasePlugin) Shutdown(ctx *Context) error {
	return nil
}

// ExecuteCommand executes a command with the given arguments.
// Routes to registered command handlers if available, otherwise returns an error.
// Override this method to provide custom command routing logic.
func (p *BasePlugin) ExecuteCommand(ctx *Context, command string, args []interface{}) (interface{}, error) {
	if p.router == nil {
		p.router = NewCommandRouter()
	}
	return p.router.Route(ctx, command, args)
}

// RegisterCommand registers a command handler with the router.
// This is a convenience method for plugins embedding BasePlugin.
func (p *BasePlugin) RegisterCommand(name string, handler CommandHandler, opts ...CommandOption) error {
	if p.router == nil {
		p.router = NewCommandRouter()
	}
	return p.router.Register(name, handler, opts...)
}

// GetCommands returns metadata for all registered commands.
func (p *BasePlugin) GetCommands() []CommandMetadata {
	if p.router == nil {
		return []CommandMetadata{}
	}
	return p.router.GetCommands()
}
