package model

import (
	"errors"
	"fmt"
	"github.com/ogama/gogen/src/configuration"
	"hash/fnv"
	"math/rand"
	"time"
)

type Refs struct {
	ref map[string]*TypeGenerator
}

func (r *Refs) PutRef(name string, generator *TypeGenerator) (err error) {
	if r.ref == nil {
		r.ref = make(map[string]*TypeGenerator)
	}
	if _, exists := r.ref[name]; exists {
		err = errors.New(fmt.Sprintf("ref name %s already exists", name))
	} else {
		r.ref[name] = generator
	}
	return err
}

func (r *Refs) GetRefValue(name string, context *GeneratorContext) (value interface{}, err error) {
	var exists bool
	var generator *TypeGenerator
	if generator, exists = r.ref[name]; !exists {
		err = errors.New(fmt.Sprintf("no values for ref name %s", name))
	} else {
		value, err = (*generator).Generate(context, Ref)
	}
	return value, err
}

type GeneratorContext struct {
	Config configuration.Configuration
	Rand   *rand.Rand
	Refs   Refs
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

func (gc GeneratorContext) Skip(context *GeneratorContext) {
	skip := context.Config.Options.Skip
	if skip > 0 {
		for i := 0; i < skip; i++ {
			_ = gc.Rand.Int()
		}
	}
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
