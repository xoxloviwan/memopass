package menu

import (
	"testing"

	"github.com/knz/catwalk"
)

func TestMenu(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(func(i int) {})
	catwalk.RunModel(t, "menu", m)
}
