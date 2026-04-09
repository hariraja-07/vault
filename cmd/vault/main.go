package main

import (
	"os"

	"vault/internal/commands"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		commands.HandleHelp()
		return
	}

	switch args[1] {
	case "set":
		commands.HandleSet(args)
	case "get":
		commands.HandleGet(args)
	case "remove":
		commands.HandleRemove(args)
	case "list":
		commands.HandleList(args)
	case "help":
		if len(args) > 2 {
			commands.HandleHelp(args[2])
		} else {
			commands.HandleHelp()
		}
	case "completion":
		commands.HandleCompletion(args)
	case "config":
		commands.HandleConfig(args)
	default:
		commands.HandleHelp()
	}
}
