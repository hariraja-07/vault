Write-Output "Installing vault..."

# Detect architecture
$Architecture = $env:PROCESSOR_ARCHITECTURE
if ($Architecture -eq "ARM64") {
    $Binary = "vault-macos-arm"
} else {
    $Binary = "vault-windows.exe"
}

# Get latest release URL
$ReleaseUrl = "https://github.com/hariraja-07/vault/releases/latest/download/$Binary"

# Download
Write-Output "Downloading $Binary..."
Invoke-WebRequest -Uri $ReleaseUrl -OutFile "$env:TEMP\vault.exe"

# Install to a location in PATH
$InstallDir = "$env:LOCALAPPDATA\Programs\vault"
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Move-Item -Path "$env:TEMP\vault.exe" -Destination "$InstallDir\vault.exe"

# Add to PATH if not already
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    Write-Output "Added $InstallDir to your PATH"
    Write-Output "Please restart your terminal for changes to take effect"
}

Write-Output "Installed to $InstallDir\vault.exe"
Write-Output "Done! Run 'vault help' to get started."
