package commands

import (
	"strings"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

type Parser func([]string) (Command, error)

func ParserWithRepository(repo data.Repository) Parser {
	return func(tokens []string) (Command, error) {
		if len(tokens) < 2 {
			return nil, types.ErrUnparsable(strings.Join(tokens, " "))
		}
		// find the corresponding command of the first token
		cmd, err := tokenToCommandWithRepository(tokens[1], repo)
		if err != nil {
			return nil, err
		}

		// in general we pass the next tokens as args
		// but in the case of help the arg is special
		args := tokens[2:]
		if _, ok := cmd.(*Help); ok {
			args = tokens[:1]
		}

		// parameterize the command with args
		if err = cmd.Parameterize(args); err != nil {
			return nil, err
		}
		return cmd, nil
	}
}

func tokenToCommandWithRepository(token string, repo data.Repository) (Command, error) {
	switch token {
	case string(ListCommandName):
		return &List{repo: repo}, nil
	case string(CreateCommandName):
		return &Create{repo: repo}, nil
	case string(DeleteCommandName):
		return &Delete{repo: repo}, nil
	case string(DoCommandName):
		return &Do{repo: repo}, nil
	case string(UndoCommandName):
		return &Undo{repo: repo}, nil
	case string(WorkCommandName):
		return &Work{repo: repo}, nil
	case string(GetCommandName):
		return &Get{repo: repo}, nil
	case string(HelpCommandName):
		return &Help{}, nil
	default:
		return nil, types.ErrUnparsable(token)
	}
}
