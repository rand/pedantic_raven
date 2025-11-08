package overlay

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/layout"
)

// --- Position Tests ---

func TestCenterPosition(t *testing.T) {
	pos := CenterPosition{}
	rect := pos.Compute(120, 30, 50, 10)

	expectedX := (120 - 50) / 2 // 35
	expectedY := (30 - 10) / 2  // 10

	if rect.X != expectedX {
		t.Errorf("Expected X %d, got %d", expectedX, rect.X)
	}

	if rect.Y != expectedY {
		t.Errorf("Expected Y %d, got %d", expectedY, rect.Y)
	}

	if rect.Width != 50 {
		t.Errorf("Expected width 50, got %d", rect.Width)
	}

	if rect.Height != 10 {
		t.Errorf("Expected height 10, got %d", rect.Height)
	}
}

func TestCenterPositionSmallTerminal(t *testing.T) {
	pos := CenterPosition{}
	rect := pos.Compute(40, 10, 50, 15)

	// Overlay larger than terminal, should clamp to 0
	if rect.X != 0 {
		t.Errorf("Expected X 0 for oversized overlay, got %d", rect.X)
	}

	if rect.Y != 0 {
		t.Errorf("Expected Y 0 for oversized overlay, got %d", rect.Y)
	}
}

func TestCursorPosition(t *testing.T) {
	pos := CursorPosition{X: 10, Y: 5}
	rect := pos.Compute(120, 30, 50, 10)

	if rect.X != 10 {
		t.Errorf("Expected X 10, got %d", rect.X)
	}

	if rect.Y != 5 {
		t.Errorf("Expected Y 5, got %d", rect.Y)
	}
}

func TestCursorPositionBoundsClipping(t *testing.T) {
	// Overlay would go off right edge
	pos := CursorPosition{X: 100, Y: 5}
	rect := pos.Compute(120, 30, 50, 10)

	expectedX := 120 - 50 // 70 (keep within bounds)
	if rect.X != expectedX {
		t.Errorf("Expected X %d (clipped to bounds), got %d", expectedX, rect.X)
	}

	// Overlay would go off bottom edge
	pos = CursorPosition{X: 10, Y: 25}
	rect = pos.Compute(120, 30, 50, 10)

	expectedY := 30 - 10 // 20 (keep within bounds)
	if rect.Y != expectedY {
		t.Errorf("Expected Y %d (clipped to bounds), got %d", expectedY, rect.Y)
	}
}

func TestCustomPosition(t *testing.T) {
	pos := CustomPosition{X: 20, Y: 15}
	rect := pos.Compute(120, 30, 50, 10)

	if rect.X != 20 {
		t.Errorf("Expected X 20, got %d", rect.X)
	}

	if rect.Y != 15 {
		t.Errorf("Expected Y 15, got %d", rect.Y)
	}

	if rect.Width != 50 {
		t.Errorf("Expected width 50, got %d", rect.Width)
	}

	if rect.Height != 10 {
		t.Errorf("Expected height 10, got %d", rect.Height)
	}
}

// --- Base Overlay Tests ---

func TestBaseOverlayAttributes(t *testing.T) {
	overlay := NewBaseOverlay("test-overlay", true, CenterPosition{}, 50, 10)

	if overlay.ID() != "test-overlay" {
		t.Errorf("Expected ID 'test-overlay', got '%s'", overlay.ID())
	}

	if !overlay.Modal() {
		t.Error("Expected overlay to be modal")
	}

	if overlay.Width() != 50 {
		t.Errorf("Expected width 50, got %d", overlay.Width())
	}

	if overlay.Height() != 10 {
		t.Errorf("Expected height 10, got %d", overlay.Height())
	}
}

func TestBaseOverlayEscDismiss(t *testing.T) {
	overlay := NewBaseOverlay("test-overlay", true, CenterPosition{}, 50, 10)

	// Send Esc key
	updatedOverlay, cmd := overlay.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = updatedOverlay

	// Should return DismissOverlay command
	if cmd == nil {
		t.Fatal("Expected DismissOverlay command")
	}

	msg := cmd()
	dismissMsg, ok := msg.(DismissOverlay)
	if !ok {
		t.Fatalf("Expected DismissOverlay message, got %T", msg)
	}

	if dismissMsg.ID != "test-overlay" {
		t.Errorf("Expected dismiss ID 'test-overlay', got '%s'", dismissMsg.ID)
	}
}

func TestBaseOverlaySetContent(t *testing.T) {
	overlay := NewBaseOverlay("test-overlay", true, CenterPosition{}, 50, 10)

	content := "Test content"
	overlay.SetContent(content)

	area := layout.Rect{X: 0, Y: 0, Width: 50, Height: 10}
	view := overlay.View(area)

	if view != content {
		t.Errorf("Expected content '%s', got '%s'", content, view)
	}
}

// --- Confirm Dialog Tests ---

func TestConfirmDialogYes(t *testing.T) {
	yesCalled := false
	onYes := func() tea.Cmd {
		yesCalled = true
		return nil
	}

	dialog := NewConfirmDialog("confirm", "Confirm", "Are you sure?", onYes, nil)

	// Default selected is Yes
	if dialog.selected != 0 {
		t.Error("Expected Yes to be selected by default")
	}

	// Press Enter
	_, cmd := dialog.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("Expected command after Enter")
	}

	// Execute the batch command
	// In a real scenario, the tea.Batch would execute both cmds
	// For testing, we verify onYes was called
	if !yesCalled {
		t.Error("Expected onYes callback to be called")
	}
}

func TestConfirmDialogNo(t *testing.T) {
	noCalled := false
	onNo := func() tea.Cmd {
		noCalled = true
		return nil
	}

	dialog := NewConfirmDialog("confirm", "Confirm", "Are you sure?", nil, onNo)

	// Move to No
	updated, _ := dialog.Update(tea.KeyMsg{Type: tea.KeyRight})
	dialog = updated.(*ConfirmDialog)

	if dialog.selected != 1 {
		t.Error("Expected No to be selected after right arrow")
	}

	// Press Enter
	_, cmd := dialog.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("Expected command after Enter")
	}

	if !noCalled {
		t.Error("Expected onNo callback to be called")
	}
}

func TestConfirmDialogEscEqualsNo(t *testing.T) {
	noCalled := false
	onNo := func() tea.Cmd {
		noCalled = true
		return nil
	}

	dialog := NewConfirmDialog("confirm", "Confirm", "Are you sure?", nil, onNo)

	// Press Esc (should trigger No)
	_, _ = dialog.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if !noCalled {
		t.Error("Expected onNo to be called when Esc is pressed")
	}
}

func TestConfirmDialogNavigation(t *testing.T) {
	dialog := NewConfirmDialog("confirm", "Confirm", "Are you sure?", nil, nil)

	// Start at Yes
	if dialog.selected != 0 {
		t.Error("Expected to start at Yes")
	}

	// Right arrow -> No
	updated, _ := dialog.Update(tea.KeyMsg{Type: tea.KeyRight})
	dialog = updated.(*ConfirmDialog)
	if dialog.selected != 1 {
		t.Error("Expected No after right arrow")
	}

	// Left arrow -> Yes
	updated, _ = dialog.Update(tea.KeyMsg{Type: tea.KeyLeft})
	dialog = updated.(*ConfirmDialog)
	if dialog.selected != 0 {
		t.Error("Expected Yes after left arrow")
	}

	// 'l' (vim-style right)
	updated, _ = dialog.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	dialog = updated.(*ConfirmDialog)
	if dialog.selected != 1 {
		t.Error("Expected No after 'l'")
	}

	// 'h' (vim-style left)
	updated, _ = dialog.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	dialog = updated.(*ConfirmDialog)
	if dialog.selected != 0 {
		t.Error("Expected Yes after 'h'")
	}
}

// --- Message Dialog Tests ---

func TestMessageDialogDismiss(t *testing.T) {
	okCalled := false
	onOK := func() tea.Cmd {
		okCalled = true
		return nil
	}

	dialog := NewMessageDialog("message", "Info", "This is a message", onOK)

	// Press Enter
	_, _ = dialog.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !okCalled {
		t.Error("Expected onOK to be called")
	}
}

func TestMessageDialogEsc(t *testing.T) {
	okCalled := false
	onOK := func() tea.Cmd {
		okCalled = true
		return nil
	}

	dialog := NewMessageDialog("message", "Info", "This is a message", onOK)

	// Press Esc (should also trigger OK)
	_, _ = dialog.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if !okCalled {
		t.Error("Expected onOK to be called on Esc")
	}
}

// --- Manager Tests ---

func TestManagerPushPop(t *testing.T) {
	manager := NewManager()

	if !manager.IsEmpty() {
		t.Error("Expected manager to be empty initially")
	}

	// Push overlay
	overlay := NewBaseOverlay("overlay1", true, CenterPosition{}, 50, 10)
	manager.Push(overlay)

	if manager.IsEmpty() {
		t.Error("Expected manager to have overlay after push")
	}

	if manager.Count() != 1 {
		t.Errorf("Expected count 1, got %d", manager.Count())
	}

	// Pop overlay
	popped, _ := manager.Pop()
	if popped == nil {
		t.Fatal("Expected to pop overlay")
	}

	if popped.ID() != "overlay1" {
		t.Errorf("Expected ID 'overlay1', got '%s'", popped.ID())
	}

	if !manager.IsEmpty() {
		t.Error("Expected manager to be empty after pop")
	}
}

func TestManagerTop(t *testing.T) {
	manager := NewManager()

	overlay1 := NewBaseOverlay("overlay1", true, CenterPosition{}, 50, 10)
	overlay2 := NewBaseOverlay("overlay2", true, CenterPosition{}, 50, 10)

	manager.Push(overlay1)
	manager.Push(overlay2)

	top := manager.Top()
	if top == nil {
		t.Fatal("Expected to get top overlay")
	}

	if top.ID() != "overlay2" {
		t.Errorf("Expected top ID 'overlay2', got '%s'", top.ID())
	}

	// Top should not remove the overlay
	if manager.Count() != 2 {
		t.Errorf("Expected count 2 after Top, got %d", manager.Count())
	}
}

func TestManagerGet(t *testing.T) {
	manager := NewManager()

	overlay1 := NewBaseOverlay("overlay1", true, CenterPosition{}, 50, 10)
	overlay2 := NewBaseOverlay("overlay2", true, CenterPosition{}, 50, 10)

	manager.Push(overlay1)
	manager.Push(overlay2)

	// Get overlay1
	found := manager.Get("overlay1")
	if found == nil {
		t.Fatal("Expected to find overlay1")
	}

	if found.ID() != "overlay1" {
		t.Errorf("Expected ID 'overlay1', got '%s'", found.ID())
	}

	// Get non-existent
	found = manager.Get("overlay3")
	if found != nil {
		t.Error("Expected nil for non-existent overlay")
	}
}

func TestManagerDismiss(t *testing.T) {
	manager := NewManager()

	overlay1 := NewBaseOverlay("overlay1", true, CenterPosition{}, 50, 10)
	overlay2 := NewBaseOverlay("overlay2", true, CenterPosition{}, 50, 10)

	manager.Push(overlay1)
	manager.Push(overlay2)

	// Dismiss overlay1 (not top)
	dismissed, _ := manager.Dismiss("overlay1")
	if dismissed == nil {
		t.Fatal("Expected to dismiss overlay1")
	}

	if dismissed.ID() != "overlay1" {
		t.Errorf("Expected dismissed ID 'overlay1', got '%s'", dismissed.ID())
	}

	// Should have only overlay2 left
	if manager.Count() != 1 {
		t.Errorf("Expected count 1 after dismiss, got %d", manager.Count())
	}

	if manager.Top().ID() != "overlay2" {
		t.Error("Expected overlay2 to be top after dismissing overlay1")
	}
}

func TestManagerClear(t *testing.T) {
	manager := NewManager()

	overlay1 := NewBaseOverlay("overlay1", true, CenterPosition{}, 50, 10)
	overlay2 := NewBaseOverlay("overlay2", true, CenterPosition{}, 50, 10)

	manager.Push(overlay1)
	manager.Push(overlay2)

	if manager.Count() != 2 {
		t.Fatalf("Expected count 2, got %d", manager.Count())
	}

	manager.Clear()

	if !manager.IsEmpty() {
		t.Error("Expected manager to be empty after clear")
	}

	if manager.Count() != 0 {
		t.Errorf("Expected count 0 after clear, got %d", manager.Count())
	}
}

func TestManagerHasModal(t *testing.T) {
	manager := NewManager()

	// No overlays
	if manager.HasModal() {
		t.Error("Expected HasModal to be false with no overlays")
	}

	// Add non-modal
	nonModal := NewBaseOverlay("nonmodal", false, CenterPosition{}, 50, 10)
	manager.Push(nonModal)

	if manager.HasModal() {
		t.Error("Expected HasModal to be false with only non-modal overlay")
	}

	// Add modal
	modal := NewBaseOverlay("modal", true, CenterPosition{}, 50, 10)
	manager.Push(modal)

	if !manager.HasModal() {
		t.Error("Expected HasModal to be true with modal overlay")
	}
}

func TestManagerTopModal(t *testing.T) {
	manager := NewManager()

	nonModal := NewBaseOverlay("nonmodal", false, CenterPosition{}, 50, 10)
	modal1 := NewBaseOverlay("modal1", true, CenterPosition{}, 50, 10)
	modal2 := NewBaseOverlay("modal2", true, CenterPosition{}, 50, 10)

	manager.Push(modal1)
	manager.Push(nonModal)
	manager.Push(modal2)

	topModal := manager.TopModal()
	if topModal == nil {
		t.Fatal("Expected to find a modal overlay")
	}

	if topModal.ID() != "modal2" {
		t.Errorf("Expected topmost modal to be 'modal2', got '%s'", topModal.ID())
	}
}

func TestManagerAllIDs(t *testing.T) {
	manager := NewManager()

	overlay1 := NewBaseOverlay("overlay1", true, CenterPosition{}, 50, 10)
	overlay2 := NewBaseOverlay("overlay2", true, CenterPosition{}, 50, 10)
	overlay3 := NewBaseOverlay("overlay3", true, CenterPosition{}, 50, 10)

	manager.Push(overlay1)
	manager.Push(overlay2)
	manager.Push(overlay3)

	ids := manager.AllIDs()
	if len(ids) != 3 {
		t.Fatalf("Expected 3 IDs, got %d", len(ids))
	}

	// Check order (bottom to top)
	if ids[0] != "overlay1" || ids[1] != "overlay2" || ids[2] != "overlay3" {
		t.Errorf("Expected IDs in order [overlay1, overlay2, overlay3], got %v", ids)
	}
}

func TestManagerUpdateDismissOverlay(t *testing.T) {
	manager := NewManager()

	overlay := NewBaseOverlay("overlay1", true, CenterPosition{}, 50, 10)
	manager.Push(overlay)

	if manager.Count() != 1 {
		t.Fatalf("Expected count 1, got %d", manager.Count())
	}

	// Send DismissOverlay message
	manager.Update(DismissOverlay{ID: "overlay1"})

	if !manager.IsEmpty() {
		t.Error("Expected overlay to be dismissed")
	}
}

func TestManagerUpdateWindowSize(t *testing.T) {
	manager := NewManager()

	manager.Update(tea.WindowSizeMsg{Width: 100, Height: 25})

	width, height := manager.TerminalSize()
	if width != 100 {
		t.Errorf("Expected width 100, got %d", width)
	}

	if height != 25 {
		t.Errorf("Expected height 25, got %d", height)
	}
}
