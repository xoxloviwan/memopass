package config

import (
	"flag"
	"fmt"
)

const (
	defaultAddress = "localhost:3000"
	defaultScheme  = "http"
)

var (
	address    = flag.String("a", defaultAddress, "server address")
	Scheme     = flag.Bool("https", false, "use https scheme")
	rootCApath = flag.String("r", "", "path to root ca")
)

type Config struct {
	Address    string
	Version    string
	RootCApath string
}

func InitConfig(version string) *Config {
	scheme := defaultScheme
	if Scheme != nil && *Scheme {
		scheme += "s"
	}
	if address == nil {
		address = new(string)
		*address = defaultAddress
	}
	if rootCApath == nil {
		rootCApath = new(string)
		*rootCApath = ""
	}
	return &Config{
		Address:    fmt.Sprintf("%s://%s", scheme, *address),
		Version:    version,
		RootCApath: *rootCApath,
	}
}
