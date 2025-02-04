package creditcard

import (
	"errors"
	"fmt"
	"iwakho/gopherkeep/internal/model"

	msgs "iwakho/gopherkeep/internal/cli/messages"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	ccn = iota
	exp
	cvv
	ccnLen = 16 + 3
	expLen = 5
	cvvLen = 3
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type Control interface {
	AddCard(p model.Card) error
}

type modelCard struct {
	inputs   []textinput.Model
	focused  int
	err      error
	nextPage int
	Control
}

func newModelCard(nextPage int, ctrl Control) *modelCard {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[ccn] = textinput.New()
	inputs[ccn].Placeholder = "4505 **** **** 1234"
	inputs[ccn].Focus()
	inputs[ccn].CharLimit = 20
	inputs[ccn].Width = 30
	inputs[ccn].Prompt = ""
	inputs[ccn].Validate = ccnValidator

	inputs[exp] = textinput.New()
	inputs[exp].Placeholder = "ММ/ГГ "
	inputs[exp].CharLimit = 5
	inputs[exp].Width = 5
	inputs[exp].Prompt = ""
	inputs[exp].Validate = expValidator

	inputs[cvv] = textinput.New()
	inputs[cvv].Placeholder = "XXX"
	inputs[cvv].CharLimit = 3
	inputs[cvv].Width = 5
	inputs[cvv].Prompt = ""
	inputs[cvv].Validate = cvvValidator

	return &modelCard{
		inputs:   inputs,
		focused:  0,
		err:      errors.New("not filled"),
		nextPage: nextPage,
		Control:  ctrl,
	}
}

func NewPage(nextPage int, ctrl Control) *modelCard {
	return newModelCard(nextPage, ctrl)
}

func (m *modelCard) Init() tea.Cmd {
	return textinput.Blink
}

func (m *modelCard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if len(m.inputs[ccn].Value()) != ccnLen ||
				len(m.inputs[exp].Value()) != expLen ||
				len(m.inputs[cvv].Value()) != cvvLen {
				m.nextInput()
				break
			}
			err := m.AddCard(model.Card{
				Number:   m.inputs[ccn].Value(),
				Exp:      m.inputs[exp].Value(),
				VerifVal: m.inputs[cvv].Value(),
			})
			if err != nil {
				return m, nil
			}
			return m, msgs.NextPageCmd(m.nextPage, nil)
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		default:
			for i := range m.inputs {
				var candidate textinput.Model
				candidate, cmds[i] = m.inputs[i].Update(msg)
				if candidate.Err != nil && len(candidate.Value()) > 1 {
					m.err = candidate.Err
				} else {
					m.inputs[i] = candidate
				}
			}
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	}
	return m, tea.Batch(cmds...)
}

func (m *modelCard) View() string {
	return fmt.Sprintf(
		` Добавить новую карту:

 %s
 %s

 %s  %s
 %s        %s

 %s
`,
		inputStyle.Width(30).Render("Номер карты"),
		m.inputs[ccn].View(),
		inputStyle.Width(12).Render("Действует до"),
		inputStyle.Width(6).Render("CVV"),
		m.inputs[exp].View(),
		m.inputs[cvv].View(),
		continueStyle.Render("Продолжить ->"),
	) + "\n"
}

// nextInput focuses the next input field
func (m *modelCard) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *modelCard) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
