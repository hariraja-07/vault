package commands

import (
	"fmt"
	"strings"

	"vault/internal/storage"
)

func HandleSet(args []string) {
	force := false
	key := ""
	value := ""

	for _, arg := range args[2:] {
		if strings.HasPrefix(arg, "--") {
			flag := strings.TrimPrefix(arg, "--")
			if flag == "force" || flag == "f" {
				force = true
			}
		} else if key == "" {
			key = arg
		} else if value == "" {
			value = arg
		}
	}

	_ = force

	if key == "" || value == "" {
		HandleHelp("set")
		return
	}

	data := storage.LoadData()

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

func init() {}
