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

type baseList struct {
	items []list.Item
	title string
}

const mainTitle = "Разделы"
const backAction = "Назад"

func mainItems() baseList {
	return baseList{
		items: []list.Item{
			item("Логины/пароли"),
			item("Заметки"),
			item("Файлы"),
			item("Данные банковских карт"),
		},
		title: mainTitle,
	}
}

func actionItems(title string) baseList {
	return baseList{
		items: []list.Item{
			item("Добавить"),
			item("Посмотреть"),
			item(backAction),
		},
		title: title,
	}
}

func listModel(lst baseList) list.Model {
	const defaultWidth = 20

	l := list.New(lst.items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = lst.title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

func NewListPage() listPage {
	l := listModel(mainItems())
	return listPage{modelList{list: l}, 0, 0}
}
