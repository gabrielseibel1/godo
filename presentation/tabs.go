package presentation

import (
	"fmt"
	"math"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/godo/logic"
	"github.com/gabrielseibel1/godo/types"
)

const allItems = "ALL"

type doer func(types.ID) error

type undoer func(types.ID) error

type deleter func(types.ID) error

type SelectableListModel interface {
	tea.Model
	Selected() types.Actionable
}

type TabbedListModel struct {
	title     lipgloss.Style
	listModel SelectableListModel
	tabs      []string
	activeTab int
	style     lipgloss.Style
	do        doer
	undo      undoer
	del       deleter
}

func NewTabbedListModel(path string, list SelectableListModel, do doer, undo undoer, del deleter) TabbedListModel {
	return TabbedListModel{
		title: lipgloss.NewStyle().
			// Bold(true).
			// Foreground(lipgloss.Color("10")).
			// Border(lipgloss.ThickBorder()).
			// BorderForeground(highlightColor).
			Padding(1, 1).
			SetString(fmt.Sprintf("%s\n%s", banner, path)),
		listModel: list,
		activeTab: 0,
		style:     lipgloss.NewStyle().Margin(2, 2),
		do:        do,
		undo:      undo,
		del:       del,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m TabbedListModel) Init() tea.Cmd {
	return m.listModel.Init()
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m TabbedListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "q", tea.KeyCtrlC.String():
			return m, tea.Quit
		case tea.KeyTab.String():
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
		case tea.KeyShiftTab.String():
			m.activeTab = int(math.Abs(float64((m.activeTab - 1) % len(m.tabs))))
		case "y":
			if selected := m.listModel.Selected(); selected != nil {
				if err := m.do(selected.Identify()); err != nil {
					panic(err)
				}
			}
		case "n":
			if selected := m.listModel.Selected(); selected != nil {
				if err := m.undo(selected.Identify()); err != nil {
					panic(err)
				}
			}
		case "d":
			if selected := m.listModel.Selected(); selected != nil {
				if err := m.del(selected.Identify()); err != nil {
					panic(err)
				}
			}
		}
		listModel, cmd := m.listModel.Update(msg)
		m.listModel = listModel.(SelectableListModel)
		return m, cmd

	case tea.WindowSizeMsg:
		_, v := m.style.GetFrameSize()
		msg.Height -= v + lipgloss.Height(m.title.String())
		listModel, cmd := m.listModel.Update(msg)
		m.listModel = listModel.(SelectableListModel)
		return m, cmd

	case itemsMsg:
		tags := apply.ToSlice(logic.ListTags(msg), types.IDToString)
		slices.Sort(tags)
		m.tabs = append([]string{allItems}, tags...)
		if m.tabs[m.activeTab] != allItems {
			msg = itemsMsg(logic.SublistWithTags([]types.ID{types.ID(m.tabs[m.activeTab])}, msg))
		}
		listModel, cmd := m.listModel.Update(msg)
		m.listModel = listModel.(SelectableListModel)
		return m, cmd
	}

	listModel, cmd := m.listModel.Update(msg)
	m.listModel = listModel.(SelectableListModel)
	return m, cmd
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m TabbedListModel) View() string {
	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		if i == m.activeTab {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	return m.style.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		m.title.String(),
		lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.JoinVertical(lipgloss.Left, renderedTabs...), m.listModel.View()),
	))
}

var (
	inactiveTabBorder = lipgloss.HiddenBorder()
	activeTabBorder   = lipgloss.RoundedBorder()
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder)
)

const (
	banner = `
  _____     ___     
 / ___/__  / _ \___ 
/ (_ / _ \/ // / _ \
\___/\___/____/\___/  your to-do list at...
  `
)
