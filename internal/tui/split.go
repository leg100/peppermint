package tui

import tea "github.com/charmbracelet/bubbletea"

// split is a division of terminal space into two parts, which can be either be
// stacked horizontally or vertically.
type split struct {
	width  int
	height int
}

func (m split) Init() tea.Cmd {
	return nil
}

func (m split) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "<":
			// decrease width
			return m, tea.Quit
		case ">":
			// increase width
			return m, tea.Quit
		case "+":
			// increase height
			return m, tea.Quit
		case "-":
			// decrease height
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m split) View() string {
	return "moto\n"
}
