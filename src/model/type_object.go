package model

import (
	"github.com/ogama/gogen/src/constants"
	"sort"
)

type ObjectType struct {
	fieldTemplates []FieldModel
}

func (o *ObjectType) Generate(context *GeneratorContext, request GenerationRequest) (result interface{}, err error) {
	generatedObject := make(map[string]interface{})
	for _, template := range o.fieldTemplates {
		var generated interface{}
		if generated, err = template.Generate(context, request); err != nil {
			return result, err
		}
		generatedObject[template.FieldName] = generated
	}
	return generatedObject, err
}

func (o *ObjectType) GetName() string {
	return ""
}

type ObjectTypeFactory struct{}

func (o ObjectTypeFactory) DefaultOptions() TypeOptions {
	return TypeOptions{}
}

func (o ObjectTypeFactory) New(parameters TypeFactoryParameter) (generator TypeGenerator, err error) {
	fields := parameters.Options[constants.ObjectFieldsTemplatesOptionName].([]FieldModel)
	sort.Slice(fields, func(i int, j int) bool {
		return fields[i].FieldName < fields[j].FieldName
	})
	return &ObjectType{
		fieldTemplates: fields,
	}, err
}
