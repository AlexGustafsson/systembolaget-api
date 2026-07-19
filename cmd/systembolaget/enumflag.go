package main

import (
	"fmt"
	"slices"

	"github.com/urfave/cli/v3"
)

type EnumConfig struct {
	Choices []string
}

type EnumFlag = cli.FlagBase[string, EnumConfig, enumValue]

var _ cli.ValueCreator[string, EnumConfig] = (*enumValue)(nil)

var _ cli.Value = (*enumValue)(nil)

type enumValue struct {
	destination *string
	config      EnumConfig
}

// Create implements [cli.ValueCreator].
func (e enumValue) Create(val string, p *string, c EnumConfig) cli.Value {
	*p = val
	return &enumValue{
		destination: p,
		config:      c,
	}
}

// ToString implements [cli.ValueCreator].
func (e enumValue) ToString(val string) string {
	e.destination = &val
	return e.String()
}

// Get implements [cli.Value].
func (e *enumValue) Get() any {
	return *e.destination
}

// Set implements [cli.Value].
func (e *enumValue) Set(val string) error {
	if !slices.Contains(e.config.Choices, val) {
		return fmt.Errorf("invalid choice")
	}

	*e.destination = val
	return nil
}

// String implements [cli.Value].
func (e *enumValue) String() string {
	if e.destination != nil && *e.destination != "" {
		return fmt.Sprintf("%q", *e.destination)
	}
	return ""
}
