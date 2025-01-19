package messages

import tea "github.com/charmbracelet/bubbletea"

type NextPage struct {
	PageID int
	tea.Msg
}

func NextPageCmd(pageID int, msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return NextPage{PageID: pageID, Msg: msg}
	}
}
