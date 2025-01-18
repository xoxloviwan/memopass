package pair

import (
	"iwakho/gopherkeep/internal/cli/views/basics/form"
	"iwakho/gopherkeep/internal/model"
)

type modelForm = form.ModelForm

type Control interface {
	AddPair(model.Pair) error
}

func InitPair(nextPage func(), ctrl Control) *modelForm {
	fc := form.FormCaller{
		FormName:    "Вход",
		InputNames:  []string{"Логин", "Пароль"},
		ButtonNames: []string{"Добавить", "Отмена"},
	}
	m := form.InitForm(&fc, ctrl.AddPair)
	m.NextPage = nextPage
	return m
}

func NewPage(nextPage func(), ctrl Control) *modelForm {
	return InitPair(nextPage, ctrl)
}
