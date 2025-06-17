# Update PR #50

## New Title
chore: remove devcontainer and update development tooling

## New Description
This PR removes the DevContainer configuration and updates various development tooling configurations for improved consistency and simplified setup.

### Major Changes

**DevContainer Removal:**
- Removed entire `.devcontainer/` directory and all associated files
- Eliminated Docker-in-Docker DevContainer setup
- Removed custom development environment configuration

**Development Tooling Updates:**
- Updated CI/CD workflows with improved Node.js versions (22, 24)
- Enhanced GitHub Actions with better error handling and test reporting
- Added Trunk integration for test result tracking
- Updated ESLint configuration from StandardJS to @typescript-eslint
- Added comprehensive `.dockerignore` for better build optimization

**Configuration Improvements:**
- Added standardized `.prettierrc` configuration
- Updated `.gitignore` with better patterns and exclusions
- Enhanced dependabot configuration with single quotes consistency
- Updated issue templates and GitHub workflows for better formatting

**CI/CD Enhancements:**
- Improved test setup with proper database initialization
- Added comprehensive test result reporting
- Enhanced Docker build process with SBOM and provenance
- Better UV and Python dependency management in workflows

### Breaking Changes

- DevContainer support removed - developers need to use local development setup
- Node.js version requirement updated to v22+

### Migration Guide

For developers previously using DevContainer:
1. Install Node.js 22+ locally
2. Install pnpm, Docker, and required tools
3. Run `pnpm install` and `pnpm docker:up` for local development
4. Use standard VS Code without DevContainer for development

This change simplifies the development setup by removing the complexity of DevContainer configuration while maintaining all core functionality through standard local development tools.