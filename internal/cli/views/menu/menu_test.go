package menu

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knz/catwalk"
)

type listPageWrapper struct {
	*listPage
}

func (m *listPageWrapper) Init() tea.Cmd {
	return nil
}

func (m *listPageWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.listPage.Update(m, msg)
}

func TestAuth(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(func(i int) {})
	catwalk.RunModel(t, "menu", &listPageWrapper{m})
}
