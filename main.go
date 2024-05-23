package main

import (
	"log/slog"
	"os"

	"github.com/gabrielseibel1/godo/commands"
	"github.com/gabrielseibel1/godo/data"
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
	parse := commands.ParserWithRepository(repo)
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
