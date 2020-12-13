package common

import (
	"github.com/ogama/gogen/src/configuration"
	"strings"
)

type GeneratedObject struct {
	Object interface{}
}

func (g GeneratedObject) GetValue(path string) (result interface{}, exists bool) {
	result, exists = g.getValuePathArray(strings.Split(path, "."))
	return result, exists
}

func (g GeneratedObject) getValuePathArray(pathArray []string) (result interface{}, exists bool) {
	if len(pathArray) > 0 {
		field := pathArray[0]
		if mapValue, isMap := g.Object.(map[string]interface{}); isMap {
			if result, exists = mapValue[field]; exists {

			}
		}
	} else {
		exists = false
	}
	return result, exists
}

type Format interface {
	Format(generatedObject GeneratedObject) (result string, err error)
}

type Builder interface {
	Build(configuration configuration.FormatConfiguration) (result Format, err error)
}
