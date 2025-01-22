package views

import (
	ctrl "iwakho/gopherkeep/internal/cli/controls"
	"testing"

	"github.com/knz/catwalk"
)

//go:generate go test . -args -rewrite

func TestPages(t *testing.T) {
	// Initialize the model to test.
	pages := InitPages(ctrl.New(nil))
	catwalk.RunModel(t, "pages", pages)
}
