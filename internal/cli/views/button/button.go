package button

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func RenderButton(sb *strings.Builder, btn string, focused bool) {
	if focused {
		btn = FocusedStyle.Render(fmt.Sprintf("[ %s ]", btn))
	} else {
		btn = fmt.Sprintf("[ %s ]", BlurredStyle.Render(btn))
	}
	fmt.Fprintf(sb, "%s\n", btn)
}
