package presentation

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type EditorModel struct {
	textArea     textarea.Model
	overrideText string
	err          error
}

type errMsg error

func NewEditorModel() *EditorModel {
	ti := textarea.New()

	return &EditorModel{
		textArea: ti,
		err:      nil,
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
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyEsc:
			if m.textArea.Focused() {
				m.textArea.Blur()
			}
			return m, nil
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textArea, cmd = m.textArea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *EditorModel) View() string {
	if m.Focused() {
		m.textArea.ShowLineNumbers = true
	} else {
		m.textArea.SetValue(m.overrideText)
		m.textArea.ShowLineNumbers = false
	}
	return fmt.Sprintf(
		"Item description:\n\n%s\n\n%s",
		m.textArea.View(),
		`("e" to start editing, Ctrl+S to save, Esc to cancel editing)`,
	) + "\n\n"
}

func (m *EditorModel) Text() string {
	return m.textArea.Value()
}

func (m *EditorModel) Clear() {
	m.textArea.Reset()
}

func (m *EditorModel) Set(text string) {
	m.overrideText = text
}

func (m *EditorModel) Focus() {
	m.textArea.Focus()
}

func (m *EditorModel) Blur() {
	m.textArea.Blur()
}

func (m *EditorModel) Focused() bool {
	return m.textArea.Focused()
}
