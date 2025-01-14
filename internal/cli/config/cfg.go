package config

import (
	"flag"
	"os"
)

const defaultAddress = "https://localhost"

var (
	address = flag.String("a", defaultAddress, "server address")
)

type Config struct {
	Address    string
	Version    string
	RootCApath string
}

func InitConfig(version string) *Config {
	rootCApath, _ := os.LookupEnv("ROOTCA_PATH")
	return &Config{
		Address:    *address,
		Version:    version,
		RootCApath: rootCApath,
	}
}
