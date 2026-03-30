package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an activity",
		Args:  argsExact(1, "godo delete <id>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := types.ID(args[0])
			if _, err := repo.Get(id); err == data.ErrNotFound {
				return fmt.Errorf("activity %q not found", args[0])
			}
			return repo.Remove(id)
		},
		ValidArgsFunction: idCompletionFunc,
	}
}
