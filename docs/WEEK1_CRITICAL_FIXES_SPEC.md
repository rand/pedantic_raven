# Week 1: Critical Fixes - Specification

**Date**: 2025-11-12
**Phase**: Production Readiness - Week 1
**Timeline**: 5 days
**Dependencies**: Phase 9 complete, critical review complete

## Overview

Week 1 focuses on fixing **critical blockers** preventing production deployment:
1. Race condition in semantic analyzer (crashes)
2. Test failures (confidence issue)
3. Deployment documentation (operational readiness)
4. Authentication planning (security readiness)

## Objectives

1. **Fix Critical Race Condition**: Resolve "close of closed channel" panic
2. **Fix Test Failures**: All tests passing
3. **Deployment Documentation**: Complete deployment guide
4. **Authentication Design**: Plan authentication approach (discussion with user)

## Success Criteria

- [ ] Zero panics/crashes in concurrent operations
- [ ] 100% test pass rate (or documented skips)
- [ ] Deployment guide complete (400+ lines)
- [ ] Authentication approach agreed upon with user
- [ ] All changes committed and tested

---

## Work Stream 1: Semantic Analyzer Race Condition Fix

### Problem Statement

**File**: `internal/editor/semantic/analyzer.go:112`
**Error**: "panic: close of closed channel"
**Impact**: Application crashes during concurrent mode switching

### Root Cause

```go
// Current buggy code (line 107-113)
defer func() {
    a.mu.Lock()
    a.running = false
    a.analysis.Duration = time.Since(startTime)
    a.mu.Unlock()
    close(updateChan)  // ← Can be called multiple times
}()
```

**Race Scenario**:
1. Goroutine A: `Analyze()` creates Chan1, captures reference
2. Goroutine B: `Analyze()` called rapidly, creates Chan2
3. Goroutine A: Cancellation triggers defer
4. Goroutine A: Attempts to close Chan1
5. **Panic**: Chan1 already closed OR channel reference stale

### Solution Design

**Option 1: sync.Once (Recommended)**
```go
type StreamingAnalyzer struct {
    mu            sync.RWMutex
    tokenizer     *Tokenizer
    extractor     Extractor
    running       bool
    cancel        context.CancelFunc
    updateChan    chan AnalysisUpdate
    channelClosed *sync.Once  // ← NEW: Ensure close called once
    analysis      *Analysis
}

// In Analyze():
a.channelClosed = &sync.Once{}  // Reset for new channel

// In defer:
a.channelClosed.Do(func() {
    close(updateChan)
})
```

**Option 2: Panic Recovery (Fallback)**
```go
defer func() {
    a.mu.Lock()
    a.running = false
    a.analysis.Duration = time.Since(startTime)
    a.mu.Unlock()

    // Safe close with recovery
    func() {
        defer func() {
            if r := recover(); r != nil {
                // Channel already closed, ignore
            }
        }()
        close(updateChan)
    }()
}()
```

**Option 3: Channel Tracking (Most Robust)**
```go
type StreamingAnalyzer struct {
    // ... existing fields
    currentChannelID uint64  // Track which channel is active
    channelMu        sync.Mutex
}

// In Analyze():
a.channelMu.Lock()
channelID := atomic.AddUint64(&a.currentChannelID, 1)
a.channelMu.Unlock()

// In defer:
a.channelMu.Lock()
if channelID == a.currentChannelID {
    close(updateChan)
}
a.channelMu.Unlock()
```

**Chosen Approach**: Option 1 (sync.Once) - Simple, thread-safe, idiomatic

### Implementation Steps

1. Add `channelClosed *sync.Once` field to `StreamingAnalyzer`
2. Initialize `channelClosed` in `Analyze()` before creating channel
3. Replace `close(updateChan)` with `a.channelClosed.Do(func() { close(updateChan) })`
4. Add tests for concurrent `Analyze()` calls
5. Run with race detector: `go test -race`

### Testing Strategy

**New Tests**:
1. `TestAnalyzerConcurrentAnalyze` - Multiple rapid calls
2. `TestAnalyzerConcurrentStop` - Analyze + Stop concurrently
3. `TestAnalyzerChannelCloseOnce` - Verify channel closed exactly once
4. `TestAnalyzerRaceDetector` - Run with -race flag

**Existing Tests**:
- Verify `TestConcurrentModeSwitching` passes
- Verify all semantic analyzer tests pass
- Verify Edit Mode integration tests pass

### Verification

```bash
# Run tests with race detector
go test -race ./internal/editor/semantic/...
go test -race ./internal/integration -run TestConcurrent

# Stress test
for i in {1..100}; do
  go test ./internal/integration -run TestConcurrentModeSwitching || break
done
```

---

## Work Stream 2: Test Failure Fixes

### 2a. mnemosyne Error Test Failures

**File**: `internal/mnemosyne/errors_test.go`

**Failure 1: TestWrapErrorDeadlineExceeded** (line 99)
```go
// Current assertion (WRONG)
if !strings.Contains(errMsg, "timed out") &&
   !strings.Contains(errMsg, "deadline exceeded") {
    t.Errorf("Expected error to contain 'timed out' or 'deadline exceeded', got: %s", errMsg)
}

// Fix: Accept "timeout" as valid
if !strings.Contains(strings.ToLower(errMsg), "timeout") &&
   !strings.Contains(strings.ToLower(errMsg), "timed out") &&
   !strings.Contains(strings.ToLower(errMsg), "deadline") {
    t.Errorf("Expected error to contain timeout indicator, got: %s", errMsg)
}
```

**Failure 2: TestCategorizeServerError/ErrUnavailable** (line 262)
```go
// Current expectation (WRONG)
want := server

// Fix: ErrUnavailable is categorized as 'connection', which is correct
want := connection  // Not server

// Reasoning: "server unavailable" is a connection problem, not a server error
```

### 2b. GLiNER Docker Test

**File**: `internal/editor/gliner_e2e_test.go`
**Test**: `TestGLiNERE2E_A2_DockerServiceLifecycle`

**Issue**: Docker service doesn't respond (environmental)

**Solution**: Skip test if Docker unavailable
```go
func TestGLiNERE2E_A2_DockerServiceLifecycle(t *testing.T) {
    // Check Docker availability
    if !isDockerAvailable(t) {
        t.Skip("Docker not available or GLiNER service unhealthy")
    }

    // ... existing test code
}

func isDockerAvailable(t *testing.T) bool {
    cmd := exec.Command("docker", "info")
    if err := cmd.Run(); err != nil {
        return false
    }

    // Try to ping GLiNER service
    client := &http.Client{Timeout: 2 * time.Second}
    resp, err := client.Get("http://localhost:8001/health")
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    return resp.StatusCode == http.StatusOK
}
```

### Testing Strategy

**Verification**:
```bash
# Run all tests
go test ./internal/mnemosyne -v
go test ./internal/editor -v
go test ./internal/integration -v

# Verify 100% pass rate
go test ./... 2>&1 | grep "^FAIL" | wc -l  # Should be 0
```

---

## Work Stream 3: Deployment Guide

### Objectives

Create comprehensive deployment guide for production deployment.

### Sections Required

1. **Prerequisites** (50 lines)
   - Go 1.21+
   - Docker (optional, for GLiNER)
   - mnemosyne server
   - Environment variables

2. **Installation** (100 lines)
   - Clone repository
   - Build from source
   - Binary installation
   - Docker installation

3. **Configuration** (100 lines)
   - Environment variables reference
   - Config file format
   - mnemosyne connection settings
   - GLiNER integration settings

4. **Deployment Options** (150 lines)
   - **Local Development**: Running on localhost
   - **Docker Compose**: Multi-container setup
   - **Systemd Service**: Linux system service
   - **Production Server**: Bare metal or VM

5. **mnemosyne Setup** (80 lines)
   - Installing mnemosyne server
   - Configuration
   - Connection testing
   - Troubleshooting

6. **Monitoring & Logging** (60 lines)
   - Log locations
   - Log levels
   - Monitoring endpoints (future)
   - Health checks

7. **Troubleshooting** (80 lines)
   - Common issues
   - Connection problems
   - GLiNER issues
   - Performance issues

8. **Security Considerations** (40 lines)
   - File permissions
   - Network security
   - Authentication (coming in Week 1)
   - Best practices

9. **Backup & Recovery** (40 lines)
   - Data locations
   - Backup procedures
   - Recovery procedures

10. **Updating** (30 lines)
    - Update procedures
    - Version compatibility
    - Migration guides

**Total**: 730 lines (exceeds 400+ target)

### Example Sections

**Docker Compose Example**:
```yaml
version: '3.8'

services:
  mnemosyne:
    image: mnemosyne:latest
    ports:
      - "50051:50051"
    volumes:
      - mnemosyne-data:/data
    environment:
      - MNEMOSYNE_HOST=0.0.0.0
      - MNEMOSYNE_PORT=50051

  gliner:
    image: gliner-service:latest
    ports:
      - "8001:8001"
    environment:
      - MODEL_PATH=/models/gliner-large

  pedantic-raven:
    build: .
    depends_on:
      - mnemosyne
      - gliner
    environment:
      - MNEMOSYNE_ADDR=mnemosyne:50051
      - GLINER_ENDPOINT=http://gliner:8001
    volumes:
      - ./workspace:/workspace

volumes:
  mnemosyne-data:
```

**Systemd Service Example**:
```ini
[Unit]
Description=Pedantic Raven TUI
After=network.target mnemosyne.service

[Service]
Type=simple
User=pedantic-raven
WorkingDirectory=/opt/pedantic-raven
ExecStart=/opt/pedantic-raven/bin/pedantic_raven
Restart=on-failure
RestartSec=10

Environment="MNEMOSYNE_ADDR=localhost:50051"
Environment="GLINER_ENDPOINT=http://localhost:8001"

[Install]
WantedBy=multi-user.target
```

---

## Work Stream 4: Authentication Planning

### Objectives

1. Discuss authentication approach with user
2. Decide on authentication method
3. Plan implementation for Week 2

### Options to Present

**Option 1: Environment Variable Token**
- **Pros**: Simple, no dependencies
- **Cons**: No user management, shared token
- **Use Case**: Single user or trusted environment

```bash
PEDANTIC_RAVEN_TOKEN=secret-token-here
```

**Option 2: mnemosyne-based Authentication**
- **Pros**: Centralized auth, user management
- **Cons**: Depends on mnemosyne server
- **Use Case**: Multi-user deployment

**Option 3: Local Password File**
- **Pros**: No external dependencies
- **Cons**: Manual user management
- **Use Case**: Small teams

**Option 4: OAuth2/OIDC**
- **Pros**: Enterprise-grade, SSO support
- **Cons**: Complex setup, overkill for CLI
- **Use Case**: Large organizations

**Recommendation**: Start with Option 1 (env token) for Week 1, plan Option 2 (mnemosyne-based) for future.

### Questions for User

1. **Target Deployment**: Single user or multi-user?
2. **Authentication Complexity**: Simple token or full user management?
3. **Integration**: Use mnemosyne for auth or separate system?
4. **Timeline**: Quick fix for Week 1 or comprehensive solution?

---

## Parallelization Strategy

### Safe Parallel Execution

**Agent 1 (Sonnet)**: Semantic Analyzer Fix (Stream 1)
- Complex race condition
- Needs careful analysis
- Critical for production
- **Estimated**: 4-6 hours

**Agent 2 (Haiku)**: Test Failure Fixes (Stream 2)
- Simple assertion updates
- Docker availability check
- Low risk
- **Estimated**: 2-3 hours

**Agent 3 (Haiku)**: Deployment Guide (Stream 3)
- Documentation only
- No code changes
- Zero conflict risk
- **Estimated**: 4-6 hours

**User Discussion**: Authentication Planning (Stream 4)
- Design decision required
- User input needed
- Cannot parallelize

### Execution Order

**Phase 1** (Parallel):
- Launch Agent 1 (Semantic analyzer fix)
- Launch Agent 2 (Test fixes)
- Launch Agent 3 (Deployment guide)

**Phase 2** (After Phase 1):
- Discuss authentication approach with user
- Plan implementation for Week 2

**Phase 3** (Integration):
- Merge all changes
- Run full test suite
- Verify all fixes
- Commit and push

---

## Risk Assessment

### Low Risk

- Test failure fixes (simple assertions)
- Deployment guide (documentation only)

### Medium Risk

- Semantic analyzer fix (critical but well-understood)
- Authentication planning (depends on user requirements)

### Mitigation

- Test with race detector after semantic analyzer fix
- Stress test concurrent operations
- Review authentication options thoroughly with user

---

## Success Metrics

**After Week 1**:
- [ ] Zero panics in concurrent operations
- [ ] 100% test pass rate
- [ ] Deployment guide complete (730+ lines)
- [ ] Authentication approach decided
- [ ] All changes committed and pushed
- [ ] CI/CD passing (if configured)

**Quality Gates**:
- [ ] All tests pass with `-race` flag
- [ ] No new compiler warnings
- [ ] No new lint warnings
- [ ] Documentation reviewed and accurate

---

## Deliverables

**Code Changes**:
1. `internal/editor/semantic/analyzer.go` (race condition fix)
2. `internal/mnemosyne/errors_test.go` (test assertions)
3. `internal/editor/gliner_e2e_test.go` (Docker availability check)
4. New tests for concurrent analyzer usage

**Documentation**:
5. `docs/DEPLOYMENT.md` (730+ lines)
6. `docs/WEEK1_COMPLETION.md` (completion report)

**Planning**:
7. Authentication approach documented
8. Week 2 plan ready

---

**Last Updated**: 2025-11-12
**Phase**: Week 1 Critical Fixes
**Status**: Ready for Execution
