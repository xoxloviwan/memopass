package views

import (
	"iwakho/gopherkeep/internal/cli/views/items/pair"
	"iwakho/gopherkeep/internal/cli/views/list"
	"iwakho/gopherkeep/internal/cli/views/login"

	tea "github.com/charmbracelet/bubbletea"
)

// inspired by github.com/sspaink/telegraf-companion

var (
	ready bool
)

const pageTotal = 3

var currentPage = 0

type Page interface {
	Init(int, int)
	Update(tea.Model, tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

type App struct {
	pages []Page
}

func nextPage(id int) func() {
	return func() {
		currentPage = id
	}
}

func NewApp() (App, error) {
	app := App{pages: make([]Page, pageTotal)}

	app.pages[0] = login.NewAuthPage(nextPage(1))

	app.pages[1] = list.NewListPage(func(id int) {
		const offset = 2
		nextPage := id + offset
		if nextPage < pageTotal {
			currentPage = nextPage
		} else {
			currentPage = 1
		}
	})
	app.pages[2] = pair.NewPairPage(nextPage(1))
	//app.pages[3] = pair.NewPairPage(nextPage(1))

	return app, nil
}

func (s App) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.EnableMouseCellMotion)
}

func (s App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !ready {
			ready = true
			for _, p := range s.pages {
				p.Init(msg.Width, msg.Height)
			}
		}
	}

	model, cmd := s.pages[currentPage].Update(s, msg)
	return model, cmd
}

func (s App) View() string {
	if !ready {
		return "\n  Initializing..."
	}
	return s.pages[currentPage].View()
}
