package show

import (
	"fmt"
	"iwakho/gopherkeep/internal/cli/views/basics/list"
	"iwakho/gopherkeep/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

type Control interface {
	GetPairs(int, int) ([]model.PairInfo, error)
}

type PairFetcher struct {
	Control
}

func (f PairFetcher) Fetch(itemsPerPage int, offset int) []list.Item {
	items := []list.Item{}
	pairs, err := f.GetPairs(itemsPerPage, offset)
	if err != nil {
		return []list.Item{{Title: "Ошибка", Description: err.Error()}}
	}

	for _, v := range pairs {
		item := list.Item{
			Title:       v.Date.Local().String(),
			Description: fmt.Sprintf("\tЛогин: %s\n\tПароль: %s", v.Login, v.Password),
		}
		items = append(items, item)
	}
	return items
}

func NewPage(nextPage int, ctrl Control) tea.Model {
	return list.New(
		"Посмотреть пары логин/пароль",
		&PairFetcher{Control: ctrl},
		nextPage)
}
