package config

import (
	"flag"
)

const (
	defaultAddress  = ":3000"
	defaultDBPath   = "db/memopass.db"
	defaultCertPath = ""
	defaultKeyPath  = ""
)

var (
	address  = flag.String("a", defaultAddress, "server address")
	dbPath   = flag.String("d", defaultDBPath, "database path")
	certPath = flag.String("c", defaultCertPath, "certificate path")
	keyPath  = flag.String("k", defaultKeyPath, "key path")
)

type Config struct {
	Address  string
	Version  string
	DBPath   string
	CertPath string
	KeyPath  string
}

func InitConfig(version string) *Config {
	if address == nil {
		address = new(string)
		*address = defaultAddress
	}
	if dbPath == nil {
		dbPath = new(string)
		*dbPath = defaultDBPath
	}
	if certPath == nil {
		certPath = new(string)
		*certPath = defaultCertPath
	}
	if keyPath == nil {
		keyPath = new(string)
		*keyPath = defaultKeyPath
	}
	return &Config{
		Address:  *address,
		Version:  version,
		DBPath:   *dbPath,
		CertPath: *certPath,
		KeyPath:  *keyPath,
	}
}
