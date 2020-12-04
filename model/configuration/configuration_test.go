package configuration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Get_should_return_the_value(t *testing.T) {
	// GIVEN
	option := Options{"field": "value"}
	// WHEN
	result := option.Get("field")
	// THEN
	assert.Equal(t, "value", result)
}

func Test_Get_should_return_the_value_at_depth_2(t *testing.T) {
	// GIVEN
	option := Options{"field": map[string]interface{}{"field2": "value"}}
	// WHEN
	result := option.Get("field.field2")
	// THEN
	assert.Equal(t, "value", result)
}

func Test_Get_should_return_empty_when_path_not_exists(t *testing.T) {
	// GIVEN
	option := Options{"field": "value"}
	// WHEN
	result := option.Get("unexisting")
	// THEN
	assert.Equal(t, "", result)
}

func Test_Get_should_return_empty_when_path_at_depth_2_not_exists(t *testing.T) {
	// GIVEN
	option := Options{"field": map[string]interface{}{"field2": "value"}}
	// WHEN
	result := option.Get("field.unexisting")
	// THEN
	assert.Equal(t, "", result)
}
