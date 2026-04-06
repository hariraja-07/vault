# vault

> A simple, CLI-based key-value storage tool written in Go.

[![Go Version](https://img.shields.io/badge/Go-1.25-blue)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

---

## Features

- **Key-Value Storage** — Store and retrieve secrets
- **Grouped Organization** — Organize secrets by project (`work/api_key`)
- **JSON Persistence** — Data stored securely in JSON format
- **ASCII Tree View** — Clean, readable list output

---

## Installation

### Quick Install (One-Line)

**Linux / Mac:**
```bash
curl -fsSL https://raw.githubusercontent.com/hariraja-07/vault/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/hariraja-07/vault/main/scripts/install.ps1 | iex
```

---

### Manual Binary

Download the binary for your platform:

| Platform | Download |
|----------|----------|
| Linux | [vault-linux](https://github.com/hariraja-07/vault/releases/latest/download/vault-linux) |
| macOS Intel | [vault-macos-intel](https://github.com/hariraja-07/vault/releases/latest/download/vault-macos-intel) |
| macOS Apple Silicon | [vault-macos-arm](https://github.com/hariraja-07/vault/releases/latest/download/vault-macos-arm) |
| Windows | [vault-windows.exe](https://github.com/hariraja-07/vault/releases/latest/download/vault-windows.exe) |

After download:

```bash
# Linux / Mac
chmod +x vault-linux
mv vault-linux /usr/local/bin/vault

# Windows
# Rename to vault.exe and add to PATH
```

---

### Build from Source

Requires Go 1.25 or later.

```bash
git clone https://github.com/hariraja-07/vault.git
cd vault
make build
sudo mv vault /usr/local/bin/
```

---

## Quick Start

### Set a secret
```bash
vault set api_key sk_live_xxxxx
```

### Get a secret
```bash
vault get api_key
```

### Delete a secret
```bash
vault remove api_key
```

### List all secrets
```bash
vault list
vault list --full    # Show nested keys
```

---

## Commands

| Command | Description |
|---------|-------------|
| `vault set <key> <value> [--force]` | Set a key-value pair |
| `vault get <key>` | Get a secret |
| `vault remove <key>` | Delete a key or group |
| `vault list [--full]` | List all secrets |
| `vault help` | Show help |

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--force` | `-F` | Force overwrite existing key, group, or subkey |
| `--full` | `-f` | Show nested keys within groups |

---

## License

See [LICENSE](LICENSE) for details.

---

## Contact

For questions or feedback, reach out at:
- Email: hariraja1976@gmail.com
- LinkedIn: [hariharasudhan-rajendiran](https://www.linkedin.com/in/hariharasudhan-rajendiran-4b5b77331)
