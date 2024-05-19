package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

const UndoCommandName CommandName = "undo"

type Undo struct {
	id   types.ID
	repo data.Repository
}

// String implements Command.
func (u Undo) String() string {
	return fmt.Sprintf("command %s %s", UndoCommandName, u.id)
}

// Execute implements Command.
func (u Undo) Execute() error {
	panic("unimplemented")
}

// Parameterize implements Command.
func (u Undo) Parameterize([]string) error {
	panic("unimplemented")
}
