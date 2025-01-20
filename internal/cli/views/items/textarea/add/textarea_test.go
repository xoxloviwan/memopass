package textarea

import (
	"testing"

	"github.com/knz/catwalk"
)

type MockController struct{}

func (c *MockController) AddText(string) error {
	return nil
}

func TestTextArea(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(func() {}, &MockController{})
	catwalk.RunModel(t, "textarea", m)
}
