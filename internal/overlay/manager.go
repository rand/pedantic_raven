package overlay

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/layout"
)

// Manager handles the stack of overlays and their lifecycle.
//
// The manager maintains a stack where the topmost overlay receives input.
// Modal overlays block interaction with underlying overlays.
type Manager struct {
	stack      []Overlay
	termWidth  int
	termHeight int
}

// NewManager creates a new overlay manager.
func NewManager() *Manager {
	return &Manager{
		stack:      make([]Overlay, 0),
		termWidth:  120,
		termHeight: 30,
	}
}

// Push adds an overlay to the top of the stack.
//
// The overlay will be initialized and will receive input focus.
func (m *Manager) Push(overlay Overlay) tea.Cmd {
	m.stack = append(m.stack, overlay)
	return overlay.Init()
}

// Pop removes the topmost overlay from the stack.
//
// Returns the removed overlay and its OnDismiss command.
// Returns nil if the stack is empty.
func (m *Manager) Pop() (Overlay, tea.Cmd) {
	if len(m.stack) == 0 {
		return nil, nil
	}

	// Get the top overlay
	top := m.stack[len(m.stack)-1]

	// Remove from stack
	m.stack = m.stack[:len(m.stack)-1]

	// Call OnDismiss
	return top, top.OnDismiss()
}

// Dismiss removes a specific overlay from the stack.
//
// Returns the dismissed overlay and its OnDismiss command.
// Returns nil if the overlay is not found.
func (m *Manager) Dismiss(id OverlayID) (Overlay, tea.Cmd) {
	for i, overlay := range m.stack {
		if overlay.ID() == id {
			// Remove from stack
			m.stack = append(m.stack[:i], m.stack[i+1:]...)
			return overlay, overlay.OnDismiss()
		}
	}
	return nil, nil
}

// Top returns the topmost overlay without removing it.
// Returns nil if the stack is empty.
func (m *Manager) Top() Overlay {
	if len(m.stack) == 0 {
		return nil
	}
	return m.stack[len(m.stack)-1]
}

// Get retrieves an overlay by ID without removing it.
// Returns nil if not found.
func (m *Manager) Get(id OverlayID) Overlay {
	for _, overlay := range m.stack {
		if overlay.ID() == id {
			return overlay
		}
	}
	return nil
}

// IsEmpty returns true if there are no overlays.
func (m *Manager) IsEmpty() bool {
	return len(m.stack) == 0
}

// Count returns the number of overlays in the stack.
func (m *Manager) Count() int {
	return len(m.stack)
}

// Clear removes all overlays from the stack.
//
// Returns a batch command of all OnDismiss commands.
func (m *Manager) Clear() tea.Cmd {
	var cmds []tea.Cmd

	for i := len(m.stack) - 1; i >= 0; i-- {
		if cmd := m.stack[i].OnDismiss(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	m.stack = make([]Overlay, 0)

	return tea.Batch(cmds...)
}

// SetTerminalSize updates the manager's knowledge of terminal dimensions.
//
// This is used for overlay positioning calculations.
func (m *Manager) SetTerminalSize(width, height int) {
	m.termWidth = width
	m.termHeight = height
}

// TerminalSize returns the current terminal dimensions.
func (m *Manager) TerminalSize() (width, height int) {
	return m.termWidth, m.termHeight
}

// Update processes a Bubble Tea message and routes it to the appropriate overlay.
//
// The manager handles:
// - DismissOverlay: Removes the specified overlay
// - WindowSizeMsg: Updates terminal size
// - Other messages: Routed to the topmost modal overlay, or all non-modal overlays
func (m *Manager) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case DismissOverlay:
		_, cmd := m.Dismiss(msg.ID)
		return cmd

	case tea.WindowSizeMsg:
		m.SetTerminalSize(msg.Width, msg.Height)
		return nil

	default:
		if len(m.stack) == 0 {
			return nil
		}

		// Route to topmost overlay
		top := m.stack[len(m.stack)-1]
		updatedOverlay, cmd := top.Update(msg)
		m.stack[len(m.stack)-1] = updatedOverlay
		return cmd
	}
}

// View renders all visible overlays.
//
// Modal overlays obscure everything beneath them.
// Non-modal overlays are rendered in stack order.
//
// Returns a string representation of all overlays and their positions.
// In a real implementation, this would use advanced terminal rendering
// to properly layer overlays.
func (m *Manager) View() string {
	if len(m.stack) == 0 {
		return ""
	}

	// Find the bottom-most modal overlay
	bottomModal := -1
	for i := len(m.stack) - 1; i >= 0; i-- {
		if m.stack[i].Modal() {
			bottomModal = i
			break
		}
	}

	// Render from bottom-most modal (or all overlays if no modal)
	startIndex := 0
	if bottomModal >= 0 {
		startIndex = bottomModal
	}

	// Render each overlay
	result := ""
	for i := startIndex; i < len(m.stack); i++ {
		overlay := m.stack[i]
		area := m.computeOverlayArea(overlay)
		content := overlay.View(area)

		// In a real implementation, we'd properly position and layer the content
		// For now, just concatenate
		result += content + "\n"
	}

	return result
}

// computeOverlayArea calculates the rectangular area for an overlay.
func (m *Manager) computeOverlayArea(overlay Overlay) layout.Rect {
	// Get overlay's desired size
	// For now, use default sizes (will be improved when we add size methods)
	overlayWidth := 50
	overlayHeight := 10

	// If it's a BaseOverlay, get actual dimensions
	if base, ok := overlay.(*BaseOverlay); ok {
		overlayWidth = base.Width()
		overlayHeight = base.Height()
	}

	// Compute position
	return overlay.Position().Compute(m.termWidth, m.termHeight, overlayWidth, overlayHeight)
}

// HasModal returns true if there's at least one modal overlay in the stack.
func (m *Manager) HasModal() bool {
	for _, overlay := range m.stack {
		if overlay.Modal() {
			return true
		}
	}
	return false
}

// TopModal returns the topmost modal overlay.
// Returns nil if there are no modal overlays.
func (m *Manager) TopModal() Overlay {
	for i := len(m.stack) - 1; i >= 0; i-- {
		if m.stack[i].Modal() {
			return m.stack[i]
		}
	}
	return nil
}

// AllIDs returns the IDs of all overlays in the stack (bottom to top).
func (m *Manager) AllIDs() []OverlayID {
	ids := make([]OverlayID, len(m.stack))
	for i, overlay := range m.stack {
		ids[i] = overlay.ID()
	}
	return ids
}
