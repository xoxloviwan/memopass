package main

import (
	"fmt"
	"iwakho/gopherkeep/internal/views/login"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	tabs := []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"}
	tabContent := []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab", "Mascara Tab", "Foundation Tab"}
	m := login.NewModelTabs(tabs, tabContent)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
