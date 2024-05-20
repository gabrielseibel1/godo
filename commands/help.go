package commands

import (
	"fmt"
	"strings"
)

const HelpCommandName CommandName = "help"

type Help struct {
	name string
}

func (h *Help) Parameterize(args []string) error {
	if len(args) != 1 {
		return errArgsCount(1, len(args))
	}
	h.name = args[0]
	return nil
}

func (h *Help) Execute() error {
	fmt.Printf("usage: %s [command] [args]\n", h.name)
	fmt.Printf("\t[command] = {%s}\n",
		strings.Join([]string{
			string(ListCommandName),
			string(GetCommandName),
			string(CreateCommandName),
			string(DeleteCommandName),
			string(DoCommandName),
			string(UndoCommandName),
			string(WorkCommandName),
			string(HelpCommandName),
		}, "|"))
	return nil
}

func (h *Help) String() string {
	return "command help"
}
