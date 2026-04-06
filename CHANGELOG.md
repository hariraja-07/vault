# Changelog

All notable changes to this project will be documented in this file.

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
