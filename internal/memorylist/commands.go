package memorylist

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// LoadMemoriesCmd creates a command to load memories from mnemosyne.
func LoadMemoriesCmd(client *mnemosyne.Client, opts LoadOptions) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return MemoriesErrorMsg{
				Err: mnemosyne.ErrNotConnected,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Build list options
		listOpts := mnemosyne.ListMemoriesOptions{
			Namespace:       opts.Namespace,
			MaxResults:      opts.Limit,
			MemoryTypes:     opts.MemoryTypes,
			Tags:            opts.Tags,
			MinImportance:   opts.MinImportance,
			IncludeArchived: opts.IncludeArchived,
		}

		memories, err := client.ListMemories(ctx, listOpts)
		if err != nil {
			return MemoriesErrorMsg{Err: err}
		}

		return MemoriesLoadedMsg{
			Memories:   memories,
			TotalCount: uint32(len(memories)), // TODO: Get actual total from response
		}
	}
}

// SearchMemoriesCmd creates a command to search memories using Recall.
func SearchMemoriesCmd(client *mnemosyne.Client, query string, opts LoadOptions) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return MemoriesErrorMsg{
				Err: mnemosyne.ErrNotConnected,
			}
		}

		if query == "" {
			// Empty query, just list memories
			return LoadMemoriesCmd(client, opts)()
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Build recall options
		recallOpts := mnemosyne.RecallOptions{
			Query:           query,
			Namespace:       opts.Namespace,
			MaxResults:      opts.Limit,
			MemoryTypes:     opts.MemoryTypes,
			Tags:            opts.Tags,
			MinImportance:   opts.MinImportance,
			IncludeArchived: opts.IncludeArchived,
		}

		results, err := client.Recall(ctx, recallOpts)
		if err != nil {
			return MemoriesErrorMsg{Err: err}
		}

		// Extract memories from search results
		memories := make([]*pb.MemoryNote, len(results))
		for i, result := range results {
			memories[i] = result.Memory
		}

		return MemoriesLoadedMsg{
			Memories:   memories,
			TotalCount: uint32(len(memories)),
		}
	}
}

// RefreshMemoriesCmd creates a command to refresh the current memory list.
func RefreshMemoriesCmd(client *mnemosyne.Client, currentOpts LoadOptions, searchQuery string) tea.Cmd {
	if searchQuery != "" {
		return SearchMemoriesCmd(client, searchQuery, currentOpts)
	}
	return LoadMemoriesCmd(client, currentOpts)
}

// LoadOptions holds options for loading memories.
type LoadOptions struct {
	Namespace       *pb.Namespace
	Limit           uint32
	MemoryTypes     []pb.MemoryType
	Tags            []string
	MinImportance   *uint32
	IncludeArchived bool
}

// DefaultLoadOptions returns default loading options.
func DefaultLoadOptions() LoadOptions {
	return LoadOptions{
		Namespace:       nil, // Load all namespaces
		Limit:           50,  // Default page size
		IncludeArchived: false,
	}
}

// LoadWithNamespace creates load options for a specific namespace.
func LoadWithNamespace(ns *pb.Namespace) LoadOptions {
	opts := DefaultLoadOptions()
	opts.Namespace = ns
	return opts
}

// LoadWithFilters creates load options with filters applied.
func LoadWithFilters(tags []string, minImportance uint32) LoadOptions {
	opts := DefaultLoadOptions()
	opts.Tags = tags
	if minImportance > 0 {
		opts.MinImportance = &minImportance
	}
	return opts
}

// ConnectAndLoadCmd creates a command to connect to mnemosyne and load memories.
func ConnectAndLoadCmd(client *mnemosyne.Client, opts LoadOptions) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return MemoriesErrorMsg{
				Err: mnemosyne.ErrNotConnected,
			}
		}

		// Check if already connected
		if !client.IsConnected() {
			// Attempt to connect
			if err := client.Connect(); err != nil {
				return MemoriesErrorMsg{Err: err}
			}
		}

		// Now load memories
		return LoadMemoriesCmd(client, opts)()
	}
}

// InitCmd returns the initial command to load memories on component init.
func InitCmd(client *mnemosyne.Client) tea.Cmd {
	return ConnectAndLoadCmd(client, DefaultLoadOptions())
}
