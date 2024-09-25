package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// horizontal is a stack of two models, one to the left, one to the right. Each
// model has a minimum width, a current width, but also a max width for
// the left model.
type horizontal struct {
	one, two tea.Model
	// minimum sizes
	minOne, minTwo int
	// current sizes
	currOne, currTwo int
	// max size of model one
	maxOne int
	// terminal dimensions
	width, height int
}

func newHorizontal(left, right tea.Model, minOne, minTwo int) tea.Model {
	return horizontal{
		one:    left,
		two:    right,
		minOne: minOne,
		minTwo: minTwo,
		maxOne: minOne,
	}
}

func (m horizontal) Init() tea.Cmd {
	return nil
}

func (m horizontal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "<":
			if currOne := m.currOne - 1; currOne < m.minOne {
				// Ignore
				return m, nil
			}
			m.currOne--
			m.maxOne--
			m.currTwo++
		case ">":
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
		m.one, _ = m.one.Update(tea.WindowSizeMsg{
			Width:  m.currOne,
			Height: m.height,
		})
		m.two, _ = m.two.Update(tea.WindowSizeMsg{
			Width:  m.currTwo,
			Height: m.height,
		})
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if msg.Width < (m.currOne + m.currTwo) {
			// User has shrunk terminal, so shrink right then left models.
			m.currTwo = max(m.minTwo, msg.Width-m.currOne)
			m.two, _ = m.two.Update(tea.WindowSizeMsg{
				Width:  m.currTwo,
				Height: msg.Height,
			})
			if (m.currOne + m.currTwo) == msg.Width {
				// Managed to get away with only shrinking right model
				return m, nil
			}
			// Nope, need to shrink left model too.
			m.currOne = max(m.minOne, msg.Width-m.currTwo)
			m.one, _ = m.one.Update(tea.WindowSizeMsg{
				Width:  m.currOne,
				Height: msg.Height,
			})
		} else if msg.Width > (m.currOne + m.currTwo) {
			// User has expanded terminal. If still smaller than sum of minimums
			// then do nothing
			if (m.minOne + m.minTwo) > msg.Width {
				return m, nil
			}
			// There is additional room to grow models. First grow the left
			// model as far as possible towards its max size.
			m.currOne = min(m.maxOne, msg.Width-m.currTwo)
			m.one, _ = m.one.Update(tea.WindowSizeMsg{
				Width:  m.currOne,
				Height: msg.Height,
			})
			// If there is any room spare, grow the right model
			m.currTwo = msg.Width - m.currOne
			m.two, _ = m.two.Update(tea.WindowSizeMsg{
				Width:  m.currTwo,
				Height: msg.Height,
			})
		}
	}
	return m, nil
}

func (m horizontal) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.one.View(), m.two.View())
}

type splitType interface {
	ShrinkKey() string
	GrowKey() string
	WindowSizeMsg(size int)
	View(tea.Model, tea.Model) string
}
