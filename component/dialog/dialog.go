package dialog

import (
    "kaban-board-plus/component/button"

    "github.com/charmbracelet/lipgloss"
)

var (
    dialogBoxStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#874BFD")).
    Padding(1, 0).
    BorderTop(true).
    BorderLeft(true).
    BorderRight(true).
    BorderBottom(true)

    titleStyle = lipgloss.NewStyle().
    Width(50).
    Align(lipgloss.Center).
    Bold(true).
    Foreground(lipgloss.Color("#FF9800"))

    messageStyle = lipgloss.NewStyle().
    Width(50).
    Align(lipgloss.Center)
)

type Dialog struct {
    title string
    message string
    width int
    height int
    buttons []button.Button
}

func NewDialog(title, message string, width, height int, buttons []button.Button) *Dialog {
    return &Dialog{
        title: title,
        message: message,
        width: width,
        height: height,
        buttons: buttons,
    }
}

func (d Dialog) Render() string {
    buttons := []string{}
    for _, b := range d.buttons {
        buttons = append(buttons, b.Render())
    }
    btnLine := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)
    content := lipgloss.JoinVertical(lipgloss.Top, titleStyle.Render(d.title), messageStyle.Render(d.message), btnLine)
    return lipgloss.Place(d.width, d.height, lipgloss.Center, lipgloss.Center, dialogBoxStyle.Render(content)) 
}
