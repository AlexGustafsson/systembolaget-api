package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

type EnumValue struct {
	Value   string
	Choices []string
}

type EnumFlag struct {
	Name        string
	Usage       string
	Aliases     []string
	Value       EnumValue
	DefaultText string
}

func (f *EnumFlag) Apply(set *flag.FlagSet) error {
	set.Var(&f.Value, f.Name, f.Usage)
	return nil
}

func (f *EnumFlag) Names() []string {
	return cli.FlagNames(f.Name, f.Aliases)
}

func (f *EnumFlag) IsSet() bool {
	return f.Value.Value != ""
}

func (f *EnumFlag) TakesValue() bool {
	return true
}

func (f *EnumFlag) GetUsage() string {
	return f.Usage
}

func (f *EnumFlag) GetValue() string {
	return f.Value.Value
}

func (f *EnumFlag) GetDefaultText() string {
	return f.DefaultText
}

func (f *EnumFlag) GetEnvVars() []string {
	return nil
}

func (f *EnumFlag) String() string {
	return cli.FlagStringer(f) + fmt.Sprintf(" (one of: %s)", strings.Join(f.Value.Choices, ", "))
}

func (v *EnumValue) Get() any {
	return v.Value
}

func (v *EnumValue) String() string {
	return v.Value
}

func (v *EnumValue) Set(value string) error {
	for _, x := range v.Choices {
		if x == value {
			v.Value = value
			return nil
		}
	}
	return fmt.Errorf("not one of %s", strings.Join(v.Choices, ", "))
}
