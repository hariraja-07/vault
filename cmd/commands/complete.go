package commands

import (
	"strings"

	"github.com/spf13/cobra"
	"vault/internal/storage"
)

// completeKeys provides completion suggestions for keys from JSON storage
func completeKeys(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Check if typing a group path (contains /)
	if strings.Contains(toComplete, "/") {
		return completeInnerKeys(toComplete)
	}

	// Default: show all keys and groups
	keys := storage.GetAllKeys()
	return keys, cobra.ShellCompDirectiveNoFileComp
}

// completeInnerKeys returns inner keys for a specific group
func completeInnerKeys(toComplete string) ([]string, cobra.ShellCompDirective) {
	parts := strings.SplitN(toComplete, "/", 2)
	group := parts[0]
	prefix := ""
	if len(parts) > 1 {
		prefix = parts[1]
	}

	data := storage.LoadData()

	// Check if group exists and is actually a group
	groupMap, ok := data[group].(map[string]interface{})
	if !ok {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Get inner keys
	var innerKeys []string
	for key := range groupMap {
		if prefix == "" || strings.HasPrefix(key, prefix) {
			innerKeys = append(innerKeys, key)
		}
	}

	return innerKeys, cobra.ShellCompDirectiveNoFileComp
}

// RegisterKeyCompletion registers key completion for a command
func RegisterKeyCompletion(cmd *cobra.Command) {
	cmd.RegisterFlagCompletionFunc("key", completeKeys)
}
