package pair

import (
	ctrl "iwakho/gopherkeep/internal/cli/controls"
	"iwakho/gopherkeep/internal/cli/views/form"
)

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
