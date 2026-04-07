#compdef vault

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
