package list

import (
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

type ModelList struct {
	List        list.Model
	takenItemID int
	choose      choice
	quitting    bool
	NextPage    func(int)
}

func (m ModelList) Init() tea.Cmd {
	return nil
}

func (m ModelList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.List.SelectedItem().(Item)
			if ok {
				if m.List.Title == mainTitle {
					m.takenItemID = m.List.Index()
					m.List = ListModel(actionItems(i.Title))
				} else {
					switch i.Title {
					case AddItem:
						m.NextPage(m.takenItemID * 2)
					case ShowItem:
						m.NextPage(m.takenItemID*2 + 1)
					default:
						m.List = ListModel(mainItems())
					}
				}
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m ModelList) View() string {
	if m.quitting {
		return quitTextStyle.Render("Ничего не нужно? Ну пока!")
	}
	return m.List.View()
}
