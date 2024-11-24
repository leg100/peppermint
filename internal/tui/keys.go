package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	ExplorerFullScreen    key.Binding
	TogglePreview         key.Binding
	ShrinkVerticalSplit   key.Binding
	GrowVerticalSplit     key.Binding
	ShrinkHorizontalSplit key.Binding
	GrowHorizontalSplit   key.Binding
	SwitchPane            key.Binding
}

var Keys = keyMap{
	ExplorerFullScreen: key.NewBinding(
		key.WithKeys("E"),
		key.WithHelp("E", "toggle full explorer"),
	),
	TogglePreview: key.NewBinding(
		key.WithKeys("P"),
		key.WithHelp("P", "toggle preview"),
	),
	ShrinkVerticalSplit: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "shink vertical split"),
	),
	GrowVerticalSplit: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "grow vertical split"),
	),
	ShrinkHorizontalSplit: key.NewBinding(
		key.WithKeys("<"),
		key.WithHelp("<", "shink horizontal split"),
	),
	GrowHorizontalSplit: key.NewBinding(
		key.WithKeys(">"),
		key.WithHelp(">", "grow horizontal split"),
	),
	SwitchPane: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch pane"),
	),
}
