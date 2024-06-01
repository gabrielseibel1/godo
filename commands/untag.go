package commands

import (
	"fmt"
	"strings"

	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

const UntagCommandName CommandName = "untag"

type Untag struct {
	repo data.Repository
	ids  []types.ID
	tag  types.ID
}

func (u *Untag) Parameterize(args []string) error {
	if len(args) < 2 {
		return errArgsTooFew(2, len(args))
	}
	for _, arg := range args[:len(args)-1] {
		u.ids = append(u.ids, types.ID(arg))
	}
	u.tag = types.ID(args[len(args)-1])
	return nil
}

func (u *Untag) Execute() error {
	for _, id := range u.ids {
		a, err := u.repo.Get(id)
		if err != nil {
			return err
		}
		a.RemoveTag(u.tag)
		if err := u.repo.Put(a); err != nil {
			return err
		}
	}
	return nil
}

func (u *Untag) String() string {
	return fmt.Sprintf("command %s %s %s", UntagCommandName, strings.Join(apply.ToSlice(u.ids, func(id types.ID) string { return string(id) }), " "), u.tag)
}
