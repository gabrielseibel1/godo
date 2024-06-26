package presentation

import (
	"errors"
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
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(0)
	d.Styles.SelectedTitle.Foreground(lipgloss.Color("9"))
	d.Styles.SelectedTitle.BorderForeground(lipgloss.Color("9"))
	m := ListModel{
		list:  list.New(make([]list.Item, 0), d, 0, 0),
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
		items := apply.ToSlice(msg, func(a types.Actionable) list.Item { return UIItem{Actionable: a, style: m.style} })
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

func (m ListModel) Selected() (int, types.Actionable) {
	i := m.list.Index()
	if sel := m.list.SelectedItem(); sel != nil {
		return i, sel.(UIItem).Actionable
	}
	return i, nil
}

func (m ListModel) Select(id types.ID) error {
	for i := 0; i < len(m.list.Items()); i++ {
		curr := m.list.Items()[i].(UIItem).Identity()
		if string(curr) == string(id) {
			m.list.Select(i)
			return nil
		}
	}
	return errors.New("not found")
}

func (m ListModel) SelectIndex(i int) error {
	m.list.Select(i)
	return nil
}
