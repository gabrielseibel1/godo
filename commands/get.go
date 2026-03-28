package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/presentation"
	"github.com/gabrielseibel1/godo/types"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an activity by ID",
		Args:  argsExact(1, "godo get <id>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := repo.Get(types.ID(args[0]))
			if err == data.ErrNotFound {
				return fmt.Errorf("activity %q not found", args[0])
			}
			if err != nil {
				return err
			}
			presentation.PrintItem(a)
			return nil
		},
		ValidArgsFunction: idCompletionFunc,
	}
}
