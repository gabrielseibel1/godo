package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
	"github.com/spf13/cobra"
)

func newWorkCmd() *cobra.Command {
	var dateStr string

	cmd := &cobra.Command{
		Use:   "work <id> <duration-or-range>",
		Short: "Log work on an activity",
		Long: `Log work on an activity using either a duration or a time range.

Examples:
  godo work my-task 3h          # duration
  godo work my-task 18:00-21:00 # time range (needs --date or ID with date)`,
		Args: argsExact(2, "godo work <id> <duration-or-range>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := types.ID(args[0])
			a, err := repo.Get(id)
			if err == data.ErrNotFound {
				return fmt.Errorf("activity %q not found", args[0])
			}
			if err != nil {
				return err
			}

			input := args[1]
			if strings.Contains(input, "-") && strings.Contains(input, ":") {
				// Parse as time range HH:MM-HH:MM
				start, end, err := parseTimeRange(input, dateStr, string(id))
				if err != nil {
					return err
				}
				a.WorkPeriod(start, end)
			} else {
				// Parse as duration
				duration, err := time.ParseDuration(input)
				if err != nil {
					return err
				}
				a.Work(duration)
			}
			return repo.Put(a)
		},
		ValidArgsFunction: idCompletionFunc,
	}

	cmd.Flags().StringVar(&dateStr, "date", "", "date for the time range (YYYY-MM-DD)")

	return cmd
}

// parseTimeRange parses "HH:MM-HH:MM" into start/end time.Time values.
// If dateStr is empty, tries to infer date from activityID (extra-hour-YYYY-MM-DD), else uses today.
// If end < start, assumes it crosses midnight.
func parseTimeRange(input, dateStr, activityID string) (time.Time, time.Time, error) {
	parts := strings.SplitN(input, "-", 2)
	// Handle HH:MM-HH:MM where HH:MM contains ':'
	// We need to split on the '-' that separates the two times
	// Find the '-' that's between two time strings
	dashIdx := -1
	for i := 1; i < len(input)-1; i++ {
		if input[i] == '-' && i > 2 {
			dashIdx = i
			break
		}
	}
	if dashIdx == -1 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid time range format: %s (expected HH:MM-HH:MM)", input)
	}
	parts = []string{input[:dashIdx], input[dashIdx+1:]}

	date := resolveDate(dateStr, activityID)

	startTime, err := time.Parse("15:04", parts[0])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time: %w", err)
	}

	// Handle "24:00" as midnight
	var endTime time.Time
	endNextDay := false
	if parts[1] == "24:00" {
		endTime, _ = time.Parse("15:04", "00:00")
		endNextDay = true
	} else {
		endTime, err = time.Parse("15:04", parts[1])
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end time: %w", err)
		}
	}

	start := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.Local)
	end := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.Local)

	// If end is before start or explicitly next day, it crosses midnight
	if endNextDay || end.Before(start) || end.Equal(start) {
		end = end.AddDate(0, 0, 1)
	}

	return start, end, nil
}

// resolveDate determines the date from explicit flag, activity ID pattern, or today.
func resolveDate(dateStr, activityID string) time.Time {
	if dateStr != "" {
		if d, err := time.Parse("2006-01-02", dateStr); err == nil {
			return d
		}
	}
	// Try to extract date from activity ID suffix (e.g. work-2026-03-02)
	if len(activityID) >= 10 {
		suffix := activityID[len(activityID)-10:]
		if d, err := time.Parse("2006-01-02", suffix); err == nil {
			return d
		}
	}
	return time.Now()
}
