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
	pair string
	card string
	file string
	text string
}

type api struct {
	login   string
	signUp  string
	add     items
	getItem string
}

type Client struct {
	token   string
	baseURL string
	api
	http.Client
}

func InitClient(certPath string, baseURL string) (*Client, error) {
	apiURL := baseURL + apiBase
	api := api{
		login:  apiURL + apiLogin,
		signUp: apiURL + apiSignUp,
		add: items{
			pair: apiURL + apiAddPair,
			card: apiURL + apiAddCard,
			file: apiURL + apiAddFile,
			text: apiURL + apiAddText,
		},
		getItem: apiURL + apiGetItem,
	}
	if certPath == "" {
		return &Client{
			Client:  http.Client{},
			baseURL: baseURL,
			api:     api,
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
		api:     api,
	}, nil
}
