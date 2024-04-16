package dialogComponent

import tea "github.com/charmbracelet/bubbletea"

type componentPreventMsg struct {}

func newComponentPreventMsg() tea.Cmd {
    return func() tea.Msg {
        return componentPreventMsg{}
    }
}
