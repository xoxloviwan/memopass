package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func (cli *Client) AddItem(url string, body *bytes.Buffer, contentHeader string) error {
	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	r.Header.Set("Authorization", cli.token)
	r.Header.Set("Content-Type", contentHeader)
	resp, err := cli.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Bad status code: %d %s", resp.StatusCode, string(data))
	}
	return nil
}
