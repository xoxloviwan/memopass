package login

import tea "github.com/charmbracelet/bubbletea"

type AuthPage struct {
	modelForm
}

func NewAuthPage() (AuthPage, error) {
	return AuthPage{
		InitSignUp(),
	}, nil
}

func (ap *AuthPage) Init(width, height int) {
	ap.modelForm.Init()
}

func (ap *AuthPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}
func (ap AuthPage) View() string {
	return ap.modelForm.View()
}
