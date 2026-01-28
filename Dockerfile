# Dockerfile for GitHub Actions Self-Hosted Runner (Linux)
# Uses official GitHub Actions runner binaries with minimal Ubuntu image

FROM ubuntu:24.04 AS base

# Prevent interactive prompts during build
ENV DEBIAN_FRONTEND=noninteractive

# Install minimal dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    jq \
    git \
    sudo \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Create a user to run the runner
RUN useradd -m -s /bin/bash runner && \
    usermod -aG sudo runner && \
    echo "runner ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

WORKDIR /home/runner

# Detect architecture and download appropriate runner
RUN RUNNER_VERSION=$(curl -s https://api.github.com/repos/actions/runner/releases/latest | jq -r '.tag_name' | sed 's/v//') && \
    ARCH=$(dpkg --print-architecture) && \
    case "$ARCH" in \
        amd64) RUNNER_ARCH="x64" ;; \
        arm64) RUNNER_ARCH="arm64" ;; \
        *) echo "Unsupported architecture: $ARCH" && exit 1 ;; \
    esac && \
    echo "Downloading runner version ${RUNNER_VERSION} for ${RUNNER_ARCH}..." && \
    curl -o actions-runner.tar.gz -L \
    "https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-linux-${RUNNER_ARCH}-${RUNNER_VERSION}.tar.gz" && \
    tar xzf actions-runner.tar.gz && \
    rm actions-runner.tar.gz && \
    chown -R runner:runner /home/runner

# Install runner dependencies
RUN ./bin/installdependencies.sh && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

USER runner

ENTRYPOINT ["/entrypoint.sh"]
