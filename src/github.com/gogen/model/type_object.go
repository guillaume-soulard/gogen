package model

import "sort"

const ObjectFieldsTemplatesOptionName = "objectFieldsTemplates"

type ObjectType struct {
	fieldTemplates []FieldModel
}

func (o *ObjectType) Generate(context *GeneratorContext) (result interface{}, err error) {
	generatedObject := make(map[string]interface{})
	for _, template := range o.fieldTemplates {
		var generated interface{}
		if generated, err = template.Generate(context); err != nil {
			return result, err
		}
		generatedObject[template.FieldName] = generated
	}
	return generatedObject, err
}

type ObjectTypeFactory struct{}

func (o ObjectTypeFactory) DefaultOptions() TypeOptions {
	return TypeOptions{}
}

func (o ObjectTypeFactory) New(parameters TypeFactoryParameter) (generator TypeGenerator, err error) {
	fields := parameters.Options[ObjectFieldsTemplatesOptionName].([]FieldModel)
	sort.Slice(fields, func(i int, j int) bool {
		return fields[i].FieldName < fields[j].FieldName
	})
	return &ObjectType{
		fieldTemplates: fields,
	}, err
}
