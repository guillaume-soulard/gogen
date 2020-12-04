package stdout

import (
	"fmt"
	"github.com/ogama/gogen/model/configuration"
	commonOutput "github.com/ogama/gogen/model/output/common"
)

type BuilderStdout struct{}

func (b BuilderStdout) Build(_ configuration.OutputConfiguration) (result commonOutput.Output, err error) {
	result = OutputStdout{}
	return result, err
}

type OutputStdout struct{}

func (o OutputStdout) Write(object string) (err error) {
	fmt.Println(object)
	return err
}
