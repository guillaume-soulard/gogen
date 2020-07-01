package gogen

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	. "github.com/ogama/gogen/model"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
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
	return ExecuteFile(args[1])
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
