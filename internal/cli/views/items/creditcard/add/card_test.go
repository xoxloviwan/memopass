package creditcard

import (
	"iwakho/gopherkeep/internal/model"
	"testing"

	"github.com/knz/catwalk"
)

type MockController struct{}

func (c *MockController) AddCard(p model.Card) error {
	return nil
}

func TestCard(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(func() {}, &MockController{})
	catwalk.RunModel(t, "card", m)
}
