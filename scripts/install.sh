#!/bin/bash

set -e

echo "Installing vault..."

# Detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

# Set binary name based on OS and architecture
case "$OS" in
    Linux*)
        BINARY="vault-linux"
        ;;
    Darwin*)
        if [ "$ARCH" = "arm64" ]; then
            BINARY="vault-macos-arm"
        else
            BINARY="vault-macos-intel"
        fi
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Get latest release URL
RELEASE_URL="https://github.com/hariraja-07/vault/releases/latest/download/${BINARY}"

# Download
echo "Downloading ${BINARY}..."
curl -fsSL "$RELEASE_URL" -o /tmp/vault

# Make executable
chmod +x /tmp/vault

# Install to /usr/local/bin
if [ -w /usr/local/bin ]; then
    mv /tmp/vault /usr/local/bin/vault
    echo "Installed to /usr/local/bin/vault"
else
    echo "Installing to ~/.local/bin..."
    mkdir -p ~/.local/bin
    mv /tmp/vault ~/.local/bin/vault
    echo "Installed to ~/.local/bin/vault"
    echo "Add ~/.local/bin to your PATH if not already added"
fi

# Install shell completion
echo ""
echo "Installing shell completion..."

install_bash_completion() {
    if command -v apt-get &>/dev/null; then
        if ! dpkg -l bash-completion &>/dev/null 2>&1; then
            echo "Installing bash-completion package..."
            sudo apt-get update -qq && sudo apt-get install -y bash-completion 2>/dev/null && return 0
        fi
    elif command -v dnf &>/dev/null; then
        if ! rpm -q bash-completion &>/dev/null 2>&1; then
            echo "Installing bash-completion package..."
            sudo dnf install -y bash-completion 2>/dev/null && return 0
        fi
    elif command -v brew &>/dev/null; then
        if ! brew list bash-completion@2 &>/dev/null 2>&1; then
            echo "Installing bash-completion package..."
            brew install bash-completion@2 2>/dev/null && return 0
        fi
    fi
    return 1
}

case "$SHELL" in
    */bash)
        if ! install_bash_completion; then
            echo "Could not auto-install bash-completion (may need sudo or package manager)"
            echo "Using standalone completion script..."
        else
            echo "Bash completion package installed"
        fi

        if [ -d /etc/bash_completion.d ] && [ -w /etc/bash_completion.d ]; then
            vault completion bash > /etc/bash_completion.d/vault 2>/dev/null && echo "Bash completion installed"
        elif [ -w ~/.bash_completion.d ] 2>/dev/null || mkdir -p ~/.bash_completion.d 2>/dev/null; then
            vault completion bash > ~/.bash_completion.d/vault && echo "Bash completion installed"
        else
            echo 'source <(vault completion bash)' >> ~/.bashrc && echo "Bash completion added to ~/.bashrc"
        fi
        ;;
    */zsh)
        if [ -d ~/.zfunc ]; then
            vault completion zsh > ~/.zfunc/_vault && echo "Zsh completion installed"
        else
            mkdir -p ~/.zfunc && vault completion zsh > ~/.zfunc/_vault && echo "Zsh completion installed"
        fi
        ;;
    */fish)
        mkdir -p ~/.config/fish/completions && vault completion fish > ~/.config/fish/completions/vault.fish && echo "Fish completion installed"
        ;;
    *)
        echo "Shell completion not supported for $SHELL"
        echo "You can manually run: vault completion <shell>"
        ;;
esac

echo ""
echo "Done! Run 'vault help' to get started."
