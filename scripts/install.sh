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

echo "Done! Run 'vault help' to get started."
