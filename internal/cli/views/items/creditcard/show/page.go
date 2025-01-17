package show

import (
	"iwakho/gopherkeep/internal/cli/views/basics/item"
	"iwakho/gopherkeep/internal/model"

	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

const listHeight = 14

func listModel() list.Model {
	const defaultWidth = 20

	l := list.New([]list.Item{}, item.ItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Посмотреть карты"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

type Control interface {
	GetCards(int, int) ([]model.CardInfo, error)
}

type cardPage struct {
	Model  modelList
	width  int
	height int
}

func NewPage(nextPage func(), ctrl Control) *cardPage {
	return &cardPage{modelList{
		list:     listModel(),
		nextPage: nextPage,
		Control:  ctrl,
	}, 0, 0}
}

func (pp *cardPage) Init(width, height int) {
	pp.width = width
}
func (pp *cardPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := pp.Model.Update(msg)
	pp.Model = model.(modelList)
	return m, cmd
}

func (pp *cardPage) View() string {
	return pp.Model.View()
}
