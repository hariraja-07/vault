package models

type Flag struct {
	Name        string
	Short       string
	Description string
}

type Command struct {
	Usage    string
	Desc     string
	Examples []string
	Flags    []Flag
}

var Commands = map[string]Command{
	"set": {
		Usage: "vault set <key> <value> [--force]",
		Desc:  "Set a key-value pair",
		Examples: []string{
			"vault set api_key sk_live_xxxxx",
			"vault set work/db_password secret123",
			"vault set db/pass1 value        # create group with subkey",
			"vault set db/pass1 newvalue     # error if exists",
			"vault set db/pass1 newvalue --force  # overwrite subkey",
		},
		Flags: []Flag{
			{Name: "force", Short: "F", Description: "Force overwrite existing key or subkey"},
		},
	},
	"get": {
		Usage: "vault get <key>",
		Desc:  "Get the value of a key",
		Examples: []string{
			"vault get api_key",
			"vault get work/db_password",
		},
	},
	"remove": {
		Usage: "vault remove <key>",
		Desc:  "Delete a key or group",
		Examples: []string{
			"vault remove api_key",
			"vault remove work/db_password",
			"vault remove db          # removes entire group",
		},
	},
	"list": {
		Usage: "vault list [group] [--full]",
		Desc:  "List all secrets or secrets in a specific group",
		Examples: []string{
			"vault list",
			"vault list --full        # show keys within groups",
			"vault list -f            # short form for --full",
			"vault list work          # show keys in 'work' group only",
		},
		Flags: []Flag{
			{Name: "full", Short: "f", Description: "Show nested keys within groups"},
		},
	},
	"help": {
		Usage: "vault help [command]",
		Desc:  "Show help for commands",
		Examples: []string{
			"vault help",
			"vault help set",
		},
	},
	"completion": {
		Usage: "vault completion <shell>",
		Desc:  "Generate shell completion script",
		Examples: []string{
			"vault completion bash",
			"vault completion zsh",
			"vault completion fish",
			"vault completion powershell",
			"vault completion cmd",
		},
	},
}
