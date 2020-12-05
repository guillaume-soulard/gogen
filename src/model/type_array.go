package model

import (
	"errors"
	"fmt"
	"github.com/ogama/gogen/src/constants"
)

type ArrayType struct {
	minBound      int
	maxBound      int
	itemTemplate  interface{}
	itemGenerator ObjectModel
}

func (a *ArrayType) Generate(context *GeneratorContext) (result interface{}, err error) {
	numberOfItemsToGenerate := context.GenerateIntegerBetween(a.minBound, a.maxBound)
	array := make([]interface{}, numberOfItemsToGenerate)
	for i := 0; i < numberOfItemsToGenerate; i++ {
		var generatedItem interface{}
		if generatedItem, err = a.itemGenerator.Generate(context); err != nil {
			return result, err
		}
		if objectMap, isMap := generatedItem.(map[string]interface{}); isMap {
			generatedItem = objectMap["itemTemplate"]
		}
		array[i] = generatedItem
	}
	return array, err
}

type ArrayTypeFactory struct{}

func (a ArrayTypeFactory) DefaultOptions() TypeOptions {
	defaultOptions := TypeOptions{}
	defaultOptions.Add("bounds.min", 0)
	defaultOptions.Add("bounds.max", 10)
	defaultOptions.Add("itemTemplate", nil)
	return defaultOptions
}

func (a ArrayTypeFactory) New(parameters TypeFactoryParameter) (generator TypeGenerator, err error) {
	itemTemplate := parameters.Options.GetOptionAsInterface("itemTemplate")
	if itemTemplate == nil {
		return generator, errors.New("itemTemplate is required")
	}
	min := parameters.Options.GetOptionAsInt("bounds.min")
	max := parameters.Options.GetOptionAsInt("bounds.max")
	if min > max {
		return generator, errors.New(fmt.Sprintf("bounds.min = %d is greater than bounds.max = %d", min, max))
	}
	return &ArrayType{
		minBound:      min,
		maxBound:      max,
		itemGenerator: parameters.Options.GetOptionAsInterface(constants.ArrayGeneratorOptionName).(ObjectModel),
	}, err
}
