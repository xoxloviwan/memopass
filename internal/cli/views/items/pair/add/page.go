package pair

import (
	"iwakho/gopherkeep/internal/cli/views/basics/form"
	"iwakho/gopherkeep/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

type modelForm = form.ModelForm

type Client interface {
	AddPair(model.Pair) error
}

func InitPair(nextPage func(), client Client) modelForm {
	fc := form.FormCaller{
		FormName:    "Вход",
		InputNames:  []string{"Логин", "Пароль"},
		ButtonNames: []string{"Добавить", "Отмена"},
	}
	m := form.InitForm(&fc, client.AddPair)
	m.NextPage = nextPage
	return *m
}

type addPairPage struct {
	Model  modelForm
	width  int
	height int
}

func NewPage(nextPage func(), client Client) *addPairPage {
	return &addPairPage{InitPair(nextPage, client), 0, 0}
}

func (pp *addPairPage) Init(width, height int) {
	pp.width = width
}
func (pp *addPairPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := pp.Model.Update(msg)
	pp.Model = *model.(*modelForm)
	return m, cmd
}

func (pp *addPairPage) View() string {
	return pp.Model.View()
}
