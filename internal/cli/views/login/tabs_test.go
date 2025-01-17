package login

import (
	"iwakho/gopherkeep/internal/model"
	"testing"

	"github.com/knz/catwalk"
)

type MockController struct{}

func (c *MockController) Login(p model.Pair) error {
	return nil
}
func (c *MockController) SignUp(p model.Pair) error {
	return nil
}

func TestLoginTab(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(func() {}, &MockController{})
	catwalk.RunModel(t, "login_tab", &m.TabContent[0])
}
