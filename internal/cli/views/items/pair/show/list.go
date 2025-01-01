package show

import (
	"fmt"
	ctrl "iwakho/gopherkeep/internal/cli/controls"
	"iwakho/gopherkeep/internal/cli/views/item"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func pairsItems() []list.Item {
	items := []list.Item{}
	show := ctrl.ShowCtrl{}
	pairs, err := show.GetPairs(10, 0)
	if err != nil {
		return []list.Item{item.Item{Title: "Ошибка", Description: err.Error()}}
	}

	for _, v := range pairs {
		item := item.Item{
			Title:       v.Date.Local().String(),
			Description: fmt.Sprintf("Логин: %s\nПароль: %s", v.Login, v.Password),
		}
		items = append(items, item)
	}
	return items
}

type modelList struct {
	list     list.Model
	ready    bool
	nextPage func()
}

func (m modelList) Init() tea.Cmd {
	return m.list.SetItems(pairsItems())
}

func (m modelList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.ready {
		cmd := m.list.SetItems(pairsItems())
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
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelList) View() string {
	return m.list.View()
}
