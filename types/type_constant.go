package types

type ConstantType struct {
	constant interface{}
}

func (c *ConstantType) Generate(context *GeneratorContext) (result interface{}, err error) {
	return c.constant, err
}

type ConstantTypeFactory struct {}

func (c ConstantTypeFactory) DefaultOptions() TypeOptions {
	return TypeOptions{}
}

func (c ConstantTypeFactory) New(parameters TypeFactoryParameter)  (generator TypeGenerator, err error) {
	return &ConstantType{
		constant: parameters.Template,
	}, err
}

