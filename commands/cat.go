package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/logic"
	"github.com/spf13/cobra"
)

func newCatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cat",
		Short: "List all tags",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			as, err := repo.List()
			if err != nil {
				return err
			}
			tags := logic.ListTags(as)
			for _, tag := range tags {
				fmt.Println(tag)
			}
			return nil
		},
	}
}
