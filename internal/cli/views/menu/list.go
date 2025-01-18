package menu

import (
	"iwakho/gopherkeep/internal/cli/views/basics/item"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type modelList struct {
	list        list.Model
	takenItemID int
	quitting    bool
	nextPage    func(int)
}

func (m *modelList) Init() tea.Cmd {
	return nil
}

func (m *modelList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			i, ok := m.list.SelectedItem().(item.Item)
			if ok {
				if m.list.Title == mainTitle {
					m.takenItemID = m.list.Index()*2 + 1
					m.list = listModel(actionItems(i.Title))
				} else if m.takenItemID > 0 {
					switch i.Title {
					case AddItem:
						m.nextPage(m.takenItemID - 1)
					case ShowItem:
						m.nextPage(m.takenItemID)
					default:
						m.list = listModel(mainItems())
					}
				} else {
					m.list = listModel(mainItems())
				}
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *modelList) View() string {
	if m.quitting {
		return quitTextStyle.Render("Ничего не нужно? Ну пока!")
	}
	return m.list.View()
}
