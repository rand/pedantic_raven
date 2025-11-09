package modes

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/memorydetail"
	"github.com/rand/pedantic-raven/internal/memorygraph"
	"github.com/rand/pedantic-raven/internal/memorylist"
	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// LayoutMode defines which layout is active.
type LayoutMode int

const (
	LayoutModeStandard LayoutMode = iota // List + Detail
	LayoutModeGraph                      // Full screen graph
)

// FocusTarget defines which component has focus in standard layout.
type FocusTarget int

const (
	FocusTargetList FocusTarget = iota
	FocusTargetDetail
)

// ExploreMode provides memory workspace with list, detail, and graph views.
type ExploreMode struct {
	*BaseMode

	// Components
	memoryList   *memorylist.Model
	memoryDetail *memorydetail.Model
	graph        *memorygraph.Model

	// Layout state
	layoutMode  LayoutMode
	focusTarget FocusTarget

	// UI state
	showHelp bool

	// Size tracking
	width  int
	height int

	// Mnemosyne client for real data
	mnemosyneClient *mnemosyne.Client
	useRealData     bool // If false, use sample data
}

// NewExploreMode creates a new explore mode with all components.
func NewExploreMode() *ExploreMode {
	return &ExploreMode{
		BaseMode: NewBaseMode(
			ModeExplore,
			"Explore",
			"Memory workspace with list, detail, and graph views",
		),
		memoryList:   nil, // Will be initialized in Init
		memoryDetail: nil, // Will be initialized in Init
		graph:        nil, // Will be initialized in Init
		layoutMode:   LayoutModeStandard,
		focusTarget:  FocusTargetList,
	}
}

// Init initializes the explore mode.
func (m *ExploreMode) Init() tea.Cmd {
	// Initialize memory list
	listModel := memorylist.NewModel()
	listModel.SetFocus(true) // List starts with focus
	m.memoryList = &listModel

	// Initialize memory detail
	detailModel := memorydetail.NewModel()
	detailModel.SetFocus(false) // Detail starts without focus
	m.memoryDetail = &detailModel

	// Initialize graph model
	graphModel := memorygraph.NewModel()
	m.graph = &graphModel

	// Initialize mnemosyne client
	cfg := mnemosyne.ConfigFromEnv()
	if cfg.Enabled {
		client, err := mnemosyne.NewClient(cfg)
		if err == nil {
			// Try to connect
			if err := client.Connect(); err == nil {
				m.mnemosyneClient = client
				m.useRealData = true

				// Set client on components
				m.memoryList.SetMnemosyneClient(client)
				m.memoryDetail.SetMnemosyneClient(client)
				m.graph.SetClient(client)
			}
		}
	}

	// Initialize base mode
	if m.BaseMode != nil {
		return m.BaseMode.Init()
	}
	return nil
}

// OnEnter is called when explore mode becomes active.
func (m *ExploreMode) OnEnter() tea.Cmd {
	var cmds []tea.Cmd

	if m.useRealData && m.mnemosyneClient != nil {
		// Load real data from mnemosyne
		if m.memoryList != nil {
			// Load memories with default filters
			filters := memorylist.Filters{
				Namespace:     "global", // Default to global namespace
				MinImportance: 0,
				MaxImportance: 0, // 0 means no max filter
			}
			cmds = append(cmds, memorylist.LoadMemoriesFromServer(m.mnemosyneClient, filters))
		}

		// Load graph - we'll use the first memory from the list as root
		// In a real implementation, you might want to track a specific root or let user choose
		if m.graph != nil {
			// For now, load a sample graph until we have a proper root selection mechanism
			// TODO: Implement root memory selection or use most important/recent memory
			cmds = append(cmds, m.loadSampleGraph())
		}
	} else {
		// Fallback to sample data
		if m.memoryList != nil {
			cmds = append(cmds, m.loadSampleMemories())
		}

		if m.graph != nil {
			cmds = append(cmds, m.loadSampleGraph())
		}
	}

	return tea.Batch(cmds...)
}

// OnExit is called when explore mode becomes inactive.
func (m *ExploreMode) OnExit() tea.Cmd {
	if m.BaseMode != nil {
		return m.BaseMode.OnExit()
	}
	return nil
}

// Update processes messages.
func (m *ExploreMode) Update(msg tea.Msg) (Mode, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle window size
	if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = wsMsg.Width
		m.height = wsMsg.Height
		cmds = append(cmds, m.handleWindowSize(wsMsg))
		return m, tea.Batch(cmds...)
	}

	// Handle global keybindings
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "?":
			// Toggle help overlay
			m.showHelp = !m.showHelp
			return m, nil

		case "esc":
			// Close help if open
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}

		case "g":
			// Toggle layout mode (not when help is showing)
			if !m.showHelp {
				m.toggleLayout()
				return m, nil
			}

		case "tab":
			// Cycle focus (only in standard layout, not when help is showing)
			if m.layoutMode == LayoutModeStandard && !m.showHelp {
				m.cycleFocus()
				return m, nil
			}
		}
	}

	// Handle component-specific messages
	switch msg := msg.(type) {
	case memorylist.MemorySelectedMsg:
		// User selected a memory in the list, show it in detail view
		if m.memoryDetail != nil && msg.Memory != nil {
			m.memoryDetail.SetMemory(msg.Memory)
		}
		return m, nil

	case memorydetail.LinkSelectedMsg:
		// User wants to navigate to a linked memory
		if m.mnemosyneClient != nil && msg.TargetID != "" {
			// Load the linked memory from server
			cmd := memorydetail.LoadMemory(m.mnemosyneClient, msg.TargetID)
			return m, cmd
		}
		return m, nil

	case memorydetail.CloseRequestMsg:
		// User wants to close detail view
		// Clear the detail view
		if m.memoryDetail != nil {
			m.memoryDetail.SetMemory(nil)
		}
		return m, nil

	case memorylist.MemoriesLoadedMsg:
		// Forward to memory list
		if m.memoryList != nil {
			updated, cmd := m.memoryList.Update(msg)
			*m.memoryList = updated
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)

	case memorygraph.GraphLoadedMsg:
		// Forward to graph
		if m.graph != nil {
			updated, cmd := m.graph.Update(msg)
			*m.graph = updated
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)
	}

	// Forward keyboard input to focused component (only if help is not showing)
	if !m.showHelp {
		if m.layoutMode == LayoutModeGraph {
			// In graph mode, all input goes to graph
			if m.graph != nil {
				updated, cmd := m.graph.Update(msg)
				*m.graph = updated
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		} else {
			// In standard mode, input goes to focused component
			switch m.focusTarget {
			case FocusTargetList:
				if m.memoryList != nil {
					updated, cmd := m.memoryList.Update(msg)
					*m.memoryList = updated
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
				}
			case FocusTargetDetail:
				if m.memoryDetail != nil {
					updated, cmd := m.memoryDetail.Update(msg)
					*m.memoryDetail = updated
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
				}
			}
		}
	}

	// Update base mode
	if m.BaseMode != nil {
		_, baseCmd := m.BaseMode.Update(msg)
		if baseCmd != nil {
			cmds = append(cmds, baseCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// handleWindowSize updates component sizes based on layout mode.
func (m *ExploreMode) handleWindowSize(msg tea.WindowSizeMsg) tea.Cmd {
	var cmds []tea.Cmd

	// Reserve space for UI chrome (title, status, help)
	contentHeight := msg.Height - 6
	if contentHeight < 10 {
		contentHeight = 10
	}

	if m.layoutMode == LayoutModeGraph {
		// Graph takes full width in graph mode
		if m.graph != nil {
			m.graph.SetSize(msg.Width, contentHeight)
		}
	} else {
		// Standard layout: list on left (40%), detail on right (60%)
		listWidth := msg.Width * 4 / 10
		detailWidth := msg.Width - listWidth - 1 // -1 for divider

		if m.memoryList != nil {
			m.memoryList.SetSize(listWidth, contentHeight)
		}
		if m.memoryDetail != nil {
			m.memoryDetail.SetSize(detailWidth, contentHeight)
		}
	}

	return tea.Batch(cmds...)
}

// toggleLayout switches between standard and graph layouts.
func (m *ExploreMode) toggleLayout() {
	if m.layoutMode == LayoutModeStandard {
		m.layoutMode = LayoutModeGraph
	} else {
		m.layoutMode = LayoutModeStandard
		// Return focus to list when switching back to standard
		m.focusTarget = FocusTargetList
		if m.memoryList != nil {
			m.memoryList.SetFocus(true)
		}
		if m.memoryDetail != nil {
			m.memoryDetail.SetFocus(false)
		}
	}

	// Trigger resize to update component dimensions
	if m.width > 0 && m.height > 0 {
		m.handleWindowSize(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	}
}

// cycleFocus cycles focus between list and detail in standard layout.
func (m *ExploreMode) cycleFocus() {
	if m.focusTarget == FocusTargetList {
		m.focusTarget = FocusTargetDetail
		if m.memoryList != nil {
			m.memoryList.SetFocus(false)
		}
		if m.memoryDetail != nil {
			m.memoryDetail.SetFocus(true)
		}
	} else {
		m.focusTarget = FocusTargetList
		if m.memoryList != nil {
			m.memoryList.SetFocus(true)
		}
		if m.memoryDetail != nil {
			m.memoryDetail.SetFocus(false)
		}
	}
}

// Styles for Explore mode
var (
	listBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")). // Blue
			Padding(0, 1)

	detailBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("35")). // Green
				Padding(0, 1)

	helpTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Padding(1, 0)

	helpTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
)

// View renders the explore mode.
func (m *ExploreMode) View() string {
	// Show help overlay if requested
	if m.showHelp {
		return m.renderHelp()
	}

	if m.layoutMode == LayoutModeGraph {
		// Full screen graph
		if m.graph == nil {
			return "Initializing graph visualization..."
		}
		return m.graph.View()
	}

	// Standard layout: list + detail side by side
	if m.memoryList == nil || m.memoryDetail == nil {
		return "Initializing memory workspace..."
	}

	// Get the views from components
	listView := m.memoryList.View()
	detailView := m.memoryDetail.View()

	// Calculate dimensions (40/60 split with borders)
	listWidth := (m.width * 4 / 10) - 4 // -4 for borders
	detailWidth := m.width - listWidth - 8 // -8 for both borders

	if listWidth < 10 {
		listWidth = 10
	}
	if detailWidth < 10 {
		detailWidth = 10
	}

	// Wrap in bordered containers
	listContainer := listBorderStyle.Width(listWidth).Render(listView)
	detailContainer := detailBorderStyle.Width(detailWidth).Render(detailView)

	// Join horizontally with a gap
	return lipgloss.JoinHorizontal(lipgloss.Top, listContainer, " ", detailContainer)
}

// renderHelp renders the help overlay with consistent styling per STYLE_GUIDE.md.
func (m *ExploreMode) renderHelp() string {
	var helpLines []string

	if m.layoutMode == LayoutModeStandard {
		helpLines = []string{
			"Explore Mode - Standard Layout",
			"",
			"NAVIGATION",
			"  ↑↓ or j/k    Select memory in list",
			"  Enter         View selected memory",
			"  Page Up/Dn    Scroll list",
			"",
			"DISPLAY",
			"  Tab           Toggle list/detail focus",
			"  g             Show graph view",
			"  m             Toggle metadata display",
			"",
			"MEMORY OPERATIONS",
			"  n             New memory",
			"  e             Edit memory",
			"  Delete        Delete memory",
			"  c             Create link",
			"  x             Delete link",
			"",
			"SEARCH & FILTER",
			"  /             Start search",
			"  Ctrl+F        Filter by importance",
			"  Ctrl+T        Filter by tags",
			"  Ctrl+X        Clear filters",
			"",
			"MODE SWITCHING",
			"  Tab / Shift+Tab  Next/previous mode",
			"  1-4              Jump to mode number",
			"",
			"OTHER",
			"  Cmd+P         Command palette",
			"  Esc           Close help",
			"  ?             Toggle this help",
		}
	} else {
		helpLines = []string{
			"Explore Mode - Graph Layout",
			"",
			"NAVIGATION",
			"  ↑↓←→        Pan graph",
			"  + / -       Zoom in/out",
			"  0           Reset view",
			"  Tab         Select next node",
			"",
			"NODE ACTIONS",
			"  Enter       Navigate to node",
			"  e           Expand node",
			"  x           Collapse node",
			"  c           Center on node",
			"  l           Show links",
			"",
			"GRAPH ACTIONS",
			"  r           Re-layout graph",
			"  Space       Single layout step",
			"  g           Back to list view",
			"",
			"MODE SWITCHING",
			"  Tab / Shift+Tab  Next/previous mode",
			"  1-4              Jump to mode number",
			"",
			"OTHER",
			"  Cmd+P       Command palette",
			"  Esc         Close help",
			"  ?           Toggle this help",
		}
	}

	// Calculate centering
	maxWidth := 0
	for _, line := range helpLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Build the help content
	var lines []string
	for _, line := range helpLines {
		if line == "" {
			lines = append(lines, "")
			continue
		}

		// Check if it's a title (first line or ends with "Layout")
		if strings.HasSuffix(line, "Layout") || strings.HasPrefix(line, "Explore Mode") {
			lines = append(lines, helpTitleStyle.Render(line))
		} else {
			lines = append(lines, helpTextStyle.Render(line))
		}
	}

	content := strings.Join(lines, "\n")

	// Create bordered box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 2).
		Width(maxWidth + 4)

	box := boxStyle.Render(content)

	// Center the box on screen
	if m.height > 0 {
		boxHeight := strings.Count(box, "\n") + 1
		verticalPadding := (m.height - boxHeight) / 2
		if verticalPadding > 0 {
			topPadding := strings.Repeat("\n", verticalPadding)
			box = topPadding + box
		}
	}

	if m.width > 0 {
		horizontalPadding := (m.width - (maxWidth + 8)) / 2
		if horizontalPadding > 0 {
			padding := strings.Repeat(" ", horizontalPadding)
			boxLines := strings.Split(box, "\n")
			for i, line := range boxLines {
				boxLines[i] = padding + line
			}
			box = strings.Join(boxLines, "\n")
		}
	}

	return box
}

// Keybindings returns the keybindings for explore mode.
func (m *ExploreMode) Keybindings() []Keybinding {
	if m.layoutMode == LayoutModeGraph {
		return []Keybinding{
			{Key: "g", Description: "Toggle to list view"},
			{Key: "h/j/k/l", Description: "Pan graph"},
			{Key: "+/-", Description: "Zoom in/out"},
			{Key: "0", Description: "Reset view"},
			{Key: "Tab", Description: "Select next node"},
			{Key: "Enter", Description: "Navigate to node"},
			{Key: "e", Description: "Expand node"},
			{Key: "x", Description: "Collapse node"},
			{Key: "c", Description: "Center on selected"},
			{Key: "r", Description: "Re-layout graph"},
			{Key: "Space", Description: "Layout step"},
		}
	}

	// Standard layout keybindings
	return []Keybinding{
		{Key: "g", Description: "Toggle to graph view"},
		{Key: "Tab", Description: "Switch focus (list/detail)"},
		{Key: "j/k", Description: "Navigate list"},
		{Key: "Enter", Description: "Select memory"},
		{Key: "/", Description: "Search"},
		{Key: "r", Description: "Refresh"},
		{Key: "?", Description: "Help"},
	}
}

// loadSampleGraph creates a sample graph for demonstration.
func (m *ExploreMode) loadSampleGraph() tea.Cmd {
	return func() tea.Msg {
		// Create a sample graph
		graph := memorygraph.NewGraph()

		// Add root node
		graph.AddNode(&memorygraph.Node{
			ID:         "root",
			IsExpanded: true,
		})

		// Add some child nodes
		graph.AddNode(&memorygraph.Node{
			ID:         "concept-a",
			IsExpanded: true,
		})
		graph.AddNode(&memorygraph.Node{
			ID:         "concept-b",
			IsExpanded: true,
		})
		graph.AddNode(&memorygraph.Node{
			ID:         "concept-c",
			IsExpanded: true,
		})

		// Add deeper nodes
		graph.AddNode(&memorygraph.Node{
			ID:         "detail-a1",
			IsExpanded: true,
		})
		graph.AddNode(&memorygraph.Node{
			ID:         "detail-a2",
			IsExpanded: true,
		})
		graph.AddNode(&memorygraph.Node{
			ID:         "detail-b1",
			IsExpanded: true,
		})

		// Add edges
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "root",
			TargetID: "concept-a",
			Strength: 1.0,
		})
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "root",
			TargetID: "concept-b",
			Strength: 1.0,
		})
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "root",
			TargetID: "concept-c",
			Strength: 1.0,
		})
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "concept-a",
			TargetID: "detail-a1",
			Strength: 1.0,
		})
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "concept-a",
			TargetID: "detail-a2",
			Strength: 1.0,
		})
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "concept-b",
			TargetID: "detail-b1",
			Strength: 1.0,
		})

		return memorygraph.GraphLoadedMsg{Graph: graph}
	}
}

// loadSampleMemories creates sample memories for the list.
func (m *ExploreMode) loadSampleMemories() tea.Cmd {
	return func() tea.Msg {
		// Create sample memories using the protobuf types
		memories := []*pb.MemoryNote{
			{
				Id:          "mem-1",
				Content:     "Architecture decision: Using event sourcing for audit trail",
				Importance:  8,
				Tags:        []string{"architecture", "patterns", "event-sourcing"},
				CreatedAt:   1704067200, // 2024-01-01
				UpdatedAt:   1704153600, // 2024-01-02
				Links:       []*pb.MemoryLink{},
				Namespace:   &pb.Namespace{Namespace: &pb.Namespace_Project{Project: &pb.ProjectNamespace{Name: "myapp"}}},
			},
			{
				Id:          "mem-2",
				Content:     "Performance optimization: Added caching layer to reduce database load",
				Importance:  7,
				Tags:        []string{"performance", "optimization", "caching"},
				CreatedAt:   1704240000, // 2024-01-03
				UpdatedAt:   1704326400, // 2024-01-04
				Links:       []*pb.MemoryLink{},
				Namespace:   &pb.Namespace{Namespace: &pb.Namespace_Project{Project: &pb.ProjectNamespace{Name: "myapp"}}},
			},
			{
				Id:          "mem-3",
				Content:     "Security review: JWT token validation needs improvement",
				Importance:  9,
				Tags:        []string{"security", "auth", "jwt"},
				CreatedAt:   1704412800, // 2024-01-05
				UpdatedAt:   1704499200, // 2024-01-06
				Links:       []*pb.MemoryLink{},
				Namespace:   &pb.Namespace{Namespace: &pb.Namespace_Global{Global: &pb.GlobalNamespace{}}},
			},
			{
				Id:          "mem-4",
				Content:     "Database schema: Created users table with proper indexes",
				Importance:  6,
				Tags:        []string{"database", "schema", "users"},
				CreatedAt:   1704585600, // 2024-01-07
				UpdatedAt:   1704672000, // 2024-01-08
				Links:       []*pb.MemoryLink{},
				Namespace:   &pb.Namespace{Namespace: &pb.Namespace_Project{Project: &pb.ProjectNamespace{Name: "myapp"}}},
			},
			{
				Id:          "mem-5",
				Content:     "API design: RESTful endpoints for user management",
				Importance:  5,
				Tags:        []string{"api", "rest", "design"},
				CreatedAt:   1704758400, // 2024-01-09
				UpdatedAt:   1704844800, // 2024-01-10
				Links:       []*pb.MemoryLink{},
				Namespace:   &pb.Namespace{Namespace: &pb.Namespace_Project{Project: &pb.ProjectNamespace{Name: "myapp"}}},
			},
		}

		return memorylist.MemoriesLoadedMsg{
			Memories:   memories,
			TotalCount: uint32(len(memories)),
		}
	}
}
