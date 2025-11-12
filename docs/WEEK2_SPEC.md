# Week 2: Authentication & Quality - Specification

**Date**: 2025-11-12
**Phase**: Production Readiness - Week 2
**Timeline**: 5 days
**Dependencies**: Week 1 complete (race condition fixed, tests passing, deployment guide ready)

## Overview

Week 2 focuses on **authentication, test coverage, and critical TODOs** to reach 95% production readiness:
1. Simple token authentication (single-user deployment)
2. mnemosyne test coverage improvement (59.9% → 75%)
3. Critical TODO resolution (3 high-priority items)
4. Security hardening

## Objectives

1. **Simple Token Authentication**: Environment variable-based auth for single-user deployment
2. **Test Coverage**: Improve mnemosyne package from 59.9% to 75%+ (add 15-20 tests)
3. **Critical TODOs**: Address 3 critical items blocking production
4. **Security**: Rate limiting, input sanitization, audit logging

## Success Criteria

- [ ] Authentication working (token validation)
- [ ] mnemosyne coverage ≥ 75%
- [ ] 3 critical TODOs resolved
- [ ] Security features implemented
- [ ] Documentation updated
- [ ] All tests passing (100% pass rate maintained)

---

## Work Stream 1: Simple Token Authentication

### Requirements

**Deployment Context**: Single-user personal use
**Authentication Method**: Environment variable token
**Security Level**: Basic (sufficient for single-user)

### Design

**Environment Variable**:
```bash
export PEDANTIC_RAVEN_TOKEN="your-secret-token-here"
```

**Token Generation**:
```bash
# Generate secure random token (32 bytes)
openssl rand -base64 32
```

**Token Validation**:
- On application startup, check if `PEDANTIC_RAVEN_TOKEN` is set
- If set, enable authentication mode
- If not set, run in "open mode" (backward compatible)
- Token stored in memory, never logged

### Implementation

**1. Create auth package** (`internal/auth/`)

**File**: `internal/auth/token.go`
```go
package auth

import (
    "crypto/subtle"
    "os"
)

type TokenAuth struct {
    token string
    enabled bool
}

func NewTokenAuth() *TokenAuth {
    token := os.Getenv("PEDANTIC_RAVEN_TOKEN")
    return &TokenAuth{
        token: token,
        enabled: token != "",
    }
}

func (a *TokenAuth) IsEnabled() bool {
    return a.enabled
}

func (a *TokenAuth) Validate(providedToken string) bool {
    if !a.enabled {
        return true // Auth disabled, allow all
    }

    // Use constant-time comparison to prevent timing attacks
    return subtle.ConstantTimeCompare(
        []byte(a.token),
        []byte(providedToken),
    ) == 1
}
```

**File**: `internal/auth/token_test.go`
```go
package auth

import (
    "os"
    "testing"
)

func TestTokenAuthDisabled(t *testing.T) {
    os.Unsetenv("PEDANTIC_RAVEN_TOKEN")
    auth := NewTokenAuth()

    if auth.IsEnabled() {
        t.Error("Expected auth to be disabled when PEDANTIC_RAVEN_TOKEN not set")
    }

    if !auth.Validate("any-token") {
        t.Error("Expected validation to pass when auth disabled")
    }
}

func TestTokenAuthEnabled(t *testing.T) {
    token := "test-secret-token"
    os.Setenv("PEDANTIC_RAVEN_TOKEN", token)
    defer os.Unsetenv("PEDANTIC_RAVEN_TOKEN")

    auth := NewTokenAuth()

    if !auth.IsEnabled() {
        t.Error("Expected auth to be enabled")
    }

    if !auth.Validate(token) {
        t.Error("Expected valid token to pass")
    }

    if auth.Validate("wrong-token") {
        t.Error("Expected invalid token to fail")
    }
}

func TestTokenAuthConstantTime(t *testing.T) {
    token := "correct-token"
    os.Setenv("PEDANTIC_RAVEN_TOKEN", token)
    defer os.Unsetenv("PEDANTIC_RAVEN_TOKEN")

    auth := NewTokenAuth()

    // Test that validation uses constant-time comparison
    // (timing attacks should not reveal token length/content)
    testCases := []struct {
        name string
        input string
        want bool
    }{
        {"Exact match", token, true},
        {"Wrong token", "wrong-token", false},
        {"Shorter token", "short", false},
        {"Longer token", token + "extra", false},
        {"Empty token", "", false},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got := auth.Validate(tc.input)
            if got != tc.want {
                t.Errorf("Validate(%q) = %v, want %v", tc.input, got, tc.want)
            }
        })
    }
}
```

**2. Integrate with application startup**

**File**: `cmd/pedantic_raven/main.go` (modifications)
```go
import (
    "github.com/rand/pedantic-raven/internal/auth"
    // ... other imports
)

func main() {
    // Initialize authentication
    tokenAuth := auth.NewTokenAuth()

    if tokenAuth.IsEnabled() {
        log.Info("Authentication enabled (PEDANTIC_RAVEN_TOKEN set)")
    } else {
        log.Warn("Authentication disabled (PEDANTIC_RAVEN_TOKEN not set)")
    }

    // Pass tokenAuth to application
    app := NewApp(tokenAuth)
    app.Run()
}
```

**3. Add token prompt on startup** (if enabled)

For TUI application, prompt for token on startup if auth is enabled:

```go
func (a *App) promptForToken() (string, error) {
    // Use bubble tea to show password input
    prompt := textinput.New()
    prompt.Placeholder = "Enter token"
    prompt.EchoMode = textinput.EchoPassword
    prompt.Focus()

    // ... bubble tea program to get input
    return token, nil
}
```

**4. Update documentation**

Update `docs/DEPLOYMENT.md` and `README.md`:
- Add "Authentication" section
- Document `PEDANTIC_RAVEN_TOKEN` environment variable
- Provide token generation command
- Explain single-user deployment model

### Testing

**Unit Tests** (6 tests):
1. `TestTokenAuthDisabled` - Auth disabled when env var not set
2. `TestTokenAuthEnabled` - Auth enabled when env var set
3. `TestTokenAuthConstantTime` - Constant-time comparison
4. `TestTokenAuthValidation` - Various token inputs
5. `TestTokenAuthEmptyToken` - Empty string handling
6. `TestTokenAuthLongToken` - Long token handling

**Integration Tests** (2 tests):
1. `TestAppStartupWithAuth` - Application starts with auth enabled
2. `TestAppStartupWithoutAuth` - Application starts with auth disabled

**Security Tests** (2 tests):
1. `TestTimingAttackResistance` - Verify constant-time comparison
2. `TestTokenNotLogged` - Verify token never appears in logs

---

## Work Stream 2: mnemosyne Test Coverage Improvement

### Current State

**Coverage**: 59.9% (target: 75%)
**Gap**: +15.1% (need ~15-20 new tests)

### Files to Test

Priority order based on importance and current coverage gaps:

1. **client.go** (high priority)
   - Connection management
   - RPC method wrappers
   - Error handling

2. **connection_pool.go** (high priority)
   - Pool lifecycle
   - Connection health checks
   - Concurrent access

3. **recall.go** (medium priority)
   - Semantic search operations
   - Query building
   - Result handling

4. **remember.go** (medium priority)
   - Memory storage operations
   - Namespace handling
   - Validation

5. **evolve.go** (low priority)
   - Memory evolution
   - Consolidation logic
   - Archival

### New Tests Required (18 tests)

**client.go Tests** (6 tests):
1. `TestClientConnectRetry` - Connection retry with backoff
2. `TestClientConnectTimeout` - Connection timeout handling
3. `TestClientDisconnectCleanup` - Resource cleanup on disconnect
4. `TestClientConcurrentRequests` - Multiple concurrent RPC calls
5. `TestClientContextCancellation` - Request cancellation
6. `TestClientReconnectAfterFailure` - Automatic reconnection

**connection_pool.go Tests** (5 tests):
1. `TestPoolExhaustion` - Behavior when pool exhausted
2. `TestPoolHealthCheck` - Health check removes bad connections
3. `TestPoolConcurrentAcquire` - Concurrent connection acquisition
4. `TestPoolGracefulShutdown` - Proper pool shutdown
5. `TestPoolConnectionReuse` - Connections properly reused

**recall.go Tests** (3 tests):
1. `TestRecallWithNamespace` - Namespace filtering
2. `TestRecallPagination` - Cursor-based pagination
3. `TestRecallEmptyResults` - No results handling

**remember.go Tests** (2 tests):
1. `TestRememberValidation` - Input validation
2. `TestRememberDuplicateHandling` - Duplicate memory handling

**evolve.go Tests** (2 tests):
1. `TestEvolveMemoryDecay` - Memory importance decay
2. `TestEvolveArchival` - Old memory archival

### Testing Strategy

**Mock Server Approach**:
- Create mock gRPC server for testing
- No external dependencies
- Fast, deterministic tests

**Example Mock Server**:
```go
type mockMnemosyneServer struct {
    pb.UnimplementedMemoryServiceServer
    memories map[string]*pb.MemoryNote
    mu sync.RWMutex
}

func (m *mockMnemosyneServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    memory, exists := m.memories[req.Id]
    if !exists {
        return nil, status.Error(codes.NotFound, "memory not found")
    }

    return &pb.GetResponse{Memory: memory}, nil
}
```

---

## Work Stream 3: Critical TODO Resolution

### TODO #1: Delete Confirmation Dialog

**File**: `internal/memorydetail/model.go`
**Line**: ~TBD (search for "TODO: Show confirmation dialog")

**Problem**: No confirmation before deleting memories (accidental deletion risk)

**Solution**:
```go
type ConfirmationDialog struct {
    message string
    confirmed bool
    width int
    height int
}

func NewConfirmationDialog(message string) *ConfirmationDialog {
    return &ConfirmationDialog{
        message: message,
        width: 50,
        height: 10,
    }
}

func (d *ConfirmationDialog) View() string {
    title := titleStyle.Render("Confirm Deletion")
    msg := textStyle.Render(d.message)
    buttons := buttonStyle.Render("[Y]es") + "  " + buttonStyle.Render("[N]o")

    return lipgloss.JoinVertical(
        lipgloss.Center,
        title,
        "",
        msg,
        "",
        buttons,
    )
}
```

**Integration**:
- Show dialog when user presses 'D' (delete)
- Wait for 'Y' or 'N' keypress
- Only proceed with deletion if 'Y' confirmed

**Testing** (2 tests):
1. `TestConfirmationDialogYes` - Deletion proceeds on 'Y'
2. `TestConfirmationDialogNo` - Deletion cancelled on 'N'

### TODO #2: Relevance Scoring for Search

**File**: `internal/memorylist/model.go`
**Line**: ~TBD (search for "TODO: Implement relevance scoring")

**Problem**: Search results not sorted by relevance

**Solution**:
```go
type RelevanceScorer struct {
    query string
}

func (s *RelevanceScorer) Score(memory *pb.MemoryNote) float64 {
    score := 0.0

    // Exact match in content (high score)
    if strings.Contains(strings.ToLower(memory.Content), strings.ToLower(s.query)) {
        score += 10.0
    }

    // Match in tags (medium score)
    for _, tag := range memory.Tags {
        if strings.Contains(strings.ToLower(tag), strings.ToLower(s.query)) {
            score += 5.0
        }
    }

    // Importance boost
    score += float64(memory.Importance)

    // Recency boost (newer memories score higher)
    age := time.Since(memory.CreatedAt.AsTime())
    recencyScore := math.Max(0, 5.0 - (age.Hours() / 24 / 30)) // Decay over months
    score += recencyScore

    return score
}

func SortByRelevance(memories []*pb.MemoryNote, query string) []*pb.MemoryNote {
    scorer := &RelevanceScorer{query: query}

    sort.Slice(memories, func(i, j int) bool {
        return scorer.Score(memories[i]) > scorer.Score(memories[j])
    })

    return memories
}
```

**Testing** (3 tests):
1. `TestRelevanceScoringExactMatch` - Exact matches score highest
2. `TestRelevanceScoringImportance` - Importance affects score
3. `TestRelevanceScoringRecency` - Newer memories score higher

### TODO #3: Root Memory Selection

**File**: `internal/modes/explore.go`
**Line**: ~TBD (search for "TODO: Implement root memory selection")

**Problem**: Graph visualization doesn't intelligently select root

**Solution**:
```go
func selectRootMemory(memories []*pb.MemoryNote) *pb.MemoryNote {
    if len(memories) == 0 {
        return nil
    }

    // Strategy: Select memory with highest combined score
    type scoredMemory struct {
        memory *pb.MemoryNote
        score float64
    }

    scored := make([]scoredMemory, len(memories))
    for i, mem := range memories {
        score := 0.0

        // Importance (0-10)
        score += float64(mem.Importance) * 2.0

        // Link count (more links = more central)
        score += float64(len(mem.OutgoingLinks)) * 1.5
        score += float64(len(mem.IncomingLinks)) * 1.0

        // Recency (newer = more relevant)
        age := time.Since(mem.CreatedAt.AsTime())
        recencyScore := math.Max(0, 10.0 - (age.Hours() / 24 / 7)) // Decay over weeks
        score += recencyScore

        scored[i] = scoredMemory{memory: mem, score: score}
    }

    // Sort by score descending
    sort.Slice(scored, func(i, j int) bool {
        return scored[i].score > scored[j].score
    })

    return scored[0].memory
}
```

**Testing** (3 tests):
1. `TestRootSelectionByImportance` - High importance selected
2. `TestRootSelectionByLinks` - High link count selected
3. `TestRootSelectionByRecency` - Recent memories favored

---

## Work Stream 4: Security Hardening

### Rate Limiting

**Objective**: Prevent abuse of mnemosyne API

**Implementation**:
```go
package security

import (
    "sync"
    "time"
)

type RateLimiter struct {
    requests map[string]*requestCounter
    limit int
    window time.Duration
    mu sync.RWMutex
}

type requestCounter struct {
    count int
    reset time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string]*requestCounter),
        limit: limit,
        window: window,
    }
}

func (r *RateLimiter) Allow(clientID string) bool {
    r.mu.Lock()
    defer r.mu.Unlock()

    now := time.Now()
    counter, exists := r.requests[clientID]

    if !exists || now.After(counter.reset) {
        r.requests[clientID] = &requestCounter{
            count: 1,
            reset: now.Add(r.window),
        }
        return true
    }

    if counter.count >= r.limit {
        return false
    }

    counter.count++
    return true
}
```

**Configuration**:
- Default: 100 requests per minute
- Configurable via environment variable: `PEDANTIC_RAVEN_RATE_LIMIT`

### Input Sanitization

**Search Query Sanitization**:
```go
func SanitizeSearchQuery(query string) string {
    // Remove null bytes
    query = strings.ReplaceAll(query, "\x00", "")

    // Limit length
    maxLength := 500
    if len(query) > maxLength {
        query = query[:maxLength]
    }

    // Remove control characters
    query = strings.Map(func(r rune) rune {
        if r < 32 && r != '\n' && r != '\t' {
            return -1
        }
        return r
    }, query)

    return query
}
```

### Audit Logging

**Log Security Events**:
- Authentication attempts (success/failure)
- Token validation failures
- Rate limit violations
- Memory deletion operations

**Implementation**:
```go
type AuditLogger struct {
    log *slog.Logger
}

func (a *AuditLogger) LogAuthAttempt(success bool, clientID string) {
    a.log.Info("auth_attempt",
        "success", success,
        "client_id", clientID,
        "timestamp", time.Now().Unix(),
    )
}

func (a *AuditLogger) LogRateLimitViolation(clientID string) {
    a.log.Warn("rate_limit_violation",
        "client_id", clientID,
        "timestamp", time.Now().Unix(),
    )
}
```

---

## Parallelization Strategy

### Safe Parallel Execution

**Agent 1 (Haiku)**: Token Authentication (Stream 1)
- Simple implementation
- Well-defined requirements
- No dependencies
- **Estimated**: 4 hours

**Agent 2 (Sonnet)**: mnemosyne Test Coverage (Stream 2)
- Complex testing scenarios
- Mock server required
- 18 new tests
- **Estimated**: 8-10 hours

**Agent 3 (Haiku)**: Critical TODO Resolution (Stream 3)
- 3 independent TODOs
- UI components
- Business logic
- **Estimated**: 6-8 hours

**Sequential**: Security Hardening (Stream 4)
- Depends on authentication (Agent 1)
- Integrates with multiple components
- **Estimated**: 4 hours (after Agent 1)

### Execution Order

**Phase 1** (Parallel):
- Launch Agent 1 (Authentication)
- Launch Agent 2 (Test Coverage)
- Launch Agent 3 (TODO Resolution)

**Phase 2** (Sequential, after Agent 1):
- Security Hardening (rate limiting, audit logging)

**Phase 3** (Integration):
- Merge all changes
- Run full test suite
- Update documentation
- Commit and push

---

## Documentation Updates

### Files to Update

1. **docs/DEPLOYMENT.md** (add Authentication section)
   - Environment variable configuration
   - Token generation
   - Security considerations

2. **README.md** (add Authentication section)
   - Quick start with authentication
   - Security features

3. **docs/SECURITY.md** (create new file)
   - Authentication model
   - Rate limiting
   - Input sanitization
   - Audit logging
   - Best practices

---

## Success Metrics

### Week 2 Targets

- [ ] Authentication: 100% working (10 tests)
- [ ] mnemosyne coverage: 59.9% → 75%+ (18 tests)
- [ ] Critical TODOs: 3/3 resolved (8 tests)
- [ ] Security: Rate limiting + audit logging (5 tests)
- [ ] **Total new tests**: 41
- [ ] **Documentation**: 3 files updated, 1 created
- [ ] **Test pass rate**: 100% maintained
- [ ] **Production readiness**: 85% → 95% (+10%)

### Quality Gates

- [ ] All tests passing (441+ tests total)
- [ ] Zero compiler warnings
- [ ] Zero lint warnings
- [ ] Zero race conditions (tested with `-race`)
- [ ] Documentation complete and accurate

---

## Risk Assessment

### Low Risk
- Token authentication (simple, well-understood)
- TODO resolution (isolated changes)

### Medium Risk
- Test coverage improvement (time-consuming, but safe)
- Security hardening (requires careful integration)

### Mitigation
- Test each component thoroughly before integration
- Use race detector for concurrent code
- Security audit after implementation

---

**Last Updated**: 2025-11-12
**Phase**: Week 2 Specification
**Status**: Ready for Execution
