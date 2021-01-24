package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imdario/mergo"
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
	refs := model.Refs{}
	if result, err = GenerateModel(config, &refs); err != nil {
		return err
	}
	var context model.GeneratorContext
	if context, err = model.NewGenerationContext(config); err != nil {
		return err
	}
	context.Refs = refs
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
	var jsons []map[string]interface{}
	if err = getJsonFromFilesRecursivly(file, &jsons); err != nil {
		return config, err
	}
	var mergedJson []byte
	if mergedJson, err = mergeJsons(&jsons); err != nil {
		return config, err
	}
	if err = json.Unmarshal(mergedJson, &config); err != nil {
		return config, err
	}
	return config, err
}

func getJsonFromFilesRecursivly(file *os.File, jsons *[]map[string]interface{}) (err error) {
	var fileContent []byte
	if fileContent, err = ioutil.ReadAll(file); err != nil {
		return err
	}
	var data map[string]interface{}
	if err = json.Unmarshal(fileContent, &data); err != nil {
		return err
	}
	*jsons = append(*jsons, data)
	var includedFilesExists bool
	var rawIncludedFiles interface{}
	if rawIncludedFiles, includedFilesExists = data["includes"]; includedFilesExists {
		if includedFiles, isArrayOrString := rawIncludedFiles.([]interface{}); isArrayOrString {
			for _, fileNameToInclude := range includedFiles {
				if fileName, isString := fileNameToInclude.(string); isString {
					var fileToInclude *os.File
					if fileToInclude, err = os.Open(fileName); err != nil {
						return err
					}
					if err = getJsonFromFilesRecursivly(fileToInclude, jsons); err != nil {
						return err
					}
				} else {
					err = errors.New(fmt.Sprintf("unable to inclide %v", fileNameToInclude))
					return err
				}
			}
		} else {
			err = errors.New(fmt.Sprintf("wrong file inclusion in %s", file.Name()))
			return err
		}
	}
	return err
}

func mergeJsons(jsons *[]map[string]interface{}) (result []byte, err error) {
	finalJson := make(map[string]interface{})
	for _, jsonContent := range *jsons {
		if err = mergo.Merge(&finalJson, jsonContent, mergo.WithOverride); err != nil {
			return result, err
		}
	}
	result, err = json.Marshal(finalJson)
	return result, err
}

func setDefaultValueInConfiguration(config *configuration.Configuration) {
	if config.Options.Amount <= 0 {
		config.Options.Amount = 10
	}
}
