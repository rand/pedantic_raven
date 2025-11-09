package memorydetail

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/mnemosyne"
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

	// Link selection state
	selectedLinkIndex int // -1 means no link selected

	// UI state
	focused bool
	err     error

	// Client integration
	mnemosyneClient *mnemosyne.Client

	// CRUD state
	editState         *EditState
	showDeleteConfirm bool
	isNewMemory       bool // True if creating a new memory
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
		memory:            nil,
		scrollOffset:      0,
		height:            20,
		width:             80,
		showMetadata:      true,
		selectedLinkIndex: -1,
		focused:           true,
		editState:         nil,
		showDeleteConfirm: false,
		isNewMemory:       false,
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
	m.selectedLinkIndex = -1
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

// SetMnemosyneClient sets the mnemosyne client for this model.
func (m *Model) SetMnemosyneClient(client *mnemosyne.Client) {
	m.mnemosyneClient = client
}

// MnemosyneClient returns the mnemosyne client, if set.
func (m Model) MnemosyneClient() *mnemosyne.Client {
	return m.mnemosyneClient
}

// CRUD operation methods

// IsEditing returns true if the model is in edit mode.
func (m Model) IsEditing() bool {
	return m.editState != nil && m.editState.isEditing
}

// HasUnsavedChanges returns true if there are unsaved changes.
func (m Model) HasUnsavedChanges() bool {
	if m.editState == nil {
		return false
	}
	return m.editState.detectChanges()
}

// EnterEditMode enters edit mode with the current memory.
func (m *Model) EnterEditMode() tea.Cmd {
	if m.memory == nil {
		return nil
	}

	return EnterEditMode(m.memory)
}

// EnterCreateMode enters edit mode for creating a new memory.
func (m *Model) EnterCreateMode(namespace *pb.Namespace) tea.Cmd {
	// Create a new empty memory
	newMemory := &pb.MemoryNote{
		Content:    "",
		Namespace:  namespace,
		Importance: 5, // Default importance
		Tags:       []string{},
	}

	m.isNewMemory = true
	return EnterEditMode(newMemory)
}

// SaveChanges saves the edited memory.
func (m *Model) SaveChanges() tea.Cmd {
	if m.editState == nil || m.editState.editedMemory == nil {
		return nil
	}

	return SaveChanges(m.mnemosyneClient, m.editState.editedMemory, m.isNewMemory)
}

// CancelEdit cancels editing and discards changes.
func (m *Model) CancelEdit() {
	m.editState = nil
	m.isNewMemory = false
}

// DeleteCurrentMemory requests deletion of the current memory.
func (m *Model) DeleteCurrentMemory() tea.Cmd {
	if m.memory == nil {
		return nil
	}

	return RequestDeleteConfirmation(m.memory)
}

// ConfirmDelete confirms and executes the deletion.
func (m *Model) ConfirmDelete() tea.Cmd {
	if m.memory == nil {
		return nil
	}

	m.showDeleteConfirm = false
	return DeleteMemory(m.mnemosyneClient, m.memory.Id)
}

// CancelDelete cancels the deletion.
func (m *Model) CancelDelete() {
	m.showDeleteConfirm = false
}

// ShowDeleteConfirmation returns true if showing delete confirmation.
func (m Model) ShowDeleteConfirmation() bool {
	return m.showDeleteConfirm
}

// EditedMemory returns the memory being edited, or nil if not editing.
func (m Model) EditedMemory() *pb.MemoryNote {
	if m.editState == nil {
		return nil
	}
	return m.editState.editedMemory
}

// SetEditedContent updates the content field in edit mode.
func (m *Model) SetEditedContent(content string) {
	if m.editState != nil && m.editState.editedMemory != nil {
		m.editState.editedMemory.Content = content
	}
}

// SetEditedTags updates the tags field in edit mode.
func (m *Model) SetEditedTags(tags []string) {
	if m.editState != nil && m.editState.editedMemory != nil {
		m.editState.editedMemory.Tags = tags
	}
}

// SetEditedImportance updates the importance field in edit mode.
func (m *Model) SetEditedImportance(importance uint32) {
	if m.editState != nil && m.editState.editedMemory != nil {
		m.editState.editedMemory.Importance = importance
	}
}

// SetEditedNamespace updates the namespace field in edit mode.
func (m *Model) SetEditedNamespace(namespace *pb.Namespace) {
	if m.editState != nil && m.editState.editedMemory != nil {
		m.editState.editedMemory.Namespace = namespace
	}
}

// GetFieldFocus returns the currently focused field in edit mode.
func (m Model) GetFieldFocus() EditField {
	if m.editState == nil {
		return FieldContent
	}
	return m.editState.fieldFocus
}

// SetFieldFocus sets the currently focused field in edit mode.
func (m *Model) SetFieldFocus(field EditField) {
	if m.editState != nil {
		m.editState.fieldFocus = field
	}
}

// CycleFieldFocus moves focus to the next field in edit mode.
func (m *Model) CycleFieldFocus() {
	if m.editState == nil {
		return
	}

	switch m.editState.fieldFocus {
	case FieldContent:
		m.editState.fieldFocus = FieldTags
	case FieldTags:
		m.editState.fieldFocus = FieldImportance
	case FieldImportance:
		m.editState.fieldFocus = FieldNamespace
	case FieldNamespace:
		m.editState.fieldFocus = FieldContent
	}
}

// Link navigation methods

// SelectNextLink selects the next link in the list.
func (m *Model) SelectNextLink() {
	if m.memory == nil || len(m.memory.Links) == 0 {
		return
	}

	m.selectedLinkIndex++
	if m.selectedLinkIndex >= len(m.memory.Links) {
		m.selectedLinkIndex = len(m.memory.Links) - 1
	}
}

// SelectPreviousLink selects the previous link in the list.
func (m *Model) SelectPreviousLink() {
	if m.memory == nil || len(m.memory.Links) == 0 {
		return
	}

	m.selectedLinkIndex--
	if m.selectedLinkIndex < -1 {
		m.selectedLinkIndex = -1
	}
}

// SelectFirstLink selects the first link.
func (m *Model) SelectFirstLink() {
	if m.memory == nil || len(m.memory.Links) == 0 {
		return
	}

	m.selectedLinkIndex = 0
}

// ClearLinkSelection clears the link selection.
func (m *Model) ClearLinkSelection() {
	m.selectedLinkIndex = -1
}

// SelectedLink returns the currently selected link, or nil if none.
func (m Model) SelectedLink() *pb.MemoryLink {
	if m.memory == nil || m.selectedLinkIndex < 0 || m.selectedLinkIndex >= len(m.memory.Links) {
		return nil
	}
	return m.memory.Links[m.selectedLinkIndex]
}

// SelectedLinkIndex returns the selected link index (-1 if none).
func (m Model) SelectedLinkIndex() int {
	return m.selectedLinkIndex
}

// HasLinks returns true if the memory has links.
func (m Model) HasLinks() bool {
	return m.memory != nil && len(m.memory.Links) > 0
}
