package controls

import (
	"bytes"
	"fmt"
	"io"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type AddFileCtrl struct{}

func (AddFileCtrl) Submit(filePath string) error {
	file, _ := os.Open(filePath)
	defer file.Close()

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
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
	q.Add("type", "2")
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
