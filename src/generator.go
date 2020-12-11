package src

import (
	"errors"
	"fmt"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/constants"
	"github.com/ogama/gogen/src/model"
	"strings"
)

func GenerateModel(configuration configuration.Configuration, refs *model.Refs) (result model.Model, err error) {
	var m model.ObjectModel
	var aliases map[string]interface{}
	if aliases, err = generateAliases(configuration); err != nil {
		return result, err
	}
	configuration.TemplateConfiguration = replaceAliasesInTemplate(aliases, configuration.TemplateConfiguration)
	if m, err = generateModel(configuration, "template", configuration.TemplateConfiguration, refs); err != nil {
		return result, err
	}
	result.ObjectModel = m
	return result, err
}

func replaceAliasesInTemplate(aliases map[string]interface{}, template interface{}) interface{} {
	if objectMap, isMap := template.(map[string]interface{}); isMap {
		newTemplate := make(map[string]interface{})
		for fieldName, fieldValue := range objectMap {
			typeName := getTypeName(fieldValue)
			if isTypeField(fieldValue) && isAlias(typeName, aliases) {
				fieldValue = aliases[typeName]
			}
			fieldValue = replaceAliasesInTemplate(aliases, fieldValue)
			newTemplate[fieldName] = fieldValue
		}
		return newTemplate
	}
	return template
}

func isTypeField(template interface{}) bool {
	if objectMap, isMap := template.(map[string]interface{}); isMap {
		for fieldName := range objectMap {
			if fieldName == "_type" {
				return true
			}
		}
	}
	return false
}

func getTypeName(template interface{}) string {
	if objectMap, isMap := template.(map[string]interface{}); isMap {
		for fieldName, fieldValue := range objectMap {
			if _, isString := fieldValue.(string); fieldName == "_type" && isString {
				return fieldValue.(string)
			}
		}
	}
	return ""
}

func isAlias(name string, aliases map[string]interface{}) bool {
	_, contains := aliases[name]
	return contains
}

func generateAliases(configuration configuration.Configuration) (aliases map[string]interface{}, err error) {
	aliases = make(map[string]interface{})
	for aliasName, alias := range configuration.Options.Alias {
		aliases[aliasName] = alias.TemplateConfiguration
	}
	return aliases, err
}

func generateModel(configuration configuration.Configuration, fieldName string, template interface{}, refs *model.Refs) (result model.ObjectModel, err error) {
	fieldType := getFieldType(template)
	if typeFactory, exists := model.TypeFactories.GetFactory(fieldType); exists {
		var options model.TypeOptions
		if options, err = getOptions(fieldType, template, configuration, refs); err != nil {
			return result, err
		}
		options = typeFactory.DefaultOptions().OverloadWith(options)
		var value model.TypeGenerator
		isRef := strings.ToLower(fieldType) == "ref"
		if value, err = typeFactory.New(model.TypeFactoryParameter{
			Configuration: configuration,
			Template:      template,
			Options:       options,
		}); err != nil {
			return result, err
		}
		value = model.AbstractTypeGenerator{TypeGenerator: value, IsRef: isRef}
		fieldModel := model.FieldModel{
			FieldName: fieldName,
			Value:     value,
		}
		if fieldModel.Value.GetName() != "" {
			if err = refs.PutRef(fieldModel.Value.GetName(), &fieldModel.Value); err != nil {
				return result, err
			}
		}
		result.Fields = append(result.Fields, fieldModel)
	} else {
		return result, errors.New(fmt.Sprintf("Unsupported type %s", fieldType))
	}
	return result, err
}

func getOptions(fieldType string, template interface{}, configuration configuration.Configuration, refs *model.Refs) (options model.TypeOptions, err error) {
	options = make(map[string]interface{})
	if fieldType == "object" {
		if err = objectSpecificGeneration(options, template, configuration, refs); err != nil {
			return options, err
		}
	} else if fieldType == "array" {
		options = getOptionField(template)
		if options[constants.ArrayGeneratorOptionName], err = generateModel(configuration, "itemTemplate", options["itemTemplate"], refs); err != nil {
			return options, err
		}
	} else {
		options = getOptionField(template)
	}
	return options, err
}

func objectSpecificGeneration(options model.TypeOptions, template interface{}, configuration configuration.Configuration, refs *model.Refs) (err error) {
	var objectFieldsTemplates []model.FieldModel
	for fieldName, fieldValue := range template.(map[string]interface{}) {
		var objectTemplate model.ObjectModel
		if objectTemplate, err = generateModel(configuration, fieldName, fieldValue, refs); err != nil {
			return err
		}
		objectFieldsTemplates = append(objectFieldsTemplates, objectTemplate.Fields...)
	}
	options[constants.ObjectFieldsTemplatesOptionName] = objectFieldsTemplates
	return err
}

func getOptionField(template interface{}) (result model.TypeOptions) {
	if fields, isMap := template.(map[string]interface{}); isMap {
		for fieldName, fieldValue := range fields {
			if strings.ToLower(fieldName) == "options" {
				if options, isOptionMap := fieldValue.(map[string]interface{}); isOptionMap {
					return model.BuildTypeOption(options)
				}
			}
		}
	}
	return result
}

func getFieldType(template interface{}) string {
	if generatedType, isGeneratedType := isGeneratedType(template); isGeneratedType {
		return generatedType
	}
	if _, isInt := template.(int); isInt {
		return "constant"
	} else if _, isString := template.(string); isString {
		return "constant"
	} else if _, isFloat := template.(float64); isFloat {
		return "constant"
	} else if _, isBool := template.(bool); isBool {
		return "constant"
	} else {
		return "object"
	}
}

func isGeneratedType(template interface{}) (string, bool) {
	if fields, isMap := template.(map[string]interface{}); isMap {
		for fieldName, fieldValue := range fields {
			if strings.ToLower(fieldName) == "_type" {
				return strings.ToLower(fieldValue.(string)), true
			}
		}
	}
	return "", false
}
