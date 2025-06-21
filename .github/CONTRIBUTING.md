# Contributing to Poo Tracker

Welcome to the Poo Tracker project! We're excited that you want to contribute. 💩

This guide will help you understand our development process and how to submit contributions that align with our project standards.

## 🚀 Quick Start

1. **Fork** the repository
2. **Clone** your fork locally
3. **Set up** the development environment
4. **Create** a feature branch
5. **Make** your changes
6. **Test** everything thoroughly
7. **Submit** a pull request

## 🛠️ Development Setup

### Prerequisites

- **Node.js** (v22 or higher)
- **pnpm** (v8 or higher) - We don't use npm, get with the program
- **Docker** and **Docker Compose**
- **Git** (obviously)
- **Python** (v3.11 or higher) for AI service

### Initial Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/poo-tracker.git
cd poo-tracker

# Install dependencies
pnpm --filter @poo-tracker/frontend install

# Copy environment file
cp .env.example .env

# Start development services
docker-compose up -d

# Backend written in Go - no migrations via pnpm

# Start all services
make dev
```

### Project Structure

```text
poo-tracker/
├── 📱 frontend/          # React + Vite + TypeScript
├── 🔧 backend/           # Node.js + Express + TypeScript
├── 🤖 ai-service/        # Python + FastAPI
├── 🐳 docker-compose.yml # Development environment
├── 📋 package.json       # Root package.json
└── 🏗️ pnpm-workspace.yaml # pnpm workspace config
```

## 📋 Coding Standards

### The Golden Rules

> **If your code stinks, it won't get merged. And yes, we'll know.**

1. **TypeScript is mandatory** - No plain JavaScript allowed
2. **Component-based architecture** - Keep files small and focused
3. **Tailwind CSS only** - No custom CSS files, utility-first or GTFO
4. **RESTful APIs** - Don't mention GraphQL, we're not that kind of project
5. **Database migrations only through Prisma** - No raw SQL scripts
6. **Stream file uploads** - No base64 blobs in the database
7. **Test everything** - No feature without tests
8. **Environment variables** - No hardcoded secrets or configuration

### Code Style

We use ESLint with the @typescript-eslint plugin and Prettier for formatting. Don't fight the config, just follow it.

#### Frontend (React + TypeScript)

```typescript
// ✅ Good - Functional component with proper TypeScript
interface ToiletProps {
  isOccupied: boolean
  onFlush: () => Promise<void>
}

export const Toilet: React.FC<ToiletProps> = ({ isOccupied, onFlush }) => {
  const [isFlushingData, setIsFlushingData] = useState(false)

  const handleFlush = async () => {
    setIsFlushingData(true)
    try {
      await onFlush()
    } finally {
      setIsFlushingData(false)
    }
  }

  return (
    <div className="flex items-center justify-center p-4 bg-white rounded-lg shadow-md">
      <button
        onClick={handleFlush}
        disabled={isOccupied || isFlushingData}
        className="bg-blue-500 hover:bg-blue-700 disabled:bg-gray-300 text-white font-bold py-2 px-4 rounded"
      >
        {isFlushingData ? 'Flushing...' : 'Flush Data'}
      </button>
    </div>
  )
}
```

#### Backend (Node.js + Express + TypeScript)

```typescript
// ✅ Good - Proper error handling and validation
import { z } from 'zod'
import { Request, Response, NextFunction } from 'express'

const createPooSchema = z.object({
  bristolScale: z.number().min(1).max(7),
  imageUrl: z.string().url().optional(),
  notes: z.string().max(500).optional()
})

export const createPoo = async (req: Request, res: Response, next: NextFunction) => {
  try {
    const data = createPooSchema.parse(req.body)

    const poo = await prisma.poo.create({
      data: {
        ...data,
        userId: req.user.id
      }
    })

    res.status(201).json(poo)
  } catch (error) {
    next(error)
  }
}
```

### Testing Standards

#### Frontend Tests (Vitest)

```typescript
// ✅ Good - Component testing with proper mocking
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { vi } from 'vitest'
import { Toilet } from '../Toilet'

describe('Toilet', () => {
  it('should flush data when button is clicked', async () => {
    const mockOnFlush = vi.fn().mockResolvedValue(undefined)

    render(<Toilet isOccupied={false} onFlush={mockOnFlush} />)

    const flushButton = screen.getByRole('button', { name: /flush data/i })
    fireEvent.click(flushButton)

    await waitFor(() => {
      expect(mockOnFlush).toHaveBeenCalledTimes(1)
    })
  })
})
```

#### Backend Tests (Jest + Supertest)

```typescript
// ✅ Good - API endpoint testing
import request from 'supertest'
import { app } from '../app'
import { createTestUser } from '../test-utils'

describe('POST /api/poos', () => {
  it('should create a new poo entry', async () => {
    const user = await createTestUser()
    const token = generateToken(user.id)

    const response = await request(app)
      .post('/api/poos')
      .set('Authorization', `Bearer ${token}`)
      .send({
        bristolScale: 4,
        notes: 'Perfect consistency'
      })
      .expect(201)

    expect(response.body).toMatchObject({
      bristolScale: 4,
      notes: 'Perfect consistency',
      userId: user.id
    })
  })
})
```

## 🎯 Contributing Guidelines

### Before You Start

1. **Check existing issues** - Don't duplicate work
2. **Read the codebase** - Understand the architecture
3. **Start small** - Don't try to rewrite everything in your first PR

### Branch Naming

Use descriptive branch names with prefixes:

```bash
feat/bristol-stool-chart-selector
fix/image-upload-memory-leak
chore/update-dependencies
docs/api-documentation-improvements
```

### Commit Messages

Follow conventional commit format:

```bash
feat: add Bristol stool chart interactive selector
fix: resolve memory leak in image upload process
chore: update Node.js dependencies to latest versions
docs: improve API documentation with examples
test: add unit tests for poo analytics service
```

### Pull Request Process

1. **Create a feature branch** from `master`
2. **Make your changes** following our coding standards
3. **Write/update tests** for your changes
4. **Update documentation** if needed
5. **Run the full test suite** locally
6. **Submit a pull request** using our template
7. **Respond to feedback** promptly

### What Makes a Good PR

- **Single focus** - One feature or fix per PR
- **Clear description** - Explain what and why
- **Tests included** - Prove your code works
- **Documentation updated** - Keep docs in sync
- **No breaking changes** - Unless absolutely necessary
- **Clean commit history** - Squash commits if needed

## 🧪 Testing

### Running Tests

```bash

# Run all tests
make test

# Frontend tests only
pnpm --filter @poo-tracker/frontend run test

# Backend tests only
go test -C backend ./...

# AI service tests only
cd ai-service && python -m pytest

# Watch mode for development
pnpm --filter @poo-tracker/frontend run test:watch
```

### Test Requirements

- **Unit tests** for all new functions and components
- **Integration tests** for API endpoints
- **E2E tests** for critical user flows (when applicable)
- **Minimum 80% code coverage** for new code

## 📝 Documentation

### What to Document

- **API endpoints** - OpenAPI/Swagger specs
- **Component props** - JSDoc comments
- **Complex logic** - Inline comments
- **Setup instructions** - README updates
- **Configuration options** - Environment variables

### Documentation Standards

````typescript
/**
 * Analyzes poo consistency using Bristol Stool Chart classification
 * @param imageBuffer - Raw image data from uploaded photo
 * @param userId - User ID for logging and analytics
 * @returns Promise resolving to Bristol scale rating (1-7)
 * @throws {ValidationError} When image is invalid or corrupted
 * @example
 * ```typescript
 * const rating = await analyzePooConsistency(buffer, 'user123')
 * console.log(`Bristol scale: ${rating}`) // Bristol scale: 4
 * ```
 */
export async function analyzePooConsistency(imageBuffer: Buffer, userId: string): Promise<number> {
  // Implementation here
}
````

## 🚫 What NOT to Do

### Code Quality

- ❌ No `console.log` statements in production code
- ❌ No hardcoded URLs or configuration
- ❌ No CSS files (use Tailwind only)
- ❌ No any types in TypeScript
- ❌ No direct database queries (use Prisma)
- ❌ No unhandled promise rejections
- ❌ No missing error boundaries

### Process

- ❌ No direct pushes to `master`
- ❌ No PRs without tests
- ❌ No large, unfocused PRs
- ❌ No copy-pasted code from Stack Overflow without understanding
- ❌ No "fix later" TODOs without tickets
- ❌ No breaking changes without discussion

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Clear description** of the issue
2. **Steps to reproduce** the bug
3. **Expected vs actual behavior**
4. **Environment details** (OS, browser, versions)
5. **Screenshots or videos** if applicable
6. **Console logs or error messages**

Use our [bug report template](.github/ISSUE_TEMPLATE/bug_report.yml).

## ✨ Feature Requests

When requesting features, please:

1. **Search existing issues** first
2. **Describe the problem** you're trying to solve
3. **Propose a solution** with details
4. **Consider the scope** - start small
5. **Offer to help** with implementation

Use our [feature request template](.github/ISSUE_TEMPLATE/feature_request.yml).

## 🏗️ Architecture Decisions

For significant changes to architecture or technology choices:

1. **Open a discussion** first
2. **Provide context** and reasoning
3. **Consider alternatives**
4. **Get feedback** from maintainers
5. **Document the decision**

## 📦 Dependencies

### Adding New Dependencies

Before adding new dependencies:

1. **Check if it's really needed** - Can we solve this another way?
2. **Research the package** - Is it maintained? Popular? Secure?
3. **Consider bundle size** - Especially for frontend
4. **Get approval** for major dependencies

### Preferred Libraries

- **UI Components**: Headless UI, Radix UI
- **Forms**: React Hook Form + Zod validation
- **HTTP Client**: Axios (backend), fetch (frontend)
- **Date/Time**: date-fns
- **Icons**: Lucide React
- **Testing**: Vitest, Testing Library, Jest
- **Styling**: Tailwind CSS only

## 🔒 Security

### Security Guidelines

- **Never commit secrets** - Use environment variables
- **Validate all inputs** - Both client and server side
- **Use HTTPS** in production
- **Implement proper authentication** - JWT with secure practices
- **Follow OWASP guidelines** - Especially for web vulnerabilities
- **Report security issues privately** - See SECURITY.md

## 📞 Getting Help

### Where to Get Help

- **GitHub Discussions** - General questions and ideas
- **Discord** - Real-time chat with the community
- **Issues** - Bug reports and feature requests
- **Email** - For sensitive or private matters

### Code Review Process

All PRs require:

1. **Automated checks passing** - CI/CD pipeline must be green
2. **Code review approval** - From at least one maintainer
3. **No unresolved comments** - Address all feedback
4. **Up-to-date with master** - Rebase or merge master first

## 🎉 Recognition

We appreciate all contributions! Contributors are recognized through:

- **GitHub contributors page**
- **Release notes mentions**
- **Community shout-outs**
- **Maintainer privileges** for consistent contributors

---

## Final Words

Remember: **Clean code = Clean colon** 💩

We're building something awesome together, and every contribution matters. Don't be afraid to ask questions, make mistakes, or suggest improvements.

Happy coding, you magnificent bastards! 🚀

---

_For questions about contributing, reach out to [@kjanat](https://github.com/kjanat) or start a discussion._
