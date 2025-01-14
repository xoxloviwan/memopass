package controls

import (
	"bytes"
	"fmt"
	"io"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/model"
	"mime/multipart"
	"net/http"
)

type AddPairCtrl struct{}

func (AddPairCtrl) Submit(p model.Pair) error {
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
	r, err := http.NewRequest("POST", iHttp.ApiAddItem, body)
	if err != nil {
		return err
	}
	r.Header.Set("Authorization", token)
	r.Header.Set("Content-Type", w.FormDataContentType())
	q := r.URL.Query()
	q.Add("type", "0")
	r.URL.RawQuery = q.Encode()
	resp, err := iHttp.Client.Do(r)
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
