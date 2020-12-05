package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/ogama/gogen/src"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/model"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Argument missing file")
		return
	}
	if err := Execute(os.Args); err != nil {
		log.Panic(err)
	}
}

func ExecuteFile(fileName string) (err error) {
	var config configuration.Configuration
	if config, err = LoadConfigurationFromFile(fileName); err != nil {
		return err
	}
	setDefaultValueInConfiguration(&config)
	var result model.Model
	if result, err = src.GenerateModel(config); err != nil {
		return err
	}
	var context model.GeneratorContext
	if context, err = model.NewGenerationContext(config); err != nil {
		return err
	}
	return result.Generate(&context)
}

func Execute(args []string) (err error) {
	for i := 1; i < len(os.Args); i++ {
		if err = ExecuteFile(args[i]); err != nil {
			return err
		}
	}
	return err
}

func LoadConfigurationFromFile(fileName string) (config configuration.Configuration, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(fileName); err != nil {
		return config, err
	}
	if err = json.Unmarshal(data, &config); err != nil {
		return config, err
	}
	return config, err
}

func setDefaultValueInConfiguration(config *configuration.Configuration) {
	if config.Options.Amount <= 0 {
		config.Options.Amount = 10
	}
}
