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
	return fmt.Sprintf("command %s %s", DoCommandName, d.id)
}

// Execute implements Command.
func (d Do) Execute() error {
	panic("unimplemented")
}

// Parameterize implements Command.
func (d Do) Parameterize([]string) error {
	panic("unimplemented")
}
