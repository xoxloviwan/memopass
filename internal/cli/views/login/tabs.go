package login

import (
	ctrl "iwakho/gopherkeep/internal/cli/controls"
	"iwakho/gopherkeep/internal/cli/views/form"
)

func InitLogin(nextPage func()) modelForm {
	fc := form.FormCaller{
		FormName:    "Вход",
		InputNames:  []string{"Логин", "Пароль"},
		ButtonNames: []string{"Войти"},
	}
	m := form.InitForm(fc)
	m.NextPage = nextPage
	m.Control = new(ctrl.LoginCtrl)
	return m
}

func InitSignUp(nextPage func()) modelForm {
	fc := form.FormCaller{
		FormName: "Регистрация",
		InputNames: []string{
			"Придумайте логин",
			"Введите пароль",
			"Повторите пароль",
		},
		ButtonNames: []string{"Зарегистрироваться"},
	}
	m := form.InitForm(fc)
	m.NextPage = nextPage
	m.Control = new(ctrl.SignUpCtrl)
	return m
}
