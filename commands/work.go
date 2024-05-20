package commands

import (
	"fmt"
	"time"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

const WorkCommandName CommandName = "work"

type Work struct {
	id       types.ID
	duration time.Duration
	repo     data.Repository
}

// String implements Command.
func (w Work) String() string {
	return fmt.Sprintf("command %s %s %s", WorkCommandName, w.id, w.duration)
}

// Execute implements Command.
func (w Work) Execute() error {
	a, err := w.repo.Get(w.id)
	if err != nil {
		return err
	}
	a.Work(w.duration)
	return w.repo.Put(a)
}

// Parameterize implements Command.
func (w *Work) Parameterize(args []string) error {
	if len(args) != 2 {
		return errArgsCount(2, len(args))
	}
	w.id = types.ID(args[0])
	duration, err := time.ParseDuration(args[1])
	if err != nil {
		return err
	}
	w.duration = duration
	return nil
}
