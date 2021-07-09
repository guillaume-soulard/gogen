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
	if tableConfig, err = configuration.Options.GetObjectAsStringOrDefault("tables", "[]"); err != nil {
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
		table.fieldMap = buildFieldMap(table.Fields)
		result[table.Name] = table
	}
	return result
}

func buildFieldMap(fields []fieldMapping) (result map[string]fieldMapping) {
	result = make(map[string]fieldMapping)
	for _, field := range fields {
		result[field.Name] = field
	}
	return result
}

type fieldMapping struct {
	Name   string `json:"name"`
	Column string `json:"column"`
}

type tableMapping struct {
	Name     string         `json:"name"`
	Table    string         `json:"table"`
	Fields   []fieldMapping `json:"fields"`
	fieldMap map[string]fieldMapping
}

type FormatSql struct {
	tables map[string]tableMapping
}

func (f *FormatSql) Begin() (err error) {
	return err
}

func (f *FormatSql) Format(generatedObject common.GeneratedObject, context *common.FormatContext) (result string, err error) {
	var objectMap map[string]interface{}
	var isMap bool
	statements := make([]string, 0)
	if objectMap, isMap = getMap(generatedObject.Object); !isMap {
		err = errors.New("root is not an object")
		return result, err
	}
	err = f.formatObject(objectMap, &statements)
	result = strings.Join(statements, "\n")
	return result, err
}

func (f *FormatSql) formatObject(object map[string]interface{}, statements *[]string) (err error) {
	for fieldName, value := range object {
		if fieldMap, isFieldMap := getMap(value); isFieldMap {
			if err = f.formatMap(fieldName, fieldMap, statements); err != nil {
				return err
			}
		}
	}
	return err
}

func getMap(object interface{}) (objectMap map[string]interface{}, isMap bool) {
	if object != nil {
		objectMap, isMap = object.(map[string]interface{})
	} else {
		isMap = false
	}
	return objectMap, isMap
}

func (f *FormatSql) formatMap(fieldName string, fieldMap map[string]interface{}, statements *[]string) (err error) {
	var tableMapping tableMapping
	if tableMapping, err = f.getTableMappingOf(fieldName); err != nil {
		return err
	}
	fields := make([]string, 0)
	fieldValues := make([]string, 0)
	for fieldName, fieldValue := range fieldMap {
		if fieldMapping, fieldMappingExists := tableMapping.fieldMap[fieldName]; fieldMappingExists {
			fields = append(fields, fieldMapping.Column)
			fieldValues = append(fieldValues, getSqlValue(fieldValue))
		}
	}
	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		tableMapping.Table,
		strings.Join(fields, ","),
		strings.Join(fieldValues, ","),
	)
	*statements = append(*statements, statement)
	return err
}

func getSqlValue(fieldValue interface{}) string {
	if _, isString := fieldValue.(string); isString {
		return fmt.Sprintf("'%v'", fieldValue)
	} else {
		return fmt.Sprintf("%v", fieldValue)
	}
}

func (f *FormatSql) getTableMappingOf(name string) (mapping tableMapping, err error) {
	var exists bool
	if mapping, exists = f.tables[name]; !exists {
		err = errors.New(fmt.Sprintf("no sql mapping found for %s", name))
	}
	return mapping, err
}

func (f *FormatSql) End() (err error) {
	return err
}
