package http

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
)

const (
	apiLogin   = "/user/login"
	apiSignUp  = "/user/signup"
	apiAddPair = "/item/add/pair"
	apiAddCard = "/item/add/card"
	apiAddFile = "/item/add/file"
	apiAddText = "/item/add/text"
	apiGetItem = "/item"
	apiBase    = "/api/v1"
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
	getItem string
}

type Client struct {
	token   string
	baseURL string
	Api
	http.Client
}

func InitClient(certPath string, baseURL string) (*Client, error) {
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
		getItem: apiURL + apiGetItem,
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
