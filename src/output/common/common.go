package common

import "github.com/ogama/gogen/src/configuration"

type Output interface {
	Begin() (err error)
	Write(object string) (err error)
	End() (err error)
}

type Builder interface {
	Build(configuration configuration.OutputConfiguration) (result Output, err error)
}
