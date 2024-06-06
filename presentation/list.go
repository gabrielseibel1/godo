package presentation

import (
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/godo/types"
)

type ListModel struct {
	list  list.Model
	style lipgloss.Style
}

func NewListModel(style lipgloss.Style) ListModel {
	m := ListModel{
		list:  list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0),
		style: style,
	}
	m.list.SetShowTitle(false)
	return m
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := m.style.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case itemsMsg:
		slices.SortFunc(msg, func(a, b types.Actionable) int {
			return strings.Compare(string(a.Identity()), string(b.Identity()))
		})
		items := apply.ToSlice(msg, func(a types.Actionable) list.Item { return UIItem{a} })
		return m, m.list.SetItems(items)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	return m.style.Render(m.list.View())
}

func (m ListModel) Items() []UIItem {
	return apply.ToSlice(m.list.Items(), func(i list.Item) UIItem { return i.(UIItem) })
}

func (m ListModel) Selected() types.Actionable {
	if sel := m.list.SelectedItem(); sel != nil {
		return sel.(UIItem).Actionable
	}
	return nil
}
