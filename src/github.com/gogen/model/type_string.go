package model

import (
	"./reggen"
)

type StringType struct {
	pattern         string
	stringGenerator *reggen.Generator
}

func (s *StringType) Generate(context *GeneratorContext) (result interface{}, err error) {
	if s.stringGenerator == nil {
		var stringGenerator *reggen.Generator
		if stringGenerator, err = reggen.NewGenerator(s.pattern, context.Rand); err != nil {
			return result, err
		}
		s.stringGenerator = stringGenerator
	}
	return s.stringGenerator.Generate(2147483647), err
}

type StringTypeFactory struct{}

func (s StringTypeFactory) DefaultOptions() TypeOptions {
	defaultOptions := TypeOptions{}
	defaultOptions.Add("pattern", "[A-Z]{1}[A-Za-z]{10,25}")
	return defaultOptions
}

func (s StringTypeFactory) New(parameters TypeFactoryParameter) (generator TypeGenerator, err error) {
	pattern := parameters.Options.GetOptionAsString("pattern")
	return &StringType{
		pattern: pattern,
	}, err
}
