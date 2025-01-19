package menu

import (
	msgs "iwakho/gopherkeep/internal/cli/messages"
	iList "iwakho/gopherkeep/internal/cli/views/basics/list"

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
	nextPage    func(int) int
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
			i, ok := m.list.SelectedItem().(iList.Item)
			if ok {
				if m.list.Title == mainTitle {
					m.takenItemID = m.list.Index()*2 + 1
					m.list = listModel(actionItems(i.Title))
				} else if m.takenItemID > 0 {
					switch i.Title {
					case AddItem:
						return m, msgs.NextPageCmd(m.nextPage(m.takenItemID-1), nil)
					case ShowItem:
						return m, msgs.NextPageCmd(m.nextPage(m.takenItemID), nil)
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
