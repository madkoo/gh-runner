# Dockerfile for GitHub Actions Self-Hosted Runner (Windows)
# Uses official GitHub Actions runner binaries
# Requires Windows containers mode in Docker Desktop

# escape=`

FROM mcr.microsoft.com/windows/servercore:ltsc2022

SHELL ["powershell", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]

# Install Git
RUN Invoke-WebRequest -Uri 'https://github.com/git-for-windows/git/releases/download/v2.43.0.windows.1/Git-2.43.0-64-bit.exe' -OutFile 'git-installer.exe'; `
    Start-Process -FilePath 'git-installer.exe' -ArgumentList '/VERYSILENT', '/NORESTART', '/NOCANCEL', '/SP-', '/CLOSEAPPLICATIONS', '/RESTARTAPPLICATIONS', '/COMPONENTS=icons,ext\reg\shellhere,assoc,assoc_sh' -Wait; `
    Remove-Item -Force 'git-installer.exe'

# Set up runner directory
WORKDIR C:\actions-runner

# Download and extract the latest runner
RUN $releases = Invoke-RestMethod -Uri 'https://api.github.com/repos/actions/runner/releases/latest'; `
    $version = $releases.tag_name -replace 'v', ''; `
    $downloadUrl = \"https://github.com/actions/runner/releases/download/v$version/actions-runner-win-x64-$version.zip\"; `
    Invoke-WebRequest -Uri $downloadUrl -OutFile 'actions-runner.zip'; `
    Expand-Archive -Path 'actions-runner.zip' -DestinationPath .; `
    Remove-Item -Force 'actions-runner.zip'

# Copy entrypoint script
COPY entrypoint.ps1 C:\entrypoint.ps1

ENTRYPOINT ["powershell", "-File", "C:\\entrypoint.ps1"]
