package menu

import (
	"testing"

	"github.com/knz/catwalk"
)

func TestMenu(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(func(i int) int { return 0 })
	catwalk.RunModel(t, "menu", m)
}
