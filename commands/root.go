package commands

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/logic"
	"github.com/gabrielseibel1/godo/presentation"
)

var repo data.Repository

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "godo",
		Short: "A todo list with a TUI and CLI",
		Long:  "GoDo is a todo list application with a terminal UI and command-line interface.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return showTUI()
		},
		SilenceUsage: true,
	}

	rootCmd.AddCommand(
		newCreateCmd(),
		newListCmd(),
		newGetCmd(),
		newDeleteCmd(),
		newDoCmd(),
		newUndoCmd(),
		newWorkCmd(),
		newAutoWorkCmd(),
		newAutoListCmd(),
		newTagCmd(),
		newUntagCmd(),
		newCatCmd(),
		newSublistCmd(),
		newVersionCmd(),
		newInitCmd(),
	)

	return rootCmd
}

func Execute() {
	path := filepath.Join(data.Dir, data.JSONFile)
	repo = data.NewJSONRepository(
		data.FileReader(path),
		data.FileWriter(path),
		data.JSONDecode,
		data.JSONEncode,
		data.Compare,
	)

	rootCmd := NewRootCmd()
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(rootCmd.ErrOrStderr(), "Error:", err)
	}
}

func showTUI() error {
	dir, err := filepath.Abs(".")
	if err != nil {
		dir = ""
	}
	mt := presentation.NewTabbedListModel(
		dir,
		logic.Version(),
		presentation.NewListModel(lipgloss.NewStyle()),
		presentation.NewEditorModel(),
		logic.DoFrom(repo),
		logic.UndoFrom(repo),
		logic.DeleteFrom(repo),
		logic.EditFrom(repo),
	)
	p := tea.NewProgram(mt, tea.WithAltScreen())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	realtime := presentation.NewRealtimeSync(p, repo, ctx, time.NewTicker(time.Millisecond*15))
	go realtime.KeepSynched()

	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

// argsExact returns a cobra.PositionalArgs that requires exactly n args, showing usage on error.
func argsExact(n int, usage string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("missing required argument(s)\n\nUsage:\n  %s", usage)
		}
		return nil
	}
}

// argsRange returns a cobra.PositionalArgs that requires between min and max args, showing usage on error.
func argsRange(min, max int, usage string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < min || len(args) > max {
			return fmt.Errorf("missing required argument(s)\n\nUsage:\n  %s", usage)
		}
		return nil
	}
}

// argsMin returns a cobra.PositionalArgs that requires at least n args, showing usage on error.
func argsMin(n int, usage string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf("missing required argument(s)\n\nUsage:\n  %s", usage)
		}
		return nil
	}
}

// idCompletionFunc provides completion for activity IDs
func idCompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if repo == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	items, err := repo.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ids []string
	for _, item := range items {
		ids = append(ids, string(item.Identity()))
	}
	return ids, cobra.ShellCompDirectiveNoFileComp
}
