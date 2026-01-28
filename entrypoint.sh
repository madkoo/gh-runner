#!/bin/bash
set -e

# Required environment variables:
# - REPO_URL: GitHub repository/org/enterprise URL
# - ACCESS_TOKEN: GitHub registration token
# - RUNNER_NAME: Name for the runner (optional, defaults to hostname)
# - LABELS: Comma-separated list of labels (optional)

if [ -z "$REPO_URL" ]; then
    echo "Error: REPO_URL environment variable is required"
    exit 1
fi

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Error: ACCESS_TOKEN environment variable is required"
    exit 1
fi

RUNNER_NAME=${RUNNER_NAME:-$(hostname)}
RUNNER_WORKDIR=${RUNNER_WORKDIR:-_work}

echo "Configuring GitHub Actions Runner..."
echo "  Repository: $REPO_URL"
echo "  Runner name: $RUNNER_NAME"

# Build configuration command
CONFIG_CMD="./config.sh --url $REPO_URL --token $ACCESS_TOKEN --name $RUNNER_NAME --work $RUNNER_WORKDIR --unattended --replace"

if [ -n "$LABELS" ]; then
    CONFIG_CMD="$CONFIG_CMD --labels $LABELS"
    echo "  Labels: $LABELS"
fi

# Configure the runner
eval $CONFIG_CMD

# Cleanup function
cleanup() {
    echo "Removing runner..."
    ./config.sh remove --token $ACCESS_TOKEN
}

trap 'cleanup; exit 130' INT
trap 'cleanup; exit 143' TERM

# Start the runner
echo "Starting runner..."
./run.sh & wait $!
