package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielseibel1/godo/commands"
	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/presentation"
)

func main() {
	// prepare persistency layer
	repo := data.FileJSONRepository()

	// if no extra arguments passed, show UI
	if len(os.Args) < 2 {
		uiList(repo)
		return
	}

	// parse command from arguments
	parse := commands.NewParser(repo, presentation.DisplayItem)
	command, err := parse(os.Args)
	if err != nil {
		panic(err)
	}

	// run command
	execute := commands.ExecutorWithoutLogs()
	if err := execute(command); err != nil {
		if err == data.ErrNotFound {
			slog.Error("error", err)
		} else {
			panic(err)
		}
	}
}

func uiList(repo data.Repository) {
	// model and program init
	m := presentation.NewModel(
		"GoDo - ToDo List",
		make([]list.Item, 0),
		lipgloss.NewStyle().Margin(1, 2),
	)
	p := tea.NewProgram(m, tea.WithAltScreen())

	// realtime data synchronization
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	realtime := presentation.NewRealtimeSync(p, repo, ctx, time.NewTicker(time.Second))
	go realtime.KeepSynched()

	// run until error or user exists
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
