package http

import (
	"bytes"
	"fmt"
	"io"
	"iwakho/gopherkeep/internal/model"
	"mime/multipart"
	"net/http"
)

func (cli *Client) AddPair(p model.Pair) error {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	encP := model.Pair{}
	var err error
	encP.Login, err = CrptMngr.Encrypt(p.Login)
	if err != nil {
		return err
	}
	err = w.WriteField("login", encP.Login)
	if err != nil {
		return err
	}
	encP.Password, err = CrptMngr.Encrypt(p.Password)
	if err != nil {
		return err
	}
	err = w.WriteField("password", encP.Password)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", cli.Api.AddItem, body)
	if err != nil {
		return err
	}
	r.Header.Set("Authorization", cli.token)
	r.Header.Set("Content-Type", w.FormDataContentType())
	q := r.URL.Query()
	q.Add("type", "0")
	r.URL.RawQuery = q.Encode()
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
