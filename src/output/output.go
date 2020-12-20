package output

import (
	"errors"
	"fmt"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/format"
	"github.com/ogama/gogen/src/format/common"
	commonOutput "github.com/ogama/gogen/src/output/common"
	"github.com/ogama/gogen/src/output/file"
	"github.com/ogama/gogen/src/output/stdout"
)

type FormatThenOutput struct {
	Format common.Format
	Output commonOutput.Output
}

func (o FormatThenOutput) FormatAndWrite(object interface{}) (err error) {
	var formattedObject string
	if formattedObject, err = o.Format.Format(common.GeneratedObject{Object: object}); err != nil {
		return err
	}
	if err = o.Output.Write(formattedObject); err != nil {
		return err
	}
	return err
}

func (o FormatThenOutput) Begin() (err error) {
	if err = o.Format.Begin(); err != nil {
		return err
	}
	if err = o.Output.Begin(); err != nil {
		return err
	}
	return err
}

func (o FormatThenOutput) End() (err error) {
	if err = o.Format.End(); err != nil {
		return err
	}
	if err = o.Output.End(); err != nil {
		return err
	}
	return err
}

type StrategyOutput struct {
	defaultOutput commonOutput.Builder
	outputs       map[string]commonOutput.Builder
}

func (s StrategyOutput) GetOutput(configuration configuration.OutputConfiguration) (result FormatThenOutput, err error) {
	if outputBuilder, exists := s.outputs[configuration.Type]; exists {
		var outputFormat common.Format
		if outputFormat, err = format.Formats.GetFormatOfDefault(configuration.FormatConfiguration); err != nil {
			return result, err
		} else {
			var output commonOutput.Output
			if output, err = outputBuilder.Build(configuration); err != nil {
				return result, err
			}
			result = FormatThenOutput{
				Format: outputFormat,
				Output: output,
			}
		}
	} else {
		err = errors.New(fmt.Sprintf("unknown output type '%s'", configuration.Type))
	}
	return result, err
}

func (s StrategyOutput) GetDefaultOutput() FormatThenOutput {
	if output, err := s.defaultOutput.Build(configuration.OutputConfiguration{
		Type:                "",
		Options:             configuration.EmptyOptions(),
		FormatConfiguration: configuration.FormatConfiguration{},
	}); err != nil {
		panic(err)
	} else {
		return FormatThenOutput{
			Format: format.Formats.GetDefaultFormat(),
			Output: output,
		}
	}
}

var Outputs = StrategyOutput{
	outputs: map[string]commonOutput.Builder{
		"stdout": stdout.BuilderStdout{},
		"file":   file.BuilderFile{},
	},
	defaultOutput: stdout.BuilderStdout{},
}
