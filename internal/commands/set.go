package commands

import (
	"fmt"
	"strings"

	"vault/internal/storage"
)

func HandleSet(args []string) {
	if len(args) < 4 {
		HandleHelp("set")
		return
	}

	data := storage.LoadData()

	key := args[2]
	value := args[3]

	if strings.Contains(key, "/") {
		parts := strings.SplitN(key, "/", 2)
		group := parts[0]
		subKey := parts[1]

		if _, ok := data[group]; !ok {
			data[group] = map[string]interface{}{}
		}

		groupMap := data[group].(map[string]interface{})
		groupMap[subKey] = value
	} else {
		data[key] = value
	}

	storage.SaveData(data)
	fmt.Println("Saved:", key)
}
