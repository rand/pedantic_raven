# GLiNER Integration Guide

**Status**: Experimental Feature (Phase 5)
**Last Updated**: 2025-11-08

---

## Overview

Pedantic Raven now supports **GLiNER2** for enhanced named entity recognition. GLiNER is a zero-shot NER model that extracts entities more accurately than pattern matching and allows custom entity types.

### What You Get

**Before (Pattern Matching)**:
- 6 hardcoded entity types
- Keyword-based classification ("Dr" → Person, "Inc" → Organization)
- Limited accuracy on ambiguous cases
- No context understanding

**After (GLiNER2)**:
- **ANY entity types** you define
- ML-based extraction (outperforms ChatGPT on NER benchmarks)
- Context-aware classification
- Handles complex, ambiguous text

### Example Improvement

**Text**: "Alice manages the authentication API at Google"

**Pattern Matcher** (old):
- ✓ Alice → Person (has capital letter)
- ✗ authentication → (missed)
- ✗ API → (missed)
- ✓ Google → Organization (proper noun)

**GLiNER2** (new):
- ✓ Alice → person (confidence: 0.95)
- ✓ authentication → security_concept (confidence: 0.88)
- ✓ API → technology (confidence: 0.92)
- ✓ Google → organization (confidence: 0.98)

---

## Quick Start

### Prerequisites

- Python 3.9+ installed
- ~1GB disk space (for model download)
- ~1GB RAM (for model inference)

### 1. Start GLiNER Service

```bash
# Navigate to service directory
cd services/gliner

# Create virtual environment (first time only)
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies (first time only)
pip install -r requirements.txt

# Start service
uvicorn main:app --host 127.0.0.1 --port 8765
```

The service will download the GLiNER2 model (~680MB) on first run. This is cached for future runs.

### 2. Run Pedantic Raven

```bash
# In a new terminal
cd ../..
go build -o pedantic_raven .
./pedantic_raven
```

### 3. Verify GLiNER is Active

In Pedantic Raven's status line (bottom of screen), you should see:
```
[GLiNER: ✓]  # Green checkmark = active
[Pattern]    # Orange = fallback to pattern matcher
```

---

## Configuration

### Environment Variables

```bash
# Enable/disable GLiNER
export GLINER_ENABLED=true

# Service URL
export GLINER_SERVICE_URL=http://localhost:8765

# Timeout (seconds)
export GLINER_TIMEOUT=5

# Fall back to pattern matcher if GLiNER unavailable
export GLINER_FALLBACK=true
```

### Config File

Create `config.toml` in project root:

```toml
[gliner]
enabled = true
service_url = "http://localhost:8765"
timeout = 5
fallback_to_pattern = true

[gliner.entity_types]
# Default entity types for extraction
default = [
    "person",
    "organization",
    "location",
    "technology",
    "concept",
    "product",
    "api_endpoint",
    "database",
    "security_concern"
]

# Add custom types for your domain
custom = [
    "architecture_pattern",
    "deployment_target",
    "performance_metric"
]
```

### Custom Entity Types

You can define domain-specific entity types:

**Technical Documentation**:
```toml
custom = [
    "api_endpoint",       # /users/login, /api/v1/posts
    "http_method",        # GET, POST, PUT, DELETE
    "database_table",     # users, posts, sessions
    "environment_var",    # DATABASE_URL, API_KEY
    "architecture_pattern" # microservices, event-driven
]
```

**Product Documentation**:
```toml
custom = [
    "feature_name",       # Dark mode, Auto-save
    "user_segment",       # Premium users, Free tier
    "metric",             # conversion rate, DAU
    "competitor"          # ProductX, ServiceY
]
```

**Security Audit**:
```toml
custom = [
    "vulnerability",      # SQL injection, XSS
    "security_tool",      # WAF, IDS, SIEM
    "compliance_standard", # GDPR, HIPAA, PCI-DSS
    "threat_actor"        # APT28, insider threat
]
```

---

## Usage

### In Edit Mode

When you type in the editor, semantic analysis runs automatically (500ms debounce). If GLiNER is active, it will extract entities using the configured entity types.

**Context Panel** will show:
- **Entities** section with GLiNER-extracted entities
- Entity types and confidence scores
- Occurrence counts

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Ctrl+G` | Toggle GLiNER on/off |
| `Ctrl+Shift+G` | Force re-analysis with GLiNER |
| `?` then `G` | Show GLiNER status and config |

### Status Indicators

**Status Line** (bottom of screen):
- `[GLiNER: ✓]` - GLiNER active and working
- `[GLiNER: ⚠]` - GLiNER service unavailable (using pattern fallback)
- `[GLiNER: ✗]` - GLiNER disabled in config
- `[Pattern]` - Using pattern matcher

---

## Deployment Options

### Option 1: Manual (Development)

**Terminal 1 - GLiNER Service**:
```bash
cd services/gliner
source venv/bin/activate
uvicorn main:app --host 127.0.0.1 --port 8765
```

**Terminal 2 - Pedantic Raven**:
```bash
./pedantic_raven
```

### Option 2: Docker Compose (Recommended)

```bash
# Start both services
docker-compose up

# Stop services
docker-compose down
```

Services:
- GLiNER: http://localhost:8765
- Pedantic Raven: Runs in terminal

### Option 3: Systemd Service (Production)

**Install GLiNER Service**:
```bash
sudo cp services/gliner/gliner.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable gliner
sudo systemctl start gliner
```

**Check Status**:
```bash
sudo systemctl status gliner
```

**Logs**:
```bash
sudo journalctl -u gliner -f
```

---

## Performance

### Expected Performance

**GLiNER2-Large (340M params)**:
- Startup: 3-5 seconds (model loading)
- Inference: 100-300ms per analysis (typical document)
- Memory: ~700MB model + ~300MB Python runtime (~1GB total)
- CPU: 100% usage during extraction (normal)

**Integration Overhead**:
- HTTP round-trip: 5-20ms (localhost)
- JSON serialization: 1-5ms
- Total latency: 110-325ms (well within 500ms debounce)

### Optimization Tips

**If extraction is slow (>500ms)**:
1. Use medium model (`gliner2-medium-v1`, 205M params)
2. Reduce entity types (fewer types = faster extraction)
3. Increase analysis debounce in config

**If memory usage is high**:
1. Use medium model (uses ~400MB instead of ~700MB)
2. Restart service periodically (Docker handles this)

**If CPU usage is high**:
- This is normal during extraction
- Service is idle when not analyzing

---

## Troubleshooting

### Service Won't Start

**Symptom**: `uvicorn` fails to start

**Solutions**:
```bash
# Check Python version (need 3.9+)
python --version

# Verify dependencies installed
pip list | grep gliner2

# Reinstall dependencies
pip install -r requirements.txt --upgrade

# Check port not in use
lsof -i :8765
```

### Model Download Fails

**Symptom**: "Failed to download model from HuggingFace"

**Solutions**:
```bash
# Manual download
python -c "from gliner2 import GLiNER; GLiNER.from_pretrained('fastino/gliner2-large-v1')"

# Check internet connection
ping huggingface.co

# Check disk space (need ~1GB)
df -h ~/.cache/huggingface
```

### Service Unavailable in Pedantic Raven

**Symptom**: Status shows `[GLiNER: ⚠]`

**Solutions**:
```bash
# Test service directly
curl http://localhost:8765/health

# Check service logs
# (in service terminal, look for errors)

# Verify URL in config
echo $GLINER_SERVICE_URL

# Try manual connection
python -c "import requests; print(requests.get('http://localhost:8765/health').json())"
```

### Extraction Returns No Entities

**Symptom**: Context panel shows no entities despite text

**Possible Causes**:
1. **Threshold too high** - Lower to 0.2 or 0.3
2. **Wrong entity types** - Verify types match your text domain
3. **Text too short** - GLiNER works best with complete sentences
4. **Model not loaded** - Check service logs for "Model loaded"

**Solutions**:
```bash
# Test extraction directly
curl -X POST http://localhost:8765/extract \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Alice works at Google",
    "entity_types": ["person", "organization"],
    "threshold": 0.2
  }'

# Should return entities
```

### Slow Performance

**Symptom**: Lag when typing in editor

**Solutions**:
1. Check GLiNER service CPU usage (`top` or `htop`)
2. Increase analysis debounce (default 500ms → 1000ms)
3. Switch to medium model for faster inference
4. Reduce number of entity types

---

## Advanced Usage

### Batch Processing

For analyzing multiple files:

```python
import requests

texts = ["doc1 content...", "doc2 content...", "doc3 content..."]
results = []

for text in texts:
    resp = requests.post("http://localhost:8765/extract", json={
        "text": text,
        "entity_types": ["person", "organization", "location"],
        "threshold": 0.3
    })
    results.append(resp.json())
```

### Custom Model

To use a different GLiNER model:

**Edit `services/gliner/model.py`**:
```python
class GLiNERModel:
    def __init__(self, model_name: str = "urchade/gliner_medium-v2.1"):  # Changed
        # ...
```

Available models:
- `fastino/gliner2-large-v1` (340M, best accuracy)
- `urchade/gliner_medium-v2.1` (205M, faster)
- `urchade/gliner_small-v2.1` (82M, fastest)

### Monitoring

**Service Health**:
```bash
watch -n 5 'curl -s http://localhost:8765/health | jq'
```

**Model Info**:
```bash
curl http://localhost:8765/model_info | jq
```

**Request Metrics** (add to service):
```python
# In main.py
from prometheus_client import Counter, Histogram

extract_requests = Counter('gliner_extract_requests_total', 'Total extraction requests')
extract_latency = Histogram('gliner_extract_latency_seconds', 'Extraction latency')
```

---

## Comparison: Pattern vs GLiNER

| Feature | Pattern Matcher | GLiNER2 |
|---------|----------------|---------|
| Entity Types | 6 fixed | Unlimited custom |
| Accuracy | ~60-70% | ~85-95% |
| Context-Aware | No | Yes |
| Speed | <1ms | 100-300ms |
| Memory | ~0MB | ~1GB |
| Setup | None | Python service |
| Offline | Yes | Yes |
| Custom Types | No | Yes |
| Ambiguity Handling | Poor | Excellent |

### When to Use Pattern Matcher

- Low-latency required (<10ms)
- Limited resources (no 1GB RAM available)
- Simple, unambiguous text
- No Python installation possible

### When to Use GLiNER

- High accuracy critical
- Custom entity types needed
- Ambiguous or complex text
- Context engineering for AI (Pedantic Raven's primary use case)

---

## FAQ

### Q: Does GLiNER send data to external servers?

**A**: No. GLiNER runs 100% locally. The model downloads from HuggingFace on first run, but after that, all processing is local.

### Q: Can I use both GLiNER and pattern matcher?

**A**: Yes. Set `fallback_to_pattern = true` in config. Pedantic Raven will use GLiNER when available, falling back to pattern matcher if GLiNER service is down.

### Q: How much does GLiNER slow down analysis?

**A**: Typical analysis goes from <1ms (pattern) to ~200ms (GLiNER). Since Pedantic Raven debounces analysis by 500ms, you won't notice the difference during typing.

### Q: Can I run GLiNER on a different machine?

**A**: Yes. Set `service_url = "http://other-machine:8765"` in config. Ensure firewall allows connections.

### Q: What if I don't have Python?

**A**: GLiNER is optional. Pedantic Raven works fine with pattern matcher. GLiNER just provides better accuracy.

### Q: Can I contribute custom entity types?

**A**: Yes! We welcome PRs with domain-specific entity type collections. Add them to `config.toml` examples.

---

## Development

### Running Tests

**Python Service**:
```bash
cd services/gliner
pytest
pytest --cov=. --cov-report=html
```

**Go Client**:
```bash
go test ./internal/gliner/...
go test ./internal/gliner/... -v
```

**Integration Tests**:
```bash
# Start service first
cd services/gliner && uvicorn main:app &

# Run integration tests
go test ./internal/gliner/... -tags=integration
```

### Adding Custom Entity Types

**Step 1**: Define types in config
**Step 2**: Test with sample text
**Step 3**: Adjust threshold if needed
**Step 4**: Document in your project README

---

## Links

- **GLiNER2 Paper**: https://arxiv.org/html/2507.18546v1
- **Model on HuggingFace**: https://huggingface.co/fastino/gliner2-large-v1
- **GLiNER GitHub**: https://github.com/urchade/GLiNER
- **Service README**: [services/gliner/README.md](../services/gliner/README.md)
- **Go Client API**: [internal/gliner/](../internal/gliner/)

---

## Support

**Issues**: [GitHub Issues](https://github.com/rand/pedantic_raven/issues)
**Questions**: [GitHub Discussions](https://github.com/rand/pedantic_raven/discussions)

When reporting issues, include:
- GLiNER service version (`pip show gliner2`)
- Service logs (from uvicorn output)
- Pedantic Raven version (`./pedantic_raven --version`)
- Sample text causing issues
- Your config.toml

---

**Status**: Experimental - Feedback Welcome!
**Last Updated**: 2025-11-08
