package model

import (
	"errors"
)

type PickerType struct {
	items              []interface{}
	probabilisticItems []ProbabilisticItem
	probabilistic      bool
}

type ProbabilisticItem struct {
	value interface{}
	min   float64
	max   float64
}

func (p *PickerType) Generate(context *GeneratorContext) (result interface{}, err error) {
	if p.probabilistic {
		var object interface{}
		random := context.GenerateFloatBetween(1, 100)
		for _, item := range p.probabilisticItems {
			if item.min <= random && item.max > random {
				object = item.value
				break
			}
		}
		return object, err
	} else {
		index := context.GenerateIntegerBetween(0, len(p.items)-1)
		return p.items[index], err
	}
}

type PickerTypeFactory struct{}

func (p PickerTypeFactory) DefaultOptions() TypeOptions {
	defaultOptions := TypeOptions{}
	defaultOptions.Add("items", nil)
	return defaultOptions
}

func (p PickerTypeFactory) New(parameters TypeFactoryParameter) (generator TypeGenerator, err error) {
	items := parameters.Options.GetOptionAsInterface("items").([]interface{})
	if items == nil {
		return generator, errors.New("items is required")
	}
	var isProbabilistic bool
	if isProbabilistic, err = isCorrectProbabilisticItems(items); err != nil {
		return generator, err
	}
	probabilisticItems := make([]ProbabilisticItem, len(items))
	if isProbabilistic {
		offset := 0.0
		for i := 0; i < len(items); i++ {
			probabilisticItems[i] = ProbabilisticItem{
				value: items[i].(map[string]interface{})["value"],
				min:   offset,
				max:   items[i].(map[string]interface{})["probability"].(float64) + offset,
			}
			offset += items[i].(map[string]interface{})["probability"].(float64)
		}
	}
	return &PickerType{
		items:              items,
		probabilisticItems: probabilisticItems,
		probabilistic:      isProbabilistic,
	}, err
}

func isCorrectProbabilisticItems(items []interface{}) (bool, error) {
	probabilisticItemsCount := 0
	nonProbabilisticItemsCount := 0
	probabilitySum := 0.0

	for _, item := range items {
		if mapItem, isMap := item.(map[string]interface{}); isMap {
			if mapItem["value"] != nil && mapItem["probability"] != nil {
				probabilisticItemsCount++
				probabilitySum += mapItem["probability"].(float64)
			} else {
				nonProbabilisticItemsCount++
			}
		}
	}
	if probabilisticItemsCount > 0 && nonProbabilisticItemsCount > 0 {
		return false, errors.New("probabilistic items structure is incorrect")
	}
	if probabilisticItemsCount > 0 && nonProbabilisticItemsCount == 0 && probabilitySum != 100 {
		return false, errors.New("probability sum is not equals to 100")
	}
	return probabilisticItemsCount > 0, nil
}
