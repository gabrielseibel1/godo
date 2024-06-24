package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
)

const ListCommandName CommandName = "list"

type List struct {
	repo    data.Repository
	display Displayer
}

// String implements Command.
func (l List) String() string {
	return fmt.Sprintf("command %s", ListCommandName)
}

// Execute implements Command.
func (l List) Execute() error {
	as, err := l.repo.List()
	if err != nil {
		return err
	}
	for _, a := range as {
		l.display(a)
	}
	return nil
}

// Parameterize implements Command.
func (l List) Parameterize(args []string) error {
	if len(args) != 0 {
		return errArgsCount(0, len(args))
	}
	return nil
}
