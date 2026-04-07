package commands

import (
	"fmt"
	"os"
	"strings"
)

var validShells = []string{"bash", "zsh", "fish", "powershell", "cmd"}

func HandleCompletion(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: vault completion <shell>")
		fmt.Println("Available shells: bash, zsh, fish, powershell, cmd")
		os.Exit(1)
	}

	shell := strings.ToLower(args[2])

	valid := false
	for _, s := range validShells {
		if s == shell {
			valid = true
			break
		}
	}

	if !valid {
		fmt.Printf("Error: Invalid shell '%s'\n", shell)
		fmt.Println("Available shells: bash, zsh, fish, powershell, cmd")
		os.Exit(1)
	}

	switch shell {
	case "bash":
		generateBashCompletion()
	case "zsh":
		generateZshCompletion()
	case "fish":
		generateFishCompletion()
	case "powershell":
		generatePowerShellCompletion()
	case "cmd":
		generateCmdCompletion()
	}
}

func generateBashCompletion() {
	completion := `#!/bin/bash

# vault Bash Completion
# Install: vault completion bash > /etc/bash_completion.d/vault
# Or: source <(vault completion bash)

_vault() {
    local cur prev words cword
    _init_completion -n=: || return

    local commands="set get remove list help completion"

    case "${words[1]}" in
        set|get|remove)
            local keys
            keys=$(vault list 2>/dev/null | grep -E '^[├─ └─]' | tr -d '├─ └─ ' | grep -v '/' | grep -v '^Vault$')
            COMPREPLY=($(compgen -W "$keys" -- "$cur"))
            ;;
        list)
            local groups
            groups=$(vault list 2>/dev/null | grep '/' | tr -d '├─ └─ /')
            COMPREPLY=($(compgen -W "$groups" -- "$cur"))
            COMPREPLY+=('--full' '-f')
            ;;
        help)
            COMPREPLY=($(compgen -W "set get remove list help completion" -- "$cur"))
            ;;
        completion)
            COMPREPLY=($(compgen -W "bash zsh fish powershell cmd" -- "$cur"))
            ;;
        *)
            COMPREPLY=($(compgen -W "$commands" -- "$cur"))
            ;;
    esac

    return 0
}

complete -F _vault vault
`
	fmt.Print(completion)
}

func generateZshCompletion() {
	completion := `#compdef vault

# vault Zsh Completion
# Install: vault completion zsh > ~/.zfunc/_vault

_vault() {
    local -a commands
    commands=('set' 'get' 'remove' 'list' 'help' 'completion')

    _arguments -C \
        '1:command:->command' \
        '2:key:->key' \
        '*:arg:->arg'

    case $state in
        command)
            _describe 'command' commands
            ;;
        key)
            case $words[1] in
                set|get|remove)
                    local keys
                    keys=($(vault list 2>/dev/null | grep -E '^[├─ └─]' | tr -d '├─ └─ ' | grep -v '/'))
                    _describe 'key' keys
                    ;;
                list)
                    local groups
                    groups=($(vault list 2>/dev/null | grep '/' | tr -d '├─ └─ /'))
                    _describe 'group' groups
                    _options \
                        '--full[show nested keys]' \
                        '-f[short for --full]'
                    ;;
            esac
            ;;
        arg)
            case $words[1] in
                help)
                    _describe 'command' commands
                    ;;
                completion)
                    _describe 'shell' 'bash' 'zsh' 'fish' 'powershell' 'cmd'
                    ;;
            esac
            ;;
    esac
}

_vault
`
	fmt.Print(completion)
}

func generateFishCompletion() {
	completion := `# vault Fish Completion
# Install: vault completion fish > ~/.config/fish/completions/vault.fish

complete -c vault -n '__fish_use_subcommand' -a 'set' -d 'Set a key-value pair'
complete -c vault -n '__fish_use_subcommand' -a 'get' -d 'Get a secret'
complete -c vault -n '__fish_use_subcommand' -a 'remove' -d 'Delete a key or group'
complete -c vault -n '__fish_use_subcommand' -a 'list' -d 'List all secrets'
complete -c vault -n '__fish_use_subcommand' -a 'help' -d 'Show help'
complete -c vault -n '__fish_use_subcommand' -a 'completion' -d 'Generate completion script'

# Key completions for set, get, remove
complete -c vault -n '__fish_seen_subcommand_from set; or __fish_seen_subcommand_from get; or __fish_seen_subcommand_from remove' -a '(__vault_keys)'

# Group completions for list
complete -c vault -n '__fish_seen_subcommand_from list' -a '(__vault_groups)'
complete -c vault -n '__fish_seen_subcommand_from list' -l full -d 'Show nested keys'
complete -c vault -n '__fish_seen_subcommand_from list' -s f -d 'Short for --full'

# Help completion
complete -c vault -n '__fish_seen_subcommand_from help' -a 'set get remove list help completion'

# Completion shell options
complete -c vault -n '__fish_seen_subcommand_from completion' -a 'bash zsh fish powershell cmd'

# Helper functions
function __vault_keys
    vault list 2>/dev/null | string replace -r '^[├─ └─ ]*' '' | string replace -r '/$' '' | string match -v ''
end

function __vault_groups
    vault list 2>/dev/null | string replace -r '^[├─ └─ ]*' '' | string match '*/'
end
`
	fmt.Print(completion)
}

func generatePowerShellCompletion() {
	completion := `# vault PowerShell completion
$vaultCompleter = {
    param($wordToComplete, $commandAst, $cursorPosition)

    $commands = @('set', 'get', 'remove', 'list', 'help', 'completion')
    $shells = @('bash', 'zsh', 'fish', 'powershell', 'cmd')

    if ($wordToComplete -match '^-') {
        return @('--force', '-F', '--full', '-f') | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_)
        }
    }

    if ($commandAst.CommandElements.Count -eq 1) {
        return $commands | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'Command', $_)
        }
    }

    if ($commandAst.CommandElements[1] -eq 'completion') {
        return $shells | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'Command', $_)
        }
    }
}

Register-ArgumentCompleter -CommandName vault -ScriptBlock $vaultCompleter
`
	fmt.Print(completion)
}

func generateCmdCompletion() {
	completion := `@echo off
rem vault CMD completion
doskey /completion:on
doskey /exename=vault

:vault_complete
if "%1"=="" goto :done
echo set get remove list help completion
:done
`
	fmt.Print(completion)
}

func init() {}
