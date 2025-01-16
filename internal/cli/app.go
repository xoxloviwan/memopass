package cli

import (
	ctrl "iwakho/gopherkeep/internal/cli/controls"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/cli/views"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	tui   *tea.Program
	pages *views.Pages
}

func (app *App) Run() error {
	_, err := app.tui.Run()
	return err
}

func New(client *iHttp.Client) *App {
	m := views.InitPages(ctrl.New(client))
	prg := tea.NewProgram(m, func(pp *tea.Program) {
		m.Sender = pp
	})
	return &App{
		tui:   prg,
		pages: m,
	}
}
