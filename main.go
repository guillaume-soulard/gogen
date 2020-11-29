package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	. "github.com/ogama/gogen/model"
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
	var config Configuration
	if config, err = LoadConfigurationFromFile(fileName); err != nil {
		return err
	}
	setDefaultValueInConfiguration(&config)
	var result Model
	if result, err = GenerateModel(config); err != nil {
		return err
	}
	var context GeneratorContext
	if context, err = NewGenerationContext(config); err != nil {
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

func LoadConfigurationFromFile(fileName string) (config Configuration, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(fileName); err != nil {
		return config, err
	}
	if err = json.Unmarshal(data, &config); err != nil {
		return config, err
	}
	return config, err
}

func setDefaultValueInConfiguration(config *Configuration) {
	if config.Options.Amount <= 0 {
		config.Options.Amount = 10
	}
}
