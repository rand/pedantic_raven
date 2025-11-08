package memorylist

import (
	"strings"
	"testing"
	"time"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// --- View Rendering Tests ---

func TestView(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	view := m.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain header
	if !strings.Contains(view, "Memory Workspace") {
		t.Error("Expected view to contain 'Memory Workspace'")
	}
}

func TestViewLoading(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)
	m.SetLoading(true)

	view := m.View()

	if !strings.Contains(view, "Loading memories") {
		t.Error("Expected view to contain 'Loading memories'")
	}
}

func TestViewError(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)
	m.SetError(ErrTestError)

	view := m.View()

	if !strings.Contains(view, "Error loading memories") {
		t.Error("Expected view to contain 'Error loading memories'")
	}
}

func TestViewEmpty(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)
	m.SetMemories([]*pb.MemoryNote{}, 0)

	view := m.View()

	if !strings.Contains(view, "No memories found") {
		t.Error("Expected view to contain 'No memories found'")
	}
}

func TestViewWithMemories(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)
	memories := []*pb.MemoryNote{
		createTestMemory("1", "Test memory content", 8, []string{"tag1"}, time.Hour),
	}
	m.SetMemories(memories, 1)

	view := m.View()

	// Should contain memory content
	if !strings.Contains(view, "Test memory") {
		t.Error("Expected view to contain memory content")
	}

	// Should contain importance
	if !strings.Contains(view, "[Imp: 8]") {
		t.Error("Expected view to contain importance indicator")
	}
}

func TestViewSelectedIndicator(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First memory", 8, nil, time.Hour),
		createTestMemory("2", "Second memory", 6, nil, 2*time.Hour),
	}
	m.SetMemories(memories, 2)

	view := m.View()

	// First memory should have selection indicator
	lines := strings.Split(view, "\n")
	var hasSelector bool
	for _, line := range lines {
		if strings.Contains(line, "First memory") && strings.Contains(line, ">") {
			hasSelector = true
			break
		}
	}

	if !hasSelector {
		t.Error("Expected view to show selection indicator")
	}
}

func TestViewFooterStats(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
		createTestMemory("2", "Second", 6, nil, 2*time.Hour),
	}
	m.SetMemories(memories, 10)

	view := m.View()

	// Should show memory count
	if !strings.Contains(view, "Showing") {
		t.Error("Expected view to contain memory count")
	}

	// Should show total count
	if !strings.Contains(view, "total: 10") {
		t.Error("Expected view to show total count")
	}
}

// --- Helper Function Tests ---

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Simple content",
			content:  "This is a title",
			expected: "This is a title",
		},
		{
			name:     "Markdown heading",
			content:  "# Main Title",
			expected: "Main Title",
		},
		{
			name:     "Multiple lines",
			content:  "First line title\nSecond line content",
			expected: "First line title",
		},
		{
			name:     "Empty content",
			content:  "",
			expected: "(Untitled)",
		},
		{
			name:     "Whitespace only",
			content:  "   \n   ",
			expected: "(Untitled)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTitle(tt.content)
			if result != tt.expected {
				t.Errorf("Expected title '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestFormatRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Just now",
			time:     now.Add(-30 * time.Second),
			expected: "just now",
		},
		{
			name:     "1 minute ago",
			time:     now.Add(-time.Minute),
			expected: "1 minute ago",
		},
		{
			name:     "30 minutes ago",
			time:     now.Add(-30 * time.Minute),
			expected: "30 minutes ago",
		},
		{
			name:     "1 hour ago",
			time:     now.Add(-time.Hour),
			expected: "1 hour ago",
		},
		{
			name:     "5 hours ago",
			time:     now.Add(-5 * time.Hour),
			expected: "5 hours ago",
		},
		{
			name:     "Yesterday",
			time:     now.Add(-25 * time.Hour),
			expected: "yesterday",
		},
		{
			name:     "3 days ago",
			time:     now.Add(-3 * 24 * time.Hour),
			expected: "3 days ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatRelativeTime(tt.time)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestFormatNamespace(t *testing.T) {
	tests := []struct {
		name      string
		namespace *pb.Namespace
		expected  string
	}{
		{
			name:      "Nil namespace",
			namespace: nil,
			expected:  "",
		},
		{
			name: "Global namespace",
			namespace: &pb.Namespace{
				Namespace: &pb.Namespace_Global{
					Global: &pb.GlobalNamespace{},
				},
			},
			expected: "global",
		},
		{
			name: "Project namespace",
			namespace: &pb.Namespace{
				Namespace: &pb.Namespace_Project{
					Project: &pb.ProjectNamespace{
						Name: "myapp",
					},
				},
			},
			expected: "project:myapp",
		},
		{
			name: "Session namespace",
			namespace: &pb.Namespace{
				Namespace: &pb.Namespace_Session{
					Session: &pb.SessionNamespace{
						Project:   "myapp",
						SessionId: "sess-123",
					},
				},
			},
			expected: "project:myapp:session:sess-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatNamespace(tt.namespace)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestRenderImportance(t *testing.T) {
	m := NewModel()

	tests := []struct {
		name       string
		importance uint32
		shouldFail bool
	}{
		{"Low importance", 1, false},
		{"Medium importance", 5, false},
		{"High importance", 8, false},
		{"Critical importance", 9, false},
		{"Max importance", 10, false},
		{"Invalid zero", 0, true},
		{"Invalid high", 11, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.renderImportance(tt.importance)

			if tt.shouldFail {
				if !strings.Contains(result, "[Imp: -]") {
					t.Errorf("Expected invalid importance indicator, got '%s'", result)
				}
			} else {
				if !strings.Contains(result, "[Imp:") {
					t.Errorf("Expected importance indicator, got '%s'", result)
				}
			}
		})
	}
}

// --- Rendering Component Tests ---

func TestRenderHeader(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	header := m.renderHeader()

	if !strings.Contains(header, "Memory Workspace") {
		t.Error("Expected header to contain 'Memory Workspace'")
	}

	if !strings.Contains(header, "Sort:") {
		t.Error("Expected header to contain sort indicator")
	}
}

func TestRenderHeaderWithSearch(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)
	m.SetSearchQuery("test query")

	header := m.renderHeader()

	if !strings.Contains(header, "test query") {
		t.Error("Expected header to contain search query")
	}
}

func TestRenderFooter(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
	}
	m.SetMemories(memories, 1)

	footer := m.renderFooter()

	if !strings.Contains(footer, "Showing") {
		t.Error("Expected footer to contain 'Showing'")
	}

	if !strings.Contains(footer, "j/k") {
		t.Error("Expected footer to contain key hints")
	}
}

func TestRenderFooterEmpty(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	footer := m.renderFooter()

	if !strings.Contains(footer, "0 memories") {
		t.Error("Expected footer to show 0 memories")
	}
}

func TestRenderMemoryRow(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)
	mem := createTestMemory("1", "# Test Title\nContent here", 8, []string{"tag1", "tag2"}, time.Hour)

	row := m.renderMemoryRow(mem, false, 0)

	// Should be 3 lines
	lines := strings.Split(row, "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// First line should contain title and importance
	if !strings.Contains(lines[0], "Test Title") {
		t.Error("Expected title in first line")
	}

	if !strings.Contains(lines[0], "[Imp: 8]") {
		t.Error("Expected importance in first line")
	}

	// Second line should contain tags
	if !strings.Contains(lines[1], "tag1") || !strings.Contains(lines[1], "tag2") {
		t.Error("Expected tags in second line")
	}

	// Third line should contain timestamp
	if !strings.Contains(lines[2], "Updated") {
		t.Error("Expected timestamp in third line")
	}
}

func TestVisibleLines(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	visible := m.visibleLines()

	// Height 10 - header (1) - footer (1) = 8
	if visible != 8 {
		t.Errorf("Expected 8 visible lines, got %d", visible)
	}
}
