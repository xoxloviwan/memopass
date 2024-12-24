package login

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

// Describes the styles for the auth page
var (
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	highlightColor    = lipgloss.AdaptiveColor{Light: "#13002D", Dark: "#22ADF6"}
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	inactiveTabStyle = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 3)
	activeTabStyle   = inactiveTabStyle.Border(activeTabBorder, true).Bold(true)
	magicOffset      = 2
	windowStyle      = lipgloss.NewStyle().
				BorderForeground(highlightColor).
				Padding(1, magicOffset).
				Align(lipgloss.Left).
				Border(lipgloss.NormalBorder()).UnsetBorderTop()
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

func NewAuthPage(onEnter Callback) AuthPage {
	ap := AuthPage{
		TabContent: []modelForm{
			InitLogin(onEnter),
			InitSignUp(onEnter),
		},
		Tabs: make(map[int]Tab),
	}
	for i := range ap.TabContent {
		ap.Tabs[i] = Tab{Name: ap.TabContent[i].name}
	}
	return ap
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
func (ap *AuthPage) renderTabs() string {
	var renderedTabs []string

	// Sort the keys to make sure tabs are in the same order everytime
	keys := make([]int, 0, len(ap.Tabs))
	for k := range ap.Tabs {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for i, k := range keys {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(keys)-1, k == ap.activatedTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(ap.Tabs[k].Name))
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderedTabs...,
	)
	return row
}

func (ap AuthPage) View() string {
	doc := strings.Builder{}

	// Tabs
	row := ap.renderTabs()
	_, err := doc.WriteString(row)
	if err != nil {
		return err.Error()
	}
	doc.WriteString("\n")

	windowWidth := lipgloss.Width(row) - windowStyle.GetHorizontalPadding() + magicOffset
	// Content
	_, err = doc.WriteString(windowStyle.Width(windowWidth).Render(ap.TabContent[ap.activatedTab].View()))
	if err != nil {
		return err.Error()
	}

	return docStyle.Render(doc.String())
}
