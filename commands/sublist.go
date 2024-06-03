package commands

import (
	"fmt"
	"strings"

	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

const SublistCommandName CommandName = "sublist"

type FilterByTags func(tags []types.ID, repo data.Repository) ([]types.Actionable, error)

type Sublist struct {
	repo         data.Repository
	tags         []types.ID
	display      Displayer
	filterByTags FilterByTags
}

func (s *Sublist) Parameterize(args []string) error {
	if len(args) < 1 {
		return errArgsTooFew(1, len(args))
	}
	s.tags = apply.ToSlice(args, func(arg string) types.ID { return types.ID(arg) })
	return nil
}

func (s *Sublist) Execute() error {
	tagged, err := s.filterByTags(s.tags, s.repo)
	if err != nil {
		return err
	}
	for _, a := range tagged {
		s.display(a)
	}
	return nil
}

func (s *Sublist) String() string {
	return fmt.Sprintf("command %s %s", SublistCommandName, strings.Join(apply.ToSlice(s.tags, types.IDToString), " "))
}
