package show

import (
	"fmt"
	"iwakho/gopherkeep/internal/cli/views/basics/list"
	"iwakho/gopherkeep/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

type Control interface {
	GetCards(int, int) ([]model.CardInfo, error)
}

type CardFetcher struct {
	Control
}

func (f CardFetcher) Fetch(itemsPerPage int, offset int) []list.Item {
	items := []list.Item{}
	cards, err := f.GetCards(itemsPerPage, offset)
	if err != nil {
		return []list.Item{{Title: "Ошибка", Description: err.Error()}}
	}

	for _, v := range cards {
		item := list.Item{
			Title:       v.Date.Local().String(),
			Description: fmt.Sprintf("\tНомер: %s\n\tДействует до: %s\n\tCVV: %s", v.Number, v.Exp, v.VerifVal),
		}
		items = append(items, item)
	}
	return items
}

func NewPage(nextPage func(), ctrl Control) tea.Model {
	return list.New(
		"Посмотреть карты",
		&CardFetcher{Control: ctrl},
		nextPage)
}
