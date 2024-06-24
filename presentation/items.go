package presentation

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/fungo/conv"
	"github.com/gabrielseibel1/godo/types"
)

func PrintItem(a types.Actionable) {
	fmt.Println(CommandItem{Actionable: a}.String())
}

type CommandItem struct {
	types.Actionable
	style lipgloss.Style
}

func (c CommandItem) String() string {
	return fmt.Sprintf("%s %s -(%s)-> %s ~ %s",
		checkbox(c.Done(), c.Worked()),
		c.title(c.Identity(), c.Done(), c.Worked()),
		c.Worked(),
		c.Description(),
		tags(c.Tags()),
	)
}

func (c CommandItem) title(id types.ID, done bool, worked time.Duration) string {
	style := c.style.SetString(fmt.Sprintf("\"%s\"", string(id))).Bold(true)
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
		slices.Sort(tagsSlice)
		return strings.Join(tagsSlice, ",")
	}
	return ""
}

func checkbox(done bool, worked time.Duration) string {
	if done {
		return "✅"
	}
	if worked > 0 {
		return "📶"
	}
	return "🆕"
}

type UIItem struct {
	types.Actionable
	style lipgloss.Style
}

func (i UIItem) Title() string {
	return i.style.Render(fmt.Sprintf("%s %s (%s)",
		checkbox(i.Done(), i.Worked()),
		fmt.Sprintf("\"%s\"", string(i.Identity())),
		i.Worked(),
	))
}

func (i UIItem) FilterValue() string { return i.Title() }
