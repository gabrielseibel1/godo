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
	return c.repo.Put(c.actionable)
}

// Parameterize implements Command.
func (c *Create) Parameterize(args []string) error {
	if len(args) != 2 {
		return errArgsCount(2, len(args))
	}
	id, description := types.ID(args[0]), args[1]
	if id == "" {
		return fmt.Errorf("no id")
	}
	if description == "" {
		return fmt.Errorf("no description")
	}
	c.actionable = types.NewActivity(id, description)
	return nil
}
