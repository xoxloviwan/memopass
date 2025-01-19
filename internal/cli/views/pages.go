package views

import (
	addCard "iwakho/gopherkeep/internal/cli/views/items/creditcard/add"
	showCards "iwakho/gopherkeep/internal/cli/views/items/creditcard/show"
	"iwakho/gopherkeep/internal/cli/views/items/file/picker"
	addPair "iwakho/gopherkeep/internal/cli/views/items/pair/add"
	showPairs "iwakho/gopherkeep/internal/cli/views/items/pair/show"
	addText "iwakho/gopherkeep/internal/cli/views/items/textarea/add"
	showTexts "iwakho/gopherkeep/internal/cli/views/items/textarea/show"
	"iwakho/gopherkeep/internal/cli/views/login"
	"iwakho/gopherkeep/internal/cli/views/menu"

	tea "github.com/charmbracelet/bubbletea"
)

// inspired by github.com/sspaink/telegraf-companion

var (
	ready bool
)

type Page interface {
	tea.Model
}

type Sender interface {
	Send(tea.Msg)
}

type Pages struct {
	pages       []Page
	currentPage int
	Sender
	width  int
	height int
}

type Controller interface {
	login.Control
	addPair.Control
	addCard.Control
	picker.Control
	showPairs.Control
	showCards.Control
	addText.Control
}

const offset = 2

func InitPages(ctrl Controller) *Pages {
	const pageTotal = 10
	p := Pages{
		pages: make([]Page, pageTotal),
	}

	p.add(0, login.NewPage(p.nextPage(1), ctrl))
	p.add(1, menu.NewPage(func(id int) {
		nextPage := id + offset
		if nextPage < pageTotal {
			p.currentPage = nextPage
		} else {
			p.currentPage = 1
		}
	}))
	p.add(offset+0, addPair.NewPage(p.nextPage(1), ctrl))
	p.add(offset+1, showPairs.NewPage(p.nextPage(1), ctrl))
	p.add(offset+2, addText.NewPage(p.nextPage(1), ctrl))
	p.add(offset+3, showTexts.NewPage(p.nextPage(1)))
	p.add(offset+4, picker.NewPage(p.nextPage(1), ctrl))
	p.add(offset+5, p.get(1)) // TODO file download
	p.add(offset+6, addCard.NewPage(p.nextPage(1), ctrl))
	p.add(offset+7, showCards.NewPage(p.nextPage(1), ctrl))

	return &p
}

func WithSender(pages *Pages) func(*tea.Program) {
	return func(pp *tea.Program) {
		pages.Sender = pp
	}
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
	return tea.Batch(tea.EnterAltScreen)
}

func (s *Pages) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		if !ready {
			ready = true
			p := s.pages[s.currentPage]
			if p != nil {
				p.Init()
			}
		}
	}
	lastPage := s.currentPage
	_, cmd := s.pages[s.currentPage].Update(msg)

	if s.currentPage != lastPage {
		cmd = tea.Batch(cmd, s.pages[s.currentPage].Init())
		if s.currentPage == offset+3 {
			_, cmd1 := s.pages[s.currentPage].Update(tea.WindowSizeMsg{Width: s.width, Height: s.height})
			cmd = tea.Batch(cmd, cmd1)
		}
	}
	return s, cmd
}

func (s *Pages) View() string {
	if !ready {
		return "\n  Initializing..."
	}
	return s.pages[s.currentPage].View()
}
