package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

type Range[T any] struct {
	Minimum T
	Maximum T
}

func (r *Range[T]) String() string {
	return fmt.Sprintf("%v-%v", r.Minimum, r.Maximum)
}

type DateRangeFlag = cli.FlagBase[*Range[time.Time], struct{}, dateRangeValue]

type dateRangeValue struct {
	destination **Range[time.Time]
}

// Create implements [cli.ValueCreator].
func (r dateRangeValue) Create(val *Range[time.Time], p **Range[time.Time], c struct{}) cli.Value {
	*p = val
	return &dateRangeValue{
		destination: p,
	}
}

// ToString implements [cli.ValueCreator].
func (r dateRangeValue) ToString(val *Range[time.Time]) string {
	r.destination = &val
	return r.String()
}

// Get implements [cli.Value].
func (r *dateRangeValue) Get() any {
	return *r.destination
}

// Set implements [cli.Value].
func (r *dateRangeValue) Set(value string) error {
	parts := strings.SplitN(value, "-", 6)
	if len(parts) != 6 {
		return fmt.Errorf("invalid range format, expected min-max")
	}

	minimumString := strings.Join(parts[0:3], "-")
	maximumString := strings.Join(parts[3:6], "-")

	minimum, minimumErr := time.Parse("2006-01-02", minimumString)
	maximum, maximumErr := time.Parse("2006-01-02", maximumString)
	err := errors.Join(minimumErr, maximumErr)
	if err != nil {
		return err
	}

	*r.destination = &Range[time.Time]{
		Minimum: minimum,
		Maximum: maximum,
	}
	return nil
}

// String implements [cli.Value].
func (r *dateRangeValue) String() string {
	if r.destination != nil && *r.destination != nil {
		return (*r.destination).String()
	}
	return ""
}

type IntRangeFlag = cli.FlagBase[*Range[int], struct{}, intRangeFlag]

type intRangeFlag struct {
	destination **Range[int]
}

// Create implements [cli.ValueCreator].
func (r intRangeFlag) Create(val *Range[int], p **Range[int], c struct{}) cli.Value {
	*p = val
	return &intRangeFlag{
		destination: p,
	}
}

// ToString implements [cli.ValueCreator].
func (r intRangeFlag) ToString(val *Range[int]) string {
	r.destination = &val
	return r.String()
}

// Get implements [cli.Value].
func (r *intRangeFlag) Get() any {
	return *r.destination
}

// Set implements [cli.Value].
func (r *intRangeFlag) Set(value string) error {
	minimumString, maximumString, ok := strings.Cut(value, "-")
	if !ok {
		return fmt.Errorf("invalid range format, expected min-max")
	}

	minimum, minimumErr := strconv.ParseInt(minimumString, 10, 32)
	maximum, maximumErr := strconv.ParseInt(maximumString, 10, 32)
	err := errors.Join(minimumErr, maximumErr)
	if err != nil {
		return err
	}

	*r.destination = &Range[int]{
		Minimum: int(minimum),
		Maximum: int(maximum),
	}
	return nil
}

// String implements [cli.Value].
func (r *intRangeFlag) String() string {
	if r.destination != nil && *r.destination != nil {
		return (*r.destination).String()
	}
	return ""
}

type FloatRangeFlag = cli.FlagBase[*Range[float32], struct{}, floatRangeFlag]

type floatRangeFlag struct {
	destination **Range[float32]
}

// Create implements [cli.ValueCreator].
func (r floatRangeFlag) Create(val *Range[float32], p **Range[float32], c struct{}) cli.Value {
	*p = val
	return &floatRangeFlag{
		destination: p,
	}
}

// ToString implements [cli.ValueCreator].
func (r floatRangeFlag) ToString(val *Range[float32]) string {
	r.destination = &val
	return r.String()
}

// Get implements [cli.Value].
func (r *floatRangeFlag) Get() any {
	return *r.destination
}

// Set implements [cli.Value].
func (r *floatRangeFlag) Set(value string) error {
	minimumString, maximumString, ok := strings.Cut(value, "-")
	if !ok {
		return fmt.Errorf("invalid range format, expected min-max")
	}

	minimum, minimumErr := strconv.ParseFloat(minimumString, 32)
	maximum, maximumErr := strconv.ParseFloat(maximumString, 32)
	err := errors.Join(minimumErr, maximumErr)
	if err != nil {
		return err
	}

	*r.destination = &Range[float32]{
		Minimum: float32(minimum),
		Maximum: float32(maximum),
	}
	return nil
}

// String implements [cli.Value].
func (r *floatRangeFlag) String() string {
	if r.destination != nil && *r.destination != nil {
		return (*r.destination).String()
	}
	return ""
}
