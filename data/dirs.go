package data

import (
	"os"
	"path/filepath"
)

const Dir = ".godo"

func DotGodoDirCreater(parent string) func() error {
	return func() error { return os.MkdirAll(filepath.Join(parent, Dir), os.ModeDir|os.ModePerm) }
}
