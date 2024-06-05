package presentation

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/fungo/conv"
	"github.com/gabrielseibel1/godo/types"
)

type UIItem struct {
	types.Actionable
}

func (i UIItem) Title() string {
	return string(i.Identify())
}

func (i UIItem) Description() string {
	return fmt.Sprintf("%s %s",
		checkbox(i.Done(), i.Worked()),
		lipgloss.NewStyle().Bold(true).SetString(i.Describe()),
	)
}

func (i UIItem) FilterValue() string { return i.Title() }

type CommandItem struct {
	types.Actionable
}

func (c CommandItem) String() string {
	return fmt.Sprintf("%s %s -(%s)-> %s ~ %s",
		checkbox(c.Done(), c.Worked()),
		title(c.Identify(), c.Done(), c.Worked()),
		c.Worked(),
		c.Describe(),
		tags(c.Tags()),
	)
}

func PrintItem(a types.Actionable) {
	fmt.Println(CommandItem{Actionable: a}.String())
}

func checkbox(done bool, worked time.Duration) string {
	if done {
		return "✅"
	}
	if worked > 0 {
		return "📶"
	} else {
		return "🆕"
	}
}

func title(id types.ID, done bool, worked time.Duration) string {
	style := lipgloss.NewStyle().SetString(fmt.Sprintf("\"%s\"", string(id))).Bold(true)
	if done {
		return style.Foreground(lipgloss.Color("10")).String() // green
	}
	if worked > 0 {
		return style.Foreground(lipgloss.Color("3")).String() // yellow
	}
	return style.String() // no color
}

func tags(t map[types.ID]struct{}) string {
	if len(t) > 0 {
		tagsSlice := apply.ToSlice(
			conv.MapKeysToSlice(t),
			func(id types.ID) string {
				return fmt.Sprintf("[%s]", types.IDToString(id))
			},
		)
		return strings.Join(tagsSlice, ",")
	}
	return ""
}
