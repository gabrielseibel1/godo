package commands

import (
	"github.com/gabrielseibel1/godo/types"
	"github.com/spf13/cobra"
)

func newUntagCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "untag <id>... <tag>",
		Short: "Remove a tag from one or more activities",
		Args:  argsMin(2, "godo untag <id>... <tag>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			tag := types.ID(args[len(args)-1])
			for _, idStr := range args[:len(args)-1] {
				a, err := repo.Get(types.ID(idStr))
				if err != nil {
					return err
				}
				a.RemoveTag(tag)
				if err := repo.Put(a); err != nil {
					return err
				}
			}
			return nil
		},
		ValidArgsFunction: idCompletionFunc,
	}
}
