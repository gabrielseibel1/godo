package commands

import "fmt"

func errArgsCount(exp, act int) error {
	return fmt.Errorf("expected %d arguments, got %d", exp, act)
}
