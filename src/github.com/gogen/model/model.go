package model

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Model struct {
	ObjectModel ObjectModel
}

func (m Model) Generate(context *GeneratorContext) (err error) {
	startTime := time.Now()
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
	}
	endTime := time.Now()
	log.Println(fmt.Sprintf("Generation end took : %s", endTime.Sub(startTime).String()))
	return err
}
