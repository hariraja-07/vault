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
		handleSet(args, data)

	case "get":
		handleGet(args, data)

	case "remove":
		handleRemove(args, data)

	case "list":
		handleList(data)

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

	data[key] = value
	saveData(data)

	fmt.Println("Saved:", key)
}

func handleGet(args []string, data map[string]interface{}) {
	if len(args) < 3 {
		handleHelp("get", "short")
		return
	}
	key := args[2]
	value, exists := data[key]
	if exists {
		fmt.Println(value)
	} else {
		fmt.Println("Key not found")
	}
}

func handleRemove(args []string, data map[string]interface{}) {
	if len(args) < 3 {
		handleHelp("remove", "short")
		return
	}

	key := args[2]

	delete(data, key)
	saveData(data)

	fmt.Println("Deleted :", key)
}

func handleList(data map[string]interface{}) {
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

		if i == len(keys)-1 {
			fmt.Printf("└── %s\n", key)
		} else {
			fmt.Printf("├── %s\n", key)
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
