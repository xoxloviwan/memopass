package pair

import (
	ctrl "iwakho/gopherkeep/internal/cli/controls"
	"iwakho/gopherkeep/internal/cli/views/basics/form"

	tea "github.com/charmbracelet/bubbletea"
)

type modelForm = form.ModelForm

func InitPair(nextPage func()) modelForm {
	fc := form.FormCaller{
		FormName:    "Вход",
		InputNames:  []string{"Логин", "Пароль"},
		ButtonNames: []string{"Добавить", "Отмена"},
	}
	m := form.InitForm(fc)
	m.NextPage = nextPage
	m.Control = new(ctrl.AddPairCtrl)
	return m
}

type addPairPage struct {
	Model  modelForm
	width  int
	height int
}

func NewPage(nextPage func()) *addPairPage {
	return &addPairPage{InitPair(nextPage), 0, 0}
}

func (pp *addPairPage) Init(width, height int) {
	pp.width = width
}
func (pp *addPairPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := pp.Model.Update(msg)
	pp.Model = model.(modelForm)
	return m, cmd
}

func (pp *addPairPage) View() string {
	return pp.Model.View()
}
