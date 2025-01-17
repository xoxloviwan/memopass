package model

import (
	"bytes"
	"encoding/json"
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

type Decryptor interface {
	Decrypt(string) (string, error)
}

func DecryptPairs(data []byte, crpt Decryptor) (pairs []PairInfo, err error) {
	err = json.Unmarshal(data, &pairs)
	if err != nil {
		return nil, err
	}

	for i := range pairs {
		login, err := crpt.Decrypt(pairs[i].Login)
		if err != nil {
			return pairs, err
		}
		pairs[i].Login = login
		password, err := crpt.Decrypt(pairs[i].Password)
		if err != nil {
			return pairs, err
		}
		pairs[i].Password = password
	}
	return pairs, nil
}
