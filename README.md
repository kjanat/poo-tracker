# 💩 Poo Tracker

Ever wondered if your gut's on a winning streak, or if your last kebab is about to take its revenge? Welcome to **Poo Tracker**, the brutally honest app that lets you log every majestic turd, rabbit pellet, and volcanic diarrhea eruption without shame or censorship. Because your bowel movements _do_ matter, and no, your GP isn't getting this level of data.

## 🚀 Quick Start

### Prerequisites

- Node.js 18+
- Docker & Docker Compose
- pnpm (because npm is for amateurs)

### Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/kjanat/poo-tracker.git
   cd poo-tracker
   ```

2. **Install dependencies:**

   ```bash
   pnpm install
   ```

3. **Set up environment variables:**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start the services:**

   ```bash
   # Start database and supporting services
   pnpm run docker:up

   # Start the development servers
   pnpm run dev
   ```

5. **Set up the database:**

   ```bash
   # Run database migrations
   pnpm run db:migrate

   # (Optional) Seed with test data
   pnpm run db:seed
   ```

6. **Open your browser:**

- Frontend: <http://localhost:5173>
- Backend API: <http://localhost:3002>
- AI Service: <http://localhost:8001>
- MinIO Console: <http://localhost:9002> (minioadmin/minioadmin123)

## 🏗️ Tech Stack

- **Frontend**: React + Vite + TypeScript + TailwindCSS
- **Backend**: Node.js + Express + TypeScript + Prisma
- **Database**: PostgreSQL
- **Storage**: MinIO (S3-compatible for photos)
- **AI Service**: Python + FastAPI + scikit-learn
- **Infrastructure**: Docker + Docker Compose

## 📖 What's the fucking point?

- **Track your shits**: Log every glorious bowel movement in a timeline worthy of its own Netflix docu-series. Date, time, volume, Bristol score, color, floaters—you name it.
- **Photo evidence**: Snap a pic for science (or just to traumatize your friends and family). Upload and catalog the full visual glory of your rectal output.
- **Meal log**: Record what you've stuffed down your throat so you can finally prove that Taco Tuesday was, in fact, a war crime against your colon.
- **AI-powered analysis**: Our cold, heartless machine learning model doesn't judge—just coldly correlates your food, stress, and fiber intake with your latest gut grenade.
- **Export & share**: PDF reports for your doctor, gastroenterologist, or weird fetish community. It's your data, so overshare as you wish.

## 🔧 Development

### Available Scripts

```bash
# Development
pnpm run dev             # Start frontend + backend
pnpm run dev:frontend    # Start frontend only
pnpm run dev:backend     # Start backend only

# Database
pnpm run db:migrate      # Run Prisma migrations
pnpm run db:seed         # Seed database

# Docker
pnpm run docker:up       # Start all services
pnpm run docker:down     # Stop all services

# Testing & Linting
pnpm run test            # Run all tests
pnpm run lint            # Run linters
pnpm run build           # Build for production
```

### Project Structure

```tree
poo-tracker/
├── frontend/          # React frontend (Vite + TypeScript + TailwindCSS)
├── backend/           # Express.js API (TypeScript + Prisma)
├── ai-service/        # Python FastAPI AI service
├── docker-compose.yml # Infrastructure setup
└── package.json       # Root workspace config
```

### API Endpoints

| Endpoint                 | Method   | Description                           |
| ------------------------ | -------- | ------------------------------------- |
| `/api/auth/register`     | `POST`   | Register a new user                   |
| `/api/auth/login`        | `POST`   | Authenticate user                     |
| `/api/entries`           | `GET`    | Get bowel movement entries            |
| `/api/entries`           | `POST`   | Create a new bowel movement entry     |
| `/api/entries/:id`       | `PUT`    | Update a bowel movement entry         |
| `/api/entries/:id`       | `DELETE` | Delete a bowel movement entry         |
| `/api/uploads/photo`     | `POST`   | Upload a photo of your bowel movement |
| `/api/analytics/summary` | `GET`    | Get AI analysis summary               |

## 📝 How it works

1. **Eat** something questionable.
2. **Shit**. (Preferably in a toilet, but we're not here to kink-shame.)
3. **Log** your experience in Poo Tracker—optionally add a photo, rate your suffering, describe the smell if you're brave enough.
4. **Repeat** until you finally realize you're lactose intolerant or have IBS.

## 🔒 Privacy, baby

We encrypt your brown notes and hide them away. Nobody's reading your logs except you—and whatever godforsaken AI wants to learn about the day-after effects of your sushi buffet.

## 👥 Who is this for?

- People with IBS, Crohn, colitis, or "nobody shits like me" syndrome.
- Control freaks and data lovers.
- Curious bastards who just want to know if spicy food really _is_ their nemesis.

## 🎯 Features

### Bristol Stool Chart Integration

- Complete 7-type Bristol Stool Chart support
- Visual indicators and descriptions
- Trend analysis over time

### AI-Powered Insights

- Pattern recognition in bowel movements
- Meal correlation analysis
- Health recommendations
- Risk factor identification

### Photo Documentation

- Secure photo upload and storage
- Image compression and optimization
- Privacy-focused MinIO storage

### Export & Analytics

- Detailed charts and visualizations
- Data export capabilities
- Comprehensive health insights

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Follow the coding standards (see copilot instructions)
4. Write tests for new features
5. Commit using conventional commits (`feat: add Bristol chart selector`)
6. Push and create a Pull Request

### Coding Standards

- TypeScript everywhere (no JS allowed)
- Component-based architecture
- TailwindCSS for styling (no CSS modules)
- RESTful API design
- Comprehensive test coverage
- ESLint + Prettier (follow the config, don't "fix" it)

## 🚀 Deployment

### Environment Variables

```env
# Database
DATABASE_URL="postgresql://poo_user:secure_password_123@localhost:5432/poo_tracker"

# Storage (MinIO)
MINIO_ENDPOINT="localhost:9000"
MINIO_ACCESS_KEY="minioadmin"
MINIO_SECRET_KEY="minioadmin123"

# API Configuration
API_PORT=3002
CORS_ORIGIN="http://localhost:5173"

# AI Service
AI_SERVICE_URL="http://localhost:8001"

# JWT
JWT_SECRET="your-super-secret-jwt-key-change-in-production"
```

### Production Build

```bash
pnpm run build
```

## 📄 License

This project is licensed under the **GNU Affero General Public License v3.0** (AGPLv3).

See the [LICENSE](LICENSE) file for full license details.

### What this means?

- ✅ Free to use, modify, and distribute
- ✅ Commercial use allowed
- ⚠️ Must share source code if you distribute
- ⚠️ Network use triggers copyleft (AGPLv3 special sauce)
- 📋 Must include license and copyright notices

**TL;DR**: Keep it open source, keep it free. Don't be a corporate asshole.

---

**Disclaimer:** Not responsible for phone screen damage caused by ill-advised photo documentation. Use with pride, shame, or scientific detachment. Up to you.

---

**Golden Rule:** If your code stinks, it won't get merged. And yes, we'll know. 💩

Ready to track some legendary logs? Get started above!
