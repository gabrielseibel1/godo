package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

const DoCommandName CommandName = "do"

type Do struct {
	id   types.ID
	repo data.Repository
}

// String implements Command.
func (d Do) String() string {
	if d.id == "" {
		return fmt.Sprintf("command %s", DoCommandName)
	}
	return fmt.Sprintf("command %s %s", DoCommandName, d.id)
}

// Execute implements Command.
func (d Do) Execute() error {
	a, err := d.repo.Get(d.id)
	if err != nil {
		return err
	}
	a.Do()
	return d.repo.Put(a)
}

// Parameterize implements Command.
func (d *Do) Parameterize(args []string) error {
	if len(args) != 1 {
		return errArgsCount(1, len(args))
	}
	d.id = types.ID(args[0])
	return nil
}
