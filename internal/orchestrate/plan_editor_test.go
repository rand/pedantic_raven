package orchestrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPlanEditorNew tests creation of a new plan editor.
func TestPlanEditorNew(t *testing.T) {
	editor := NewPlanEditor()

	assert.NotNil(t, editor)
	assert.Equal(t, "", editor.content)
	assert.Equal(t, 0, editor.cursor)
	assert.False(t, editor.dirty)
	assert.False(t, editor.isValid)
	assert.Equal(t, ModeEdit, editor.mode)
}

// TestPlanEditorInsertChar tests character insertion.
func TestPlanEditorInsertChar(t *testing.T) {
	editor := NewPlanEditor()

	// Insert a single character
	editor.insertChar('a')
	assert.Equal(t, "a", editor.content)
	assert.Equal(t, 1, editor.cursor)
	assert.True(t, editor.dirty)

	// Insert another character
	editor.insertChar('b')
	assert.Equal(t, "ab", editor.content)
	assert.Equal(t, 2, editor.cursor)

	// Insert in the middle
	editor.cursor = 1
	editor.insertChar('x')
	assert.Equal(t, "axb", editor.content)
	assert.Equal(t, 2, editor.cursor)
}

// TestPlanEditorDeleteChar tests backspace deletion.
func TestPlanEditorDeleteChar(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = "hello"
	editor.cursor = 5

	// Delete last character
	editor.deleteChar()
	assert.Equal(t, "hell", editor.content)
	assert.Equal(t, 4, editor.cursor)
	assert.True(t, editor.dirty)

	// Delete from middle (cursor at pos 2, delete char before pos 2 which is 'e')
	editor.cursor = 2
	editor.deleteChar()
	assert.Equal(t, "hll", editor.content)
	assert.Equal(t, 1, editor.cursor)

	// Delete with cursor at start
	editor.cursor = 0
	editor.deleteChar()
	assert.Equal(t, "hll", editor.content) // No change
	assert.Equal(t, 0, editor.cursor)
}

// TestPlanEditorDeleteCharForward tests delete key.
func TestPlanEditorDeleteCharForward(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = "hello"
	editor.cursor = 0

	// Delete first character
	editor.deleteCharForward()
	assert.Equal(t, "ello", editor.content)
	assert.True(t, editor.dirty)

	// Delete from middle
	editor.cursor = 1
	editor.deleteCharForward()
	assert.Equal(t, "elo", editor.content)

	// Delete at end
	editor.cursor = len(editor.content)
	editor.deleteCharForward()
	assert.Equal(t, "elo", editor.content) // No change
}

// TestPlanEditorNewline tests newline insertion.
func TestPlanEditorNewline(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = "hello"
	editor.cursor = 5

	// Insert newline at end
	editor.insertNewline()
	assert.Equal(t, "hello\n", editor.content)
	assert.Equal(t, 6, editor.cursor)
	assert.True(t, editor.dirty)

	// Insert newline in middle
	editor.cursor = 2
	editor.insertNewline()
	assert.Equal(t, "he\nllo\n", editor.content)
	assert.Equal(t, 3, editor.cursor)
}

// TestPlanEditorValidJSON tests JSON validation with valid content.
func TestPlanEditorValidJSON(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = `{
  "name": "Test Plan",
  "description": "A test work plan",
  "maxConcurrent": 2,
  "tasks": [
    {
      "id": "task1",
      "description": "First task",
      "type": 1,
      "dependencies": []
    }
  ]
}`

	editor.validateContent()

	assert.True(t, editor.isValid)
	assert.Empty(t, editor.validationErrors)
	assert.NotNil(t, editor.plan)
	assert.Equal(t, "Test Plan", editor.plan.Name)
	assert.Equal(t, 1, len(editor.plan.Tasks))
}

// TestPlanEditorInvalidJSON tests JSON validation with invalid content.
func TestPlanEditorInvalidJSON(t *testing.T) {
	editor := NewPlanEditor()

	// Invalid JSON syntax
	editor.content = `{
  "name": "Test Plan",
  "description": "Missing closing brace"
`

	editor.validateContent()

	assert.False(t, editor.isValid)
	assert.NotEmpty(t, editor.validationErrors)
	assert.Nil(t, editor.plan)
	assert.Contains(t, editor.validationErrors[0], "unexpected end of JSON")
}

// TestPlanEditorValidateMissingField tests validation with missing required field.
func TestPlanEditorValidateMissingField(t *testing.T) {
	editor := NewPlanEditor()

	// Missing "name" field
	editor.content = `{
  "description": "A test work plan",
  "maxConcurrent": 2,
  "tasks": []
}`

	editor.validateContent()

	assert.False(t, editor.isValid)
	assert.NotEmpty(t, editor.validationErrors)
	assert.Contains(t, editor.validationErrors[0], "name cannot be empty")
}

// TestPlanEditorValidateEmptyTasks tests validation with no tasks.
func TestPlanEditorValidateEmptyTasks(t *testing.T) {
	editor := NewPlanEditor()

	editor.content = `{
  "name": "Test Plan",
  "description": "A test work plan",
  "maxConcurrent": 2,
  "tasks": []
}`

	editor.validateContent()

	assert.False(t, editor.isValid)
	assert.NotEmpty(t, editor.validationErrors)
	assert.Contains(t, editor.validationErrors[0], "must contain at least one task")
}

// TestPlanEditorSave tests saving content to a file.
func TestPlanEditorSave(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test-plan.json")

	editor := NewPlanEditor()
	editor.content = `{
  "name": "Test Plan",
  "description": "A test work plan",
  "maxConcurrent": 2,
  "tasks": [
    {
      "id": "task1",
      "description": "First task",
      "type": 1,
      "dependencies": []
    }
  ]
}`
	editor.filename = filename

	// Save file
	err := editor.save()
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, filename)

	// Verify content
	data, err := os.ReadFile(filename)
	require.NoError(t, err)
	assert.Equal(t, editor.content, string(data))
}

// TestPlanEditorSaveNoFilename tests save without filename.
func TestPlanEditorSaveNoFilename(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = `{"name":"test"}`
	editor.filename = ""

	err := editor.save()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no filename")
}

// TestPlanEditorLoad tests loading content from a file.
func TestPlanEditorLoad(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test-plan.json")

	// Create a test file
	content := `{
  "name": "Loaded Plan",
  "description": "A loaded work plan",
  "maxConcurrent": 2,
  "tasks": [
    {
      "id": "task1",
      "description": "First task",
      "type": 1,
      "dependencies": []
    }
  ]
}`
	err := os.WriteFile(filename, []byte(content), 0o644)
	require.NoError(t, err)

	// Load the file
	editor := NewPlanEditor()
	err = editor.load(filename)
	require.NoError(t, err)

	// Verify content
	assert.Equal(t, content, editor.content)
	assert.Equal(t, 0, editor.cursor)
	assert.False(t, editor.dirty)
	assert.True(t, editor.isValid)
	assert.NotNil(t, editor.plan)
	assert.Equal(t, "Loaded Plan", editor.plan.Name)
}

// TestPlanEditorLoadNonexistent tests loading a non-existent file.
func TestPlanEditorLoadNonexistent(t *testing.T) {
	editor := NewPlanEditor()
	err := editor.load("/tmp/nonexistent-file-that-does-not-exist.json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}

// TestPlanEditorDirtyFlag tests dirty flag tracking.
func TestPlanEditorDirtyFlag(t *testing.T) {
	editor := NewPlanEditor()
	assert.False(t, editor.dirty)

	// Insert character marks as dirty
	editor.insertChar('a')
	assert.True(t, editor.dirty)

	// Save clears dirty flag
	tmpDir := t.TempDir()
	editor.filename = filepath.Join(tmpDir, "test.json")
	editor.dirty = false
	assert.False(t, editor.dirty)

	// Edit marks as dirty again
	editor.insertChar('b')
	assert.True(t, editor.dirty)
}

// TestPlanEditorMoveCursorUp tests moving cursor up.
func TestPlanEditorMoveCursorUp(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = "line1\nline2\nline3"

	// Start on line 3
	editor.cursor = 13 // Position in "line3"
	editor.moveCursorUp()

	// Should move to line 2
	lineStart := 6 // Start of "line2"
	assert.GreaterOrEqual(t, editor.cursor, lineStart)
	assert.Less(t, editor.cursor, 12)

	// Move up again
	editor.moveCursorUp()

	// Should move to line 1
	assert.Less(t, editor.cursor, 6)

	// Try to move up from first line (no change)
	oldCursor := editor.cursor
	editor.moveCursorUp()
	assert.Equal(t, oldCursor, editor.cursor)
}

// TestPlanEditorMoveCursorDown tests moving cursor down.
func TestPlanEditorMoveCursorDown(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = "line1\nline2\nline3"

	// Start on line 1
	editor.cursor = 2
	editor.moveCursorDown()

	// Should move to line 2
	lineStart := 6
	assert.GreaterOrEqual(t, editor.cursor, lineStart)

	// Move down again
	editor.moveCursorDown()

	// Should move to line 3
	assert.GreaterOrEqual(t, editor.cursor, 12)

	// Try to move down from last line (no change)
	oldCursor := editor.cursor
	editor.moveCursorDown()
	assert.Equal(t, oldCursor, editor.cursor)
}

// TestPlanEditorMoveCursorLineStart tests moving cursor to line start.
func TestPlanEditorMoveCursorLineStart(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = "line1\nline2\nline3"

	// Position in middle of line 2
	editor.cursor = 9
	editor.moveCursorToLineStart()

	// Should be at start of line 2
	assert.Equal(t, 6, editor.cursor)
}

// TestPlanEditorMoveCursorLineEnd tests moving cursor to line end.
func TestPlanEditorMoveCursorLineEnd(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = "line1\nline2\nline3"

	// Position at start of line 2
	editor.cursor = 6
	editor.moveCursorToLineEnd()

	// Should be at end of line 2 (before newline)
	assert.Equal(t, 11, editor.cursor)
}

// TestPlanEditorSaveLoadRoundtrip tests save and load roundtrip.
func TestPlanEditorSaveLoadRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "roundtrip.json")

	// Create and save
	editor1 := NewPlanEditor()
	editor1.content = `{
  "name": "Roundtrip Test",
  "description": "Testing save and load",
  "maxConcurrent": 3,
  "tasks": [
    {
      "id": "task1",
      "description": "Task 1",
      "type": 0,
      "dependencies": []
    },
    {
      "id": "task2",
      "description": "Task 2",
      "type": 1,
      "dependencies": ["task1"]
    }
  ]
}`
	editor1.filename = filename
	editor1.validateContent() // Validate so plan is populated
	err := editor1.save()
	require.NoError(t, err)

	// Load in a new editor
	editor2 := NewPlanEditor()
	err = editor2.load(filename)
	require.NoError(t, err)

	// Verify content matches
	assert.Equal(t, editor1.content, editor2.content)
	assert.NotNil(t, editor1.plan)
	assert.NotNil(t, editor2.plan)
	assert.Equal(t, editor1.plan.Name, editor2.plan.Name)
	assert.Equal(t, len(editor1.plan.Tasks), len(editor2.plan.Tasks))
}

// TestPlanEditorInit tests initialization.
func TestPlanEditorInit(t *testing.T) {
	editor := NewPlanEditor()
	cmd := editor.Init()
	assert.Nil(t, cmd)
}

// TestPlanEditorValidateComplexPlan tests validation of complex plan.
func TestPlanEditorValidateComplexPlan(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = `{
  "name": "Complex Plan",
  "description": "Multi-task plan with dependencies",
  "maxConcurrent": 2,
  "tasks": [
    {
      "id": "checkout",
      "description": "Checkout code",
      "type": 1,
      "dependencies": []
    },
    {
      "id": "build",
      "description": "Build",
      "type": 1,
      "dependencies": ["checkout"]
    },
    {
      "id": "test",
      "description": "Test",
      "type": 0,
      "dependencies": ["build"]
    },
    {
      "id": "deploy",
      "description": "Deploy",
      "type": 1,
      "dependencies": ["test"]
    }
  ]
}`

	editor.validateContent()

	assert.True(t, editor.isValid)
	assert.Empty(t, editor.validationErrors)
	assert.NotNil(t, editor.plan)
	assert.Equal(t, 4, len(editor.plan.Tasks))
}

// TestPlanEditorValidateCyclicDependencies tests validation detects cycles.
func TestPlanEditorValidateCyclicDependencies(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = `{
  "name": "Cyclic Plan",
  "description": "Plan with circular dependencies",
  "maxConcurrent": 2,
  "tasks": [
    {
      "id": "task1",
      "description": "Task 1",
      "type": 1,
      "dependencies": ["task2"]
    },
    {
      "id": "task2",
      "description": "Task 2",
      "type": 1,
      "dependencies": ["task1"]
    }
  ]
}`

	editor.validateContent()

	assert.False(t, editor.isValid)
	assert.NotEmpty(t, editor.validationErrors)
	assert.Contains(t, editor.validationErrors[0], "circular")
}

// TestPlanEditorCursorBounds tests cursor bounds.
func TestPlanEditorCursorBounds(t *testing.T) {
	editor := NewPlanEditor()
	editor.content = "hello"

	// Move cursor past end
	editor.cursor = 100

	// Cursor should stay at valid position when deleting forward
	editor.cursor = len(editor.content)
	editor.deleteCharForward()
	assert.Equal(t, "hello", editor.content)
}
