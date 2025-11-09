package memorydetail

import (
	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Model represents the memory detail view component state.
type Model struct {
	// Current memory being displayed
	memory *pb.MemoryNote

	// Viewport state
	scrollOffset int
	height       int
	width        int

	// Panel visibility
	showMetadata bool

	// UI state
	focused bool
	err     error

	// Client integration (optional)
	client interface{} // Can hold *mnemosyne.Client
}

// Messages for the memory detail component.
type (
	// MemoryLoadedMsg is sent when a memory is loaded.
	MemoryLoadedMsg struct {
		Memory *pb.MemoryNote
	}

	// MemoryErrorMsg is sent when memory loading fails.
	MemoryErrorMsg struct {
		Err error
	}

	// CloseRequestMsg is sent when the user wants to close the detail view.
	CloseRequestMsg struct{}

	// LinkSelectedMsg is sent when a user selects a link to navigate.
	LinkSelectedMsg struct {
		TargetID string
	}
)

// NewModel creates a new memory detail model.
func NewModel() Model {
	return Model{
		memory:       nil,
		scrollOffset: 0,
		height:       20,
		width:        80,
		showMetadata: true,
		focused:      true,
	}
}

// NewModelWithMemory creates a new model with a memory pre-loaded.
func NewModelWithMemory(memory *pb.MemoryNote) Model {
	m := NewModel()
	m.memory = memory
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetSize sets the component size.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetFocus sets the focus state.
func (m *Model) SetFocus(focused bool) {
	m.focused = focused
}

// IsFocused returns whether the component is focused.
func (m Model) IsFocused() bool {
	return m.focused
}

// SetMemory sets the memory to display.
func (m *Model) SetMemory(memory *pb.MemoryNote) {
	m.memory = memory
	m.scrollOffset = 0
	m.err = nil
}

// Memory returns the current memory.
func (m Model) Memory() *pb.MemoryNote {
	return m.memory
}

// SetError sets the error state.
func (m *Model) SetError(err error) {
	m.err = err
}

// Error returns the current error, if any.
func (m Model) Error() error {
	return m.err
}

// ToggleMetadata toggles the metadata panel visibility.
func (m *Model) ToggleMetadata() {
	m.showMetadata = !m.showMetadata
}

// ShowMetadata returns whether the metadata panel is visible.
func (m Model) ShowMetadata() bool {
	return m.showMetadata
}

// SetClient sets the mnemosyne client for this model.
func (m *Model) SetClient(client interface{}) {
	m.client = client
}

// Client returns the mnemosyne client, if set.
func (m Model) Client() interface{} {
	return m.client
}
