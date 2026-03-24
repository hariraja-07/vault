package main

import (
	"fmt"
	"os"
	"sort"
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
		if len(args) < 4 {
			fmt.Println("Usage: vault set <key> <value>")
			return
		}

		key := args[2]
		value := args[3]

		data[key] = value
		saveData(data)

		fmt.Println("Saved:", key)

	case "get":
		if len(args) < 3 {
			fmt.Println("Usage: vault get <key>")
			return
		}
		key := args[3]
		value, exists := data[key]
		if exists {
			fmt.Println(value)
		} else {
			fmt.Println("Key not found")
		}

	case "remove":
		if len(args) < 3 {
			fmt.Println("Usage: vault remove <key>")
			return
		}

		key := args[2]

		delete(data, key)
		saveData(data)

		fmt.Println("Deleted :", key)

	case "list":
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

			if i == len(keys)-1 {
				fmt.Printf("└── %s: %v\n", key, value)
			} else {
				fmt.Printf("├── %s: %v\n", key, value)
			}
		}

	default:
		fmt.Println("Unknown command")
	}
}
