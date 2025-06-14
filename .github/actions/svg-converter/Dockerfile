FROM node:22-alpine

# Install system dependencies for image conversion
RUN apk add --no-cache \
    librsvg \
    imagemagick \
    bash \
    jq

# Install Node.js dependencies globally
RUN npm install -g @svgr/cli

# Copy the entrypoint script
COPY entrypoint.sh /entrypoint.sh

# Make it executable
RUN chmod +x /entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/entrypoint.sh"]
