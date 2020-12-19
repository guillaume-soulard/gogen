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
	var maxFileSizeInByte int64
	if maxFileSizeInByte, err =
	result = OutputFile{
		FileName:        fileName,
		MaxFileSizeByte: 0,
		file:            os.File{},
	}
	return result, err
}

type OutputFile struct {
	FileName        string `json:"fileName"`
	MaxFileSizeByte int64 `json:"maxFileSize"`
	file            os.File
}

func (o OutputFile) Write(object string) (err error) {
	dfgh
	return err
}
