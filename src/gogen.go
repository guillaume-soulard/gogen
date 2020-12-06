package src

import (
	"encoding/json"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/model"
	"io/ioutil"
	"os"
)

func ExecuteFile(file *os.File) (err error) {
	var config configuration.Configuration
	if config, err = LoadConfigurationFromFile(file); err != nil {
		return err
	}
	setDefaultValueInConfiguration(&config)
	var result model.Model
	if result, err = GenerateModel(config); err != nil {
		return err
	}
	var context model.GeneratorContext
	if context, err = model.NewGenerationContext(config); err != nil {
		return err
	}
	return result.Generate(&context)
}

func Execute(args Args) (err error) {
	for i := 0; i < len(*args.Files); i++ {
		if err = ExecuteFile(&(*args.Files)[i]); err != nil {
			return err
		}
	}
	return err
}

func LoadConfigurationFromFile(file *os.File) (config configuration.Configuration, err error) {
	var data []byte
	if data, err = ioutil.ReadAll(file); err != nil {
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
