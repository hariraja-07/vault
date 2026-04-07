# vault Fish Completion
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
