package http

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
)

const (
	apiLogin       = "/user/login"
	apiSignUp      = "/user/signup"
	apiAddPair     = "/item/add/pair"
	apiAddCard     = "/item/add/card"
	apiAddFile     = "/item/add/file"
	apiAddText     = "/item/add/text"
	apiGetPairs    = "/item/pairs"
	apiGetCards    = "/item/cards"
	apiGetFiles    = "/item/files"
	apiGetTexts    = "/item/texts"
	apiGetTextByID = "/item/text"
	apiGetFileByID = "/item/file"
	apiBase        = "/api/v1"
)

type items struct {
	Pair string
	Card string
	File string
	Text string
}

type Api struct {
	login   string
	signUp  string
	Add     items
	Get     items
	GetById items
}

type Client struct {
	token   string
	baseURL string
	Api
	http.Client
}

func New(certPath string, baseURL string) (*Client, error) {
	apiURL := baseURL + apiBase
	api := Api{
		login:  apiURL + apiLogin,
		signUp: apiURL + apiSignUp,
		Add: items{
			Pair: apiURL + apiAddPair,
			Card: apiURL + apiAddCard,
			File: apiURL + apiAddFile,
			Text: apiURL + apiAddText,
		},
		Get: items{
			Pair: apiURL + apiGetPairs,
			Card: apiURL + apiGetCards,
			File: apiURL + apiGetFiles,
			Text: apiURL + apiGetTexts,
		},
		GetById: items{
			File: apiURL + apiGetFileByID,
			Text: apiURL + apiGetTextByID,
		},
	}
	if certPath == "" {
		return &Client{
			Client:  http.Client{},
			baseURL: baseURL,
			Api:     api,
		}, nil
	}
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return &Client{
		Client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		},
		baseURL: baseURL,
		Api:     api,
	}, nil
}
