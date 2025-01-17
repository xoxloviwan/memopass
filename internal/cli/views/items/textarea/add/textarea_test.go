package textarea

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knz/catwalk"
)

type MockController struct{}

func (c *MockController) AddText(string) error {
	return nil
}

type TextAreaPageWrapper struct {
	*TextareaPage
}

func (m *TextAreaPageWrapper) Init() tea.Cmd {
	return nil
}

func (m *TextAreaPageWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.TextareaPage.Update(m, msg)
}

func TestTextArea(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(func() {}, &MockController{})
	catwalk.RunModel(t, "textarea", &TextAreaPageWrapper{m})
}
