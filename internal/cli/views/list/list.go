package list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type choice struct {
	item   string
	action string
}

type modelList struct {
	list     list.Model
	choose   choice
	quitting bool
}

func (m modelList) Init() tea.Cmd {
	return nil
}

func (m modelList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				if m.list.Title == mainTitle {
					m.list = listModel(actionItems(string(i)))
				} else {
					if string(i) == backAction {
						m.list = listModel(mainItems())
					} else {
						m.choose = choice{
							item:   m.list.Title,
							action: string(i),
						}
					}
				}
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m modelList) View() string {

	if m.choose.action != "" {
		return quitTextStyle.Render(fmt.Sprintf("Выбор сделан! %s %s", m.choose.item, m.choose.action))
	}
	if m.quitting {
		return quitTextStyle.Render("Ничего не нужно? Ну пока!")
	}
	return m.list.View()
}
