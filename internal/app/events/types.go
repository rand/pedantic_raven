// Package events provides a pub/sub event system for decoupled component communication.
//
// This event system enables reactive updates across the application without tight coupling
// between components. Components publish events when their state changes, and other components
// subscribe to receive those events.
//
// Inspired by the crush TUI's pubsub.Broker pattern.
package events

import (
	"time"
)

// EventType identifies the category of event.
type EventType int

const (
	// Semantic analysis events
	SemanticAnalysisStarted EventType = iota
	SemanticAnalysisProgress
	SemanticAnalysisComplete
	SemanticAnalysisFailed

	// Memory events
	MemoryRecalled
	MemoryCreated
	MemoryUpdated
	MemoryDeleted
	MemorySearchStarted
	MemorySearchComplete

	// Agent orchestration events
	AgentStarted
	AgentProgress
	AgentCompleted
	AgentFailed
	AgentPaused
	AgentResumed

	// Collaboration events
	UserJoined
	UserLeft
	CursorMoved
	EditMade
	ChatMessageReceived

	// Proposal events
	ProposalGenerated
	ProposalAccepted
	ProposalRejected
	ProposalModified

	// Diagnostic events
	DiagnosticsUpdated
	ValidationStarted
	ValidationComplete

	// File system events
	FileChanged
	FileOpened
	FileClosed
	FileHistoryUpdated

	// UI events
	ModeChanged
	LayoutChanged
	PaneFocused
	OverlayOpened
	OverlayClosed

	// System events
	ServerConnected
	ServerDisconnected
	ServerReconnecting
	ErrorOccurred
)

// String returns the human-readable name of the event type.
func (e EventType) String() string {
	switch e {
	case SemanticAnalysisStarted:
		return "SemanticAnalysisStarted"
	case SemanticAnalysisProgress:
		return "SemanticAnalysisProgress"
	case SemanticAnalysisComplete:
		return "SemanticAnalysisComplete"
	case SemanticAnalysisFailed:
		return "SemanticAnalysisFailed"
	case MemoryRecalled:
		return "MemoryRecalled"
	case MemoryCreated:
		return "MemoryCreated"
	case MemoryUpdated:
		return "MemoryUpdated"
	case MemoryDeleted:
		return "MemoryDeleted"
	case MemorySearchStarted:
		return "MemorySearchStarted"
	case MemorySearchComplete:
		return "MemorySearchComplete"
	case AgentStarted:
		return "AgentStarted"
	case AgentProgress:
		return "AgentProgress"
	case AgentCompleted:
		return "AgentCompleted"
	case AgentFailed:
		return "AgentFailed"
	case AgentPaused:
		return "AgentPaused"
	case AgentResumed:
		return "AgentResumed"
	case UserJoined:
		return "UserJoined"
	case UserLeft:
		return "UserLeft"
	case CursorMoved:
		return "CursorMoved"
	case EditMade:
		return "EditMade"
	case ChatMessageReceived:
		return "ChatMessageReceived"
	case ProposalGenerated:
		return "ProposalGenerated"
	case ProposalAccepted:
		return "ProposalAccepted"
	case ProposalRejected:
		return "ProposalRejected"
	case ProposalModified:
		return "ProposalModified"
	case DiagnosticsUpdated:
		return "DiagnosticsUpdated"
	case ValidationStarted:
		return "ValidationStarted"
	case ValidationComplete:
		return "ValidationComplete"
	case FileChanged:
		return "FileChanged"
	case FileOpened:
		return "FileOpened"
	case FileClosed:
		return "FileClosed"
	case FileHistoryUpdated:
		return "FileHistoryUpdated"
	case ModeChanged:
		return "ModeChanged"
	case LayoutChanged:
		return "LayoutChanged"
	case PaneFocused:
		return "PaneFocused"
	case OverlayOpened:
		return "OverlayOpened"
	case OverlayClosed:
		return "OverlayClosed"
	case ServerConnected:
		return "ServerConnected"
	case ServerDisconnected:
		return "ServerDisconnected"
	case ServerReconnecting:
		return "ServerReconnecting"
	case ErrorOccurred:
		return "ErrorOccurred"
	default:
		return "Unknown"
	}
}

// Event represents a domain event that occurred in the application.
type Event struct {
	Type      EventType
	Timestamp time.Time
	Data      interface{}
}

// NewEvent creates a new event with the given type and data.
func NewEvent(eventType EventType, data interface{}) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

// --- Specific Event Data Types ---

// SemanticAnalysisData contains data for semantic analysis events.
type SemanticAnalysisData struct {
	BufferID string
	Progress float32 // 0.0 to 1.0
	Analysis *SemanticAnalysis
	Error    error
}

// SemanticAnalysis represents the result of semantic analysis.
type SemanticAnalysis struct {
	Triples  []Triple
	Holes    []TypedHole
	Entities []Entity
}

// Triple represents a subject-predicate-object relationship.
type Triple struct {
	Subject   string
	Predicate string
	Object    string
	Line      int
}

// TypedHole represents a ?placeholder that needs to be filled.
type TypedHole struct {
	Placeholder string
	Context     string
	Line        int
	Column      int
}

// Entity represents an extracted entity (capitalized word, @symbol, #file).
type Entity struct {
	Name      string
	Type      EntityType
	Frequency int
	Locations []Location
}

// EntityType categorizes entities.
type EntityType int

const (
	EntityPerson EntityType = iota
	EntityFile
	EntitySymbol
	EntityConcept
)

// Location represents a position in a document.
type Location struct {
	Line   int
	Column int
}

// MemoryData contains data for memory events.
type MemoryData struct {
	MemoryID  string
	MemoryIDs []string
	Query     string
	Results   []SearchResult
	Error     error
}

// SearchResult represents a memory search result with score.
type SearchResult struct {
	MemoryID       string
	Content        string
	Score          float32
	SemanticScore  *float32
	FTSScore       *float32
	GraphScore     *float32
}

// AgentData contains data for agent orchestration events.
type AgentData struct {
	AgentID   string
	AgentName string
	Task      string
	Status    string
	Progress  float32
	Result    interface{}
	Error     error
}

// UserData contains data for collaboration events.
type UserData struct {
	UserID   string
	Username string
	Position *Position
	Edit     *Edit
	Message  string
}

// Position represents a cursor position.
type Position struct {
	Line   int
	Column int
}

// Edit represents a document edit operation.
type Edit struct {
	Type   EditType
	Start  Position
	End    Position
	Text   string
}

// EditType categorizes edit operations.
type EditType int

const (
	EditInsert EditType = iota
	EditDelete
	EditReplace
)

// ProposalData contains data for proposal events.
type ProposalData struct{
	ProposalID  string
	Proposal    *ChangeProposal
	Accepted    bool
	Modifications string
}

// ChangeProposal represents an AI-generated suggestion.
type ChangeProposal struct {
	ID          string
	Title       string
	Description string
	Confidence  float32
	Changes     []Edit
	Rationale   string
}

// DiagnosticData contains data for diagnostic events.
type DiagnosticData struct {
	BufferID    string
	Diagnostics []Diagnostic
}

// Diagnostic represents a validation issue.
type Diagnostic struct {
	Severity Severity
	Message  string
	Location Location
	Code     string
}

// Severity categorizes diagnostic severity.
type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityHint
	SeverityInfo
)

// FileData contains data for file system events.
type FileData struct {
	Path  string
	Paths []string
}

// UIData contains data for UI events.
type UIData struct {
	Mode     string
	Layout   string
	PaneID   string
	Overlay  string
}

// ServerData contains data for server connection events.
type ServerData struct {
	ServerURL string
	Error     error
	Latency   time.Duration
}

// ErrorData contains data for error events.
type ErrorData struct {
	Error   error
	Context string
	Fatal   bool
}
