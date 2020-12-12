package model

import (
	"strings"
)

type TypeOptions map[string]interface{}

func BuildTypeOption(json map[string]interface{}) (result TypeOptions) {
	return getOptionsRecursively(json, []string{}, TypeOptions{})
}

func getOptionsRecursively(json map[string]interface{}, parentPath []string, options TypeOptions) TypeOptions {
	for key, value := range json {
		newParentPath := append(parentPath, key)
		if key == "itemTemplate" {
			options[strings.Join(newParentPath, ".")] = value
		} else {
			if fieldMap, isMap := value.(map[string]interface{}); isMap {
				options = getOptionsRecursively(fieldMap, newParentPath, options)
			} else {
				options[strings.Join(newParentPath, ".")] = value
			}
		}
	}
	return options
}

func (to TypeOptions) GetOptionAsString(optionPath string) string {
	if value, exists := to[optionPath]; exists {
		return value.(string)
	} else {
		return ""
	}
}

func (to TypeOptions) GetOptionAsInt(optionPath string) int {
	switch to[optionPath].(type) {
	case float64:
		return int(to[optionPath].(float64))
	case nil:
		return 0
	default:
		return to[optionPath].(int)
	}
}

func (to TypeOptions) GetOptionAsBool(optionPath string) bool {
	switch to[optionPath].(type) {
	case nil:
		return false
	default:
		return to[optionPath].(bool)
	}
}

func (to TypeOptions) GetOptionAsFloat(optionPath string) float64 {
	switch to[optionPath].(type) {
	case int:
		return float64(to[optionPath].(int))
	case nil:
		return 0
	default:
		return to[optionPath].(float64)
	}
}

func (to TypeOptions) GetOptionAsInterface(optionPath string) interface{} {
	return to[optionPath]
}

func (to TypeOptions) OverloadWith(options TypeOptions) TypeOptions {
	newOptions := TypeOptions{}
	for k, v := range to {
		newOptions[k] = v
	}
	for k, v := range options {
		newOptions[k] = v
	}
	return newOptions
}

func (to TypeOptions) Add(path string, value interface{}) {
	to[path] = value
}
