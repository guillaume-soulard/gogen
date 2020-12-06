package main

import (
	"github.com/labstack/gommon/log"
	"github.com/ogama/gogen/src"
)

func main() {
	if args, err := src.GetArgs(); err == nil {
		if err := src.Execute(args); err != nil {
			log.Panic(err)
		}
	}
}
