# Support

Thanks for using Poo Tracker! 💩

This document explains how to get help when you need it.

## 🤔 Before You Ask

Please check these resources first:

- 📚 **[README](../README.md)** - Basic setup and usage information
- 📖 **[Documentation](https://github.com/kjanat/poo-tracker/wiki)** - Detailed guides and tutorials
- 🔍 **[Existing Issues](https://github.com/kjanat/poo-tracker/issues)** - Someone might have already asked
- 💬 **[Discussions](https://github.com/kjanat/poo-tracker/discussions)** - Community Q&A

## 🆘 Getting Help

### For Users

If you're using Poo Tracker and need help:

#### 1. 💬 GitHub Discussions (Recommended)

- **Best for**: General questions, usage help, feature discussions
- **Response time**: Usually within 24-48 hours
- **Audience**: Community members and maintainers
- **[Start a discussion →](https://github.com/kjanat/poo-tracker/discussions)**

#### 2. 🐛 GitHub Issues

- **Best for**: Bug reports, feature requests
- **Use templates**: We have specific templates for different types of issues
- **Be specific**: Include steps to reproduce, environment details
- **[Create an issue →](https://github.com/kjanat/poo-tracker/issues/new/choose)**

#### 3. 💬 Discord Community

- **Best for**: Real-time chat, quick questions
- **Response time**: Varies (community-driven)
- **[Join Discord →](https://discord.gg/your-server)**

### For Developers

If you're developing with or contributing to Poo Tracker:

#### 1. 📋 Contributing Guide

- Read our **[CONTRIBUTING.md](CONTRIBUTING.md)** first
- Includes setup instructions, coding standards, and PR process

#### 2. 🏗️ Development Setup Issues

- Check the **[Development Setup](CONTRIBUTING.md#-development-setup)** section
- Common issues are documented with solutions

#### 3. 🧪 Testing Problems

- See our **[Testing Guide](CONTRIBUTING.md#-testing)** for help
- Include test output and environment details when asking

#### 4. 🔧 Architecture Questions

- Use **[GitHub Discussions](https://github.com/kjanat/poo-tracker/discussions)** for design questions
- Tag discussions with appropriate labels

## 📞 Contact Methods

### Community Support (Public)

| Method                | Best For                         | Response Time |
| --------------------- | -------------------------------- | ------------- |
| 💬 GitHub Discussions | General questions, feature ideas | 24-48 hours   |
| 🐛 GitHub Issues      | Bug reports, specific problems   | 24-72 hours   |
| 💬 Discord            | Quick questions, real-time help  | Varies        |

### Direct Contact (Private)

| Method      | Best For                         | Contact                              |
| ----------- | -------------------------------- | ------------------------------------ |
| 📧 Email    | Business inquiries, partnerships | <poo-tracker@kjanat.com>             |
| 🔒 Security | Security vulnerabilities         | <poo-tracker+security@kjanat.com>    |
| 👤 Personal | Direct questions to maintainer   | [@kjanat](https://github.com/kjanat) |

## 🏷️ Issue Labels

When creating issues, these labels help us categorize and prioritize:

### Type Labels

- `bug` - Something isn't working
- `enhancement` - New feature or improvement
- `question` - General question
- `documentation` - Documentation improvement
- `good first issue` - Good for newcomers

### Priority Labels

- `critical` - System is broken
- `high` - Important issue
- `medium` - Standard priority
- `low` - Nice to have

### Component Labels

- `frontend` - React/UI related
- `backend` - Node.js/API related
- `ai-service` - Python/ML related
- `database` - PostgreSQL/Prisma related
- `docker` - Containerization issues
- `ci/cd` - Build/deployment issues

## 📋 How to Write Good Support Requests

### Include These Details

#### For Bug Reports

```markdown
**Environment:**

- OS: Windows 11 / macOS 14 / Ubuntu 22.04
- Node.js version: v22.17.0
- Browser: Chrome 119 (if frontend issue)
- Docker version: 24.0.6 (if relevant)

**Steps to Reproduce:**

1. Go to...
2. Click on...
3. See error...

**Expected Behavior:**
What should have happened

**Actual Behavior:**
What actually happened

**Screenshots:**
[Attach if helpful]

**Console Output:**
[Copy and paste any error messages]
```

#### For Feature Requests

```markdown
**Problem:**
What problem are you trying to solve?

**Proposed Solution:**
How would you like to see it solved?

**Use Case:**
Describe your specific use case

**Alternatives:**
What alternatives have you considered?
```

#### For Questions

```markdown
**What I'm trying to do:**
Clear description of your goal

**What I've tried:**
Steps you've already taken

**Specific question:**
What exactly are you stuck on?

**Context:**
Your setup, environment, etc.
```

## 🚀 Self-Help Resources

### Common Issues

#### Setup Problems

**Docker not starting:**

```bash
# Check Docker status
docker --version
docker-compose --version

# Reset Docker
docker-compose down
docker-compose up -d
```

**Dependencies not installing:**

```bash
# Clear cache and reinstall
rm -rf node_modules pnpm-lock.yaml
pnpm install
```

**Database connection issues:**

```bash
# Check database is running
docker-compose ps

# Reset database
cd backend
pnpm prisma migrate reset
```

#### Development Issues

**TypeScript errors:**

- Check your `tsconfig.json` configuration
- Make sure all dependencies are installed
- Restart your IDE's TypeScript service

**Build failures:**

- Check for TypeScript errors first
- Ensure all environment variables are set
- Try clearing build cache

**Test failures:**

- Run tests in isolation: `pnpm test -- --run`
- Check test database configuration
- Make sure test data is properly cleaned up

### Debugging Tips

1. **Check the logs:**

   ```bash
   # Docker container logs
   docker-compose logs frontend
   docker-compose logs backend
   docker-compose logs ai-service
   ```

2. **Use the debugger:**
   - Frontend: Browser dev tools
   - Backend: VS Code debugger or `console.log`
   - AI Service: Python debugger or print statements

3. **Isolate the problem:**
   - Test individual components
   - Use minimal reproduction cases
   - Check network requests in dev tools

## 📚 Learning Resources

### Getting Started

- **[React Documentation](https://react.dev/)** - Frontend framework
- **[Node.js Guides](https://nodejs.org/en/docs/guides/)** - Backend runtime
- **[FastAPI Tutorial](https://fastapi.tiangolo.com/tutorial/)** - API framework
- **[Docker Getting Started](https://docs.docker.com/get-started/)** - Containerization

### Advanced Topics

- **[TypeScript Handbook](https://www.typescriptlang.org/docs/)** - Type system
- **[Prisma Docs](https://www.prisma.io/docs/)** - Database ORM
- **[Tailwind CSS](https://tailwindcss.com/docs)** - Utility-first CSS
- **[Testing Library](https://testing-library.com/docs/)** - Testing utilities

## 🎯 Response Time Expectations

| Issue Type          | Initial Response   | Resolution       |
| ------------------- | ------------------ | ---------------- |
| **Critical Bug**    | Within 4 hours     | 24-48 hours      |
| **Bug**             | Within 24 hours    | 3-7 days         |
| **Feature Request** | Within 48 hours    | Varies           |
| **Question**        | Within 24-48 hours | Usually same day |
| **Documentation**   | Within 48 hours    | 1-2 weeks        |

_Note: These are goals, not guarantees. Response times may vary based on complexity and maintainer availability._

## 🤝 Community Guidelines

When asking for help:

- **Be respectful** - We're all here to help each other
- **Be patient** - Maintainers are volunteers with day jobs
- **Be specific** - Vague questions get vague answers
- **Search first** - Don't ask questions that have been answered
- **Give back** - Help others when you can

## 🔄 Escalation Process

If you're not getting the help you need:

1. **Wait for response time** - Give it at least 48 hours for non-critical issues
2. **Provide more details** - Add additional context to your request
3. **Try different channels** - Move from Discord to GitHub Issues for example
4. **Tag maintainers** - Use @mentions sparingly and only when necessary
5. **Contact directly** - Email for urgent or sensitive matters

---

## 💡 Pro Tips

- **Use markdown formatting** in GitHub issues for better readability
- **Include links** to relevant code or documentation
- **Be clear about your environment** - versions matter
- **Show your work** - explain what you've already tried
- **Follow up** - let us know if solutions work

---

Remember: There are no stupid questions, only stupid answers! 💩

We're here to help you succeed with Poo Tracker. Don't hesitate to reach out! 🚀
