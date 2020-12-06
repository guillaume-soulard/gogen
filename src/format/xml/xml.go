package xml

import (
	"encoding/xml"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/format/common"
)

type BuilderXml struct{}

func (b BuilderXml) Build(configuration configuration.FormatConfiguration) (result common.Format, err error) {
	var pretty bool
	pretty, err = configuration.Options.GetBoolOrDefault("pretty", false)
	result = FormatXml{pretty: pretty}
	return result, err
}

type FormatXml struct {
	pretty bool
}

func (f FormatXml) Format(object interface{}) (result string, err error) {
	var marshalResult []byte
	if f.pretty {
		marshalResult, err = xml.MarshalIndent(object, "", "  ")
	} else {
		marshalResult, err = xml.Marshal(object)
	}
	result = string(marshalResult)
	return result, err
}
