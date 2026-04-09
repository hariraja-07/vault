package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var validShells = []string{"bash", "zsh", "fish", "powershell", "cmd"}

func getCurrentShell() string {
	shell := detectShellFromPS()
	if shell != "" {
		return shell
	}

	shell = detectShellFromEnv()
	if shell != "" {
		return shell
	}

	return os.Getenv("SHELL")
}

func detectShellFromPS() string {
	pid := os.Getpid()
	for i := 0; i < 10; i++ {
		cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", pid), "-o", "comm=")
		output, _ := cmd.Output()
		name := strings.TrimSpace(string(output))

		if name == "bash" || name == "zsh" || name == "fish" {
			return name
		}

		cmd = exec.Command("ps", "-o", "ppid=", "-p", fmt.Sprintf("%d", pid))
		output, _ = cmd.Output()
		newPid, _ := strconv.Atoi(strings.TrimSpace(string(output)))
		if newPid == 0 || newPid == 1 {
			break
		}
		pid = newPid
	}
	return ""
}

func detectShellFromEnv() string {
	shells := []string{"fish", "zsh", "bash"}
	for _, shell := range shells {
		if path, err := exec.LookPath(shell); err == nil && path != "" {
			return shell
		}
	}
	return ""
}

func HandleCompletion(args []string) {
	if len(args) < 3 {
		installForCurrentShell()
		return
	}

	subcmd := strings.ToLower(args[2])
	userHome, _ := os.UserHomeDir()

	switch subcmd {
	case "install":
		if len(args) >= 4 {
			shell := strings.ToLower(args[3])
			installForShell(userHome, shell)
		} else {
			installForCurrentShell()
		}
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
	case "help", "--help", "-h":
		showCompletionHelp()
	default:
		fmt.Println("Usage: vault completion [bash|zsh|fish|powershell|cmd]")
		fmt.Println("Run without arguments to auto-install for your shell.")
		fmt.Println("Examples:")
		fmt.Println("  vault completion          # Auto-install for current shell")
		fmt.Println("  vault completion bash     # Output Bash completion script")
		fmt.Println("  vault completion zsh     # Output Zsh completion script")
		fmt.Println("  vault completion fish     # Output Fish completion script")
	}
}

func installForShell(home string, shell string) {
	switch shell {
	case "bash":
		installBashCompletion(home)
	case "zsh":
		installZshCompletion(home)
	case "fish":
		installFishCompletion(home)
	case "powershell":
		installPowerShellCompletion(home)
	case "cmd":
		installCmdCompletion(home)
	default:
		fmt.Printf("Error: Unsupported shell: %s\n", shell)
	}
}

func installForCurrentShell() {
	shell := getCurrentShell()
	userHome, _ := os.UserHomeDir()

	switch {
	case strings.Contains(shell, "bash"):
		installBashCompletion(userHome)
	case strings.Contains(shell, "zsh"):
		installZshCompletion(userHome)
	case strings.Contains(shell, "fish"):
		installFishCompletion(userHome)
	default:
		fmt.Printf("Error: Unsupported shell: %s\n", shell)
		fmt.Println("Supported shells: Bash, Zsh, Fish")
		fmt.Println("For manual installation, use: vault completion <shell>")
	}
}

func installBashCompletion(home string) {
	bashCompletionDir := filepath.Join(home, ".bash_completion.d")
	bashCompletionFile := filepath.Join(bashCompletionDir, "vault")

	fmt.Println("Installing Bash completion...")

	os.MkdirAll(bashCompletionDir, 0755)
	os.WriteFile(bashCompletionFile, []byte(generateBashCompletionString()), 0644)

	bashrcPath := filepath.Join(home, ".bashrc")
	bashrcContent, _ := os.ReadFile(bashrcPath)
	bashrc := string(bashrcContent)

	sourceLine := "\n# vault shell completion\n[ -f ~/.bash_completion.d/vault ] && source ~/.bash_completion.d/vault\n"

	if !strings.Contains(bashrc, "bash_completion.d/vault") {
		os.WriteFile(bashrcPath, []byte(bashrc+sourceLine), 0644)
	}

	fmt.Println("✓ Installed to ~/.bash_completion.d/vault")
	fmt.Println("✓ Added to ~/.bashrc")
	fmt.Println("✓ Shell completion setup complete!")
	fmt.Println("→ Restart your terminal or run: source ~/.bashrc")
}

func installZshCompletion(home string) {
	zfuncDir := filepath.Join(home, ".zfunc")
	zshCompletionFile := filepath.Join(zfuncDir, "_vault")

	fmt.Println("Installing Zsh completion...")

	os.MkdirAll(zfuncDir, 0755)
	os.WriteFile(zshCompletionFile, []byte(generateZshCompletionString()), 0644)

	zshrcPath := filepath.Join(home, ".zshrc")
	if _, err := os.Stat(zshrcPath); os.IsNotExist(err) {
		f, _ := os.Create(zshrcPath)
		f.Close()
	}
	zshrcContent, _ := os.ReadFile(zshrcPath)
	zshrc := string(zshrcContent)

	setupLines := ""

	if !strings.Contains(zshrc, "compinit") {
		setupLines += "\n# Initialize zsh completion system\nautoload -Uz compinit\ncompinit\n"
	}

	if !strings.Contains(zshrc, ".zfunc") {
		setupLines += "\n# vault shell completion\nfpath+=(~/.zfunc)\nautoload -Uz _vault\ncompdef _vault vault\n"
	}

	if setupLines != "" {
		os.WriteFile(zshrcPath, []byte(zshrc+setupLines), 0644)
	}

	fmt.Println("✓ Installed to ~/.zfunc/_vault")
	fmt.Println("✓ Updated ~/.zshrc")
	fmt.Println("✓ Shell completion setup complete!")
	fmt.Println("→ Restart your terminal or run: source ~/.zshrc")
}

func installFishCompletion(home string) {
	fishCompletionDir := filepath.Join(home, ".config", "fish", "completions")
	fishCompletionFile := filepath.Join(fishCompletionDir, "vault.fish")

	fmt.Println("Installing Fish completion...")

	os.MkdirAll(fishCompletionDir, 0755)
	os.WriteFile(fishCompletionFile, []byte(generateFishCompletionString()), 0644)

	fmt.Println("✓ Installed to ~/.config/fish/completions/vault.fish")
	fmt.Println("✓ Shell completion setup complete!")
	fmt.Println("→ Restart your terminal or run: fish")
}

func installPowerShellCompletion(home string) {
	psProfileDir := filepath.Join(home, "Documents", "WindowsPowerShell")
	psProfileFile := filepath.Join(psProfileDir, "Microsoft.PowerShell_profile.ps1")

	fmt.Println("Installing PowerShell completion...")

	os.MkdirAll(psProfileDir, 0755)

	var psProfileContent string
	if _, err := os.Stat(psProfileFile); os.IsNotExist(err) {
		psProfileContent = "# PowerShell profile\n# vault completion\nInvoke-Expression -Command $(vault completion powershell)\n"
	} else {
		existing, _ := os.ReadFile(psProfileFile)
		content := string(existing)
		if !strings.Contains(content, "vault completion") {
			psProfileContent = content + "\n# vault completion\nInvoke-Expression -Command $(vault completion powershell)\n"
		} else {
			psProfileContent = content
		}
	}

	os.WriteFile(psProfileFile, []byte(psProfileContent), 0644)

	fmt.Println("✓ Installed to PowerShell profile")
	fmt.Println("✓ Shell completion setup complete!")
	fmt.Println("→ Restart PowerShell to use completions")
}

func installCmdCompletion(home string) {
	cmdCompletionDir := filepath.Join(home, ".config", "vault")
	cmdCompletionFile := filepath.Join(cmdCompletionDir, "vault_complete.bat")

	fmt.Println("Installing CMD completion...")

	os.MkdirAll(cmdCompletionDir, 0755)
	os.WriteFile(cmdCompletionFile, []byte(generateCmdCompletionString()), 0644)

	fmt.Println("✓ Installed to %USERPROFILE%\\.config\\vault\\vault_complete.bat")
	fmt.Println("✓ Registry configured for CMD autostart")
	fmt.Println("✓ Shell completion setup complete!")
	fmt.Println("→ Restart CMD to use completions")
}

func showCompletionHelp() {
	fmt.Println("vault completion - Install shell completion")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  vault completion              Auto-install for current shell")
	fmt.Println("  vault completion bash        Install for Bash")
	fmt.Println("  vault completion zsh        Install for Zsh")
	fmt.Println("  vault completion fish       Install for Fish")
	fmt.Println("  vault completion powershell  Install for PowerShell")
	fmt.Println("  vault completion cmd        Install for CMD")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  vault completion            # Auto-detect and install")
	fmt.Println("  vault completion bash      # Install for Bash")
	fmt.Println("  vault completion zsh       # Install for Zsh")
	fmt.Println("  vault completion fish      # Install for Fish")
}

func execCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, _ := cmd.Output()
	return string(output)
}

func generateBashCompletion() {
	fmt.Print(generateBashCompletionString())
}

func generateBashCompletionString() string {
	return `#!/bin/bash

# vault Bash Completion (standalone - no dependencies)
# Install: vault completion bash > ~/.bash_completion.d/vault
# Or: source <(vault completion bash)

_vault() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    local word_count=${#COMP_WORDS[@]}
    local cmd="${COMP_WORDS[1]}"

    # Define valid commands
    local commands="set get remove list help completion config"
    local list_cmds="bash zsh fish powershell cmd"

    # Helper function to get recent keys
    _vault_get_keys() {
        local output
        output=$(vault list --recent 2>/dev/null)
        if [ -z "$output" ]; then
            output=$(vault list 2>/dev/null)
        fi
        echo "$output" | grep -E '^[[:space:]]*├|^[[:space:]]*└' | awk '{print $NF}' | grep -v '/' | tr '\n' ' '
    }

    # Helper function to get all keys
    _vault_get_all_keys() {
        local output
        output=$(vault list 2>/dev/null)
        echo "$output" | grep -E '^[[:space:]]*├|^[[:space:]]*└' | awk '{print $NF}' | grep -v '/' | tr '\n' ' '
    }

    # Helper function to get groups
    _vault_get_groups() {
        local output
        output=$(vault list 2>/dev/null)
        echo "$output" | grep -E '^[[:space:]]*├|^[[:space:]]*└' | awk '{print $NF}' | grep '/' | sed 's/\/$//' | tr '\n' ' '
    }

    # Helper function to combine recent + all keys for completion
    _vault_get_completion_keys() {
        local keys=$(_vault_get_keys)
        if [ -z "$keys" ]; then
            keys=$(_vault_get_all_keys)
        fi
        echo "$keys"
    }

    # Level 1: vault <TAB> - show commands only
    # Only show commands if position 1 is empty, partial match, or not a valid command
    if [ $COMP_CWORD -eq 1 ]; then
        # cmd is COMP_WORDS[1] - the actual command being typed
        # Check if cmd is a valid command
        if [ -n "$cmd" ]; then
            local match=$(echo "$commands" | tr ' ' '\n' | grep -w "^$cmd")
            if [ -n "$match" ]; then
                # cmd IS a valid command, show Level 2 (keys + flags)
                :
            else
                # cmd is partial command name, show matching commands
                COMPREPLY=($(compgen -W "$commands" -- "$cur"))
                return
            fi
        else
            # cmd is empty, show all commands
            COMPREPLY=($(compgen -W "$commands" -- "$cur"))
            return
        fi
    fi

    # Level 3: vault <cmd> <key> <TAB> - show flags or filter keys
    if [ $word_count -ge 3 ]; then
        case "$cmd" in
            set)
                COMPREPLY=($(compgen -W "--force -F" -- "$cur"))
                ;;
            get|remove)
                local keys=$(_vault_get_completion_keys)
                COMPREPLY=($(compgen -W "$keys" -- "$cur"))
                ;;
        esac
        return
    fi

    # Level 2: vault <cmd> <TAB> - show keys + flags
    # If cur equals cmd, the user just typed the command - show all completions
    # Otherwise, cur contains a partial key for filtering
    local filter_cur="$cur"
    if [ "$cur" = "$cmd" ]; then
        filter_cur=""
    fi

    case "$cmd" in
        set)
            local keys=$(_vault_get_completion_keys)
            COMPREPLY=($(compgen -W "$keys --force -F" -- "$filter_cur"))
            ;;
        get|remove)
            local keys=$(_vault_get_completion_keys)
            COMPREPLY=($(compgen -W "$keys" -- "$filter_cur"))
            ;;
        list)
            local groups=$(_vault_get_groups)
            COMPREPLY=($(compgen -W "$groups --full -f --recent" -- "$filter_cur"))
            ;;
        help)
            COMPREPLY=($(compgen -W "$commands" -- "$filter_cur"))
            ;;
        completion)
            COMPREPLY=($(compgen -W "$list_cmds" -- "$filter_cur"))
            ;;
        config)
            COMPREPLY=($(compgen -W "--recent" -- "$filter_cur"))
            ;;
    esac
}

complete -F _vault vault
`
}

func generateZshCompletion() {
	fmt.Print(generateZshCompletionString())
}

func generateZshCompletionString() string {
	return `#compdef vault

# vault Zsh Completion
# Install: vault completion zsh > ~/.zfunc/_vault

_vault() {
    local -a commands
    commands=(set get remove list help completion config)

    # Get keys from vault list
    _get_vault_keys() {
        local output
        output=$(vault list --recent 2>/dev/null)
        if [ -z "$output" ]; then
            output=$(vault list 2>/dev/null)
        fi
        echo "$output" | grep -E '^[[:space:]]*├|^[[:space:]]*└' | awk '{print $NF}' | grep -v '/'
    }

    # Get groups from vault list
    _get_vault_groups() {
        local output
        output=$(vault list 2>/dev/null)
        echo "$output" | grep -E '^[[:space:]]*├|^[[:space:]]*└' | awk '{print $NF}' | grep '/' | sed 's/\/$//'
    }

    local curcontext="$curcontext" state
    typeset -A opt_args

    _arguments -C \
        '1: :->cmds' \
        '2: :->keys' \
        '3: :->args'

    case $state in
        cmds)
            _describe 'command' commands
            ;;
        keys)
            local cmd=$words[1]
            case $cmd in
                set)
                    local -a keys
                    keys=($_get_vault_keys)
                    if [ ${#keys[@]} -gt 0 ]; then
                        _describe 'key' keys
                    fi
                    _arguments -s \
                        '--force[force overwrite]' \
                        '-F[force overwrite]'
                    ;;
                get|remove)
                    local -a keys
                    keys=($_get_vault_keys)
                    if [ ${#keys[@]} -gt 0 ]; then
                        _describe 'key' keys
                    fi
                    ;;
                list)
                    local -a groups
                    groups=($_get_vault_groups)
                    if [ ${#groups[@]} -gt 0 ]; then
                        _describe 'group' groups
                    fi
                    _arguments -s \
                        '--full[show nested keys]' \
                        '--recent[show recent keys]'
                    ;;
                help)
                    _describe 'command' commands
                    ;;
                completion)
                    _describe 'shell' bash zsh fish powershell cmd
                    ;;
            esac
            ;;
        args)
            local cmd=$words[1]
            case $cmd in
                set)
                    _arguments -s \
                        '--force[force overwrite]' \
                        '-F[force overwrite]'
                    ;;
            esac
            ;;
    esac
}
`
}

func generateFishCompletion() {
	fmt.Print(generateFishCompletionString())
}

func generateFishCompletionString() string {
	return `# vault Fish Completion
# Install: vault completion fish > ~/.config/fish/completions/vault.fish

# Level 1: vault <TAB> - commands only (exclude files)
complete -c vault -n '__fish_use_subcommand' -f -a 'set' -d 'Set a key-value pair'
complete -c vault -n '__fish_use_subcommand' -f -a 'get' -d 'Get a secret'
complete -c vault -n '__fish_use_subcommand' -f -a 'remove' -d 'Delete a key or group'
complete -c vault -n '__fish_use_subcommand' -f -a 'list' -d 'List all secrets'
complete -c vault -n '__fish_use_subcommand' -f -a 'help' -d 'Show help'
complete -c vault -n '__fish_use_subcommand' -f -a 'completion' -d 'Generate completion script'
complete -c vault -n '__fish_use_subcommand' -f -a 'config' -d 'Manage configuration'

# Level 2: vault set/get/remove <TAB> - keys + flags
complete -c vault -n '__fish_seen_subcommand_from set' -a '(__vault_all_keys)' -f
complete -c vault -n '__fish_seen_subcommand_from set' -l force -d 'Force overwrite existing key'
complete -c vault -n '__fish_seen_subcommand_from set' -s F -d 'Force overwrite existing key'
complete -c vault -n '__fish_seen_subcommand_from get' -a '(__vault_all_keys)' -f
complete -c vault -n '__fish_seen_subcommand_from remove' -a '(__vault_all_keys)' -f

# Level 2: vault list <TAB> - groups + flags
complete -c vault -n '__fish_seen_subcommand_from list' -a '(__vault_groups)' -f
complete -c vault -n '__fish_seen_subcommand_from list' -l full -d 'Show nested keys'
complete -c vault -n '__fish_seen_subcommand_from list' -s f -d 'Short for --full'
complete -c vault -n '__fish_seen_subcommand_from list' -l recent -d 'Show recent keys'

# Level 3: vault set <key> <TAB> - flags only
complete -c vault -n '__fish_seen_subcommand_from set; and __fish_nth_token_since 2' -l force -d 'Force overwrite'
complete -c vault -n '__fish_seen_subcommand_from set; and __fish_nth_token_since 2' -s F -d 'Force overwrite'

# Help and completion
complete -c vault -n '__fish_seen_subcommand_from help' -f -a 'set get remove list help completion config'
complete -c vault -n '__fish_seen_subcommand_from completion' -f -a 'bash zsh fish powershell cmd'

# Helper functions - parse vault list output to extract clean keys
function __vault_all_keys
    set -l cmd (commandline -opc)
    if test (count $cmd) -ge 2
        set -l subcmd $cmd[2]
        # Skip if already at key position (3rd word)
        if test (count $cmd) -ge 3
            return
        end
    end
    
    set -l output (vault list --recent 2>/dev/null)
    if test -z "$output"
        set output (vault list 2>/dev/null)
    end
    
    for line in $output
        if string match -q '*├*' $line
            set -l key (string replace -r '.*[├]─[[:space:]]*' '' $line)
            if test -n "$key"
                and string match -qv '/' $key
                echo $key
            end
        else if string match -q '*└*' $line
            set -l key (string replace -r '.*[└]─[[:space:]]*' '' $line)
            if test -n "$key"
                and string match -qv '/' $key
                echo $key
            end
        end
    end
end

function __vault_groups
    set -l output (vault list 2>/dev/null)
    
    for line in $output
        if string match -q '*├*' $line
            set -l key (string replace -r '.*[├]─[[:space:]]*' '' $line)
            if test -n "$key"
                and string match -q '*/*' $key
                string replace -r '/$' '' $key
            end
        else if string match -q '*└*' $line
            set -l key (string replace -r '.*[└]─[[:space:]]*' '' $line)
            if test -n "$key"
                and string match -q '*/*' $key
                string replace -r '/$' '' $key
            end
        end
    end
end
`
}

func generatePowerShellCompletion() {
	completion := `# vault PowerShell completion
$vaultCompleter = {
    param($wordToComplete, $commandAst, $cursorPosition)

    $commands = @('set', 'get', 'remove', 'list', 'help', 'completion', 'config')
    $shells = @('bash', 'zsh', 'fish', 'powershell', 'cmd')

    $wordCount = $commandAst.CommandElements.Count

    # Helper function to parse keys from vault list output
    function Get-VaultKeys {
        $output = vault list --recent 2>$null
        if (-not $output) {
            $output = vault list 2>$null
        }
        $keys = @()
        $lines = $output -split [Environment]::NewLine
        foreach ($line in $lines) {
            if ($line -match '├|└') {
                $key = ($line -split '├' -split '└' | Select-Object -Last 1).Trim()
                if ($key -and $key -notmatch '/') {
                    $keys += $key
                }
            }
        }
        return $keys
    }

    # Helper function to get all keys
    function Get-VaultAllKeys {
        $output = vault list 2>$null
        $keys = @()
        $lines = $output -split [Environment]::NewLine
        foreach ($line in $lines) {
            if ($line -match '├|└') {
                $key = ($line -split '├' -split '└' | Select-Object -Last 1).Trim()
                if ($key -and $key -notmatch '/') {
                    $keys += $key
                }
            }
        }
        return $keys
    }

    # Helper function to get groups
    function Get-VaultGroups {
        $output = vault list 2>$null
        $groups = @()
        $lines = $output -split [Environment]::NewLine
        foreach ($line in $lines) {
            if ($line -match '├|└') {
                $key = ($line -split '├' -split '└' | Select-Object -Last 1).Trim()
                if ($key -match '/') {
                    $groups += ($key -replace '/$', '')
                }
            }
        }
        return $groups
    }

    # Check if word at position 1 is already a known command
    if ($wordCount -ge 2) {
        $cmd = $commandAst.CommandElements[1].Value
        $cmdIsKnown = $commands -contains $cmd
        
        # If command is known and we're at word position 2, show keys
        if ($cmdIsKnown -and $wordCount -eq 2) {
            switch ($cmd) {
                'set' {
                    $completions = @()
                    $keys = Get-VaultKeys
                    if ($keys.Count -eq 0) { $keys = Get-VaultAllKeys }
                    $keys | ForEach-Object {
                        $completions += [System.Management.Automation.CompletionResult]::new($_, $_, 'Argument', $_)
                    }
                    $completions += [System.Management.Automation.CompletionResult]::new('--force', '--force', 'ParameterName', 'Force overwrite existing key')
                    $completions += [System.Management.Automation.CompletionResult]::new('-F', '-F', 'ParameterName', 'Force overwrite existing key')
                    return $completions
                }
                'get' {
                    $keys = Get-VaultKeys
                    if ($keys.Count -eq 0) { $keys = Get-VaultAllKeys }
                    return $keys | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'Argument', $_) }
                }
                'remove' {
                    $keys = Get-VaultKeys
                    if ($keys.Count -eq 0) { $keys = Get-VaultAllKeys }
                    return $keys | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'Argument', $_) }
                }
                'list' {
                    $completions = @()
                    $groups = Get-VaultGroups
                    $groups | ForEach-Object {
                        $completions += [System.Management.Automation.CompletionResult]::new($_, $_, 'Argument', $_)
                    }
                    $completions += [System.Management.Automation.CompletionResult]::new('--full', '--full', 'ParameterName', 'Show nested keys')
                    $completions += [System.Management.Automation.CompletionResult]::new('-f', '-f', 'ParameterName', 'Short for --full')
                    $completions += [System.Management.Automation.CompletionResult]::new('--recent', '--recent', 'ParameterName', 'Show recent keys')
                    return $completions
                }
                'help' {
                    return $commands | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'Command', $_) }
                }
                'completion' {
                    return $shells | ForEach-Object { [System.Management.Automation.CompletionResult]::new($_, $_, 'Command', $_) }
                }
            }
        }
        
        # If command is known and we're at word position 3, show flags only
        if ($cmdIsKnown -and $wordCount -ge 3) {
            if ($cmd -eq 'set') {
                return @('--force', '-F') | ForEach-Object {
                    [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterName', $_)
                }
            }
            return
        }
    }

    # Level 1: vault <TAB> - show commands only
    if ($wordCount -eq 1) {
        return $commands | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'Command', $_)
        }
    }
}

Register-ArgumentCompleter -CommandName vault -ScriptBlock $vaultCompleter
`
	fmt.Print(completion)
}

func generateCmdCompletion() {
	fmt.Print(generateCmdCompletionString())
}

func generateCmdCompletionString() string {
	return `@echo off
rem vault CMD completion script
rem This script registers vault completions in the Windows registry
rem For full experience, run this on CMD startup
rem
rem NOTE: CMD completion uses a static keys.txt file for key suggestions.
rem To update keys, run: vault completion cmd

setlocal enabledelayedexpansion

set "VAULT_CONFIG=%USERPROFILE%\.config\vault"
set "VAULT_KEYS=%VAULT_CONFIG%\keys.txt"
set "COMPLETION_REGISTRY=HKEY_CURRENT_USER\Software\Microsoft\Command Processor\Completion"

rem Create config directory if not exists
if not exist "%VAULT_CONFIG%" mkdir "%VAULT_CONFIG%"

rem Update keys.txt from vault list
for /f "tokens=*" %%k in ('vault list 2^>nul') do (
    echo %%k >> "%VAULT_KEYS%.tmp"
)

rem Process and deduplicate keys
if exist "%VAULT_KEYS%.tmp" (
    sort "%VAULT_KEYS%.tmp" /unique > "%VAULT_KEYS%" 2>nul
    del "%VAULT_KEYS%.tmp" 2>nul
)

rem Clear old vault completions
reg query "%COMPLETION_REGISTRY%" /v "vault_cmd" 2>nul >nul
if !errorlevel! equ 0 (
    for /f "tokens=1,2*" %%a in ('reg query "%COMPLETION_REGISTRY%" /s 2^>nul ^| findstr /i "vault_"') do (
        reg delete "%COMPLETION_REGISTRY%" /v "%%a" /f >nul 2>&1
    )
)

rem Register command words
reg add "%COMPLETION_REGISTRY%\vault_cmd" /ve /d "set get remove list help completion config" /f >nul 2>&1

rem Register keys from keys.txt (exclude groups with /)
if exist "%VAULT_KEYS%" (
    for /f "usebackq tokens=*" %%k in ("%VAULT_KEYS%") do (
        echo %%k | findstr "/" >nul
        if !errorlevel! neq 0 (
            reg add "%COMPLETION_REGISTRY%" /v "vault_%%k" /d "%%k" /f >nul 2>&1
        )
    )
)

rem Register flags
reg add "%COMPLETION_REGISTRY%" /v "vault_--force" /d "--force" /f >nul 2>&1
reg add "%COMPLETION_REGISTRY%" /v "vault_-F" /d "-F" /f >nul 2>&1
reg add "%COMPLETION_REGISTRY%" /v "vault_--full" /d "--full" /f >nul 2>&1
reg add "%COMPLETION_REGISTRY%" /v "vault_-f" /d "-f" /f >nul 2>&1
reg add "%COMPLETION_REGISTRY%" /v "vault_--recent" /d "--recent" /f >nul 2>&1

rem Set autorun
reg add "HKEY_CURRENT_USER\Software\Microsoft\Command Processor" /v AutoRun /t REG_SZ /d "call \"%USERPROFILE%\.config\vault\vault_complete.bat\"" /f >nul 2>&1

echo vault CMD completion registered successfully
echo Restart CMD to enable completions
endlocal
`
}

func init() {}
