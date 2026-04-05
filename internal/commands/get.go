package commands

import (
	"fmt"
	"strings"

	"vault/internal/storage"
)

func HandleGet(args []string) {
	if len(args) < 3 {
		HandleHelp("get")
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
		val, ok := groupMap[subKey]
		if !ok {
			fmt.Println("Key not found")
			return
		}

		fmt.Println(val)
		return
	}

	val, ok := data[key]
	if !ok {
		fmt.Println("Key not found")
		return
	}
	fmt.Println(val)
}
