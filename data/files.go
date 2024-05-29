package data

import "os"

func FileCreater(path string) func() error {
	return func() error {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		return f.Close()
	}
}
