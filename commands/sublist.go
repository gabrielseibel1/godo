package commands

import (
	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/godo/logic"
	"github.com/gabrielseibel1/godo/presentation"
	"github.com/gabrielseibel1/godo/types"
	"github.com/spf13/cobra"
)

func newSublistCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sublist <tag>...",
		Short: "List activities filtered by tags",
		Args:  argsMin(1, "godo sublist <tag>..."),
		RunE: func(cmd *cobra.Command, args []string) error {
			tags := apply.ToSlice(args, func(arg string) types.ID { return types.ID(arg) })
			as, err := repo.List()
			if err != nil {
				return err
			}
			tagged := logic.SublistWithTags(tags, as)
			for _, a := range tagged {
				presentation.PrintItem(a)
			}
			return nil
		},
	}
}
