package login

import (
	"iwakho/gopherkeep/internal/model"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knz/catwalk"
)

type MockController struct{}

func (c *MockController) Login(p model.Pair) error {
	return nil
}
func (c *MockController) SignUp(p model.Pair) error {
	return nil
}

type AuthPageWrapper struct {
	*AuthPage
}

func (m *AuthPageWrapper) Init() tea.Cmd {
	return nil
}

func (m *AuthPageWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.AuthPage.Update(m, msg)
}

func TestAuth(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(func() {}, &MockController{})
	catwalk.RunModel(t, "login_tab", &AuthPageWrapper{m})
}
