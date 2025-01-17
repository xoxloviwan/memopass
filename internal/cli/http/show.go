package http

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (cli *Client) GetItems(url string, limit int, offset int) (data []byte, err error) {
	r, err := http.NewRequest("GET", url, nil)
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
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
