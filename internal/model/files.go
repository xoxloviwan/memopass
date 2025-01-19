package model

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type File struct {
	Name string `json:"name"`
	Blob []byte `json:"blob"`
	ID   int    `json:"id"`
}

type FileInfo struct {
	File
	Metainfo `json:"meta"`
}

func FillTextForm(text string) (body *bytes.Buffer, header string, err error) {
	body = new(bytes.Buffer)
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile("text", "text.txt")
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(part, bytes.NewBuffer([]byte(text)))
	if err != nil {
		return nil, "", err
	}
	err = w.Close()
	if err != nil {
		return nil, "", err
	}
	return body, w.FormDataContentType(), nil
}

func FillFileForm(file *os.File) (body *bytes.Buffer, header string, err error) {
	body = new(bytes.Buffer)
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, "", err
	}
	err = w.Close()
	if err != nil {
		return nil, "", err
	}
	return body, w.FormDataContentType(), nil
}
