package noteadd

import (
	"fmt"

	msgs "iwakho/gopherkeep/internal/cli/messages"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type textAreaModel struct {
	textarea textarea.Model
	err      error
	nextPage int
	Control
}

type Control interface {
	AddText(string) error
}

func newTextareaModel(nextPage int, ctrl Control) *textAreaModel {
	ti := textarea.New()
	ti.Placeholder = "Однажды это случилось ..."
	ti.Focus()

	return &textAreaModel{
		textarea: ti,
		err:      nil,
		nextPage: nextPage,
		Control:  ctrl,
	}
}

func NewPage(nextPage int, ctrl Control) *textAreaModel {
	return newTextareaModel(nextPage, ctrl)
}

func (m *textAreaModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *textAreaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlEnd:
			err := m.AddText(m.textarea.Value())
			if err != nil {
				m.err = err
				return m, nil
			}
			return m, msgs.NextPageCmd(m.nextPage, nil)
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *textAreaModel) View() string {
	return fmt.Sprintf(
		"Введите текст:\n\n%s\n\n%s",
		m.textarea.View(),
		"Для завершения ввода нажмите Ctrl+End",
	) + "\n\n"
}
