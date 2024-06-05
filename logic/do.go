package logic

import (
	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

func DoFrom(repo data.Repository) func(types.ID) error {
	return func(id types.ID) error {
		a, err := repo.Get(id)
		if err != nil {
			return err
		}
		a.Do()
		return repo.Put(a)
	}
}

func UndoFrom(repo data.Repository) func(types.ID) error {
	return func(id types.ID) error {
		a, err := repo.Get(id)
		if err != nil {
			return err
		}
		a.Undo()
		return repo.Put(a)
	}
}
