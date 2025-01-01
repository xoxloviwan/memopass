package show

import (
	"fmt"
	"iwakho/gopherkeep/internal/cli/views/item"

	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

const listHeight = 14

func listModel(lst []list.Item) list.Model {
	const defaultWidth = 20

	l := list.New(lst, item.ItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Посмотреть пары логин/пароль"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

type pairPage struct {
	Model  modelList
	width  int
	height int
}

func pairsItems() []list.Item {
	items := []list.Item{}
	for i := range 10 {
		item := item.Item{
			Title:       fmt.Sprintf("aaa %d", i),
			Description: fmt.Sprintf("bbbbbb\n dfdfd \ndrf dfr %d", i),
		}
		items = append(items, item)
	}
	return items
}

func NewPage(nextPage func()) *pairPage {
	lm := listModel(pairsItems())
	return &pairPage{modelList{list: lm, nextPage: nextPage}, 0, 0}
}

func (pp *pairPage) Init(width, height int) {
	pp.width = width
}
func (pp *pairPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := pp.Model.Update(msg)
	pp.Model = model.(modelList)
	return m, cmd
}

func (pp *pairPage) View() string {
	return pp.Model.View()
}
