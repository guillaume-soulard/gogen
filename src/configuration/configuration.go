package configuration

import (
	"encoding/json"
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
	DateFormat string     `json:"dateFormat"`
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

func (o Options) GetBoolOrDefault(option string, defaultValue bool) (result bool, err error) {
	var optionValue interface{}
	var exists bool
	optionValue, exists = o.Get(option)
	if !exists {
		result = defaultValue
		return result, err
	}
	if optionBool, isBool := optionValue.(bool); isBool {
		result = optionBool
	} else if optionString, isString := optionValue.(string); isString {
		if optionString == "" {
			result = defaultValue
		} else {
			result, err = strconv.ParseBool(optionString)
		}
	}
	return result, err
}

func (o Options) GetStringOrDefault(option string, defaultValue string) (result string, err error) {
	optionValue, exists := o.Get(option)
	if !exists {
		result = defaultValue
		return result, err
	}
	if optionValueString, isString := optionValue.(string); isString && optionValueString != "" {
		result = optionValueString
	} else {
		result = defaultValue
	}
	return result, err
}

func (o Options) GetObjectAsStringOrDefault(option string, defaultValue string) (result string, err error) {
	optionValue, exists := o.Get(option)
	if !exists {
		result = defaultValue
		return result, err
	}
	var bytes []byte
	if bytes, err = json.Marshal(optionValue); err != nil {
		return result, err
	}
	result = string(bytes)
	return result, err
}

func (o Options) Get(option string) (result interface{}, exists bool) {
	pathElement := strings.Split(option, ".")
	var currentObject map[string]interface{}
	currentObject = o
	for i, path := range pathElement {
		if pathValue, exists := currentObject[path]; exists {
			if mapValue, isMap := pathValue.(map[string]interface{}); isMap {
				currentObject = mapValue
			} else if stringValue := pathValue; i == len(pathElement)-1 {
				result = stringValue
				exists = true
				return result, exists
			}
		}
	}
	return result, exists
}

func EmptyOptions() Options {
	return Options{}
}
