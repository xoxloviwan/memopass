package notelist

import (
	"fmt"
	"iwakho/gopherkeep/internal/cli/views/basics/list"
	"iwakho/gopherkeep/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

type Control interface {
	GetTexts(int, int) ([]model.FileInfo, error)
	GetTextById(int) (string, error)
}

type TextFetcher struct {
	Control
}

func (f TextFetcher) Fetch(itemsPerPage int, offset int) []list.Item {
	items := []list.Item{}
	files, err := f.GetTexts(itemsPerPage, offset)
	if err != nil {
		return []list.Item{{Title: "Ошибка", Description: err.Error()}}
	}

	for _, v := range files {
		startText := string(v.Blob)
		if len(startText) > 10 {
			startText = startText[:5] + "..."
		}
		item := list.Item{
			Title:       v.Date.Local().String(),
			Description: fmt.Sprintf("\tЗаметка: %s", startText),
		}
		items = append(items, item)
	}
	return items
}

func NewPage(nextPage int, ctrl Control) tea.Model {
	return list.New(
		"Посмотреть заметки",
		&TextFetcher{Control: ctrl},
		nextPage)
}
