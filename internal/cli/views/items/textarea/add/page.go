package textarea

import (
	tea "github.com/charmbracelet/bubbletea"
)

type TextareaPage struct {
	Model  textAreaModel
	width  int
	height int
}

func NewPage(nextPage func(), ctrl Control) *TextareaPage {
	return &TextareaPage{NewTextareaModel(nextPage, ctrl), 0, 0}
}

func (pp *TextareaPage) Init(width, height int) {
	pp.width = width
}
func (pp *TextareaPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := pp.Model.Update(msg)
	pp.Model = model.(textAreaModel)
	return m, cmd
}

func (pp *TextareaPage) View() string {
	return pp.Model.View()
}
