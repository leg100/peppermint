package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pane struct {
	width, height int
	content       string
}

func (m pane) Init() tea.Cmd {
	return nil
}

func (m pane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m pane) View() string {
	return lipgloss.NewStyle().Border(
		lipgloss.RoundedBorder(),
	).Render(
		lipgloss.Place(m.width-2, m.height-2, lipgloss.Center, lipgloss.Center, fmt.Sprintf("%dx%d", m.width, m.height)),
	)
}
