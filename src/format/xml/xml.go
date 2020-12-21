package xml

import (
	"github.com/clbanning/anyxml"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/format/common"
)

type BuilderXml struct{}

func (b BuilderXml) Build(configuration configuration.FormatConfiguration) (result common.Format, err error) {
	var pretty bool
	if pretty, err = configuration.Options.GetBoolOrDefault("pretty", false); err != nil {
		return result, err
	}
	var objectRootName string
	objectRootName, err = configuration.Options.GetStringOrDefault("objectRootName", "object")
	result = FormatXml{
		pretty:         pretty,
		objectRootName: objectRootName,
	}
	return result, err
}

type FormatXml struct {
	pretty         bool
	objectRootName string
}

func (f FormatXml) Begin() (err error) {
	return err
}

func (f FormatXml) Format(generatedObject common.GeneratedObject, context *common.FormatContext) (result string, err error) {
	var marshalResult []byte
	if f.pretty {
		marshalResult, err = anyxml.XmlIndent(generatedObject.Object, "", "  ", f.objectRootName)
	} else {
		marshalResult, err = anyxml.Xml(generatedObject.Object, f.objectRootName)
	}
	result = string(marshalResult)
	return result, err
}

func (f FormatXml) End() (err error) {
	return err
}
