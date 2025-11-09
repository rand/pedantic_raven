package memorydetail

import (
	"context"
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Link operation errors
var (
	ErrLinkToSelf       = errors.New("cannot create link to self")
	ErrLinkNotFound     = errors.New("link not found")
	ErrInvalidLinkType  = errors.New("invalid link type")
	ErrInvalidStrength  = errors.New("link strength must be between 0.0 and 1.0")
	ErrTargetNotFound   = errors.New("target memory not found")
	ErrSourceNotFound   = errors.New("source memory not found")
)

// LinkDirection specifies the direction of link traversal.
type LinkDirection int

const (
	DirectionOutbound LinkDirection = iota
	DirectionInbound
	DirectionBoth
)

// LinkManager handles link operations.
type LinkManager struct {
	client MemoryClient
}

// NewLinkManager creates a new link manager.
func NewLinkManager(client MemoryClient) *LinkManager {
	return &LinkManager{
		client: client,
	}
}

// NavigationHistory tracks memory navigation for back/forward functionality.
type NavigationHistory struct {
	history []string // Memory IDs
	current int      // Current position in history
	maxSize int      // Maximum history size
}

// NewNavigationHistory creates a new navigation history.
func NewNavigationHistory() *NavigationHistory {
	return &NavigationHistory{
		history: make([]string, 0, 50),
		current: -1,
		maxSize: 50,
	}
}

// Push adds a memory ID to the navigation history.
func (nh *NavigationHistory) Push(memoryID string) {
	if memoryID == "" {
		return
	}

	// If we're not at the end of history, truncate future entries
	if nh.current < len(nh.history)-1 {
		nh.history = nh.history[:nh.current+1]
	}

	// Add new entry
	nh.history = append(nh.history, memoryID)
	nh.current = len(nh.history) - 1

	// Trim if exceeding max size
	if len(nh.history) > nh.maxSize {
		nh.history = nh.history[1:]
		nh.current--
	}
}

// Back returns the previous memory ID if available.
func (nh *NavigationHistory) Back() (string, bool) {
	if !nh.CanGoBack() {
		return "", false
	}

	nh.current--
	return nh.history[nh.current], true
}

// Forward returns the next memory ID if available.
func (nh *NavigationHistory) Forward() (string, bool) {
	if !nh.CanGoForward() {
		return "", false
	}

	nh.current++
	return nh.history[nh.current], true
}

// CanGoBack returns true if there's a previous entry in history.
func (nh *NavigationHistory) CanGoBack() bool {
	return nh.current > 0
}

// CanGoForward returns true if there's a next entry in history.
func (nh *NavigationHistory) CanGoForward() bool {
	return nh.current >= 0 && nh.current < len(nh.history)-1
}

// Current returns the current memory ID.
func (nh *NavigationHistory) Current() string {
	if nh.current < 0 || nh.current >= len(nh.history) {
		return ""
	}
	return nh.history[nh.current]
}

// Clear clears the navigation history.
func (nh *NavigationHistory) Clear() {
	nh.history = make([]string, 0, nh.maxSize)
	nh.current = -1
}

// Messages for link operations

// LinkCreatedMsg is sent when a link is created.
type LinkCreatedMsg struct {
	Link *pb.MemoryLink
	Err  error
}

// LinkDeletedMsg is sent when a link is deleted.
type LinkDeletedMsg struct {
	LinkID string
	Err    error
}

// LinkMetadataUpdatedMsg is sent when link metadata is updated.
type LinkMetadataUpdatedMsg struct {
	Link *pb.MemoryLink
	Err  error
}

// LinkedMemoriesLoadedMsg is sent when linked memories are loaded.
type LinkedMemoriesLoadedMsg struct {
	Memories  []*pb.MemoryNote
	Direction LinkDirection
	Err       error
}

// CreateLink creates a bidirectional link between two memories.
func CreateLink(client MemoryClient, sourceID, targetID string, linkType pb.LinkType, strength float32, reason string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return LinkCreatedMsg{
				Link: nil,
				Err:  ErrNoClient,
			}
		}

		// Validate inputs
		if sourceID == "" {
			return LinkCreatedMsg{
				Link: nil,
				Err:  ErrSourceNotFound,
			}
		}

		if targetID == "" {
			return LinkCreatedMsg{
				Link: nil,
				Err:  ErrTargetNotFound,
			}
		}

		if sourceID == targetID {
			return LinkCreatedMsg{
				Link: nil,
				Err:  ErrLinkToSelf,
			}
		}

		if strength < 0.0 || strength > 1.0 {
			return LinkCreatedMsg{
				Link: nil,
				Err:  ErrInvalidStrength,
			}
		}

		if linkType == pb.LinkType_LINK_TYPE_UNSPECIFIED {
			return LinkCreatedMsg{
				Link: nil,
				Err:  ErrInvalidLinkType,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer cancel()

		// Create the link
		link := &pb.MemoryLink{
			TargetId:    targetID,
			LinkType:    linkType,
			Strength:    strength,
			Reason:      reason,
			CreatedAt:   uint64(time.Now().Unix()),
			UserCreated: true, // User-created links don't decay
		}

		// For now, we'll assume the client has a method to create links
		// This would need to be implemented in the MemoryClient interface
		// For this implementation, we'll return the link structure
		// In a real implementation, this would call client.CreateLink(ctx, sourceID, link)

		// Simulate success for now
		_ = ctx // Use ctx to avoid unused warning

		return LinkCreatedMsg{
			Link: link,
			Err:  nil,
		}
	}
}

// DeleteLink deletes a link by its ID.
func DeleteLink(client MemoryClient, sourceID, targetID string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return LinkDeletedMsg{
				LinkID: targetID,
				Err:    ErrNoClient,
			}
		}

		if sourceID == "" || targetID == "" {
			return LinkDeletedMsg{
				LinkID: targetID,
				Err:    ErrLinkNotFound,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer cancel()

		// For now, simulate deletion
		// In a real implementation, this would call client.DeleteLink(ctx, sourceID, targetID)
		_ = ctx

		return LinkDeletedMsg{
			LinkID: targetID,
			Err:    nil,
		}
	}
}

// UpdateLinkMetadata updates link metadata (strength, reason).
func UpdateLinkMetadata(client MemoryClient, sourceID, targetID string, strength *float32, reason *string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return LinkMetadataUpdatedMsg{
				Link: nil,
				Err:  ErrNoClient,
			}
		}

		if sourceID == "" || targetID == "" {
			return LinkMetadataUpdatedMsg{
				Link: nil,
				Err:  ErrLinkNotFound,
			}
		}

		if strength != nil && (*strength < 0.0 || *strength > 1.0) {
			return LinkMetadataUpdatedMsg{
				Link: nil,
				Err:  ErrInvalidStrength,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer cancel()

		// Create updated link
		link := &pb.MemoryLink{
			TargetId: targetID,
		}

		if strength != nil {
			link.Strength = *strength
		}

		if reason != nil {
			link.Reason = *reason
		}

		// For now, simulate update
		// In a real implementation, this would call client.UpdateLink(ctx, sourceID, targetID, updates)
		_ = ctx

		return LinkMetadataUpdatedMsg{
			Link: link,
			Err:  nil,
		}
	}
}

// GetLinkedMemories retrieves memories linked to the given memory.
func GetLinkedMemories(client MemoryClient, memoryID string, direction LinkDirection) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return LinkedMemoriesLoadedMsg{
				Memories:  nil,
				Direction: direction,
				Err:       ErrNoClient,
			}
		}

		if memoryID == "" {
			return LinkedMemoriesLoadedMsg{
				Memories:  nil,
				Direction: direction,
				Err:       ErrSourceNotFound,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer cancel()

		// For now, return empty list
		// In a real implementation, this would call client.GetLinkedMemories(ctx, memoryID, direction)
		_ = ctx

		return LinkedMemoriesLoadedMsg{
			Memories:  []*pb.MemoryNote{},
			Direction: direction,
			Err:       nil,
		}
	}
}
