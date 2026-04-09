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

# Install shell completions for all shells
echo ""
echo "Installing shell completions..."

# Bash completion
if [ -f /bin/bash ] || [ -f /usr/bin/bash ]; then
    mkdir -p ~/.bash_completion.d
    vault completion bash > ~/.bash_completion.d/vault 2>/dev/null && echo "✓ Bash completion installed"
    if ! grep -q "bash_completion.d/vault" ~/.bashrc 2>/dev/null; then
        echo 'source ~/.bash_completion.d/vault' >> ~/.bashrc 2>/dev/null && echo "✓ Added to ~/.bashrc"
    fi
fi

# Zsh completion
if [ -d ~/.zfunc ] || mkdir -p ~/.zfunc 2>/dev/null; then
    vault completion zsh > ~/.zfunc/_vault 2>/dev/null && echo "✓ Zsh completion installed"
fi

# Fish completion
if [ -d ~/.config/fish ] || mkdir -p ~/.config/fish/completions 2>/dev/null; then
    vault completion fish > ~/.config/fish/completions/vault.fish 2>/dev/null && echo "✓ Fish completion installed"
fi

echo ""
echo "Shell completions installed! Restart your shell or run:"
echo "  source ~/.bashrc  (for Bash)"
echo "  source ~/.zshrc   (for Zsh, if using oh-my-zsh)"
echo "  fish             (for Fish)"

echo ""
echo "Done! Run 'vault help' to get started."
