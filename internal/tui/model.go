package tui

import (
	"context"
	"database/sql"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func run() error {
	ctx := context.TODO()
	db, err := sql.Open("pgx", "postgres:///lr")
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, "select price, city from land_registry_price_paid_uk where city=$1 and price>=$2 and price<=$3 limit 1000", "FOLKESTONE", "300000", "310000")
	if err != nil {
		return err
	}
	defer rows.Close()

	m := model{}
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
	go func() {
		for rows.Next() {
			var res result
			if err := rows.Scan(&res.price, &res.city); err != nil {
				p.Send(err)
			}
			p.Send(res)
		}
		if err := rows.Err(); err != nil {
			p.Send(err)
		}
	}()
	_, err = p.Run()
	return err
}

type result struct {
	price int
	city  string
}

type model struct {
	width   int
	height  int
	results []result
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
	case result:
		m.results = append(m.results, msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		Headers("PRICE", "CITY")
	for _, res := range m.results {
		t.Row(strconv.Itoa(res.price), res.city)
	}
	return "peppermint\n" + t.String()
}
