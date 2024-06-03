package presentation

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielseibel1/fungo/apply"
)

type ListModel struct {
	title string
	list  list.Model
	style lipgloss.Style
}

func NewModel(title string, items []list.Item, style lipgloss.Style) ListModel {
	m := ListModel{
		title: title,
		list:  list.New(items, list.NewDefaultDelegate(), 0, 0),
		style: style,
	}
	m.list.Title = m.title
	return m
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := m.style.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case itemsMsg:
		return m, m.list.SetItems(msg)
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
