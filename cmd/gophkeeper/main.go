package main

import (
	"flag"
	"fmt"
	"iwakho/gopherkeep/internal/cli/views"
	"log"
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
		log.Fatal("Too many arguments")
	}
	if *vers {
		fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
		os.Exit(0)
	}
}

func main() {
	m, err := views.NewApp()
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
