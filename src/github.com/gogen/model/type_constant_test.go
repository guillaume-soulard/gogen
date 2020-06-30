package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Generate_should_return_1(t *testing.T) {
	// GIVEN
	generator := ConstantType{constant: 1}
	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})
	// THEN
	assert.Equal(t, 1, result)
}

func Test_Generate_should_return_1_dot_1(t *testing.T) {
	// GIVEN
	generator := ConstantType{constant: 1.1}
	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})
	// THEN
	assert.Equal(t, 1.1, result)
}

func Test_Generate_should_return_true(t *testing.T) {
	// GIVEN
	generator := ConstantType{constant: true}
	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})
	// THEN
	assert.Equal(t, true, result)
}

func Test_Generate_should_return_test(t *testing.T) {
	// GIVEN
	generator := ConstantType{constant: "test"}
	// WHEN
	result, _ := generator.Generate(&GeneratorContext{})
	// THEN
	assert.Equal(t, "test", result)
}
