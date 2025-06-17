# DevContainer Configuration for Poo Tracker

This directory contains the DevContainer configuration for the Poo Tracker project, providing a complete development environment with all necessary tools and services.

## Features

- **Docker-in-Docker**: Full Docker support inside the container
- **Multi-language support**: Node.js, Python, TypeScript
- **Package managers**: pnpm, uv (Python)
- **Database**: PostgreSQL with automatic migrations
- **Storage**: MinIO (S3-compatible)
- **Caching**: Redis
- **AI Service**: Python FastAPI with ML capabilities
- **Code quality**: ESLint, Prettier, Husky pre-commit hooks
- **VS Code extensions**: Pre-configured for optimal development experience
- **GPG Commit Signing**: Seamless integration with host GPG keys
- **Python Virtual Environments**: Isolated UV-managed environments

## Port Configuration

The DevContainer forwards the following ports:

- **3000**: Frontend (React/Vite)
- **3001**: Backend (Node.js/Express)
- **5432**: PostgreSQL Database
- **6379**: Redis Cache
- **8001**: AI Service (Python FastAPI)
- **9000**: MinIO S3 API
- **9002**: MinIO Web Console

## Environment Variables

All services are pre-configured with development-friendly environment variables:

- Database connections
- S3 storage configuration
- Redis cache settings
- AI service endpoints
- CORS settings for local development

## Automatic Setup

The DevContainer automatically:

1. **Post-Create**: Installs all dependencies, sets up environment files
2. **Post-Start**: Starts infrastructure services, runs database migrations
3. **Service Health Checks**: Waits for all services to be ready

## Usage

1. Open the project in VS Code
2. When prompted, click "Reopen in Container"
3. Wait for the setup to complete
4. Run `pnpm dev` to start development servers

For detailed setup guides:

- **GPG Signing**: See [GPG_SETUP.md](./GPG_SETUP.md)
- **Python UV**: See [UV_SETUP.md](./UV_SETUP.md)

## Manual Commands

If you need to manage services manually:

```bash
# Start infrastructure services
pnpm docker:up

# Stop infrastructure services
pnpm docker:down

# Run database migrations
pnpm db:migrate

# Seed the database
pnpm db:seed

# Start development servers
pnpm dev
```

## Troubleshooting

### Port Conflicts

If you encounter port conflicts, the DevContainer uses isolated networking to prevent conflicts with your host system. Each service runs in its own container with mapped ports.

### Docker Issues

If Docker-in-Docker isn't working:

1. Ensure Docker Desktop is running on your host
2. Check that the DevContainer has `--privileged` access
3. Verify the Docker socket mount is working

### Database Issues

If database connections fail:

1. Check that PostgreSQL container is running: `docker ps`
2. Verify database credentials in environment files
3. Run migrations manually: `pnpm db:migrate`

### Performance

For better performance:

1. Use WSL2 backend on Windows
2. Allocate sufficient memory to Docker Desktop (8GB+ recommended)
3. Consider using volume mounts for node_modules

## Clean Code Reminders

Remember the project's coding standards:

- TypeScript everywhere
- ESLint + Prettier for formatting
- Component-based architecture
- Test-driven development
- Proper commit messages
- No console.log in production code

Now go code, you magnificent bastard! 💩
