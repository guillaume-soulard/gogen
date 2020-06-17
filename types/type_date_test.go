package types

import (
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
	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})
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
			"truncate": "milliseconds",
		},
	})
	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})
	// THEN
	assert.Equal(t, time.Date(2099, time.December, 31, 23, 59, 59, 0, time.Local), result)
}