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

// ParamType represents the type of a command parameter.
type ParamType string

// Parameter types.
const (
	ParamTypeString ParamType = "string"
	ParamTypeInt    ParamType = "int"
	ParamTypeFloat  ParamType = "float"
	ParamTypeBool   ParamType = "bool"
	ParamTypeObject ParamType = "object"
	ParamTypeArray  ParamType = "array"
)

// CommandMetadata describes a command with its parameters and return type.
type CommandMetadata struct {
	Name        string
	Description string
	Parameters  []ParameterMetadata
	ReturnType  *ReturnTypeMetadata
	Examples    []CommandExample
}

// ParameterMetadata describes a command parameter.
type ParameterMetadata struct {
	Name        string
	Type        ParamType
	Description string
	Required    bool
	Default     interface{}
}

// ReturnTypeMetadata describes the return type of a command.
type ReturnTypeMetadata struct {
	Name        string
	Description string
	Schema      map[string]ParamType
}

// CommandExample provides a usage example for a command.
type CommandExample struct {
	Description string
	Args        interface{}
	Result      interface{}
}

// CommandOption is a functional option for configuring command metadata.
type CommandOption func(*CommandMetadata)

// WithDescription sets the command description.
func WithDescription(desc string) CommandOption {
	return func(m *CommandMetadata) {
		m.Description = desc
	}
}

// WithParameters sets the command parameters.
func WithParameters(params ...ParameterMetadata) CommandOption {
	return func(m *CommandMetadata) {
		m.Parameters = params
	}
}

// WithReturnType sets the command return type metadata.
func WithReturnType(name string, schema map[string]ParamType) CommandOption {
	return func(m *CommandMetadata) {
		m.ReturnType = &ReturnTypeMetadata{
			Name:   name,
			Schema: schema,
		}
	}
}

// WithExamples sets usage examples for the command.
func WithExamples(examples ...CommandExample) CommandOption {
	return func(m *CommandMetadata) {
		m.Examples = examples
	}
}

// ParamOption is a functional option for configuring parameters.
type ParamOption func(*ParameterMetadata)

// Param creates a parameter definition with the given options.
func Param(name string, paramType ParamType, desc string, opts ...ParamOption) ParameterMetadata {
	p := ParameterMetadata{
		Name:        name,
		Type:        paramType,
		Description: desc,
		Required:    true, // default to required
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// Required marks a parameter as required.
func Required() ParamOption {
	return func(p *ParameterMetadata) {
		p.Required = true
	}
}

// Optional marks a parameter as optional.
func Optional() ParamOption {
	return func(p *ParameterMetadata) {
		p.Required = false
	}
}

// Default sets a default value for an optional parameter.
func Default(val interface{}) ParamOption {
	return func(p *ParameterMetadata) {
		p.Required = false
		p.Default = val
	}
}
