package model

type FieldModel struct {
	FieldName string
	Value     TypeGenerator
}

func (fieldTemplate FieldModel) Generate(context *GeneratorContext, request GenerationRequest) (result interface{}, err error) {
	result, err = fieldTemplate.Value.Generate(context, request)
	return result, err
}
