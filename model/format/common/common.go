package common

import "github.com/ogama/gogen/model/configuration"

type Format interface {
	Format(object interface{}) (result string, err error)
}

type Builder interface {
	Build(configuration configuration.FormatConfiguration) (result Format, err error)
}
