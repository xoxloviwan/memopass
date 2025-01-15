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

type encryptWriter struct {
	*multipart.Writer
}

func (w *encryptWriter) encryptField(name string, value string) error {
	ecnrypted, err := CrptMngr.Encrypt(value)
	if err != nil {
		return err
	}
	return w.WriteField(name, ecnrypted)
}

func newEncryptWriter(body *bytes.Buffer) *encryptWriter {
	return &encryptWriter{Writer: multipart.NewWriter(body)}
}

func fillCardForm(card model.Card, body *bytes.Buffer) (*multipart.Writer, error) {
	w := newEncryptWriter(body)
	err := w.encryptField("ccn", card.Number)
	if err != nil {
		return nil, err
	}
	err = w.encryptField("exp", card.Exp)
	if err != nil {
		return nil, err
	}
	err = w.encryptField("cvv", card.VerifVal)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return w.Writer, nil
}

func (cli *Client) AddCard(card model.Card) error {
	body := new(bytes.Buffer)
	w, err := fillCardForm(card, body)
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
	q.Add("type", strconv.Itoa(model.ItemTypeCard))
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
