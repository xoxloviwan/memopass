package show

import (
	"fmt"
	"iwakho/gopherkeep/internal/cli/views/basics/list"
	"iwakho/gopherkeep/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

type Control interface {
	GetFiles(int, int) ([]model.FileInfo, error)
	GetFileById(int) error
}

type FileFetcher struct {
	Control
}

func (f FileFetcher) Fetch(itemsPerPage int, offset int) []list.Item {
	items := []list.Item{}
	files, err := f.GetFiles(itemsPerPage, offset)
	if err != nil {
		return []list.Item{{Title: "Ошибка", Description: err.Error()}}
	}

	for _, v := range files {
		item := list.Item{
			Title:       v.Date.Local().String(),
			Description: fmt.Sprintf("\tФайл: %s", v.Name),
		}
		items = append(items, item)
	}
	return items
}

func NewPage(nextPage int, ctrl Control) tea.Model {
	return list.New(
		"Посмотреть файлы",
		&FileFetcher{Control: ctrl},
		nextPage)
}
