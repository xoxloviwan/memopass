package creditcard

import (
	"iwakho/gopherkeep/internal/model"
	"testing"

	"github.com/knz/catwalk"
)

//go:generate go test . -args -rewrite

type MockController struct{}

func (c *MockController) AddCard(p model.Card) error {
	return nil
}

func TestCard(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(0, &MockController{})
	catwalk.RunModel(t, "card", m)
}
