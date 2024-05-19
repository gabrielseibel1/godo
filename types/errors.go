package types

import "fmt"

func ErrUnparsable(s string) error {
	return fmt.Errorf("cannot parse %s", s)
}
