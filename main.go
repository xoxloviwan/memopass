package main

import (
	"fmt"
	"iwakho/gopherkeep/internal/views/login"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m, err := login.NewSampleConfigUI()
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
