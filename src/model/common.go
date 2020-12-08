package model

import (
	"github.com/ogama/gogen/src/configuration"
	"hash/fnv"
	"math/rand"
	"time"
)

type GeneratedValuesByType map[string]interface{}

type GeneratorContext struct {
	Config                configuration.Configuration
	Rand                  *rand.Rand
	GeneratedValuesByType GeneratedValuesByType
}

func (gc GeneratorContext) GenerateIntegerBetween(min int, max int) int {
	return min + gc.Rand.Intn(max-min+1)
}

func (gc GeneratorContext) GenerateInteger64Between(min int64, max int64) int64 {
	return min + gc.Rand.Int63n(max-min+1)
}

func (gc GeneratorContext) GenerateFloatBetween(min float64, max float64) float64 {
	return min + gc.Rand.Float64()*(max-min)
}

func NewGenerationContext(config configuration.Configuration) (result GeneratorContext, err error) {
	h := fnv.New32a()
	var seed int64
	if config.Options.Seed != "" {
		if _, err = h.Write([]byte(config.Options.Seed)); err != nil {
			return result, err
		}
		seed = int64(h.Sum32())
	} else {
		seed = time.Now().UnixNano()
	}
	result.Rand = rand.New(rand.NewSource(seed))
	result.Config = config
	return result, err
}

type TypeFactory interface {
	New(parameters TypeFactoryParameter) (generator TypeGenerator, err error)
	DefaultOptions() TypeOptions
}
