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
	statements := make([]string, 0)
	if mapValue, isMap := generatedObject.Object.(map[string]interface{}); isMap {
		if err = f.generateSqlForObject("", mapValue, &statements, context); err != nil {
			return result, err
		}
		result = fmt.Sprintf("%s\n", strings.Join(statements, "\n"))
	} else {
		err = errors.New("root is not an object")
	}
	return result, err
}

func (f *FormatSql) End() (err error) {
	return err
}

func (f *FormatSql) generateSqlForObject(fieldPath string, object map[string]interface{}, statements *[]string, context *common.FormatContext) (err error) {
	if tableMapping, exists := f.tables[fieldPath]; exists {
		if err = f.generateSqlForObjectWithMapping(tableMapping, object, statements, context); err != nil {
			return err
		}
	} else {
		if err = f.generateSqlForObjectWithoutMapping(fieldPath, object, statements); err != nil {
			return err
		}
	}
	return err
}

func (f *FormatSql) generateSqlForObjectWithMapping(mapping tableMapping, object map[string]interface{}, statements *[]string, context *common.FormatContext) (err error) {
	columnNames := make([]string, len(mapping.Fields))
	values := make([]string, len(mapping.Fields))
	var exists bool
	var value interface{}
	for i, field := range mapping.Fields {
		columnNames[i] = field.Column
		if value, exists = common.GetValue(field, []string{field.Name}); !exists {
			err = errors.New(fmt.Sprintf("field %s not found in generated object", field.Name))
			return err
		}
		values[i] = f.getSqlValueFor(value, context)
	}
	buildSqlStatement(columnNames, values, statements)
	return err
}

func (f *FormatSql) getSqlValueFor(value interface{}, context *common.FormatContext) string {
	if stringValue, isString := value.(string); isString {
		return fmt.Sprintf("'%s'", stringValue)
	} else {
		return fmt.Sprintf("%v", value)
	}
}

func (f *FormatSql) generateSqlForObjectWithoutMapping(fieldPath string, object map[string]interface{}, statements *[]string) (err error) {
	columnNames := make([]string, 0)
	values := make([]string, 0)
	for fieldName, fieldValue := range object {
		if mapValue, isMap := fieldValue.(map[string]interface{}); isMap {
			if err = f.generateSqlForObject(fmt.Sprintf("%s.%s", fieldPath, fieldName), mapValue, statements, nil); err != nil {
				return err
			}
		} else {
			columnNames = append(columnNames, fieldName)
			values = append(values, f.getSqlValueFor(fieldValue, nil))
		}
	}
	buildSqlStatement(columnNames, values, statements)
	return err
}

func buildSqlStatement(columns []string, values []string, statements *[]string) {
	if len(columns) > 0 && len(values) > 0 {
		statement := fmt.Sprintf("INSERT INTO (%s) VALUES (%s);", strings.Join(columns, ","), strings.Join(values, ","))
		*statements = append([]string{statement}, *statements...)
	}
}
