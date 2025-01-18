package show

import (
	"iwakho/gopherkeep/internal/cli/views/basics/item"
	"iwakho/gopherkeep/internal/model"

	"github.com/charmbracelet/bubbles/list"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

const listHeight = 14

func listModel() list.Model {
	const defaultWidth = 100

	l := list.New([]list.Item{}, item.ItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Посмотреть пары логин/пароль"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

type Control interface {
	GetPairs(int, int) ([]model.PairInfo, error)
}

func NewPage(nextPage func(), ctrl Control) *modelList {
	return &modelList{
		list:     listModel(),
		nextPage: nextPage,
		Control:  ctrl,
	}
}
