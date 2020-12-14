package common

import (
	"github.com/ogama/gogen/src/configuration"
)

type GeneratedObject struct {
	Object interface{}
}

func (g GeneratedObject) GetValue(path []string) (result interface{}, exists bool) {
	result, exists = getValue(g.Object, path)
	return result, exists
}

func getValue(object interface{}, path []string) (result interface{}, exists bool) {
	exists = false
	pathLen := len(path)
	if pathLen > 0 {
		field := path[0]
		if mapValue, isMap := object.(map[string]interface{}); isMap {
			if result, exists = mapValue[field]; exists {
				if pathLen > 1 {
					newPath := path[1:]
					result, exists = getValue(result, newPath)
				}
			}
		}
	}
	return result, exists
}

type Format interface {
	Format(generatedObject GeneratedObject) (result string, err error)
}

type Builder interface {
	Build(configuration configuration.FormatConfiguration) (result Format, err error)
}
