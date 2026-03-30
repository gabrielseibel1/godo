package commands

import (
	"os"
	"path/filepath"

	"github.com/gabrielseibel1/godo/data"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a .godo directory and godo.json file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := filepath.Join(data.Dir, data.JSONFile)
			if err := os.MkdirAll(data.Dir, os.ModeDir|os.ModePerm); err != nil {
				return err
			}
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			return f.Close()
		},
	}
}
