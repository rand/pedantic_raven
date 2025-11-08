// Package overlay provides modal dialogs and non-modal overlays for the TUI.
//
// The overlay system supports:
// - Modal dialogs that block interaction with the main UI
// - Non-modal popups (completions, quick references)
// - Stacking multiple overlays
// - Flexible positioning (center, cursor, custom)
// - Dismissal via Esc, click outside, or programmatic close
//
// Inspired by the crush TUI's dialog system.
package overlay

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/layout"
)

// OverlayID uniquely identifies an overlay.
type OverlayID string

// Overlay represents a UI element that appears above the main interface.
//
// Overlays can be modal (blocking interaction) or non-modal (allowing
// interaction with underlying UI). They handle their own rendering,
// input, and lifecycle.
type Overlay interface {
	// ID returns the unique identifier for this overlay.
	ID() OverlayID

	// Modal returns true if this overlay should block interaction with underlying UI.
	Modal() bool

	// Position returns the position strategy for this overlay.
	Position() Position

	// Init initializes the overlay and returns an initial command.
	Init() tea.Cmd

	// Update processes a Bubble Tea message and returns updated state.
	// Return DismissOverlay command to close the overlay.
	Update(msg tea.Msg) (Overlay, tea.Cmd)

	// View renders the overlay content.
	// The area parameter specifies the available space.
	View(area layout.Rect) string

	// OnDismiss is called when the overlay is being dismissed.
	// Use this to clean up resources or trigger actions on close.
	OnDismiss() tea.Cmd
}

// Position specifies how an overlay should be positioned.
type Position interface {
	// Compute calculates the overlay's rectangular area.
	// termWidth/termHeight are the terminal dimensions.
	// overlayWidth/overlayHeight are the overlay's desired size.
	Compute(termWidth, termHeight, overlayWidth, overlayHeight int) layout.Rect
}

// CenterPosition centers the overlay in the terminal.
type CenterPosition struct{}

func (p CenterPosition) Compute(termWidth, termHeight, overlayWidth, overlayHeight int) layout.Rect {
	x := (termWidth - overlayWidth) / 2
	y := (termHeight - overlayHeight) / 2

	// Ensure non-negative
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	return layout.Rect{
		X:      x,
		Y:      y,
		Width:  overlayWidth,
		Height: overlayHeight,
	}
}

// CursorPosition positions the overlay at a specific cursor location.
type CursorPosition struct {
	X int
	Y int
}

func (p CursorPosition) Compute(termWidth, termHeight, overlayWidth, overlayHeight int) layout.Rect {
	x := p.X
	y := p.Y

	// Ensure overlay stays within terminal bounds
	if x+overlayWidth > termWidth {
		x = termWidth - overlayWidth
	}
	if y+overlayHeight > termHeight {
		y = termHeight - overlayHeight
	}

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	return layout.Rect{
		X:      x,
		Y:      y,
		Width:  overlayWidth,
		Height: overlayHeight,
	}
}

// CustomPosition uses explicit coordinates.
type CustomPosition struct {
	X int
	Y int
}

func (p CustomPosition) Compute(termWidth, termHeight, overlayWidth, overlayHeight int) layout.Rect {
	return layout.Rect{
		X:      p.X,
		Y:      p.Y,
		Width:  overlayWidth,
		Height: overlayHeight,
	}
}

// --- Overlay Messages ---

// DismissOverlay is a command that closes the specified overlay.
type DismissOverlay struct {
	ID OverlayID
}

// OverlayDismissed is a message sent when an overlay is dismissed.
type OverlayDismissed struct {
	ID OverlayID
}

// --- Base Overlay Implementation ---

// BaseOverlay provides a default implementation of the Overlay interface.
//
// Concrete overlays can embed this to get default implementations
// and only override the methods they need to customize.
type BaseOverlay struct {
	id       OverlayID
	modal    bool
	position Position
	width    int
	height   int
	content  string
}

// NewBaseOverlay creates a new base overlay with the given parameters.
func NewBaseOverlay(id OverlayID, modal bool, position Position, width, height int) *BaseOverlay {
	return &BaseOverlay{
		id:       id,
		modal:    modal,
		position: position,
		width:    width,
		height:   height,
		content:  "",
	}
}

// ID implements Overlay.
func (o *BaseOverlay) ID() OverlayID {
	return o.id
}

// Modal implements Overlay.
func (o *BaseOverlay) Modal() bool {
	return o.modal
}

// Position implements Overlay.
func (o *BaseOverlay) Position() Position {
	return o.position
}

// Init implements Overlay.
func (o *BaseOverlay) Init() tea.Cmd {
	return nil
}

// Update implements Overlay.
func (o *BaseOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Esc dismisses by default
		if msg.String() == "esc" {
			return o, func() tea.Msg {
				return DismissOverlay{ID: o.id}
			}
		}
	}
	return o, nil
}

// View implements Overlay.
func (o *BaseOverlay) View(area layout.Rect) string {
	return o.content
}

// OnDismiss implements Overlay.
func (o *BaseOverlay) OnDismiss() tea.Cmd {
	return nil
}

// SetContent updates the overlay's content.
func (o *BaseOverlay) SetContent(content string) {
	o.content = content
}

// Width returns the overlay's desired width.
func (o *BaseOverlay) Width() int {
	return o.width
}

// Height returns the overlay's desired height.
func (o *BaseOverlay) Height() int {
	return o.height
}

// --- Common Overlay Types ---

// ConfirmDialog is a modal dialog asking for yes/no confirmation.
type ConfirmDialog struct {
	*BaseOverlay
	title    string
	message  string
	onYes    func() tea.Cmd
	onNo     func() tea.Cmd
	selected int // 0 = Yes, 1 = No
}

// NewConfirmDialog creates a new confirmation dialog.
func NewConfirmDialog(id OverlayID, title, message string, onYes, onNo func() tea.Cmd) *ConfirmDialog {
	return &ConfirmDialog{
		BaseOverlay: NewBaseOverlay(id, true, CenterPosition{}, 50, 10),
		title:       title,
		message:     message,
		onYes:       onYes,
		onNo:        onNo,
		selected:    0,
	}
}

// Update implements Overlay.
func (d *ConfirmDialog) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			d.selected = 0
			return d, nil
		case "right", "l":
			d.selected = 1
			return d, nil
		case "enter":
			var cmd tea.Cmd
			if d.selected == 0 && d.onYes != nil {
				cmd = d.onYes()
			} else if d.selected == 1 && d.onNo != nil {
				cmd = d.onNo()
			}
			// Dismiss the dialog
			dismissCmd := func() tea.Msg {
				return DismissOverlay{ID: d.id}
			}
			return d, tea.Batch(cmd, dismissCmd)
		case "esc":
			// Esc = No
			var cmd tea.Cmd
			if d.onNo != nil {
				cmd = d.onNo()
			}
			dismissCmd := func() tea.Msg {
				return DismissOverlay{ID: d.id}
			}
			return d, tea.Batch(cmd, dismissCmd)
		}
	}
	return d, nil
}

// View implements Overlay.
func (d *ConfirmDialog) View(area layout.Rect) string {
	// Simple text-based rendering for now
	// TODO: Use lipgloss for proper styling
	yesIndicator := " "
	noIndicator := " "
	if d.selected == 0 {
		yesIndicator = ">"
	} else {
		noIndicator = ">"
	}

	return d.title + "\n\n" +
		d.message + "\n\n" +
		yesIndicator + " [Yes]  " + noIndicator + " [No]"
}

// MessageDialog is a modal dialog displaying a message.
type MessageDialog struct {
	*BaseOverlay
	title   string
	message string
	onOK    func() tea.Cmd
}

// NewMessageDialog creates a new message dialog.
func NewMessageDialog(id OverlayID, title, message string, onOK func() tea.Cmd) *MessageDialog {
	return &MessageDialog{
		BaseOverlay: NewBaseOverlay(id, true, CenterPosition{}, 50, 10),
		title:       title,
		message:     message,
		onOK:        onOK,
	}
}

// Update implements Overlay.
func (d *MessageDialog) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc":
			var cmd tea.Cmd
			if d.onOK != nil {
				cmd = d.onOK()
			}
			dismissCmd := func() tea.Msg {
				return DismissOverlay{ID: d.id}
			}
			return d, tea.Batch(cmd, dismissCmd)
		}
	}
	return d, nil
}

// View implements Overlay.
func (d *MessageDialog) View(area layout.Rect) string {
	return d.title + "\n\n" +
		d.message + "\n\n" +
		"[OK]"
}
