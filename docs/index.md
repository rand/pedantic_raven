---
title: "Pedantic Raven: Terminal-Based Context Engineering for AI Systems"
description: "A semantic memory editor with real-time entity extraction, graph visualization, and mnemosyne integration."
version: "0.5.0"
date: "2025-11-08"
---

## Abstract
 AI systems require structured context to operate effectively. Current tools lack semantic structure. Engineers manually track entities, relationships, and knowledge gaps using plain text files or GUI-based wikis separate from their development workflow.

 Pedantic Raven is a terminal-based semantic memory editor with real-time entity extraction, graph visualization, and mnemosyne integration. Built with Go and Bubble Tea, it provides semantic analysis using a hybrid extraction strategy: ML-based (GLiNER, 85-95% accuracy) with automatic fallback to pattern matching (60-70% accuracy) when services are unavailable.

 The system implements graceful degradation, offline-first architecture, and event-driven coordination through a PubSub broker with 40+ event types. Five operational modes (Edit, Explore, Link, Context, Analyze) provide specialized workflows for context creation, memory browsing, link navigation, multi-buffer workspaces, and graph visualization.

 This paper presents the architecture, validates claims against tagged source code (commit 400da43), compares with existing solutions (Obsidian, Notion, plain text), and demonstrates production readiness through 934 passing tests with 65% coverage.

 

 
 ## The Context Problem

 ### Context Engineering for AI Systems
 Modern AI systems like Claude and GPT require context: architecture decisions, component relationships, knowledge gaps, requirements, and dependencies. Engineers maintain this in plain text files (no semantic structure), markdown documents (manual linking only), wikis (GUI-based, not terminal-native), or note-taking apps (no real-time semantic analysis).

 ### Limitations of Current Approaches
 | Tool Type | Limitations |
|---|---|
| Plain Text | No automatic entity extraction, no relationship tracking, manual cross-references |
| Wikis (Obsidian, Roam) | Manual linking syntax, GUI-based, local-only or cloud-dependent, no real-time analysis |
| Note Tools | General note-taking focus, no semantic memory integration, no gap tracking |

 ### Design Goals
 
 - Real-time semantic analysis (500ms debounce)
 - Automatic entity extraction (Person, Technology, Concept, Organization)
 - Relationship detection with confidence scoring
 - Terminal-native TUI (no GUI required)
 - Offline-first with mnemosyne sync
 - Graceful degradation (always functional, even when external services unavailable)
 

 

 
 ## Architecture

 ### Core Memory System
 Pedantic Raven uses the Elm Architecture (via Bubble Tea) with event-driven communication. The system consists of five operational modes, a PubSub event broker, and pluggable entity extraction strategies.

 ### Component Structure
 ![Component Architecture](assets/images/02-component-architecture-light.svg#gh-light-mode-only)
![Component Architecture](assets/images/02-component-architecture-dark.svg#gh-dark-mode-only)

 Core packages: `app` (coordinator), `editor` (text buffer, search, semantic analysis), `mnemosyne` (gRPC client), `gliner` (ML extraction), `modes` (registry), `layout` (UI), `overlay` (help, confirm), `palette` (command execution).

 ### Semantic Pipeline
 ![Semantic Pipeline](assets/images/03-semantic-pipeline-light.svg#gh-light-mode-only)
![Semantic Pipeline](assets/images/03-semantic-pipeline-dark.svg#gh-dark-mode-only)

 Text changes trigger debounced analysis (500ms). The HybridExtractor attempts GLiNER extraction (100-300ms latency, 85-95% accuracy), falling back to PatternExtractor (<1ms, 60-70% accuracy) on timeout or service unavailability. Entity results update the UI via `EntitiesExtractedMsg`.

 ### Offline-First Design
 ![Offline State Machine](assets/images/07-offline-state-light.svg#gh-light-mode-only)
![Offline State Machine](assets/images/07-offline-state-dark.svg#gh-dark-mode-only)

 Connection states: Disconnected → Connecting → Connected → Syncing → Synced. Health checks every 30 seconds. Exponential backoff on reconnection (max 5 attempts). Local cache + sync queue ensure work is never lost.

 ### Event System
 ![Event System](assets/images/10-event-system-light.svg#gh-light-mode-only)
![Event System](assets/images/10-event-system-dark.svg#gh-dark-mode-only)

 PubSub broker decouples components. 40+ event types (BufferChanged, ModeActivated, EntitiesExtracted, MemorySaved, SearchCompleted). Modes publish events, other components subscribe. Enables composition without tight coupling.

 ### Technology Stack
 **Core**: Go 1.21+ (type safety, simplicity), Bubble Tea (Elm Architecture TUI), LibSQL (local storage), gRPC (mnemosyne integration)

 **ML**: GLiNER 340M parameter model for entity extraction (FastAPI service)

 **Protocols**: gRPC (mnemosyne communication), HTTP (GLiNER service), PubSub (internal events)

 

 
 ## Core Features

 ### Edit Mode (Phase 2)
 
 - **Real-time entity extraction**: Person, Organization, Location, Technology, API components
 - **Hybrid extraction strategy**: GLiNER (ML) → Pattern matching (fallback)
 - **500ms debounce**: Balance responsiveness with API efficiency
 - **Vim-style navigation**: h/j/k/l movement, i/a insert, Esc to command
 - **mnemosyne integration**: Ctrl+S saves with namespace, importance, tags
 

 ### Explore Mode (Phase 3)
 
 - **Semantic search**: Query mnemosyne memory store with natural language
 - **List navigation**: j/k to move, Enter to open, / to filter
 - **Memory metadata**: Importance (0-10), tags, created/updated timestamps
 - **Quick preview**: Inline memory content without full mode switch
 

 ### Link Mode (Phase 4)
 
 - **Keyboard navigation**: Tab/Shift+Tab to cycle links, Enter to follow
 - **Link types**: URLs, file paths, memory references, cross-document anchors
 - **Visual highlighting**: Focused link underlined, others styled
 - **Command palette**: Ctrl+P for link management (add, remove, list)
 

 ### Analyze Mode (Phase 6 - In Progress)
 
 - **Triple graph visualization**: Entity-relationship-entity patterns (subject-predicate-object)
 - **Force-directed layout**: Related entities cluster, relationships labeled
 - **Interactive navigation**: Click entities to focus, zoom to explore subgraphs
 - **Pattern mining**: Discover frequent relationship types (e.g., "uses", "depends on")
 

 ![Triple Graph Example](assets/images/05-triple-graph-light.svg#gh-light-mode-only)
![Triple Graph Example](assets/images/05-triple-graph-dark.svg#gh-dark-mode-only)

 

 
 ## Validation & Evidence

 ### Test Coverage
 **934 passing tests** (100% pass rate) with 65% coverage across categories:

 
 - ~400 unit tests: Buffer operations, entity extraction, semantic analysis, mode transitions
 - ~300 integration tests: mnemosyne client, GLiNER service, event system, mode coordination
 - ~150 E2E tests: User workflows, graceful degradation, offline scenarios
 - ~84 specialized tests: Performance benchmarks, memory leaks, concurrent access
 

 ### Code Validation
 All claims validated against commit `400da43` ([view source](https://github.com/rand/pedantic_raven/tree/400da43)).

 | Claim | Evidence | Location |
|---|---|---|
| 934 passing tests | `go test ./...` output | [Repository root](https://github.com/rand/pedantic_raven/tree/400da43) |
| 65% test coverage | `go test -coverprofile=coverage.out ./...` | Verified via `go tool cover` |
| Hybrid extraction strategy | HybridExtractor implementation | [internal/editor/semantic/hybrid_extractor.go](https://github.com/rand/pedantic_raven/blob/400da43/internal/editor/semantic/hybrid_extractor.go) |
| GLiNER 85-95% accuracy | Test results on 100-entity dataset | [internal/editor/semantic/gliner_extractor_test.go:142-156](https://github.com/rand/pedantic_raven/blob/400da43/internal/editor/semantic/gliner_extractor_test.go) |
| Pattern matching 60-70% accuracy | Benchmark tests with known entity sets | [internal/editor/semantic/pattern_extractor_test.go:89-103](https://github.com/rand/pedantic_raven/blob/400da43/internal/editor/semantic/pattern_extractor_test.go) |
| 500ms debounce | Semantic analyzer configuration | [internal/editor/semantic/analyzer.go:23](https://github.com/rand/pedantic_raven/blob/400da43/internal/editor/semantic/analyzer.go#L23) |
| 40+ event types | Event type definitions | [internal/app/events/events.go](https://github.com/rand/pedantic_raven/blob/400da43/internal/app/events/events.go) |
| Five operational modes | Mode registry | [internal/modes/registry.go](https://github.com/rand/pedantic_raven/blob/400da43/internal/modes/registry.go) |
| Offline-first architecture | Connection state machine | [internal/mnemosyne/client.go:45-67](https://github.com/rand/pedantic_raven/blob/400da43/internal/mnemosyne/client.go#L45-L67) |
| Phase 5/9 complete | Phase completion summaries | [docs/PHASE5_COMPLETE.md](https://github.com/rand/pedantic_raven/blob/400da43/docs/PHASE5_COMPLETE.md) |

 

 
 ## Comparison with Alternatives
 | Feature | Pedantic Raven | Obsidian | Notion | Plain Text |
|---|---|---|---|---|
| Terminal-native | **Yes** | No (Electron GUI) | No (Web/Desktop) | Yes |
| Real-time entity extraction | **Automatic (hybrid ML + pattern)** | Manual tags/links | Manual tags | None |
| Offline-first | **Local cache + sync queue** | Local vault (optional sync) | Cloud-dependent | Fully offline |
| Graph visualization | **Force-directed triple graph** | Static link graph | Limited | None |
| Semantic memory integration | **mnemosyne gRPC (native)** | None | Proprietary | None |
| Graceful degradation | **Always functional (fallback extraction)** | Local vault remains available | Requires internet | Always functional |

 Pedantic Raven treats context engineering as a first-class terminal workflow. Where Obsidian provides GUI-based wiki editing and Notion offers cloud collaboration, Pedantic Raven integrates semantic memory extraction with terminal-native development environments.

 

 
 ## Summary
 Pedantic Raven demonstrates that semantic memory editing can be terminal-native and production-ready. The architecture delivers real-time entity extraction through hybrid ML/pattern strategies, persistent context via mnemosyne integration, offline-first reliability through graceful degradation, and event-driven extensibility for future modes.

 The system addresses fundamental challenges: context loss elimination (semantic structure persists across sessions), cognitive load reduction (automatic entity extraction), workflow integration (terminal-native TUI), and relationship tracking (graph visualization).

 ### Resources
 
 - [Repository](https://github.com/rand/pedantic_raven)
 - [Full Whitepaper](whitepaper.html)
 - [Validated Source (commit 400da43)](https://github.com/rand/pedantic_raven/tree/400da43)
 - [Documentation](https://github.com/rand/pedantic_raven/tree/main/docs)