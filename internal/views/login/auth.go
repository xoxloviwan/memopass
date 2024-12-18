package login

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Describes the styles for the plugins page
// Could be re-used by other pages in the future to keep a consisten look
var (
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}
	highlight = lipgloss.AdaptiveColor{Light: "#13002D", Dark: "#22ADF6"}
	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1)
	activeTab = tab.Border(activeTabBorder, true)

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	docStyle = lipgloss.NewStyle().Padding(1, 2, 0, 2)

	special = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	checked = lipgloss.NewStyle().SetString("✓").Foreground(special).PaddingRight(1).String()
)

type Tab struct {
	Name string
}

type AuthPage struct {
	activatedTab int
	Tabs         map[int]Tab
	TabContent   []modelForm

	width  int
	height int
}

func NewAuthPage() (AuthPage, error) {
	InitLogin()
	InitSignUp()
	ap := AuthPage{
		TabContent: []modelForm{
			InitLogin(),
			InitSignUp(),
		},
		Tabs: make(map[int]Tab),
	}
	for i := range ap.TabContent {
		ap.Tabs[i] = Tab{Name: ap.TabContent[i].name}
	}
	return ap, nil
}

func (ap *AuthPage) Init(width, height int) {
	ap.width = width
}

func (ap *AuthPage) Update(m tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println("AuthPage Exiting...")
			return m, tea.Quit
		case "left":
			if ap.activatedTab < len(ap.Tabs)-1 {
				ap.activatedTab++
			} else {
				ap.activatedTab--
			}
			return m, nil
		case "right":
			if ap.activatedTab > 0 {
				ap.activatedTab--
			} else {
				ap.activatedTab++
			}
			return m, nil
		}
	}
	tabModel, cmd := ap.TabContent[ap.activatedTab].Update(msg)
	ap.TabContent[ap.activatedTab] = tabModel.(modelForm)
	return m, cmd
}

// renderTabs will create the view for the tabs
// counting the new lines can help determine the height for other components
func (ap *AuthPage) renderTabs(width int) string {
	var renderedTabs []string

	// Sort the keys to make sure tabs are in the same order everytime
	// Using a map helps with organizing selected plugins with plugin type
	keys := make([]int, 0, len(ap.Tabs))
	for k := range ap.Tabs {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		if k == ap.activatedTab {
			renderedTabs = append(renderedTabs, activeTab.Render(ap.Tabs[k].Name))
		} else {
			renderedTabs = append(renderedTabs, tab.Render(ap.Tabs[k].Name))
		}
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderedTabs...,
	)
	gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap) + "\n\n"
}

func (ap AuthPage) View() string {
	doc := strings.Builder{}

	// Tabs
	row := ap.renderTabs(ap.width)
	_, err := doc.WriteString(row)
	if err != nil {
		return err.Error()
	}

	// Content
	_, err = doc.WriteString(ap.TabContent[ap.activatedTab].View())
	if err != nil {
		return err.Error()
	}

	return docStyle.Render(doc.String())
}
