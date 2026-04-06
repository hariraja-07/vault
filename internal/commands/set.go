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

	if key == "" || value == "" {
		HandleHelp("set")
		return
	}

	data := storage.LoadData()

	if strings.Contains(key, "/") {
		parts := strings.SplitN(key, "/", 2)
		group := parts[0]
		subKey := parts[1]

		if existingGroup, exists := data[group]; exists {
			if !storage.IsGroup(existingGroup) {
				if !force {
					fmt.Printf("Error: Key '%s' already exists.\n", group)
					fmt.Println("Use --force to overwrite.")
					return
				}
				delete(data, group)
				fmt.Printf("Warning: overwriting key '%s'\n", group)
			}
		}

		if _, exists := data[group]; !exists {
			data[group] = map[string]interface{}{}
		}

		groupMap := data[group].(map[string]interface{})
		groupMap[subKey] = value
	} else {
		if existingValue, exists := data[key]; exists {
			if !force {
				if storage.IsGroup(existingValue) {
					groupMap := existingValue.(map[string]interface{})
					count := len(groupMap)

					fmt.Printf("Error: Group '%s' already exists with %d nested key(s).\n", key, count)
					fmt.Println("Use --force to delete all nested keys and overwrite.")
					return
				}
				fmt.Printf("Error: Key '%s' already exists.\n", key)
				fmt.Println("Use --force to overwrite.")
				return
			}

			if storage.IsGroup(existingValue) {
				groupMap := existingValue.(map[string]interface{})
				count := len(groupMap)
				fmt.Printf("Warning: overwriting group '%s' (%d nested key(s) will be deleted)\n", key, count)
			} else {
				fmt.Printf("Warning: overwriting key '%s'\n", key)
			}
			delete(data, key)
		}

		data[key] = value
	}

	storage.SaveData(data)
	fmt.Println("Saved:", key)
}

func init() {}
