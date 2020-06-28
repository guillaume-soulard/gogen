package model

import (
	"errors"
	"fmt"
	"strings"
)

func GenerateModel(configuration Configuration) (result Model, err error) {
	var model ObjectModel
	var aliases map[string]interface{}
	if aliases, err = generateAliases(configuration); err != nil {
		return result, err
	}
	configuration.TemplateConfiguration = replaceAliasesInTemplate(aliases, configuration.TemplateConfiguration)
	if model, err = generateModel(configuration, "template", configuration.TemplateConfiguration); err != nil {
		return result, err
	}
	result.ObjectModel = model
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
		for fieldName, _ := range objectMap {
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

func generateAliases(configuration Configuration) (aliases map[string]interface{}, err error) {
	aliases = make(map[string]interface{})
	for aliasName, alias := range configuration.Options.Alias {
		aliases[aliasName] = alias.TemplateConfiguration
	}
	return aliases, err
}

func generateModel(configuration Configuration, fieldName string, template interface{}) (result ObjectModel, err error) {
	fieldType := getFieldType(template)
	if typeFactory, exists := TypeFactories.GetFactory(fieldType); exists {
		var options TypeOptions
		if options, err = getOptions(fieldType, template, configuration); err != nil {
			return result, err
		}
		options = typeFactory.DefaultOptions().OverloadWith(options)
		var value TypeGenerator
		if value, err = typeFactory.New(TypeFactoryParameter{
			Configuration: configuration,
			Template:      template,
			Options:       options,
		}); err != nil {
			return result, err
		}
		result.Fields = append(result.Fields, FieldModel{
			FieldName: fieldName,
			Value:     value,
		})
	} else {
		return result, errors.New(fmt.Sprintf("Unsupported type %s", fieldType))
	}
	return result, err
}

func getOptions(fieldType string, template interface{}, configuration Configuration) (options TypeOptions, err error) {
	options = make(map[string]interface{})
	if fieldType == "object" {
		if err = objectSpecificGeneration(options, template, configuration); err != nil {
			return options, err
		}
	} else if fieldType == "array" {
		options = getOptionField(template)
		if options[ArrayGeneratorOptionName], err = generateModel(configuration, "itemTemplate", options["itemTemplate"]); err != nil {
			return options, err
		}
	} else {
		options = getOptionField(template)
	}
	return options, err
}

func objectSpecificGeneration(options TypeOptions, template interface{}, configuration Configuration) (err error) {
	var objectFieldsTemplates []FieldModel
	for fieldName, fieldValue := range template.(map[string]interface{}) {
		var objectTemplate ObjectModel
		if objectTemplate, err = generateModel(configuration, fieldName, fieldValue); err != nil {
			return err
		}
		objectFieldsTemplates = append(objectFieldsTemplates, objectTemplate.Fields...)
	}
	options[ObjectFieldsTemplatesOptionName] = objectFieldsTemplates
	return err
}

func getOptionField(template interface{}) (result TypeOptions) {
	if fields, isMap := template.(map[string]interface{}); isMap {
		for fieldName, fieldValue := range fields {
			if strings.ToLower(fieldName) == "options" {
				if options, isOptionMap := fieldValue.(map[string]interface{}); isOptionMap {
					return BuildTypeOption(options)
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
