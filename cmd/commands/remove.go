package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"vault/internal/storage"
)

var RemoveCmd = &cobra.Command{
	Use:   "remove <key>",
	Short: "Delete a key or group",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeKeys(cmd, args, toComplete)
	},
	Run: func(cmd *cobra.Command, args []string) {
		data := storage.LoadData()
		key := args[0]

		if strings.Contains(key, "/") {
			if strings.Count(key, "/") > 1 {
				fmt.Println("Error: Only one level grouping allowed (group/key)")
				return
			}

			parts := strings.SplitN(key, "/", 2)
			group := parts[0]
			subKey := parts[1]

			groupMap, ok := data[group].(map[string]interface{})
			if !ok {
				fmt.Println("Group not found")
				return
			}

			if _, exists := groupMap[subKey]; !exists {
				fmt.Println("Key not found")
				return
			}
			delete(groupMap, subKey)

			if len(groupMap) == 0 {
				delete(data, group)
			}

			storage.SaveData(data)
			storage.TrackKeyUsage(key)
			fmt.Println("Deleted:", key)
			return
		}

		if _, exists := data[key]; !exists {
			fmt.Println("Key not found")
			return
		}
		delete(data, key)
		storage.SaveData(data)
		storage.TrackKeyUsage(key)

		fmt.Println("Deleted:", key)
	},
}
