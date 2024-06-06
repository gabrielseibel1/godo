package logic

import (
	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

func EditFrom(repo data.Repository) func(types.Actionable) error {
	return repo.Put
}
