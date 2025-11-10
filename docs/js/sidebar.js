// Dynamic sidebar content based on scroll position
(function() {
    // Detect which page we're on
    function getSectionComments() {
        const path = window.location.pathname;
        const isWhitepaper = path.includes('whitepaper.html');

        if (isWhitepaper) {
            return {
                sectionComments: whitepaperComments,
                subsectionComments: whitepaperSubsections
            };
        } else {
            return {
                sectionComments: indexComments,
                subsectionComments: indexSubsections
            };
        }
    }

    // Index page comments (original)
    const indexComments = {
        'abstract': '// Markdown + Obsidian + AST = knowledge graph',
        'challenge': '// Context engineering: the art of information architecture',
        'architecture': '// DuckDB + GLiNER + RocksDB + SQLite = semantic layer',
        'features': '// 4 modes: Edit, Explore, Link, Analyze',
        'validation': '// 934 tests • 65% coverage • v0.5.0 tagged',
        'comparison': '// Obsidian for structure, Notion for UX, minus the cloud',
        'conclusion': '// Offline-first semantic knowledge graphs ✓'
    };

    const indexSubsections = {
        'Context Engineering for AI Systems': '// Information architecture for LLMs',
        'Limitations of Current Approaches': '// Markdown lacks semantics, Obsidian lacks inference',
        'Design Goals': '// Offline-first, semantic-rich, zero vendor lock-in',
        'Core Memory System': '// RocksDB for graphs, SQLite for facts',
        'Component Structure': '// Modular Rust: Store, SemanticLayer, CLI, REPL',
        'Semantic Pipeline': '// Parse → Extract → Infer → Validate → Store',
        'Offline-First Design': '// No cloud, no tracking, no subscription fees',
        'Event System': '// Pub/sub for plugin extensibility',
        'Technology Stack': '// Rust + DuckDB + GLiNER + tree-sitter',
        'Edit Mode (Phase 2)': '// Create and modify notes with semantic validation',
        'Explore Mode (Phase 3)': '// Traverse knowledge graph, find connections',
        'Link Mode (Phase 4)': '// Bidirectional links with inference',
        'Analyze Mode (Phase 6 - In Progress)': '// NER + entity extraction + graph queries',
        'Test Coverage': '// 934 passing tests across integration + unit suites',
        'Code Validation': '// Every claim backed by tagged source code',
        'Resources': '// Source, docs, API reference, usage guide'
    };

    // Whitepaper page comments
    const whitepaperComments = {
        'executive-summary': '// Real-time semantic analysis as you type',
        'introduction-the-context-problem': '// AI needs structure, engineers maintain chaos',
        'system-architecture': '// Go TUI + mnemosyne gRPC + GLiNER ML',
        'core-features': '// Edit (analyze) → Explore (search) → Analyze (patterns)',
        'technical-deep-dive': '// Force-directed graphs + exponential backoff + PubSub',
        'user-workflows': '// Creating context → Exploring memory → Analyzing patterns',
        'quality--testing': '// 934 tests • <16ms renders • 500ms analysis',
        'comparison': '// ICS → Obsidian → Roam → Notion: Fight!',
        'getting-started': '// make build && ./pedantic_raven',
        'roadmap': '// Phase 5 complete (55%) → v1.0 April 2026',
        'conclusion': '// Offline-first semantic memory for developers'
    };

    const whitepaperSubsections = {
        // Introduction
        'Context Engineering for AI Systems': '// Architecture, components, gaps, requirements, dependencies',
        'Limitations of Current Approaches': '// No auto-extraction, manual links, GUI-only, no semantic',
        'The Vision for Pedantic Raven': '// 9 features: real-time, entities, relationships, holes, deps, graph, mnemosyne, offline, terminal',

        // System Architecture
        'High-Level Overview': '// Go TUI ↔ mnemosyne Server ↔ GLiNER Service',
        'Technology Stack': '// Go 1.25 + Bubble Tea + gRPC + GLiNER (optional)',
        'Component Architecture': '// 15 internal packages, domain-focused organization',
        'Design Patterns': '// Elm Architecture + Event-Driven + Strategy + Graceful Degradation',
        'Data Flow: Edit → mnemosyne → Explore': '// 4 phases: Create, Store, Retrieve, Visualize',

        // Core Features
        'Edit Mode (Phase 2 - Complete)': '// Full editor + semantic analysis + terminal',
        'Entity Extraction: Hybrid Strategy': '// Pattern (60-70%, <1ms) + GLiNER (85-95%, 100-300ms)',
        'Explore Mode (Phase 4-5 - Complete)': '// 4 search modes, offline support, graph viz',
        'Analyze Mode (Phase 6 - In Progress, ~30%)': '// Triple graph + entity frequency + patterns + holes',
        'Mode Switching': '// 5 modes: Edit, Explore, Analyze, Orchestrate, Collaborate',

        // Technical Deep-Dive
        'Semantic Analyzer Implementation': '// Unicode tokens + pattern keywords + ML inference',
        'Graph Algorithms': '// Fruchterman-Reingold + spatial grid + <16ms convergence',
        'Connection Management': '// 30s health checks + exponential backoff + sync queue',
        'Event System': '// 40+ events, PubSub broker, thread-safe, reactive UI',
        'Error Handling': '// Network/Server/Validation categories, retry strategies, graceful degradation',

        // User Workflows
        'Workflow 1: Creating Context Documents': '// Type → Analyze → Mark gaps → Save → Store',
        'Workflow 2: Exploring Semantic Memory': '// Search → Navigate → Link → Back/Forward',
        'Workflow 3: Analyzing Semantic Patterns': '// Triples → Entities → Relationships → Holes → Export',
        'Workflow 4: Offline Work': '// Offline mode → Edit/Create → Auto-reconnect → Sync',

        // Quality & Testing
        'Test Coverage': '// 934 tests: events (18) + buffer (52) + semantic (63) + graph (134)',
        'Design Goals (Not Hype)': '// <16ms renders, <500ms analysis, ~10-20MB RAM',
        'Quality Gates': '// 934/934 passing, no races, no leaks, docs updated',

        // Comparison
        'vs ICS (Legacy Python Tool)': '// Rich TUI vs CLI, Real-time vs Basic, 934 tests vs Minimal',
        'vs Obsidian (Note-Taking)': '// Auto vs Manual, Terminal vs Electron, mnemosyne vs None',
        'vs Roam Research (Graph Notes)': '// Auto extraction vs Manual, Offline vs Cloud, Free vs $15/mo',
        'vs Notion (Workspace)': '// Real-time semantic vs None, Terminal vs Web/mobile, Graph+vector vs Tables',

        // Getting Started
        'Installation': '// Prerequisites: Go 1.25+, Docker (optional), mnemosyne (optional)',
        'Configuration': '// config.toml or env vars: GLINER_ENABLED, MNEMOSYNE_ENABLED',
        'Optional: GLiNER Service Setup': '// Docker: docker-compose up -d OR Manual: uvicorn',
        'Optional: mnemosyne Server Setup': '// cargo install mnemosyne-server',
        'Keyboard Shortcuts': '// Global: ?/Ctrl+Q/Ctrl+P, Edit: Ctrl+S/F/H, Explore: /j/k/e/n/l',

        // Roadmap
        'Current Status: Phase 5 Complete (55%)': '// Foundation → Semantic → Advanced → Client → Real Integration ✓',
        'Completed Phases': '// Phase 1-5: 934 tests, Edit + Explore + Offline + Graph complete',
        'Planned Phases': '// Phase 6: Analyze (30% done), Phase 7: Orchestrate, Phase 8: Collaborate, Phase 9: Polish',
        'Timeline': '// Phase 6: Late Nov 2025, Phase 7-8: Dec-Feb 2026, Phase 9: Mar 2026, v1.0: Apr 2026',
        'Community & Contributions': '// MIT license, part of mnemosyne ecosystem, replacing ICS',

        // Conclusion
        'Key Achievements': '// Real-time semantic + graceful degradation + offline-first + 934 tests + terminal-native',
        'Call to Action': '// Try: make build && ./pedantic_raven, Contribute: GitHub issues/PRs',
        'Acknowledgments': '// Built on Bubble Tea, mnemosyne, GLiNER, Go community'
    };

    function updateSidebarContent() {
        const sidebar = document.querySelector('.sidebar-tagline');
        if (!sidebar) return;

        // Get page-specific comments
        const { sectionComments, subsectionComments } = getSectionComments();

        // Get all sections and headings
        const sections = [...document.querySelectorAll('section[id]')];
        const headings = [...document.querySelectorAll('h2, h3')];

        // Account for navbar height
        const navbarHeight = 80;
        const scrollPosition = window.scrollY + navbarHeight + 50;

        // Find current section
        let currentSection = null;
        for (let i = sections.length - 1; i >= 0; i--) {
            if (scrollPosition >= sections[i].offsetTop) {
                currentSection = sections[i].id;
                break;
            }
        }

        // Find nearest h3 for more granular commentary
        let nearestH3 = null;
        let minDistance = Infinity;

        for (const heading of headings) {
            if (heading.tagName === 'H3') {
                const distance = Math.abs(scrollPosition - heading.offsetTop);
                if (distance < minDistance && scrollPosition >= heading.offsetTop - 100) {
                    minDistance = distance;
                    nearestH3 = heading.textContent.trim();
                }
            }
        }

        // Prioritize subsection commentary if we're close to an h3
        if (nearestH3 && subsectionComments[nearestH3] && minDistance < 300) {
            sidebar.textContent = subsectionComments[nearestH3];
        } else if (currentSection && sectionComments[currentSection]) {
            sidebar.textContent = sectionComments[currentSection];
        } else {
            sidebar.textContent = '// Context Engineering';
        }
    }

    // Initialize on page load
    function init() {
        updateSidebarContent();

        // Update on scroll with throttling
        let ticking = false;
        window.addEventListener('scroll', function() {
            if (!ticking) {
                window.requestAnimationFrame(function() {
                    updateSidebarContent();
                    ticking = false;
                });
                ticking = true;
            }
        });
    }

    // Run on DOMContentLoaded
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
