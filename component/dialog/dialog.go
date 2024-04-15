package dialog

import (
	"errors"
	mesg "kaban-board-plus/common/msg"
	"kaban-board-plus/component/button"

	"github.com/charmbracelet/bubbles/help"
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

type DialogComponent interface {
    Update(msg tea.Msg) tea.Cmd 
    Render() string
}

type Dialog struct {
    title string
    content []DialogComponent
    help help.Model
    width int
    height int
    buttons []button.Button
    focusedButton int
}

func NewDialog(title string, components []DialogComponent, width, height, focusedButton int, buttons []button.Button) *Dialog {
    return &Dialog{
        title: title,
        content: components,
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
    var cmds []tea.Cmd
    for _, c := range d.content {
        cmd := c.Update(msg)
        if cmd != nil {
            cmds = append(cmds, cmd)
        }
    }
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keys.Enter):
            any, err := d.Click()
            if err != nil {
                cmd := mesg.NewErrorMsg(err)
                cmds = append(cmds, cmd)
                return tea.Batch(cmds...)
            }
            cmd, ok := any.(tea.Cmd)
            if !ok{
                cmd = mesg.NewErrorMsg(errors.New("Invalid command"))
                cmds = append(cmds, cmd)
                return tea.Batch(cmds...)
            }
            cmds = append(cmds, cmd)
            return tea.Batch(cmds...)
        case key.Matches(msg, keys.Left):
            d.PrevButton()
        case key.Matches(msg, keys.Right):
            d.NextButton()
        }
    }
    return tea.Batch(cmds...)
}

func (d Dialog) Render() string {
    buttons := []string{}
    for _, b := range d.buttons {
        buttons = append(buttons, b.Render())
        buttons = append(buttons, "    ")
    }
    content := []string{}
    content = append(content, titleStyle.Render(d.title))
    content = append(content, " ")
    for _, c := range d.content {
        content = append(content, messageStyle.Render(c.Render()))
    }
    content = append(content, " ")
    content = append(content, " ")
        btnLine := lipgloss.JoinHorizontal(
        lipgloss.Top, 
        buttons...
    )
    content = append(content, btnLine)
    renderedContent := lipgloss.JoinVertical(
        lipgloss.Center, 
        content..., 
    )
    return lipgloss.Place(d.width, d.height, 
        lipgloss.Center, lipgloss.Center, 
        dialogBoxStyle.Render(renderedContent),
    ) 
}
