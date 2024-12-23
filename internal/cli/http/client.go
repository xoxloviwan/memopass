package http

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
)

var (
	Client  http.Client
	BaseURL string
)

func InitClient(certPath string, baseURL string) error {
	BaseURL = baseURL + "/api/v1"
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
