package picker

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	ctrl "iwakho/gopherkeep/internal/cli/controls"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type fileSender interface {
	Submit(string) error
}

type modelPicker struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
	nextPage     func()
	is1stUpdate  bool
	renderCnt    int
	updateCnt    int
	updates      []tea.Msg
	stack        []byte
	control      fileSender
}

func newModelPicker(nextPage func()) modelPicker {
	fp := filepicker.New()
	fp.ShowHidden = true
	fp.AutoHeight = true
	fp.Height = 10
	fp.CurrentDirectory, _ = os.Getwd()
	return modelPicker{
		filepicker:  fp,
		nextPage:    nextPage,
		is1stUpdate: true,
		control:     new(ctrl.AddFileCtrl),
	}
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

type ReadDirMsg struct {
	id      int
	entries []os.DirEntry
}

func (m modelPicker) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m modelPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.is1stUpdate {
		m.is1stUpdate = false
		return m, m.Init()
	}
	m.updateCnt++
	m.updates = append(m.updates, msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)
	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	if m.selectedFile != "" {
		err := m.control.Submit(m.selectedFile)
		if err != nil {
			m.err = err
			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}
		m.selectedFile = ""
		m.nextPage()
	}

	return m, cmd
}

func (m modelPicker) View() string {
	if m.quitting {
		return ""
	}
	m.renderCnt++
	var s strings.Builder
	s.WriteString("Выберите файл:\n  ")
	// s.WriteString("Stack:" + string(m.stack))
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString("" + m.filepicker.CurrentDirectory)
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}

	view := m.filepicker.View()
	s.WriteString("\n\n" + view + "\n")

	s.WriteString("renders: " + strconv.Itoa(m.renderCnt) + " updates: " + strconv.Itoa(m.updateCnt) + " []updates:" + fmt.Sprintf("%v", m.updates)) // debug
	return s.String()
}
