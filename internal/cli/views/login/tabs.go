package login

import (
	"iwakho/gopherkeep/internal/cli/views/basics/form"
)

func InitLogin(nextPage func(), client Client) modelForm {
	fc := form.FormCaller{
		FormName:    "Вход",
		InputNames:  []string{"Логин", "Пароль"},
		ButtonNames: []string{"Войти"},
	}
	m := form.InitForm(&fc, client.Login)
	m.NextPage = nextPage
	return *m
}

func InitSignUp(nextPage func(), client Client) modelForm {
	fc := form.FormCaller{
		FormName: "Регистрация",
		InputNames: []string{
			"Придумайте логин",
			"Введите пароль",
			"Повторите пароль",
		},
		ButtonNames: []string{"Зарегистрироваться"},
	}
	m := form.InitForm(&fc, client.SignUp)
	m.NextPage = nextPage
	return *m
}
