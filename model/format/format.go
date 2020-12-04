package format

import (
	"errors"
	"fmt"
	"github.com/ogama/gogen/model/configuration"
	"github.com/ogama/gogen/model/format/common"
	"github.com/ogama/gogen/model/format/json"
)

type StrategyFormat struct {
	defaultFormat common.Builder
	formats       map[string]common.Builder
}

func (s StrategyFormat) GetFormatOfDefault(configuration configuration.FormatConfiguration) (result common.Format, err error) {
	if configuration.Type == "" {
		result = s.GetDefaultFormat()
		return result, err
	}
	if builder, exists := s.formats[configuration.Type]; exists {
		result, err = builder.Build(configuration)
	} else {
		err = errors.New(fmt.Sprintf("unknown builder '%s'", configuration.Type))
	}
	return result, err
}

func (s StrategyFormat) GetDefaultFormat() common.Format {
	if result, err := s.defaultFormat.Build(configuration.FormatConfiguration{Options: configuration.EmptyOptions()}); err != nil {
		panic(err)
	} else {
		return result
	}
}

var Formats = StrategyFormat{
	formats: map[string]common.Builder{
		"json": json.BuilderJson{},
	},
	defaultFormat: json.BuilderJson{},
}
