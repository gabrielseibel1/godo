package commands

import (
	"fmt"
	"strings"

	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

const TagCommandName CommandName = "tag"

type Tag struct {
	repo data.Repository
	ids  []types.ID
	tag  types.ID
}

func (t *Tag) Parameterize(args []string) error {
	if len(args) < 2 {
		return errArgsTooFew(2, len(args))
	}
	for _, arg := range args[:len(args)-1] {
		t.ids = append(t.ids, types.ID(arg))
	}
	t.tag = types.ID(args[len(args)-1])
	return nil
}

func (t *Tag) Execute() error {
	for _, id := range t.ids {
		a, err := t.repo.Get(id)
		if err != nil {
			return err
		}
		a.Tag(t.tag)
		if err := t.repo.Put(a); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tag) String() string {
	return fmt.Sprintf("command %s %s %s",
		TagCommandName,
		strings.Join(apply.ToSlice(t.ids, func(id types.ID) string { return string(id) }), " "),
		t.tag,
	)
}
