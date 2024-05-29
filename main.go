package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
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
	path := filepath.Join(data.Dir, data.JSONFile)

	repo := data.NewJSONRepository(
		data.FileReader(path),
		data.FileWriter(path),
		data.JSONDecode,
		data.JSONEncode,
		data.Compare,
	)

	// if no extra arguments passed, show UI
	if len(os.Args) < 2 {
		showUI(repo)
		return
	}

	// something was passed as arg, interpret as command
	runCommand(repo, path)
}

func runCommand(repo data.Repository, path string) {
	// parse command from arguments
	parse := commands.NewParser(commands.Deps{
		Repo:         repo,
		Displayer:    presentation.DisplayItem,
		Initializers: []commands.Initializer{data.DotGodoDirCreater(filepath.Base("")), data.FileCreater(path)},
	})
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

func showUI(repo data.Repository) {
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
