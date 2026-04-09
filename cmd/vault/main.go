package main

import (
	"vault/cmd"
	_ "vault/cmd/commands"
)

func main() {
	cmd.Execute()
}
