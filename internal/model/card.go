package model

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
)

type Card struct {
	Number   string `json:"ccn"`
	Exp      string `json:"exp"`
	VerifVal string `json:"cvv"`
}

type CardInfo struct {
	Card
	Metainfo `json:"meta"`
}

type Encryptor interface {
	Encrypt(string) (string, error)
}

type encryptWriter struct {
	*multipart.Writer
	Encryptor
}

func newEncryptWriter(crpt Encryptor, body *bytes.Buffer) *encryptWriter {
	return &encryptWriter{Encryptor: crpt, Writer: multipart.NewWriter(body)}
}

func (w *encryptWriter) encryptField(name string, value string) error {
	ecnrypted, err := w.Encrypt(value)
	if err != nil {
		return err
	}
	return w.WriteField(name, ecnrypted)
}

func FillCardForm(card Card, crpt Encryptor) (body *bytes.Buffer, header string, err error) {
	body = new(bytes.Buffer)
	w := newEncryptWriter(crpt, body)
	err = w.encryptField("ccn", card.Number)
	if err != nil {
		return nil, "", err
	}
	err = w.encryptField("exp", card.Exp)
	if err != nil {
		return nil, "", err
	}
	err = w.encryptField("cvv", card.VerifVal)
	if err != nil {
		return nil, "", err
	}
	err = w.Close()
	if err != nil {
		return nil, "", err
	}
	return body, w.FormDataContentType(), nil
}

func DecryptCards(data []byte, crpt Decryptor) (cards []CardInfo, err error) {
	err = json.Unmarshal(data, &cards)
	if err != nil {
		return nil, err
	}
	for i := range cards {
		num, err := crpt.Decrypt(cards[i].Number)
		if err != nil {
			return cards, err
		}
		cards[i].Number = num
		exp, err := crpt.Decrypt(cards[i].Exp)
		if err != nil {
			return cards, err
		}
		cards[i].Exp = exp
		cvv, err := crpt.Decrypt(cards[i].VerifVal)
		if err != nil {
			return cards, err
		}
		cards[i].VerifVal = cvv
	}
	return cards, nil
}
