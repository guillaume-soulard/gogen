package json

import (
	jsonEncode "encoding/json"
	"github.com/ogama/gogen/model/configuration"
	"github.com/ogama/gogen/model/format/common"
)

type BuilderJson struct{}

func (b BuilderJson) Build(configuration configuration.FormatConfiguration) (result common.Format, err error) {
	var pretty bool
	pretty, err = configuration.Options.GetBoolOrDefault("pretty")
	result = FormatJson{pretty: pretty}
	return result, err
}

type FormatJson struct {
	pretty bool
}

func (f FormatJson) Format(object interface{}) (result string, err error) {
	var marshalResult []byte
	if f.pretty {
		marshalResult, err = jsonEncode.MarshalIndent(object, "", "  ")
	} else {
		marshalResult, err = jsonEncode.Marshal(object)
	}
	result = string(marshalResult)
	return result, err
}
