package show

import (
	"fmt"
	"iwakho/gopherkeep/internal/model"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knz/catwalk"
)

type MockController struct {
	cnt int
}

func (c *MockController) GetPairs(limit int, offset int) ([]model.PairInfo, error) {
	c.cnt++
	fmt.Println(c.cnt, offset)
	if offset == 2 {
		return []model.PairInfo{
			{
				Pair: model.Pair{
					Login:    "yyy",
					Password: "zzz",
				},
				Metainfo: model.Metainfo{
					Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			{
				Pair: model.Pair{
					Login:    "xxx",
					Password: "yyy",
				},
				Metainfo: model.Metainfo{
					Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		}, nil
	}
	if offset > 1 {
		return []model.PairInfo{}, nil
	}
	return []model.PairInfo{
		{
			Pair: model.Pair{
				Login:    "aaa",
				Password: "bbb",
			},
			Metainfo: model.Metainfo{
				Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Text: "pair from yandex",
			},
		},
		{
			Pair: model.Pair{
				Login:    "ccc",
				Password: "ddddd",
			},
			Metainfo: model.Metainfo{
				Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Text: "pair from google",
			},
		},
	}, nil
}

type PairPageWrapper struct {
	*pairPage
}

func (m *PairPageWrapper) Init() tea.Cmd {
	return nil
}

func (m *PairPageWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.pairPage.Update(m, msg)
}

func TestShowPairs(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(func() {}, &MockController{})
	catwalk.RunModel(t, "pairs", &PairPageWrapper{m})
}
