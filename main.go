package main

import (
	"log/slog"
	"os"

	"github.com/gabrielseibel1/godo/commands"
	"github.com/gabrielseibel1/godo/data"
)

func main() {
	if len(os.Args) < 1 {
		panic("no program name")
	}
	name := os.Args[0]
	if len(os.Args) < 2 {
		help := &commands.Help{}
		if err := help.Parameterize([]string{name}); err != nil {
			panic(err)
		}
		if err := help.Execute(); err != nil {
			panic(err)
		}
		return
	}

	// prepare persistency layer
	repo := data.LoggerOverRepository(data.FileJSONRepository())

	// parse command
	parse := commands.ParserWithRepository(repo)
	command, err := parse(os.Args)
	if err != nil {
		panic(err)
	}

	// run command
	execute := commands.ExecutorWithLog()
	if err := execute(command); err != nil {
		if err == data.ErrNotFound {
			slog.Error("error", err)
		} else {
			panic(err)
		}
	}
}
