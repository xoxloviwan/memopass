package show

import (
	"fmt"
	"iwakho/gopherkeep/internal/cli/views/list"

	bubbleList "github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
)

type modelList = list.ModelList
type baseList = list.BaseList

type pairPage struct {
	Model  modelList
	width  int
	height int
}

func pairsItems() baseList {
	items := []bubbleList.Item{}
	for i := range 10 {
		item := list.Item{
			Title:       fmt.Sprintf("aaa %d", i),
			Description: fmt.Sprintf("bbbbbb\n dfdfd \ndrf dfr %d", i),
		}
		items = append(items, item)
	}
	bl := baseList{Items: items}
	return bl
}

func NewPairPage(nextPage func()) *pairPage {
	NextPage := func(index int) {
		nextPage()
	}
	lm := list.ListModel(pairsItems())
	return &pairPage{modelList{List: lm, NextPage: NextPage}, 0, 0}
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
