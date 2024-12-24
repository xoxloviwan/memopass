package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listPage struct {
	Model  modelList
	width  int
	height int
}

func (lp *listPage) Init(width, height int) {
	lp.width = width
}
func (lp *listPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	return lp.Model.Update(msg)
}

func (lp *listPage) View() string {
	return lp.Model.View()
}

func NewListPage() listPage {
	items := []list.Item{
		item("Логины/пароли"),
		item("Заметки"),
		item("Файлы"),
		item("Данные банковских карт"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Разделы"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return listPage{modelList{list: l}, 0, 0}
}
