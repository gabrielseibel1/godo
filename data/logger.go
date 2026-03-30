package data

import (
	"log/slog"

	"github.com/gabrielseibel1/godo/types"
)

type Logger struct {
	inner Repository
}

func LoggerOverRepository(repo Repository) Repository {
	return Logger{inner: repo}
}

// Get implements Repository.
func (l Logger) Get(id types.ID) (types.Actionable, error) {
	a, err := l.inner.Get(id)
	if err != nil {
		slog.Error("repository", "op", "get", "err", err)
	} else {
		slog.Info("repository", "get", "success")
	}
	return a, err
}

// List implements Repository.
func (l Logger) List() ([]types.Actionable, error) {
	a, err := l.inner.List()
	if err != nil {
		slog.Error("repository", "op", "list", "err", err)
	} else {
		slog.Info("repository", "list", "success")
	}
	return a, err
}

// Put implements Repository.
func (l Logger) Put(a types.Actionable) error {
	err := l.inner.Put(a)
	if err != nil {
		slog.Error("repository", "op", "put", "err", err)
	} else {
		slog.Info("repository", "put", "success")
	}
	return err
}

// Remove implements Repository.
func (l Logger) Remove(id types.ID) error {
	err := l.inner.Remove(id)
	if err != nil {
		slog.Error("repository", "op", "remove", "err", err)
	} else {
		slog.Info("repository", "remove", "success")
	}
	return err
}
