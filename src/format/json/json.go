package json

import (
	jsonEncode "encoding/json"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/format/common"
)

type BuilderJson struct{}

func (b BuilderJson) Build(configuration configuration.FormatConfiguration) (result common.Format, err error) {
	var pretty bool
	pretty, err = configuration.Options.GetBoolOrDefault("pretty", false)
	result = FormatJson{pretty: pretty}
	return result, err
}

type FormatJson struct {
	pretty bool
}

func (f FormatJson) Format(generatedObject common.GeneratedObject) (result string, err error) {
	var marshalResult []byte
	if f.pretty {
		marshalResult, err = jsonEncode.MarshalIndent(generatedObject.Object, "", "  ")
	} else {
		marshalResult, err = jsonEncode.Marshal(generatedObject.Object)
	}
	result = string(marshalResult)
	return result, err
}
