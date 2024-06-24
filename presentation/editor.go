package presentation

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EditorModel struct {
	title       textarea.Model
	description textarea.Model
	err         error
	state       editorModelState
}

type editorModelState int

const (
	editorModelStateNothing editorModelState = iota
	editorModelStateTitle
	editorModelStateDescription
)

type errMsg error

func NewEditorModel() *EditorModel {
	title := textarea.New()
	title.SetHeight(1)
	title.ShowLineNumbers = false
	description := textarea.New()
	description.ShowLineNumbers = false
	return &EditorModel{
		title:       title,
		description: description,
		err:         nil,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *EditorModel) Init() tea.Cmd {
	return textarea.Blink
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case editorModelStateNothing:
		return m, nil
	case editorModelStateTitle:
		return m.updateEditingTitle(msg)
	case editorModelStateDescription:
		return m.updatEditingDescription(msg)
	}
	return m, nil
}

func (m *EditorModel) updateEditingTitle(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.title.Focused() {
				m.title.Blur()
			}
			m.state = editorModelStateNothing
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	var cmd tea.Cmd
	m.title, cmd = m.title.Update(msg)
	return m, cmd
}

func (m *EditorModel) updatEditingDescription(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.description.Focused() {
				m.description.Blur()
			}
			m.state = editorModelStateNothing
			return m, nil
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	var cmd tea.Cmd
	m.description, cmd = m.description.Update(msg)
	return m, cmd
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *EditorModel) View() string {
	return lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, false).Padding(1, 2).Render(
		fmt.Sprintf("\n%s\n\n%s\n\n%s", m.title.View(), m.description.View(), `("r" to rename, "e" to edit, Ctrl+S to save, Esc to cancel)`),
	)
}

func (m *EditorModel) Description() string {
	return m.description.Value()
}

func (m *EditorModel) Title() string {
	return m.title.Value()
}

func (m *EditorModel) Clear() {
	m.title.Reset()
	m.description.Reset()
}

func (m *EditorModel) SetDescription(text string) {
	m.description.SetValue(text)
}

func (m *EditorModel) SetTitle(title string) {
	m.title.SetValue(title)
}

func (m *EditorModel) FocusTitle() {
	m.title.Focus()
	m.state = editorModelStateTitle
}

func (m *EditorModel) FocusDescription() {
	m.description.Focus()
	m.state = editorModelStateDescription
}

func (m *EditorModel) Blur() {
	m.title.Blur()
	m.description.Blur()
	m.state = editorModelStateNothing
}

func (m EditorModel) FocusedTitle() bool {
	return m.state == editorModelStateTitle
}

func (m EditorModel) FocusedDescription() bool {
	return m.state == editorModelStateDescription
}
