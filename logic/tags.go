package logic

import (
	"github.com/gabrielseibel1/fungo/check"
	"github.com/gabrielseibel1/fungo/filter"
	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
	"golang.org/x/exp/maps"
)

func ListTags(repo data.Repository) ([]types.ID, error) {
	as, err := repo.List()
	if err != nil {
		return nil, err
	}
	tags := make(map[types.ID]struct{})
	for _, a := range as {
		maps.Copy(tags, a.Tags())
	}
	return maps.Keys(tags), nil
}

func SublistWithTags(tags []types.ID, repo data.Repository) ([]types.Actionable, error) {
	as, err := repo.List()
	if err != nil {
		return nil, err
	}
	tagged := filter.Slice(as, func(a types.Actionable) bool {
		return check.Some(tags, func(tag types.ID) bool {
			_, ok := a.Tags()[tag]
			return ok
		})
	})
	return tagged, nil
}
