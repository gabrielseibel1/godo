package commands

import "fmt"

const VersionCommandName CommandName = "version"

type Version struct {
	version string
}

func (v *Version) Parameterize(args []string) error {
	if len(args) > 0 {
		return errArgsCount(0, len(args))
	}
	return nil
}

func (v *Version) Execute() error {
	fmt.Println(v.version)
	return nil
}

func (v *Version) String() string {
	return fmt.Sprintf("command %s", VersionCommandName)
}
