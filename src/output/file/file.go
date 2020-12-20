package file

import (
	"github.com/ogama/gogen/src/configuration"
	"github.com/ogama/gogen/src/output/common"
	"os"
)

type BuilderFile struct{}

func (b BuilderFile) Build(configuration configuration.OutputConfiguration) (result common.Output, err error) {
	var fileName string
	if fileName, err = configuration.Options.GetStringOrDefault("fileName", ""); err != nil {
		return result, err
	}
	var file *os.File
	if file, err = os.Open(fileName); err != nil {
		return result, err
	}
	result = OutputFile{
		file: file,
	}
	return result, err
}

type OutputFile struct {
	file *os.File
}

func (o OutputFile) Begin() (err error) {
	return err
}

func (o OutputFile) Write(object string) (err error) {
	dfgh
	return err
}

func (o OutputFile) End() (err error) {
	return err
}
