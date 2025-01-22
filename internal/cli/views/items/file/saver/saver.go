package saver

import (
	msgs "iwakho/gopherkeep/internal/cli/messages"
	md "iwakho/gopherkeep/internal/model"
	"os"
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Control interface {
	GetFileById(id int) (*md.File, error)
}

type model struct {
	path      string
	savedPath string
	err       string
	Control
	nextPage int
}

func NewPage(nextPage int, ctrl Control) *model {
	return &model{
		path:      os.TempDir(),
		savedPath: "",
		nextPage:  nextPage,
		Control:   ctrl,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgs.LoadData:
		file, err := m.GetFileById(msg.ID)
		if err != nil {
			m.err = err.Error()
			return m, nil
		}
		m.savedPath = path.Join(m.path, file.Name)
		err = os.WriteFile(m.savedPath, file.Blob, 0644)
		if err != nil {
			m.err = err.Error()
			return m, nil
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m, msgs.NextPageCmd(m.nextPage, nil)
		}
	}
	return m, nil
}

func (m *model) View() string {
	var s strings.Builder
	if m.err != "" {
		s.WriteString("Ошибка: " + m.err)
	} else if m.savedPath == "" {
		s.WriteString("Загрузка...")
	} else {
		s.WriteString("Сохранено: " + m.savedPath)
	}
	return s.String()
}
