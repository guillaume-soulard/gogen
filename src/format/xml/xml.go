package xml

import (
	"github.com/clbanning/anyxml"
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/format/common"
	"time"
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

func (f FormatXml) Format(generatedObject common.GeneratedObject) (result string, err error) {
	var marshalResult []byte
	if err = formatTimeRecursively(&generatedObject.Object); err != nil {
		return result, err
	}
	if f.pretty {
		marshalResult, err = anyxml.XmlIndent(generatedObject.Object, "", "  ", f.objectRootName)
	} else {
		marshalResult, err = anyxml.Xml(generatedObject.Object, f.objectRootName)
	}
	result = string(marshalResult)
	return result, err
}

func formatTimeRecursively(object *interface{}) (err error) {
	if mapObject, isMap := (*object).(map[string]interface{}); isMap {
		for fieldName, fieldValue := range mapObject {
			if timeField, isTime := fieldValue.(time.Time); isTime {
				mapObject[fieldName] = timeField.Format(time.RFC3339)
			}
		}
	}
	return err
}
