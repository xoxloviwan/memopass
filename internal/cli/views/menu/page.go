package menu

import (
	item "iwakho/gopherkeep/internal/cli/views/basics/list"

	"github.com/charmbracelet/bubbles/list"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

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

func NewPage(nextPage func(int) int) *modelList {
	l := listModel(mainItems())
	return &modelList{list: l, nextPage: nextPage}
}
