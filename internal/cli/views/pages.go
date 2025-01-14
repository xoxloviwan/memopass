package views

import (
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/cli/views/items/creditcard"
	"iwakho/gopherkeep/internal/cli/views/items/file/picker"
	addPair "iwakho/gopherkeep/internal/cli/views/items/pair/add"
	showPairs "iwakho/gopherkeep/internal/cli/views/items/pair/show"
	"iwakho/gopherkeep/internal/cli/views/login"
	"iwakho/gopherkeep/internal/cli/views/menu"

	tea "github.com/charmbracelet/bubbletea"
)

// inspired by github.com/sspaink/telegraf-companion

var (
	ready bool
)

const pageTotal = 10

var currentPage = 0

type Page interface {
	Init(int, int)
	Update(tea.Model, tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

type Sender interface {
	Send(tea.Msg)
}

type App struct {
	pages []Page
	Sender
	*iHttp.Client
}

func nextPage(id int) func() {
	return func() {
		currentPage = id
	}
}

func NewApp(client *iHttp.Client) (*App, error) {
	app := App{
		pages:  make([]Page, pageTotal),
		Client: client,
	}

	const offset = 2
	app.pages[0] = login.NewPage(nextPage(1), app.Client)

	app.pages[1] = menu.NewPage(func(id int) {
		nextPage := id + offset
		if nextPage < pageTotal {
			currentPage = nextPage
		} else {
			currentPage = 1
		}
		// fix refresh for list and file picker
		if (id == 1 || id == 4) && app.Sender != nil {
			go app.Sender.Send(new(tea.Msg))
		}
	})
	app.pages[offset+0] = addPair.NewPage(nextPage(1), app.Client)
	app.pages[offset+1] = showPairs.NewPage(nextPage(1), app.Client)
	app.pages[offset+2] = app.pages[1] // TODO text editor
	app.pages[offset+3] = app.pages[1] // TODO text viewer
	app.pages[offset+4] = picker.NewPage(nextPage(1), app.Client)
	app.pages[offset+5] = app.pages[1] // TODO file download
	app.pages[offset+6] = creditcard.NewPage(nextPage(1), app.Client)
	app.pages[offset+7] = app.pages[1] // TODO credit card list view

	return &app, nil
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
				if p != nil {
					p.Init(msg.Width, msg.Height)
				}
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
