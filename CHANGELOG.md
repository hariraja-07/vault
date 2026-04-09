# Changelog

All notable changes to this project will be documented in this file.

## [v0.2.0] - 2026-04-08

### Added
- **Recent Keys Tracking** - Automatically tracks recently used keys for smart suggestions
- **Configuration System** - New `vault config` command to manage settings
- **`--recent` flag** - `vault list --recent [n]` to show recent keys
- **Smart Shell Completion** - Context-aware completions with recent keys support
- **CMD Full Completion** - Registry-based tab completion for Windows CMD

### Changes
- Improved shell completions for Bash, Zsh, Fish, PowerShell, and CMD
- Commands-only completion at first level (no file suggestions)
- Keys + flags shown at appropriate completion levels
- Partial key matching supported in completions

### Configuration Options
- `recent-limit` - Set number of recent keys to show in completion (default: 10)

### Examples
```bash
vault config get recent-limit
vault config set recent-limit 20
vault list --recent
vault list --recent 5
```

## [v0.1.3] - 2026-04-06
### Changes
- `-F` short form for `--force` flag
- `-f` short form for `--full` flag

## [v0.1.2] - 2026-04-06
### Added
- Require `--force` flag to update existing subkeys in `vault set` command
- Consistent behavior with flat keys

### Changes
- `vault set group/key value` → Creates new subkey
- `vault set group/key value` → Error if subkey exists
- `vault set group/key value --force` → Updates with warning

## [v0.1.1] - 2026-04-06
### Added
- `--force` flag for `vault set` command
- Conflict detection for keys and groups
- Safe overwrite to prevent accidental data loss

### Changes
- `vault set key value` → Error if key exists
- `vault set key value --force` → Allows overwrite
- Helpful warnings show how many nested keys will be deleted

## [v0.1] - 2026-04-05
### Added
- Basic CRUD operations (set, get, remove, list)
- Group support for organizing secrets
- JSON persistence in `~/.vault/`
- ASCII tree view for listing
