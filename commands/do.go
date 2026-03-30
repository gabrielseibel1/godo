package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
	"github.com/spf13/cobra"
)

func newDoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "do <id>",
		Short: "Mark an activity as done",
		Args:  argsExact(1, "godo do <id>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := types.ID(args[0])
			a, err := repo.Get(id)
			if err == data.ErrNotFound {
				return fmt.Errorf("activity %q not found", args[0])
			}
			if err != nil {
				return err
			}
			a.Do()
			return repo.Put(a)
		},
		ValidArgsFunction: idCompletionFunc,
	}
}
