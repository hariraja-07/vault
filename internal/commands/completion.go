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

_vault_completion() {
    local cur prev words cword
    _init_completion -n=: || return

    commands="set get remove list help completion"

    if [[ $cword -eq 1 ]]; then
        COMPREPLY=($(compgen -W "$commands" -- "$cur"))
    elif [[ $cword -eq 2 ]]; then
        case "${words[1]}" in
            set|get|remove)
                COMPREPLY=($(compgen -W "$(vault list 2>/dev/null | grep -v 'Vault' | grep -v '/' | tr -d '├─ └─ ')" -- "$cur"))
                ;;
        esac
    fi
}

complete -F _vault_completion vault
`
	fmt.Print(completion)
}

func generateZshCompletion() {
	completion := `#compdef vault

_vault() {
    local -a commands
    commands=('set' 'get' 'remove' 'list' 'help' 'completion')

    _arguments -C \\
        '1: :->command' \\
        '2: :->arg' \\
        '*:: :->args'

    case $state in
        command)
            _describe 'command' commands
            ;;
        arg)
            case $words[1] in
                set|get|remove)
                    _message "Enter key name"
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
	completion := `complete -c vault -n '__fish_use_subcommand' -a 'set get remove list help completion'
complete -c vault -n '__fish_seen_subcommand_from set' -a '(__vault_keys)'
complete -c vault -n '__fish_seen_subcommand_from get' -a '(__vault_keys)'
complete -c vault -n '__fish_seen_subcommand_from remove' -a '(__vault_keys)'
complete -c vault -n '__fish_seen_subcommand_from completion' -a 'bash zsh fish powershell cmd'

function __vault_keys
    vault list 2>/dev/null | string replace -r '^[├─ └─ ]*' '' | string replace -r '/$' ''
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
