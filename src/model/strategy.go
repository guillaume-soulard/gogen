package model

import (
	"github.com/ogama/gogen/src/configuration"
)

var TypeFactories = TypeGeneratorFactory{
	factories: map[string]TypeFactory{
		"number":   NumberTypeFactory{},
		"array":    ArrayTypeFactory{},
		"string":   StringTypeFactory{},
		"date":     DateTypeFactory{},
		"picker":   PickerTypeFactory{},
		"object":   ObjectTypeFactory{},
		"constant": ConstantTypeFactory{},
	},
}

type TypeGeneratorFactory struct {
	factories map[string]TypeFactory
}

func (tf TypeGeneratorFactory) GetFactory(typeName string) (TypeFactory, bool) {
	if factory, exists := tf.factories[typeName]; exists {
		return factory, true
	} else {
		return factory, false
	}
}

type TypeFactoryParameter struct {
	Configuration configuration.Configuration
	Template      interface{}
	Options       TypeOptions
}

type TypeGenerator interface {
	Generate(context *GeneratorContext) (result interface{}, err error)
	GetName() string
}
