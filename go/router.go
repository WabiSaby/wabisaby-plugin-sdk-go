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
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

// CommandHandler represents a command handler function.
// Supported signatures:
//   - func(ctx *Context) (interface{}, error)
//   - func(ctx *Context, args *T) (interface{}, error)
//   - func(ctx *Context, args *T) (*R, error)
type CommandHandler interface{}

// registeredCommand holds a registered command with its metadata and handler.
type registeredCommand struct {
	metadata CommandMetadata
	handler  CommandHandler
	argType  reflect.Type // nil if handler takes no args beyond Context
}

// CommandRouter manages command registration and routing.
type CommandRouter struct {
	mu       sync.RWMutex
	commands map[string]*registeredCommand
}

// NewCommandRouter creates a new command router.
func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		commands: make(map[string]*registeredCommand),
	}
}

// Register registers a command with its handler and options.
// Returns an error if the command is already registered or the handler signature is invalid.
func (r *CommandRouter) Register(name string, handler CommandHandler, opts ...CommandOption) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.commands[name]; exists {
		return fmt.Errorf("command %q already registered", name)
	}

	// Validate handler signature
	handlerVal := reflect.ValueOf(handler)
	handlerType := handlerVal.Type()

	if handlerType.Kind() != reflect.Func {
		return fmt.Errorf("handler must be a function, got %s", handlerType.Kind())
	}

	// Must have at least 1 input (*Context) and 2 outputs (result, error)
	if handlerType.NumIn() < 1 || handlerType.NumIn() > 2 {
		return fmt.Errorf("handler must have 1 or 2 inputs, got %d", handlerType.NumIn())
	}
	if handlerType.NumOut() != 2 {
		return fmt.Errorf("handler must have 2 outputs (result, error), got %d", handlerType.NumOut())
	}

	// First input must be *Context
	ctxType := handlerType.In(0)
	if ctxType.Kind() != reflect.Ptr || ctxType.Elem().Name() != "Context" {
		return fmt.Errorf("first handler argument must be *Context, got %s", ctxType)
	}

	// Second output must implement error
	errType := handlerType.Out(1)
	if !errType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return fmt.Errorf("second handler return must be error, got %s", errType)
	}

	// Build metadata
	metadata := CommandMetadata{Name: name}
	for _, opt := range opts {
		opt(&metadata)
	}

	cmd := &registeredCommand{
		metadata: metadata,
		handler:  handler,
	}

	// Detect argument type if handler takes typed args
	if handlerType.NumIn() == 2 {
		argType := handlerType.In(1)
		if argType.Kind() == reflect.Ptr {
			cmd.argType = argType.Elem()
		} else {
			cmd.argType = argType
		}
	}

	r.commands[name] = cmd
	return nil
}

// Route routes a command to its handler and returns the result.
func (r *CommandRouter) Route(ctx *Context, command string, args []interface{}) (interface{}, error) {
	r.mu.RLock()
	cmd, exists := r.commands[command]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown command: %s", command)
	}

	return r.invokeHandler(ctx, cmd, args)
}

// invokeHandler calls the handler with proper argument marshaling.
func (r *CommandRouter) invokeHandler(ctx *Context, cmd *registeredCommand, args []interface{}) (interface{}, error) {
	handlerVal := reflect.ValueOf(cmd.handler)
	handlerType := handlerVal.Type()

	var callArgs []reflect.Value
	callArgs = append(callArgs, reflect.ValueOf(ctx))

	// If handler expects typed arguments, unmarshal them
	if cmd.argType != nil && handlerType.NumIn() == 2 {
		argPtr := reflect.New(cmd.argType)

		if len(args) > 0 {
			var argsMap map[string]interface{}

			// Check if first arg is already a map
			if m, ok := args[0].(map[string]interface{}); ok {
				argsMap = m
			} else {
				// Construct map from positional args using parameter metadata
				argsMap = make(map[string]interface{})
				for i, param := range cmd.metadata.Parameters {
					if i < len(args) {
						argsMap[param.Name] = args[i]
					} else if param.Default != nil {
						argsMap[param.Name] = param.Default
					}
				}
			}

			// Marshal to JSON then unmarshal to typed struct
			jsonBytes, err := json.Marshal(argsMap)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal arguments: %w", err)
			}
			if err := json.Unmarshal(jsonBytes, argPtr.Interface()); err != nil {
				return nil, fmt.Errorf("failed to unmarshal arguments to %s: %w", cmd.argType.Name(), err)
			}
		}

		// Pass pointer to the handler
		callArgs = append(callArgs, argPtr)
	}

	// Call the handler
	results := handlerVal.Call(callArgs)

	// Extract return values
	var result interface{}
	var err error

	if len(results) >= 1 && !results[0].IsNil() {
		result = results[0].Interface()
	}
	if len(results) >= 2 && !results[1].IsNil() {
		if errVal, ok := results[1].Interface().(error); ok {
			err = errVal
		}
	}

	return result, err
}

// GetCommands returns metadata for all registered commands.
func (r *CommandRouter) GetCommands() []CommandMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commands := make([]CommandMetadata, 0, len(r.commands))
	for _, cmd := range r.commands {
		commands = append(commands, cmd.metadata)
	}
	return commands
}

// HasCommand checks if a command is registered.
func (r *CommandRouter) HasCommand(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.commands[name]
	return exists
}

// CommandCount returns the number of registered commands.
func (r *CommandRouter) CommandCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.commands)
}
