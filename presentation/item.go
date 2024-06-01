package presentation

import (
	"fmt"
	"strings"

	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/fungo/conv"
	"github.com/gabrielseibel1/godo/types"
)

type Item struct {
	title, desc string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title + i.desc }
func (i Item) String() string      { return fmt.Sprintf("%s -> %s", i.title, i.desc) }

func ItemFromActionable(actionable types.Actionable) Item {
	var checkbox string
	if actionable.Done() {
		checkbox = "✅"
	} else {
		if actionable.Worked() > 0 {
			checkbox = "📶"
		} else {
			checkbox = "🆕"
		}
	}
	var worked string
	if duration := actionable.Worked(); duration != 0 {
		worked = fmt.Sprintf("(Worked = %s) ", duration.String())
	}
	var tags string
	if tagsMap := actionable.Tags(); len(tagsMap) > 0 {
		tagsSlice := apply.ToSlice(
			conv.MapKeysToSlice(tagsMap),
			func(id types.ID) string {
				return fmt.Sprintf("[%s]", types.IDToString(id))
			},
		)
		tags = fmt.Sprintf("(Tags = %s) ", strings.Join(tagsSlice, ","))
	}
	var description string
	if desc := actionable.Describe(); desc != "" {
		if worked != "" || tags != "" {
			description = fmt.Sprintf("(Description = %s)", desc)
		} else {
			description = desc
		}
	}
	return Item{
		title: fmt.Sprintf("%s |%s|", checkbox, types.IDToString(actionable.Identify())),
		desc:  fmt.Sprintf("%s%s%s", worked, tags, description),
	}
}

func DisplayItem(a types.Actionable) string { return ItemFromActionable(a).String() }
