# GLiNER NER Service Dockerfile

FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Pre-download GLiNER model (optional, comment out to download on first run)
# RUN python -c "from gliner2 import GLiNER; GLiNER.from_pretrained('fastino/gliner2-large-v1')"

# Copy application code
COPY main.py model.py ./

# Expose port
EXPOSE 8765

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=60s --retries=3 \
  CMD curl -f http://localhost:8765/health || exit 1

# Run service
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8765"]
