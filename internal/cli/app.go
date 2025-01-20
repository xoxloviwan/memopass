package cli

import (
	ctrl "iwakho/gopherkeep/internal/cli/controls"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/cli/views"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	tui   *tea.Program
	pages tea.Model
}

func (app *App) Run() error {
	_, err := app.tui.Run()
	return err
}

func New(client *iHttp.Client) *App {
	pages := views.InitPages(ctrl.New(client))
	prg := tea.NewProgram(pages)
	return &App{
		tui:   prg,
		pages: pages,
	}
}
