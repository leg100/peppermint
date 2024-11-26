package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type model struct {
	paner tea.Model
}

func New() model {
	return model{
		paner: newPaner(
			pane{n: 0},
			pane{n: 1},
			pane{n: 2},
		),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	m.paner, _ = m.paner.Update(msg)
	return m, nil
}

func (m model) View() string {
	return m.paner.View()
}
