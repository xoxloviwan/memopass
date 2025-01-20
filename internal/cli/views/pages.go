package views

import (
	msgs "iwakho/gopherkeep/internal/cli/messages"
	addCard "iwakho/gopherkeep/internal/cli/views/items/creditcard/add"
	showCards "iwakho/gopherkeep/internal/cli/views/items/creditcard/show"
	listFiles "iwakho/gopherkeep/internal/cli/views/items/file/list"
	"iwakho/gopherkeep/internal/cli/views/items/file/picker"
	addFile "iwakho/gopherkeep/internal/cli/views/items/file/picker"
	getFile "iwakho/gopherkeep/internal/cli/views/items/file/saver"
	addText "iwakho/gopherkeep/internal/cli/views/items/note/add"
	listTexts "iwakho/gopherkeep/internal/cli/views/items/note/list"
	showText "iwakho/gopherkeep/internal/cli/views/items/note/show"
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
	listTexts.Control
	listFiles.Control
	showText.Control
	getFile.Control
}

const offset = 2

func InitPages(ctrl Controller) *Pages {
	const pageTotal = 12
	p := Pages{
		pages: make([]Page, pageTotal),
	}

	p.add(0, login.NewPage(1, ctrl))
	p.add(1, menu.NewPage(func(id int) int {
		nextPage := id + offset
		if nextPage < pageTotal {
			return nextPage
		} else {
			return 1
		}
	}))
	p.add(offset+0, addPair.NewPage(1, ctrl))
	p.add(offset+1, showPairs.NewPage(1, ctrl))
	p.add(offset+2, addText.NewPage(1, ctrl))
	p.add(offset+3, listTexts.NewPage(offset+8, ctrl))
	p.add(offset+4, addFile.NewPage(1, ctrl))
	p.add(offset+5, listFiles.NewPage(offset+9, ctrl))
	p.add(offset+6, addCard.NewPage(1, ctrl))
	p.add(offset+7, showCards.NewPage(1, ctrl))
	p.add(offset+8, showText.NewPage(1, ctrl))
	p.add(offset+9, getFile.NewPage(1, ctrl))

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

func (*Pages) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (s *Pages) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdChain []tea.Cmd
	p := s.pages[s.currentPage]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		if !ready {
			ready = true
			if p != nil {
				p.Init()
			}
		}
	case msgs.NextPage:
		s.currentPage = msg.PageID
		p = s.pages[s.currentPage]
		cmd := p.Init()
		cmdChain = append(cmdChain, cmd)
		_, cmd = p.Update(tea.WindowSizeMsg{Width: s.width, Height: s.height})
		cmdChain = append(cmdChain, cmd)
		if msg.Msg != nil {
			_, cmd := p.Update(msg.Msg)
			cmdChain = append(cmdChain, cmd)
		}
	}
	_, cmd := p.Update(msg)
	cmdChain = append(cmdChain, cmd)
	return s, tea.Batch(cmdChain...)
}

func (s *Pages) View() string {
	if !ready {
		return "\n  Initializing..."
	}
	return s.pages[s.currentPage].View()
}
