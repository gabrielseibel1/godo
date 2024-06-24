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
	Selected() (int, types.Actionable)
	Select(types.ID) error
	SelectIndex(int) error
}

type editableTextModel interface {
	tea.Model
	Title() string
	SetTitle(string)
	Description() string
	SetDescription(string)
	Clear()
	FocusTitle()
	FocusDescription()
	Blur()
	FocusedTitle() bool
	FocusedDescription() bool
}

type tabbedListModelState int

const (
	tabbedListModelStateBrowsing tabbedListModelState = iota
	tabbedListModelStateEditTitle
	tabbedListModelStateEditDescription
)

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
	state       tabbedListModelState
}

func NewTabbedListModel(
	path string,
	version string,
	list selectableListModel,
	editor editableTextModel,
	do doer,
	undo undoer,
	del deleter,
	edit editor,
) *TabbedListModel {
	return &TabbedListModel{
		title:       lipgloss.NewStyle().Padding(1, 1).SetString(fmt.Sprintf("GoDo (%s) your todo list at...\n%s", version, path)),
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

func (m *TabbedListModel) updateEditingTitle(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case tea.KeyCtrlS.String():
			if _, selected := m.listModel.Selected(); selected != nil {
				if err := m.del(selected.Identity()); err != nil {
					panic(err)
				}
				selected.Identify(types.ID(m.editorModel.Title()))
				if err := m.edit(selected); err != nil {
					panic(err)
				}
				if err := m.listModel.Select(selected.Identity()); err != nil {
					panic(err)
				}
				m.editorModel.Blur()
				m.state = tabbedListModelStateBrowsing
			}
		case tea.KeyEsc.String():
			m.editorModel.Blur()
			m.state = tabbedListModelStateBrowsing
		}
	}

	editor, cmd := m.editorModel.Update(msg)
	m.editorModel = editor.(editableTextModel)
	return m, cmd
}

func (m *TabbedListModel) updateEditingDescription(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case tea.KeyCtrlS.String():
			if _, selected := m.listModel.Selected(); selected != nil {
				selected.Describe(m.editorModel.Description())
				if err := m.edit(selected); err != nil {
					panic(err)
				}
				m.editorModel.Blur()
				m.state = tabbedListModelStateBrowsing
			}
		case tea.KeyEsc.String():
			m.editorModel.Blur()
			m.state = tabbedListModelStateBrowsing
		}
	}
	editor, cmd := m.editorModel.Update(msg)
	m.editorModel = editor.(editableTextModel)
	return m, cmd
}

func (m *TabbedListModel) updateBrowsing(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if _, selected := m.listModel.Selected(); selected != nil {
				if err := m.do(selected.Identity()); err != nil {
					panic(err)
				}
			}
		case "n":
			if _, selected := m.listModel.Selected(); selected != nil {
				if err := m.undo(selected.Identity()); err != nil {
					panic(err)
				}
			}
		case "d":
			if _, selected := m.listModel.Selected(); selected != nil {
				if err := m.del(selected.Identity()); err != nil {
					panic(err)
				}
				// m.listModel.SelectIndex(max(0, i-1))
			}
		case "c":
			a := types.NewActivity("New Item", "No description yet...")
			if err := m.edit(a); err != nil {
				panic(err)
			}
		case "e":
			m.state = tabbedListModelStateEditDescription
			defer m.editorModel.FocusDescription()
		case "r":
			m.state = tabbedListModelStateEditTitle
			defer m.editorModel.FocusTitle()

		}
	}
	return m.updateSubModels(msg)
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *TabbedListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		_, v := m.style.GetFrameSize()
		msg.Height -= v + lipgloss.Height(m.title.String())
		msg.Height -= lipgloss.Height(m.editorModel.View())
		return m.updateSubModels(msg)
	}

	if msg, ok := msg.(itemsMsg); ok {
		tags := apply.ToSlice(logic.ListTags(msg), types.IDToString)
		slices.Sort(tags)
		m.tabs = append([]string{allItems}, tags...)
		if m.tabs[m.activeTab] != allItems {
			msg = itemsMsg(logic.SublistWithTags([]types.ID{types.ID(m.tabs[m.activeTab])}, msg))
		}
		return m.updateSubModels(msg)
	}

	switch m.state {
	case tabbedListModelStateEditDescription:
		return m.updateEditingDescription(msg)

	case tabbedListModelStateEditTitle:
		return m.updateEditingTitle(msg)

	case tabbedListModelStateBrowsing:
		return m.updateBrowsing(msg)
	}

	return m, nil
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

	if !m.editorModel.FocusedTitle() && !m.editorModel.FocusedDescription() {
		if _, selected := m.listModel.Selected(); selected != nil {
			m.editorModel.SetTitle(string(selected.Identity()))
			m.editorModel.SetDescription(selected.Description())
		}
	}

	// TODO turn tabs into paged list
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.title.String(),
		lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...), m.listModel.View()),
		m.editorModel.View(),
	)
}

var (
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle = lipgloss.NewStyle().PaddingRight(4).PaddingLeft(4)
	activeTabStyle   = inactiveTabStyle.Copy().Underline(true)
)
