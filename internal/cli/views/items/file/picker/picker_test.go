package picker

import (
	"testing"

	"github.com/knz/catwalk"
)

type MockController struct{}

func (c *MockController) AddFile(string) error {
	return nil
}

func NewPickerPageWrapper() *modelPicker {
	mp := newModelPicker(func() {}, &MockController{})
	mp.testMode = true
	return mp
}

func TestPicker(t *testing.T) {
	// Initialize the model to test.
	pp := NewPickerPageWrapper()
	catwalk.RunModel(t, "picker", pp)
}
