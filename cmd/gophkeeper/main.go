package main

import (
	"errors"
	"flag"
	"fmt"
	"iwakho/gopherkeep/internal/cli/config"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/cli/views"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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
	cli, err := iHttp.InitClient(cfg.RootCApath, cfg.Address)
	fatal(err)
	m, err := views.NewApp(cli)
	fatal(err)
	p := tea.NewProgram(m, func(pp *tea.Program) {
		m.Sender = pp
	})
	_, err = p.Run()
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
