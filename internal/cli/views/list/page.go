package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listPage struct {
	Model  ModelList
	width  int
	height int
}

func (lp *listPage) Init(width, height int) {
	lp.width = width
}
func (lp *listPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := lp.Model.Update(msg)
	lp.Model = model.(ModelList)
	return m, cmd
}

func (lp *listPage) View() string {
	return lp.Model.View()
}

type BaseList struct {
	Items []list.Item
	Title string
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

func mainItems() BaseList {
	return BaseList{
		Items: []list.Item{
			Item{Title: Pairs},
			Item{Title: Notes},
			Item{Title: Files},
			Item{Title: Cards},
		},
		Title: mainTitle,
	}
}

func actionItems(title string) BaseList {
	return BaseList{
		Items: []list.Item{
			Item{Title: AddItem},
			Item{Title: ShowItem},
			Item{Title: backAction},
		},
		Title: title,
	}
}

func ListModel(lst BaseList) list.Model {
	const defaultWidth = 20

	l := list.New(lst.Items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = lst.Title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

func NewListPage(nextPage func(int)) *listPage {
	l := ListModel(mainItems())
	return &listPage{ModelList{List: l, NextPage: nextPage}, 0, 0}
}
