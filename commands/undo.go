package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
	"github.com/spf13/cobra"
)

func newUndoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "undo <id>",
		Short: "Mark an activity as not done",
		Args:  argsExact(1, "godo undo <id>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := types.ID(args[0])
			a, err := repo.Get(id)
			if err == data.ErrNotFound {
				return fmt.Errorf("activity %q not found", args[0])
			}
			if err != nil {
				return err
			}
			a.Undo()
			return repo.Put(a)
		},
		ValidArgsFunction: idCompletionFunc,
	}
}
