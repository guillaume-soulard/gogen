package common

import "github.com/ogama/gogen/model/configuration"

type Output interface {
	Write(object string) (err error)
}

type Builder interface {
	Build(configuration configuration.OutputConfiguration) (result Output, err error)
}
