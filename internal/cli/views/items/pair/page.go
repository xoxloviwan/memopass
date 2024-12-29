package pair

import (
	tea "github.com/charmbracelet/bubbletea"
)

type pairPage struct {
	Model  modelPair
	width  int
	height int
}

func NewPairPage(nextPage func()) *pairPage {
	return &pairPage{InitPair(nextPage), 0, 0}
}

func (pp *pairPage) Init(width, height int) {
	pp.width = width
}
func (pp *pairPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := pp.Model.Update(msg)
	pp.Model = model.(modelPair)
	return m, cmd
}

func (pp *pairPage) View() string {
	return pp.Model.View()
}
