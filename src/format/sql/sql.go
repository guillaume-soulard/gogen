package sql

import (
	"encoding/json"
	"errors"
	"fmt"
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
		f.generateSqlForObjectWithMapping(tableMapping, object, statements)
	} else {
		f.generateSqlForObjectWithoutMapping(object, statements)
	}
}

func (f *FormatSql) generateSqlForObjectWithMapping(mapping tableMapping, object map[string]interface{}, statements *[]string) (err error) {
	columnNames := make([]string, len(mapping.fields))
	values := make([]string, len(mapping.fields))
	var exists bool
	var value interface
	for i, field := range mapping.fields {
		columnNames[i] = field.column
		if value, exists = common.GetValue(field, []string { field.name }); !exists {
			err = errors.New(fmt.Sprintf("field %s not found in generated object", field.name))
			return err
		}
		values[i] = f.getSqlValueFor(value)
	}
	return err
}

func (f *FormatSql) getSqlValueFor(value interface{}) string {
	if stringValue, isString := value.(string); isString {
		return fmt.Sprintf("'%s'", stringValue)
	} else {
		return fmt.Sprintf("%v", value)
	}
}

func (f *FormatSql) generateSqlForObjectWithoutMapping(object map[string]interface{}, statements *[]string) {
	for fieldName, fieldValue := range object {
		if mapValue, isMap := fieldValue.(map[string]interface{}); isMap {
			f.generateSqlForObject()
		} else {

		}
	}
}
