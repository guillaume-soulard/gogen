package model

import (
	"github.com/ogama/gogen/src/reggen"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func Test_Generate_should_generate_single_character_string_generation(t *testing.T) {
	// GIVEN
	generator, _ := reggen.NewGenerator("a", rand.New(rand.NewSource(0)))
	stringType := StringType{
		stringGenerator: generator,
	}
	// WHEN
	result, _ := stringType.Generate(&GeneratorContext{})
	// THEN
	assert.Equal(t, "a", result)
}

func Test_Generate_should_generate_empty_string_when_pattern_is_empty(t *testing.T) {
	// GIVEN
	generator, _ := reggen.NewGenerator("", rand.New(rand.NewSource(0)))
	stringType := StringType{
		stringGenerator: generator,
	}
	// WHEN
	result, _ := stringType.Generate(&GeneratorContext{})
	// THEN
	assert.Equal(t, "", result)
}

func Test_Generate_should_generate_string_with_length_10(t *testing.T) {
	// GIVEN
	generator, _ := reggen.NewGenerator("[a-z]{10}", rand.New(rand.NewSource(0)))
	stringType := StringType{
		stringGenerator: generator,
	}
	// WHEN
	result, _ := stringType.Generate(&GeneratorContext{})
	// THEN
	assert.Len(t, result, 10)
}

func Test_Generate_should_generate_string_with_length_10_and_full_of_a(t *testing.T) {
	// GIVEN
	generator, _ := reggen.NewGenerator("[a]{10}", rand.New(rand.NewSource(0)))
	stringType := StringType{
		stringGenerator: generator,
	}
	// WHEN
	result, _ := stringType.Generate(&GeneratorContext{})
	// THEN
	assert.Equal(t, "aaaaaaaaaa", result)
}
