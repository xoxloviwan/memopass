package menu

import (
	"iwakho/gopherkeep/internal/cli/views/basics/item"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
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
			item.Item{Title: Pairs},
			item.Item{Title: Notes},
			item.Item{Title: Files},
			item.Item{Title: Cards},
		},
		Title: mainTitle,
	}
}

func actionItems(title string) BaseList {
	return BaseList{
		Items: []list.Item{
			item.Item{Title: AddItem},
			item.Item{Title: ShowItem},
			item.Item{Title: backAction},
		},
		Title: title,
	}
}

func listModel(lst BaseList) list.Model {
	const defaultWidth = 20

	l := list.New(lst.Items, item.ItemDelegate{}, defaultWidth, listHeight)
	l.Title = lst.Title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

func NewPage(nextPage func(int)) *listPage {
	l := listModel(mainItems())
	return &listPage{modelList{list: l, nextPage: nextPage}, 0, 0}
}
