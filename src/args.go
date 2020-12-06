package src

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)

type Args struct {
	Files *[]os.File
}

func GetArgs() (result Args, err error) {
	result = Args{}
	parser := argparse.NewParser("gogen", "A seeded generation tool based on json file")
	parser.ExitOnHelp(true)
	parser.Usage("test")
	result.Files = parser.FileList("f", "file", os.O_RDONLY, 0600, &argparse.Options{
		Required: true,
		Help:     "the json files to use for the generation",
	})
	if err = parser.Parse(os.Args); err != nil {
		fmt.Println(parser.Usage(err))
	}
	return result, err
}
