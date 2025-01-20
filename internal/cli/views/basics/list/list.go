package list

import (
	msgs "iwakho/gopherkeep/internal/cli/messages"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type Fetcher interface {
	Fetch(int, int) []Item
}

type modelList struct {
	list       list.Model
	allFetched bool
	Fetcher
	nextPage int
	sendID   bool
}

func New(title string, f Fetcher, nextPage int, sendID bool) *modelList {
	return &modelList{
		list:     newModel(title),
		nextPage: nextPage,
		Fetcher:  f,
		sendID:   sendID,
	}
}

const listHeight = 14

func newModel(title string) list.Model {
	const defaultWidth = 100

	l := list.New([]list.Item{}, ItemDelegate{}, defaultWidth, listHeight)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

func (m *modelList) fetchItems() []list.Item {
	itemsPerPage := m.list.Paginator.PerPage
	offset := len(m.list.Items())
	gotItems := m.Fetch(itemsPerPage, offset)
	var items []list.Item
	for _, v := range gotItems {
		items = append(items, v)
	}
	return items
}

func (m *modelList) Init() tea.Cmd {
	return m.list.SetItems(m.fetchItems())

}

func (m *modelList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.sendID {
				item := m.list.SelectedItem()
				if item == nil {
					return m, msgs.NextPageCmd(m.nextPage, nil)
				}
				it := item.(Item)
				if it.ID != 0 {
					return m, msgs.NextPageCmd(m.nextPage, msgs.LoadData{ID: it.ID})
				} else {
					return m, msgs.NextPageCmd(m.nextPage, nil)
				}
			}
			return m, nil
			// return m, msgs.NextPageCmd(m.nextPage, nil)
		case "down":
			if m.list.Paginator.OnLastPage() && m.list.Index() == len(m.list.Items())-1 && !m.allFetched {
				items := m.list.Items()
				newItems := m.fetchItems()
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

func (m *modelList) View() string {
	return m.list.View()
}
