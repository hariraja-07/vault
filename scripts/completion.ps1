# vault PowerShell Completion
# Install: Add to your PowerShell profile
# Add this line to your profile: Invoke-Expression $(vault completion powershell)

using namespace System.Management.Automation
using namespace System.Management.Automation.Language

$script:VaultCommands = @('set', 'get', 'remove', 'list', 'help', 'completion')
$script:VaultFlags = @('--force', '-F', '--full', '-f')
$script:VaultShells = @('bash', 'zsh', 'fish', 'powershell', 'cmd')

function Get-VaultKeys {
    $output = vault list 2>$null
    if ($output) {
        $output | Where-Object { $_ -match '^[├─ └─]' } | ForEach-Object {
            $_ -replace '^[├─ └─]\s*', '' -replace '/.*$', ''
        } | Where-Object { $_ -and $_ -notmatch '/' }
    }
}

function Get-VaultGroups {
    $output = vault list 2>$null
    if ($output) {
        $output | Where-Object { $_ -match '/' } | ForEach-Object {
            $_ -replace '^[├─ └─]\s*', '' -replace '/.*', ''
        } | Sort-Object -Unique
    }
}

Register-ArgumentCompleter -CommandName vault -ParameterName key -ScriptBlock {
    param($wordToComplete, $commandName, $cursorPosition)
    Get-VaultKeys | ForEach-Object {
        [CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}

Register-ArgumentCompleter -CommandName vault -ParameterName group -ScriptBlock {
    param($wordToComplete, $commandName, $cursorPosition)
    Get-VaultGroups | ForEach-Object {
        [CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}

Register-ArgumentCompleter -CommandName vault -NativeArgument '--force' -ScriptBlock {
    param($wordToComplete, $commandName, $cursorPosition)
    $VaultFlags | ForEach-Object {
        [CompletionResult]::new($_, $_, 'Flag', $_)
    }
}

Register-ArgumentCompleter -CommandName vault -NativeArgument '-F' -ScriptBlock {
    param($wordToComplete, $commandName, $cursorPosition)
    $VaultFlags | ForEach-Object {
        [CompletionResult]::new($_, $_, 'Flag', $_)
    }
}

Register-ArgumentCompleter -CommandName vault -NativeArgument '--full' -ScriptBlock {
    param($wordToComplete, $commandName, $cursorPosition)
    $VaultFlags | ForEach-Object {
        [CompletionResult]::new($_, $_, 'Flag', $_)
    }
}

Register-ArgumentCompleter -CommandName vault -NativeArgument '-f' -ScriptBlock {
    param($wordToComplete, $commandName, $cursorPosition)
    $VaultFlags | ForEach-Object {
        [CompletionResult]::new($_, $_, 'Flag', $_)
    }
}
