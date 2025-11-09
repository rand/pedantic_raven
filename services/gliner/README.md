# GLiNER NER Service

FastAPI service providing Named Entity Recognition using GLiNER2 model.

## What is GLiNER?

GLiNER (Generalist Model for Named Entity Recognition) is a zero-shot NER model that can extract **any entity types** you specify at runtime, without retraining. It outperforms ChatGPT on NER benchmarks while running entirely locally on CPU.

**Model**: `fastino/gliner2-large-v1`
- 340M parameters
- CPU-optimized (no GPU required)
- Apache 2.0 license
- Supports custom entity types

## Quick Start

### Installation

```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
```

### Run Service

```bash
# Option 1: Direct Python
python main.py

# Option 2: Uvicorn (recommended)
uvicorn main:app --host 127.0.0.1 --port 8765 --reload

# Option 3: Production (with workers)
uvicorn main:app --host 0.0.0.0 --port 8765 --workers 4
```

The service will start at `http://localhost:8765`

### First Run

On first run, the service will download the GLiNER2 model (~680MB) from HuggingFace. The model is cached in `~/.cache/huggingface/` for subsequent runs.

## API Endpoints

### Health Check

```bash
curl http://localhost:8765/health
```

Response:
```json
{
  "status": "healthy",
  "model_loaded": true,
  "model_name": "fastino/gliner2-large-v1"
}
```

### Extract Entities

```bash
curl -X POST http://localhost:8765/extract \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Alice works at Google in San Francisco. She uses PostgreSQL for databases.",
    "entity_types": ["person", "organization", "location", "technology"],
    "threshold": 0.3
  }'
```

Response:
```json
{
  "entities": [
    {
      "text": "Alice",
      "label": "person",
      "start": 0,
      "end": 5,
      "score": 0.95
    },
    {
      "text": "Google",
      "label": "organization",
      "start": 16,
      "end": 22,
      "score": 0.98
    },
    {
      "text": "San Francisco",
      "label": "location",
      "start": 26,
      "end": 39,
      "score": 0.94
    },
    {
      "text": "PostgreSQL",
      "label": "technology",
      "start": 51,
      "end": 61,
      "score": 0.92
    }
  ],
  "entity_count": 4,
  "text_length": 81
}
```

### Model Info

```bash
curl http://localhost:8765/model_info
```

Response:
```json
{
  "model_name": "fastino/gliner2-large-v1",
  "loaded": true,
  "model_type": "GLiNER2",
  "parameters": "340M",
  "license": "Apache 2.0"
}
```

## Usage Examples

### Custom Entity Types

Extract domain-specific entities:

```json
{
  "text": "The API endpoint /users/login accepts POST requests. It returns a JWT token.",
  "entity_types": ["api_endpoint", "http_method", "security_token"],
  "threshold": 0.3
}
```

### Technical Documentation

Extract architecture concepts:

```json
{
  "text": "The system uses microservices with event-driven architecture. Each service has its own PostgreSQL database.",
  "entity_types": ["architecture_pattern", "technology", "database"],
  "threshold": 0.3
}
```

### Adjust Threshold

Lower threshold for more entities (may include false positives):

```json
{
  "text": "...",
  "entity_types": ["..."],
  "threshold": 0.2
}
```

Higher threshold for higher confidence (may miss some entities):

```json
{
  "text": "...",
  "entity_types": ["..."],
  "threshold": 0.5
}
```

## Testing

```bash
# Run tests
pytest

# Run with coverage
pytest --cov=. --cov-report=html

# Test specific file
pytest tests/test_api.py -v
```

## Performance

**Typical Performance (CPU):**
- Startup: ~3-5 seconds (model loading)
- Extraction: ~100-300ms per request (varies with text length)
- Memory: ~700MB (model) + ~300MB (Python runtime)

**Tips for Better Performance:**
- Use batch processing for multiple texts
- Adjust `threshold` to reduce false positives
- Use medium model (`gliner2-medium-v1`, 205M params) for faster inference
- Run on server with more cores for parallel requests

## Troubleshooting

### Model Download Fails

If download fails, you can manually download the model:

```bash
python -c "from gliner2 import GLiNER; GLiNER.from_pretrained('fastino/gliner2-large-v1')"
```

### Out of Memory

If service crashes with OOM, try:
1. Use medium model (205M params instead of 340M)
2. Reduce batch size in requests
3. Increase system swap space

### Slow Inference

If extraction is slow (>1s per request):
1. Check CPU usage (should be near 100% during extraction)
2. Consider using smaller texts (split long documents)
3. Switch to medium model for faster inference

## Development

### Project Structure

```
services/gliner/
├── main.py           # FastAPI application
├── model.py          # GLiNER2 model wrapper
├── requirements.txt  # Python dependencies
├── tests/            # Test suite
└── README.md         # This file
```

### Adding Features

- **New endpoints**: Add to `main.py`
- **Model changes**: Update `model.py`
- **Dependencies**: Update `requirements.txt`

## License

This service uses:
- GLiNER2 model: Apache 2.0
- FastAPI: MIT
- Service code: Same as Pedantic Raven project

## Links

- **GLiNER2 Paper**: https://arxiv.org/html/2507.18546v1
- **Model on HuggingFace**: https://huggingface.co/fastino/gliner2-large-v1
- **GLiNER GitHub**: https://github.com/urchade/GLiNER
- **FastAPI Docs**: https://fastapi.tiangolo.com
