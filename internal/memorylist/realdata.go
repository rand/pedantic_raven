package memorylist

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// LoadMemoriesFromServer creates a Bubble Tea command that loads memories
// from the mnemosyne server with the specified filters.
func LoadMemoriesFromServer(client *mnemosyne.Client, filters Filters) tea.Cmd {
	return func() tea.Msg {
		if client == nil || !client.IsConnected() {
			return MemoriesLoadedMsg{
				Memories:   nil,
				TotalCount: 0,
				Err:        mnemosyne.ErrNotConnected,
			}
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Build namespace filter if specified
		var namespace *pb.Namespace
		if filters.Namespace != "" {
			namespace = parseNamespaceString(filters.Namespace)
		}

		// Build minimum importance filter
		var minImportance *uint32
		if filters.MinImportance > 0 {
			imp := uint32(filters.MinImportance)
			minImportance = &imp
		}

		// Create list options
		opts := mnemosyne.ListMemoriesOptions{
			Namespace:       namespace,
			MaxResults:      100, // Default limit for initial load
			Tags:            filters.Tags,
			MinImportance:   minImportance,
			IncludeArchived: false,
		}

		// Load memories from server
		memories, err := client.ListMemories(ctx, opts)
		if err != nil {
			return MemoriesLoadedMsg{
				Memories:   nil,
				TotalCount: 0,
				Err:        err,
			}
		}

		// Apply client-side MaxImportance filter if needed
		if filters.MaxImportance > 0 {
			memories = filterByMaxImportance(memories, filters.MaxImportance)
		}

		return MemoriesLoadedMsg{
			Memories:   memories,
			TotalCount: uint32(len(memories)),
			Err:        nil,
		}
	}
}

// SearchMemories creates a Bubble Tea command that performs a semantic search
// using the Recall endpoint.
func SearchMemories(client *mnemosyne.Client, query string) tea.Cmd {
	return func() tea.Msg {
		if client == nil || !client.IsConnected() {
			return MemoriesLoadedMsg{
				Memories:   nil,
				TotalCount: 0,
				Err:        mnemosyne.ErrNotConnected,
			}
		}

		if query == "" {
			// Empty query - return empty result
			return MemoriesLoadedMsg{
				Memories:   []*pb.MemoryNote{},
				TotalCount: 0,
				Err:        nil,
			}
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Create recall options
		opts := mnemosyne.RecallOptions{
			Query:      query,
			MaxResults: 100, // Default limit
		}

		// Perform semantic search
		results, err := client.Recall(ctx, opts)
		if err != nil {
			return MemoriesLoadedMsg{
				Memories:   nil,
				TotalCount: 0,
				Err:        err,
			}
		}

		// Extract memories from search results
		memories := make([]*pb.MemoryNote, 0, len(results))
		for _, result := range results {
			if result.Memory != nil {
				memories = append(memories, result.Memory)
			}
		}

		return MemoriesLoadedMsg{
			Memories:   memories,
			TotalCount: uint32(len(memories)),
			Err:        nil,
		}
	}
}

// RefreshMemories creates a Bubble Tea command that refreshes the current
// memory list, invalidating any cached results.
func RefreshMemories(client *mnemosyne.Client) tea.Cmd {
	return func() tea.Msg {
		if client == nil || !client.IsConnected() {
			return MemoriesLoadedMsg{
				Memories:   nil,
				TotalCount: 0,
				Err:        mnemosyne.ErrNotConnected,
			}
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Load all memories without filters
		opts := mnemosyne.ListMemoriesOptions{
			MaxResults:      100,
			IncludeArchived: false,
		}

		memories, err := client.ListMemories(ctx, opts)
		if err != nil {
			return MemoriesLoadedMsg{
				Memories:   nil,
				TotalCount: 0,
				Err:        err,
			}
		}

		return MemoriesLoadedMsg{
			Memories:   memories,
			TotalCount: uint32(len(memories)),
			Err:        nil,
		}
	}
}

// Helper functions

// parseNamespaceString converts a namespace string to a protobuf Namespace.
// Supported formats:
//   - "global" -> GlobalNamespace
//   - "project:name" -> ProjectNamespace
//   - "project:name:session:id" -> SessionNamespace
func parseNamespaceString(ns string) *pb.Namespace {
	if ns == "global" {
		return mnemosyne.GlobalNamespace()
	}

	// Simple parsing - could be enhanced
	if len(ns) > 8 && ns[:8] == "project:" {
		projectName := ns[8:]
		return mnemosyne.ProjectNamespace(projectName)
	}

	// Default to treating as project namespace
	return mnemosyne.ProjectNamespace(ns)
}

// filterByMaxImportance filters memories to only include those with
// importance <= maxImportance.
func filterByMaxImportance(memories []*pb.MemoryNote, maxImportance int) []*pb.MemoryNote {
	filtered := make([]*pb.MemoryNote, 0, len(memories))
	for _, mem := range memories {
		if int(mem.Importance) <= maxImportance {
			filtered = append(filtered, mem)
		}
	}
	return filtered
}
