package presentation

import (
	"fmt"
	"strings"

	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/godo/types"
)

type Item struct {
	title, desc string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title + i.desc }
func (i Item) String() string      { return i.title + " : " + i.desc }

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
	return Item{
		title: string(actionable.Identify()),
		desc: fmt.Sprintf("%s (%s) %s %s",
			checkbox,
			actionable.Worked().String(),
			strings.Join(apply.ToSlice(actionable.Tags(), func(id types.ID) string { return fmt.Sprintf("[%s]", string(id)) }), ","),
			actionable.Describe(),
		),
	}
}

func DisplayItem(a types.Actionable) string { return ItemFromActionable(a).String() }
