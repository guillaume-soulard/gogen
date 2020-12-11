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
		"ref":      RefTypeFactory{},
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
	Generate(context *GeneratorContext, request GenerationRequest) (result interface{}, err error)
	GetName() string
}

type AbstractTypeGenerator struct {
	lastGenerated   interface{}
	generationCount int64
	lastRequest     GenerationRequest
	TypeGenerator   TypeGenerator
	IsRef           bool
}

type GenerationRequest string

const (
	Generate GenerationRequest = "generate"
	Ref      GenerationRequest = "ref"
)

func (a AbstractTypeGenerator) Generate(context *GeneratorContext, request GenerationRequest) (result interface{}, err error) {
	sdfgsdg
	if a.IsRef {
		result, err = a.TypeGenerator.Generate(context, request)
	} else if request == Generate {
		if a.generationCount == 1 && a.lastRequest == Ref {
			result = a.lastGenerated
		} else {
			result, err = a.TypeGenerator.Generate(context, request)
			a.lastGenerated = result
		}
	} else if request == Ref {
		if a.generationCount == 0 {
			result, err = a.TypeGenerator.Generate(context, request)
			a.lastGenerated = result
		}
		result = a.lastGenerated
	}
	a.lastRequest = request
	a.generationCount++
	return result, err
}

func (a AbstractTypeGenerator) GetName() string {
	return a.TypeGenerator.GetName()
}
