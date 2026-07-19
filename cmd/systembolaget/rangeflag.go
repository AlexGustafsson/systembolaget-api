package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
)

type RangeFlag[T any] = cli.FlagBase[*Range[T], struct{}, rangeValue[T]]

type Range[T any] struct {
	Minimum T
	Maximum T
}

func (r *Range[T]) String() string {
	return fmt.Sprintf("%v-%v", r.Minimum, r.Maximum)
}

type rangeValue[T any] struct {
	destination **Range[T]
}

// Create implements [cli.ValueCreator].
func (r rangeValue[T]) Create(val *Range[T], p **Range[T], c struct{}) cli.Value {
	*p = val
	return &rangeValue[T]{
		destination: p,
	}
}

// ToString implements [cli.ValueCreator].
func (r rangeValue[T]) ToString(val *Range[T]) string {
	r.destination = &val
	return r.String()
}

// Get implements [cli.Value].
func (r *rangeValue[T]) Get() any {
	return *r.destination
}

// Set implements [cli.Value].
func (r *rangeValue[T]) Set(value string) error {
	minimumString, maximumString, ok := strings.Cut(value, "-")
	if !ok {
		return fmt.Errorf("invalid range format, expected min-max")
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
	default:
		panic("unsupported range type")
	}

	*r.destination = &Range[T]{
		Minimum: minimumValue.Elem().Interface().(T),
		Maximum: maximumValue.Elem().Interface().(T),
	}
	return nil
}

// String implements [cli.Value].
func (r *rangeValue[T]) String() string {
	if r.destination != nil && *r.destination != nil {
		return (*r.destination).String()
	}
	return ""
}
