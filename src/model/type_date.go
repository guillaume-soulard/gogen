package model

import (
	"errors"
	"fmt"
	"github.com/ogama/gogen/src/constants"
	"strings"
	"time"
)

type DateType struct {
	min      time.Time
	max      time.Time
	truncate string
	name     string
}

func (d *DateType) Generate(context *GeneratorContext, _ GenerationRequest) (result interface{}, err error) {
	date := time.Unix(context.GenerateInteger64Between(d.min.Unix(), d.max.Unix()), 0).In(time.UTC)
	var duration time.Duration
	if duration, err = getDurationFrom(d.truncate); err != nil {
		return result, err
	}
	date = date.Truncate(duration)
	date = date.In(time.UTC)
	return date, err
}

func (d *DateType) GetName() string {
	return ""
}

type DateTypeFactory struct{}

func (d DateTypeFactory) DefaultOptions() TypeOptions {
	defaultOptions := TypeOptions{}
	defaultOptions.Add("bounds.min", "1970-01-01T00:00:00")
	defaultOptions.Add("bounds.max", "2099-12-31T23:59:59")
	defaultOptions.Add("truncate", "milliseconds")
	defaultOptions.Add("name", "")
	return defaultOptions
}

func (d DateTypeFactory) New(parameters TypeFactoryParameter) (generator TypeGenerator, err error) {
	var min, max time.Time
	if min, err = time.ParseInLocation(constants.DateFormat, parameters.Options.GetOptionAsString("bounds.min"), time.UTC); err != nil {
		return generator, err
	}
	if max, err = time.ParseInLocation(constants.DateFormat, parameters.Options.GetOptionAsString("bounds.max"), time.UTC); err != nil {
		return generator, err
	}
	if min.After(max) {
		return generator, errors.New(fmt.Sprintf("bounds.min = %s is greater than bounds.max = %s", min.Format(constants.DateFormat), max.Format(constants.DateFormat)))
	}
	return &DateType{
		min:      min,
		max:      max,
		truncate: parameters.Options.GetOptionAsString("truncate"),
		name:     parameters.Options.GetOptionAsString("name"),
	}, err
}

func getDurationFrom(truncation string) (result time.Duration, err error) {
	switch strings.ToLower(truncation) {
	case "milliseconds":
		result = time.Millisecond
	case "seconds":
		result = time.Second
	case "minutes":
		result = time.Minute
	case "hours":
		result = time.Hour
	default:
		err = errors.New(fmt.Sprintf("unsupported truncate %s", truncation))
	}
	return result, err
}
