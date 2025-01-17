package views

import (
	addCard "iwakho/gopherkeep/internal/cli/views/items/creditcard/add"
	showCards "iwakho/gopherkeep/internal/cli/views/items/creditcard/show"
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

type Page interface {
	Init(int, int)
	Update(tea.Model, tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

type Sender interface {
	Send(tea.Msg)
}

type Pages struct {
	pages       []Page
	currentPage int
	Sender
}

type Controller interface {
	login.Control
	addPair.Control
	addCard.Control
	picker.Control
	showPairs.Control
	showCards.Control
}

func InitPages(ctrl Controller) *Pages {
	const pageTotal = 10
	p := Pages{
		pages: make([]Page, pageTotal),
	}
	const offset = 2

	p.add(0, login.NewPage(p.nextPage(1), ctrl))
	p.add(1, menu.NewPage(func(id int) {
		nextPage := id + offset
		if nextPage < pageTotal {
			p.currentPage = nextPage
		} else {
			p.currentPage = 1
		}
		// fix refresh for list and file picker
		if (id == 1 || id == 4) && p.Sender != nil {
			go p.Send(new(tea.Msg))
		}
		if id == 1 {
			p.add(offset+id, showPairs.NewPage(p.nextPage(1), ctrl)) // reset list
		}
	}))
	p.add(offset+0, addPair.NewPage(p.nextPage(1), ctrl))
	p.add(offset+1, showPairs.NewPage(p.nextPage(1), ctrl))
	p.add(offset+2, p.get(1)) // TODO text editor
	p.add(offset+3, p.get(1)) // TODO text viewer
	p.add(offset+4, picker.NewPage(p.nextPage(1), ctrl))
	p.add(offset+5, p.get(1)) // TODO file download
	p.add(offset+6, addCard.NewPage(p.nextPage(1), ctrl))
	p.add(offset+7, showCards.NewPage(p.nextPage(1), ctrl))

	return &p
}

func (ps *Pages) add(id int, page Page) {
	ps.pages[id] = page
}

func (ps *Pages) get(id int) Page {
	return ps.pages[id]
}

func (ps *Pages) nextPage(id int) func() {
	return func() {
		ps.currentPage = id
	}
}

func (*Pages) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.EnableMouseCellMotion)
}

func (s *Pages) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	model, cmd := s.pages[s.currentPage].Update(s, msg)
	return model, cmd
}

func (s *Pages) View() string {
	if !ready {
		return "\n  Initializing..."
	}
	return s.pages[s.currentPage].View()
}
