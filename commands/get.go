package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

const GetCommandName CommandName = "get"

type Get struct {
	id      types.ID
	repo    data.Repository
	display Displayer
}

func (g *Get) Parameterize(args []string) error {
	if len(args) != 1 {
		return errArgsCount(1, len(args))
	}
	g.id = types.ID(args[0])
	return nil
}

func (g *Get) Execute() error {
	a, err := g.repo.Get(g.id)
	if err != nil {
		return err
	}
	fmt.Println(g.display(a))
	return nil
}

func (g *Get) String() string {
	return fmt.Sprintf("command %s %s", GetCommandName, g.id)
}
