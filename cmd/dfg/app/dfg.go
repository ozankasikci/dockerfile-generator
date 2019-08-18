package app

import (
	"github.com/ozankasikci/dockerfile-generator/cmd/dfg/cmd"
)

func Run() error {
	cmd := cmd.NewDfgCommand()
	return cmd.Execute()
}
