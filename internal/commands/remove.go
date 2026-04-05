package commands

import (
	"fmt"
	"strings"

	"vault/internal/storage"
)

func HandleRemove(args []string) {
	if len(args) < 3 {
		HandleHelp("remove")
		return
	}

	data := storage.LoadData()
	key := args[2]

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
		fmt.Println("Deleted:", key)
		return
	}

	if _, exists := data[key]; !exists {
		fmt.Println("Key not found")
		return
	}
	delete(data, key)
	storage.SaveData(data)

	fmt.Println("Deleted:", key)
}
