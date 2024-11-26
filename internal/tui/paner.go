package tui

import (
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type panePosition int

const (
	leftPane panePosition = iota
	topRightPane
	bottomRightPane
)

// paner manages the layout of the three panes that make up the Pug full screen terminal app.
type paner struct {
	active  panePosition
	models  map[panePosition]tea.Model
	widths  map[panePosition]int
	heights map[panePosition]int
	visible map[panePosition]bool
	// total width and height of the terminal space available to panes.
	width, height int
	// minimum width and heights for panes
	minWidth, minHeight int
	// maximum width of left pane and maximum height of top right pane
	maxLeftPaneWidth, maxTopRightHeight int
}

func newPaner(
	leftModel tea.Model,
	topRightModel tea.Model,
	bottomRightModel tea.Model,
) paner {
	p := paner{
		active: leftPane,
		models: map[panePosition]tea.Model{
			leftPane:        leftModel,
			topRightPane:    topRightModel,
			bottomRightPane: bottomRightModel,
		},
		widths: map[panePosition]int{
			leftPane:        0,
			topRightPane:    0,
			bottomRightPane: 0,
		},
		heights: map[panePosition]int{
			leftPane:        0,
			topRightPane:    0,
			bottomRightPane: 0,
		},
		visible: map[panePosition]bool{
			leftPane:        true,
			topRightPane:    true,
			bottomRightPane: true,
		},
		minWidth:          10,
		minHeight:         10,
		maxLeftPaneWidth:  10,
		maxTopRightHeight: 10,
	}
	return p
}

func (p paner) Init() tea.Cmd {
	return nil
}

func (p paner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.ShrinkPaneWidth):
			p.changeActivePaneWidth(-1)
		case key.Matches(msg, Keys.GrowPaneWidth):
			p.changeActivePaneWidth(1)
		case key.Matches(msg, Keys.ShrinkPaneHeight):
			p.changeActivePaneHeight(-1)
		case key.Matches(msg, Keys.GrowPaneHeight):
			p.changeActivePaneHeight(1)
		case key.Matches(msg, Keys.SwitchPane):
			p.switchActivePane()
		case key.Matches(msg, Keys.ClosePane):
			p.closeActivePane()
		}
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		p.setPaneWidths()
		p.setPaneHeights()
		p.updateChildSizes()
	}
	return p, nil
}

func (p *paner) switchActivePane() {
	var visiblePanes []panePosition
	for position, visible := range p.visible {
		if visible {
			visiblePanes = append(visiblePanes, position)
		}
	}
	slices.Sort(visiblePanes)
	newActivePanePosition := (p.active + 1) % panePosition(len(visiblePanes))
	p.active = visiblePanes[newActivePanePosition]
}

func (p *paner) closeActivePane() {
	var numVisible int
	for _, visible := range p.visible {
		if visible {
			numVisible++
		}
	}
	if numVisible == 1 {
		// cannot close last visible pane
		return
	}
	p.visible[p.active] = false
	p.switchActivePane()

	p.setPaneHeights()
	p.setPaneWidths()
	p.updateChildSizes()
}

func (p *paner) setPaneWidths() {
	if p.visible[topRightPane] || p.visible[bottomRightPane] {
		p.widths[leftPane] = clamp(p.widths[leftPane], p.minWidth, p.width-p.minWidth)
	} else {
		p.widths[leftPane] = p.width
	}
	if p.visible[leftPane] {
		p.widths[topRightPane] = max(p.minWidth, p.width-p.widths[leftPane])
		p.widths[bottomRightPane] = max(p.minWidth, p.width-p.widths[leftPane])
	} else {
		p.widths[topRightPane] = p.width
		p.widths[bottomRightPane] = p.width
	}
}

func (p *paner) setPaneHeights() {
	p.heights[leftPane] = p.height
	if p.visible[bottomRightPane] {
		p.heights[topRightPane] = clamp(p.heights[topRightPane], p.minHeight, p.height-p.minHeight)
	} else {
		p.heights[topRightPane] = p.height
	}
	if p.visible[topRightPane] {
		p.heights[bottomRightPane] = max(p.minHeight, p.height-p.heights[topRightPane])
	} else {
		p.heights[bottomRightPane] = p.height
	}
}

func (p *paner) changeActivePaneWidth(delta int) {
	switch p.active {
	case topRightPane, bottomRightPane:
		// on the right panes, shrink width is actually grow width, and vice
		// versa
		delta = -delta
	}
	for position := range p.models {
		if position == p.active {
			p.widths[position] = clamp(p.widths[position]+delta, p.minWidth, p.width-p.minWidth)
		} else {
			p.widths[position] = clamp(p.widths[position]-delta, p.minWidth, p.width-p.minWidth)
		}
	}
	p.maxLeftPaneWidth = p.widths[leftPane]
	p.setPaneWidths()
	p.updateChildSizes()
}

func (p *paner) changeActivePaneHeight(delta int) {
	if p.active == leftPane {
		// Cannot change height of left pane because it occupies the full height
		// already.
		return
	}
	for position := range p.models {
		if position == p.active {
			p.heights[position] = clamp(p.heights[position]+delta, p.minHeight, p.height-p.minHeight)
		} else {
			p.heights[position] = clamp(p.heights[position]-delta, p.minHeight, p.height-p.minHeight)
		}
	}
	p.maxTopRightHeight = p.heights[topRightPane]
	p.setPaneHeights()
	p.updateChildSizes()
}

func (p *paner) updateChildSizes() {
	for position, model := range p.models {
		p.models[position], _ = model.Update(tea.WindowSizeMsg{
			Width:  p.widths[position] - 2,
			Height: p.heights[position] - 2,
		})
	}
}

var borderStyle = map[bool]lipgloss.Style{
	true:  lipgloss.NewStyle().Border(lipgloss.ThickBorder()),
	false: lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
}

func (m paner) View() string {
	renderPane := func(position panePosition) string {
		if !m.visible[position] {
			return ""
		}
		rendered := m.models[position].View()
		isActive := position == m.active
		return borderStyle[isActive].Render(rendered)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top,
		removeEmptyStrings(
			renderPane(leftPane),
			lipgloss.JoinVertical(lipgloss.Top,
				removeEmptyStrings(
					renderPane(topRightPane),
					renderPane(bottomRightPane),
				)...,
			),
		)...,
	)
}

func removeEmptyStrings(strs ...string) []string {
	n := 0
	for _, s := range strs {
		if s != "" {
			strs[n] = s
			n++
		}
	}
	return strs[:n]
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
