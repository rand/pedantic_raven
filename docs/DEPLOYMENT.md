# Pedantic Raven Deployment Guide

**Version**: 2.0
**Last Updated**: November 12, 2025
**Status**: Complete

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Authentication](#authentication)
5. [Deployment Options](#deployment-options)
6. [mnemosyne Setup](#mnemosyne-setup)
7. [GLiNER Integration](#gliner-integration)
8. [Monitoring & Logging](#monitoring--logging)
9. [Troubleshooting](#troubleshooting)
10. [Security Considerations](#security-considerations)
11. [Backup & Recovery](#backup--recovery)
12. [Updating](#updating)

---

## Prerequisites

### System Requirements

**Minimum**:
- CPU: 2 cores
- RAM: 512 MB
- Disk: 100 MB
- Terminal: 80x24 (recommended 120x30+)

**Recommended**:
- CPU: 4+ cores
- RAM: 2+ GB
- Disk: 1 GB (for models if using GLiNER)
- Terminal: 120x30+
- OS: Linux, macOS, or Windows with WSL2

### Software Requirements

**Required**:
- **Go 1.21+** - For building from source
  ```bash
  go version
  # Expected output: go version go1.21.x or higher
  ```

**Optional But Recommended**:
- **Docker** - For containerized deployment
  ```bash
  docker --version
  # Expected output: Docker version 20.10+ or higher
  ```

- **mnemosyne Server** - For persistent memory integration
  - Repository: https://github.com/rand/mnemosyne
  - Requires Rust toolchain
  - Default port: 50051

### Supported Platforms

| Platform | Go | Docker | Status |
|----------|----|---------|-|
| Linux (x86_64) | ✅ | ✅ | Fully supported |
| macOS (Intel) | ✅ | ✅ | Fully supported |
| macOS (ARM64) | ✅ | ✅ | Fully supported |
| Windows (WSL2) | ✅ | ✅ | Fully supported |
| Windows (Native) | ✅ | ✅ | Requires terminal emulator |

---

## Installation

### Option 1: Build from Source

**Clone the repository**:
```bash
git clone https://github.com/rand/pedantic_raven.git
cd pedantic_raven
git checkout main
```

**Download dependencies**:
```bash
go mod download
go mod tidy
```

**Build the binary**:
```bash
# Development build
go build -o pedantic_raven .

# Optimized production build (smaller binary)
go build -ldflags="-s -w" -o pedantic_raven .
```

**Verify the build**:
```bash
./pedantic_raven --version
# Should output version information
```

### Option 2: Binary Installation

**Download pre-built binary**:
```bash
# Download the appropriate binary for your system
# Check https://github.com/rand/pedantic_raven/releases

# Extract
tar xzf pedantic_raven-v1.0.0-linux-amd64.tar.gz

# Make executable
chmod +x pedantic_raven

# Verify
./pedantic_raven --help
```

### Option 3: Docker Installation

**Build Docker image locally**:
```bash
# Clone repository
git clone https://github.com/rand/pedantic_raven.git
cd pedantic_raven

# Build image
docker build -t pedantic-raven:latest .

# Verify image
docker images pedantic-raven
```

**Pull from registry** (future):
```bash
docker pull rand/pedantic-raven:latest
```

### Option 4: Package Manager

**Homebrew** (macOS/Linux):
```bash
brew tap rand/pedantic-raven
brew install pedantic-raven

# Verify
pedantic_raven --help
```

---

## Configuration

### Environment Variables Reference

**Core Settings**:

| Variable | Default | Format | Description |
|----------|---------|--------|-------------|
| `MNEMOSYNE_ENABLED` | `true` | bool | Enable/disable mnemosyne integration |
| `MNEMOSYNE_ADDR` | `localhost:50051` | host:port | mnemosyne server address |
| `MNEMOSYNE_TIMEOUT` | `30` | seconds | Default operation timeout |
| `MNEMOSYNE_MAX_RETRIES` | `3` | int | Maximum retry attempts |
| `GLINER_ENABLED` | `true` | bool | Enable/disable GLiNER entity extraction |
| `GLINER_SERVICE_URL` | `http://localhost:8765` | URL | GLiNER service endpoint |
| `GLINER_TIMEOUT` | `5` | seconds | GLiNER request timeout |
| `GLINER_MAX_RETRIES` | `2` | int | GLiNER retry attempts |
| `GLINER_FALLBACK_TO_PATTERN` | `true` | bool | Fallback to pattern matching if GLiNER unavailable |
| `GLINER_SCORE_THRESHOLD` | `0.3` | 0.0-1.0 | Minimum confidence for entities |

**Environment Variable Examples**:

```bash
# Run with custom mnemosyne server
export MNEMOSYNE_ADDR=192.168.1.100:50051
export MNEMOSYNE_TIMEOUT=60
./pedantic_raven

# Disable GLiNER, use pattern matching only
export GLINER_ENABLED=false
./pedantic_raven

# Use remote GLiNER service
export GLINER_SERVICE_URL=http://gliner-server.example.com:8765
./pedantic_raven

# Stricter entity scoring
export GLINER_SCORE_THRESHOLD=0.7
./pedantic_raven
```

### Configuration File

**Optional TOML configuration** (`config.toml`):

```toml
# GLiNER Configuration
[gliner]
enabled = true
service_url = "http://localhost:8765"
timeout = 5
max_retries = 2
fallback_to_pattern = true
score_threshold = 0.3

# Default entity types
[gliner.entity_types]
default = [
    "person",
    "organization",
    "location",
    "technology",
    "concept",
    "product"
]
custom = [
    "security_concern",
    "api_endpoint",
    "database_type"
]
```

**Load custom config**:
```bash
# Place config.toml in working directory
./pedantic_raven
# or specify path
./pedantic_raven --config /etc/pedantic_raven/config.toml
```

---

## Authentication

### Overview

Pedantic Raven supports simple token-based authentication for single-user deployments. Authentication is **disabled by default** for backward compatibility. When enabled, it requires a secret token to be provided by the application user.

**Key Features**:
- Environment variable-based configuration
- Constant-time token comparison (timing attack resistant)
- Backward compatible (disabled when not configured)
- No performance overhead when disabled
- Suitable for single-user, trusted environments

### Enabling Authentication

**Generate a secure token**:
```bash
# Generate 32-byte random token encoded in base64
openssl rand -base64 32
# Output: yB3xK9pL2m5nQ8rT1uV4wX7yZ0aB3cD6eF9gH2jK5lM8nP0qR3sT6uV9wX2yZ5
```

**Set environment variable**:
```bash
export PEDANTIC_RAVEN_TOKEN="yB3xK9pL2m5nQ8rT1uV4wX7yZ0aB3cD6eF9gH2jK5lM8nP0qR3sT6uV9wX2yZ5"
./pedantic_raven
```

**Verify authentication is enabled**:

The application will log on startup:
```
2025-11-12T10:30:45.123Z INFO  auth: Authentication enabled (PEDANTIC_RAVEN_TOKEN set)
```

### Token Configuration

**Environment variable**:
```bash
export PEDANTIC_RAVEN_TOKEN="your-secret-token-here"
```

**Configuration file** (optional):

While authentication settings are primarily environment-based, you can document them in your deployment:

```toml
# config.toml (informational only, actual token in PEDANTIC_RAVEN_TOKEN env var)
[auth]
# Token is read from environment variable: PEDANTIC_RAVEN_TOKEN
# If not set, authentication is disabled
enabled_via_env = "PEDANTIC_RAVEN_TOKEN"
```

### Best Practices

**Token Generation**:
```bash
# Recommended: Use openssl for cryptographically secure randomness
openssl rand -base64 32

# Alternative: Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"

# Alternative: Go
go run -c 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"; ); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'
```

**Token Storage**:
- Store token in environment variable on deployment server
- Use systemd EnvironmentFile for production deployments
- Never commit token to version control
- Consider using a secrets manager (e.g., HashiCorp Vault, AWS Secrets Manager)

**Token Rotation**:
```bash
# 1. Generate new token
NEW_TOKEN=$(openssl rand -base64 32)

# 2. Update environment variable
export PEDANTIC_RAVEN_TOKEN="$NEW_TOKEN"

# 3. Restart application
sudo systemctl restart pedantic-raven

# 4. Verify in logs
sudo journalctl -u pedantic-raven | grep "Authentication enabled"
```

### Systemd Integration

**Store token securely in systemd service**:

Create `/etc/pedantic_raven/auth.env`:
```bash
# This file contains sensitive data - protect with restricted permissions
PEDANTIC_RAVEN_TOKEN="yB3xK9pL2m5nQ8rT1uV4wX7yZ0aB3cD6eF9gH2jK5lM8nP0qR3sT6uV9wX2yZ5"
```

Set restrictive permissions:
```bash
sudo chmod 600 /etc/pedantic_raven/auth.env
sudo chown pedantic-raven:pedantic-raven /etc/pedantic_raven/auth.env
```

Update systemd service file:
```ini
[Service]
# Load authentication token from environment file
EnvironmentFile=/etc/pedantic_raven/auth.env
ExecStart=/opt/pedantic-raven/bin/pedantic_raven
```

Reload and restart:
```bash
sudo systemctl daemon-reload
sudo systemctl restart pedantic-raven
```

### Docker Integration

**Pass token to container**:
```bash
# Option 1: Via environment variable
docker run -e PEDANTIC_RAVEN_TOKEN="your-token" pedantic-raven:latest

# Option 2: Via .env file
echo "PEDANTIC_RAVEN_TOKEN=your-token" > .env
docker run --env-file .env pedantic-raven:latest

# Option 3: In docker-compose.yml
# .env file
PEDANTIC_RAVEN_TOKEN=your-token-here

# docker-compose.yml
services:
  pedantic-raven:
    build: .
    env_file: .env
    environment:
      - MNEMOSYNE_ADDR=mnemosyne:50051
      - GLINER_SERVICE_URL=http://gliner:8765
```

**Do NOT commit .env files to version control**:
```bash
# .gitignore
.env
.env.local
auth.env
```

### Security Properties

**Timing Attack Resistance**:
- Uses `crypto/subtle.ConstantTimeCompare` for token comparison
- Comparison time is constant regardless of token length or position of mismatch
- Prevents attackers from using timing measurements to guess tokens

**Token Confidentiality**:
- Token is stored in memory (not persisted to disk)
- Token is never logged (private struct field)
- Token not included in debug output
- Environment variable not exposed to child processes unless explicitly passed

**Authentication Disabled by Default**:
- No breaking changes to existing deployments
- Backward compatible with applications not using authentication
- Can be enabled for specific deployments without code changes

### Disabling Authentication

If you need to temporarily disable authentication:
```bash
# Simply don't set the environment variable
unset PEDANTIC_RAVEN_TOKEN
./pedantic_raven

# Or explicitly empty it (same effect as not set)
export PEDANTIC_RAVEN_TOKEN=""
./pedantic_raven
```

Application will log:
```
2025-11-12T10:30:45.123Z INFO  auth: Authentication disabled (PEDANTIC_RAVEN_TOKEN not set)
```

### Future Enhancements

**Planned (not yet implemented)**:
- Multiple token support for key rotation
- Token expiration
- Rate limiting based on authentication failures
- Audit logging of authentication attempts
- TLS client certificate authentication
- OAuth2 integration

**Current Status**: Suitable for single-user, trusted environments on trusted networks.

---

### mnemosyne Configuration

**Connection settings** (via environment):
```bash
export MNEMOSYNE_ADDR=mnemosyne-server:50051
export MNEMOSYNE_TIMEOUT=45
export MNEMOSYNE_MAX_RETRIES=5
```

**Test connection**:
```bash
# Check if mnemosyne server is reachable
curl -v localhost:50051 2>&1 | head -20
# or use nc
nc -zv localhost 50051
```

### GLiNER Configuration

**Service endpoint**:
```bash
export GLINER_SERVICE_URL=http://localhost:8765
export GLINER_TIMEOUT=10
export GLINER_FALLBACK_TO_PATTERN=true
```

**Custom entity types** (via config.toml):
```toml
[gliner.entity_types]
default = ["person", "organization", "location"]
custom = [
    "api_endpoint",
    "database_name",
    "security_vulnerability",
    "performance_metric"
]
```

---

## Deployment Options

### Option 1: Local Development

**Start everything on localhost**:

```bash
# Terminal 1: Start mnemosyne server
cd ../mnemosyne
cargo run --bin mnemosyne-rpc

# Terminal 2: Start GLiNER service (optional)
cd pedantic_raven/services/gliner
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
uvicorn main:app --host 127.0.0.1 --port 8765

# Terminal 3: Run Pedantic Raven
cd pedantic_raven
go build -o pedantic_raven .
./pedantic_raven
```

**Verify all services**:
```bash
# Check mnemosyne
nc -zv localhost 50051 && echo "✓ mnemosyne OK"

# Check GLiNER (if running)
curl -s http://localhost:8765/health | jq . && echo "✓ GLiNER OK"
```

### Option 2: Docker Compose

**Full multi-container setup** with mnemosyne, GLiNER, and Pedantic Raven.

**Create docker-compose.yml**:
```yaml
version: '3.8'

services:
  # mnemosyne Memory Server
  mnemosyne:
    image: rand/mnemosyne:latest
    container_name: pedantic-raven-mnemosyne
    ports:
      - "50051:50051"
    volumes:
      - mnemosyne-data:/data
    environment:
      - MNEMOSYNE_HOST=0.0.0.0
      - MNEMOSYNE_PORT=50051
      - LOG_LEVEL=info
    healthcheck:
      test: ["CMD", "nc", "-zv", "localhost", "50051"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
    restart: unless-stopped
    networks:
      - pedantic-raven-net

  # GLiNER NER Service
  gliner:
    build:
      context: ./services/gliner
      dockerfile: Dockerfile
    container_name: pedantic-raven-gliner
    ports:
      - "8765:8765"
    environment:
      - LOG_LEVEL=info
    volumes:
      # Cache models across restarts
      - huggingface-cache:/root/.cache/huggingface
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8765/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 360s
    restart: unless-stopped
    networks:
      - pedantic-raven-net

  # Pedantic Raven TUI
  # Note: TUI apps should run on host terminal, not in container
  # This shows how to containerize if needed
  pedantic-raven:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: pedantic-raven-app
    depends_on:
      mnemosyne:
        condition: service_healthy
      gliner:
        condition: service_healthy
    environment:
      - MNEMOSYNE_ADDR=mnemosyne:50051
      - GLINER_SERVICE_URL=http://gliner:8765
      - LOG_LEVEL=info
    volumes:
      - ./workspace:/workspace
    stdin_open: true
    tty: true
    networks:
      - pedantic-raven-net

volumes:
  mnemosyne-data:
    driver: local
  huggingface-cache:
    driver: local

networks:
  pedantic-raven-net:
    driver: bridge
```

**Usage**:
```bash
# Start all services
docker-compose up -d

# View service logs
docker-compose logs -f mnemosyne
docker-compose logs -f gliner
docker-compose logs -f pedantic-raven

# Check service status
docker-compose ps

# Stop all services
docker-compose down

# Clean up everything (including data)
docker-compose down -v

# Rebuild images after code changes
docker-compose up -d --build

# Scale services (if applicable)
docker-compose up -d --scale gliner=2
```

**Advanced Docker Compose options**:

```yaml
# Add resource limits
services:
  gliner:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G

  mnemosyne:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G

# Add custom networks for isolation
networks:
  backend:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### Option 3: Systemd Service (Linux)

**Create service unit** (`/etc/systemd/system/pedantic-raven.service`):

```ini
[Unit]
Description=Pedantic Raven TUI Application
Documentation=https://github.com/rand/pedantic_raven
After=network.target mnemosyne.service gliner.service
Wants=mnemosyne.service gliner.service

[Service]
Type=simple
User=pedantic-raven
Group=pedantic-raven
WorkingDirectory=/opt/pedantic-raven
ExecStart=/opt/pedantic-raven/bin/pedantic_raven
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=pedantic-raven

# Environment variables
Environment="MNEMOSYNE_ADDR=localhost:50051"
Environment="GLINER_SERVICE_URL=http://localhost:8765"
Environment="GLINER_ENABLED=true"
Environment="LOG_LEVEL=info"

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/opt/pedantic-raven/workspace

# Resource limits
MemoryLimit=512M
CPUQuota=50%

[Install]
WantedBy=multi-user.target
```

**Create mnemosyne service** (`/etc/systemd/system/mnemosyne.service`):

```ini
[Unit]
Description=mnemosyne Memory Server
Documentation=https://github.com/rand/mnemosyne
After=network.target

[Service]
Type=simple
User=mnemosyne
Group=mnemosyne
WorkingDirectory=/opt/mnemosyne
ExecStart=/opt/mnemosyne/bin/mnemosyne-rpc
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=mnemosyne

Environment="MNEMOSYNE_HOST=127.0.0.1"
Environment="MNEMOSYNE_PORT=50051"

[Install]
WantedBy=multi-user.target
```

**Setup instructions**:

```bash
# Create user and group
sudo useradd -r -s /bin/false pedantic-raven
sudo useradd -r -s /bin/false mnemosyne

# Create application directories
sudo mkdir -p /opt/pedantic-raven/bin
sudo mkdir -p /opt/pedantic-raven/workspace
sudo mkdir -p /opt/mnemosyne/bin
sudo mkdir -p /var/lib/mnemosyne

# Copy binaries
sudo cp pedantic_raven /opt/pedantic-raven/bin/
sudo cp ../mnemosyne/target/release/mnemosyne-rpc /opt/mnemosyne/bin/

# Set permissions
sudo chown -R pedantic-raven:pedantic-raven /opt/pedantic-raven
sudo chown -R mnemosyne:mnemosyne /opt/mnemosyne
sudo chmod 755 /opt/pedantic-raven/bin/pedantic_raven
sudo chmod 755 /opt/mnemosyne/bin/mnemosyne-rpc

# Install service files
sudo cp pedantic-raven.service /etc/systemd/system/
sudo cp mnemosyne.service /etc/systemd/system/

# Enable and start services
sudo systemctl daemon-reload
sudo systemctl enable mnemosyne.service
sudo systemctl enable pedantic-raven.service
sudo systemctl start mnemosyne.service
sudo systemctl start pedantic-raven.service

# Check status
sudo systemctl status mnemosyne.service
sudo systemctl status pedantic-raven.service

# View logs
sudo journalctl -u pedantic-raven -f
sudo journalctl -u mnemosyne -f
```

### Option 4: Production Server Deployment

**Bare metal or VM deployment**:

```bash
# 1. System setup
sudo apt-get update
sudo apt-get install -y curl wget git

# 2. Install Go
wget https://go.dev/dl/go1.25.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version

# 3. Build from source
git clone https://github.com/rand/pedantic_raven.git
cd pedantic_raven
go build -ldflags="-s -w" -o pedantic_raven .

# 4. Create installation directory
sudo mkdir -p /opt/pedantic-raven/{bin,config,workspace}
sudo cp pedantic_raven /opt/pedantic-raven/bin/

# 5. Install mnemosyne
git clone https://github.com/rand/mnemosyne.git
cd mnemosyne
cargo build --release
sudo cp target/release/mnemosyne-rpc /opt/mnemosyne/bin/

# 6. Configure systemd services
sudo cp pedantic-raven.service /etc/systemd/system/
sudo cp mnemosyne.service /etc/systemd/system/

# 7. Start services
sudo systemctl daemon-reload
sudo systemctl enable mnemosyne pedantic-raven
sudo systemctl start mnemosyne pedantic-raven

# 8. Verify deployment
curl localhost:50051 2>&1 | head
pedantic_raven --help
```

**Production hardening checklist**:
- [ ] Run services with non-root users
- [ ] Configure firewall rules
- [ ] Set up log rotation
- [ ] Configure backup strategy
- [ ] Set resource limits (CPU, memory, disk)
- [ ] Enable SELinux/AppArmor policies
- [ ] Configure monitoring and alerting
- [ ] Document recovery procedures

---

## mnemosyne Setup

### Installing mnemosyne Server

**Option 1: Build from Source**:
```bash
# Clone mnemosyne repository
git clone https://github.com/rand/mnemosyne.git
cd mnemosyne

# Build (requires Rust 1.70+)
cargo build --release

# Binary location
./target/release/mnemosyne-rpc
```

**Option 2: Using Existing Binary**:
```bash
# Download from releases
wget https://github.com/rand/mnemosyne/releases/download/v1.0.0/mnemosyne-rpc-linux-amd64
chmod +x mnemosyne-rpc-linux-amd64

# Move to standard location
sudo mv mnemosyne-rpc-linux-amd64 /usr/local/bin/mnemosyne-rpc
```

### Configuration

**Default settings**:
```bash
# Start with defaults
mnemosyne-rpc

# Output should show
# mnemosyne RPC server listening on 127.0.0.1:50051
```

**Custom configuration** (environment variables):
```bash
# Custom host and port
export MNEMOSYNE_HOST=0.0.0.0
export MNEMOSYNE_PORT=5005
mnemosyne-rpc

# Data directory
export MNEMOSYNE_DATA_DIR=/var/lib/mnemosyne
mnemosyne-rpc

# Log level
export LOG_LEVEL=debug
mnemosyne-rpc
```

### Connection Testing

**Test from Pedantic Raven**:
```bash
# Pedantic Raven will attempt connection on startup
# Status shown in bottom status bar:
# [mnemosyne: ✓] - Connected (green)
# [mnemosyne: ✗] - Disconnected (red)

# Fallback behavior
# - If mnemosyne unavailable: Uses sample/demo mode
# - All memory operations work with sample data
# - No data persists to mnemosyne
```

**Test from command line**:
```bash
# Check if port is open
nc -zv localhost 50051
# Expected: Connection successful

# Check with gRPC
grpcurl -plaintext localhost:50051 mnemosyne.v1.HealthService/Check

# Check with simple TCP connection
timeout 5 bash -c 'echo > /dev/tcp/localhost/50051' && echo "OK"
```

### Troubleshooting Connection Issues

**Issue: Connection refused**
```bash
# 1. Verify mnemosyne is running
ps aux | grep mnemosyne-rpc

# 2. Check if port is listening
netstat -tlnp | grep 50051
lsof -i :50051

# 3. Verify firewall
sudo ufw status
sudo iptables -L

# 4. Restart mnemosyne
systemctl restart mnemosyne

# 5. Check logs
journalctl -u mnemosyne -n 50
```

**Issue: Timeout errors**
```bash
# Increase timeout
export MNEMOSYNE_TIMEOUT=60
./pedantic_raven

# Check network latency
ping -c 5 localhost

# Verify no resource issues on mnemosyne server
top -p $(pgrep mnemosyne-rpc)
```

**Issue: Connection succeeds but no data**
```bash
# Check data directory permissions
ls -la /var/lib/mnemosyne/

# Verify mnemosyne is accepting connections
grpcurl -plaintext localhost:50051 list

# Check mnemosyne logs
journalctl -u mnemosyne -f
```

---

## GLiNER Integration

### Docker Setup for GLiNER

**Using existing docker-compose.yml**:
```bash
# Start only GLiNER
docker-compose up -d gliner

# View startup logs (model download)
docker-compose logs -f gliner

# Check health
curl http://localhost:8765/health | jq .
```

**Manual Docker setup**:
```bash
# Build GLiNER image
cd services/gliner
docker build -t gliner-service:latest .

# Run container
docker run -d \
  --name gliner \
  -p 8765:8765 \
  -v huggingface-cache:/root/.cache/huggingface \
  -e LOG_LEVEL=info \
  gliner-service:latest

# View logs
docker logs -f gliner
```

### Configuration

**Environment variables**:
```bash
export GLINER_ENABLED=true
export GLINER_SERVICE_URL=http://localhost:8765
export GLINER_TIMEOUT=10
export GLINER_MAX_RETRIES=2
export GLINER_FALLBACK_TO_PATTERN=true
export GLINER_SCORE_THRESHOLD=0.5
```

**Custom entity types** (config.toml):
```toml
[gliner.entity_types]
default = ["person", "organization", "location", "technology", "concept", "product"]
custom = [
    "api_endpoint",
    "database_name",
    "security_vulnerability",
    "performance_metric",
    "compliance_requirement"
]
```

**GLiNER service configuration** (services/gliner/.env):
```bash
# Model configuration
MODEL_NAME=ner-large
MODEL_CACHE_DIR=/root/.cache/huggingface

# Server configuration
HOST=0.0.0.0
PORT=8765
LOG_LEVEL=info

# Performance tuning
BATCH_SIZE=32
MAX_WORKERS=4
```

### Testing Entity Extraction

**Health check**:
```bash
curl http://localhost:8765/health
# Expected response: {"status": "ok", "version": "1.0"}
```

**Test entity extraction**:
```bash
curl -X POST http://localhost:8765/extract \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Alice works at Google on machine learning",
    "entity_types": ["person", "organization", "technology"]
  }' | jq .

# Expected response:
# {
#   "entities": [
#     {"text": "Alice", "type": "person", "score": 0.95, "start": 0, "end": 5},
#     {"text": "Google", "type": "organization", "score": 0.98, "start": 15, "end": 21},
#     {"text": "machine learning", "type": "technology", "score": 0.92, "start": 25, "end": 40}
#   ]
# }
```

**Status check in Pedantic Raven**:
```
# Look at bottom status bar
[GLiNER: ✓]      # Green - service running
[GLiNER: ✗]      # Red - service unavailable (fallback active)
```

### Fallback Behavior

**When GLiNER is unavailable**:
1. Pedantic Raven attempts connection on startup
2. If unavailable, falls back to pattern-based extraction
3. Extracts 6 default entity types (person, organization, location, etc.)
4. Lower accuracy but continues to function
5. Periodically retries connection to GLiNER

**Configuration**:
```bash
# Disable fallback (require GLiNER)
export GLINER_FALLBACK_TO_PATTERN=false
./pedantic_raven

# With fallback enabled (default)
export GLINER_FALLBACK_TO_PATTERN=true
./pedantic_raven
```

---

## Monitoring & Logging

### Log Locations

**System logs** (if using systemd):
```bash
# Real-time logs
sudo journalctl -u pedantic-raven -f

# Last 100 lines
sudo journalctl -u pedantic-raven -n 100

# Last hour
sudo journalctl -u pedantic-raven --since "1 hour ago"

# mnemosyne logs
sudo journalctl -u mnemosyne -f
```

**Application logs** (if running directly):
```bash
# Logs written to stdout
./pedantic_raven 2>&1 | tee app.log

# Or redirect to file
./pedantic_raven > app.log 2>&1 &
```

**Docker logs**:
```bash
# View container logs
docker-compose logs pedantic-raven
docker-compose logs mnemosyne
docker-compose logs gliner

# Stream logs
docker-compose logs -f pedantic-raven

# Last N lines
docker-compose logs --tail 100 pedantic-raven
```

### Log Levels

**Configure log level**:
```bash
export LOG_LEVEL=debug   # Verbose output
export LOG_LEVEL=info    # Standard output (default)
export LOG_LEVEL=warn    # Warnings and errors only
export LOG_LEVEL=error   # Errors only
```

**Example log output** (INFO level):
```
2025-11-12T10:30:45.123Z INFO  app: Starting Pedantic Raven
2025-11-12T10:30:45.456Z INFO  mnemosyne: Connected to server at localhost:50051
2025-11-12T10:30:45.789Z INFO  gliner: Service available at http://localhost:8765
2025-11-12T10:30:46.012Z INFO  editor: File opened: /home/user/context.md
2025-11-12T10:30:48.345Z DEBUG semantic: Analyzed 1024 bytes, extracted 5 entities
```

### Health Check Endpoints

**Pedantic Raven** (checks dependencies):
```bash
# Access from within Pedantic Raven
:health

# Output shown in terminal
mnemosyne: ✓ Connected (latency: 5ms)
GLiNER: ✓ Available (latency: 120ms)
Memory: 45MB / 512MB (8%)
Uptime: 2h 30m
```

**mnemosyne gRPC health check**:
```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

**GLiNER HTTP health check**:
```bash
curl -s http://localhost:8765/health | jq .
```

### Debugging Tips

**Enable debug logging**:
```bash
export LOG_LEVEL=debug
./pedantic_raven
```

**Check resource usage**:
```bash
# Monitor Pedantic Raven
watch -n 1 'ps aux | grep pedantic_raven'

# Monitor mnemosyne
watch -n 1 'ps aux | grep mnemosyne-rpc'

# Monitor GLiNER container
docker stats gliner
```

**Network debugging**:
```bash
# Monitor connections
netstat -tlnp | grep -E '50051|8765'

# Trace network requests
tcpdump -i lo port 50051 or port 8765

# Test latency
ping -c 5 localhost
```

**Performance profiling**:
```bash
# Build with profiling
go build -o pedantic_raven_prof .

# Profile CPU
./pedantic_raven_prof -cpuprofile=cpu.prof
# Ctrl+C to stop
go tool pprof cpu.prof

# Profile memory
./pedantic_raven_prof -memprofile=mem.prof
go tool pprof mem.prof
```

---

## Troubleshooting

### Common Issues and Solutions

**Issue: "mnemosyne server unreachable"**

Solution:
```bash
# 1. Verify server is running
systemctl status mnemosyne

# 2. Check if port is listening
netstat -tlnp | grep 50051

# 3. Verify address and port
echo $MNEMOSYNE_ADDR

# 4. Test connection
nc -zv localhost 50051

# 5. Restart server
systemctl restart mnemosyne
```

**Issue: "GLiNER service unavailable"**

Solution:
```bash
# 1. Verify GLiNER is running
docker-compose ps gliner

# 2. Check container logs
docker-compose logs gliner

# 3. Test endpoint
curl http://localhost:8765/health

# 4. Restart container
docker-compose restart gliner

# 5. Check if fallback is working
# Should see "[GLiNER: ✗]" in status bar but app still works
```

**Issue: "Terminal too small - compact layout"**

Solution:
```bash
# Resize terminal to at least 120x30
# Check current size
echo "Cols: $(tput cols), Rows: $(tput lines)"

# Resize terminal or drag window edge
# Pedantic Raven auto-detects and switches to normal layout
```

**Issue: High memory usage / slow performance**

Solution:
```bash
# 1. Check memory usage
top -p $(pgrep pedantic_raven)

# 2. Close other applications
kill $(pgrep "slack\|chrome\|docker")

# 3. Reduce file size (limit editing to <10MB files)
wc -l large_file.txt

# 4. Disable GLiNER if not needed
export GLINER_ENABLED=false
./pedantic_raven

# 5. Restart with fresh state
killall pedantic_raven
./pedantic_raven
```

**Issue: Semantic analysis not running**

Solution:
```bash
# 1. Check editor mode
# Press '1' to ensure in Edit mode

# 2. Verify semantic analyzer is active
# Type some text and wait 500ms for results in context panel

# 3. Check if GLiNER is needed
# If GLiNER unavailable, should still use pattern matcher

# 4. Check logs
export LOG_LEVEL=debug
./pedantic_raven

# 5. Try simpler text
# Some patterns may not trigger extraction
```

**Issue: Can't save files**

Solution:
```bash
# 1. Verify directory exists and is writable
ls -la /path/to/workspace

# 2. Check file permissions
chmod 755 /path/to/workspace

# 3. Try saving to different location
# Press Ctrl+S and choose different path

# 4. Check disk space
df -h

# 5. Verify user has write permissions
id
groups
```

**Issue: Connection timeout**

Solution:
```bash
# 1. Increase timeout
export MNEMOSYNE_TIMEOUT=60
export GLINER_TIMEOUT=15
./pedantic_raven

# 2. Check network latency
ping -c 5 localhost

# 3. Check server load
top -p $(pgrep mnemosyne-rpc)

# 4. Verify no firewall issues
sudo ufw allow 50051/tcp
sudo ufw allow 8765/tcp

# 5. Check for network errors
netstat -i
dmesg | tail
```

---

## Security Considerations

### File Permissions

**Application directory**:
```bash
# Restrict access to application directory
sudo chmod 755 /opt/pedantic-raven
sudo chown -R pedantic-raven:pedantic-raven /opt/pedantic-raven

# Workspace directory (user-writable)
sudo chmod 775 /opt/pedantic-raven/workspace

# Data directory (only root/service user)
sudo chmod 700 /var/lib/mnemosyne
sudo chown mnemosyne:mnemosyne /var/lib/mnemosyne
```

**Log files**:
```bash
# Restrict log access
sudo chmod 600 /var/log/pedantic-raven.log
sudo chown pedantic-raven:pedantic-raven /var/log/pedantic-raven.log
```

### Network Security

**Firewall configuration**:
```bash
# UFW (Ubuntu)
sudo ufw allow 50051/tcp from 127.0.0.1  # mnemosyne (localhost only)
sudo ufw allow 8765/tcp from 127.0.0.1   # GLiNER (localhost only)

# iptables
sudo iptables -A INPUT -p tcp --dport 50051 -s 127.0.0.1 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 8765 -s 127.0.0.1 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 50051 -j DROP
sudo iptables -A INPUT -p tcp --dport 8765 -j DROP
```

**SSH tunneling** (for remote access):
```bash
# Create secure tunnel to remote mnemosyne
ssh -L 50051:127.0.0.1:50051 user@remote-server

# Now connect locally
export MNEMOSYNE_ADDR=localhost:50051
./pedantic_raven
```

### Future Authentication

**Planned security features** (not yet implemented):
- TLS/SSL for gRPC connections
- API key authentication for mnemosyne
- Session-based access control
- Audit logging

**Current status**: Runs in trusted environment (localhost or private network).

### Best Practices

**Security checklist**:
- [ ] Run services on non-root user accounts
- [ ] Restrict file permissions to minimum necessary
- [ ] Configure firewall to limit access to localhost
- [ ] Use SSH tunneling for remote connections
- [ ] Disable mnemosyne/GLiNER when not needed
- [ ] Regularly update dependencies
- [ ] Monitor logs for suspicious activity
- [ ] Back up sensitive memory data

---

## Backup & Recovery

### Data Locations

**mnemosyne data**:
```
/var/lib/mnemosyne/          # Primary data directory
  └── data/                  # Memory database
      ├── vectors/           # Embedding vectors
      ├── graph/             # Memory relationships
      └── metadata/          # Memory metadata
```

**Pedantic Raven workspace**:
```
/opt/pedantic-raven/workspace/  # User-created files
```

**Docker volumes**:
```
mnemosyne-data             # mnemosyne database volume
huggingface-cache          # GLiNER model cache
```

### Backup Procedures

**Manual backup**:
```bash
# Backup mnemosyne data
sudo tar -czf mnemosyne-backup-$(date +%Y%m%d).tar.gz \
  /var/lib/mnemosyne/

# Backup workspace
tar -czf workspace-backup-$(date +%Y%m%d).tar.gz \
  /opt/pedantic-raven/workspace/

# Verify backups
ls -lh *-backup-*.tar.gz
tar -tzf mnemosyne-backup-*.tar.gz | head
```

**Automated backup** (cron job):
```bash
# Create backup script
sudo cat > /usr/local/bin/backup-pedantic-raven.sh << 'EOF'
#!/bin/bash
set -e

BACKUP_DIR="/backups/pedantic-raven"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Backup mnemosyne
echo "Backing up mnemosyne..."
tar -czf "$BACKUP_DIR/mnemosyne_$TIMESTAMP.tar.gz" \
  /var/lib/mnemosyne/ 2>/dev/null || true

# Backup workspace
echo "Backing up workspace..."
tar -czf "$BACKUP_DIR/workspace_$TIMESTAMP.tar.gz" \
  /opt/pedantic-raven/workspace/ 2>/dev/null || true

# Keep only last 7 days
find "$BACKUP_DIR" -mtime +7 -delete

echo "Backup complete: $BACKUP_DIR"
EOF

sudo chmod +x /usr/local/bin/backup-pedantic-raven.sh

# Schedule daily at 2 AM
sudo crontab -e
# Add: 0 2 * * * /usr/local/bin/backup-pedantic-raven.sh
```

**Docker volume backup**:
```bash
# Backup Docker volumes
docker run --rm -v mnemosyne-data:/data -v $(pwd):/backup \
  alpine tar czf /backup/mnemosyne-docker-backup.tar.gz -C /data .

# Backup HuggingFace cache
docker run --rm -v huggingface-cache:/cache -v $(pwd):/backup \
  alpine tar czf /backup/huggingface-backup.tar.gz -C /cache .
```

### Recovery Procedures

**Restore from backup**:
```bash
# Stop services
sudo systemctl stop pedantic-raven mnemosyne

# Clear existing data
sudo rm -rf /var/lib/mnemosyne/data/*

# Restore from backup
sudo tar -xzf mnemosyne-backup-20251112.tar.gz -C /

# Restore workspace
tar -xzf workspace-backup-20251112.tar.gz -C /opt/pedantic-raven/

# Fix permissions
sudo chown -R mnemosyne:mnemosyne /var/lib/mnemosyne/
sudo chown -R pedantic-raven:pedantic-raven /opt/pedantic-raven/

# Start services
sudo systemctl start mnemosyne
sudo systemctl start pedantic-raven

# Verify
sudo systemctl status mnemosyne pedantic-raven
```

**Restore Docker volume**:
```bash
# Stop containers
docker-compose down

# Remove volume
docker volume rm mnemosyne-data

# Restore from backup
docker run --rm -v mnemosyne-data:/data -v $(pwd):/backup \
  alpine tar xzf /backup/mnemosyne-docker-backup.tar.gz -C /data

# Start containers
docker-compose up -d

# Verify
docker-compose logs mnemosyne
```

---

## Updating

### Update Procedures

**Before updating**:
1. Create backup of mnemosyne data
2. Document current version
3. Review changelog for breaking changes

```bash
# Check current version
./pedantic_raven --version

# View changelog
cat docs/CHANGELOG.md

# Backup data
sudo systemctl stop pedantic-raven
tar -czf mnemosyne-backup-$(date +%Y%m%d).tar.gz /var/lib/mnemosyne/
```

**Update from source**:
```bash
# Pull latest code
git fetch origin
git checkout main

# Build new version
go mod download
go build -o pedantic_raven .

# Run tests to verify
go test ./...

# Stop old version
sudo systemctl stop pedantic-raven

# Install new binary
sudo cp pedantic_raven /opt/pedantic-raven/bin/
sudo chown pedantic-raven:pedantic-raven /opt/pedantic-raven/bin/pedantic_raven

# Start new version
sudo systemctl start pedantic-raven

# Verify
sudo systemctl status pedantic-raven
```

**Update Docker images**:
```bash
# Pull latest images
docker-compose pull

# Rebuild if needed
docker-compose build --no-cache

# Start with new images
docker-compose up -d

# Verify versions
docker-compose exec pedantic-raven pedantic_raven --version
```

**Update mnemosyne**:
```bash
# Same procedure as main application
# Stop service, backup data, build/pull new version, start

# Or use pre-built binaries
wget https://github.com/rand/mnemosyne/releases/download/v1.1.0/mnemosyne-rpc-linux-amd64
chmod +x mnemosyne-rpc-linux-amd64
sudo mv mnemosyne-rpc-linux-amd64 /usr/local/bin/mnemosyne-rpc

# Verify
mnemosyne-rpc --version
```

### Version Compatibility

**Go dependencies**:
```bash
# Check current versions
go mod graph

# Update dependencies
go get -u ./...

# Update specific package
go get -u github.com/charmbracelet/bubbletea@latest

# Verify compatibility
go mod tidy
go test ./...
```

**mnemosyne API compatibility**:
- Pedantic Raven requires mnemosyne v1.0+
- Protocol buffer definitions in `proto/mnemosyne/v1/`
- Breaking changes documented in CHANGELOG

**Supported versions**:
| Pedantic Raven | mnemosyne | Go | Status |
|---|---|---|---|
| v1.5+ | v1.0+ | 1.21+ | Current |
| v1.3-1.4 | v0.9 | 1.20+ | Deprecated |
| v1.0-1.2 | v0.8 | 1.20+ | End of Life |

### Migration Guides

**Migrating from v1.0 to v1.1**:
```bash
# No database migration needed
# Just update binary and restart
# New entity types automatically supported

# Update config.toml if using custom entity types
# See Configuration section for details
```

**Migrating data between servers**:
```bash
# 1. Backup source mnemosyne
tar -czf mnemosyne-export.tar.gz /var/lib/mnemosyne/

# 2. Transfer to new server
scp mnemosyne-export.tar.gz user@new-server:/tmp/

# 3. Restore on new server
ssh user@new-server
sudo tar -xzf /tmp/mnemosyne-export.tar.gz -C /

# 4. Adjust ownership
sudo chown -R mnemosyne:mnemosyne /var/lib/mnemosyne/

# 5. Restart mnemosyne
sudo systemctl restart mnemosyne
```

---

## Quick Reference

### Common Commands

```bash
# Build
go build -o pedantic_raven .

# Run
./pedantic_raven

# Test
go test ./...

# Check version
./pedantic_raven --version

# Start systemd services
sudo systemctl start mnemosyne pedantic-raven

# View logs
sudo journalctl -u pedantic-raven -f

# Docker operations
docker-compose up -d
docker-compose logs -f
docker-compose down

# Backup
tar -czf backup.tar.gz /var/lib/mnemosyne/
```

### Environment Variables

```bash
# mnemosyne
export MNEMOSYNE_ADDR=localhost:50051
export MNEMOSYNE_TIMEOUT=30
export MNEMOSYNE_MAX_RETRIES=3

# GLiNER
export GLINER_ENABLED=true
export GLINER_SERVICE_URL=http://localhost:8765
export GLINER_TIMEOUT=5
export GLINER_FALLBACK_TO_PATTERN=true

# Logging
export LOG_LEVEL=info
```

### Useful Files

- `/opt/pedantic-raven/bin/pedantic_raven` - Application binary
- `/etc/systemd/system/pedantic-raven.service` - Systemd service
- `/var/lib/mnemosyne/` - mnemosyne data directory
- `/opt/pedantic-raven/workspace/` - User files
- `config.toml` - Configuration file

---

**Last Updated**: November 12, 2025
**Status**: Complete
**Next Review**: February 2026
