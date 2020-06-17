package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Generate_should_generate_a_float_number(t *testing.T) {
	// GIVEN
	generator := NumberType{}

	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})

	// THEN
	assert.IsType(t, float64(0), result)
}

func Test_Generate_should_generate_0(t *testing.T) {
	// GIVEN
	generator := NumberType{
		minBound: 0,
		maxBound: 0,
	}

	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})

	// THEN
	assert.Equal(t, float64(0), result)
}

func Test_Generate_should_generate_1(t *testing.T) {
	// GIVEN
	generator := NumberType{
		minBound: 1,
		maxBound: 1,
	}

	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})

	// THEN
	assert.Equal(t, float64(1), result)
}

func Test_Generate_should_generate_1_dot_1(t *testing.T) {
	// GIVEN
	generator := NumberType{
		minBound: 1.1,
		maxBound: 1.1,
		decimal:  1,
	}

	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})

	// THEN
	assert.Equal(t, 1.1, result)
}

func Test_Generate_should_generate_1_dot_111(t *testing.T) {
	// GIVEN
	generator := NumberType{
		minBound: 1.11111,
		maxBound: 1.11111,
		decimal:  3,
	}

	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})

	// THEN
	assert.Equal(t, 1.111, result)
}

func Test_Generate_should_generate_minus_1(t *testing.T) {
	// GIVEN
	generator := NumberType{
		minBound: -1,
		maxBound: -1,
	}

	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})

	// THEN
	assert.Equal(t, float64(-1), result)
}

func Test_Generate_should_generate_number_between_1_and_10(t *testing.T) {
	// GIVEN
	generator := NumberType{
		minBound: 1,
		maxBound: 10,
	}

	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})

	// THEN
	assert.Contains(t, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, result)
}

func Test_Generate_should_generate_number_trough_1_to_10_in_order(t *testing.T) {
	// GIVEN
	generator := NumberType{
		minBound:          1,
		maxBound:          10,
		decimal:           0,
		sequenceEnable:    true,
		sequenceIncrement: 1,
		currentSequence:   1,
	}

	context := GeneratorContext{}
	for i := 1; i <= 10; i++ {
		// WHEN
		result, _ := generator.Generate(&context)
		// THEN
		assert.Equal(t, float64(i), result)
	}
}

func Test_Generate_should_generate_number_trough_1_to_5_in_order_twice(t *testing.T) {
	// GIVEN
	generator := NumberType{
		minBound:          1,
		maxBound:          5,
		decimal:           0,
		sequenceEnable:    true,
		sequenceIncrement: 1,
		currentSequence:   1,
		sequenceCycle:     true,
	}

	// WHEN
	context := GeneratorContext{}
	result := make([]float64, 0)
	for i := 1; i <= 10; i++ {
		generated, _ := generator.Generate(&context)
		result = append(result, generated.(float64))
	}

	// THEN
	assert.Equal(t, []float64{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}, result)
}

func Test_New_should_not_return_an_error(t *testing.T) {
	// GIVEN
	parameters := TypeFactoryParameter{
		Options: map[string]interface{}{
			"bounds.min": 1,
			"bounds.max": 1,
		},
	}

	// WHEN
	_, err := NumberTypeFactory{}.New(parameters)

	// THEN
	assert.NoError(t, err)
}

func Test_New_should_return_error_because_min_is_less_than_max(t *testing.T) {
	// GIVEN
	parameters := TypeFactoryParameter{
		Options: map[string]interface{}{
			"bounds.min": 2,
			"bounds.max": 1,
		},
	}

	// WHEN
	_, err := NumberTypeFactory{}.New(parameters)

	// THEN
	assert.Error(t, err)
}
