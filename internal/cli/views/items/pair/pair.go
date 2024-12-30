package pair

import (
	"fmt"
	"strings"

	ctrl "iwakho/gopherkeep/internal/cli/controls"
	btn "iwakho/gopherkeep/internal/cli/views/button"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = btn.FocusedStyle
	cursorStyle         = btn.FocusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = btn.BlurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

type button string

type modelPair struct {
	focusIndex  int
	inputs      []textinput.Model
	cursorMode  cursor.Mode
	nextPage    func()
	buttons     []button
	indexMax    int
	failMessage string
}

const (
	submitButton = "Добавить"
	backButton   = "Назад"
)

func InitPair(nextPage func()) modelPair {

	m := modelPair{
		inputs:   make([]textinput.Model, 2),
		nextPage: nextPage,
		buttons:  make([]button, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Логин"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Пароль"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}
	m.buttons[0] = submitButton
	m.buttons[1] = backButton
	m.indexMax = len(m.inputs) + len(m.buttons) - 1

	return m
}

func (m modelPair) Init() tea.Cmd {
	return textinput.Blink
}

func (m modelPair) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			switch m.inputs[1].EchoMode {
			case textinput.EchoNormal:
				m.inputs[1].EchoMode = textinput.EchoPassword
			case textinput.EchoPassword:
				m.inputs[1].EchoMode = textinput.EchoNormal
			}
			return m, nil

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				err := ctrl.AddPair(m.inputs[0].Value(), m.inputs[1].Value())
				if err != nil {
					m.failMessage = err.Error()
				}
				return m, nil
			}

			if s == "enter" && m.focusIndex == len(m.inputs)+1 {
				m.nextPage()
				return m, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > m.indexMax {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = m.indexMax
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *modelPair) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m modelPair) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	for i := range m.buttons {
		b.WriteRune('\n')
		btn.RenderButton(&b, string(m.buttons[i]), m.focusIndex == i+len(m.inputs))
	}

	if m.failMessage != "" {
		b.WriteString(btn.ErrorStyle.Render(m.failMessage))
	}
	b.WriteRune('\n')

	b.WriteString(helpStyle.Render("echoMode is "))
	b.WriteString(cursorModeHelpStyle.Render(fmt.Sprintf("%d", m.inputs[1].EchoMode)))
	b.WriteString(helpStyle.Render(" (ctrl+r to change mode)"))

	return b.String()
}
