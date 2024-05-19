package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

const CreateCommandName CommandName = "create"

type Create struct {
	actionable types.Actionable
	repo       data.Repository
}

// String implements Command.
func (c Create) String() string {
	return fmt.Sprintf("command %s %+v", CreateCommandName, c.actionable)
}

// Execute implements Command.
func (c Create) Execute() error {
	panic("unimplemented")
}

// Parameterize implements Command.
func (c Create) Parameterize([]string) error {
	panic("unimplemented")
}
