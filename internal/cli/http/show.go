package http

import (
	"encoding/json"
	"fmt"
	"io"
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"strconv"
)

func (cli *Client) GetPairs(limit int, offset int) ([]model.PairInfo, error) {
	pairs := []model.PairInfo{}
	r, err := http.NewRequest("GET", cli.Api.Get.Pair, nil)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	r.URL.RawQuery = q.Encode()
	r.Header.Set("Authorization", cli.token)
	resp, err := cli.Do(r)
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
	for i := range pairs {
		login, err := CrptMngr.Decrypt(pairs[i].Login)
		if err == nil {
			pairs[i].Login = login
		}
		password, err := CrptMngr.Decrypt(pairs[i].Password)
		if err == nil {
			pairs[i].Password = password
		}
	}

	return pairs, nil
}
