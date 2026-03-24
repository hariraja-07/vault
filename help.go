package main

type Command struct {
	Usage string
	Desc  string
}

var commands = map[string]Command{
	"set": {
		Usage: "vault set <key> <value>",
		Desc:  "Set a key with a value",
	},
	"get": {
		Usage: "vault get <key>",
		Desc:  "Get the value of a key",
	},
	"remove": {
		Usage: "vault remove <key>",
		Desc:  "Delete a key",
	},
	"list": {
		Usage: "vault list",
		Desc:  "List all stored keys",
	},
	"help": {
		Usage: "vault help",
		Desc:  "show help for all commands",
	},
}
