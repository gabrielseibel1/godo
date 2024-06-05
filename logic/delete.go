package logic

import (
	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

func DeleteFrom(repo data.Repository) func(types.ID) error { return repo.Remove }
