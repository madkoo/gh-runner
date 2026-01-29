# gh-runner

> A GitHub CLI extension for effortlessly adding self-hosted runners to your repositories, organizations, or enterprises.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 📋 Table of Contents

- [What is gh-runner?](#what-is-gh-runner)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
  - [Basic Commands](#basic-commands)
  - [Advanced Usage](#advanced-usage)
- [Examples](#examples)
- [How It Works](#how-it-works)
- [Managing Runners](#managing-runners)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [Development](#development)
- [Roadmap](#roadmap)
- [License](#license)

## What is gh-runner?

`gh-runner` is a command-line tool that simplifies the process of setting up GitHub Actions self-hosted runners. Instead of manually downloading runner binaries, configuring tokens, and managing deployments, this extension automates the entire process with a single command.

**Why use self-hosted runners?**
- Run workflows on your own hardware with specific configurations
- Access to GPUs, specialized software, or network resources
- Better control over the execution environment
- Cost savings for high-volume CI/CD pipelines
- Compliance with security and privacy requirements

## Features

- ✅ **Simple CLI** - Add runners with a single command
- ✅ **Multiple Scopes** - Support for repositories, organizations, and enterprises
- ✅ **Docker Integration** - Automatically deploy runners in Docker containers
- ✅ **Windows Support** - Deploy Windows runners with `--windows` flag
- ✅ **Custom Configuration** - Set custom names and labels
- ✅ **Secure** - Uses GitHub CLI authentication and token management
- ✅ **Cross-Platform** - Works on Linux, macOS, and Windows
- ✅ **Official Runner Binaries** - Uses official GitHub Actions runner directly
- ⏳ **Kubernetes Support** - Coming soon for k8s/minikube deployments

## Prerequisites

### Required
- **GitHub CLI** (`gh`) - [Installation guide](https://cli.github.com/)
- **Git** - For version control
- **Go 1.21+** - For building from source

### Optional
- **Docker** - Required for `--docker` flag ([Installation guide](https://www.docker.com/))
- **Minikube** - Required for `--k8s` flag (coming soon)

### Permissions
You need **admin access** to the repository, organization, or enterprise where you want to add runners.

## Installation

### Method 1: Install from GitHub (Recommended)

```bash
gh extension install madkoo/gh-runner
```

### Method 2: Install Locally

```bash
# Clone the repository
git clone https://github.com/madkoo/gh-runner.git
cd gh-runner

# Install the extension
gh extension install .
```

### Method 3: Build from Source

```bash
# Clone and build
git clone https://github.com/madkoo/gh-runner.git
cd gh-runner
make build

# Install
gh extension install .
```

### Verify Installation

```bash
gh runner --help
```

## Usage

### Basic Commands

#### Add a Runner to a Repository

```bash
gh runner add OWNER/REPO --docker
```

Example:
```bash
gh runner add mycompany/my-private-repo --docker
```

#### Add a Runner to an Organization

```bash
gh runner add org/ORG-NAME --docker
```

Example:
```bash
gh runner add org/my-company --docker
```

#### Add a Runner to an Enterprise

```bash
gh runner add enterprise/ENTERPRISE-NAME --docker
```

Example:
```bash
gh runner add enterprise/my-enterprise --docker
```

### Advanced Usage

#### Custom Runner Name

```bash
gh runner add owner/repo --docker --name production-runner-01
```

#### Add Custom Labels

```bash
gh runner add owner/repo --docker --labels linux,x64,gpu,cuda,large-disk
```

Labels help you target specific runners in your workflows:
```yaml
jobs:
  build:
    runs-on: [self-hosted, linux, gpu, cuda]
```

#### Combine Options

```bash
gh runner add org/myorg --docker \
  --name gpu-runner-01 \
  --labels linux,x64,gpu,nvidia-rtx-4090,128gb-ram
```

#### Windows Runner

```bash
gh runner add myorg/myrepo --docker --windows
```

**Note:** Windows containers require:
- Docker Desktop with Windows containers enabled
- Running on a Windows host
- Switch to Windows containers via Docker Desktop settings

#### Manual Setup (Without Docker)

```bash
gh runner add owner/repo
```

This displays manual installation instructions with the registration token.

## Examples

### Example 1: Quick Docker Runner for a Private Repo

```bash
gh runner add mycompany/private-api --docker
```

Output:
```
✓ Registration token obtained
✓ Runner scope: Repository
✓ Target: mycompany/private-api

🐳 Setting up Docker runner...
🐧 Using Linux runner configuration (slim Ubuntu image)
✓ Docker container started: gh-runner-MacBook-Pro
✓ Runner name: MacBook-Pro
✓ Platform: Linux

To view logs: docker logs -f gh-runner-MacBook-Pro
To stop: docker stop gh-runner-MacBook-Pro
To remove: docker rm -f gh-runner-MacBook-Pro
```

### Example 2: Organization Runner with Custom Labels

```bash
gh runner add org/myorg --docker \
  --name build-server-01 \
  --labels linux,x64,build,production
```

### Example 3: Development Runner

```bash
gh runner add myuser/test-project --docker \
  --name dev-runner \
  --labels linux,x64,development
```

### Example 4: Windows Runner

```bash
gh runner add mycompany/windows-app --docker --windows \
  --name windows-runner \
  --labels windows,x64,dotnet
```

### Example 5: Using the Runner in a Workflow

After adding your runner, use it in `.github/workflows/ci.yml`:

```yaml
name: CI

on: [push, pull_request]

jobs:
  build:
    runs-on: [self-hosted, linux, x64]
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Build
        run: make build
      
      - name: Test
        run: make test
```

## How It Works

1. **Authentication** - Uses your existing GitHub CLI credentials
2. **Token Generation** - Requests a registration token from GitHub API
3. **Deployment** - Either:
   - Starts a Docker container with the runner pre-configured (with `--docker`)
   - Provides manual setup instructions with the token
4. **Registration** - The runner automatically registers with GitHub
5. **Ready** - Runner appears online in your Settings → Actions → Runners

### Architecture

```
┌─────────────┐
│   gh CLI    │
│   (User)    │
└──────┬──────┘
       │
       │ gh runner add owner/repo --docker
       │
┌──────▼──────────┐
│   gh-runner     │
│   Extension     │
└──────┬──────────┘
       │
       │ 1. Get Registration Token
       ├─────────────────────────────┐
       │                             │
┌──────▼──────────┐         ┌────────▼────────┐
│  GitHub API     │         │     Docker      │
│  (REST)         │         │                 │
│                 │         │  ┌───────────┐  │
│  Returns Token  │────────▶│  │  Runner   │  │
└─────────────────┘         │  │ Container │  │
                            │  │ (Ubuntu   │  │
                            │  │  or Win)  │  │
                            │  └───────────┘  │
                            └─────────────────┘
                                    │
                                    │ Registers & Listens
                                    │
                            ┌───────▼─────────┐
                            │  GitHub Actions │
                            │   Job Queue     │
                            └─────────────────┘
```

## Managing Runners

### View Runner Logs

```bash
docker logs -f gh-runner-<runner-name>
```

### Stop a Runner

```bash
docker stop gh-runner-<runner-name>
```

### Remove a Runner

```bash
# Stop and remove container
docker rm -f gh-runner-<runner-name>

# Then remove from GitHub UI
# Go to Settings → Actions → Runners → Click runner → Remove
```

### List Running Runners

```bash
docker ps --filter "name=gh-runner-"
```

### Restart a Runner

```bash
docker restart gh-runner-<runner-name>
```

## Troubleshooting

### "Docker is not available"

**Problem:** Docker is not installed or not running.

**Solution:**
```bash
# Install Docker Desktop (macOS/Windows)
# or Docker Engine (Linux)
# Then start Docker and try again
docker version
```

### "failed to get registration token"

**Problem:** Authentication or permission issues.

**Solutions:**
1. Authenticate with GitHub CLI:
   ```bash
   gh auth login
   ```

2. Verify you have admin access to the repo/org/enterprise

3. Check your token has the correct scopes:
   ```bash
   gh auth status
   ```

### "Container name already exists"

**Problem:** A runner with the same name is already running.

**Solution:**
```bash
# Remove the existing container
docker rm -f gh-runner-<name>

# Or use a different name
gh runner add owner/repo --docker --name runner-02
```

### Runner shows "Offline" in GitHub

**Problem:** Runner container stopped or crashed.

**Solution:**
```bash
# Check container status
docker ps -a --filter "name=gh-runner-"

# View logs to see what happened
docker logs gh-runner-<name>

# Restart the container
docker restart gh-runner-<name>
```

### Permission Denied Errors

**Problem:** Need sudo for Docker commands.

**Solution:**
```bash
# Add your user to docker group (Linux)
sudo usermod -aG docker $USER
# Then log out and back in
```

## Contributing

We welcome contributions! Here's how you can help:

### Reporting Bugs

1. Check if the issue already exists in [Issues](https://github.com/madkoo/gh-runner/issues)
2. If not, create a new issue with:
   - Clear description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Your environment (OS, Go version, Docker version)

### Suggesting Features

1. Open an issue with the `enhancement` label
2. Describe the feature and use case
3. Explain how it would benefit users

### Submitting Pull Requests

1. **Fork the repository**
   ```bash
   gh repo fork madkoo/gh-runner --clone
   cd gh-runner
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/my-awesome-feature
   ```

3. **Make your changes**
   - Follow Go best practices
   - Keep code simple and readable
   - Add comments for complex logic

4. **Test your changes**
   ```bash
   make build
   ./gh-runner add --help
   ```

5. **Commit with clear messages**
   ```bash
   git commit -m "feat: add support for runner groups"
   ```

6. **Push and create PR**
   ```bash
   git push origin feature/my-awesome-feature
   gh pr create
   ```

### Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Keep functions focused and small
- Add error handling for all operations
- Write descriptive commit messages

### Development Setup

```bash
# Clone the repo
git clone https://github.com/madkoo/gh-runner.git
cd gh-runner

# Install dependencies
go mod download

# Build
make build

# Run locally (without installing as extension)
./gh-runner add owner/repo --docker
```

### Testing Your Changes

```bash
# Build the extension
make build

# Install locally
gh extension install .

# Test with a real repository (use a test repo!)
gh runner add your-test-org/test-repo --docker --name test-runner

# Verify in GitHub UI
# Settings → Actions → Runners

# Clean up
docker rm -f gh-runner-test-runner
```

## Development

### Project Structure

```
gh-runner/
├── main.go              # Main application code
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
├── Makefile             # Build automation
├── runners/             # Runner Docker configurations
│   ├── linux/           # Linux runner files
│   │   ├── Dockerfile   # Ubuntu-based runner image
│   │   └── entrypoint.sh
│   └── windows/         # Windows runner files
│       ├── Dockerfile   # Windows Server Core runner image
│       └── entrypoint.ps1
├── README.md            # This file
├── QUICKSTART.md        # Quick start guide
├── LICENSE              # MIT license
├── examples.sh          # Usage examples
└── .gitignore           # Git ignore rules
```

### Building

```bash
# Using make
make build

# Or manually
CGO_ENABLED=0 go build -trimpath -o gh-runner
```

### Installing Locally

```bash
# Install the extension
gh extension install .

# Or upgrade if already installed
gh extension upgrade runner
```

### Uninstalling

```bash
gh extension remove runner
```

## Roadmap

### Current Features (v1.0)
- ✅ Add runners to repositories
- ✅ Add runners to organizations
- ✅ Add runners to enterprises
- ✅ Docker deployment (Linux)
- ✅ Docker deployment (Windows)
- ✅ Custom runner names
- ✅ Custom labels
- ✅ Manual setup instructions
- ✅ Official GitHub Actions runner binaries

### Planned Features

#### v1.1 - Kubernetes Support
- ⏳ Deploy runners to minikube
- ⏳ Deploy runners to existing k8s clusters
- ⏳ Helm chart support
- ⏳ Auto-scaling configuration

#### v1.2 - Runner Management
- ⏳ List existing runners (`gh runner list`)
- ⏳ Remove runners (`gh runner remove`)
- ⏳ Update runner configuration
- ⏳ Show runner status and health

#### v1.3 - Advanced Features
- ⏳ Runner groups support
- ⏳ Multiple runners deployment
- ⏳ Configuration file support
- ⏳ Runner templates
- ⏳ Automatic updates
- ⏳ Runner monitoring and alerts

### Contributing Ideas
Have an idea? [Open an issue](https://github.com/madkoo/gh-runner/issues/new) or submit a PR!

## FAQ

### Q: Do I need Docker?
**A:** No, but it's recommended. Without Docker, the tool provides manual setup instructions with the registration token.

### Q: Can I run multiple runners?
**A:** Yes! Use different `--name` values for each runner:
```bash
gh runner add owner/repo --docker --name runner-01
gh runner add owner/repo --docker --name runner-02
```

### Q: How do I update a runner?
**A:** Currently, remove the old container and create a new one:
```bash
docker rm -f gh-runner-old-name
gh runner add owner/repo --docker --name new-name
```

### Q: Can I use this in production?
**A:** Yes! The Docker runner uses official GitHub Actions runner binaries with minimal Ubuntu 24.04 (Linux) or Windows Server Core LTSC2022 base images. For production, consider:
- Using specific labels for different environments
- Setting up monitoring
- Using `--restart unless-stopped` (automatic with this tool)

### Q: What about runner security?
**A:** Runners use ephemeral registration tokens from GitHub. The Docker container runs with minimal privileges. Always review the [GitHub security best practices](https://docs.github.com/en/actions/hosting-your-own-runners/managing-self-hosted-runners/about-self-hosted-runners#self-hosted-runner-security).

### Q: Does this work with GitHub Enterprise Server?
**A:** Currently supports GitHub.com and GitHub Enterprise Cloud. GHES support is planned.

## Related Projects

- [actions/runner](https://github.com/actions/runner) - Official GitHub Actions runner (used by this extension)
- [actions-runner-controller](https://github.com/actions/actions-runner-controller) - Kubernetes controller for runners
- [github/gh](https://github.com/cli/cli) - GitHub CLI

## Support

- 📖 [Quick Start Guide](QUICKSTART.md)
- 💬 [GitHub Discussions](https://github.com/madkoo/gh-runner/discussions)
- 🐛 [Report Issues](https://github.com/madkoo/gh-runner/issues)
- 📧 Contact the maintainer

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### MIT License Summary

```
Copyright (c) 2026 Madis Koosaar

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## Acknowledgments

- GitHub CLI team for the excellent `gh` tool and [go-gh](https://github.com/cli/go-gh) library
- GitHub Actions team for the [official runner](https://github.com/actions/runner) binaries
- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- The Go and GitHub Actions communities

## Star History

If you find this project useful, please consider giving it a ⭐!

---

**Made with ❤️ by [Madis Koosaar](https://github.com/madkoo)**
