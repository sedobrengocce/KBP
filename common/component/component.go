package component

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Screen, tea.Cmd)
    View() string
    ToggleHelp()
    SetSize(width, height int)
}
