package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

const CatCommandName CommandName = "cat"

type TagLister func(as []types.Actionable) []types.ID

type Cat struct {
	repo     data.Repository
	listTags TagLister
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
	tags := c.listTags(as)
	for tag := range tags {
		fmt.Println(tag)
	}
	return nil
}

func (c *Cat) String() string {
	return fmt.Sprintf("command %s", CatCommandName)
}
