package main

import (
	"flag"
	"fmt"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/cli/views"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const defaultAddress = "https://localhost"

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
	vers         = flag.Bool("v", false, "version")
	address      = flag.String("a", defaultAddress, "server address")
)

func init() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		fmt.Printf("Too many arguments")
		os.Exit(1)
	}
	if vers != nil && *vers {
		fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
		os.Exit(0)
	}
	rootCApath, _ := os.LookupEnv("ROOTCA_PATH")
	if address == nil {
		address = new(string)
		*address = defaultAddress
	}
	err := iHttp.InitClient(rootCApath, *address)
	fatal(err)
}

func main() {
	m, err := views.NewApp()
	fatal(err)
	_, err = tea.NewProgram(m).Run()
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
