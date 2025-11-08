// Package palette provides a command palette with fuzzy search.
//
// The command palette allows users to:
// - Discover available commands (Ctrl+K)
// - Search with fuzzy matching
// - Execute commands by name
// - View keybindings and descriptions
//
// Commands are organized by category and can be mode-specific.
package palette

import (
	tea "github.com/charmbracelet/bubbletea"
)

// CommandID uniquely identifies a command.
type CommandID string

// Category groups related commands.
type Category string

const (
	CategoryFile        Category = "File"
	CategoryEdit        Category = "Edit"
	CategoryView        Category = "View"
	CategoryMode        Category = "Mode"
	CategoryMemory      Category = "Memory"
	CategoryOrchestrate Category = "Orchestrate"
	CategoryHelp        Category = "Help"
)

// Command represents an executable command with metadata.
type Command struct {
	ID          CommandID
	Name        string
	Description string
	Keybinding  string // e.g., "Ctrl+S", "Ctrl+K"
	Category    Category
	Execute     func() tea.Cmd
}

// CommandRegistry manages the collection of available commands.
//
// Commands can be registered globally or for specific modes.
type CommandRegistry struct {
	commands map[CommandID]Command
}

// NewCommandRegistry creates a new command registry.
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[CommandID]Command),
	}
}

// Register adds a command to the registry.
//
// If a command with the same ID exists, it will be replaced.
func (r *CommandRegistry) Register(cmd Command) {
	r.commands[cmd.ID] = cmd
}

// Unregister removes a command from the registry.
func (r *CommandRegistry) Unregister(id CommandID) {
	delete(r.commands, id)
}

// Get retrieves a command by ID.
// Returns nil-Command if not found.
func (r *CommandRegistry) Get(id CommandID) (Command, bool) {
	cmd, ok := r.commands[id]
	return cmd, ok
}

// All returns all registered commands.
func (r *CommandRegistry) All() []Command {
	commands := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// ByCategory returns all commands in a specific category.
func (r *CommandRegistry) ByCategory(category Category) []Command {
	var commands []Command
	for _, cmd := range r.commands {
		if cmd.Category == category {
			commands = append(commands, cmd)
		}
	}
	return commands
}

// Count returns the number of registered commands.
func (r *CommandRegistry) Count() int {
	return len(r.commands)
}

// --- Fuzzy Matching ---

// MatchScore represents how well a command matches a query.
type MatchScore struct {
	Command Command
	Score   int // Higher is better
}

// FuzzyMatch finds commands matching the query using fuzzy matching.
//
// Scoring:
// - Exact name match: +100
// - Name contains query: +50
// - Description contains query: +20
// - Category match: +10
//
// Returns matches sorted by score (highest first).
func (r *CommandRegistry) FuzzyMatch(query string) []MatchScore {
	if query == "" {
		// No query - return all commands with score 0
		matches := make([]MatchScore, 0, len(r.commands))
		for _, cmd := range r.commands {
			matches = append(matches, MatchScore{
				Command: cmd,
				Score:   0,
			})
		}
		return matches
	}

	// Lowercase for case-insensitive matching
	queryLower := toLower(query)

	var matches []MatchScore
	for _, cmd := range r.commands {
		score := 0

		nameLower := toLower(cmd.Name)
		descLower := toLower(cmd.Description)
		categoryLower := toLower(string(cmd.Category))

		// Exact name match
		if nameLower == queryLower {
			score += 100
		}

		// Name contains query
		if contains(nameLower, queryLower) {
			score += 50
		}

		// Description contains query
		if contains(descLower, queryLower) {
			score += 20
		}

		// Category contains query
		if contains(categoryLower, queryLower) {
			score += 10
		}

		// Subsequence matching (fuzzy)
		if subsequenceMatch(nameLower, queryLower) {
			score += 30
		}

		if score > 0 {
			matches = append(matches, MatchScore{
				Command: cmd,
				Score:   score,
			})
		}
	}

	// Sort by score (highest first)
	sortByScore(matches)

	return matches
}

// --- Helper Functions ---

func toLower(s string) string {
	// Simple ASCII lowercasing
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}

	return false
}

// subsequenceMatch checks if query is a subsequence of s.
// Example: "fop" matches "file open" (f, o, p appear in order)
func subsequenceMatch(s, query string) bool {
	if len(query) == 0 {
		return true
	}
	if len(query) > len(s) {
		return false
	}

	queryIdx := 0
	for i := 0; i < len(s) && queryIdx < len(query); i++ {
		if s[i] == query[queryIdx] {
			queryIdx++
		}
	}

	return queryIdx == len(query)
}

// sortByScore sorts matches by score (highest first), then by name.
func sortByScore(matches []MatchScore) {
	// Simple bubble sort (fine for small lists)
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			// Sort by score descending
			if matches[i].Score < matches[j].Score {
				matches[i], matches[j] = matches[j], matches[i]
			} else if matches[i].Score == matches[j].Score {
				// Same score - sort by name ascending
				if matches[i].Command.Name > matches[j].Command.Name {
					matches[i], matches[j] = matches[j], matches[i]
				}
			}
		}
	}
}
