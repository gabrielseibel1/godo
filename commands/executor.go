package commands

import "log/slog"

type Executor func(Command) error

func ExecutorWithLog() Executor {
	return func(c Command) error {
		if err := c.Execute(); err != nil {
			slog.Error("failed", c.String(), "error", err)
			return err
		}
		slog.Info("executed", c.String(), "success")
		return nil
	}
}
