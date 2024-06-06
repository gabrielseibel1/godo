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

type editor func(types.Actionable) error

type selectableListModel interface {
	tea.Model
	Selected() types.Actionable
}

type editableTextModel interface {
	tea.Model
	Text() string
	Clear()
	Focus()
	Blur()
	Focused() bool
	Set(string)
}

type TabbedListModel struct {
	title       lipgloss.Style
	listModel   selectableListModel
	tabs        []string
	activeTab   int
	style       lipgloss.Style
	do          doer
	undo        undoer
	del         deleter
	edit        editor
	editorModel editableTextModel
}

func NewTabbedListModel(
	path string,
	list selectableListModel,
	editor editableTextModel,
	do doer,
	undo undoer,
	del deleter,
	edit editor,
) *TabbedListModel {
	return &TabbedListModel{
		title: lipgloss.NewStyle().
			// Bold(true).
			// Foreground(lipgloss.Color("10")).
			// Border(lipgloss.ThickBorder()).
			// BorderForeground(highlightColor).
			Padding(1, 1).
			SetString(fmt.Sprintf("%s\n%s", banner, path)),
		listModel:   list,
		editorModel: editor,
		activeTab:   0,
		style:       lipgloss.NewStyle().Margin(2, 2),
		do:          do,
		undo:        undo,
		del:         del,
		edit:        edit,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m TabbedListModel) Init() tea.Cmd {
	return m.listModel.Init()
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *TabbedListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editorModel.Focused() {

		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case tea.KeyCtrlS.String():
				if selected := m.listModel.Selected(); m.editorModel.Focused() && selected != nil {
					selected.Describe(m.editorModel.Text())
					if err := m.edit(selected); err != nil {
						panic(err)
					}
					m.editorModel.Clear()
					m.editorModel.Blur()
					return m, nil
				}
			}
		}

		editor, cmd := m.editorModel.Update(msg)
		m.editorModel = editor.(editableTextModel)
		return m, cmd
	}

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
				if err := m.do(selected.Identity()); err != nil {
					panic(err)
				}
			}
		case "n":
			if selected := m.listModel.Selected(); selected != nil {
				if err := m.undo(selected.Identity()); err != nil {
					panic(err)
				}
			}
		case "d":
			if selected := m.listModel.Selected(); selected != nil {
				if err := m.del(selected.Identity()); err != nil {
					panic(err)
				}
			}
		case "e":
			defer m.editorModel.Focus()
		}
		return m.updateSubModels(msg)

	case tea.WindowSizeMsg:
		// TODO probably worth reviewing if need to resize other comps
		_, v := m.style.GetFrameSize()
		msg.Height -= v + lipgloss.Height(m.title.String())
		msg.Height -= lipgloss.Height(m.editorModel.View())
		return m.updateSubModels(msg)

	case itemsMsg:
		tags := apply.ToSlice(logic.ListTags(msg), types.IDToString)
		slices.Sort(tags)
		m.tabs = append([]string{allItems}, tags...)
		if m.tabs[m.activeTab] != allItems {
			msg = itemsMsg(logic.SublistWithTags([]types.ID{types.ID(m.tabs[m.activeTab])}, msg))
		}
		return m.updateSubModels(msg)

	default:
		return m.updateSubModels(msg)
	}
}

func (m *TabbedListModel) updateSubModels(msg tea.Msg) (tea.Model, tea.Cmd) {
	listModel, cmd1 := m.listModel.Update(msg)
	m.listModel = listModel.(selectableListModel)
	editorModel, cmd2 := m.editorModel.Update(msg)
	m.editorModel = editorModel.(editableTextModel)
	return m, tea.Batch(cmd1, cmd2)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *TabbedListModel) View() string {
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

	if !m.editorModel.Focused() {
		if selected := m.listModel.Selected(); selected != nil {
			m.editorModel.Set(selected.Description())
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.title.String(),
		lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.JoinVertical(lipgloss.Left, renderedTabs...), m.listModel.View()),
		m.editorModel.View(),
	)
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
