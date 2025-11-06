# Poo Tracker AI Service

> AI-powered analysis for bowel movement patterns and correlations

This is a FastAPI-based microservice that provides intelligent analysis of bowel movement data, meal
correlations, and generates personalized health insights and recommendations.

## 🚀 Features

- **Bristol Stool Chart Analysis**: Comprehensive analysis of stool types and patterns
- **Meal Correlation Tracking**: Links dietary choices to digestive outcomes
- **Pattern Recognition**: Identifies timing patterns and frequency trends
- **Health Risk Assessment**: Flags potential digestive health concerns
- **Personalized Recommendations**: AI-generated suggestions for digestive health improvement
- **Redis Caching**: Efficient caching of analysis results
- **RESTful API**: Clean, well-documented API endpoints

## 🛠 Tech Stack

- **Framework**: FastAPI 0.115+
- **Language**: Python 3.9+
- **Data Analysis**: pandas, numpy, scikit-learn
- **Caching**: Redis (optional, with graceful degradation)
- **API Documentation**: Built-in Swagger/OpenAPI
- **Testing**: pytest with asyncio support
- **Code Quality**: ruff for linting and formatting
- **Package Manager**: uv (fast Python package manager)
- **Containerization**: Docker

## 📋 Prerequisites

- Python 3.9 or higher
- Redis server (optional, service degrades gracefully without it)
- uv package manager (recommended) or pip

## 🔧 Installation

### Using uv (Recommended)

```bash
# Install uv if you haven't already
curl -LsSf https://astral.sh/uv/install.sh | sh

# Install dependencies
uv sync

# Activate virtual environment
source .venv/bin/activate # Linux/Mac
# or
.venv\Scripts\activate # Windows
```

### Using pip

```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate # Linux/Mac
# or
venv\Scripts\activate # Windows

# Install dependencies
pip install -e .
```

### Development Dependencies

```bash
# Install with development dependencies
uv sync --group dev
```

## 🏃‍♂️ Running the Service

### Development Server

```bash
# Using uvicorn directly
uvicorn ai_service.main:app --reload --host 0.0.0.0 --port 8000

# Or using Python
python main.py
```

### Production Server

```bash
# Using uvicorn with production settings
uvicorn ai_service.main:app --host 0.0.0.0 --port 8000 --workers 4
```

### Docker

```bash
# Build the image
docker build -t poo-tracker-ai .

# Run the container
docker run -p 8000:8000 poo-tracker-ai

# With Redis connection
docker run -p 8000:8000 -e REDIS_URL=redis://your-redis-host:6379 poo-tracker-ai
```

#### Docker + Redis Quickstart

Running the container with a Redis instance lets the service cache results and improves overall
performance. Specify the `REDIS_URL` and `PORT` environment variables when starting the image.

```bash
# Launch Redis (optional)
docker run -d --name poo-redis -p 6379:6379 redis:7

# Build the AI service image
docker build -t poo-tracker-ai .

# Run the service with environment variables
docker run --rm -p 8000:8000 \
  -e REDIS_URL=redis://localhost:6379 \
  -e PORT=8000 \
  poo-tracker-ai
```

Once running, the API exposes the following endpoints:

- `http://localhost:8000/health`
- `http://localhost:8000/analyze`
- `http://localhost:8000/docs`

## 🌐 API Endpoints

### Health Check

- **GET** `/health` - Service health status and Redis connectivity

### Analysis

- **POST** `/analyze` - Analyze bowel movement patterns and meal correlations

### API Documentation

- **GET** `/docs` - Interactive Swagger UI documentation
- **GET** `/redoc` - Alternative ReDoc documentation

## 📊 API Usage Examples

### Health Status

```bash
curl http://localhost:8000/health
```

```json
{
  "status": "healthy",
  "timestamp": "2025-06-16T23:30:00",
  "redis_connected": true
}
```

### Pattern Analysis

```bash
curl -X POST http://localhost:8000/analyze \\
-H "Content-Type: application/json" \\
-d '{
    "entries": [
      {
        "id": "entry-1",
        "userId": "user-123",
        "bristolType": 4,
        "volume": "medium",
        "color": "brown",
        "consistency": "normal",
        "floaters": false,
        "pain": 2,
        "strain": 1,
        "satisfaction": 8,
        "createdAt": "2025-06-16T08:30:00Z"
      }
    ],
    "meals": [
      {
        "id": "meal-1",
        "userId": "user-123",
        "name": "Oatmeal with berries",
        "mealTime": "2025-06-15T07:00:00Z",
        "category": "breakfast",
        "spicyLevel": 0,
        "fiberRich": true,
        "dairy": false,
        "gluten": true
      }
    ]
  }'
```

## 🧪 Testing

```bash
# Run all tests
pytest

# Run with coverage
pytest --cov=main

# Run specific test file
pytest test_main.py

# Run tests in verbose mode
pytest -v
```

## 🔍 Code Quality

```bash
# Format code
black .
isort .

# Lint code
ruff check .

# Type checking
mypy main.py

# Run all quality checks
ruff check . && black --check . && isort --check . && mypy main.py
```

## 🗂 Project Structure

```text
ai-service/
├── main.py              # FastAPI application and analysis logic
├── test_main.py         # Test suite
├── pyproject.toml       # Project configuration and dependencies
├── Dockerfile           # Container configuration
└── README.md           # This file
```

## ⚙️ Configuration

### Environment Variables

Copy `ai-service/.env.example` to `.env` and update values as needed.

- `REDIS_URL` - Redis connection URL (default: `redis://localhost:6379`)
- `PORT` - Service port (default: 8000)

### Redis Configuration

The service gracefully handles Redis unavailability:

- ✅ **With Redis**: Analysis results are cached for improved performance
- ✅ **Without Redis**: Service operates normally but without caching

## 📈 Analysis Features

### Bristol Stool Chart Analysis

- Distribution of stool types 1-7
- Health assessment based on patterns
- Trend analysis over time

### Meal Correlation Analysis

- Links meals to bowel movements within 6-48 hour digestion window
- Analyzes impact of food categories, spiciness, fiber content
- Identifies problematic foods or beneficial patterns

### Pattern Recognition

- Timing patterns (hourly and daily distributions)
- Frequency statistics and trends
- Consistency analysis over time

### Health Insights

- Personalized recommendations based on data patterns
- Risk factor identification
- Digestive health scoring

## 🚨 Health Disclaimers

This AI service provides **educational insights only** and should not replace professional medical
advice. Users experiencing persistent digestive issues should consult healthcare providers.

## 🤝 Contributing

1. Follow the coding standards defined in `pyproject.toml`
1. Run tests and quality checks before submitting
1. Update tests for new features
1. Follow the existing code style and patterns

## 📝 License

This project is part of the Poo Tracker application suite.

## 🔗 Integration

This service is designed to work with:

- **Backend API**: Node.js/Express backend at `../backend/`
- **Frontend**: React application at `../frontend/`
- **Database**: PostgreSQL with Prisma ORM
- **Caching**: Redis for performance optimization

---

_Built with ❤️ and a healthy sense of humor for the Poo Tracker project._
