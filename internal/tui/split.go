package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// split model composed of two child models, horizontally or verticalled
// stacked.
type split struct {
	one, two tea.Model
	// minimum sizes
	minOne, minTwo int
	// current sizes
	currOne, currTwo int
	// max size of model one
	maxOne int
	// stack algorithm: horizontal or vertical
	stack stack
}

func newSplit(one, two tea.Model, minOne, minTwo int, stack stack) tea.Model {
	return split{
		one:    one,
		two:    two,
		minOne: minOne,
		minTwo: minTwo,
		maxOne: minOne,
		stack:  stack,
	}
}

func (m split) Init() tea.Cmd {
	return nil
}

func (m split) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case m.stack.ShrinkKey():
			if currOne := m.currOne - 1; currOne < m.minOne {
				// Ignore
				return m, nil
			}
			m.currOne--
			m.maxOne--
			m.currTwo++
		case m.stack.GrowKey():
			if currTwo := m.currTwo - 1; currTwo < m.minTwo {
				// Ignore
				return m, nil
			}
			m.currOne++
			m.maxOne++
			m.currTwo--
		default:
			// Ignore
			return m, nil
		}
	case tea.WindowSizeMsg:
		size := m.stack.Update(msg)
		if size < (m.currOne + m.currTwo) {
			// User has shrunk terminal, so shrink child models to fit.
			m.currOne = max(m.minOne, size-m.currTwo)
			m.currTwo = max(m.minTwo, size-m.currOne)
		} else if size > (m.currOne + m.currTwo) {
			// User has expanded terminal.
			if (m.minOne + m.minTwo) > size {
				// Still smaller than sum of minimums, so ignore
				return m, nil
			}
			m.currOne = min(m.maxOne, size-m.currTwo)
			m.currTwo = size - m.currOne
		} else {
			// no change
			return m, nil
		}
	default:
		return m, nil
	}
	// Either terminal has been resized, or user has moved split, so update
	// child model dimensions accordingly.
	m.one = m.stack.UpdateChildModel(m.one, m.currOne)
	m.two = m.stack.UpdateChildModel(m.two, m.currTwo)
	return m, nil
}

func (m split) View() string {
	return m.stack.View(m.one, m.two)
}

type stack interface {
	ShrinkKey() string
	GrowKey() string
	// Update stack with new terminal size and return relevant terminal
	// dimension to resize split on.
	Update(tea.WindowSizeMsg) int
	UpdateChildModel(tea.Model, int) tea.Model
	View(tea.Model, tea.Model) string
}

type horizontal struct {
	height int
}

func (m *horizontal) ShrinkKey() string { return "<" }
func (m *horizontal) GrowKey() string   { return ">" }

func (m *horizontal) Update(msg tea.WindowSizeMsg) int {
	m.height = msg.Height
	return msg.Width
}

func (m *horizontal) UpdateChildModel(child tea.Model, width int) tea.Model {
	child, _ = child.Update(tea.WindowSizeMsg{
		Width:  width,
		Height: m.height,
	})
	return child
}

func (m *horizontal) View(one, two tea.Model) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, one.View(), two.View())
}

type vertical struct {
	width int
}

func (m *vertical) ShrinkKey() string { return "-" }
func (m *vertical) GrowKey() string   { return "+" }

func (m *vertical) Update(msg tea.WindowSizeMsg) int {
	m.width = msg.Width
	return msg.Height
}

func (m *vertical) UpdateChildModel(child tea.Model, height int) tea.Model {
	child, _ = child.Update(tea.WindowSizeMsg{
		Width:  m.width,
		Height: height,
	})
	return child
}

func (m *vertical) View(one, two tea.Model) string {
	return lipgloss.JoinVertical(lipgloss.Top, one.View(), two.View())
}
