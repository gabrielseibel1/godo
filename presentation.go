package main

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabrielseibel1/fungo/apply"
	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title + i.desc }

type itemsMsg []list.Item

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case itemsMsg:
		return m, m.list.SetItems(msg)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func actionableToItem(actionable types.Actionable) list.Item {
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
	return item{
		title: string(actionable.Identify()),
		desc: fmt.Sprintf("%s (%s) %s",
			checkbox, actionable.Worked().String(), actionable.Describe()),
	}
}

func uiList(repo data.Repository) {
	items := make([]list.Item, 0, 1000)
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "GoDo - ToDo List"

	valuesChannel := make(chan list.Item, 10_000)
	updateChannel := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := tea.NewProgram(m, tea.WithAltScreen())

	go func(vCh chan<- list.Item, uCh chan<- int) {
		ticker := time.NewTicker(time.Second).C
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				// get updated list from repo
				abstracts, err := repo.List()
				if err != nil {
					panic(err)
				}
				// don't send update if no changes occurred
				slices.SortFunc(abstracts, func(a1, a2 types.Actionable) int {
					i := strings.Compare(string(a1.Identify()), string(a2.Identify()))
					if i != 0 {
						return i
					}
					return strings.Compare(a1.Describe(), a2.Describe())
				})
				listItems := apply.ToSlice(abstracts, actionableToItem)
				if slices.Equal(listItems, m.list.Items()) {
					break
				}
				// send update
				for _, li := range listItems {
					valuesChannel <- li
				}
				updateChannel <- len(abstracts)
			}
		}
	}(valuesChannel, updateChannel)

	go func(vCh <-chan list.Item, uCh <-chan int) {
		for {
			select {
			case <-ctx.Done():
				return
			case n := <-uCh:
				// trigger item changes from update
				msg := make(itemsMsg, n)
				for i := 0; i < n; i++ {
					msg[i] = <-vCh
				}
				p.Send(msg)
			}
		}
	}(valuesChannel, updateChannel)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
