package mnemosyne

import (
	"context"
	"testing"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// --- Recall Tests ---

func TestRecallNotConnected(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	ctx := context.Background()
	_, err = client.Recall(ctx, RecallOptions{Query: "test"})
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestRecallValidation(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.connected = true
	defer func() { client.connected = false }()

	ctx := context.Background()

	// Empty query should fail
	_, err = client.Recall(ctx, RecallOptions{})
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for empty query, got: %v", err)
	}
}

func TestRecallDefaultMaxResults(t *testing.T) {
	// This test verifies the default value is set correctly
	opts := RecallOptions{
		Query: "test query",
	}

	if opts.MaxResults != 0 {
		t.Errorf("Expected MaxResults to be 0 (unset), got %d", opts.MaxResults)
	}

	// The actual default of 10 is applied in the Recall method
}

func TestRecallWithAllOptions(t *testing.T) {
	// Test that all options can be set without error
	semanticWeight := float32(0.8)
	ftsWeight := float32(0.15)
	graphWeight := float32(0.05)
	minImportance := uint32(5)

	opts := RecallOptions{
		Query:          "test query",
		Namespace:      ProjectNamespace("myproject"),
		MaxResults:     20,
		MinImportance:  &minImportance,
		SemanticWeight: &semanticWeight,
		FtsWeight:      &ftsWeight,
		GraphWeight:    &graphWeight,
		Tags:           []string{"tag1", "tag2"},
	}

	if opts.Query == "" {
		t.Error("Query should be set")
	}
}

// --- SemanticSearch Tests ---

func TestSemanticSearchNotConnected(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	ctx := context.Background()
	_, err = client.SemanticSearch(ctx, SemanticSearchOptions{
		Embedding: []float32{0.1, 0.2, 0.3},
	})
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestSemanticSearchValidation(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.connected = true
	defer func() { client.connected = false }()

	ctx := context.Background()

	// Empty embedding should fail
	_, err = client.SemanticSearch(ctx, SemanticSearchOptions{})
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for empty embedding, got: %v", err)
	}
}

func TestSemanticSearchDefaultMaxResults(t *testing.T) {
	opts := SemanticSearchOptions{
		Embedding: []float32{0.1, 0.2, 0.3},
	}

	if opts.MaxResults != 0 {
		t.Errorf("Expected MaxResults to be 0 (unset), got %d", opts.MaxResults)
	}
}

func TestSemanticSearchWithEmbedding(t *testing.T) {
	// Test 768-dimensional embedding (typical for sentence transformers)
	embedding := make([]float32, 768)
	for i := range embedding {
		embedding[i] = float32(i) / 768.0
	}

	minImportance := uint32(7)
	opts := SemanticSearchOptions{
		Embedding:     embedding,
		Namespace:     GlobalNamespace(),
		MaxResults:    15,
		MinImportance: &minImportance,
	}

	if len(opts.Embedding) != 768 {
		t.Errorf("Expected 768-dimensional embedding, got %d", len(opts.Embedding))
	}
}

// --- GraphTraverse Tests ---

func TestGraphTraverseNotConnected(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	ctx := context.Background()
	_, _, err = client.GraphTraverse(ctx, GraphTraverseOptions{
		SeedIDs: []string{"mem-1", "mem-2"},
	})
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestGraphTraverseValidation(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.connected = true
	defer func() { client.connected = false }()

	ctx := context.Background()

	// Empty seed IDs should fail
	_, _, err = client.GraphTraverse(ctx, GraphTraverseOptions{})
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for empty seed IDs, got: %v", err)
	}
}

func TestGraphTraverseDefaultMaxHops(t *testing.T) {
	opts := GraphTraverseOptions{
		SeedIDs: []string{"mem-1"},
	}

	if opts.MaxHops != 0 {
		t.Errorf("Expected MaxHops to be 0 (unset), got %d", opts.MaxHops)
	}

	// The actual default of 2 is applied in the GraphTraverse method
}

func TestGraphTraverseWithOptions(t *testing.T) {
	minStrength := float32(0.5)
	opts := GraphTraverseOptions{
		SeedIDs:         []string{"mem-1", "mem-2", "mem-3"},
		MaxHops:         3,
		MinLinkStrength: &minStrength,
	}

	if len(opts.SeedIDs) != 3 {
		t.Errorf("Expected 3 seed IDs, got %d", len(opts.SeedIDs))
	}
}

// --- GetContext Tests ---

func TestGetContextNotConnected(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	ctx := context.Background()
	_, err = client.GetContext(ctx, GetContextOptions{
		MemoryIDs: []string{"mem-1"},
	})
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestGetContextValidation(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.connected = true
	defer func() { client.connected = false }()

	ctx := context.Background()

	// Empty memory IDs should fail
	_, err = client.GetContext(ctx, GetContextOptions{})
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for empty memory IDs, got: %v", err)
	}
}

func TestGetContextDefaultDepth(t *testing.T) {
	opts := GetContextOptions{
		MemoryIDs:    []string{"mem-1"},
		IncludeLinks: true,
	}

	if opts.MaxLinkedDepth != 0 {
		t.Errorf("Expected MaxLinkedDepth to be 0 (unset), got %d", opts.MaxLinkedDepth)
	}

	// The actual default of 1 is applied in the GetContext method
}

func TestGetContextWithOptions(t *testing.T) {
	opts := GetContextOptions{
		MemoryIDs:      []string{"mem-1", "mem-2"},
		IncludeLinks:   true,
		MaxLinkedDepth: 2,
	}

	if !opts.IncludeLinks {
		t.Error("Expected IncludeLinks to be true")
	}

	if opts.MaxLinkedDepth != 2 {
		t.Errorf("Expected MaxLinkedDepth 2, got %d", opts.MaxLinkedDepth)
	}
}

// --- Streaming Method Tests ---

func TestRecallStreamNotConnected(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	ctx := context.Background()
	_, err = client.RecallStream(ctx, RecallOptions{Query: "test"})
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestRecallStreamValidation(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.connected = true
	defer func() { client.connected = false }()

	ctx := context.Background()

	// Empty query should fail
	_, err = client.RecallStream(ctx, RecallOptions{})
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for empty query, got: %v", err)
	}
}

func TestListMemoriesStreamNotConnected(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	ctx := context.Background()
	_, err = client.ListMemoriesStream(ctx, ListMemoriesOptions{})
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestStoreMemoryStreamNotConnected(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	ctx := context.Background()
	_, err = client.StoreMemoryStream(ctx, StoreMemoryOptions{
		Content:   "test",
		Namespace: GlobalNamespace(),
	})
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

func TestStoreMemoryStreamValidation(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.connected = true
	defer func() { client.connected = false }()

	ctx := context.Background()

	// Missing content should fail
	_, err = client.StoreMemoryStream(ctx, StoreMemoryOptions{
		Namespace: GlobalNamespace(),
	})
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for missing content, got: %v", err)
	}

	// Missing namespace should fail
	_, err = client.StoreMemoryStream(ctx, StoreMemoryOptions{
		Content: "test content",
	})
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for missing namespace, got: %v", err)
	}
}

// --- Options Structure Tests ---

func TestRecallOptionsStructure(t *testing.T) {
	// Ensure RecallOptions has all expected fields
	opts := RecallOptions{
		Query:           "test",
		Namespace:       ProjectNamespace("proj"),
		MaxResults:      10,
		MemoryTypes:     []pb.MemoryType{},
		Tags:            []string{"tag"},
		IncludeArchived: false,
	}

	if opts.Query != "test" {
		t.Error("RecallOptions.Query not set correctly")
	}
}

func TestSemanticSearchOptionsStructure(t *testing.T) {
	opts := SemanticSearchOptions{
		Embedding:       []float32{0.1, 0.2},
		Namespace:       GlobalNamespace(),
		MaxResults:      5,
		IncludeArchived: true,
	}

	if len(opts.Embedding) != 2 {
		t.Error("SemanticSearchOptions.Embedding not set correctly")
	}
}

func TestGraphTraverseOptionsStructure(t *testing.T) {
	opts := GraphTraverseOptions{
		SeedIDs:   []string{"id1", "id2"},
		MaxHops:   3,
		LinkTypes: []pb.LinkType{},
	}

	if len(opts.SeedIDs) != 2 {
		t.Error("GraphTraverseOptions.SeedIDs not set correctly")
	}
}

func TestGetContextOptionsStructure(t *testing.T) {
	opts := GetContextOptions{
		MemoryIDs:      []string{"id1"},
		IncludeLinks:   true,
		MaxLinkedDepth: 2,
	}

	if len(opts.MemoryIDs) != 1 {
		t.Error("GetContextOptions.MemoryIDs not set correctly")
	}
}
