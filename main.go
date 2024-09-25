package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/leg100/moto/internal/tui"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	m := tui.New()
	p := tea.NewProgram(m,
		// Use the full size of the terminal with its "alternate screen buffer"
		tea.WithAltScreen(),
		// Enabling mouse cell motion removes the ability to "blackboard" text
		// with the mouse, which is useful for then copying text into the
		// clipboard. Therefore we've decided to disable it and leave it
		// commented out for posterity.
		//
		// tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}
