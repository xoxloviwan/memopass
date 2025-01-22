package http

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func (cli *Client) GetFileById(url string, id int) (data []byte, name string, err error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}
	r.Header.Set("Authorization", cli.token)
	q := r.URL.Query()
	q.Add("id", strconv.Itoa(id))
	r.URL.RawQuery = q.Encode()
	resp, err := cli.Do(r)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	cd := resp.Header.Get("Content-Disposition")
	const prefix = "filename="
	if !strings.Contains(cd, prefix) {
		return nil, "", fmt.Errorf("Content-Disposition header not found")
	}
	filename := cd[strings.Index(cd, prefix)+len(prefix)+1 : len(cd)-1]
	return data, filename, nil
}
