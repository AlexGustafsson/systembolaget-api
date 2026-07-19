package main

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

type Range[T any] struct {
	Minimum T
	Maximum T
}

func (r *Range[T]) String() string {
	return fmt.Sprintf("%v-%v", r.Minimum, r.Maximum)
}

type RangeFlag[T any] struct {
	Name        string
	Usage       string
	Aliases     []string
	Value       *Range[T]
	DefaultText string
}

func (f *RangeFlag[T]) Apply(set *flag.FlagSet) error {
	set.Var(f, f.Name, f.Usage)
	return nil
}

func (f *RangeFlag[T]) Names() []string {
	return cli.FlagNames(f.Name, f.Aliases)
}

func (f *RangeFlag[T]) IsSet() bool {
	return f.Value != nil
}

func (f *RangeFlag[T]) TakesValue() bool {
	return true
}

func (f *RangeFlag[T]) GetUsage() string {
	return f.Usage
}

func (f *RangeFlag[T]) GetValue() string {
	return f.Value.String()
}

func (f *RangeFlag[T]) GetDefaultText() string {
	return f.DefaultText
}

func (f *RangeFlag[T]) GetEnvVars() []string {
	return nil
}

func (f *RangeFlag[T]) String() string {
	return cli.FlagStringer(f)
}

func (f *RangeFlag[T]) Get() any {
	return f.Value
}

func (f *RangeFlag[T]) Set(value string) error {
	minimumString, maximumString, ok := strings.Cut(value, ",")
	if !ok {
		return fmt.Errorf("invalid range format, expected min,max")
	}

	minimumValue := reflect.ValueOf(new(T))
	maximumValue := reflect.ValueOf(new(T))

	switch minimumValue.Type().Elem().Kind() {
	case reflect.String:
		minimumValue.Elem().Set(reflect.ValueOf(minimumString))
		maximumValue.Elem().Set(reflect.ValueOf(maximumString))
	case reflect.Int:
		minimum, err := strconv.ParseInt(minimumString, 10, 32)
		if err != nil {
			return err
		}

		maximum, err := strconv.ParseInt(maximumString, 10, 32)
		if err != nil {
			return err
		}

		minimumValue.Elem().Set(reflect.ValueOf(int(minimum)))
		maximumValue.Elem().Set(reflect.ValueOf(int(maximum)))
	case reflect.Float32:
		minimum, err := strconv.ParseFloat(minimumString, 32)
		if err != nil {
			return err
		}

		maximum, err := strconv.ParseFloat(maximumString, 32)
		if err != nil {
			return err
		}

		minimumValue.Elem().Set(reflect.ValueOf(float32(minimum)))
		maximumValue.Elem().Set(reflect.ValueOf(float32(maximum)))
	default:
		return fmt.Errorf("unsupported range type")
	}

	f.Value = &Range[T]{minimumValue.Elem().Interface().(T), maximumValue.Elem().Interface().(T)}
	return nil
}
