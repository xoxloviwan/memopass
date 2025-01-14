package form

import (
	"errors"
	"fmt"
	"strings"

	btn "iwakho/gopherkeep/internal/cli/views/basics/button"
	"iwakho/gopherkeep/internal/model"

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

type submitFunc func(model.Pair) error

type ModelForm struct {
	Name        string
	focusIndex  int
	inputs      []textinput.Model
	cursorMode  cursor.Mode
	NextPage    func()
	buttons     []string
	indexMax    int
	failMessage string
	isUpdated   bool
	submitCall  submitFunc
}

func (m *ModelForm) Submit(p model.Pair) error {
	return m.submitCall(p)
}

type FormCaller struct {
	FormName    string
	InputNames  []string
	ButtonNames []string
}

func InitForm(fc *FormCaller, submitCall submitFunc) *ModelForm {
	m := ModelForm{
		Name:       fc.FormName,
		inputs:     make([]textinput.Model, len(fc.InputNames)),
		submitCall: submitCall,
	}
	m.buttons = fc.ButtonNames
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		if i == 0 {
			t.Placeholder = fc.InputNames[i]
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		} else {
			t.Placeholder = fc.InputNames[i]
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}
		m.inputs[i] = t
	}
	m.indexMax = len(m.inputs) + len(m.buttons) - 1
	return &m
}

func (m *ModelForm) Init() tea.Cmd {
	return textinput.Blink
}

func (m *ModelForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

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
			if s == "enter" && m.focusIndex == len(m.inputs) && m.isUpdated {
				m.isUpdated = false
				login := m.inputs[0].Value()
				password := m.inputs[1].Value()
				var err error
				if len(m.inputs) == 3 {
					passwordRepeated := m.inputs[2].Value()
					if password != passwordRepeated {
						err = errors.New("Пароли не совпадают")
					}
				}
				if err == nil {
					err = m.Submit(model.Pair{Login: login, Password: password})
				}
				if err != nil {
					m.failMessage = err.Error()
					return m, nil
				}
				if m.NextPage != nil {
					m.NextPage()
				}
				return m, nil
			}

			if s == "enter" && m.focusIndex == len(m.inputs)+1 { // back
				if m.NextPage != nil {
					m.NextPage()
				}
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

func (m *ModelForm) updateInputs(msg tea.Msg) tea.Cmd {
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

func (m *ModelForm) View() string {
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
