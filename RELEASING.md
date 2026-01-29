# Release Process

This document explains how to create a new release of gh-runner with precompiled binaries.

## Overview

The `gh-runner` extension uses GitHub Actions to automatically build and publish precompiled binaries for multiple platforms when a new release is created.

## Creating a Release

### 1. Decide on the Version Number

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality in a backwards compatible manner
- **PATCH** version for backwards compatible bug fixes

Example: `v1.0.0`, `v1.1.0`, `v1.1.1`

### 2. Create and Push a Version Tag

```bash
# Make sure you're on the main branch with latest changes
git checkout main
git pull

# Create a new tag
git tag v1.0.0

# Push the tag to GitHub
git push origin v1.0.0
```

### 3. Automated Build Process

Once the tag is pushed, the GitHub Actions workflow (`.github/workflows/release.yml`) will automatically:
1. Checkout the code
2. Build binaries for multiple platforms:
   - Linux (x64, arm64)
   - macOS (x64, arm64)
   - Windows (x64, arm64)
3. Create a GitHub release
4. Attach all precompiled binaries to the release

### 4. Edit the Release Notes (Optional)

After the workflow completes:
1. Go to https://github.com/madkoo/gh-runner/releases
2. Find the newly created release
3. Click "Edit release"
4. Add release notes describing what's new, what's changed, and any bug fixes
5. Save the release

## Manual Release (Alternative)

You can also trigger the release workflow manually:

1. Go to https://github.com/madkoo/gh-runner/actions/workflows/release.yml
2. Click "Run workflow"
3. Select the branch
4. Click "Run workflow"

This is useful for testing or creating a release without a tag.

## Installation by Users

Once the release is published with binaries, users can install the extension with:

```bash
gh extension install madkoo/gh-runner
```

The `gh` CLI will automatically:
1. Detect the latest release
2. Download the appropriate precompiled binary for their platform
3. Install it without requiring Go or any build tools

## Verifying a Release

After creating a release, verify it works correctly:

```bash
# Remove any existing installation
gh extension remove runner

# Install from the new release
gh extension install madkoo/gh-runner

# Verify the version and functionality
gh runner --help
```

## Troubleshooting

### Workflow Fails

If the release workflow fails:
1. Check the workflow logs in the Actions tab
2. Common issues:
   - Invalid `go.mod` file
   - Insufficient permissions (ensure `contents: write` is set)
   - Build errors in the Go code

### Binaries Not Attached

If binaries are not attached to the release:
1. Ensure the workflow completed successfully
2. Check that the `cli/gh-extension-precompile` action ran
3. Verify the release was created (not just a draft)

### Installation Still Fails

If users still get "missing executable" errors after a release:
1. Ensure the release has binaries attached (check the release page)
2. Users may need to remove cached extension: `gh extension remove runner`
3. Try installing again: `gh extension install madkoo/gh-runner`
