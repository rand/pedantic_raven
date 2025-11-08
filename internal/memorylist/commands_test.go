package memorylist

import (
	"testing"

	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// --- LoadOptions Tests ---

func TestDefaultLoadOptions(t *testing.T) {
	opts := DefaultLoadOptions()

	if opts.Namespace != nil {
		t.Error("Expected nil namespace for default options")
	}

	if opts.Limit != 50 {
		t.Errorf("Expected default limit 50, got %d", opts.Limit)
	}

	if opts.IncludeArchived {
		t.Error("Expected IncludeArchived to be false by default")
	}
}

func TestLoadWithNamespace(t *testing.T) {
	ns := &pb.Namespace{
		Namespace: &pb.Namespace_Global{
			Global: &pb.GlobalNamespace{},
		},
	}

	opts := LoadWithNamespace(ns)

	if opts.Namespace != ns {
		t.Error("Expected namespace to be set")
	}

	if opts.Limit != 50 {
		t.Error("Expected default limit to be preserved")
	}
}

func TestLoadWithFilters(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	minImportance := uint32(7)

	opts := LoadWithFilters(tags, minImportance)

	if len(opts.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(opts.Tags))
	}

	if opts.MinImportance == nil {
		t.Fatal("Expected MinImportance to be set")
	}

	if *opts.MinImportance != 7 {
		t.Errorf("Expected MinImportance 7, got %d", *opts.MinImportance)
	}
}

func TestLoadWithFiltersZeroImportance(t *testing.T) {
	opts := LoadWithFilters(nil, 0)

	if opts.MinImportance != nil {
		t.Error("Expected MinImportance to be nil when zero")
	}
}

// --- Command Tests ---

func TestLoadMemoriesCmd(t *testing.T) {
	// Test with nil client
	cmd := LoadMemoriesCmd(nil, DefaultLoadOptions())
	msg := cmd()

	if _, ok := msg.(MemoriesErrorMsg); !ok {
		t.Error("Expected MemoriesErrorMsg when client is nil")
	}
}

func TestSearchMemoriesCmd(t *testing.T) {
	// Test with nil client
	cmd := SearchMemoriesCmd(nil, "test query", DefaultLoadOptions())
	msg := cmd()

	if _, ok := msg.(MemoriesErrorMsg); !ok {
		t.Error("Expected MemoriesErrorMsg when client is nil")
	}
}

func TestSearchMemoriesCmdEmptyQuery(t *testing.T) {
	// Empty query should behave like LoadMemoriesCmd
	// This test just verifies it doesn't crash
	cmd := SearchMemoriesCmd(nil, "", DefaultLoadOptions())
	_ = cmd() // Just ensure it executes
}

func TestConnectAndLoadCmd(t *testing.T) {
	// Test with nil client
	cmd := ConnectAndLoadCmd(nil, DefaultLoadOptions())
	msg := cmd()

	if _, ok := msg.(MemoriesErrorMsg); !ok {
		t.Error("Expected MemoriesErrorMsg when client is nil")
	}
}

func TestInitCmd(t *testing.T) {
	// Test with nil client
	cmd := InitCmd(nil)
	msg := cmd()

	if _, ok := msg.(MemoriesErrorMsg); !ok {
		t.Error("Expected MemoriesErrorMsg when client is nil")
	}
}

// --- Model Integration Tests ---

func TestNewModelWithClient(t *testing.T) {
	client := &mnemosyne.Client{}
	m := NewModelWithClient(client)

	if m.client == nil {
		t.Error("Expected client to be set")
	}

	if !m.autoReload {
		t.Error("Expected autoReload to be true")
	}
}

func TestSetClient(t *testing.T) {
	m := NewModel()
	client := &mnemosyne.Client{}

	m.SetClient(client)

	if m.Client() == nil {
		t.Error("Expected client to be set")
	}
}

func TestLoadOptions(t *testing.T) {
	m := NewModel()
	opts := LoadWithFilters([]string{"tag1"}, 5)

	m.SetLoadOptions(opts)

	retrieved := m.LoadOptions()

	if len(retrieved.Tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(retrieved.Tags))
	}

	if retrieved.MinImportance == nil || *retrieved.MinImportance != 5 {
		t.Error("Expected MinImportance to be 5")
	}
}

func TestSearchQuery(t *testing.T) {
	m := NewModel()
	m.searchQuery = "test query"

	if m.SearchQuery() != "test query" {
		t.Errorf("Expected search query 'test query', got '%s'", m.SearchQuery())
	}
}

// --- Message Handling Tests ---

func TestHandleLoadRequestMsg(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	m, cmd := m.Update(LoadRequestMsg{})

	if !m.IsLoading() {
		t.Error("Expected loading state after LoadRequestMsg")
	}

	// Cmd should be nil since we don't have a client
	if cmd != nil {
		t.Error("Expected nil cmd when no client is set")
	}
}

func TestHandleReloadRequestMsg(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	m, cmd := m.Update(ReloadRequestMsg{})

	if !m.IsLoading() {
		t.Error("Expected loading state after ReloadRequestMsg")
	}

	if cmd != nil {
		t.Error("Expected nil cmd when no client is set")
	}
}

func TestInitWithClient(t *testing.T) {
	client := &mnemosyne.Client{}
	m := NewModelWithClient(client)

	cmd := m.Init()

	if cmd == nil {
		t.Error("Expected non-nil cmd from Init when client is set")
	}

	// Execute the cmd to see what message it produces
	msg := cmd()

	if _, ok := msg.(LoadRequestMsg); !ok {
		t.Errorf("Expected LoadRequestMsg from Init, got %T", msg)
	}
}

func TestInitWithoutClient(t *testing.T) {
	m := NewModel()

	cmd := m.Init()

	if cmd != nil {
		t.Error("Expected nil cmd from Init when no client is set")
	}
}
