package show

import (
	"fmt"
	"iwakho/gopherkeep/internal/model"
	"reflect"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/knz/catwalk"
)

//go:generate go test . -args -rewrite

type MockController struct {
	cnt int
}

func (c *MockController) GetPairs(limit int, offset int) ([]model.PairInfo, error) {
	c.cnt++
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}
	fmt.Println("LoadLocation", loc, loc == nil, reflect.TypeOf(loc))
	fmt.Println("LoadLocation time1", time.Date(2020, 1, 1, 3, 0, 0, 0, loc))
	fmt.Println("LoadLocation time2", time.Date(2020, 1, 1, 3, 0, 0, 0, loc).In(loc))
	fmt.Println("LoadLocation time3", time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC).In(loc))
	if offset == 2 {
		return []model.PairInfo{
			{
				Pair: model.Pair{
					Login:    "yyy",
					Password: "zzz",
				},
				Metainfo: model.Metainfo{
					Date: time.Date(2020, 1, 1, 1, 0, 0, 0, loc),
				},
			},
			{
				Pair: model.Pair{
					Login:    "xxx",
					Password: "yyy",
				},
				Metainfo: model.Metainfo{
					Date: time.Date(2020, 1, 1, 0, 0, 0, 0, loc),
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
				Date: time.Date(2020, 1, 1, 3, 0, 0, 0, loc),
			},
		},
		{
			Pair: model.Pair{
				Login:    "ccc",
				Password: "ddddd",
			},
			Metainfo: model.Metainfo{
				Date: time.Date(2020, 1, 1, 2, 0, 0, 0, loc),
			},
		},
	}, nil
}

func TestShowPairs(t *testing.T) {
	// Initialize the model to test.
	m := NewPage(0, &MockController{})
	catwalk.RunModel(t, "pairs", m)
}
