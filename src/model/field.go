package model

type FieldModel struct {
	FieldName string
	Value     TypeGenerator
}

func (fieldTemplate FieldModel) Generate(context *GeneratorContext) (result interface{}, err error) {
	return fieldTemplate.Value.Generate(context)
}
