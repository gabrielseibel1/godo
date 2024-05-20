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
	a, err := u.repo.Get(u.id)
	if err != nil {
		return err
	}
	a.Undo()
	return u.repo.Put(a)
}

// Parameterize implements Command.
func (u *Undo) Parameterize(args []string) error {
	if len(args) != 1 {
		return errArgsCount(1, len(args))
	}
	u.id = types.ID(args[0])
	return nil
}
