package picker

import (
	"testing"

	"github.com/knz/catwalk"
)

//go:generate go test . -args -rewrite

type MockController struct{}

func (c *MockController) AddFile(string) error {
	return nil
}

func NewPickerPageWrapper() *modelPicker {
	mp := newModelPicker(0, &MockController{})
	mp.testMode = true
	mp.filepicker.ShowPermissions = false
	mp.filepicker.ShowSize = false
	return mp
}

func TestPicker(t *testing.T) {
	// Initialize the model to test.
	pp := NewPickerPageWrapper()
	catwalk.RunModel(t, "picker", pp)
}
