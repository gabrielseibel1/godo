package commands

import (
	"github.com/gabrielseibel1/godo/presentation"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all activities",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			as, err := repo.List()
			if err != nil {
				return err
			}
			for _, a := range as {
				presentation.PrintItem(a)
			}
			return nil
		},
	}
}
