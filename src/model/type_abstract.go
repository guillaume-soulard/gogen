package model

type AbstractTypeGenerator struct {
	lastGenerated   interface{}
	generationCount int64
	lastRequest     GenerationRequest
	TypeGenerator   TypeGenerator
	IsRef           bool
}

func (a *AbstractTypeGenerator) Generate(context *GeneratorContext, request GenerationRequest) (result interface{}, err error) {
	if a.IsRef {
		result, err = a.TypeGenerator.Generate(context, request)
	} else if request == Generate {
		if a.generationCount == 1 && a.lastRequest == Ref {
			result = a.lastGenerated
		} else {
			result, err = a.TypeGenerator.Generate(context, request)
			a.lastGenerated = result
		}
	} else if request == Ref {
		if a.generationCount == 0 {
			result, err = a.TypeGenerator.Generate(context, request)
			a.lastGenerated = result
		}
		result = a.lastGenerated
	}
	a.lastRequest = request
	a.generationCount++
	return result, err
}

func (a *AbstractTypeGenerator) GetName() string {
	return a.TypeGenerator.GetName()
}
