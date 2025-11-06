# Poo Tracker Frontend

> Modern React application for tracking bowel movements with style and zero shame

The frontend is a responsive, type-safe React application built with modern web technologies. It provides an intuitive interface for logging bowel movements, tracking meals, viewing analytics, and managing health data.

## 🚀 Features

- **Intuitive Bowel Movement Logging**: Easy-to-use forms with Bristol Stool Chart integration
- **Photo Upload & Management**: S3-compatible storage for visual documentation
- **Meal Tracking**: Comprehensive food logging with correlation analysis
- **Analytics Dashboard**: Beautiful charts and insights powered by AI analysis
- **User Authentication**: Secure login and registration system
- **Responsive Design**: Works seamlessly on desktop, tablet, and mobile
- **Real-time Updates**: Live data synchronization with backend
- **Offline Support**: Service worker for offline functionality

## 🛠 Tech Stack

- **Framework**: React 19 with hooks and functional components
- **Build Tool**: Vite 6 with TypeScript support
- **Styling**: TailwindCSS v4 with utility-first approach
- **Routing**: React Router DOM v7 for navigation
- **State Management**: TanStack Query for server state + Zustand for client state
- **Forms**: React Hook Form with validation
- **Icons**: Lucide React for beautiful SVG icons
- **Charts**: Recharts for data visualizations
- **HTTP Client**: Axios for API communication
- **Testing**: Vitest + React Testing Library
- **Type Safety**: TypeScript 5.8 throughout

## 📋 Prerequisites

- Node.js 22+
- pnpm 9+ (installed at workspace root)
- Backend API running on `http://localhost:3002`

## 🔧 Installation & Setup

### Using Workspace Commands (Recommended)

```bash
# From the root directory
pnpm install

# Start development server
pnpm --filter @poo-tracker/frontend dev

# Or use the workspace shortcut
pnpm dev:frontend
```

### Manual Setup (if needed)

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
pnpm install

# Start development server
pnpm dev
```

## 🏃‍♂️ Development

### Development Server

```bash
# Using workspace commands (recommended)
pnpm dev:frontend

# Or directly
pnpm --filter @poo-tracker/frontend dev

# Manual approach
cd frontend && pnpm dev
```

Access the application at: <http://localhost:5173>

### Build & Preview

```bash
# Build for production
pnpm --filter @poo-tracker/frontend build

# Preview production build
pnpm --filter @poo-tracker/frontend preview
```

## 🧪 Testing

```bash
# Run all tests
pnpm --filter @poo-tracker/frontend test

# Watch mode for development
pnpm --filter @poo-tracker/frontend test:watch

# Coverage report
pnpm --filter @poo-tracker/frontend test:coverage

# Interactive test UI
pnpm --filter @poo-tracker/frontend test:ui
```

## 🎨 Code Quality

```bash
# Lint code
pnpm --filter @poo-tracker/frontend lint

# Fix linting issues
pnpm --filter @poo-tracker/frontend lint:fix

# Format code (handled by Prettier in workspace)
pnpm format
```

## 📁 Project Structure

```text
frontend/
├── public/                 # Static assets
│   ├── favicon.ico
│   ├── logo_*.png         # Generated favicons
│   └── manifest.json      # PWA manifest
├── src/
│   ├── components/        # Reusable React components
│   │   ├── Logo.tsx
│   │   ├── Navbar.tsx
│   │   └── ProtectedRoute.tsx
│   ├── pages/            # Page components
│   │   ├── HomePage.tsx
│   │   ├── LoginPage.tsx
│   │   ├── DashboardPage.tsx
│   │   ├── NewEntryPage.tsx
│   │   ├── MealsPage.tsx
│   │   ├── AnalyticsPage.tsx
│   │   └── ProfilePage.tsx
│   ├── stores/           # State management
│   │   └── authStore.ts  # Authentication state
│   ├── utils/            # Utility functions
│   │   └── branding.ts   # Brand assets utilities
│   ├── test/             # Test utilities
│   │   ├── setup.ts
│   │   └── vitest-setup.d.ts
│   ├── App.tsx           # Main application component
│   ├── main.tsx          # Application entry point
│   └── index.css         # Global styles and Tailwind imports
├── index.html            # HTML template
├── package.json          # Dependencies and scripts
├── vite.config.ts        # Vite configuration
├── vitest.config.ts      # Test configuration
├── tailwind.config.js    # TailwindCSS configuration
├── tsconfig.json         # TypeScript configuration
└── tsconfig.node.json    # Node.js TypeScript config
```

## 🔌 API Integration

The frontend communicates with:

- **Backend API**: `http://localhost:3002` - Main REST API
- **AI Service**: `http://localhost:8001` (via backend proxy) - Analysis endpoints

### Key API Endpoints Used

- `POST /auth/login` - User authentication
- `POST /auth/register` - User registration
- `GET /entries` - Fetch bowel movement entries
- `POST /entries` - Create new entry
- `GET /meals` - Fetch meal data
- `POST /meals` - Create new meal
- `GET /analytics` - Fetch analysis data
- `POST /uploads` - Photo uploads

## 🎯 Key Features

### Bristol Stool Chart Integration

- Interactive chart for stool type selection
- Visual representations with descriptions
- Health indicators and recommendations

### Photo Management

- Drag & drop photo uploads
- Image preview and cropping
- S3-compatible storage integration
- Secure file access with signed URLs

### Analytics & Insights

- Beautiful data visualizations
- AI-powered pattern recognition
- Meal correlation analysis
- Health trend tracking
- Exportable reports

### Responsive Design

- Mobile-first approach
- Tablet and desktop optimizations
- Touch-friendly interfaces
- Accessible design patterns

## ⚙️ Configuration

### Environment Variables

Copy `frontend/.env.example` to `.env.local` and adjust the values:

```env
# API Configuration
VITE_API_URL=http://localhost:3002
VITE_AI_SERVICE_URL=http://localhost:8001

# Feature Flags
VITE_ENABLE_ANALYTICS=true
VITE_ENABLE_PHOTO_UPLOAD=true

# S3/MinIO Configuration (for direct uploads if needed)
VITE_S3_BUCKET=poo-photos
VITE_S3_ENDPOINT=http://localhost:9000
```

### TailwindCSS Configuration

The frontend uses TailwindCSS v4 with:

- Custom color scheme for the poo theme
- Responsive breakpoints
- Custom components and utilities
- Dark mode support (future enhancement)

## 🚨 Health & Safety

This application handles sensitive health data:

- All API calls use HTTPS in production
- User authentication with JWT tokens
- Secure file uploads with validation
- Privacy-focused design (no analytics tracking)
- HIPAA-consideration ready architecture

## 🤝 Contributing

1. Follow React and TypeScript best practices
2. Use functional components with hooks
3. Write tests for new components and features
4. Follow the established file structure
5. Use workspace commands for consistency
6. Ensure responsive design for all new features

### Component Guidelines

- Use TypeScript interfaces for all props
- Implement proper error boundaries
- Use React Query for server state management
- Follow the existing naming conventions
- Add proper accessibility attributes

## 🔗 Integration

This frontend integrates with:

- **Backend API**: `../backend/` - Authentication, data storage, business logic
- **AI Service**: `../ai-service/` - Pattern analysis and recommendations
- **Docker Services**: PostgreSQL, Redis, MinIO for complete functionality

---

_Built with ❤️ and a healthy sense of humor for tracking life's most honest moments._
