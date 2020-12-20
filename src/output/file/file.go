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
	result = &OutputFile{
		fileName: fileName,
	}
	return result, err
}

type OutputFile struct {
	fileName string
	file     *os.File
}

func (o *OutputFile) Begin() (err error) {
	if _, err := os.Stat(o.fileName); err == nil {
		if err = os.Remove(o.fileName); err != nil {
			return err
		}
	}
	o.file, err = os.Create(o.fileName)
	return err
}

func (o *OutputFile) Write(object string) (err error) {
	_, err = o.file.Write([]byte(object))
	return err
}

func (o *OutputFile) End() (err error) {
	err = o.file.Close()
	return err
}
