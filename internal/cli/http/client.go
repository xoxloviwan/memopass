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
	apiAddItem = "/item/add"
	apiGetItem = "/item"
	apiBase    = "/api/v1"
)

type Api struct {
	Login   string
	SignUp  string
	AddItem string
	GetItem string
}

type Client struct {
	token   string
	baseURL string
	Api     Api
	http.Client
}

func InitClient(certPath string, baseURL string) (*Client, error) {
	apiURL := baseURL + apiBase
	api := Api{
		Login:   apiURL + apiLogin,
		SignUp:  apiURL + apiSignUp,
		AddItem: apiURL + apiAddItem,
		GetItem: apiURL + apiGetItem,
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
	// Плохо понимаю нужно ли вручную грузить сертификат,
	// т.к. Wireshark показывает TLS и данные выглядят зашифрованными.
	// Ошибка "failed to verify certificate" возникает, если создать пустой CertPool и не грузить туда ничего.
	// Если создать дефолтный http.Client{} без CertPool, то по дефолту какой-то сертификат видимо есть.
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
