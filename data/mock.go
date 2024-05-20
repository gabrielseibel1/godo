package data

import (
	"github.com/gabrielseibel1/godo/types"
	"golang.org/x/exp/maps"
)

// Mock is a data persistency strategy that has predefined values
type Mock struct {
	data map[types.ID]types.Actionable
}

func MockWithData(data map[types.ID]types.Actionable) Repository {
	return Mock{data: data}
}

// Initialize implements Repository.
func (m Mock) Initialize() error {
	return nil
}

// Get implements Repository.
func (m Mock) Get(id types.ID) (types.Actionable, error) {
	a, ok := m.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

// List implements Repository.
func (m Mock) List() ([]types.Actionable, error) {
	return maps.Values(m.data), nil
}

// Put implements Repository.
func (m Mock) Put(types.Actionable) error {
	return nil
}

// Remove implements Repository.
func (m Mock) Remove(types.ID) error {
	return nil
}
