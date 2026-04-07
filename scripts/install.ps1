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

# Install PowerShell completion
Write-Output ""
Write-Output "Installing PowerShell completion..."

$ProfilePath = $PROFILE
if (-not (Test-Path $ProfilePath)) {
    New-Item -Path $ProfilePath -ItemType File -Force | Out-Null
}

$CompletionLine = "Invoke-Expression -Command `$(vault completion powershell)"
if (-not (Select-String -Path $ProfilePath -Pattern "vault completion" -Quiet)) {
    Add-Content -Path $ProfilePath -Value ""
    Add-Content -Path $ProfilePath -Value "# vault completion"
    Add-Content -Path $ProfilePath -Value $CompletionLine
    Write-Output "PowerShell completion added to your profile"
}

# Install CMD completion
Write-Output ""
Write-Output "Installing CMD completion..."

$CmdFile = "$env:USERPROFILE\vault_complete.bat"
$CmdContent = @"
@echo off
rem vault CMD completion
doskey /exename=vault.exe

if "%%1"=="set" goto :set
if "%%1"=="get" goto :get
if "%%1"=="remove" goto :remove
if "%%1"=="list" goto :list
if "%%1"=="help" goto :help
if "%%1"=="completion" goto :completion
goto :end

:set
:get
:remove
:list
:help
:completion
echo set get remove list help completion --force -F --full -f

:end
"@

Set-Content -Path $CmdFile -Value $CmdContent -Force
Write-Output "CMD completion installed to $CmdFile"
Write-Output "To enable, run: $CmdFile"

# Add to CMD autorun
$AutorunPath = "HKCU:\Software\Microsoft\Command Processor"
$AutorunValue = Get-ItemProperty -Path $AutorunPath -Name AutoRun -ErrorAction SilentlyContinue
if (-not $AutorunValue) {
    Set-ItemProperty -Path $AutorunPath -Name AutoRun -Value "CALL `"$CmdFile`""
    Write-Output "CMD completion enabled for all CMD sessions"
} else {
    Write-Output "CMD autorun already set. Completion file installed to $CmdFile"
}

Write-Output ""
Write-Output "Done! Run 'vault help' to get started."
