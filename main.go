package main

import (
	"fmt"
	"github.com/ogama/gogen/src"
)

func main() {
	if args, err := src.GetArgs(); err == nil {
		if err := src.Execute(args); err != nil {
			fmt.Println(err.Error())
		}
	}
}
