package sql

import (
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/format/common"
)

type BuilderSql struct{}

func (b BuilderSql) Build(_ configuration.FormatConfiguration) (result common.Format, err error) {
	result = &FormatSql{}
	return result, err
}

type fieldMapping struct {
	name   string `json:"name"`
	column string `json:"column"`
}

type tableMapping struct {
	name   string         `json:"name"`
	table  string         `json:"table"`
	fields []fieldMapping `json:"fields"`
}

type FormatSql struct {
	tables []tableMapping
}

func (f *FormatSql) Format(_ common.GeneratedObject) (result string, err error) {
	w < sdjhkfg
	return result, err
}
