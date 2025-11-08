package overlay

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/layout"
)

var (
	filePickerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	filePickerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("170"))

	selectedFileStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Bold(true)

	fileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	dirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)
)

// FilePickerResult is sent when a file is selected.
type FilePickerResult struct {
	FilePath string
	Canceled bool
}

// FilePicker is an overlay for selecting files from the filesystem.
type FilePicker struct {
	*BaseOverlay
	currentDir  string
	searchQuery string
	files       []fileEntry
	selected    int
	offset      int
	onSelect    func(path string) tea.Cmd
}

type fileEntry struct {
	name  string
	path  string
	isDir bool
}

// NewFilePicker creates a new file picker overlay.
func NewFilePicker(id OverlayID, startDir string, onSelect func(path string) tea.Cmd) *FilePicker {
	if startDir == "" {
		startDir, _ = os.Getwd()
	}

	fp := &FilePicker{
		BaseOverlay: NewBaseOverlay(id, true, CenterPosition{}, 60, 20),
		currentDir:  startDir,
		searchQuery: "",
		selected:    0,
		offset:      0,
		onSelect:    onSelect,
	}

	// Load initial directory
	fp.loadDirectory()

	return fp
}

// loadDirectory reads the current directory and populates the file list.
func (fp *FilePicker) loadDirectory() {
	entries, err := os.ReadDir(fp.currentDir)
	if err != nil {
		fp.files = []fileEntry{}
		return
	}

	// Reset selection
	fp.selected = 0
	fp.offset = 0

	// Build file list with parent directory option
	fp.files = []fileEntry{
		{name: "..", path: filepath.Dir(fp.currentDir), isDir: true},
	}

	for _, entry := range entries {
		// Skip hidden files unless we want to show them
		if strings.HasPrefix(entry.Name(), ".") && entry.Name() != ".." {
			continue
		}

		fp.files = append(fp.files, fileEntry{
			name:  entry.Name(),
			path:  filepath.Join(fp.currentDir, entry.Name()),
			isDir: entry.IsDir(),
		})
	}

	// Sort: directories first, then files
	sort.Slice(fp.files, func(i, j int) bool {
		// Keep ".." at top
		if fp.files[i].name == ".." {
			return true
		}
		if fp.files[j].name == ".." {
			return false
		}

		// Directories before files
		if fp.files[i].isDir && !fp.files[j].isDir {
			return true
		}
		if !fp.files[i].isDir && fp.files[j].isDir {
			return false
		}

		// Alphabetical within same type
		return fp.files[i].name < fp.files[j].name
	})
}

// filteredFiles returns files matching the search query.
func (fp *FilePicker) filteredFiles() []fileEntry {
	if fp.searchQuery == "" {
		return fp.files
	}

	query := strings.ToLower(fp.searchQuery)
	filtered := []fileEntry{}

	for _, file := range fp.files {
		if strings.Contains(strings.ToLower(file.name), query) {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

// Update implements Overlay.
func (fp *FilePicker) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		filtered := fp.filteredFiles()

		switch msg.String() {
		case "esc":
			// Cancel and dismiss
			dismissCmd := func() tea.Msg {
				return DismissOverlay{ID: fp.id}
			}
			resultCmd := func() tea.Msg {
				return FilePickerResult{Canceled: true}
			}
			return fp, tea.Batch(resultCmd, dismissCmd)

		case "down", "j":
			if fp.selected < len(filtered)-1 {
				fp.selected++
				// Adjust offset for scrolling
				visibleLines := fp.Height() - 6 // Account for title, search, padding
				if fp.selected-fp.offset >= visibleLines {
					fp.offset++
				}
			}

		case "up", "k":
			if fp.selected > 0 {
				fp.selected--
				// Adjust offset for scrolling
				if fp.selected < fp.offset {
					fp.offset--
				}
			}

		case "enter":
			if len(filtered) > 0 {
				selected := filtered[fp.selected]

				if selected.isDir {
					// Navigate into directory
					fp.currentDir = selected.path
					fp.loadDirectory()
					fp.searchQuery = "" // Clear search when changing directory
				} else {
					// Select file
					var cmd tea.Cmd
					if fp.onSelect != nil {
						cmd = fp.onSelect(selected.path)
					}

					dismissCmd := func() tea.Msg {
						return DismissOverlay{ID: fp.id}
					}
					resultCmd := func() tea.Msg {
						return FilePickerResult{FilePath: selected.path, Canceled: false}
					}

					return fp, tea.Batch(cmd, resultCmd, dismissCmd)
				}
			}

		case "backspace":
			if len(fp.searchQuery) > 0 {
				fp.searchQuery = fp.searchQuery[:len(fp.searchQuery)-1]
				fp.selected = 0
				fp.offset = 0
			}

		default:
			// Add to search query if it's a printable character
			if msg.Type == tea.KeyRunes {
				fp.searchQuery += string(msg.Runes)
				fp.selected = 0
				fp.offset = 0
			}
		}
	}

	return fp, nil
}

// View implements Overlay.
func (fp *FilePicker) View(area layout.Rect) string {
	filtered := fp.filteredFiles()
	visibleLines := fp.Height() - 6 // Account for borders, title, search, padding

	// Build header
	title := filePickerTitleStyle.Render("Open File")
	dirPath := "Directory: " + fp.currentDir
	searchLine := "Search: " + fp.searchQuery + "█"

	if fp.searchQuery == "" {
		searchLine = "Search: (type to filter)"
	}

	var content strings.Builder
	content.WriteString(title + "\n")
	content.WriteString(dirPath + "\n")
	content.WriteString(searchLine + "\n")
	content.WriteString(strings.Repeat("─", area.Width-6) + "\n")

	// Build file list
	start := fp.offset
	end := start + visibleLines
	if end > len(filtered) {
		end = len(filtered)
	}

	for i := start; i < end; i++ {
		file := filtered[i]
		indicator := "  "
		if i == fp.selected {
			indicator = "> "
		}

		fileName := file.name
		if file.isDir {
			fileName = dirStyle.Render(fileName + "/")
		} else {
			fileName = fileStyle.Render(fileName)
		}

		if i == fp.selected {
			fileName = selectedFileStyle.Render(fileName)
		}

		content.WriteString(indicator + fileName + "\n")
	}

	// Show help at bottom
	help := "↑↓: navigate | Enter: select/open | Esc: cancel"
	content.WriteString(strings.Repeat("─", area.Width-6) + "\n")
	content.WriteString(help)

	return filePickerStyle.
		Width(area.Width - 4).
		Height(area.Height - 4).
		Render(content.String())
}
