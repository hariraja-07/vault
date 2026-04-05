package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: vault <command> [arguments]")
		return
	}

	data := loadData()

	switch args[1] {

	case "set":
		handleSet(args, data)

	case "get":
		handleGet(args, data)

	case "remove":
		handleRemove(args, data)

	case "list":
		full := false
		group := ""

		for _, arg := range args[2:] {
			if arg == "--full" {
				full = true
			} else {
				group = arg
			}
		}
		handleList(data, full, group)

	case "help":
		if len(args) > 2 {
			handleHelp(args[2])
		} else {
			handleHelp()
		}

	default:
		fmt.Println("Unknown command")
	}
}

func handleSet(args []string, data map[string]interface{}) {
	if len(args) < 4 {
		handleHelp("set", "short")
		return
	}

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

	saveData(data)
	fmt.Println("Saved:", key)
}

func handleGet(args []string, data map[string]interface{}) {
	if len(args) < 3 {
		handleHelp("get", "short")
		return
	}
	key := args[2]

	// grouped key
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

	// only key
	val, ok := data[key]
	if !ok {
		fmt.Println("Key not found")
		return
	}
	fmt.Println(val)
}

func handleRemove(args []string, data map[string]interface{}) {
	if len(args) < 3 {
		handleHelp("remove", "short")
		return
	}
	key := args[2]

	//grouped key
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

		saveData(data)
		fmt.Println("Deleted:", key)
		return
	}

	//normal key
	if _, exists := data[key]; !exists {
		fmt.Println("Key not found")
		return
	}
	delete(data, key)
	saveData(data)

	fmt.Println("Deleted :", key)
}

func handleList(data map[string]interface{}, full bool, group string) {
	if group != "" {
		groupMap, ok := data[group].(map[string]interface{})
		if !ok {
			fmt.Println("Group not found")
			return
		}

		if len(groupMap) == 0 {
			fmt.Println("Group is empty")
			return
		}

		var keys []string
		for key := range groupMap {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		fmt.Println(group + "/")

		for i, key := range keys {
			isLast := i == len(keys)-1

			prefix := "├── "
			if isLast {
				prefix = "└── "
			}

			fmt.Printf("%s%s\n", prefix, key)
		}

		return
	}

	if len(data) == 0 {
		fmt.Println("No data stored")
		return
	}

	var keys []string
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Println("Vault")
	for i, key := range keys {
		value := data[key]

		isLast := i == len(keys)-1

		prefix := "├── "
		if isLast {
			prefix = "└── "
		}

		if groupMap, ok := value.(map[string]interface{}); ok {
			fmt.Printf("%s %s/\n", prefix, key)
			if full {
				var subKeys []string
				for subKey := range groupMap {
					subKeys = append(subKeys, subKey)
				}
				sort.Strings(subKeys)

				for j, subKey := range subKeys {
					isSubLast := j == len(subKeys)-1

					if isLast {
						if isSubLast {
							fmt.Printf("    └── %s\n", subKey)
						} else {
							fmt.Printf("    ├── %s\n", subKey)
						}
					} else {
						if isSubLast {
							fmt.Printf("│   └── %s\n", subKey)
						} else {
							fmt.Printf("│   ├── %s\n", subKey)
						}
					}
				}
			}
		} else {
			fmt.Printf("%s %s\n", prefix, key)
		}
	}
}

func handleHelp(command ...string) {
	if len(command) == 0 {
		fmt.Println("Vault commmands:")
		for name, cmd := range commands {
			fmt.Printf("%-10s : %s\n", name, cmd.Desc)
		}
		return
	}

	cmdName := command[0]

	cmd, ok := commands[cmdName]

	if !ok {
		fmt.Println("Unknown command:", cmdName)
		return
	}

	if len(command) > 1 && command[1] == "short" {
		fmt.Printf("Usage: %s\n", cmd.Usage)
		return
	}

	fmt.Printf("Command: %s\n", cmdName)
	fmt.Printf("Description: %s\n", cmd.Desc)
	fmt.Printf("Usage: %s\n", cmd.Usage)
}
