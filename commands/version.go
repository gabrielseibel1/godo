package commands

import (
	"fmt"

	"github.com/gabrielseibel1/godo/logic"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(logic.Version())
		},
	}
}
