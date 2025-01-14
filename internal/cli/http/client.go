package http

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
)

var (
	Client     http.Client
	ApiLogin   string
	ApiSignUp  string
	ApiAddItem string
	ApiGetItem string
)

func InitClient(certPath string, baseURL string) error {
	apiURL := baseURL + "/api/v1"
	ApiLogin = apiURL + "/user/login"
	ApiSignUp = apiURL + "/user/signup"
	ApiAddItem = apiURL + "/item/add"
	ApiGetItem = apiURL + "/item"
	if certPath == "" {
		Client = http.Client{}
		return nil
	}
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	// Плохо понимаю нужно ли вручную грузить сертификат,
	// т.к. Wireshark показывает TLS и данные выглядят зашифрованными.
	// Ошибка "failed to verify certificate" возникает, если создать пустой CertPool и не грузить туда ничего.
	// Если создать дефолтный http.Client{} без CertPool, то по дефолту какой-то сертификат видимо есть.

	Client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	return nil
}
