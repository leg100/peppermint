package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type model struct {
	split tea.Model
}

func New() model {
	return model{
		split: newSplit(
			pane{},
			newSplit(
				pane{},
				pane{},
				10, 10,
				&vertical{},
			),
			10, 10,
			&horizontal{},
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
	m.split, _ = m.split.Update(msg)
	return m, nil
}

func (m model) View() string {
	return m.split.View()
}
