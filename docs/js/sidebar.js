// Dynamic sidebar content based on scroll position
(function() {
    // Section-specific comments for pedantic_raven (mix of technical + dry wit)
    const sectionComments = {
        'abstract': '// Markdown + Obsidian + AST = knowledge graph',
        'challenge': '// Context engineering: the art of information architecture',
        'architecture': '// DuckDB + GLiNER + RocksDB + SQLite = semantic layer',
        'features': '// 4 modes: Edit, Explore, Link, Analyze',
        'validation': '// 934 tests • 65% coverage • v0.5.0 tagged',
        'comparison': '// Obsidian for structure, Notion for UX, minus the cloud',
        'conclusion': '// Offline-first semantic knowledge graphs ✓'
    };

    // Subsection commentary (detected via nearest h3)
    const subsectionComments = {
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

    function updateSidebarContent() {
        const sidebar = document.querySelector('.sidebar-tagline');
        if (!sidebar) return;

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
            sidebar.textContent = '// Semantic knowledge graph for Markdown';
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
