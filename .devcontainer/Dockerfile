# Use Ubuntu 24.04 LTS as base image
FROM ubuntu:24.04

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Update package lists and install essential packages
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    build-essential \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    gnupg \
    gnupg2 \
    pinentry-curses \
    lsb-release \
    sudo \
    zip \
    unzip \
    jq \
    vim \
    nano \
    htop \
    tree \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
ARG USERNAME=kjanat
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# Install pnpm
RUN curl -fsSL https://get.pnpm.io/install.sh | env PNPM_VERSION=9.15.0 sh -
ENV PNPM_HOME="/root/.local/share/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

# Install Python UV
RUN curl -LsSf https://astral.sh/uv/install.sh | sh
ENV PATH="/root/.local/bin:$PATH"

# Install Docker CLI (will be used with Docker-in-Docker feature)
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get update \
    && apt-get install -y docker-ce-cli docker-compose-plugin \
    && rm -rf /var/lib/apt/lists/*

# Install Docker Compose standalone
RUN curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose \
    && chmod +x /usr/local/bin/docker-compose

# Install PostgreSQL client
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
    && echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" | tee /etc/apt/sources.list.d/pgdg.list \
    && apt-get update \
    && apt-get install -y postgresql-client-17 \
    && rm -rf /var/lib/apt/lists/*

# Set up user environment
USER $USERNAME
WORKDIR /home/$USERNAME

# Install pnpm for the user
RUN curl -fsSL https://get.pnpm.io/install.sh | env PNPM_VERSION=9.15.0 sh -
ENV PNPM_HOME="/home/$USERNAME/.local/share/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

# Install Python UV for the user
RUN curl -LsSf https://astral.sh/uv/install.sh | sh
ENV PATH="/home/$USERNAME/.local/bin:$PATH"

# Set up shell environment with GPG and UV configuration
RUN echo 'export PATH="$HOME/.local/share/pnpm:$HOME/.local/bin:$PATH"' >> ~/.bashrc \
    && echo 'alias ll="ls -alF"' >> ~/.bashrc \
    && echo 'alias la="ls -A"' >> ~/.bashrc \
    && echo 'alias l="ls -CF"' >> ~/.bashrc \
    && echo 'export EDITOR=code' >> ~/.bashrc \
    && echo 'export GPG_TTY=$(tty)' >> ~/.bashrc \
    && echo 'export UV_CACHE_DIR="$HOME/.cache/uv"' >> ~/.bashrc \
    && echo 'export VIRTUAL_ENV="/workspaces/poo-tracker/.venv"' >> ~/.bashrc \
    && echo 'if [ -d "$VIRTUAL_ENV" ]; then export PATH="$VIRTUAL_ENV/bin:$PATH"; fi' >> ~/.bashrc

# Configure GPG for proper signing
RUN mkdir -p ~/.gnupg \
    && echo "use-agent" >> ~/.gnupg/gpg.conf \
    && echo "pinentry-mode loopback" >> ~/.gnupg/gpg.conf \
    && chmod 700 ~/.gnupg

# Add user to docker group (will be created by Docker-in-Docker feature)
USER root
RUN usermod -aG docker $USERNAME || true

# Switch back to non-root user
USER $USERNAME

# Set working directory
WORKDIR /workspaces/poo-tracker

# Keep container running
CMD ["sleep", "infinity"]
