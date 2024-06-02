package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
	"golang.org/x/exp/maps"
)

const CatCommandName CommandName = "cat"

type Cat struct {
	repo data.Repository
}

func (c *Cat) Parameterize(args []string) error {
	if len(args) != 0 {
		return errArgsCount(0, len(args))
	}
	return nil
}

func (c *Cat) Execute() error {
	as, err := c.repo.List()
	if err != nil {
		return err
	}
	tags := make(map[types.ID]struct{})
	for _, a := range as {
		maps.Copy(tags, a.Tags())
	}
	for tag := range tags {
		fmt.Println(tag)
	}
	return nil
}

func (c *Cat) String() string {
	return fmt.Sprintf("command %s", CatCommandName)
}
