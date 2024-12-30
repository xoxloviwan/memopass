package login

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"errors"
	"strings"

	ctrls "iwakho/gopherkeep/internal/cli/controls"
	btn "iwakho/gopherkeep/internal/cli/views/button"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = btn.FocusedStyle
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

type button struct {
	title   string
	focused bool
}

type Callback func()

type modelForm struct {
	name         string
	focusIndex   int
	inputs       []textinput.Model
	cursorMode   cursor.Mode
	submitButton button
	failMessage  string
	isLogin      bool
	onEnter      Callback
	isUpdated    bool
}

func InitLogin(onEnter Callback) modelForm {
	m := modelForm{
		name:   "Вход",
		inputs: make([]textinput.Model, 2),
		submitButton: button{
			title: "Войти",
		},
		isLogin: true,
		onEnter: onEnter,
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

	return m
}

func InitSignUp(onEnter Callback) modelForm {
	m := modelForm{
		name:   "Регистрация",
		inputs: make([]textinput.Model, 3),
		submitButton: button{
			title: "Зарегистрироваться",
		},
		onEnter: onEnter,
	}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Придумайте логин"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Введите пароль"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 2:
			t.Placeholder = "Повторите пароль"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m modelForm) Init() tea.Cmd {
	return textinput.Blink
}

func (m modelForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) && m.isUpdated {
				m.isUpdated = false
				login := m.inputs[0].Value()
				password := m.inputs[1].Value()
				var err error
				if m.isLogin {
					err = ctrls.TryLogin(login, password)
				} else {
					passwordRepeated := m.inputs[2].Value()
					if password == passwordRepeated {
						err = ctrls.SignUp(login, password)
					} else {
						err = errors.New("Пароли не совпадают")
					}
				}
				if err != nil {
					m.failMessage = err.Error()
					return m, nil
				}
				if m.onEnter != nil {
					m.onEnter()
				}
				return m, nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
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

func (m *modelForm) updateInputs(msg tea.Msg) tea.Cmd {
	if keymsg, ok := msg.(tea.KeyMsg); ok && keymsg.Type == tea.KeyRunes {
		m.failMessage = ""
		m.isUpdated = true
	}
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m modelForm) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	b.WriteString("\n\n")
	btn.RenderButton(&b, m.submitButton.title, m.focusIndex == len(m.inputs))

	if m.failMessage != "" {
		b.WriteString(btn.ErrorStyle.Render(m.failMessage))
	}
	if m.isLogin {
		b.WriteString("\n")
	}
	return b.String()
}
