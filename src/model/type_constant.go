package model

type ConstantType struct {
	constant interface{}
}

func (c *ConstantType) Generate(_ *GeneratorContext, _ GenerationRequest) (result interface{}, err error) {
	return c.constant, err
}

func (c *ConstantType) GetName() string {
	return ""
}

type ConstantTypeFactory struct{}

func (c ConstantTypeFactory) DefaultOptions() TypeOptions {
	return TypeOptions{}
}

func (c ConstantTypeFactory) New(parameters TypeFactoryParameter) (generator TypeGenerator, err error) {
	return &ConstantType{
		constant: parameters.Template,
	}, err
}
