package mnemosyne

import (
	"testing"
	"time"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// TestStoreMemoryOptionsStructure verifies StoreMemoryOptions field setting
func TestStoreMemoryOptionsStructure(t *testing.T) {
	importance := uint32(8)
	memoryType := pb.MemoryType_MEMORY_TYPE_INSIGHT

	opts := StoreMemoryOptions{
		Content:           "test content",
		Namespace:         ProjectNamespace("test"),
		Importance:        &importance,
		Context:           "test context",
		Tags:              []string{"tag1", "tag2"},
		MemoryType:        &memoryType,
		SkipLLMEnrichment: true,
	}

	if opts.Content != "test content" {
		t.Errorf("expected content 'test content', got %q", opts.Content)
	}

	if *opts.Importance != 8 {
		t.Errorf("expected importance 8, got %d", *opts.Importance)
	}

	if len(opts.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(opts.Tags))
	}

	if !opts.SkipLLMEnrichment {
		t.Error("expected SkipLLMEnrichment to be true")
	}
}

// TestUpdateMemoryOptionsStructure verifies UpdateMemoryOptions field setting
func TestUpdateMemoryOptionsStructure(t *testing.T) {
	content := "updated content"
	importance := uint32(9)

	opts := UpdateMemoryOptions{
		MemoryID:   "mem-123",
		Content:    &content,
		Importance: &importance,
		Tags:       []string{"new-tag"},
		AddTags:    []string{"add-tag"},
		RemoveTags: []string{"remove-tag"},
	}

	if opts.MemoryID != "mem-123" {
		t.Errorf("expected MemoryID 'mem-123', got %q", opts.MemoryID)
	}

	if *opts.Content != "updated content" {
		t.Errorf("expected content 'updated content', got %q", *opts.Content)
	}

	if len(opts.AddTags) != 1 {
		t.Errorf("expected 1 add tag, got %d", len(opts.AddTags))
	}

	if len(opts.RemoveTags) != 1 {
		t.Errorf("expected 1 remove tag, got %d", len(opts.RemoveTags))
	}
}

// TestListMemoriesOptionsStructure verifies ListMemoriesOptions field setting
func TestListMemoriesOptionsStructure(t *testing.T) {
	minImportance := uint32(7)

	opts := ListMemoriesOptions{
		Namespace:       ProjectNamespace("test"),
		MemoryTypes:     []pb.MemoryType{pb.MemoryType_MEMORY_TYPE_INSIGHT},
		Tags:            []string{"tag1", "tag2"},
		MinImportance:   &minImportance,
		MaxResults:      50,
		IncludeArchived: true,
	}

	if opts.MaxResults != 50 {
		t.Errorf("expected MaxResults 50, got %d", opts.MaxResults)
	}

	if !opts.IncludeArchived {
		t.Error("expected IncludeArchived to be true")
	}

	if *opts.MinImportance != 7 {
		t.Errorf("expected MinImportance 7, got %d", *opts.MinImportance)
	}
}

// TestRecallOptionsWeights verifies RecallOptions weight settings
func TestRecallOptionsWeights(t *testing.T) {
	semanticWeight := float32(0.8)
	ftsWeight := float32(0.15)
	graphWeight := float32(0.05)

	opts := RecallOptions{
		Query:          "test",
		SemanticWeight: &semanticWeight,
		FtsWeight:      &ftsWeight,
		GraphWeight:    &graphWeight,
	}

	if *opts.SemanticWeight != 0.8 {
		t.Errorf("expected SemanticWeight 0.8, got %f", *opts.SemanticWeight)
	}

	if *opts.FtsWeight != 0.15 {
		t.Errorf("expected FtsWeight 0.15, got %f", *opts.FtsWeight)
	}

	if *opts.GraphWeight != 0.05 {
		t.Errorf("expected GraphWeight 0.05, got %f", *opts.GraphWeight)
	}
}

// TestSemanticSearchOptionsFields verifies SemanticSearchOptions field setting
func TestSemanticSearchOptionsFields(t *testing.T) {
	embedding := make([]float32, 768)
	minImportance := uint32(5)

	opts := SemanticSearchOptions{
		Embedding:       embedding,
		Namespace:       GlobalNamespace(),
		MaxResults:      15,
		MinImportance:   &minImportance,
		IncludeArchived: true,
	}

	if len(opts.Embedding) != 768 {
		t.Errorf("expected 768-dim embedding, got %d", len(opts.Embedding))
	}

	if opts.MaxResults != 15 {
		t.Errorf("expected MaxResults 15, got %d", opts.MaxResults)
	}
}

// TestGraphTraverseOptionsFields verifies GraphTraverseOptions field setting
func TestGraphTraverseOptionsFields(t *testing.T) {
	minStrength := float32(0.7)

	opts := GraphTraverseOptions{
		SeedIDs:         []string{"mem-1", "mem-2", "mem-3"},
		MaxHops:         3,
		MinLinkStrength: &minStrength,
	}

	if len(opts.SeedIDs) != 3 {
		t.Errorf("expected 3 seed IDs, got %d", len(opts.SeedIDs))
	}

	if opts.MaxHops != 3 {
		t.Errorf("expected MaxHops 3, got %d", opts.MaxHops)
	}

	if *opts.MinLinkStrength != 0.7 {
		t.Errorf("expected MinLinkStrength 0.7, got %f", *opts.MinLinkStrength)
	}
}

// TestGetContextOptionsFields verifies GetContextOptions field setting
func TestGetContextOptionsFields(t *testing.T) {
	opts := GetContextOptions{
		MemoryIDs:      []string{"mem-1", "mem-2"},
		IncludeLinks:   true,
		MaxLinkedDepth: 3,
	}

	if len(opts.MemoryIDs) != 2 {
		t.Errorf("expected 2 memory IDs, got %d", len(opts.MemoryIDs))
	}

	if !opts.IncludeLinks {
		t.Error("expected IncludeLinks to be true")
	}

	if opts.MaxLinkedDepth != 3 {
		t.Errorf("expected MaxLinkedDepth 3, got %d", opts.MaxLinkedDepth)
	}
}

// TestConfigFromEnvInvalidValues verifies environment parsing edge cases
func TestConfigFromEnvInvalidValues(t *testing.T) {
	t.Setenv("MNEMOSYNE_TIMEOUT", "invalid")
	t.Setenv("MNEMOSYNE_MAX_RETRIES", "invalid")

	config := ConfigFromEnv()

	// Should use defaults for invalid values
	if config.Timeout != 30*time.Second {
		t.Errorf("expected default timeout for invalid input, got %v", config.Timeout)
	}

	if config.MaxRetries != 3 {
		t.Errorf("expected default max retries for invalid input, got %d", config.MaxRetries)
	}
}

// TestConfigFromEnvNegativeValues verifies negative values are ignored
func TestConfigFromEnvNegativeValues(t *testing.T) {
	t.Setenv("MNEMOSYNE_TIMEOUT", "-10")
	t.Setenv("MNEMOSYNE_MAX_RETRIES", "-5")

	config := ConfigFromEnv()

	// Should use defaults for negative values
	if config.Timeout != 30*time.Second {
		t.Errorf("expected default timeout for negative input, got %v", config.Timeout)
	}

	if config.MaxRetries != 3 {
		t.Errorf("expected default max retries for negative input, got %d", config.MaxRetries)
	}
}

// TestConfigFromEnvZeroTimeout verifies zero timeout is ignored
func TestConfigFromEnvZeroTimeout(t *testing.T) {
	t.Setenv("MNEMOSYNE_TIMEOUT", "0")

	config := ConfigFromEnv()

	// Should use default for zero timeout
	if config.Timeout != 30*time.Second {
		t.Errorf("expected default timeout for zero input, got %v", config.Timeout)
	}
}

// TestNewClientWithZeroTimeout verifies zero timeout gets default
func TestNewClientWithZeroTimeout(t *testing.T) {
	cfg := Config{
		ServerAddr: "localhost:50051",
		Timeout:    0,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Timeout should be set to default
	if client.defaultCtx == nil {
		t.Error("expected default context to be set")
	}
}
