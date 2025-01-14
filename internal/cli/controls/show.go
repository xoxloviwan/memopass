package controls

import (
	"encoding/json"
	"fmt"
	"io"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"strconv"
)

type ShowCtrl struct{}

func (s ShowCtrl) GetPairs(limit int, offset int) ([]model.PairInfo, error) {
	pairs := []model.PairInfo{}
	r, err := http.NewRequest("GET", iHttp.ApiGetItem, nil)
	q := r.URL.Query()
	q.Add("type", "0")
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	r.URL.RawQuery = q.Encode()
	r.Header.Set("Authorization", token)
	resp, err := iHttp.Client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &pairs)
	if err != nil {
		return nil, err
	}
	return pairs, nil
}
