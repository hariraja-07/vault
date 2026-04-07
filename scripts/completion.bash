#!/bin/bash

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
