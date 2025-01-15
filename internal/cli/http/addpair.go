package http

import (
	"bytes"
	"fmt"
	"io"
	"iwakho/gopherkeep/internal/model"
	"mime/multipart"
	"net/http"
	"strconv"
)

func fillPairForm(p model.Pair, body *bytes.Buffer) (*multipart.Writer, error) {
	w := newEncryptWriter(body)
	err := w.encryptField("login", p.Login)
	if err != nil {
		return nil, err
	}
	err = w.encryptField("password", p.Password)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return w.Writer, nil
}

func (cli *Client) AddPair(p model.Pair) error {
	body := new(bytes.Buffer)
	w, err := fillPairForm(p, body)
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
	q.Add("type", strconv.Itoa(model.ItemTypeLoginPass))
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
