package pair

import (
	"iwakho/gopherkeep/internal/model"
	"testing"

	"github.com/knz/catwalk"
)

type MockController struct{}

func (c *MockController) AddPair(p model.Pair) error {
	return nil
}

//go:generate go test . -args -rewrite

func TestPages(t *testing.T) {
	// Initialize the model to test.
	add := NewPage(0, &MockController{})
	catwalk.RunModel(t, "pair", add)
}
