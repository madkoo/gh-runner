# PowerShell entrypoint for GitHub Actions Self-Hosted Runner (Windows)

$ErrorActionPreference = "Stop"

# Required environment variables
$RepoUrl = $env:REPO_URL
$AccessToken = $env:ACCESS_TOKEN
$RunnerName = if ($env:RUNNER_NAME) { $env:RUNNER_NAME } else { $env:COMPUTERNAME }
$RunnerWorkdir = if ($env:RUNNER_WORKDIR) { $env:RUNNER_WORKDIR } else { "_work" }
$Labels = $env:LABELS

if (-not $RepoUrl) {
    Write-Error "Error: REPO_URL environment variable is required"
    exit 1
}

if (-not $AccessToken) {
    Write-Error "Error: ACCESS_TOKEN environment variable is required"
    exit 1
}

Write-Host "Configuring GitHub Actions Runner..."
Write-Host "  Repository: $RepoUrl"
Write-Host "  Runner name: $RunnerName"

# Build configuration command
$configArgs = @(
    "--url", $RepoUrl,
    "--token", $AccessToken,
    "--name", $RunnerName,
    "--work", $RunnerWorkdir,
    "--unattended",
    "--replace"
)

if ($Labels) {
    $configArgs += @("--labels", $Labels)
    Write-Host "  Labels: $Labels"
}

# Configure the runner
& .\config.cmd @configArgs

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to configure runner"
    exit 1
}

# Cleanup function
$cleanup = {
    Write-Host "Removing runner..."
    & .\config.cmd remove --token $AccessToken
}

# Register cleanup on exit
Register-EngineEvent -SourceIdentifier PowerShell.Exiting -Action $cleanup

Write-Host "Starting runner..."
& .\run.cmd
