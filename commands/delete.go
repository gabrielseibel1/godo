package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

const DeleteCommandName CommandName = "delete"

type Delete struct {
	id   types.ID
	repo data.Repository
}

// String implements Command.
func (d Delete) String() string {
	return fmt.Sprintf("command %s %s", DeleteCommandName, d.id)
}

// Execute implements Command.
func (d Delete) Execute() error {
	return d.repo.Remove(d.id)
}

// Parameterize implements Command.
func (d *Delete) Parameterize(args []string) error {
	if len(args) != 1 {
		return errArgsCount(1, len(args))
	}
	d.id = types.ID(args[0])
	return nil
}
