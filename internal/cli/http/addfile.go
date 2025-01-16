package http

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (cli *Client) AddFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
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
	r, err := http.NewRequest("POST", cli.api.add.file, body)
	if err != nil {
		return err
	}
	r.Header.Set("Authorization", cli.token)
	r.Header.Set("Content-Type", w.FormDataContentType())
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
