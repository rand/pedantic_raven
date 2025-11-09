# Contributing to Pedantic Raven

Welcome! We're excited you want to contribute to Pedantic Raven. This guide will help you get started and ensure your contributions align with the project's vision and standards.

**Table of Contents**

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Workflow](#development-workflow)
4. [Issue Guidelines](#issue-guidelines)
5. [Pull Request Process](#pull-request-process)
6. [Code Style Guidelines](#code-style-guidelines)
7. [Testing Requirements](#testing-requirements)
8. [Commit Guidelines](#commit-guidelines)
9. [Review Process](#review-process)
10. [Release Process](#release-process)

---

## Code of Conduct

### Our Commitment

We are committed to providing a welcoming, inclusive, and respectful environment for all contributors regardless of:
- Age, body size, disability, ethnicity, gender identity, experience level
- National origin, political beliefs, race, religion, sexual identity, sexual orientation

### Expected Behavior

- Use welcoming and inclusive language
- Be respectful of differing opinions and experiences
- Accept constructive criticism gracefully
- Focus on what's best for the community
- Show empathy towards other community members

### Unacceptable Behavior

- Harassment or discrimination in any form
- Offensive comments about identity or experience
- Public or private harassment
- Publishing private information without consent
- Other conduct that violates professional norms

### Reporting Issues

If you experience or witness unacceptable behavior, please contact maintainers at:
- Open a confidential issue
- Email: [maintainer contact]

All reports will be reviewed and kept confidential.

---

## Getting Started

### Prerequisites

- Go 1.25+
- git
- Basic familiarity with Go and terminal applications

### Setup Development Environment

```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/YOUR_USERNAME/pedantic-raven.git
cd pedantic-raven

# Add upstream remote
git remote add upstream https://github.com/rand/pedantic-raven.git

# Install dependencies
go mod download

# Verify setup
go build -o pedantic_raven .
./pedantic_raven --help
```

### Optional: mnemosyne Server

For testing mnemosyne integration:

```bash
# Clone mnemosyne
git clone https://github.com/rand/mnemosyne.git
cd mnemosyne

# Build and run
make build
./mnemosyne serve

# In another terminal
MNEMOSYNE_ADDR=localhost:50051 go test ./...
```

---

## Development Workflow

### 1. Create a Feature Branch

```bash
# Update main branch
git checkout main
git pull upstream main

# Create feature branch
# Use format: feature/description or fix/description
git checkout -b feature/my-feature-name
```

### 2. Make Changes

```bash
# Edit files in your favorite editor
# Follow code style guidelines (see below)
# Add tests for new functionality

# Test locally
go test ./...
go build -o pedantic_raven .
```

### 3. Commit Changes

```bash
# Stage changes
git add internal/mypackage/*.go

# Commit with descriptive message
git commit -m "Add support for new entity type

- Implement EntityExtractor for custom types
- Add tests for extraction logic
- Update documentation with examples"

# Push to your fork
git push origin feature/my-feature-name
```

### 4. Create Pull Request

See [Pull Request Process](#pull-request-process) below.

---

## Issue Guidelines

### Reporting Bugs

Before creating a bug report, check if it's already been reported. When creating one, include:

**Title**: Clear, descriptive summary
```
"Editor crashes when opening large files (>10MB)"
```

**Description**:
```markdown
## Description
[Describe what happened]

## Steps to Reproduce
1. Open file larger than 10MB
2. Switch to Edit mode
3. [What happens]

## Expected Behavior
[What should happen]

## Actual Behavior
[What actually happened]

## Environment
- OS: macOS 14.1
- Go version: 1.25
- mnemosyne: Connected

## Logs
[Any error messages or output]

## Additional Context
[Screenshots, config files, etc.]
```

### Requesting Features

Describe the use case and desired behavior:

```markdown
## Summary
[One sentence description]

## Motivation
[Why do we need this feature?]
[What problem does it solve?]

## Proposed Solution
[How should this work?]
[Any design considerations?]

## Alternative Approaches
[Other ways to solve this problem?]

## Additional Context
[Examples, mockups, etc.]
```

### Asking Questions

Use GitHub Discussions (if available) or GitHub Issues marked as questions:

```markdown
## Question
[Your question here]

## Context
[What are you trying to accomplish?]
[What have you already tried?]
```

---

## Pull Request Process

### Before Opening PR

1. **Ensure code works**:
   ```bash
   go test ./...
   go build -o pedantic_raven .
   go vet ./...
   gofmt -w ./internal
   ```

2. **Check for conflicts**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

3. **Update documentation**:
   - Update relevant docs in `docs/`
   - Update CHANGELOG.md with brief description
   - Add examples if applicable

4. **Verify tests**:
   ```bash
   # All tests must pass
   go test ./... -v

   # Check coverage
   go test ./... -cover
   ```

### Opening the PR

1. Push your branch to your fork
2. Go to GitHub and create Pull Request
3. Fill out the PR template completely

**PR Title Format**:
```
[Type] Description (affects: packages)

Examples:
- [Feature] Add graph export to PDF (affects: analyze, export)
- [Fix] Prevent race condition in event broker (affects: app/events)
- [Docs] Update architecture guide (affects: docs)
```

**PR Description Template**:
```markdown
## Summary
[1-2 sentences describing the changes]

## Type of Change
- [ ] Bug fix (non-breaking change fixing issue)
- [ ] New feature (non-breaking, adds functionality)
- [ ] Breaking change (existing functionality changes)
- [ ] Documentation (docs or example updates)
- [ ] Performance improvement
- [ ] Refactoring (no behavior change)

## Related Issues
Closes #123
Related to #456

## Changes
- [Brief description of each change]
- [How it works]
- [Why this approach]

## Testing
- [ ] Added tests for new functionality
- [ ] All tests pass: `go test ./...`
- [ ] Tested with mnemosyne server: [Yes/No]
- [ ] Manual testing: [Describe]

## Documentation
- [ ] Updated relevant docs in `docs/`
- [ ] Updated CHANGELOG.md
- [ ] Added code examples if needed

## Checklist
- [ ] Code follows style guidelines
- [ ] No new warnings generated
- [ ] Changes have adequate test coverage
- [ ] Documentation is updated
- [ ] No breaking changes to public APIs
- [ ] Commits follow guidelines

## Performance Impact
- [ ] No performance impact
- [ ] Performance improved: [measurements]
- [ ] Performance impact acceptable: [justification]

## Additional Context
[Any additional information reviewers should know]
```

---

## Code Style Guidelines

### Go Formatting

**All code must be formatted with gofmt**:

```bash
# Format all Go files
go fmt ./...

# Check formatting
gofmt -l ./...  # List unformatted files
```

**Follow Go conventions**:
- Package names: lowercase, single word
- Type names: Exported PascalCase, unexported camelCase
- Function names: Same as type names
- Constants: UPPER_CASE
- Variables: camelCase

### Naming Conventions

```go
// Package (lowercase, single word)
package editor

// Exported type (PascalCase)
type Model struct {
    field int
}

// Unexported type (camelCase)
type editorModel struct {
    field int
}

// Exported function (PascalCase)
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    return m, nil
}

// Unexported function (camelCase)
func (m Model) update() {
}

// Constants (UPPER_CASE)
const DefaultTimeout = 30 * time.Second
const maxRetries = 3

// Variables (camelCase)
var buffer bytes.Buffer
```

### Documentation

**Document all exported symbols**:

```go
// Package editor provides text editing functionality with semantic analysis.
//
// It supports syntax highlighting, search/replace, and real-time entity extraction.
package editor

// Model represents the state of the editor mode.
//
// It manages the text buffer, viewport, and semantic analysis results.
type Model struct {
    buffer *buffer.Buffer
    viewport *viewport.Model
}

// Update processes a Bubble Tea message and returns updated model and command.
//
// It handles keyboard input, buffer operations, and semantic analysis events.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    // ...
}
```

### Error Handling

**Always check and handle errors**:

```go
// Good - explicit error check
if err := operation(); err != nil {
    return fmt.Errorf("failed to perform operation: %w", err)
}

// Good - when intentionally ignoring, add comment
if _, err := operation(); err != nil {
    // Error logged elsewhere, safe to ignore
}

// Bad - ignoring errors
_ = operation()
```

### Comments

**Comment exported functions and types, not implementations**:

```go
// Good - explains what and why
// Model represents the editor state.
type Model struct {}

// Bad - explains implementation details
// Set the buffer field
type Model struct {}
```

### Imports

**Organize imports in groups**:

```go
import (
    // Standard library
    "fmt"
    "os"

    // External packages (alphabetical)
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"

    // Internal packages
    "github.com/rand/pedantic-raven/internal/buffer"
    "github.com/rand/pedantic-raven/internal/editor"
)
```

---

## Testing Requirements

### Coverage Requirements

- **Overall**: 70%+ code coverage
- **Critical Path**: 90%+ coverage
- **Business Logic**: 80%+ coverage
- **UI Components**: 60%+ coverage (rendering hard to test)

### Writing Tests

**Test naming**:

```go
// Test function: TestFunctionName
func TestBuffer_Insert(t *testing.T) {
    b := buffer.New()
    b.Insert(0, "hello")
    assert.Equal(t, "hello", b.String())
}

// Table-driven tests
func TestBuffer_Insert_Cases(t *testing.T) {
    cases := []struct {
        name     string
        initial  string
        pos      int
        text     string
        expected string
    }{
        {"empty buffer", "", 0, "hello", "hello"},
        {"middle insert", "hllo", 1, "e", "hello"},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            b := buffer.New()
            b.SetText(tc.initial)
            b.Insert(tc.pos, tc.text)
            assert.Equal(t, tc.expected, b.String())
        })
    }
}
```

### Running Tests

```bash
# All tests
go test ./...

# With verbose output
go test ./... -v

# With coverage
go test ./... -cover

# With race detection
go test ./... -race

# Specific test
go test ./internal/editor -run TestBuffer
```

### Benchmark Tests

```go
func BenchmarkBuffer_Insert(b *testing.B) {
    buf := buffer.New()
    for i := 0; i < b.N; i++ {
        buf.Insert(0, "x")
    }
}

// Run benchmarks
// go test -bench=. -benchmem ./internal/buffer
```

### Integration Tests

Place integration tests in `internal/integration/`:

```go
// internal/integration/workflow_test.go
func TestEditAnalyzeWorkflow(t *testing.T) {
    app := NewTestApp(t)
    defer app.Cleanup()

    // Test cross-mode workflow
    editor := app.EditMode()
    // ... test code
}
```

---

## Commit Guidelines

### Commit Message Format

```
[Type] Subject (component)

Body explaining what and why (not how).

Fixes #123
Related-To #456
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation
- **style**: Code style (formatting, missing semicolons)
- **refactor**: Code refactoring (no behavior change)
- **perf**: Performance improvement
- **test**: Adding or updating tests
- **chore**: Build, dependencies, tooling

### Examples

```bash
# Feature commit
git commit -m "feat: Add PDF export for analysis results (analyze/export)

- Implement PDFExporter interface
- Add table and chart rendering
- Support custom styling via config"

# Bug fix commit
git commit -m "fix: Prevent race condition in event broker (app/events)

The broker's subscriber map was being accessed concurrently without
proper locking. Added RWMutex to protect reads and writes."

# Documentation commit
git commit -m "docs: Update architecture guide with new section

- Add data flow diagrams
- Document design patterns
- Include code examples"

# Test commit
git commit -m "test: Add comprehensive buffer operation tests (editor/buffer)

- Table-driven tests for insert/delete/replace
- Edge case coverage (empty, boundary)
- Performance benchmarks"
```

### Commit Best Practices

- **Atomic commits**: One logical change per commit
- **Descriptive messages**: Explain why, not how
- **Reference issues**: Use "Fixes #123", "Relates-To #456"
- **Readable diffs**: Don't mix formatting and logic changes
- **Small batches**: 5-10 meaningful commits per PR

---

## Review Process

### Code Review Standards

All PRs must be reviewed and approved before merging.

**Reviewers will check**:

1. **Correctness**: Does the code do what it claims?
2. **Testing**: Are tests adequate? Do they pass?
3. **Style**: Does it follow guidelines?
4. **Performance**: Are there performance concerns?
5. **Documentation**: Is it documented?
6. **Security**: Are there security issues?

### Responding to Review Comments

- **Address all comments**: Each comment requires a response
- **Ask for clarification**: If you don't understand, ask
- **Acknowledge**: Appreciate feedback even if you disagree
- **Explain reasoning**: If you can't implement suggestion, explain why
- **Push updates**: Force-push updates to same branch
- **Request re-review**: Comment "@reviewer" when ready

### Merging PRs

- Minimum 1 approval required
- All tests must pass
- No merge conflicts
- CI checks must pass
- Squash and merge preferred for cleaner history

---

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

Format: `v1.2.3`

### Release Checklist

1. **Update Version**:
   ```bash
   # Update version in main.go or version.go
   const Version = "1.2.3"
   ```

2. **Update CHANGELOG**:
   ```markdown
   ## [1.2.3] - 2025-11-09

   ### Added
   - New feature X

   ### Fixed
   - Bug fix for issue #123

   ### Changed
   - Breaking change: Y

   ### Deprecated
   - Feature Z (will be removed in 2.0)
   ```

3. **Create Git Tag**:
   ```bash
   git tag -a v1.2.3 -m "Release version 1.2.3"
   git push origin v1.2.3
   ```

4. **Create GitHub Release**:
   - Go to Releases
   - Create release from tag
   - Copy CHANGELOG section to release notes
   - Attach binaries if applicable

5. **Announce**:
   - Post to discussions
   - Update documentation
   - Email subscribers (if applicable)

---

## Getting Help

### Documentation

- [Architecture Guide](architecture.md) - System design and patterns
- [Developer Guide](DEVELOPMENT.md) - Setup and workflow
- [Edit Mode Guide](edit-mode-guide.md) - Edit mode documentation
- [Analyze Mode Guide](analyze-mode-guide.md) - Analyze mode documentation
- [Orchestrate Mode Guide](orchestrate-mode-guide.md) - Orchestrate mode documentation

### Community

- **GitHub Issues**: Ask questions, report bugs, request features
- **GitHub Discussions**: General questions and ideas
- **Pull Requests**: Collaborate on specific changes

### Maintainers

- **Code Review**: Expect feedback within 48 hours
- **Questions**: Comment on issues/PRs
- **Urgent**: Contact maintainers directly

---

## Summary

**Contributing to Pedantic Raven**:

1. **Setup**: Fork, clone, install dependencies
2. **Develop**: Create feature branch, make changes, add tests
3. **Test**: Run `go test ./...` and verify coverage
4. **Commit**: Write clear, atomic commits
5. **PR**: Open PR with complete description and checklist
6. **Review**: Respond to feedback, update code
7. **Merge**: Approved PRs are merged by maintainers

**Key Principles**:
- ✓ Write tests for all new code
- ✓ Follow Go conventions and style guidelines
- ✓ Document public APIs and functions
- ✓ Keep commits atomic and descriptive
- ✓ Respond to review feedback promptly
- ✓ Be respectful and inclusive

Thank you for contributing to Pedantic Raven!

---

*Last updated: November 9, 2025*
