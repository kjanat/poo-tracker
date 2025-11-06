# Poo Tracker Backend

> Node.js REST API for bowel movement tracking with Express v5, TypeScript, and Prisma ORM

The backend provides a robust, type-safe REST API for managing users, bowel movement entries, meals, file uploads, and analytics. Built with modern Node.js practices and comprehensive testing.

## 🚀 Features

- **RESTful API Design**: Clean, predictable endpoints following REST conventions
- **Authentication & Authorization**: JWT-based auth with secure user management
- **Database Management**: PostgreSQL with Prisma ORM for type-safe queries
- **File Upload Handling**: S3-compatible storage for photos with streaming uploads
- **Data Validation**: Comprehensive input validation and sanitization
- **Error Handling**: Structured error responses with proper HTTP status codes
- **API Documentation**: Comprehensive endpoint documentation
- **Health Monitoring**: Health check endpoints for monitoring and deployment
- **Analytics Integration**: Seamless connection to AI service for pattern analysis

## 🛠 Tech Stack

- **Runtime**: Node.js 22+ with TypeScript 5.9
- **Framework**: Express v5 with modern middleware
- **Database**: PostgreSQL with Prisma ORM v6.10
- **Authentication**: JWT tokens with bcrypt password hashing
- **File Storage**: S3-compatible (MinIO for development, AWS S3 for production)
- **Image Processing**: Sharp for image optimization
- **Validation**: Zod for runtime type checking and validation
- **Testing**: Vitest + Supertest for comprehensive API testing
- **Development**: tsx for TypeScript execution with hot reload
- **Code Quality**: ESLint + Prettier for consistent code style

## 📋 Prerequisites

- Node.js 22+
- pnpm 9+ (installed at workspace root)
- PostgreSQL database (or Docker for development)
- Redis (optional, for caching and sessions)
- S3-compatible storage service

## 🔧 Installation & Setup

### Using Workspace Commands (Recommended)

```bash
# From the root directory
pnpm install

# Set up environment variables
cp .env.example .env
# Edit .env with your database and service configurations

# Run database migrations
pnpm --filter @poo-tracker/backend db:migrate

# Start development server
pnpm --filter @poo-tracker/backend dev

# Or use the workspace shortcut
pnpm dev:backend
```

> **Security Note**: The default usernames and passwords in
> [`docker-compose.yml`](../docker-compose.yml) are for **local development
> only**. Generate unique, strong credentials for PostgreSQL, MinIO and Redis
> before deploying to production.

### Manual Setup (if needed)

```bash
# Navigate to backend directory
cd backend

# Install dependencies
pnpm install

# Set up environment
cp ../.env.example .env

# Run migrations
pnpm db:migrate

# Start development server
pnpm dev
```

## 🏃‍♂️ Development

### Development Server

```bash
# Using workspace commands (recommended)
pnpm dev:backend

# Or directly
pnpm --filter @poo-tracker/backend dev

# Manual approach
cd backend && pnpm dev
```

The API will be available at: <http://localhost:3002>

### Database Operations

```bash
# Run migrations
pnpm --filter @poo-tracker/backend db:migrate

# Reset database (⚠️ destructive!)
pnpm --filter @poo-tracker/backend db:reset

# Seed database with test data
pnpm --filter @poo-tracker/backend db:seed

# Open Prisma Studio (database GUI)
pnpm --filter @poo-tracker/backend db:studio

# Generate Prisma client
pnpm --filter @poo-tracker/backend db:generate
```

### Build & Production

```bash
# Build for production
pnpm --filter @poo-tracker/backend build

# Start production server
pnpm --filter @poo-tracker/backend start
```

## 🧪 Testing

```bash
# Run all tests
pnpm --filter @poo-tracker/backend test

# Watch mode for development
pnpm --filter @poo-tracker/backend test:watch

# Coverage report
pnpm --filter @poo-tracker/backend test:coverage
```

### Test Coverage Goals

- **Unit Tests**: All utilities, services, and middleware (100% coverage)
- **Integration Tests**: All API endpoints with various scenarios
- **Database Tests**: Prisma models and queries
- **Authentication Tests**: JWT handling and user management

## 🎨 Code Quality

```bash
# Lint code
pnpm --filter @poo-tracker/backend lint

# Fix linting issues
pnpm --filter @poo-tracker/backend lint:fix

# Format code (handled by Prettier in workspace)
pnpm format
```

## 📁 Project Structure

```text
backend/
├── prisma/
│   ├── schema.prisma           # Database schema definition
│   ├── migrations/             # Database migration files
│   └── seed.sql/              # Database seeding scripts
├── public/                     # Static assets (logos, etc.)
├── src/
│   ├── middleware/            # Express middleware
│   │   ├── auth.ts           # JWT authentication
│   │   └── errorHandler.ts   # Global error handling
│   ├── routes/               # API route handlers
│   │   ├── auth.ts          # Authentication endpoints
│   │   ├── entries.ts       # Bowel movement CRUD
│   │   ├── meals.ts         # Meal tracking
│   │   ├── uploads.ts       # File upload handling
│   │   └── analytics.ts     # Data analysis endpoints
│   ├── services/            # Business logic services
│   │   └── ImageProcessingService.ts
│   ├── utils/              # Utility functions
│   │   ├── filename.ts     # File naming utilities
│   │   └── seed.ts         # Database seeding
│   ├── config.ts          # Application configuration
│   └── index.ts           # Application entry point
├── package.json           # Dependencies and scripts
├── tsconfig.json         # TypeScript configuration
├── jest.config.js       # Test configuration (legacy)
└── vitest.config.ts     # Vitest configuration
```

## 🌐 API Endpoints

### Authentication

- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `GET /auth/me` - Get current user profile

### Bowel Movement Entries

- `GET /entries` - List user's entries (with pagination)
- `POST /entries` - Create new entry
- `GET /entries/:id` - Get specific entry
- `PUT /entries/:id` - Update entry
- `DELETE /entries/:id` - Delete entry

### Meals

- `GET /meals` - List user's meals
- `POST /meals` - Create new meal
- `GET /meals/:id` - Get specific meal
- `PUT /meals/:id` - Update meal
- `DELETE /meals/:id` - Delete meal

### File Uploads

- `POST /uploads/photos` - Upload entry photos
- `GET /uploads/photos/:filename` - Get photo (signed URL)
- `DELETE /uploads/photos/:filename` - Delete photo

### Analytics

- `GET /analytics/overview` - User's health overview
- `GET /analytics/patterns` - Pattern analysis
- `GET /analytics/correlations` - Meal-entry correlations
- `POST /analytics/analyze` - Request AI analysis

### Health & Monitoring

- `GET /health` - Service health check
- `GET /health/db` - Database connectivity check

## 🗄️ Database Schema

### Core Tables

- **User**: User accounts and profiles
- **UserAuth**: Authentication credentials (separate for security)
- **Entry**: Bowel movement records
- **Meal**: Meal tracking data
- **MealEntry**: Many-to-many relationship between meals and entries
- **Photo**: File metadata for uploaded images

### Key Relationships

- Users → Entries (one-to-many)
- Users → Meals (one-to-many)
- Meals ↔ Entries (many-to-many via MealEntry)
- Entries → Photos (one-to-many)

## ⚙️ Configuration

### Environment Variables

Copy `backend/.env.example` to `.env` and adjust the values:

```env
# Database
DATABASE_URL="postgresql://poo_user:secure_password_123@localhost:5432/poo_tracker"

# JWT Authentication
JWT_SECRET="your-super-secret-jwt-key"
JWT_EXPIRES_IN="7d"

# S3/MinIO Storage
S3_ENDPOINT="http://localhost:9000"
S3_BUCKET="poo-photos"
S3_ACCESS_KEY="minioadmin"
S3_SECRET_KEY="minioadmin123"
S3_REGION="us-east-1"

# AI Service Integration
AI_SERVICE_URL="http://localhost:8001"

# Server Configuration
PORT=3002
NODE_ENV="development"

# Optional: Redis for caching
REDIS_URL="redis://localhost:6379"
```

### Prisma Configuration

The database schema is defined in `prisma/schema.prisma` with:

- Type-safe database queries
- Automatic migration generation
- Rich relationship modeling
- Built-in validation

## 🔒 Security Features

### Authentication & Authorization

- JWT-based stateless authentication
- Bcrypt password hashing with salt rounds
- Protected routes with middleware
- User-scoped data access

### Data Protection

- Input validation with Zod schemas
- SQL injection prevention (Prisma ORM)
- XSS protection with sanitization
- CORS configuration for cross-origin requests

### File Upload Security

- File type validation
- Size limits and streaming uploads
- Secure filename generation
- S3 signed URLs for controlled access

## 🚀 Performance Features

### Database Optimization

- Connection pooling with Prisma
- Efficient queries with proper indexing
- Pagination for large datasets
- Selective field loading

### Caching Strategy

- Redis integration for session caching
- API response caching for analytics
- S3 signed URL caching
- Database query result caching

## 🧪 Testing Strategy

### Unit Tests

- All utility functions (`utils/`)
- Service layer logic (`services/`)
- Middleware functions (`middleware/`)

### Integration Tests

- Complete API endpoint testing
- Database operations
- Authentication flows
- File upload scenarios

### Test Coverage

Current test coverage: **17.79%** (actively improving)

- ✅ **middleware/auth.ts**: 100% coverage
- ✅ **middleware/errorHandler.ts**: 100% coverage
- ✅ **services/ImageProcessingService.ts**: 100% coverage
- ✅ **utils/filename.ts**: 100% coverage
- ✅ **routes/auth.ts**: 53.65% coverage (15 integration tests)

## 🔗 Service Integration

### AI Service Communication

```typescript
// Example AI service integration
const analysisResult = await fetch(`${AI_SERVICE_URL}/analyze`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ entries, meals })
})
```

### S3/MinIO Integration

```typescript
// Example file upload
const uploadResult = await s3Client.upload({
  Bucket: process.env.S3_BUCKET,
  Key: generateSecureFilename(file),
  Body: fileStream,
  ContentType: file.mimetype
})
```

## 🚨 Health Considerations

This API handles sensitive health data:

- HIPAA-consideration ready architecture
- Comprehensive audit logging
- Data encryption at rest and in transit
- User data export and deletion capabilities
- Privacy-first design principles

## 🤝 Contributing

1. Follow TypeScript and Node.js best practices
2. Write comprehensive tests for all new features
3. Use Prisma for all database operations
4. Implement proper error handling
5. Add API documentation for new endpoints
6. Follow the established project structure

### API Development Guidelines

- Use proper HTTP status codes
- Implement comprehensive validation
- Add proper error responses
- Include request/response examples
- Follow RESTful conventions
- Use TypeScript interfaces for all data models

## 🔗 Integration Points

This backend integrates with:

- **Frontend**: `../frontend/` - Serves the React application API needs
- **AI Service**: `../ai-service/` - Forwards analysis requests and processes results
- **Database**: PostgreSQL for persistent data storage
- **Storage**: S3/MinIO for file storage and management
- **Cache**: Redis for performance optimization

---

_Built with 💩 and professional engineering standards for reliable health data management._
