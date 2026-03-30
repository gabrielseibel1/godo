package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/types"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <id> [description]",
		Short: "Create a new activity",
		Args:  argsRange(1, 2, "godo create <id> [description]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := types.ID(args[0])
			if id == "" {
				return fmt.Errorf("id cannot be empty")
			}
			description := ""
			if len(args) == 2 {
				description = args[1]
			}
			a := types.NewActivity(id, description)
			return repo.Put(a)
		},
	}
}
