package noteshow

import (
	md "iwakho/gopherkeep/internal/model"
	"testing"

	"github.com/knz/catwalk"
)

type MockController struct{}

func (c *MockController) GetTextById(int) (*md.File, error) {
	return &md.File{Blob: []byte("многострочный текст\nмногострочный текст"), Name: "example.txt"}, nil
}

//go:generate go test . -args -rewrite

func TestPages(t *testing.T) {
	// Initialize the model to test.
	pager := NewPage(0, &MockController{})
	pager.content = "многострочный текст1\nмногострочный текст2"
	pager.title = "example2.txt"
	catwalk.RunModel(t, "pager", pager)
}
