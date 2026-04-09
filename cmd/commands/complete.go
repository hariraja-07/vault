package commands

import (
	"github.com/spf13/cobra"
	"vault/internal/storage"
)

// completeKeys provides completion suggestions for keys from JSON storage
func completeKeys(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	keys := storage.GetAllKeys()
	return keys, cobra.ShellCompDirectiveNoFileComp
}

// RegisterKeyCompletion registers key completion for a command
func RegisterKeyCompletion(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc("key", completeKeys)
}
