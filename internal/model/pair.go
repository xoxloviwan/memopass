package model

import (
	"bytes"
)

type Pair struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type PairInfo struct {
	Pair
	Metainfo `json:"meta"`
}

func FillPairForm(p Pair, crpt Encryptor) (body *bytes.Buffer, header string, err error) {
	body = new(bytes.Buffer)
	w := newEncryptWriter(crpt, body)
	err = w.encryptField("login", p.Login)
	if err != nil {
		return nil, "", err
	}
	err = w.encryptField("password", p.Password)
	if err != nil {
		return nil, "", err
	}
	err = w.Close()
	if err != nil {
		return nil, "", err
	}
	return body, w.FormDataContentType(), nil
}
