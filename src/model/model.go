package model

import (
	"fmt"
	"github.com/ogama/gogen/src/format/common"
	"github.com/ogama/gogen/src/output"
	"time"
)

type Model struct {
	ObjectModel ObjectModel
}

func (m Model) Generate(context *GeneratorContext) (err error) {
	startTime := time.Now()
	interval := context.Config.Options.Generation.Options.Interval
	var outputs []output.FormatThenOutput
	if outputs, err = getOutputsFrom(context); err != nil {
		return err
	}
	if err = doBeginOutput(&outputs); err != nil {
		return err
	}
	context.Skip(context)
	formatContext := common.FormatContext{
		Config: context.Config,
	}
	for i := 1; i <= context.Config.Options.Amount; i++ {
		var generatedObject interface{}
		if generatedObject, err = m.ObjectModel.Generate(context, Generate); err != nil {
			return err
		}
		if objectMap, isMap := generatedObject.(map[string]interface{}); isMap {
			generatedObject = objectMap["template"]
		}
		if err = doOutput(&outputs, &formatContext, generatedObject); err != nil {
			return err
		}
		if interval > 0 {
			time.Sleep(time.Millisecond * time.Duration(interval))
		}
	}
	if err = doEndOutput(&outputs); err != nil {
		return err
	}
	endTime := time.Now()
	fmt.Println(fmt.Sprintf("Generation end took : %s", endTime.Sub(startTime).String()))
	return err
}

func getOutputsFrom(context *GeneratorContext) (outputs []output.FormatThenOutput, err error) {
	outputConfig := context.Config.OutputConfigurations
	if outputConfig == nil {
		outputs = []output.FormatThenOutput{output.Outputs.GetDefaultOutput()}
	} else {
		outputs = make([]output.FormatThenOutput, len(outputConfig))
		for i, configuration := range outputConfig {
			if outputs[i], err = output.Outputs.GetOutput(configuration); err != nil {
				return outputs, err
			}
		}
	}
	return outputs, err
}

func doBeginOutput(outputs *[]output.FormatThenOutput) (err error) {
	for _, o := range *outputs {
		if err = o.Begin(); err != nil {
			return err
		}
	}
	return err
}

func doOutput(outputs *[]output.FormatThenOutput, context *common.FormatContext, object interface{}) (err error) {
	for _, o := range *outputs {
		if err = o.FormatAndWrite(object, context); err != nil {
			return err
		}
	}
	return err
}

func doEndOutput(outputs *[]output.FormatThenOutput) (err error) {
	for _, o := range *outputs {
		if err = o.End(); err != nil {
			return err
		}
	}
	return err
}
