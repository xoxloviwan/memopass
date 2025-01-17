package picker

import (
	tea "github.com/charmbracelet/bubbletea"
)

type PickerPage struct {
	Model  modelPicker
	width  int
	height int
}

func NewPage(nextPage func(), ctrl Control) *PickerPage {
	return &PickerPage{newModelPicker(nextPage, ctrl), 0, 0}
}

func (pp *PickerPage) Init(width, height int) {
	pp.width = width
}
func (pp *PickerPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := pp.Model.Update(msg)
	pp.Model = model.(modelPicker)
	return m, cmd
}

func (pp *PickerPage) View() string {
	return pp.Model.View()
}
