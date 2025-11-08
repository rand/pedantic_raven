package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		PaddingTop(1).
		PaddingBottom(1)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		PaddingLeft(2)
)

type model struct {
	ready bool
}

func initialModel() model {
	return model{
		ready: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	title := titleStyle.Render("Pedantic Raven")
	subtitle := "Interactive Context Engineering Environment"
	help := helpStyle.Render("Press 'q' or 'ctrl+c' to quit")

	return fmt.Sprintf(
		"\n%s\n\n  %s\n\n%s\n\n",
		title,
		subtitle,
		help,
	)
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
