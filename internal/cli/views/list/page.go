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
	model, cmd := lp.Model.Update(msg)
	lp.Model = model.(modelList)
	return m, cmd
}

func (lp *listPage) View() string {
	return lp.Model.View()
}

type baseList struct {
	items []list.Item
	title string
}

const (
	mainTitle  = "Разделы"
	backAction = "Назад"

	Pairs = "Логины/пароли"
	Notes = "Заметки"
	Files = "Файлы"
	Cards = "Данные банковских карт"

	AddItem  = "Добавить"
	ShowItem = "Посмотреть"
)

func mainItems() baseList {
	return baseList{
		items: []list.Item{
			item(Pairs),
			item(Notes),
			item(Files),
			item(Cards),
		},
		title: mainTitle,
	}
}

func actionItems(title string) baseList {
	return baseList{
		items: []list.Item{
			item(AddItem),
			item(ShowItem),
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

func NewListPage(nextPage func(int)) *listPage {
	l := listModel(mainItems())
	return &listPage{modelList{list: l, nextPage: nextPage}, 0, 0}
}
