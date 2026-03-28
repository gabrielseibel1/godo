package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
	"github.com/spf13/cobra"
)

func newAutoWorkCmd() *cobra.Command {
	var dateStr string
	var yesterday bool

	cmd := &cobra.Command{
		Use:   "auto-work <HH:MM-HH:MM>",
		Short: "Log work automatically (auto-creates activity for the date)",
		Long: `Log a work period without specifying an activity. Automatically creates a date-based activity if it doesn't exist.

The activity ID is derived from the date: work-YYYY-MM-DD.

Date resolution:
  --date YYYY-MM-DD   use explicit date
  --yesterday          use yesterday's date
  (default)            if before 05:00, uses yesterday; otherwise today

Examples:
  godo auto-work 18:00-21:00                    # today (or yesterday if before 5am)
  godo auto-work 17:00..                        # from 17:00 until now
  godo auto-work 17:00-now                      # same as above
  godo auto-work 18:00-21:00 --date 2026-03-02  # specific date
  godo auto-work 19:00-01:00 --date 2026-03-13  # crosses midnight
  godo auto-work 18:00-23:00 --yesterday        # yesterday's date`,
		Args: argsExact(1, "godo auto-work <HH:MM-HH:MM>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			date := resolveAutoWorkDate(dateStr, yesterday)
			activityID := fmt.Sprintf("work-%s", date.Format("2006-01-02"))

			input := normalizeNow(args[0])
			start, end, err := parseTimeRange(input, date.Format("2006-01-02"), activityID)
			if err != nil {
				return err
			}

			// Get or create activity
			a, err := repo.Get(types.ID(activityID))
			if err == data.ErrNotFound {
				a = types.NewActivity(types.ID(activityID), "")
			} else if err != nil {
				return err
			}

			a.WorkPeriod(start, end)
			if err := repo.Put(a); err != nil {
				return err
			}

			duration := end.Sub(start)
			fmt.Printf("Logged %s on %s (%s-%s)\n",
				formatDuration(duration),
				date.Format("2006-01-02"),
				start.Format("15:04"),
				end.Format("15:04"),
			)
			return nil
		},
	}

	cmd.Flags().StringVar(&dateStr, "date", "", "explicit date (YYYY-MM-DD)")
	cmd.Flags().BoolVar(&yesterday, "yesterday", false, "use yesterday's date")

	return cmd
}

// resolveAutoWorkDate determines the date for auto-work.
// Priority: --date > --yesterday > smart-today (before 5am = yesterday)
func resolveAutoWorkDate(dateStr string, yesterday bool) time.Time {
	if dateStr != "" {
		if d, err := time.Parse("2006-01-02", dateStr); err == nil {
			return d
		}
	}
	now := time.Now()
	if yesterday || now.Hour() < 5 {
		return now.AddDate(0, 0, -1)
	}
	return now
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh%02dm", h, m)
}

// normalizeNow replaces "HH:MM.." and "HH:MM-now" with "HH:MM-<current time>".
func normalizeNow(input string) string {
	now := time.Now().Format("15:04")
	if strings.HasSuffix(input, "..") {
		return strings.TrimSuffix(input, "..") + "-" + now
	}
	if strings.HasSuffix(input, "-now") {
		return strings.TrimSuffix(input, "-now") + "-" + now
	}
	return input
}
