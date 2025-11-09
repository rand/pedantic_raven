package memorydetail

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// EditField represents which field is currently being edited.
type EditField int

const (
	FieldContent EditField = iota
	FieldTags
	FieldImportance
	FieldNamespace
)

// EditState tracks edit mode state.
type EditState struct {
	isEditing    bool
	editedMemory *pb.MemoryNote
	fieldFocus   EditField
	originalHash string // For change detection
	hasChanges   bool
}

// Validation constraints
const (
	MaxContentLength   = 10000
	MaxTags            = 20
	MinImportance      = 1
	MaxImportance      = 10
	OperationTimeout   = 30 * time.Second
)

// Validation errors with clear, actionable messages
var (
	ErrContentRequired   = errors.New("✗ Validation: Content is required. Enter at least 1 character.")
	ErrContentTooLong    = errors.New("✗ Validation: Content exceeds 10000 characters. Truncate and try again.")
	ErrImportanceInvalid = errors.New("✗ Validation: Importance must be 1-10. Enter a number in this range.")
	ErrTooManyTags       = errors.New("✗ Validation: Too many tags (max 20). Remove some and retry.")
	ErrNamespaceRequired = errors.New("✗ Validation: Namespace is required. Use format: project:name.")
	ErrNoClient          = errors.New("✗ Config: Mnemosyne client not configured. Check config.toml and restart.")
)

// Messages for CRUD operations

// EditModeEnteredMsg is sent when edit mode is entered.
type EditModeEnteredMsg struct {
	Memory *pb.MemoryNote
}

// MemorySavedMsg is sent when a memory is successfully saved.
type MemorySavedMsg struct {
	Memory *pb.MemoryNote
	Err    error
}

// MemoryCreatedMsg is sent when a memory is successfully created.
type MemoryCreatedMsg struct {
	Memory *pb.MemoryNote
	Err    error
}

// MemoryUpdatedMsg is sent when a memory is successfully updated.
type MemoryUpdatedMsg struct {
	Memory *pb.MemoryNote
	Err    error
}

// MemoryDeletedMsg is sent when a memory is successfully deleted.
type MemoryDeletedMsg struct {
	MemoryID string
	Err      error
}

// DeleteConfirmationRequestMsg is sent when delete confirmation is needed.
type DeleteConfirmationRequestMsg struct {
	Memory *pb.MemoryNote
}

// EnterEditMode creates a command to enter edit mode.
func EnterEditMode(memory *pb.MemoryNote) tea.Cmd {
	return func() tea.Msg {
		if memory == nil {
			return EditModeEnteredMsg{Memory: nil}
		}

		// Create a deep copy of the memory for editing
		editedMemory := cloneMemory(memory)

		return EditModeEnteredMsg{
			Memory: editedMemory,
		}
	}
}

// MemoryClient defines the interface for memory CRUD operations.
type MemoryClient interface {
	StoreMemory(context.Context, mnemosyne.StoreMemoryOptions) (*pb.MemoryNote, error)
	UpdateMemory(context.Context, mnemosyne.UpdateMemoryOptions) (*pb.MemoryNote, error)
	DeleteMemory(context.Context, string) error
}

// SaveChanges saves the edited memory to the server.
func SaveChanges(client MemoryClient, memory *pb.MemoryNote, isNew bool) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return MemorySavedMsg{
				Memory: nil,
				Err:    ErrNoClient,
			}
		}

		// Validate the memory
		if err := validateMemory(memory); err != nil {
			return MemorySavedMsg{
				Memory: nil,
				Err:    err,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer cancel()

		var savedMemory *pb.MemoryNote
		var err error

		if isNew {
			// Create new memory
			savedMemory, err = createMemory(ctx, client, memory)
		} else {
			// Update existing memory
			savedMemory, err = updateMemory(ctx, client, memory)
		}

		return MemorySavedMsg{
			Memory: savedMemory,
			Err:    err,
		}
	}
}

// CreateMemory creates a new memory on the server.
func CreateMemory(client MemoryClient, memory *pb.MemoryNote) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return MemoryCreatedMsg{
				Memory: nil,
				Err:    ErrNoClient,
			}
		}

		// Validate the memory
		if err := validateMemory(memory); err != nil {
			return MemoryCreatedMsg{
				Memory: nil,
				Err:    err,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer cancel()

		savedMemory, err := createMemory(ctx, client, memory)

		return MemoryCreatedMsg{
			Memory: savedMemory,
			Err:    err,
		}
	}
}

// UpdateMemory updates an existing memory on the server.
func UpdateMemory(client MemoryClient, memory *pb.MemoryNote) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return MemoryUpdatedMsg{
				Memory: nil,
				Err:    ErrNoClient,
			}
		}

		// Validate the memory
		if err := validateMemory(memory); err != nil {
			return MemoryUpdatedMsg{
				Memory: nil,
				Err:    err,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer cancel()

		savedMemory, err := updateMemory(ctx, client, memory)

		return MemoryUpdatedMsg{
			Memory: savedMemory,
			Err:    err,
		}
	}
}

// DeleteMemory deletes a memory from the server.
func DeleteMemory(client MemoryClient, memoryID string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return MemoryDeletedMsg{
				MemoryID: memoryID,
				Err:      ErrNoClient,
			}
		}

		if memoryID == "" {
			return MemoryDeletedMsg{
				MemoryID: memoryID,
				Err:      mnemosyne.ErrInvalidArgument,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer cancel()

		err := client.DeleteMemory(ctx, memoryID)

		return MemoryDeletedMsg{
			MemoryID: memoryID,
			Err:      err,
		}
	}
}

// RequestDeleteConfirmation creates a message to request delete confirmation.
func RequestDeleteConfirmation(memory *pb.MemoryNote) tea.Cmd {
	return func() tea.Msg {
		return DeleteConfirmationRequestMsg{
			Memory: memory,
		}
	}
}

// Helper functions

// createMemory creates a new memory on the server.
func createMemory(ctx context.Context, client MemoryClient, memory *pb.MemoryNote) (*pb.MemoryNote, error) {
	importance := uint32(memory.Importance)
	memoryType := memory.MemoryType

	opts := mnemosyne.StoreMemoryOptions{
		Content:           memory.Content,
		Namespace:         memory.Namespace,
		Importance:        &importance,
		Context:           memory.Context,
		Tags:              memory.Tags,
		MemoryType:        &memoryType,
		SkipLLMEnrichment: false,
	}

	return client.StoreMemory(ctx, opts)
}

// updateMemory updates an existing memory on the server.
func updateMemory(ctx context.Context, client MemoryClient, memory *pb.MemoryNote) (*pb.MemoryNote, error) {
	content := memory.Content
	importance := uint32(memory.Importance)

	opts := mnemosyne.UpdateMemoryOptions{
		MemoryID:   memory.Id,
		Content:    &content,
		Importance: &importance,
		Tags:       memory.Tags,
	}

	return client.UpdateMemory(ctx, opts)
}

// validateMemory validates a memory before saving.
func validateMemory(memory *pb.MemoryNote) error {
	if memory == nil {
		return ErrContentRequired
	}

	// Validate content
	if strings.TrimSpace(memory.Content) == "" {
		return ErrContentRequired
	}

	if len(memory.Content) > MaxContentLength {
		return ErrContentTooLong
	}

	// Validate importance
	if memory.Importance < MinImportance || memory.Importance > MaxImportance {
		return ErrImportanceInvalid
	}

	// Validate tags
	if len(memory.Tags) > MaxTags {
		return ErrTooManyTags
	}

	// Validate namespace (only for new memories)
	if memory.Id == "" && memory.Namespace == nil {
		return ErrNamespaceRequired
	}

	return nil
}

// hashMemory creates a SHA256 hash of memory content for change detection.
func hashMemory(m *pb.MemoryNote) string {
	if m == nil {
		return ""
	}

	// Create a string representation of all editable fields
	data := fmt.Sprintf("%s|%d|%v|%s",
		m.Content,
		m.Importance,
		m.Tags,
		formatNamespace(m.Namespace),
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// cloneMemory creates a deep copy of a memory note.
func cloneMemory(m *pb.MemoryNote) *pb.MemoryNote {
	if m == nil {
		return nil
	}

	// Clone basic fields
	clone := &pb.MemoryNote{
		Id:          m.Id,
		Content:     m.Content,
		Summary:     m.Summary,
		Context:     m.Context,
		MemoryType:  m.MemoryType,
		Importance:  m.Importance,
		Confidence:  m.Confidence,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		AccessCount: m.AccessCount,
		IsArchived:  m.IsArchived,
	}

	// Clone tags
	if m.Tags != nil {
		clone.Tags = make([]string, len(m.Tags))
		copy(clone.Tags, m.Tags)
	}

	// Clone keywords
	if m.Keywords != nil {
		clone.Keywords = make([]string, len(m.Keywords))
		copy(clone.Keywords, m.Keywords)
	}

	// Clone related files
	if m.RelatedFiles != nil {
		clone.RelatedFiles = make([]string, len(m.RelatedFiles))
		copy(clone.RelatedFiles, m.RelatedFiles)
	}

	// Clone related entities
	if m.RelatedEntities != nil {
		clone.RelatedEntities = make([]string, len(m.RelatedEntities))
		copy(clone.RelatedEntities, m.RelatedEntities)
	}

	// Clone namespace
	clone.Namespace = cloneNamespace(m.Namespace)

	// Clone links
	if m.Links != nil {
		clone.Links = make([]*pb.MemoryLink, len(m.Links))
		for i, link := range m.Links {
			clone.Links[i] = cloneLink(link)
		}
	}

	// Clone optional fields
	if m.LastAccessedAt > 0 {
		clone.LastAccessedAt = m.LastAccessedAt
	}

	if m.ExpiresAt != nil {
		expiresAt := *m.ExpiresAt
		clone.ExpiresAt = &expiresAt
	}

	if m.SupersededBy != nil {
		supersededBy := *m.SupersededBy
		clone.SupersededBy = &supersededBy
	}

	// Clone embedding
	if m.Embedding != nil {
		clone.Embedding = make([]float32, len(m.Embedding))
		copy(clone.Embedding, m.Embedding)
	}

	clone.EmbeddingModel = m.EmbeddingModel

	return clone
}

// cloneNamespace creates a deep copy of a namespace.
func cloneNamespace(ns *pb.Namespace) *pb.Namespace {
	if ns == nil {
		return nil
	}

	clone := &pb.Namespace{}

	switch n := ns.Namespace.(type) {
	case *pb.Namespace_Global:
		clone.Namespace = &pb.Namespace_Global{
			Global: &pb.GlobalNamespace{},
		}

	case *pb.Namespace_Project:
		clone.Namespace = &pb.Namespace_Project{
			Project: &pb.ProjectNamespace{
				Name: n.Project.Name,
			},
		}

	case *pb.Namespace_Session:
		clone.Namespace = &pb.Namespace_Session{
			Session: &pb.SessionNamespace{
				Project:   n.Session.Project,
				SessionId: n.Session.SessionId,
			},
		}
	}

	return clone
}

// cloneLink creates a deep copy of a memory link.
func cloneLink(link *pb.MemoryLink) *pb.MemoryLink {
	if link == nil {
		return nil
	}

	clone := &pb.MemoryLink{
		TargetId:    link.TargetId,
		LinkType:    link.LinkType,
		Strength:    link.Strength,
		Reason:      link.Reason,
		CreatedAt:   link.CreatedAt,
		UserCreated: link.UserCreated,
	}

	if link.LastTraversedAt != nil {
		lastTraversed := *link.LastTraversedAt
		clone.LastTraversedAt = &lastTraversed
	}

	return clone
}

// detectChanges checks if the memory has been modified.
func (es *EditState) detectChanges() bool {
	if es.editedMemory == nil {
		return false
	}

	currentHash := hashMemory(es.editedMemory)
	return currentHash != es.originalHash
}

// parseNamespaceString parses a namespace string into a protobuf Namespace.
// Supports formats: "global" or "project:name"
func parseNamespaceString(ns string) *pb.Namespace {
	if ns == "global" {
		return mnemosyne.GlobalNamespace()
	}

	// Simple parsing - checks for project: prefix
	if len(ns) > 8 && ns[:8] == "project:" {
		projectName := ns[8:]
		return mnemosyne.ProjectNamespace(projectName)
	}

	// Default to treating as project namespace
	return mnemosyne.ProjectNamespace(ns)
}

// LoadMemory loads a memory by ID from the mnemosyne server.
// This is a Bubble Tea command that returns MemoryLoadedMsg with the result.
func LoadMemory(client *mnemosyne.Client, memoryID string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return MemoryLoadedMsg{
				Memory: nil,
			}
		}

		if !client.IsConnected() {
			return MemoryErrorMsg{
				Err: mnemosyne.ErrNotConnected,
			}
		}

		if memoryID == "" {
			return MemoryErrorMsg{
				Err: fmt.Errorf("%w: memory ID is required", mnemosyne.ErrInvalidArgument),
			}
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer cancel()

		// Load memory from server
		memory, err := client.GetMemory(ctx, memoryID)
		if err != nil {
			return MemoryErrorMsg{
				Err: err,
			}
		}

		return MemoryLoadedMsg{
			Memory: memory,
		}
	}
}
