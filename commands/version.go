package commands

import "fmt"

const VersionCommandName CommandName = "version"

const version = "v0.2.0"

type Version struct{}

func (v *Version) Parameterize(args []string) error {
	if len(args) > 0 {
		return errArgsCount(0, len(args))
	}
	return nil
}

func (v *Version) Execute() error {
	fmt.Println(version)
	return nil
}

func (v *Version) String() string {
	return fmt.Sprintf("command %s", VersionCommandName)
}
