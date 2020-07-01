package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func Test_samples(t *testing.T) {
	samples, err := getSamples()
	assert.NoError(t, err)
	for _, sample := range samples {
		fmt.Println(sample)
		assert.NoError(t, ExecuteFile("samples/"+sample), "Sample %s", sample)
	}
}

func getSamples() (samples []string, err error) {
	samples = make([]string, 0)
	var fileInfos []os.FileInfo
	fileInfos, err = ioutil.ReadDir("samples")
	for _, file := range fileInfos {
		samples = append(samples, file.Name())
	}
	return samples, err
}
