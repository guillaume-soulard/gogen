package model

import (
	"errors"
	"fmt"
)

type RefType struct {
	refName string
}

func (n *RefType) Generate(context *GeneratorContext, _ GenerationRequest) (result interface{}, err error) {
	result, err = context.Refs.GetRefValue(n.refName, context)
	return result, err
}

func (n *RefType) GetName() string {
	return ""
}

type RefTypeFactory struct{}

func (n RefTypeFactory) DefaultOptions() TypeOptions {
	defaultOptions := TypeOptions{}
	defaultOptions.Add("name", "")
	return defaultOptions
}

func (n RefTypeFactory) New(parameters TypeFactoryParameter) (generator TypeGenerator, err error) {
	refName := parameters.Options.GetOptionAsString("name")
	if refName == "" {
		return generator, errors.New(fmt.Sprintf("name must be specify"))
	}
	return &RefType{
		refName: refName,
	}, err
}
