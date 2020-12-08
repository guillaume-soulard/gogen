package model

type FieldModel struct {
	FieldName string
	Value     TypeGenerator
}

func (fieldTemplate FieldModel) Generate(context *GeneratorContext) (result interface{}, err error) {
	result, err = fieldTemplate.Value.Generate(context)
	name := fieldTemplate.Value.GetName()
	if name != "" {
		context.GeneratedValuesByType[name] = result
	}
	return result, err
}
