package model

import (
	"github.com/ogama/gogen/src/configuration"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_New_should_generate_type_without_error(t *testing.T) {
	// GIVEN
	factory := DateTypeFactory{}
	// WHEN
	_, err := factory.New(TypeFactoryParameter{
		Options: factory.DefaultOptions(),
	})
	// THEN
	assert.NoError(t, err)
}

func Test_Generate_should_generate_a_type_date(t *testing.T) {
	// GIVEN
	factory := DateTypeFactory{}
	generator, _ := factory.New(TypeFactoryParameter{
		Options: factory.DefaultOptions(),
	})
	context, err := NewGenerationContext(configuration.Configuration{})
	assert.NoError(t, err)
	// WHEN
	result, _ := generator.Generate(&context, Generate)
	// THEN
	assert.IsType(t, time.Time{}, result)
}

func Test_Generate_should_generate_a_date(t *testing.T) {
	// GIVEN
	factory := DateTypeFactory{}
	generator, _ := factory.New(TypeFactoryParameter{
		Options: map[string]interface{}{
			"bounds.min": "2020-06-11T14:32:24",
			"bounds.max": "2020-06-11T14:32:24",
			"truncate":   "milliseconds",
		},
	})
	context, err := NewGenerationContext(configuration.Configuration{})
	assert.NoError(t, err)
	// WHEN
	result, _ := generator.Generate(&context, Generate)
	// THEN
	assert.Equal(t, time.Date(2020, time.June, 11, 14, 32, 24, 0, time.UTC), result)
}
