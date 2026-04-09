package commands

import (
	"fmt"
	"strconv"
	"strings"

	"vault/internal/storage"
)

func HandleConfig(args []string) {
	if len(args) < 3 {
		showConfigHelp()
		return
	}

	action := strings.ToLower(args[2])

	switch action {
	case "set":
		handleConfigSet(args[3:])
	case "get":
		handleConfigGet(args[3:])
	default:
		fmt.Println("Unknown config action:", action)
		showConfigHelp()
	}
}

func handleConfigSet(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: vault config set <key> <value>")
		fmt.Println("Example: vault config set recent-limit 10")
		return
	}

	key := strings.ToLower(args[0])
	value := args[1]

	switch key {
	case "recent-limit", "recentlimit":
		limit, err := strconv.Atoi(value)
		if err != nil || limit < 1 {
			fmt.Println("Error: recent-limit must be a positive number")
			return
		}
		storage.SetRecentLimit(limit)
		fmt.Printf("Set recent-limit to %d\n", limit)
	default:
		fmt.Printf("Error: Unknown config key '%s'\n", key)
		fmt.Println("Available keys: recent-limit")
	}
}

func handleConfigGet(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: vault config get <key>")
		fmt.Println("Example: vault config get recent-limit")
		return
	}

	key := strings.ToLower(args[0])

	switch key {
	case "recent-limit", "recentlimit":
		limit := storage.GetRecentLimit()
		fmt.Printf("recent-limit: %d\n", limit)
	default:
		fmt.Printf("Error: Unknown config key '%s'\n", key)
		fmt.Println("Available keys: recent-limit")
	}
}

func showConfigHelp() {
	fmt.Println("vault config - Manage vault configuration")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  vault config get <key>       Get a config value")
	fmt.Println("  vault config set <key> <value>  Set a config value")
	fmt.Println()
	fmt.Println("Available config keys:")
	fmt.Println("  recent-limit    Number of recent keys to show in completion (default: 10)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  vault config get recent-limit")
	fmt.Println("  vault config set recent-limit 20")
}
