package main

import (
	"errors"
	"flag"
	"fmt"
	app "iwakho/gopherkeep/internal/cli"
	"iwakho/gopherkeep/internal/cli/config"
	httpClient "iwakho/gopherkeep/internal/cli/http"
	"os"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
	vers         = flag.Bool("v", false, "version")
)

func init() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		fatal(errors.New("too many arguments"))
	}
	if vers != nil && *vers {
		fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
		os.Exit(0)
	}
}

func main() {
	cfg := config.InitConfig(buildVersion)
	client, err := httpClient.New(cfg.RootCApath, cfg.Address)
	fatal(err)
	err = app.New(client).Run()
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
