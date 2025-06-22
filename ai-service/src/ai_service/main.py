"""
Enhanced FastAPI application for Poo Tracker AI Service.

This is the main entry point for the AI service that provides intelligent
analysis of bowel movement patterns, meal correlations, and health insights.
"""

import time
from contextlib import asynccontextmanager
from datetime import datetime

from fastapi import FastAPI, HTTPException, Request, status
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from fastapi.responses import JSONResponse

from .config.logging import get_logger, setup_logging
from .config.settings import get_settings
from .models.database import BowelMovementData, MealData, SymptomData
from .models.requests import AnalysisRequest
from .models.responses import AnalysisResponse, ErrorResponse, HealthResponse
from .services.analyzer import AnalyzerService
from .services.health_assessor import HealthAssessorService
from .services.recommender import RecommenderService
from .utils.cache import CacheManager
from .utils.validators import DataValidator

# Initialize logging
setup_logging()
logger = get_logger("main")

# Get settings
settings = get_settings()


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager."""
    logger.info("🚀 Starting Poo Tracker AI Service")

    # Initialize services
    app.state.cache_manager = CacheManager()
    app.state.analyzer = AnalyzerService()
    app.state.health_assessor = HealthAssessorService()
    app.state.recommender = RecommenderService()
    app.state.validator = DataValidator()

    # Test Redis connection
    try:
        await app.state.cache_manager.ping()
        logger.info("✅ Redis connection established")
    except Exception as e:
        logger.warning(f"⚠️ Redis connection failed: {e}")

    logger.info("🎉 AI Service startup complete")
    yield

    # Cleanup
    logger.info("🛑 Shutting down AI Service")
    if hasattr(app.state, "cache_manager"):
        await app.state.cache_manager.close()


# Create FastAPI app
app = FastAPI(
    title=settings.app_name,
    description="AI-powered analysis for bowel movement patterns and correlations",
    version=settings.app_version,
    docs_url="/docs" if not settings.is_production else None,
    redoc_url="/redoc" if not settings.is_production else None,
    lifespan=lifespan,
)

# Add middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"] if settings.is_development else [settings.backend_url],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.add_middleware(GZipMiddleware, minimum_size=1000)


@app.middleware("http")
async def add_process_time_header(request: Request, call_next):
    """Add processing time to response headers."""
    start_time = time.time()
    response = await call_next(request)
    process_time = time.time() - start_time
    response.headers["X-Process-Time"] = str(process_time)
    return response


@app.exception_handler(HTTPException)
async def http_exception_handler(request: Request, exc: HTTPException):
    """Custom HTTP exception handler."""
    return JSONResponse(
        status_code=exc.status_code,
        content=ErrorResponse(
            error=exc.detail,
            timestamp=datetime.now(),
            request_id=getattr(request.state, "request_id", None),
        ).model_dump(),
    )


@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception):
    """General exception handler."""
    logger.error(f"Unhandled exception: {exc}", exc_info=True)

    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content=ErrorResponse(
            error="Internal server error",
            detail=str(exc) if settings.is_development else None,
            timestamp=datetime.now(),
        ).model_dump(),
    )


@app.get("/health", response_model=HealthResponse)
async def health_check():
    """
    Health check endpoint.

    Returns the current health status of the AI service including:
    - Service status
    - Redis connectivity
    - ML model availability
    - Performance metrics
    """
    start_time = time.time()

    # Check Redis connection
    redis_connected = False
    try:
        redis_connected = await app.state.cache_manager.ping()
    except Exception as e:
        logger.warning(f"Redis health check failed: {e}")

    # Check ML models (placeholder for now)
    ml_models_loaded = True  # Would check actual model loading status

    # Calculate response time
    response_time_ms = (time.time() - start_time) * 1000

    # Determine overall status
    if redis_connected and ml_models_loaded:
        status_val = "healthy"
    elif ml_models_loaded:
        status_val = "degraded"  # Service works but caching is unavailable
    else:
        status_val = "unhealthy"

    return HealthResponse(
        status=status_val,
        redis_connected=redis_connected,
        ml_models_loaded=ml_models_loaded,
        response_time_ms=response_time_ms,
        version=settings.app_version,
    )


@app.get("/")
async def root():
    """Root endpoint with service information."""
    return {
        "service": settings.app_name,
        "version": settings.app_version,
        "status": "running",
        "docs": "/docs",
        "health": "/health",
        "analysis": "/analyze",
    }


@app.post("/analyze", response_model=AnalysisResponse)
async def analyze_patterns(request: AnalysisRequest):
    """
    Analyze bowel movement patterns and correlations with meals.

    This endpoint performs comprehensive analysis including:
    - Bristol Stool Chart pattern analysis
    - Timing and frequency patterns
    - Meal correlation analysis
    - Health risk assessment
    - Personalized recommendations

    The analysis considers digestion time windows (6-48 hours) when
    correlating meals with bowel movements.
    """
    logger.info(f"Starting analysis for {len(request.entries)} entries")

    try:
        # Validate input data
        validation_result = app.state.validator.validate_analysis_request(request)
        if not validation_result.is_valid:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail=f"Validation failed: {validation_result.errors}",
            )

        # Check cache first
        cache_key = await app.state.cache_manager.generate_analysis_cache_key(request)
        cached_result = await app.state.cache_manager.get_analysis_result(cache_key)

        if cached_result:
            logger.info("Returning cached analysis result")
            return AnalysisResponse(**cached_result)

        # Convert request models to internal data models
        bowel_movements = [
            BowelMovementData(
                id=entry.id,
                user_id=entry.user_id,
                bristol_type=entry.bristol_type,
                volume=entry.volume,
                color=entry.color,
                consistency=entry.consistency,
                floaters=entry.floaters,
                pain=entry.pain,
                strain=entry.strain,
                satisfaction=entry.satisfaction,
                created_at=entry.created_at,
                recorded_at=entry.recorded_at,
            )
            for entry in request.entries
        ]

        meals = [
            MealData(
                id=meal.id,
                user_id=meal.user_id,
                name=meal.name,
                meal_time=meal.meal_time,
                category=meal.category,
                cuisine=meal.cuisine,
                spicy_level=meal.spicy_level,
                fiber_rich=meal.fiber_rich,
                dairy=meal.dairy,
                gluten=meal.gluten,
                created_at=meal.created_at,
            )
            for meal in request.meals
        ]

        symptoms = [
            SymptomData(
                id=symptom.id,
                user_id=symptom.user_id,
                bowel_movement_id=symptom.bowel_movement_id,
                type=symptom.type,
                severity=symptom.severity,
                notes=symptom.notes,
                created_at=symptom.created_at,
                recorded_at=symptom.recorded_at,
            )
            for symptom in request.symptoms
        ]

        # Perform core analysis
        analysis_result = await app.state.analyzer.analyze_comprehensive_patterns(
            bowel_movements=bowel_movements,
            meals=meals if meals else None,
            symptoms=symptoms if symptoms else None,
            user_id=bowel_movements[0].user_id if bowel_movements else None,
        )

        # Generate health assessment
        health_score = None
        if request.include_predictions:
            health_score = await app.state.health_assessor.calculate_health_score(
                bowel_movements=bowel_movements,
                symptoms=symptoms if symptoms else None,
            )

        # Generate recommendations
        recommendations = []
        risk_factors = []
        if request.include_recommendations:
            recommendations = await app.state.recommender.generate_recommendations(
                analysis_result=analysis_result,
                bowel_movements=bowel_movements,
                meals=meals if meals else None,
            )

            risk_factors = await app.state.recommender.identify_risk_factors(
                bowel_movements=bowel_movements,
                analysis_result=analysis_result,
            )

        # Create response
        response = AnalysisResponse(
            patterns=analysis_result["patterns"],
            correlations=analysis_result["correlations"],
            bristol_trends=analysis_result["bristol_analysis"],
            recommendations=recommendations,
            risk_factors=risk_factors,
            health_score=health_score,
            predictions=None,  # Would be populated by ML models
            analysis_metadata=analysis_result["metadata"],
        )

        # Cache the result
        await app.state.cache_manager.cache_analysis_result(
            cache_key, response.model_dump(), ttl=settings.cache_ttl
        )

        logger.info("Analysis completed successfully")
        return response

    except ValueError as e:
        logger.error(f"Validation error: {e}")
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail=str(e)
        ) from e
    except Exception as e:
        logger.error(f"Analysis failed: {e}", exc_info=True)
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Analysis failed: {str(e)}",
        ) from e


@app.get("/metrics")
async def get_metrics():
    """Get service metrics (for monitoring)."""
    if not settings.is_development:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Metrics endpoint not available in production",
        )

    try:
        cache_stats = await app.state.cache_manager.get_cache_stats()

        return {
            "service_name": settings.app_name,
            "version": settings.app_version,
            "uptime": time.time(),  # Would calculate actual uptime
            "cache_stats": cache_stats,
            "settings": {
                "environment": settings.environment,
                "debug": settings.debug,
                "ml_enabled": settings.enable_ml_features,
            },
        }
    except Exception as e:
        logger.error(f"Metrics retrieval failed: {e}")
        return {"error": "Metrics unavailable"}


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "ai_service.main:app",
        host=settings.host,
        port=settings.port,
        reload=settings.is_development,
        log_level=settings.log_level.lower(),
        workers=1 if settings.is_development else settings.workers,
    )
