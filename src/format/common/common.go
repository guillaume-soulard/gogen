package common

import "github.com/ogama/gogen/src/configuration"

type GeneratedObject struct {
	Object interface{}
}

func (g GeneratedObject) GetValue() (result interface{}) {
	if mapValue, isMap := g.Object.(map[string]interface{}); isMap {

	}
	return result
}

type Format interface {
	Format(generatedObject GeneratedObject) (result string, err error)
}

type Builder interface {
	Build(configuration configuration.FormatConfiguration) (result Format, err error)
}
