package cmd

import (
	"github.com/spf13/cobra"
)

// Generates a dfg command
func NewDfgCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "dfg",
		Short: "dfg: a dockerfile generator",
		Long:  "dfg: a dockerfile generator",
	}

	cmds.ResetFlags()
	cmds.AddCommand(NewCmdGenerate())

	return cmds
}
