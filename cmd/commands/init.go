package commands

import (
	"vault/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(SetCmd)
	cmd.RootCmd.AddCommand(GetCmd)
	cmd.RootCmd.AddCommand(RemoveCmd)
	cmd.RootCmd.AddCommand(ListCmd)
	cmd.RootCmd.AddCommand(CompletionCmd)
}
