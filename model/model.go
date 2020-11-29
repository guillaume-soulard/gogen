package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Model struct {
	ObjectModel ObjectModel
}

func (m Model) Generate(context *GeneratorContext) (err error) {
	startTime := time.Now()
	interval := context.Config.Options.Generation.Options.Interval
	for i := 1; i <= context.Config.Options.Amount; i++ {
		var generatedObject interface{}
		if generatedObject, err = m.ObjectModel.Generate(context); err != nil {
			return err
		}
		if objectMap, isMap := generatedObject.(map[string]interface{}); isMap {
			generatedObject = objectMap["template"]
		}
		if jsonObject, err := json.MarshalIndent(generatedObject, "", "  "); err != nil {
			return err
		} else {
			fmt.Println(string(jsonObject))
		}
		if interval > 0 {
			time.Sleep(time.Millisecond * time.Duration(interval))
		}
	}
	endTime := time.Now()
	fmt.Println(fmt.Sprintf("Generation end took : %s", endTime.Sub(startTime).String()))
	return err
}
