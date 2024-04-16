package dialog

import (
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
)

type DialogComponent interface {
    Focus() tea.Cmd
    Blur() tea.Cmd
    Update(msg tea.Msg) tea.Cmd 
    Render() string
    IsFocused() bool
}

type Dialog struct {
    title string
    content []DialogComponent
    help help.Model
    width int
    height int
    buttons []button.Button
    focusedButton int
    focusedComponent int
}

func NewDialog(title string, components []DialogComponent, width, height int, buttons []button.Button) *Dialog {
    return &Dialog{
        title: title,
        content: components,
        width: width,
        height: height,
        focusedComponent: 0,
        buttons: buttons,
    }
}

func (d *Dialog) nextButton() {
    if d.focusedComponent < len(d.content) {
        return
    }
    d.buttons[d.focusedButton].Blur()
    d.focusedButton = (d.focusedButton + 1) % len(d.buttons)
    d.buttons[d.focusedButton].Focus()
}

func (d *Dialog) prevButton() {
    if d.focusedComponent < len(d.content) {
        return
    }
    d.buttons[d.focusedButton].Blur()
    d.focusedButton = (d.focusedButton - 1 + len(d.buttons)) % len(d.buttons)
    d.buttons[d.focusedButton].Focus()
}

func (d Dialog) click() tea.Cmd {
    return d.buttons[d.focusedButton].Click()
}

func (d *Dialog) nextComponent() tea.Cmd {
    var cmd tea.Cmd
    if d.focusedComponent != len(d.content) {
        cmd = d.content[d.focusedComponent].Blur()
    }
    d.focusedComponent = (d.focusedComponent + 1) % (len(d.content) + 1)
    if d.focusedComponent == len(d.content) {
        d.buttons[d.focusedButton].Focus()
        return cmd
    } else {
        for i := range d.buttons {
            d.buttons[i].Blur()
        }
        cmds := []tea.Cmd{
            cmd,
            d.content[d.focusedComponent].Focus(),
        }
        return tea.Batch(cmds...)
    }
}  

func (d *Dialog) Update(msg tea.Msg) tea.Cmd {
    for i := range d.content {
        if d.content[i].IsFocused() {
            cmd :=d.content[i].Update(msg)
            if cmd != nil {
                return cmd
            }
        }
    }
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keys.Enter):
            cmd := d.click()
            return cmd
        case key.Matches(msg, keys.Left):
            d.prevButton()
            return nil
        case key.Matches(msg, keys.Right):
            d.nextButton()
            return nil
        case key.Matches(msg, keys.Tab):
           return d.nextComponent()
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
    content := []string{}
    content = append(content, titleStyle.Render(d.title))
    content = append(content, " ")
    for _, c := range d.content {
        content = append(content, c.Render())
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
