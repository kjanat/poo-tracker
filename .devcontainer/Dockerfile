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

# Create group and user (handle existing conflicts gracefully)
RUN if ! getent group $USERNAME; then \
    groupadd --gid $USER_GID $USERNAME 2>/dev/null || groupadd $USERNAME; \
    fi \
    && if ! id -u $USERNAME > /dev/null 2>&1; then \
    useradd --uid $USER_UID --gid $USERNAME -m -s /bin/bash $USERNAME 2>/dev/null || \
    useradd -m -s /bin/bash -g $USERNAME $USERNAME; \
    fi \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# Install pnpm
RUN wget -qO- https://get.pnpm.io/install.sh | ENV="/home/$USERNAME/.bashrc" SHELL="$(which bash)" bash -

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

# Install Python UV for the user
RUN curl -LsSf https://astral.sh/uv/install.sh | sh
ENV PATH="/home/$USERNAME/.local/bin:$PATH"

# Set up user environment
USER $USERNAME
WORKDIR /home/$USERNAME

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
