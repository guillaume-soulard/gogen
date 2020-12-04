package configuration

import (
	"strconv"
	"strings"
)

type Configuration struct {
	TemplateConfiguration interface{}           `json:"template"`
	Options               GlobalOptions         `json:"options"`
	OutputConfigurations  []OutputConfiguration `json:"outputs"`
}

type GlobalOptions struct {
	Seed       string `json:"seed"`
	Amount     int    `json:"amount"`
	Skip       int    `json:"skip"`
	Alias      map[string]Configuration
	Generation Generation `json:"generation"`
}

type Generation struct {
	Type    string                `json:"type"`
	Options GenerationOptions     `json:"options"`
	Format  []OutputConfiguration `json:"outputs"`
}

type OutputConfiguration struct {
	Type                string              `json:"type"`
	Options             Options             `json:"options"`
	FormatConfiguration FormatConfiguration `json:"format"`
}

type FormatConfiguration struct {
	Type    string  `json:"type"`
	Options Options `json:"options"`
}

type GenerationOptions struct {
	Interval int64 `json:"interval"`
}

type Options map[string]interface{}

func (o Options) GetBoolOrDefault(option string) (result bool, err error) {
	var optionValue interface{}
	optionValue = o.Get(option)
	if optionBool, isBool := optionValue.(bool); isBool {
		result = optionBool
	} else if optionString, isString := optionValue.(string); isString {
		result, err = strconv.ParseBool(optionString)
	}
	return result, err
}

func (o Options) Get(option string) (result interface{}) {
	result = ""
	pathElement := strings.Split(option, ".")
	var currentObject map[string]interface{}
	currentObject = o
	for i, path := range pathElement {
		if pathValue, exists := currentObject[path]; exists {
			if mapValue, isMap := pathValue.(map[string]interface{}); isMap {
				currentObject = mapValue
			} else if stringValue := pathValue; i == len(pathElement)-1 {
				result = stringValue
				return result
			}
		}
	}
	return result
}

func EmptyOptions() Options {
	return Options{}
}
