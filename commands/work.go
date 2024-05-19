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
	panic("unimplemented")
}

// Parameterize implements Command.
func (w Work) Parameterize([]string) error {
	panic("unimplemented")
}
