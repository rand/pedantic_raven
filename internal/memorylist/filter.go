package memorylist

import (
	"regexp"
	"strings"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// FilterOptions contains all filter criteria for memory filtering.
type FilterOptions struct {
	SearchQuery    string   // Text search query
	CaseSensitive  bool     // Whether search is case-sensitive
	UseRegex       bool     // Whether to use regex for search
	Tags           []string // Tags to filter by
	TagsAND        bool     // If true, ALL tags must match; if false, ANY tag matches
	MinImportance  int32    // Minimum importance level
	MaxImportance  int32    // Maximum importance level
	NamespaceType  string   // "global", "project", "session", or ""
	NamespaceName  string   // Specific namespace name (for project/session)
}

// SearchMemoriesLocal searches memory content with optional regex support.
// This is client-side filtering distinct from server-side search.
// Returns memories that match the search query.
func SearchMemoriesLocal(memories []*pb.MemoryNote, query string, caseSensitive bool, useRegex bool) []*pb.MemoryNote {
	if query == "" {
		return memories
	}

	var results []*pb.MemoryNote
	for _, mem := range memories {
		if matchesQuery(mem.Content, query, caseSensitive, useRegex) {
			results = append(results, mem)
		}
	}
	return results
}

// matchesQuery checks if content matches the query string or regex.
func matchesQuery(content, query string, caseSensitive bool, useRegex bool) bool {
	if !caseSensitive {
		content = strings.ToLower(content)
		query = strings.ToLower(query)
	}

	if useRegex {
		re, err := regexp.Compile(query)
		if err != nil {
			return false // Invalid regex
		}
		return re.MatchString(content)
	}

	return strings.Contains(content, query)
}

// FilterByTags filters memories by tags with AND/OR logic.
// If andLogic is true, memory must have ALL filter tags.
// If andLogic is false, memory must have AT LEAST ONE filter tag.
func FilterByTags(memories []*pb.MemoryNote, tags []string, andLogic bool) []*pb.MemoryNote {
	if len(tags) == 0 {
		return memories
	}

	var results []*pb.MemoryNote
	for _, mem := range memories {
		if matchesTags(mem.Tags, tags, andLogic) {
			results = append(results, mem)
		}
	}
	return results
}

// matchesTags checks if memory tags match the filter tags.
func matchesTags(memoryTags, filterTags []string, andLogic bool) bool {
	if andLogic {
		// Memory must have ALL filter tags
		for _, filterTag := range filterTags {
			found := false
			for _, memTag := range memoryTags {
				if memTag == filterTag {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	} else {
		// Memory must have AT LEAST ONE filter tag
		for _, filterTag := range filterTags {
			for _, memTag := range memoryTags {
				if memTag == filterTag {
					return true
				}
			}
		}
		return false
	}
}

// FilterByImportance filters memories by importance range (inclusive).
func FilterByImportance(memories []*pb.MemoryNote, minImportance, maxImportance int32) []*pb.MemoryNote {
	var results []*pb.MemoryNote
	for _, mem := range memories {
		if mem.Importance >= uint32(minImportance) && mem.Importance <= uint32(maxImportance) {
			results = append(results, mem)
		}
	}
	return results
}

// FilterByNamespace filters memories by namespace type and optionally name.
// nsType can be "global", "project", "session", or "" (no filter).
// If nsName is empty, matches any namespace of the specified type.
func FilterByNamespace(memories []*pb.MemoryNote, nsType string, nsName string) []*pb.MemoryNote {
	if nsType == "" {
		return memories // No namespace filter
	}

	var results []*pb.MemoryNote
	for _, mem := range memories {
		if matchesNamespaceFilter(mem.Namespace, nsType, nsName) {
			results = append(results, mem)
		}
	}
	return results
}

// matchesNamespaceFilter checks if a memory's namespace matches the filter criteria.
func matchesNamespaceFilter(ns *pb.Namespace, nsType string, nsName string) bool {
	if ns == nil {
		return nsType == ""
	}

	switch nsType {
	case "global":
		return ns.GetGlobal() != nil

	case "project":
		project := ns.GetProject()
		if project == nil {
			return false
		}
		if nsName == "" {
			return true // Any project matches
		}
		return project.Name == nsName

	case "session":
		session := ns.GetSession()
		if session == nil {
			return false
		}
		if nsName == "" {
			return true // Any session matches
		}
		return session.SessionId == nsName

	default:
		return false
	}
}

// FilterMemories applies all filters in sequence to a memory list.
// Filters are applied in this order: search, tags, importance, namespace.
// This order ensures efficient filtering (search is typically most selective).
func FilterMemories(memories []*pb.MemoryNote, opts FilterOptions) []*pb.MemoryNote {
	result := memories

	// Apply text search first (typically most selective)
	if opts.SearchQuery != "" {
		result = SearchMemoriesLocal(result, opts.SearchQuery, opts.CaseSensitive, opts.UseRegex)
	}

	// Apply tag filter
	if len(opts.Tags) > 0 {
		result = FilterByTags(result, opts.Tags, opts.TagsAND)
	}

	// Apply importance filter
	if opts.MinImportance > 0 || opts.MaxImportance < 10 {
		result = FilterByImportance(result, opts.MinImportance, opts.MaxImportance)
	}

	// Apply namespace filter
	if opts.NamespaceType != "" {
		result = FilterByNamespace(result, opts.NamespaceType, opts.NamespaceName)
	}

	return result
}

// CombinedFilter represents a combined filter operation that can be applied to memories.
// It provides a fluent interface for building complex filters.
type CombinedFilter struct {
	opts FilterOptions
}

// NewFilter creates a new filter builder.
func NewFilter() *CombinedFilter {
	return &CombinedFilter{
		opts: FilterOptions{
			MinImportance: 0,
			MaxImportance: 10,
		},
	}
}

// WithSearch sets the text search query.
func (f *CombinedFilter) WithSearch(query string, caseSensitive, useRegex bool) *CombinedFilter {
	f.opts.SearchQuery = query
	f.opts.CaseSensitive = caseSensitive
	f.opts.UseRegex = useRegex
	return f
}

// WithTags sets the tag filter (AND logic by default).
func (f *CombinedFilter) WithTags(tags []string, andLogic bool) *CombinedFilter {
	f.opts.Tags = tags
	f.opts.TagsAND = andLogic
	return f
}

// WithImportance sets the importance range.
func (f *CombinedFilter) WithImportance(min, max int32) *CombinedFilter {
	f.opts.MinImportance = min
	f.opts.MaxImportance = max
	return f
}

// WithNamespace sets the namespace filter.
func (f *CombinedFilter) WithNamespace(nsType, nsName string) *CombinedFilter {
	f.opts.NamespaceType = nsType
	f.opts.NamespaceName = nsName
	return f
}

// Apply applies the filter to the given memories.
func (f *CombinedFilter) Apply(memories []*pb.MemoryNote) []*pb.MemoryNote {
	return FilterMemories(memories, f.opts)
}

// Options returns the current filter options.
func (f *CombinedFilter) Options() FilterOptions {
	return f.opts
}

// Reset clears all filters.
func (f *CombinedFilter) Reset() *CombinedFilter {
	f.opts = FilterOptions{
		MinImportance: 0,
		MaxImportance: 10,
	}
	return f
}
