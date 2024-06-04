package logic

import (
	"github.com/gabrielseibel1/fungo/check"
	"github.com/gabrielseibel1/fungo/filter"
	"github.com/gabrielseibel1/godo/types"
	"golang.org/x/exp/maps"
)

func ListTags(as []types.Actionable) []types.ID {
	tags := make(map[types.ID]struct{})
	for _, a := range as {
		maps.Copy(tags, a.Tags())
	}
	return maps.Keys(tags)
}

func SublistWithTags(tags []types.ID, as []types.Actionable) []types.Actionable {
	tagged := filter.Slice(as, func(a types.Actionable) bool {
		return check.Some(tags, func(tag types.ID) bool {
			_, ok := a.Tags()[tag]
			return ok
		})
	})
	return tagged
}
