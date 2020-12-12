package csv

import (
	"fmt"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/format/common"
	"strings"
)

type BuilderCsv struct{}

func (b BuilderCsv) Build(configuration configuration.FormatConfiguration) (result common.Format, err error) {
	var separator string
	if separator, err = configuration.Options.GetStringOrDefault("separator", ","); err != nil {
		return result, err
	}
	var quoteChar string
	if quoteChar, err = configuration.Options.GetStringOrDefault("quoteChar", `"`); err != nil {
		return result, err
	}
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
	result = FormatCsv{
		separator: separator,
		quoteChar: quoteChar,
		headers:   headers,
		columns:   finalColumnArray,
	}
	return result, err
}

type FormatCsv struct {
	separator string
	quoteChar string
	headers   bool
	columns   []string
	config    []csvConfig
}

type csvConfig struct {
	headerName string
	valuePath  []string
}

func (f FormatCsv) Format(object interface{}) (result string, err error) {
	if f.config == nil {
		f.config = []csvConfig{}
		getCsvConfigFrom(object, &f, []string{})
		filterAndOrderColumns(&f)
	}
	result = ""
	if f.headers {
		result = fmt.Sprintln(f.doFormatHeader())
	}
	result = fmt.Sprintf("%s%s", result, f.doFormatCsv(object))
	return result, err
}

func (f FormatCsv) doFormatCsv(object interface{}) (result string) {
	temp := make([]string, len(f.config))
	for i, fieldConfig := range f.config {
		temp[i] = fmt.Sprintf("%v", getValue(object, fieldConfig))
	}
	return result
}

func getValue(object interface{}, config csvConfig) interface{} {

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
