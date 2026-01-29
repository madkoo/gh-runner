# Quick Start Guide

## Prerequisites

1. **GitHub CLI** - Install from https://cli.github.com/
2. **Docker** (optional but recommended) - Install from https://www.docker.com/
3. **Authentication** - Run `gh auth login` to authenticate

## Installation

### Option 1: Install from this directory
```bash
gh extension install .
```

### Option 2: Build and install manually
```bash
make build
gh extension install .
```

## Usage

### 1. Add a runner to your private repository (Docker)

```bash
gh runner add yourusername/your-private-repo --docker
```

This will:
- Generate a registration token from GitHub
- Pull the official GitHub runner Docker image
- Start a container with the runner configured and running

### 2. Add a runner to your organization

```bash
gh runner add org/yourorg --docker --labels linux,x64,production
```

### 3. Add a runner to your enterprise

```bash
gh runner add enterprise/yourenterprise --docker
```
### 4. Custom runner name

```bash
gh runner add owner/repo --docker --name my-custom-runner
```

## Managing Docker Runners

After starting a runner, you can:

```bash
# View logs
docker logs -f gh-runner-<name>

# Stop the runner
docker stop gh-runner-<name>

# Remove the runner
docker rm -f gh-runner-<name>
```

## Troubleshooting

### "Docker is not available"
Install Docker Desktop or Docker Engine on your machine.

### "failed to get registration token"
Make sure you're authenticated with `gh auth login` and have admin access to the repo/org/enterprise.

### Check runner status
Go to your repository/org/enterprise Settings → Actions → Runners to see your runner status.

## What's Next?

The current version supports:
- ✅ Repository runners
- ✅ Organization runners  
- ✅ Enterprise runners
- ✅ Docker deployment
- ✅ Custom labels and names

Coming soon:
- ⏳ Kubernetes/Minikube deployment
- ⏳ Runner removal command
- ⏳ List existing runners
- ⏳ Update runner configuration

## Example Workflow

1. Install the extension:
   ```bash
   gh extension install .
   ```

2. Add a runner to your private repo:
   ```bash
   gh runner add myorg/private-repo --docker --labels linux,gpu
   ```

3. Check the runner in GitHub:
   - Go to your repository → Settings → Actions → Runners
   - You should see your new runner online!

4. Use it in a workflow:
   ```yaml
   jobs:
     build:
       runs-on: [self-hosted, linux, gpu]
       steps:
         - uses: actions/checkout@v6
         - run: echo "Running on self-hosted runner!"
   ```
