package data

import (
	"errors"

	"github.com/gabrielseibel1/godo/types"
)

type Getter interface {
	Get(types.ID) (types.Actionable, error)
}

type Lister interface {
	List() ([]types.Actionable, error)
}

type Putter interface {
	Put(types.Actionable) error
}

type Remover interface {
	Remove(types.ID) error
}

type Repository interface {
	Getter
	Lister
	Putter
	Remover
}

var ErrNotFound = errors.New("not found")
