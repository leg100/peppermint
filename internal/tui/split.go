package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type secondaryStatus int

const (
	secondaryUnfocused secondaryStatus = iota
	secondaryHidden
	secondaryFocused
)

type childModel interface {
	tea.Model

	// Focus toggles the focus of the model.
	Focus(bool)
}

// split model composed of two child models, horizontally or vertically stacked.
type split struct {
	// The two models: primary goes on top if vertically stacked or to the left
	// if horizontally stacked.
	primary, secondary childModel
	// minimum sizes
	primaryMinSize, secondaryMinSize int
	// current sizes
	primaryCurrentSize, secondaryCurrentSize int
	// max size of primary model
	primaryMaxSize int
	// totalCurrentSize is the total size of the split dimension
	totalCurrentSize int
	// otherTotalCurrentSize is the total size of the non-split dimension
	otherTotalCurrentSize int
	// true if split horizontally, false if split vertically
	horizontal bool
	// key to shrink model primary and grow model secondary
	shrinkKey key.Binding
	// key to grow model primary and shrink model secondary
	growKey key.Binding
	// key to hide/show the secondary model
	toggleHideSecondaryKey key.Binding
	// secondaryStatus is the status of the secondary model.
	secondaryStatus secondaryStatus
	// primaryBorderStyle is the style of border for the primary model
	primaryBorderStyle lipgloss.Style
	// secondaryBorderStyle is the style of border for the secondary model
	secondaryBorderStyle lipgloss.Style
}

func horizontalSplit(primary, secondary tea.Model, primaryMinSize, secondaryMinSize int) split {
	return split{
		primary:                primary,
		secondary:              secondary,
		primaryMinSize:         primaryMinSize,
		secondaryMinSize:       secondaryMinSize,
		primaryMaxSize:         primaryMinSize,
		shrinkKey:              Keys.ShrinkHorizontalSplit,
		growKey:                Keys.GrowHorizontalSplit,
		horizontal:             true,
		toggleHideSecondaryKey: Keys.ExplorerFullScreen,
	}
}

func verticalSplit(primary, secondary tea.Model, primaryMinSize, secondaryMinSize int) split {
	return split{
		primary:                primary,
		secondary:              secondary,
		primaryMinSize:         primaryMinSize,
		secondaryMinSize:       secondaryMinSize,
		primaryMaxSize:         primaryMinSize,
		shrinkKey:              Keys.ShrinkVerticalSplit,
		growKey:                Keys.GrowVerticalSplit,
		toggleHideSecondaryKey: Keys.TogglePreview,
	}
}

func (m split) Init() tea.Cmd {
	return nil
}

func (m split) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.shrinkKey):
			return m.move(-1)
		case key.Matches(msg, m.growKey):
			return m.move(1)
		case key.Matches(msg, m.toggleHideSecondaryKey):
			switch m.secondaryStatus {
			case secondaryUnfocused, secondaryFocused:
				m.secondaryStatus = secondaryHidden
			case secondaryHidden:
				m.secondaryStatus = secondaryUnfocused
			}
			return m.setSizes()
		case key.Matches(msg, Keys.SwitchPane):
			switch m.secondaryStatus {
			case secondaryHidden:
			case secondaryUnfocused:
				m.secondaryStatus = secondaryFocused
				m.secondary.Focus(true)
			case secondaryFocused:
				// Secondary model should be a split because this parent model
				// capturing the switch pane key is the explorer, and its
				// secondary model is always a split.
				splitModel, ok := m.secondary.(split)
				if ok && !splitModel.cycleFocus() {
					// Secondary split model has cycled to its own secondary
					// model already, so this model can now focus its primary
					// model.
					m.secondaryStatus = secondaryUnfocused
				}
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		if m.horizontal {
			m.totalCurrentSize = msg.Width
			m.otherTotalCurrentSize = msg.Height
		} else {
			m.totalCurrentSize = msg.Height
			m.otherTotalCurrentSize = msg.Width
		}
		return m.setSizes()
	}
	return m, m.updateChildModels(msg, msg)
}

func (m split) ToggleFocus() bool {
	switch m.secondaryStatus {
	case secondaryFocused:
		return false
	case secondaryHidden:
		return m.primary.ToggleFocus()
	case secondaryUnfocused:
		m.secondaryStatus = secondaryFocused
		return m.secondary.ToggleFocus()
	}
	return false
}

// cycleFocus cycles the focus to the next model. If the primary model is
// currently focused, the focus is switched to the secondary model. If the
// secondary model is a split model too then cycleFocus is called on its model.
// cycleFocus returns false if its secondary model is already focused
func (m split) cycleFocus() bool {
	switch m.secondaryStatus {
	case secondaryHidden:
		return false
	case secondaryFocused:
		splitModel, ok := m.secondary.(split)
		if ok {
			return splitModel.cycleFocus(false)
		} else {
			if root {
				m.secondaryStatus = secondaryUnfocused
			} else {
				// Parent split cycles focus to its primary
				return false
			}
		}
	case secondaryUnfocused:
		splitModel, ok := m.secondary.(split)
		if ok {
			return splitModel.cycleFocus(false)
		}
	}
}

func (m *split) setSizes() (tea.Model, tea.Cmd) {
	if m.secondaryStatus {
		m.primaryCurrentSize = m.totalCurrentSize
		return m.updateChildSizes()
	}
	m.primaryCurrentSize = clamp(m.totalCurrentSize-m.secondaryCurrentSize, m.primaryMinSize, m.primaryMaxSize)
	m.secondaryCurrentSize = max(m.secondaryMinSize, m.totalCurrentSize-m.primaryCurrentSize)

	return m.updateChildSizes()
}

// move the split division by delta, shrinking and growing child models
// accordingly.
func (m *split) move(delta int) (tea.Model, tea.Cmd) {
	if m.secondaryStatus {
		// Secondary model is hidden so there is no split to move.
		return m, nil
	}
	m.primaryCurrentSize = clamp(m.primaryCurrentSize+delta, m.primaryMinSize, m.totalCurrentSize-m.secondaryMinSize)
	m.secondaryCurrentSize = clamp(m.secondaryCurrentSize-delta, m.secondaryMinSize, m.totalCurrentSize-m.primaryMinSize)
	m.primaryMaxSize = m.primaryCurrentSize

	return m.updateChildSizes()
}

func (m *split) updateChildSizes() (tea.Model, tea.Cmd) {
	newWindowSizeMsg := func(size int) (msg tea.WindowSizeMsg) {
		if m.horizontal {
			msg.Width = size
			msg.Height = m.otherTotalCurrentSize
		} else {
			msg.Width = m.otherTotalCurrentSize
			msg.Height = size
		}
		return
	}
	return m, m.updateChildModels(
		newWindowSizeMsg(m.primaryCurrentSize),
		newWindowSizeMsg(m.secondaryCurrentSize),
	)
}

func (m *split) updateChildModels(msg1, msg2 tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.primary, cmd = m.primary.Update(msg1)
	cmds = append(cmds, cmd)
	m.secondary, cmd = m.secondary.Update(msg2)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *split) setBorderStyles() {
	if m.secondaryFocused {
		if m.previewFocused {
			m.Table.SetBorderStyle(lipgloss.NormalBorder(), tui.InactivePreviewBorder)
			m.previewBorder = lipgloss.ThickBorder()
			m.previewBorderColor = tui.Blue
		} else {
			m.Table.SetBorderStyle(lipgloss.ThickBorder(), tui.Blue)
			m.previewBorder = lipgloss.NormalBorder()
			m.previewBorderColor = tui.InactivePreviewBorder
		}
	} else {
		m.Table.SetBorderStyle(lipgloss.NormalBorder(), lipgloss.NoColor{})
	}
}

func (m split) View() string {
	if m.horizontal {
		return lipgloss.JoinHorizontal(lipgloss.Top, m.primary.View(), m.secondary.View())
	} else {
		return lipgloss.JoinVertical(lipgloss.Top, m.primary.View(), m.secondary.View())
	}
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
