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
	// true if split horizontally, false if split vertically
	h bool
	// key to shrink model one and grow model two
	shrinkKey string
	// key to grow model one and shrink model two
	growKey string
	// d is the size of the primary dimension which is split
	d int
	// dOther is the size of the other dimension
	dOther int
}

func horizontalSplit(one, two tea.Model, minOne, minTwo int) split {
	return split{
		one:       one,
		two:       two,
		minOne:    minOne,
		minTwo:    minTwo,
		maxOne:    minOne,
		shrinkKey: "<",
		growKey:   ">",
		h:         true,
	}
}

func verticalSplit(one, two tea.Model, minOne, minTwo int) split {
	return split{
		one:       one,
		two:       two,
		minOne:    minOne,
		minTwo:    minTwo,
		maxOne:    minOne,
		shrinkKey: "-",
		growKey:   "+",
	}
}

func (m split) Init() tea.Cmd {
	return nil
}

func (m split) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case m.shrinkKey:
			return m.move(-1)
		case m.growKey:
			return m.move(1)
		}
	case tea.WindowSizeMsg:
		if m.h {
			m.d = msg.Width
			m.dOther = msg.Height
		} else {
			m.d = msg.Height
			m.dOther = msg.Width
		}
		m.currOne = clamp(m.d-m.currTwo, m.minOne, m.maxOne)
		m.currTwo = max(m.minTwo, m.d-m.currOne)
		return m.updateChildSizes()
	}
	return m, m.updateChildModels(msg, msg)
}

// move the split division by delta d, shrinking and growing child models
// accordingly.
func (m *split) move(d int) (tea.Model, tea.Cmd) {
	m.currOne = clamp(m.currOne+d, m.minOne, m.d-m.minTwo)
	m.currTwo = clamp(m.currTwo-d, m.minTwo, m.d-m.minOne)
	m.maxOne = m.currOne

	return m.updateChildSizes()
}

func (m *split) updateChildSizes() (tea.Model, tea.Cmd) {
	newWindowSizeMsg := func(size int) (msg tea.WindowSizeMsg) {
		if m.h {
			msg.Width = size
			msg.Height = m.dOther
		} else {
			msg.Width = m.dOther
			msg.Height = size
		}
		return
	}
	return m, m.updateChildModels(
		newWindowSizeMsg(m.currOne),
		newWindowSizeMsg(m.currTwo),
	)
}

func (m *split) updateChildModels(msg1, msg2 tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.one, cmd = m.one.Update(msg1)
	cmds = append(cmds, cmd)
	m.two, cmd = m.two.Update(msg2)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m split) View() string {
	if m.h {
		return lipgloss.JoinHorizontal(lipgloss.Top, m.one.View(), m.two.View())
	} else {
		return lipgloss.JoinVertical(lipgloss.Top, m.one.View(), m.two.View())
	}
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
