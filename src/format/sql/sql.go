package sql

import (
	"encoding/json"
	"errors"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/format/common"
	"strings"
)

type BuilderSql struct{}

func (b BuilderSql) Build(configuration configuration.FormatConfiguration) (result common.Format, err error) {
	var tables []tableMapping
	var tableConfig string
	if tableConfig, err = configuration.Options.GetStringOrDefault("tables", "[]"); err != nil {
		return result, err
	}
	if err = json.Unmarshal([]byte(tableConfig), &tables); err != nil {
		return result, err
	}
	result = &FormatSql{
		tables: buildTableMap(tables),
	}
	return result, err
}

func buildTableMap(tables []tableMapping) (result map[string]tableMapping) {
	result = make(map[string]tableMapping)
	for _, table := range tables {
		table.fieldMap = buildFieldMap(table.fields)
		result[table.name] = table
	}
	return result
}

func buildFieldMap(fields []fieldMapping) (result map[string]fieldMapping) {
	result = make(map[string]fieldMapping)
	for _, field := range fields {
		result[field.name] = field
	}
	return result
}

type fieldMapping struct {
	name   string `json:"name"`
	column string `json:"column"`
}

type tableMapping struct {
	name     string         `json:"name"`
	table    string         `json:"table"`
	fields   []fieldMapping `json:"fields"`
	fieldMap map[string]fieldMapping
}

type FormatSql struct {
	tables map[string]tableMapping
}

func (f *FormatSql) Format(generatedObject common.GeneratedObject) (result string, err error) {
	statements := make([]string, 0)
	if mapValue, isMap := generatedObject.Object.(map[string]interface{}); isMap {
		f.generateSqlForObject("", mapValue, &statements)
		result = strings.Join(statements, "\n")
	} else {
		err = errors.New("root is not an object")
	}
	return result, err
}

func (f *FormatSql) generateSqlForObject(fieldPath string, object map[string]interface{}, statements *[]string) {
	if tableMapping, exists := f.tables[fieldPath]; exists {
		generateSqlForObjectAndMapping(tableMapping, object, statements)
	} else {
		generateSqlForObjectWithoutMapping(object, statements)
	}
}
