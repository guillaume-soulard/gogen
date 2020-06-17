package types

import (
	"hash/fnv"
	"math/rand"
	"time"
)

type ObjectModel struct {
	FieldName string
	Fields    []FieldModel
}

func (objectTemplate ObjectModel) Generate(context *GeneratorContext) (result interface{}, err error) {
	generatedObject := make(map[string]interface{})
	for _, field := range objectTemplate.Fields {
		var generated interface{}
		if generated, err = field.Generate(context); err != nil {
			return result, err
		}
		generatedObject[field.FieldName] = generated
	}
	return generatedObject, err
}

type FieldModel struct {
	FieldName string
	Value     TypeGenerator
}

func (fieldTemplate FieldModel) Generate(context *GeneratorContext) (result interface{}, err error) {
	return fieldTemplate.Value.Generate(context)
}

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

type GeneratorContext struct{
	Config Configuration
	Rand *rand.Rand
}

func (gc GeneratorContext) GenerateIntegerBetween(min int, max int) int {
	return min + gc.Rand.Intn(max - min + 1)
}

func (gc GeneratorContext) GenerateInteger64Between(min int64, max int64) int64 {
	return min + gc.Rand.Int63n(max - min + 1)
}

func (gc GeneratorContext) GenerateFloatBetween(min float64, max float64) float64 {
	return min + gc.Rand.Float64() * (max - min)
}

func NewGenerationContext(config Configuration) (result GeneratorContext, err error)  {
	h := fnv.New32a()
	var seed int64
	if config.Options.Seed != "" {
		if _, err = h.Write([]byte(config.Options.Seed)); err != nil {
			return result, err
		}
		seed = int64(h.Sum32())
	} else {
		seed = time.Now().UnixNano()
	}
	result.Rand = rand.New(rand.NewSource(seed))
	result.Config = config
	return result, err
}

type TypeGenerator interface {
	Generate(context *GeneratorContext) (result interface{}, err error)
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

type TypeFactory interface {
	New(parameters TypeFactoryParameter) (generator TypeGenerator, err error)
	DefaultOptions() TypeOptions
}

type TypeFactoryParameter struct {
	Configuration Configuration
	Template      interface{}
	Options       TypeOptions
}
