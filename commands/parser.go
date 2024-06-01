package commands

import (
	"strings"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

type Parser func([]string) (Command, error)

type Deps struct {
	Repo         data.Repository
	Displayer    Displayer
	Initializers []Initializer
}

func NewParser(deps Deps) Parser {
	return func(tokens []string) (Command, error) {
		// we should have the program name and an argument (command)
		if len(tokens) < 2 {
			return nil, types.ErrUnparsable(strings.Join(tokens, " "))
		}

		// find the corresponding command of the first token (after prog name)
		var cmd Command
		switch tokens[1] {
		case string(InitCommandName):
			cmd = &Init{initializers: deps.Initializers}
		case string(ListCommandName):
			cmd = &List{repo: deps.Repo, display: deps.Displayer}
		case string(SublistCommandName):
			cmd = &Sublist{repo: deps.Repo, display: deps.Displayer}
		case string(GetCommandName):
			cmd = &Get{repo: deps.Repo, display: deps.Displayer}
		case string(CreateCommandName):
			cmd = &Create{repo: deps.Repo}
		case string(DeleteCommandName):
			cmd = &Delete{repo: deps.Repo}
		case string(DoCommandName):
			cmd = &Do{repo: deps.Repo}
		case string(UndoCommandName):
			cmd = &Undo{repo: deps.Repo}
		case string(WorkCommandName):
			cmd = &Work{repo: deps.Repo}
		case string(TagCommandName):
			cmd = &Tag{repo: deps.Repo}
		case string(UntagCommandName):
			cmd = &Untag{repo: deps.Repo}
		case string(HelpCommandName):
			cmd = &Help{}
		default:
			return nil, types.ErrUnparsable(tokens[1])
		}

		// in general we pass the next tokens as args
		// but in the case of help the arg is special
		// TODO review this
		args := tokens[2:]
		if _, ok := cmd.(*Help); ok {
			args = tokens[:1]
		}

		// parameterize the command with args
		if err := cmd.Parameterize(args); err != nil {
			return nil, err
		}
		return cmd, nil
	}
}
