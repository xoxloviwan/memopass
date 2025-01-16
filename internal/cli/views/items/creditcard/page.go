package creditcard

import (
	tea "github.com/charmbracelet/bubbletea"
)

type CardPage struct {
	Model  modelCard
	width  int
	height int
}

func NewPage(nextPage func(), ctrl Control) *CardPage {
	return &CardPage{newModelCard(nextPage, ctrl), 0, 0}
}

func (pp *CardPage) Init(width, height int) {
	pp.width = width
}
func (pp *CardPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := pp.Model.Update(msg)
	pp.Model = model.(modelCard)
	return m, cmd
}

func (pp *CardPage) View() string {
	return pp.Model.View()
}
