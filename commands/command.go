package commands

import "fmt"

type Executable interface {
	Execute() error
}

type Parameterizable interface {
	Parameterize([]string) error
}

type Command interface {
	Parameterizable
	Executable
	fmt.Stringer
}

type CommandName string
