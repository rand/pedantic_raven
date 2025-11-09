"""
GLiNER2 NER Service - FastAPI application.

Provides HTTP REST API for named entity recognition using GLiNER2 model.
"""

import logging
from contextlib import asynccontextmanager
from typing import List, Dict

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field

from model import get_model

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


# Request/Response models
class ExtractRequest(BaseModel):
    """Request for entity extraction."""
    text: str = Field(..., description="Text to analyze", min_length=1)
    entity_types: List[str] = Field(
        ...,
        description="Entity types to extract (e.g., ['person', 'organization'])",
        min_items=1
    )
    threshold: float = Field(
        0.3,
        description="Confidence threshold (0.0-1.0)",
        ge=0.0,
        le=1.0
    )


class Entity(BaseModel):
    """Extracted entity."""
    text: str = Field(..., description="Entity text")
    label: str = Field(..., description="Entity type/label")
    start: int = Field(..., description="Start character index")
    end: int = Field(..., description="End character index")
    score: float = Field(..., description="Confidence score")


class ExtractResponse(BaseModel):
    """Response from entity extraction."""
    entities: List[Entity]
    entity_count: int
    text_length: int


class HealthResponse(BaseModel):
    """Health check response."""
    status: str
    model_loaded: bool
    model_name: str


class ModelInfoResponse(BaseModel):
    """Model information response."""
    model_name: str
    loaded: bool
    model_type: str
    parameters: str
    license: str


# Application lifespan
@asynccontextmanager
async def lifespan(app: FastAPI):
    """Load model on startup, cleanup on shutdown."""
    logger.info("Starting GLiNER service...")

    # Warmup: Load model on startup
    try:
        model = get_model()
        model.load()
        logger.info("Model loaded and ready")
    except Exception as e:
        logger.error(f"Failed to load model on startup: {e}")
        # Continue anyway - model will load on first request

    yield

    logger.info("Shutting down GLiNER service...")


# Create FastAPI app
app = FastAPI(
    title="GLiNER NER Service",
    description="Named Entity Recognition service using GLiNER2 model",
    version="1.0.0",
    lifespan=lifespan
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Allow all origins (adjust for production)
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health", response_model=HealthResponse)
async def health_check():
    """
    Health check endpoint.

    Returns service status and model loading state.
    """
    model = get_model()
    return HealthResponse(
        status="healthy",
        model_loaded=model.is_loaded(),
        model_name=model.model_name
    )


@app.get("/model_info", response_model=ModelInfoResponse)
async def model_info():
    """
    Get model metadata.

    Returns information about the loaded GLiNER2 model.
    """
    model = get_model()
    info = model.model_info()
    return ModelInfoResponse(**info)


@app.post("/extract", response_model=ExtractResponse)
async def extract_entities(request: ExtractRequest):
    """
    Extract named entities from text.

    Args:
        request: ExtractRequest with text and entity_types

    Returns:
        ExtractResponse with list of entities

    Raises:
        HTTPException: If extraction fails
    """
    try:
        model = get_model()

        # Extract entities
        entities = model.extract_entities(
            text=request.text,
            entity_types=request.entity_types,
            threshold=request.threshold
        )

        return ExtractResponse(
            entities=[Entity(**e) for e in entities],
            entity_count=len(entities),
            text_length=len(request.text)
        )

    except Exception as e:
        logger.error(f"Entity extraction failed: {e}")
        raise HTTPException(status_code=500, detail=f"Extraction failed: {str(e)}")


@app.get("/")
async def root():
    """Root endpoint with service info."""
    return {
        "service": "GLiNER NER Service",
        "version": "1.0.0",
        "model": "fastino/gliner2-large-v1",
        "endpoints": {
            "health": "/health",
            "model_info": "/model_info",
            "extract": "/extract (POST)"
        }
    }


if __name__ == "__main__":
    import uvicorn

    # Run with: python main.py
    # Or: uvicorn main:app --host 127.0.0.1 --port 8765 --reload
    uvicorn.run(
        "main:app",
        host="127.0.0.1",
        port=8765,
        log_level="info"
    )
