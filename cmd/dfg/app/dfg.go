package app

import (
	"github.com/ozankasikci/dockerfile-generator/cmd/dfg/cmd"
)

// Entrypoint for the dfg command
func Run() error {
	cmd := cmd.NewDfgCommand()
	return cmd.Execute()
}
