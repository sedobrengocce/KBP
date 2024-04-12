package dialog

import (
	"errors"
	mesg "kaban-board-plus/common/msg"
	"kaban-board-plus/component/button"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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
    focusedButton int
}

func NewDialog(title, message string, width, height, focusedButton int, buttons []button.Button) *Dialog {
    return &Dialog{
        title: title,
        message: message,
        width: width,
        height: height,
        buttons: buttons,
        focusedButton: focusedButton,
    }
}

func (d *Dialog) NextButton() {
    d.buttons[d.focusedButton].Blur()
    d.focusedButton = (d.focusedButton + 1) % len(d.buttons)
    d.buttons[d.focusedButton].Focus()
}

func (d *Dialog) PrevButton() {
    d.buttons[d.focusedButton].Blur()
    d.focusedButton = (d.focusedButton - 1 + len(d.buttons)) % len(d.buttons)
    d.buttons[d.focusedButton].Focus()
}

func (d Dialog) Click() (any, error) {
    return d.buttons[d.focusedButton].Click()
}

func (d *Dialog) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keys.Enter):
            any, err := d.Click()
            if err != nil {
                return mesg.NewErrorMsg(err)
            }
            cmd, ok := any.(tea.Cmd)
            if !ok{
                return mesg.NewErrorMsg(errors.New("Invalid command"))
            }
            return cmd
        case key.Matches(msg, keys.Left):
            d.PrevButton()
        case key.Matches(msg, keys.Right):
            d.NextButton()
        }
    }
    return nil
}

func (d Dialog) Render() string {
    buttons := []string{}
    for _, b := range d.buttons {
        buttons = append(buttons, b.Render())
        buttons = append(buttons, "    ")
    }
    btnLine := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)
    content := lipgloss.JoinVertical(lipgloss.Center, titleStyle.Render(d.title)," ", messageStyle.Render(d.message), " ", " ", btnLine)
    return lipgloss.Place(d.width, d.height, lipgloss.Center, lipgloss.Center, dialogBoxStyle.Render(content)) 
}
