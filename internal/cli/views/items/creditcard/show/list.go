package show

import (
	"fmt"

	"iwakho/gopherkeep/internal/cli/views/basics/item"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type modelList struct {
	list       list.Model
	ready      bool
	nextPage   func()
	allFetched bool
	Control
}

func (m modelList) getMoreItems() []list.Item {
	items := []list.Item{}
	itemsPerPage := m.list.Paginator.PerPage
	offset := len(m.list.Items())
	cards, err := m.GetCards(itemsPerPage, offset)
	if err != nil {
		return []list.Item{item.Item{Title: "Ошибка", Description: err.Error()}}
	}

	for _, v := range cards {
		item := item.Item{
			Title:       v.Date.Local().String(),
			Description: fmt.Sprintf("\tНомер: %s\n\tДействует до: %s\n\tCVV: %s", v.Number, v.Exp, v.VerifVal),
		}
		items = append(items, item)
	}
	return items
}

func (m modelList) Init() tea.Cmd {
	return m.list.SetItems(m.getMoreItems())
}

func (m modelList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.ready {
		cmd := m.list.SetItems(m.getMoreItems())
		m.ready = true
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.ready = false
			m.nextPage()
			return m, nil
		case "down":
			if m.list.Paginator.OnLastPage() && m.list.Index() == len(m.list.Items())-1 && !m.allFetched {
				items := m.list.Items()
				newItems := m.getMoreItems()
				if len(newItems) == 0 {
					m.allFetched = true
					return m, nil
				}
				items = append(items, newItems...)
				cmd := m.list.SetItems(items)
				var cmd2 tea.Cmd
				m.list, cmd2 = m.list.Update(msg)
				return m, tea.Batch(cmd, cmd2)
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelList) View() string {
	return m.list.View()
}
