package types

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const dateFormat = "2006-01-02T15:04:05"

type DateType struct {
	min time.Time
	max time.Time
	truncate string
}

func (d *DateType) Generate(context *GeneratorContext) (result interface{}, err error) {
	date := time.Unix(context.GenerateInteger64Between(d.min.Unix(), d.max.Unix()), 0).In(time.UTC)
	var duration time.Duration
	if duration, err = getDurationFrom(d.truncate); err != nil {
		return result, err
	}
	date = date.Truncate(duration)
	return date, err
}

type DateTypeFactory struct {}

func (d DateTypeFactory) DefaultOptions() TypeOptions {
	defaultOptions := TypeOptions{}
	defaultOptions.Add("bounds.min", "1970-01-01T00:00:00")
	defaultOptions.Add("bounds.max", "2099-12-31T23:59:59")
	defaultOptions.Add("truncate", "milliseconds")
	return defaultOptions
}

func (d DateTypeFactory) New(parameters TypeFactoryParameter) (generator TypeGenerator, err error) {
	var min, max time.Time
	if min, err = time.ParseInLocation(dateFormat, parameters.Options.GetOptionAsString("bounds.min"), time.UTC); err != nil {
		return generator, err
	}
	if max, err = time.ParseInLocation(dateFormat, parameters.Options.GetOptionAsString("bounds.max"), time.UTC); err != nil {
		return generator, err
	}
	if min.After(max) {
		return generator, errors.New(fmt.Sprintf("bounds.min = %s is greater than bounds.max = %s", min.Format(dateFormat), max.Format(dateFormat)))
	}
	return &DateType{
		min: min,
		max: max,
		truncate: parameters.Options.GetOptionAsString("truncate"),
	}, err
}

func getDurationFrom(truncation string) (result time.Duration, err error) {
	switch strings.ToLower(truncation) {
	case "milliseconds": result = time.Millisecond
	case "seconds": result =  time.Second
	case "minutes": result =  time.Minute
	case "hours": result =  time.Hour
	default: err = errors.New(fmt.Sprintf("unsupported truncate %s", truncation))
	}
	return result, err
}