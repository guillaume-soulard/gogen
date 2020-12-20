package stdout

import (
	"fmt"
	"github.com/ogama/gogen/src/configuration"
	commonOutput "github.com/ogama/gogen/src/output/common"
)

type BuilderStdout struct{}

func (b BuilderStdout) Build(_ configuration.OutputConfiguration) (result commonOutput.Output, err error) {
	result = OutputStdout{}
	return result, err
}

type OutputStdout struct{}

func (o OutputStdout) Begin() (err error) {
	return err
}

func (o OutputStdout) Write(object string) (err error) {
	fmt.Println(object)
	return err
}

func (o OutputStdout) End() (err error) {
	return err
}
