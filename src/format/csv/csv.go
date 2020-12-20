package csv

import (
	"errors"
	"fmt"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/format/common"
	"strings"
)

type BuilderCsv struct{}

const (
	quoteBehaviourNever   = "never"
	quoteBehaviourAlways  = "always"
	quoteBehaviourStrings = "strings"
)

func (b BuilderCsv) Build(configuration configuration.FormatConfiguration) (result common.Format, err error) {
	var separator string
	if separator, err = configuration.Options.GetStringOrDefault("separator", ","); err != nil {
		return result, err
	}
	var quoteChar string
	if quoteChar, err = configuration.Options.GetStringOrDefault("quoteChar", ""); err != nil {
		return result, err
	}
	var quoteBehaviour string
	if quoteBehaviour, err = configuration.Options.GetStringOrDefault("quoteBehaviour", "strings"); err != nil {
		return result, err
	}
	if strings.ToLower(quoteBehaviour) != quoteBehaviourAlways &&
		strings.ToLower(quoteBehaviour) != quoteBehaviourNever &&
		strings.ToLower(quoteBehaviour) != quoteBehaviourStrings {
		err = errors.New(fmt.Sprintf("quoteBehaviour value '%s' not supported", quoteBehaviour))
		return result, err
	}
	quoteBehaviour = strings.ToLower(quoteBehaviour)
	var headers bool
	if headers, err = configuration.Options.GetBoolOrDefault("headers", true); err != nil {
		return result, err
	}
	var columns string
	if columns, err = configuration.Options.GetStringOrDefault("columns", ""); err != nil {
		return result, err
	}
	columnArray := strings.Split(columns, ",")
	finalColumnArray := make([]string, 0)
	for _, column := range columnArray {
		if column != "" {
			finalColumnArray = append(finalColumnArray, column)
		}
	}
	result = &FormatCsv{
		separator:      separator,
		quoteChar:      quoteChar,
		headers:        headers,
		columns:        finalColumnArray,
		quoteBehaviour: quoteBehaviour,
	}
	return result, err
}

type FormatCsv struct {
	separator      string
	quoteChar      string
	headers        bool
	columns        []string
	config         []csvConfig
	quoteBehaviour string
}

type csvConfig struct {
	headerName string
	valuePath  []string
}

func (f *FormatCsv) Begin() (err error) {
	return err
}

func (f *FormatCsv) Format(generatedObject common.GeneratedObject) (result string, err error) {
	if f.config == nil {
		f.config = []csvConfig{}
		getCsvConfigFrom(generatedObject.Object, f, []string{})
		filterAndOrderColumns(f)
	}
	result = ""
	if f.headers {
		result = fmt.Sprintln(f.doFormatHeader())
		f.headers = false
	}
	var line string
	if line, err = f.doFormatCsv(generatedObject); err != nil {
		return result, err
	}
	result = fmt.Sprintf("%s%s", result, line)
	return result, err
}

func (f *FormatCsv) End() (err error) {
	return err
}

func (f *FormatCsv) doFormatCsv(generatedObject common.GeneratedObject) (result string, err error) {
	temp := make([]string, len(f.config))
	for i, fieldConfig := range f.config {
		if value, exists := generatedObject.GetValue(fieldConfig.valuePath); exists {
			temp[i] = f.quoteIfNeeded(value)
		} else {
			err = errors.New(fmt.Sprintf("field %s not found", strings.Join(fieldConfig.valuePath, ".")))
			return result, err
		}
	}
	result = strings.Join(temp, f.separator)
	return result, err
}

func (f *FormatCsv) doFormatHeader() (result string) {
	headers := make([]string, len(f.config))
	for i, c := range f.config {
		headers[i] = f.quoteIfNeeded(c.headerName)
	}
	result = strings.Join(headers, f.separator)
	return result
}

func (f *FormatCsv) quoteIfNeeded(value interface{}) string {
	_, isString := value.(string)
	shouldQuote := f.quoteBehaviour == quoteBehaviourAlways || (f.quoteBehaviour == quoteBehaviourStrings && isString)
	if shouldQuote {
		return fmt.Sprintf("%s%v%s", f.quoteChar, value, f.quoteChar)
	} else {
		return fmt.Sprintf("%v", value)
	}
}

func getCsvConfigFrom(object interface{}, f *FormatCsv, path []string) {
	if objectMap, isMap := object.(map[string]interface{}); isMap {
		for fieldName, fieldValue := range objectMap {
			getCsvConfigFrom(fieldValue, f, append(path, fieldName))
		}
	} else {
		f.config = append(f.config, csvConfig{
			headerName: strings.Join(path, "."),
			valuePath:  path,
		})
	}
}

func filterAndOrderColumns(f *FormatCsv) {
	if len(f.columns) > 0 {
		finalCsvConfig := make([]csvConfig, 0)
		for _, column := range f.columns {
			for _, configItem := range f.config {
				if configItem.headerName == column {
					finalCsvConfig = append(finalCsvConfig, configItem)
				}
			}
		}
		f.config = finalCsvConfig
	}
}
