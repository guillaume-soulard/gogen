package types

import (
	"errors"
	"fmt"
	"math"
)

type NumberType struct {
	minBound          float64
	maxBound          float64
	decimal           int
	sequenceEnable    bool
	sequenceCycle     bool
	sequenceIncrement float64
	currentSequence   float64
}

func (n *NumberType) Generate(context *GeneratorContext) (result interface{}, err error) {
	if n.sequenceEnable {
		if !n.sequenceCycle && n.currentSequence > n.maxBound {
			return result, errors.New(fmt.Sprintf("sequence reach the maximum value of %f set sequence.cycle to true to allow sequence to restart from bounds.min", n.maxBound))
		} else if n.sequenceCycle && n.currentSequence > n.maxBound {
			n.currentSequence = n.minBound
		}
		value := n.currentSequence
		n.currentSequence = n.currentSequence + n.sequenceIncrement
		return value, err
	} else {
		randomNumber := context.GenerateFloatBetween(n.minBound, n.maxBound)
		p := math.Pow10(n.decimal)
		value := float64(int(randomNumber*p)) / p
		return value, err
	}
}

type NumberTypeFactory struct{}

func (n NumberTypeFactory) DefaultOptions() TypeOptions {
	defaultOptions := TypeOptions{}
	defaultOptions.Add("bounds.min", 0)
	defaultOptions.Add("bounds.max", 10)
	defaultOptions.Add("decimal", 0)
	defaultOptions.Add("sequence.enable", false)
	defaultOptions.Add("sequence.cycle", true)
	defaultOptions.Add("sequence.increment", 1)
	return defaultOptions
}

func (n NumberTypeFactory) New(parameters TypeFactoryParameter)  (generator TypeGenerator, err error) {
	min := parameters.Options.GetOptionAsFloat("bounds.min")
	max := parameters.Options.GetOptionAsFloat("bounds.max")
	if min > max {
		return generator, errors.New(fmt.Sprintf("bounds.min = %f is greater than bounds.max = %f", min, max))
	}
	return &NumberType{
		minBound:          min,
		maxBound:          max,
		decimal:           parameters.Options.GetOptionAsInt("decimal"),
		sequenceEnable:    parameters.Options.GetOptionAsBool("sequence.enable"),
		sequenceCycle:     parameters.Options.GetOptionAsBool("sequence.cycle"),
		sequenceIncrement: parameters.Options.GetOptionAsFloat("sequence.increment"),
		currentSequence:   min,
	}, err
}
