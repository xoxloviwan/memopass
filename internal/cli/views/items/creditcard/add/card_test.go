package creditcard

import (
	"iwakho/gopherkeep/internal/model"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knz/catwalk"
)

type MockController struct{}

func (c *MockController) AddCard(p model.Card) error {
	return nil
}

type CardPageWrapper struct {
	*CardPage
}

func (m *CardPageWrapper) Init() tea.Cmd {
	return nil
}

func (m *CardPageWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.CardPage.Update(m, msg)
}

func TestCard(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(func() {}, &MockController{})
	catwalk.RunModel(t, "card", &CardPageWrapper{m})
}
