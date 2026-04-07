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

case "$SHELL" in
    */bash)
        if [ -d /etc/bash_completion.d ]; then
            vault completion bash > /etc/bash_completion.d/vault 2>/dev/null && echo "Bash completion installed"
        elif [ -w ~/.bash_completion ]; then
            vault completion bash >> ~/.bash_completion && echo "Bash completion installed"
        elif [ -f ~/.bashrc ]; then
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
