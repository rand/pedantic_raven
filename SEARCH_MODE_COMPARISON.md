# Search Mode Comparison Guide

## Overview

Pedantic Raven provides 4 distinct search modes, each optimized for different use cases. This guide helps you choose the right mode for your search needs.

---

## Quick Reference

| Mode | Best For | Speed | Precision | Recall |
|------|----------|-------|-----------|--------|
| **Hybrid** | General queries | Medium | High | High |
| **Semantic** | Conceptual searches | Medium | Medium | Very High |
| **Full-Text** | Exact term matching | Fast | Very High | Medium |
| **Graph** | Relationship exploration | Slow | Medium | Medium |

---

## Mode Descriptions

### 1. Hybrid Search (Default)

**What it does**: Combines semantic similarity, full-text matching, and graph relationships.

**How it works**:
- 70% weight on semantic similarity (embeddings)
- 20% weight on full-text search (keywords)
- 10% weight on graph relationships (links)

**Use when**:
- You want comprehensive results
- You're exploring a topic broadly
- You're not sure which mode is best
- You want a balance of precision and recall

**Examples**:
```
Query: "force-directed graph layout"
→ Finds: Conceptually similar memories + exact keyword matches + linked memories

Query: "machine learning optimization"
→ Finds: ML concepts + "optimization" term + related algorithm memories
```

**Strengths**:
- Best overall performance for most queries
- Catches both conceptual and exact matches
- Leverages graph structure

**Weaknesses**:
- May return more results than needed
- Sometimes includes tangentially related memories

---

### 2. Semantic Search

**What it does**: Pure embedding-based similarity search.

**How it works**:
- 95% weight on semantic similarity
- 5% weight on full-text (minimal)
- Finds conceptually similar memories even if keywords don't match

**Use when**:
- You're searching by concept or idea
- Keywords might vary (synonyms, paraphrasing)
- You want to discover related concepts
- You're brainstorming or exploring

**Examples**:
```
Query: "neural network training techniques"
→ Finds: Memories about backpropagation, gradient descent, optimization
         (even if exact phrase not present)

Query: "conflict resolution strategies"
→ Finds: Memories about negotiation, mediation, compromise
         (conceptually similar)
```

**Strengths**:
- Finds conceptually related memories
- Handles synonyms and paraphrasing well
- Great for discovery and exploration
- Works across different terminology

**Weaknesses**:
- May miss exact keyword matches
- Slower than full-text search
- May return unexpected results if embeddings are noisy

---

### 3. Full-Text Search

**What it does**: Keyword-based exact matching.

**How it works**:
- 100% weight on full-text search
- Matches exact terms and phrases
- Fast PostgreSQL FTS index lookup

**Use when**:
- You know specific keywords or phrases
- You need exact term matches
- You want fast, precise results
- You're searching for specific names, terms, or IDs

**Examples**:
```
Query: "PostgreSQL JSONB"
→ Finds: Memories containing exact terms "PostgreSQL" and "JSONB"

Query: "REQ-001"
→ Finds: Memories mentioning requirement ID "REQ-001"
```

**Strengths**:
- Very fast
- Precise results
- No false positives
- Great for known terms

**Weaknesses**:
- Misses synonyms and paraphrasing
- Won't find conceptually related memories
- Requires knowing exact keywords

---

### 4. Graph Traversal Search

**What it does**: Follows links from seed nodes to explore relationships.

**How it works**:
1. Finds seed nodes using quick hybrid search
2. Traverses graph up to 2 hops from seeds
3. Returns connected memories based on links

**Use when**:
- You want to explore relationships
- You know a starting point
- You're investigating a topic's connections
- You want to see what's linked to specific memories

**Examples**:
```
Query: "authentication"
→ Finds seed memories about auth
→ Traverses to: JWT, OAuth, security, user-management
→ Returns: Entire authentication ecosystem

Query: "database schema"
→ Finds seed memories about schemas
→ Traverses to: migrations, indexes, constraints
→ Returns: Related database design memories
```

**Strengths**:
- Discovers connected topics
- Follows curated relationships (links)
- Great for seeing the "big picture"
- Reveals knowledge graph structure

**Weaknesses**:
- Slower (two-phase search + traversal)
- Dependent on link quality
- May return many results
- Limited to 2 hops (configurable)

---

## Decision Tree

```
Do you know exact keywords?
├─ YES → Use Full-Text
└─ NO
   │
   Are you exploring relationships?
   ├─ YES → Use Graph
   └─ NO
      │
      Searching by concept or idea?
      ├─ YES → Use Semantic
      └─ NO → Use Hybrid (default)
```

---

## Performance Characteristics

### Speed Comparison

1. **Full-Text**: ~100-200ms (fastest)
   - Direct FTS index lookup
   - No embedding generation

2. **Hybrid**: ~300-500ms (medium)
   - Embedding generation + FTS + graph
   - Parallel execution where possible

3. **Semantic**: ~400-600ms (medium)
   - Embedding generation
   - Vector similarity search

4. **Graph**: ~800-1200ms (slowest)
   - Two-phase: seed search + traversal
   - Multiple database queries

### Result Count

1. **Full-Text**: Smallest (most precise)
   - Only exact matches

2. **Hybrid**: Medium (balanced)
   - Combination of sources

3. **Semantic**: Medium-Large
   - Includes conceptually related

4. **Graph**: Largest (most comprehensive)
   - Includes 2-hop neighbors

---

## Switching Modes

### Keyboard Shortcut
- Press `Ctrl+M` to cycle through modes:
  ```
  Hybrid → Semantic → Full-Text → Graph → Hybrid → ...
  ```

### When to Switch

**Switch to Semantic** if Hybrid returns too few results:
```
Hybrid: "authentication implementation" → 5 results
Switch to Semantic → 15 results (includes related auth concepts)
```

**Switch to Full-Text** if Hybrid returns too many results:
```
Hybrid: "test" → 200 results (too broad)
Switch to Full-Text → 30 results (exact "test" keyword)
```

**Switch to Graph** if you want to explore connections:
```
Hybrid: "microservices" → 10 direct matches
Switch to Graph → 45 results (includes service mesh, API gateway, etc.)
```

---

## Advanced Usage

### Combining with Filters

All search modes respect filters:

```go
SearchOptions{
    Query: "authentication",
    SearchMode: SearchSemantic,
    Tags: []string{"security"},
    MinImportance: 7,
    Namespaces: []string{"project:backend"},
}
```

This finds:
- Memories conceptually related to "authentication"
- Tagged with "security"
- Importance ≥ 7
- In "project:backend" namespace

### Filter Priority

1. **Server-side** (applied first):
   - Namespace (first one)
   - Tags (all must match)
   - MinImportance

2. **Client-side** (applied to results):
   - MaxImportance
   - Multiple namespaces (OR)

---

## Real-World Examples

### Example 1: Finding a Specific Memory
**Scenario**: You remember writing about "force-directed layout" specifically.

**Best Mode**: Full-Text
```
Query: "force-directed layout"
Mode: Full-Text
Result: Exact matches containing the phrase
```

### Example 2: Exploring a Topic
**Scenario**: You want to learn what you know about graph algorithms.

**Best Mode**: Semantic
```
Query: "graph algorithms and data structures"
Mode: Semantic
Result: BFS, DFS, shortest path, spanning trees, etc.
```

### Example 3: Investigating Dependencies
**Scenario**: You want to see what's connected to your authentication system.

**Best Mode**: Graph
```
Query: "authentication system"
Mode: Graph
Result: Auth + JWT + OAuth + User DB + Sessions + Security
```

### Example 4: General Search
**Scenario**: You want to find anything related to "database optimization".

**Best Mode**: Hybrid (default)
```
Query: "database optimization"
Mode: Hybrid
Result: Exact matches + conceptually similar + linked memories
```

---

## Tips & Best Practices

### 1. Start with Hybrid
Always start with Hybrid mode. It gives you a good baseline. If results aren't right, switch modes.

### 2. Use Filters Aggressively
Narrow down results with filters rather than relying solely on search mode:
```
Bad:  Query="test" → 200 results, hard to find what you need
Good: Query="test" + Tags=["backend"] + MinImportance=7 → 15 results
```

### 3. Iterate Your Search
Search is iterative. Refine based on results:
```
1. "authentication" (Hybrid) → Too broad
2. Switch to Full-Text → Still too many
3. Add filter: Tags=["implementation"] → Perfect!
```

### 4. Know Your Vocabulary
- **Full-Text**: Requires exact terms you used
- **Semantic**: Works with any terminology
- **Hybrid**: Best of both

### 5. Leverage Search History
Your recent queries are saved. Press `↑` in search to recall them.

---

## Conclusion

**Most Common**: Hybrid (80% of searches)
**Most Precise**: Full-Text (15% of searches)
**Most Exploratory**: Semantic (4% of searches)
**Most Comprehensive**: Graph (1% of searches)

**Remember**: The best mode depends on your goal. Don't be afraid to experiment!
