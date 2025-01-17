package picker

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knz/catwalk"
)

type MockController struct{}

func (c *MockController) AddFile(string) error {
	return nil
}

type PickerPageWrapper struct {
	*PickerPage
}

func NewPickerPageWrapper() *PickerPageWrapper {
	mp := newModelPicker(func() {}, &MockController{})
	mp.testMode = true
	pp := &PickerPage{mp, 0, 0}
	return &PickerPageWrapper{pp}
}

func (m *PickerPageWrapper) Init() tea.Cmd {
	return nil
}

func (m *PickerPageWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.PickerPage.Update(m, msg)
}

func TestPicker(t *testing.T) {
	// Initialize the model to test.
	pp := NewPickerPageWrapper()
	catwalk.RunModel(t, "picker", pp)
}
