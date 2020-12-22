package model

import (
	"errors"
	"fmt"
	"github.com/ogama/gogen/src/configuration"
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
	if err = d.truncateDate(&date); err != nil {
		return result, err
	}
	date = date.In(time.UTC)
	result = date.Format(getDateFormatOrDefault(context.Config.Options.DateFormat))
	return result, err
}

func (d *DateType) GetName() string {
	return ""
}

const defaultStartDate = "1970-01-01T00:00:00"
const defaultEndDate = "2099-12-31T23:59:59"
const defaultDateFormat = "2006-01-02T15:04:05"

type DateTypeFactory struct{}

func (d DateTypeFactory) DefaultOptions() TypeOptions {
	defaultOptions := TypeOptions{}
	defaultOptions.Add("bounds.min", defaultStartDate)
	defaultOptions.Add("bounds.max", defaultEndDate)
	defaultOptions.Add("truncate", "milliseconds")
	defaultOptions.Add("name", "")
	return defaultOptions
}

func (d DateTypeFactory) New(parameters TypeFactoryParameter) (generator TypeGenerator, err error) {
	var min, max time.Time
	minDate := parameters.Options.GetOptionAsString("bounds.min")
	maxDate := parameters.Options.GetOptionAsString("bounds.max")
	if min, err = time.ParseInLocation(getFormatToUseForDate(minDate, parameters.Configuration), minDate, time.UTC); err != nil {
		return generator, err
	}
	if max, err = time.ParseInLocation(getFormatToUseForDate(maxDate, parameters.Configuration), maxDate, time.UTC); err != nil {
		return generator, err
	}
	if min.After(max) {
		return generator, errors.New(fmt.Sprintf("bounds.min = %s is greater than bounds.max = %s", min.Format(parameters.Configuration.Options.DateFormat), max.Format(parameters.Configuration.Options.DateFormat)))
	}
	return &DateType{
		min:      min,
		max:      max,
		truncate: parameters.Options.GetOptionAsString("truncate"),
		name:     parameters.Options.GetOptionAsString("name"),
	}, err
}

func getFormatToUseForDate(date string, configuration configuration.Configuration) string {
	if date == defaultStartDate || date == defaultEndDate {
		return defaultDateFormat
	} else {
		return getDateFormatOrDefault(configuration.Options.DateFormat)
	}
}

func getDateFormatOrDefault(dateFormat string) string {
	if dateFormat == "" {
		return defaultDateFormat
	}
	return dateFormat
}

func (d *DateType) truncateDate(date *time.Time) (err error) {
	switch strings.ToLower(d.truncate) {
	case "milliseconds":
		*date = date.Truncate(time.Millisecond)
	case "seconds":
		*date = date.Truncate(time.Second)
	case "minutes":
		*date = date.Truncate(time.Minute)
	case "hours":
		*date = date.Truncate(time.Hour)
	case "days":
		*date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	case "months":
		*date = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	case "years":
		*date = time.Date(date.Year(), time.January, 1, 0, 0, 0, 0, date.Location())
	default:
		err = errors.New(fmt.Sprintf("unsupported truncate %s", d.truncate))
	}
	return err
}
