package commands

import (
	"fmt"

	"vault/internal/models"
)

func HandleHelp(command ...string) {
	if len(command) == 0 || (len(command) > 0 && command[0] == "help") {
		fmt.Println()
		fmt.Println("vault - A simple CLI-based key-value storage tool")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  vault <command> [arguments]")
		fmt.Println()
		fmt.Println("Commands:")
		for name, cmd := range models.Commands {
			fmt.Printf("  %-10s %s\n", name, cmd.Desc)
		}
		fmt.Println()
		fmt.Println(`Use "vault help <command>" for more information about a command.`)
		fmt.Println()
		return
	}

	cmdName := command[0]

	cmd, ok := models.Commands[cmdName]

	if !ok {
		fmt.Println("Unknown command:", cmdName)
		return
	}

	if len(command) > 1 && command[1] == "short" {
		fmt.Printf("Usage: %s\n", cmd.Usage)
		return
	}

	fmt.Println()
	fmt.Println(cmd.Desc)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s\n", cmd.Usage)

	if len(cmd.Flags) > 0 {
		fmt.Println()
		fmt.Println("Flags:")
		for _, flag := range cmd.Flags {
			shortForm := ""
			if flag.Short != "" {
				shortForm = fmt.Sprintf(" (-%s)", flag.Short)
			}
			fmt.Printf("  --%s%s    %s\n", flag.Name, shortForm, flag.Description)
		}
	}

	if len(cmd.Examples) > 0 {
		fmt.Println()
		fmt.Println("Examples:")
		for _, example := range cmd.Examples {
			fmt.Printf("  %s\n", example)
		}
	}
	fmt.Println()
}
