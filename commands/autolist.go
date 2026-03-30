package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func newAutoListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auto-list [since-date]",
		Short: "Show auto-work progress",
		Long: `Show logged auto-work entries with period details.

Defaults to today. Pass a date (YYYY-MM-DD) or "month" to show entries from a given date onward.

Examples:
  godo auto-list                # today only
  godo auto-list month          # from the 1st of the current month
  godo auto-list 2026-03-01    # from March 1st onward`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			now := time.Now()
			since := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

			if len(args) == 1 {
				if args[0] == "month" {
					since = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
				} else {
					d, err := time.Parse("2006-01-02", args[0])
					if err != nil {
						return fmt.Errorf("invalid date %q (expected YYYY-MM-DD or \"month\")", args[0])
					}
					since = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
				}
			}

			as, err := repo.List()
			if err != nil {
				return err
			}

			var totalDuration time.Duration
			found := false

			for _, a := range as {
				id := string(a.Identity())
				if !strings.HasPrefix(id, "work-") {
					continue
				}
				datePart := id[len("work-"):]
				actDate, err := time.ParseInLocation("2006-01-02", datePart, time.Local)
				if err != nil {
					continue
				}
				if actDate.Before(since) {
					continue
				}

				periods := a.Periods()
				if len(periods) == 0 && a.Worked() == 0 {
					continue
				}

				found = true
				fmt.Printf("%s  %s\n", datePart, formatDuration(a.Worked()))
				for _, p := range periods {
					dur := p.End.Sub(p.Start)
					fmt.Printf("  %s-%s  %s\n", p.Start.Format("15:04"), p.End.Format("15:04"), formatDuration(dur))
				}
				totalDuration += a.Worked()
			}

			if !found {
				if len(args) == 1 {
					fmt.Printf("No auto-work entries since %s\n", since.Format("2006-01-02"))
				} else {
					fmt.Println("No auto-work entries today")
				}
				return nil
			}

			// Show total if multiple days
			if len(args) == 1 {
				fmt.Printf("\nTotal: %s\n", formatDuration(totalDuration))
			}

			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			// Suggest recent dates
			now := time.Now()
			var suggestions []string
			for i := 0; i < 7; i++ {
				d := now.AddDate(0, 0, -i)
				suggestions = append(suggestions, d.Format("2006-01-02"))
			}
			// Also suggest "month" keyword
			suggestions = append(suggestions, "month")
			// Also suggest first of current month
			firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
			suggestions = append(suggestions, firstOfMonth.Format("2006-01-02"))
			return dedupStrings(suggestions), cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

func dedupStrings(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	var result []string
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
