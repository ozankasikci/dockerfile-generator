package main

import (
	"github.com/ozankasikci/dockerfile-generator/cmd/dfg/app"
	"os"
)

func main() {
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
