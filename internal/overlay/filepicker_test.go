package overlay

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/layout"
)

func TestNewFilePicker(t *testing.T) {
	fp := NewFilePicker("test-picker", "", nil)

	if fp == nil {
		t.Fatal("Expected file picker to be created")
	}

	if fp.ID() != "test-picker" {
		t.Errorf("Expected ID 'test-picker', got %v", fp.ID())
	}

	if !fp.Modal() {
		t.Error("Expected file picker to be modal")
	}

	if fp.currentDir == "" {
		t.Error("Expected current directory to be set")
	}
}

func TestFilePickerLoadDirectory(t *testing.T) {
	// Create a temporary directory with some files
	tmpDir, err := os.MkdirTemp("", "filepicker-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFile1 := filepath.Join(tmpDir, "file1.txt")
	testFile2 := filepath.Join(tmpDir, "file2.md")
	testDir := filepath.Join(tmpDir, "subdir")

	os.WriteFile(testFile1, []byte("test"), 0644)
	os.WriteFile(testFile2, []byte("test"), 0644)
	os.Mkdir(testDir, 0755)

	// Create file picker
	fp := NewFilePicker("test", tmpDir, nil)

	// Should have loaded files (+ parent directory entry)
	if len(fp.files) < 3 {
		t.Errorf("Expected at least 3 entries (.. + 2 files + 1 dir), got %d", len(fp.files))
	}

	// First entry should be ".."
	if fp.files[0].name != ".." {
		t.Errorf("Expected first entry to be '..', got %s", fp.files[0].name)
	}

	// Check that directory appears before files
	foundFile := false
	for _, entry := range fp.files[1:] { // Skip ".."
		if entry.isDir {
			if foundFile {
				t.Error("Expected directories to appear before files")
			}
		} else {
			foundFile = true
		}
	}
}

func TestFilePickerNavigation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filepicker-nav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("b"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "c.txt"), []byte("c"), 0644)

	fp := NewFilePicker("test", tmpDir, nil)

	// Initial selection should be 0
	if fp.selected != 0 {
		t.Errorf("Expected initial selection 0, got %d", fp.selected)
	}

	// Navigate down
	keyDown := tea.KeyMsg{Type: tea.KeyDown}
	fp.Update(keyDown)

	if fp.selected != 1 {
		t.Errorf("Expected selection 1 after down, got %d", fp.selected)
	}

	// Navigate up
	keyUp := tea.KeyMsg{Type: tea.KeyUp}
	fp.Update(keyUp)

	if fp.selected != 0 {
		t.Errorf("Expected selection 0 after up, got %d", fp.selected)
	}

	// Can't go negative
	fp.Update(keyUp)
	if fp.selected != 0 {
		t.Errorf("Expected selection to stay at 0, got %d", fp.selected)
	}
}

func TestFilePickerSearch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filepicker-search")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	os.WriteFile(filepath.Join(tmpDir, "apple.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpDir, "banana.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tmpDir, "apricot.md"), []byte(""), 0644)

	fp := NewFilePicker("test", tmpDir, nil)

	// Initially all files visible
	filtered := fp.filteredFiles()
	totalFiles := len(filtered)

	// Type 'a' to search
	keyA := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	fp.Update(keyA)

	filtered = fp.filteredFiles()
	if len(filtered) == totalFiles {
		t.Error("Expected filtering to reduce file count")
	}

	// All remaining files should contain 'a'
	for _, file := range filtered {
		if file.name == ".." {
			continue // Skip parent directory
		}
		if !containsIgnoreCase(file.name, "a") {
			t.Errorf("File %s should not be in filtered results", file.name)
		}
	}
}

func TestFilePickerSelectFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filepicker-select")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	selectedPath := ""
	onSelect := func(path string) tea.Cmd {
		selectedPath = path
		return nil
	}

	fp := NewFilePicker("test", tmpDir, onSelect)

	// Navigate to the file (skip parent directory)
	keyDown := tea.KeyMsg{Type: tea.KeyDown}
	fp.Update(keyDown)

	// Select the file
	keyEnter := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := fp.Update(keyEnter)

	// Execute commands to trigger callback
	if cmd != nil {
		msgs := executeBatchCmd(cmd)
		// Check for FilePickerResult
		foundResult := false
		for _, msg := range msgs {
			if result, ok := msg.(FilePickerResult); ok {
				if !result.Canceled && result.FilePath != "" {
					foundResult = true
				}
			}
		}
		if !foundResult {
			t.Error("Expected FilePickerResult message")
		}
	}

	// Callback should have been called
	if selectedPath == "" {
		t.Error("Expected onSelect callback to be called")
	}
}

func TestFilePickerCancel(t *testing.T) {
	fp := NewFilePicker("test", "", nil)

	// Press Esc to cancel
	keyEsc := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := fp.Update(keyEsc)

	if cmd == nil {
		t.Fatal("Expected command on cancel")
	}

	// Execute commands
	msgs := executeBatchCmd(cmd)

	// Should have FilePickerResult with Canceled=true
	foundCanceled := false
	for _, msg := range msgs {
		if result, ok := msg.(FilePickerResult); ok {
			if result.Canceled {
				foundCanceled = true
			}
		}
	}

	if !foundCanceled {
		t.Error("Expected canceled FilePickerResult")
	}
}

func TestFilePickerView(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "filepicker-view")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fp := NewFilePicker("test", tmpDir, nil)

	area := layout.Rect{X: 0, Y: 0, Width: 60, Height: 20}
	view := fp.View(area)

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain title
	if !containsIgnoreCase(view, "Open File") {
		t.Error("Expected view to contain title")
	}

	// Should contain directory path
	if !containsIgnoreCase(view, tmpDir) {
		t.Error("Expected view to contain directory path")
	}
}

// Helper functions

func containsIgnoreCase(s, substr string) bool {
	return containsString(toLower(s), toLower(substr))
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

func executeBatchCmd(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}

	msg := cmd()
	if batchMsg, ok := msg.(tea.BatchMsg); ok {
		msgs := []tea.Msg{}
		for _, c := range batchMsg {
			if c != nil {
				msgs = append(msgs, c())
			}
		}
		return msgs
	}

	return []tea.Msg{msg}
}
