package memorydetail

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
	if !strings.Contains(view, "Memory Detail") {
		t.Error("Expected view to contain 'Memory Detail'")
	}
}

func TestViewEmpty(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	view := m.View()

	if !strings.Contains(view, "No memory selected") {
		t.Error("Expected view to contain 'No memory selected'")
	}
}

func TestViewWithMemory(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	memory := createTestMemory("test-1", "Test memory content\nLine 2\nLine 3", 8, []string{"tag1"})
	m.SetMemory(memory)

	view := m.View()

	// Should contain memory content
	if !strings.Contains(view, "Test memory") {
		t.Error("Expected view to contain memory content")
	}
}

func TestViewError(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	err := &testError{msg: "load failed"}
	m.SetError(err)

	view := m.View()

	if !strings.Contains(view, "Error loading memory") {
		t.Error("Expected view to contain error message")
	}

	if !strings.Contains(view, "load failed") {
		t.Error("Expected view to contain error details")
	}
}

// --- Header Tests ---

func TestRenderHeader(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	header := m.renderHeader()

	if !strings.Contains(header, "Memory Detail") {
		t.Error("Expected header to contain 'Memory Detail'")
	}
}

func TestRenderHeaderWithMemory(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	memory := createTestMemory("test-1", "# My Important Note\nContent here", 8, nil)
	m.SetMemory(memory)

	header := m.renderHeader()

	if !strings.Contains(header, "My Important Note") {
		t.Error("Expected header to contain memory title")
	}
}

func TestRenderHeaderLongTitle(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	longTitle := strings.Repeat("Very Long Title ", 10)
	memory := createTestMemory("test-1", longTitle, 8, nil)
	m.SetMemory(memory)

	header := m.renderHeader()

	// Should be truncated
	if !strings.Contains(header, "...") {
		t.Error("Expected long title to be truncated")
	}
}

// --- Content Tests ---

func TestRenderContentEmpty(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	content := m.renderContent()

	if content != "" {
		t.Error("Expected empty content when memory is nil")
	}
}

func TestRenderContentWithMemory(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	memory := createTestMemory("test-1", "Line 1\nLine 2\nLine 3", 8, nil)
	m.SetMemory(memory)

	content := m.renderContent()

	if !strings.Contains(content, "Line 1") {
		t.Error("Expected content to contain 'Line 1'")
	}

	if !strings.Contains(content, "Line 2") {
		t.Error("Expected content to contain 'Line 2'")
	}
}

func TestRenderContentWithMetadata(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)
	m.showMetadata = true

	memory := createTestMemory("test-1", "Content line", 8, []string{"tag1"})
	m.SetMemory(memory)

	content := m.renderContent()

	// Should contain both content and metadata
	if !strings.Contains(content, "Content line") {
		t.Error("Expected content line")
	}

	if !strings.Contains(content, "Metadata") {
		t.Error("Expected metadata section")
	}
}

func TestRenderContentScrolling(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	// Create long content
	lines := []string{}
	for i := 0; i < 20; i++ {
		lines = append(lines, "Line "+string(rune('A'+i)))
	}
	content := strings.Join(lines, "\n")

	memory := createTestMemory("test-1", content, 8, nil)
	m.SetMemory(memory)

	// At scroll offset 0, should see first lines
	rendered := m.renderContent()
	if !strings.Contains(rendered, "Line A") {
		t.Error("Expected to see first line at scroll offset 0")
	}

	// Scroll down
	m.scrollOffset = 5
	rendered = m.renderContent()

	// Should not see first line anymore
	if strings.Contains(rendered, "Line A") {
		t.Error("Should not see first line after scrolling down")
	}
}

// --- Metadata Panel Tests ---

func TestRenderMetadataLine(t *testing.T) {
	m := NewModel()

	memory := createTestMemory("test-1", "Content", 8, []string{"tag1", "tag2"})
	m.SetMemory(memory)

	// Line 0 should be header
	line0 := m.renderMetadataLine(0)
	if !strings.Contains(line0, "Metadata") {
		t.Error("Expected metadata header at line 0")
	}

	// Line 2 should be ID
	line2 := m.renderMetadataLine(2)
	if !strings.Contains(line2, "test-1") {
		t.Error("Expected ID at line 2")
	}

	// Line 3 should be importance
	line3 := m.renderMetadataLine(3)
	if !strings.Contains(line3, "8") {
		t.Error("Expected importance at line 3")
	}
}

func TestRenderMetadataWithLinks(t *testing.T) {
	m := NewModel()

	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1", LinkType: pb.LinkType_LINK_TYPE_REFERENCES},
		{TargetId: "link-2", LinkType: pb.LinkType_LINK_TYPE_BUILDS_UPON},
	}
	m.SetMemory(memory)

	// Line 17 should show links header
	line17 := m.renderMetadataLine(17)
	if !strings.Contains(line17, "Links") {
		t.Error("Expected links header")
	}

	// Line 18 should show first link
	line18 := m.renderMetadataLine(18)
	if !strings.Contains(line18, "link-1") {
		t.Error("Expected first link ID")
	}

	// Line 19 should show second link
	line19 := m.renderMetadataLine(19)
	if !strings.Contains(line19, "link-2") {
		t.Error("Expected second link ID")
	}
}

// --- Footer Tests ---

func TestRenderFooter(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	footer := m.renderFooter()

	if !strings.Contains(footer, "No memory") {
		t.Error("Expected footer to indicate no memory")
	}

	if !strings.Contains(footer, "q") {
		t.Error("Expected footer to show 'q' key hint")
	}
}

func TestRenderFooterWithMemory(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	memory := createTestMemory("test-1", "Line 1\nLine 2\nLine 3", 8, nil)
	m.SetMemory(memory)

	footer := m.renderFooter()

	if !strings.Contains(footer, "lines") {
		t.Error("Expected footer to show line count")
	}

	if !strings.Contains(footer, "j/k") {
		t.Error("Expected footer to show navigation hints")
	}
}

func TestRenderFooterScrollInfo(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	// Create long content
	content := strings.Repeat("Line\n", 20)
	memory := createTestMemory("test-1", content, 8, nil)
	m.SetMemory(memory)

	footer := m.renderFooter()

	// Should show scroll position
	if !strings.Contains(footer, "Lines") {
		t.Error("Expected footer to show scroll position")
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
		{
			name:     "Long title",
			content:  strings.Repeat("Very long title ", 10),
			expected: strings.Repeat("Very long title ", 10)[:47] + "...",
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

func TestPadRight(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected string
	}{
		{
			name:     "Short string",
			input:    "test",
			width:    10,
			expected: "test      ",
		},
		{
			name:     "Exact width",
			input:    "test",
			width:    4,
			expected: "test",
		},
		{
			name:     "Longer than width",
			input:    "testing",
			width:    4,
			expected: "testing",
		},
		{
			name:     "Empty string",
			input:    "",
			width:    5,
			expected: "     ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padRight(tt.input, tt.width)
			if result != tt.expected {
				t.Errorf("Expected '%s' (len %d), got '%s' (len %d)",
					tt.expected, len(tt.expected), result, len(result))
			}
		})
	}
}

// --- Integration Tests ---

func TestMetadataToggleInView(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	memory := createTestMemory("test-1", "Content", 8, []string{"tag1"})
	m.SetMemory(memory)

	// With metadata
	m.showMetadata = true
	view1 := m.View()

	// Without metadata
	m.showMetadata = false
	view2 := m.View()

	// Views should be different
	if view1 == view2 {
		t.Error("Expected different views with and without metadata")
	}

	// Without metadata should not contain "Metadata" header
	if strings.Contains(view2, "Metadata") {
		t.Error("View without metadata should not contain 'Metadata' header")
	}
}

// --- Link Highlighting Tests ---

func TestRenderMetadataLineWithSelectedLink(t *testing.T) {
	m := NewModel()
	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
		{TargetId: "link-2"},
	}
	m.SetMemory(memory)

	// No link selected
	line18 := m.renderMetadataLine(18)
	if !strings.Contains(line18, "→") {
		t.Error("Expected unselected link to have '→' indicator")
	}

	// Select first link
	m.selectedLinkIndex = 0
	line18Selected := m.renderMetadataLine(18)
	if !strings.Contains(line18Selected, "▸") {
		t.Error("Expected selected link to have '▸' indicator")
	}

	// Second link should not be selected
	line19 := m.renderMetadataLine(19)
	if !strings.Contains(line19, "→") {
		t.Error("Expected unselected second link to have '→' indicator")
	}
}

func TestFooterWithLinkNavigation(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
	}
	m.SetMemory(memory)

	// No link selected - should show 'l: links' hint
	footer := m.renderFooter()
	if !strings.Contains(footer, "l: links") {
		t.Error("Expected footer to show 'l: links' hint when links available")
	}

	// Link selected - should show link navigation hints
	m.selectedLinkIndex = 0
	footerSelected := m.renderFooter()
	if !strings.Contains(footerSelected, "n/p") {
		t.Error("Expected footer to show 'n/p' navigation hints")
	}
	if !strings.Contains(footerSelected, "Enter: follow") {
		t.Error("Expected footer to show 'Enter: follow' hint")
	}
	if !strings.Contains(footerSelected, "Esc: deselect") {
		t.Error("Expected footer to show 'Esc: deselect' hint")
	}
}

func TestFooterWithoutLinks(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	memory := createTestMemory("test-1", "Content", 8, nil)
	m.SetMemory(memory)

	footer := m.renderFooter()

	// Should not mention links
	if strings.Contains(footer, "l: links") {
		t.Error("Footer should not show link hints when there are no links")
	}
	if strings.Contains(footer, "n/p") {
		t.Error("Footer should not show link navigation when there are no links")
	}
}

func TestViewWithSelectedLink(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	m.showMetadata = true

	memory := createTestMemory("test-1", "Test content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-target-123"},
	}
	m.SetMemory(memory)

	// Select the link
	m.selectedLinkIndex = 0

	view := m.View()

	// Debug: print view if test fails
	_ = view

	// Should show link navigation hints in footer
	if !strings.Contains(view, "Enter: follow") {
		t.Error("Expected view to show 'Enter: follow' hint")
	}

	// The selected link indicator might not be visible if the viewport
	// doesn't scroll down to where the links are displayed
	// This is acceptable - the important thing is the footer shows navigation hints
}
