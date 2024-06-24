package commands

import (
	"fmt"
)

const InitCommandName CommandName = "init"

type Init struct {
	initializers []Initializer
}

func (i Init) Parameterize(args []string) error {
	if len(args) != 0 {
		return errArgsCount(0, len(args))
	}
	return nil
}

func (i Init) Execute() error {
	for _, initialize := range i.initializers {
		if err := initialize(); err != nil {
			return err
		}
	}
	return nil
}

func (i Init) String() string {
	return fmt.Sprintf("command %s", InitCommandName)
}

type DirCreater func() error
