package show

import (
	"iwakho/gopherkeep/internal/model"
	"testing"
	"time"

	"github.com/knz/catwalk"
)

//go:generate go test . -args -rewrite

type MockController struct {
	cnt int
}

func (c *MockController) GetPairs(limit int, offset int) ([]model.PairInfo, error) {
	c.cnt++
	if offset == 2 {
		return []model.PairInfo{
			{
				Pair: model.Pair{
					Login:    "yyy",
					Password: "zzz",
				},
				Metainfo: model.Metainfo{
					Date: time.Date(2020, 1, 1, 1, 0, 0, 0, time.FixedZone("MSK", 3*3600)),
				},
			},
			{
				Pair: model.Pair{
					Login:    "xxx",
					Password: "yyy",
				},
				Metainfo: model.Metainfo{
					Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.FixedZone("MSK", 3*3600)),
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
				Date: time.Date(2020, 1, 1, 3, 0, 0, 0, time.FixedZone("MSK", 3*3600)),
				Text: "pair from yandex",
			},
		},
		{
			Pair: model.Pair{
				Login:    "ccc",
				Password: "ddddd",
			},
			Metainfo: model.Metainfo{
				Date: time.Date(2020, 1, 1, 2, 0, 0, 0, time.FixedZone("MSK", 3*3600)),
				Text: "pair from google",
			},
		},
	}, nil
}

func TestShowPairs(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(0, &MockController{})
	catwalk.RunModel(t, "pairs", m)
}
