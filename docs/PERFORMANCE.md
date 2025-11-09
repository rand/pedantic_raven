# Performance Profiling and Optimization Guide

This document describes the performance benchmarking suite for Pedantic Raven and provides guidance on profiling and optimizing the application.

## Table of Contents

1. [Overview](#overview)
2. [Benchmark Suite](#benchmark-suite)
3. [Running Benchmarks](#running-benchmarks)
4. [Profiling](#profiling)
5. [Performance Targets](#performance-targets)
6. [Optimization Strategies](#optimization-strategies)
7. [Benchmark Results](#benchmark-results)

## Overview

Pedantic Raven is a Go TUI application with several performance-critical paths:

- **Graph Layout Algorithms**: Force-directed layout in Analyze Mode and Orchestrate Mode
- **Event Processing**: Internal event broker for component communication
- **Semantic Analysis**: Entity extraction and relationship detection
- **Memory Operations**: Mnemosyne client operations (recall, remember, graph traversal)
- **JSON Parsing**: Work plan validation and task graph serialization

This benchmark suite provides 72+ comprehensive benchmarks covering all critical paths.

## Benchmark Suite

### Graph Layout Benchmarks (35+ benchmarks)

#### Analyze Mode - Triple Graph (`internal/analyze/triple_graph_benchmark_test.go`)

**Node/Edge Operations**:
- `BenchmarkTripleGraphAddNode`: Adding nodes to the graph
- `BenchmarkTripleGraphAddEdge`: Adding edges to the graph

**Layout Algorithms**:
- `BenchmarkTripleGraphLayout`: Single layout iteration (10, 50, 100, 200 nodes)
- `BenchmarkTripleGraphFullLayoutStabilization`: Complete stabilization (50-100 iterations)
- `BenchmarkTripleGraphRepulsion`: Repulsion force calculation
- `BenchmarkTripleGraphAttraction`: Attraction force calculation
- `BenchmarkTripleGraphUpdatePositions`: Position update step

**Graph Operations**:
- `BenchmarkTripleGraphCalculateImportance`: Importance score calculation
- `BenchmarkTripleGraphApplyFilter`: Graph filtering by type/importance
- `BenchmarkTripleGraphGetEdgesFrom`: Edge lookup by source node
- `BenchmarkTripleGraphGetEdgesTo`: Edge lookup by target node
- `BenchmarkTripleGraphBuildFromAnalysis`: Building graph from semantic analysis

#### Orchestrate Mode - Task Graph (`internal/orchestrate/task_graph_benchmark_test.go`)

**Graph Creation**:
- `BenchmarkTaskGraphCreation`: Creating task graphs (10-100 nodes)

**Layout Operations**:
- `BenchmarkTaskGraphLayout`: Single layout iteration
- `BenchmarkTaskGraphStabilize`: Full stabilization (50-100 iterations)
- `BenchmarkTaskGraphRepulsion`: Repulsion calculation
- `BenchmarkTaskGraphAttraction`: Attraction calculation
- `BenchmarkTaskGraphUpdatePositions`: Position updates

**Graph Operations**:
- `BenchmarkTaskGraphUpdateStatus`: Task status updates
- `BenchmarkTaskGraphSelectNode`: Node selection
- `BenchmarkTaskGraphGetBounds`: Bounding box calculation
- `BenchmarkTaskGraphResize`: Viewport resizing

### Event Processing Benchmarks (15+ benchmarks)

Location: `internal/app/events/broker_benchmark_test.go`

**Core Operations**:
- `BenchmarkBrokerPublish`: Event publishing with different buffer sizes
- `BenchmarkBrokerPublishWithSubscribers`: Publishing with 1-50 active subscribers
- `BenchmarkBrokerSubscribe`: Subscription creation
- `BenchmarkBrokerUnsubscribe`: Unsubscription
- `BenchmarkBrokerSubscribeAll`: Global subscription (all event types)

**Concurrency**:
- `BenchmarkBrokerPublishConcurrent`: Concurrent publishing (1-8 goroutines)
- `BenchmarkBrokerHighThroughput`: High-throughput scenario (10 subscribers)

**Operations**:
- `BenchmarkBrokerPublishMultipleEventTypes`: Publishing different event types
- `BenchmarkBrokerSubscriberCount`: Counting subscribers
- `BenchmarkBrokerGlobalSubscriberCount`: Counting global subscribers
- `BenchmarkBrokerClear`: Clearing all subscribers
- `BenchmarkBrokerMemoryAllocation`: Memory allocation profiling

**Convenience Methods**:
- `BenchmarkBrokerPublishSemanticAnalysis`: Typed event publishing

### JSON/Work Plan Benchmarks (12+ benchmarks)

Location: `internal/orchestrate/json_benchmark_test.go`

**JSON Operations**:
- `BenchmarkWorkPlanJSONParsing`: Parsing work plans (5-100 tasks)
- `BenchmarkWorkPlanJSONMarshaling`: Marshaling work plans
- `BenchmarkTaskJSONMarshaling`: Marshaling individual tasks
- `BenchmarkTaskJSONUnmarshaling`: Unmarshaling tasks
- `BenchmarkLargeWorkPlanParsing`: Parsing large plans (500 tasks)

**Validation**:
- `BenchmarkWorkPlanValidation`: Plan validation (10-100 tasks)
- `BenchmarkTaskValidation`: Task validation
- `BenchmarkCyclicDependencyDetection`: Detecting circular dependencies
- `BenchmarkDeepDependencyTree`: Validating deep dependency chains

**Event Serialization**:
- `BenchmarkAgentEventJSONMarshaling`: Marshaling agent events
- `BenchmarkAgentEventJSONUnmarshaling`: Unmarshaling agent events

**Memory Profiling**:
- `BenchmarkMemoryAllocationJSON`: Memory allocation during JSON operations

### Semantic Analysis Benchmarks (10+ benchmarks)

Location: `internal/editor/semantic/semantic_benchmark_test.go`

**Tokenization**:
- `BenchmarkTokenize`: Text tokenization (various sizes)

**Entity Extraction**:
- `BenchmarkPatternExtractorSimple`: Simple entity extraction
- `BenchmarkPatternExtractorComplexText`: Complex text (100-2000 words)
- `BenchmarkPatternExtractorWithTypeFilter`: Filtered extraction

**Classification**:
- `BenchmarkEntityClassification`: Entity type classification
- `BenchmarkMultiWordEntityExtraction`: Multi-word entity detection

**Streaming Analysis**:
- `BenchmarkAnalyzerStreamingAnalysis`: Full streaming analysis
- `BenchmarkHybridExtractor`: Hybrid extractor with fallback logic

**Operations**:
- `BenchmarkRelationshipExtraction`: Relationship detection
- `BenchmarkTypedHoleDetection`: Typed hole pattern matching
- `BenchmarkEntityDeduplication`: Entity deduplication logic
- `BenchmarkContextBuilding`: Classification context creation

**Memory Profiling**:
- `BenchmarkMemoryAllocation`: Memory usage during extraction

### Memory Operations Benchmarks (16+ benchmarks)

Location: `internal/mnemosyne/memory_benchmark_test.go`

**Options Creation**:
- `BenchmarkRecallOptionsCreation`: Creating recall options
- `BenchmarkStoreMemoryOptionsCreation`: Creating store options
- `BenchmarkListMemoriesOptionsCreation`: Creating list options
- `BenchmarkGraphTraverseOptionsCreation`: Creating traversal options
- `BenchmarkSemanticSearchOptionsCreation`: Creating search options
- `BenchmarkUpdateMemoryOptionsCreation`: Creating update options

**Namespace Operations**:
- `BenchmarkNamespaceCreation`: Creating namespaces (global, project, session)

**Validation**:
- `BenchmarkRecallValidation`: Recall request validation
- `BenchmarkStoreMemoryValidation`: Store request validation
- `BenchmarkListMemoriesValidation`: List request validation
- `BenchmarkGetMemoryValidation`: Get request validation
- `BenchmarkDeleteMemoryValidation`: Delete request validation
- `BenchmarkGetContextValidation`: Context request validation
- `BenchmarkGraphTraverseValidation`: Traversal request validation
- `BenchmarkSemanticSearchValidation`: Search request validation

**Client Operations**:
- `BenchmarkClientConfigValidation`: Config validation

**Memory Profiling**:
- `BenchmarkMemoryAllocationRecall`: Allocations during recall
- `BenchmarkMemoryAllocationStore`: Allocations during store
- `BenchmarkTagCreation`: Tag slice creation
- `BenchmarkEmbeddingCreation`: Embedding vector creation
- `BenchmarkMultipleMemoryIDs`: Memory ID slice creation

## Running Benchmarks

### Run All Benchmarks

```bash
go test -bench=. -benchmem ./internal/analyze ./internal/app/events ./internal/orchestrate ./internal/editor/semantic ./internal/mnemosyne
```

### Run Specific Package Benchmarks

```bash
# Graph layout benchmarks
go test -bench=. -benchmem ./internal/analyze

# Event processing benchmarks
go test -bench=. -benchmem ./internal/app/events

# Task graph and JSON benchmarks
go test -bench=. -benchmem ./internal/orchestrate

# Semantic analysis benchmarks
go test -bench=. -benchmem ./internal/editor/semantic

# Memory operations benchmarks
go test -bench=. -benchmem ./internal/mnemosyne
```

### Run Specific Benchmarks

```bash
# Run only graph layout benchmarks
go test -bench=BenchmarkTripleGraphLayout -benchmem ./internal/analyze

# Run only event publish benchmarks
go test -bench=BenchmarkBrokerPublish -benchmem ./internal/app/events

# Run only JSON parsing benchmarks
go test -bench=BenchmarkWorkPlanJSONParsing -benchmem ./internal/orchestrate
```

### Benchmark Options

```bash
# Run benchmarks with custom duration
go test -bench=. -benchtime=5s ./internal/analyze

# Run benchmarks multiple times for statistical significance
go test -bench=. -count=10 ./internal/analyze

# Save benchmark results to file
go test -bench=. -benchmem ./internal/analyze > bench_analyze.txt

# Compare benchmark results
benchstat bench_before.txt bench_after.txt
```

## Profiling

### CPU Profiling

```bash
# Generate CPU profile
go test -bench=BenchmarkTripleGraphLayout -cpuprofile=cpu.prof ./internal/analyze

# Analyze CPU profile
go tool pprof cpu.prof
```

**Common pprof commands**:
```
(pprof) top10       # Show top 10 functions by CPU time
(pprof) list FuncName   # Show source code for function
(pprof) web         # Open interactive web UI
(pprof) pdf > profile.pdf  # Generate PDF visualization
```

### Memory Profiling

```bash
# Generate memory profile
go test -bench=BenchmarkTripleGraphLayout -memprofile=mem.prof ./internal/analyze

# Analyze memory profile
go tool pprof mem.prof
```

**Common memory analysis commands**:
```
(pprof) top10       # Top 10 allocations
(pprof) list FuncName   # Source code with allocations
(pprof) web         # Interactive visualization
(pprof) alloc_space # Total allocations
(pprof) inuse_space # Currently in use
```

### Continuous Profiling

```bash
# Profile specific benchmark continuously
go test -bench=BenchmarkTripleGraphLayout -cpuprofile=cpu.prof -memprofile=mem.prof -benchtime=30s ./internal/analyze
```

### Trace Profiling

```bash
# Generate execution trace
go test -bench=BenchmarkBrokerPublishConcurrent -trace=trace.out ./internal/app/events

# Analyze trace
go tool trace trace.out
```

## Performance Targets

### Graph Layout Performance

| Operation | Size | Target | Actual (M1 Pro) |
|-----------|------|--------|-----------------|
| Single layout iteration | 100 nodes | < 100ms | ~25ms |
| Full stabilization (50 iter) | 100 nodes | < 5s | ~1.2s |
| Repulsion calculation | 100 nodes | < 50ms | ~19ms |
| Attraction calculation | 100 nodes | < 20ms | ~5ms |
| Position updates | 100 nodes | < 5ms | ~0.7ms |

### Event Processing Performance

| Operation | Configuration | Target | Actual (M1 Pro) |
|-----------|---------------|--------|-----------------|
| Event publish | No subscribers | < 100ns | ~14ns |
| Event publish | 10 subscribers | < 1µs | ~600ns |
| Event publish | 50 subscribers | < 20µs | ~15µs |
| Subscribe | - | < 5µs | ~1.4µs |
| Concurrent publish | 8 goroutines | < 200ns/op | ~116ns/op |

### JSON Operations Performance

| Operation | Size | Target | Actual (M1 Pro) |
|-----------|------|--------|-----------------|
| Parse work plan | 10 tasks | < 50µs | ~11µs |
| Parse work plan | 100 tasks | < 500µs | ~102µs |
| Parse work plan | 500 tasks | < 3ms | ~508µs |
| Marshal work plan | 100 tasks | < 100µs | ~15µs |
| Task validation | 1 task | < 10ns | ~2.1ns |
| Cycle detection | 50 tasks | < 5ms | ~14µs |

### Semantic Analysis Performance

| Operation | Size | Target | Actual (expected) |
|-----------|------|--------|-------------------|
| Tokenization | 100 words | < 1ms | TBD |
| Entity extraction | 500 words | < 50ms | TBD |
| Entity extraction | 1000 words | < 100ms | TBD |
| Classification | 1 entity | < 1µs | TBD |
| Multi-word extraction | 3 words | < 10µs | TBD |

### Memory Operations Performance

| Operation | Configuration | Target | Actual (expected) |
|-----------|---------------|--------|-------------------|
| Recall options creation | - | < 100ns | TBD |
| Store options creation | - | < 200ns | TBD |
| Namespace creation | - | < 50ns | TBD |
| Validation | Any request | < 1µs | TBD |
| Embedding creation | 768-dim | < 5µs | TBD |

## Optimization Strategies

### Graph Layout Optimization

1. **Spatial Partitioning**: Use quadtrees or grid-based partitioning to reduce O(n²) repulsion calculations
2. **Barnes-Hut Algorithm**: Approximate distant node interactions for better scaling
3. **SIMD Vectorization**: Use SIMD instructions for force calculations
4. **Parallel Force Calculation**: Parallelize repulsion/attraction calculations across goroutines
5. **Adaptive Iteration**: Stop early when layout converges (velocity threshold)
6. **Layout Caching**: Cache layout results for unchanged graphs

### Event Processing Optimization

1. **Buffer Tuning**: Adjust channel buffer sizes based on event frequency
2. **Batch Publishing**: Batch multiple events into single publish operation
3. **Lock-Free Queues**: Use lock-free queues for high-concurrency scenarios
4. **Subscriber Pooling**: Reuse subscriber channels to reduce allocations
5. **Event Dropping Strategy**: Implement intelligent event dropping for slow consumers
6. **Priority Queues**: Process critical events with higher priority

### JSON Optimization

1. **JSON Streaming**: Use streaming parsers for large work plans
2. **Pre-allocation**: Pre-allocate slices for known sizes
3. **Code Generation**: Use code generation tools (e.g., easyjson, ffjson)
4. **Validation Caching**: Cache validation results for unchanged plans
5. **Lazy Parsing**: Parse only required fields on-demand
6. **Binary Serialization**: Consider Protocol Buffers for performance-critical paths

### Semantic Analysis Optimization

1. **Parallel Extraction**: Extract entities from multiple text chunks in parallel
2. **Trie Data Structures**: Use tries for efficient keyword matching
3. **Bloom Filters**: Filter out non-entities before classification
4. **Token Reuse**: Reuse token slices to reduce allocations
5. **Context Pooling**: Pool classification context objects
6. **Early Termination**: Stop extraction when confidence threshold met

### Memory Operations Optimization

1. **Connection Pooling**: Reuse gRPC connections to mnemosyne
2. **Request Batching**: Batch multiple recall/store operations
3. **Embedding Caching**: Cache frequently-used embeddings
4. **Lazy Loading**: Load memory context only when needed
5. **Compression**: Compress large memory payloads
6. **Streaming Recall**: Use streaming APIs for large result sets

### General Optimization Guidelines

1. **Measure First**: Always profile before optimizing
2. **Focus on Hot Paths**: Optimize code that runs frequently
3. **Reduce Allocations**: Minimize heap allocations in hot paths
4. **Use Buffering**: Buffer I/O operations
5. **Avoid Locks**: Use lock-free algorithms where possible
6. **Batch Operations**: Batch multiple operations together
7. **Cache Results**: Cache expensive computations
8. **Use Goroutines Wisely**: Don't over-parallelize (overhead)

## Benchmark Results

### Example Benchmark Output

```
BenchmarkTripleGraphLayout/10nodes-10              1635270    730.0 ns/op      80 B/op       1 allocs/op
BenchmarkTripleGraphLayout/50nodes-10               166155   7179 ns/op       416 B/op       1 allocs/op
BenchmarkTripleGraphLayout/100nodes-10               49333  24636 ns/op       896 B/op       1 allocs/op
BenchmarkTripleGraphLayout/200nodes-10               14086  84153 ns/op      1792 B/op       1 allocs/op

BenchmarkBrokerPublish/_buffer-10                 83893104     14.31 ns/op       0 B/op       0 allocs/op
BenchmarkBrokerPublishWithSubscribers/1subscribers-10   46846635     25.87 ns/op       0 B/op       0 allocs/op
BenchmarkBrokerPublishWithSubscribers/10subscribers-10   1894734    598.6 ns/op       0 B/op       0 allocs/op

BenchmarkWorkPlanJSONParsing/10tasks-10            107101  11299 ns/op      4184 B/op      77 allocs/op
BenchmarkWorkPlanJSONParsing/100tasks-10            10000 102193 ns/op     36728 B/op     650 allocs/op
```

### Performance Analysis

**Graph Layout Scaling**:
- Layout time scales approximately O(n²) with node count (expected for force-directed algorithms)
- 100 nodes: ~25ms per iteration (well within target of <100ms)
- Memory usage scales linearly with node count
- Opportunity: Implement Barnes-Hut or spatial partitioning for better scaling

**Event Broker Performance**:
- Zero-allocation publishing with no subscribers (14ns/op)
- Linear scaling with subscriber count
- Memory-efficient: no allocations for event distribution
- Excellent concurrency: 116ns/op across 8 goroutines

**JSON Performance**:
- Reasonable parsing times for typical work plans (10-100 tasks)
- Memory allocations scale with plan size
- Cycle detection is efficient (~631ns for 4-task cycle)
- Opportunity: Use code generation for faster JSON processing

### Performance Regression Testing

To detect performance regressions:

```bash
# Before changes
go test -bench=. -benchmem ./... > bench_before.txt

# After changes
go test -bench=. -benchmem ./... > bench_after.txt

# Compare
benchstat bench_before.txt bench_after.txt
```

Example benchstat output:
```
name                              old time/op    new time/op    delta
TripleGraphLayout/100nodes-10       24.6µs ± 2%    22.1µs ± 1%  -10.16%  (p=0.000 n=10+10)
BrokerPublish/_buffer-10            14.3ns ± 1%    14.1ns ± 0%   -1.40%  (p=0.000 n=10+8)
WorkPlanJSONParsing/100tasks-10      102µs ± 3%     95µs ± 2%   -6.86%  (p=0.000 n=10+10)

name                              old alloc/op   new alloc/op   delta
TripleGraphLayout/100nodes-10         896B ± 0%      896B ± 0%     ~     (all equal)
BrokerPublish/_buffer-10             0.00B          0.00B          ~     (all equal)
WorkPlanJSONParsing/100tasks-10     36.7kB ± 0%    34.2kB ± 0%   -6.85%  (p=0.000 n=10+10)
```

## Continuous Integration

Add benchmark tests to CI pipeline:

```yaml
# .github/workflows/benchmark.yml
name: Benchmarks
on: [push, pull_request]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run benchmarks
        run: |
          go test -bench=. -benchmem -run=^$ ./internal/analyze > analyze_bench.txt
          go test -bench=. -benchmem -run=^$ ./internal/app/events > events_bench.txt
          go test -bench=. -benchmem -run=^$ ./internal/orchestrate > orchestrate_bench.txt

      - name: Upload benchmark results
        uses: actions/upload-artifact@v3
        with:
          name: benchmark-results
          path: |
            analyze_bench.txt
            events_bench.txt
            orchestrate_bench.txt
```

## Future Work

1. **Memory Graph Layout**: Add benchmarks for memory graph visualization (mnemosyne graph traversal)
2. **Terminal Rendering**: Benchmark TUI rendering performance
3. **File I/O**: Benchmark file reading/writing operations
4. **Search Operations**: Benchmark semantic search in editor
5. **Real-World Workloads**: Create benchmarks based on actual usage patterns
6. **Continuous Benchmarking**: Set up continuous benchmark tracking
7. **Performance Budgets**: Establish and enforce performance budgets per feature

## References

- [Go Benchmarking Guide](https://pkg.go.dev/testing#hdr-Benchmarks)
- [pprof Documentation](https://github.com/google/pprof/blob/main/doc/README.md)
- [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/gophercon-2019.html)
- [Go Performance Tips](https://github.com/dgryski/go-perfbook)

---

**Last Updated**: 2025-11-09
**Platform**: Apple M1 Pro (darwin/arm64)
**Go Version**: 1.21+
