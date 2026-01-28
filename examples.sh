#!/usr/bin/env bash
set -e

# Example 1: Add a Docker runner to a repository
echo "Example 1: Adding Docker runner to a repository"
echo "Command: gh runner add repo/owner/myrepo --docker --name my-runner --labels linux,x64"
echo ""

# Example 2: Add a Docker runner to an organization
echo "Example 2: Adding Docker runner to an organization"
echo "Command: gh runner add org/myorg --docker"
echo ""

# Example 3: Show local runner instructions
echo "Example 3: Getting manual setup instructions"
echo "Command: gh runner add repo/owner/myrepo"
echo ""

echo "Note: These are example commands. Run them without 'echo' to actually create runners."
echo "Make sure you have Docker installed for --docker flag."
