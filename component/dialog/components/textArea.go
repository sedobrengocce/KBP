package dialogComponent

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type textArea struct {
    area textarea.Model
    title string
}

func NewTextArea(prompt, placeholder string) *textArea {
    area := textarea.New()
    area.ShowLineNumbers = false
    area.CharLimit = 200
    area.Placeholder = placeholder
    return &textArea{
        area: area,
        title: prompt,
    }
}

func (t *textArea) Focus() tea.Cmd {
    cmds := []tea.Cmd{
        t.area.Focus(),
        t.area.Cursor.SetMode(cursor.CursorBlink),
    }
    return tea.Batch(cmds...)
}

func (t *textArea) Blur() tea.Cmd {
    t.area.Blur()
    return t.area.Cursor.SetMode(cursor.CursorHide)
}

func (t textArea) IsFocused() bool {
    return t.area.Focused()
}

func (t textArea) GetText() string {
    return t.area.Value()
}

func (t *textArea) Update (msg tea.Msg) tea.Cmd {
    ta, cmd := t.area.Update(msg)
    t.area = ta
    return cmd
}

func (t *textArea) Render() string {
    return lipgloss.JoinVertical(lipgloss.Left,
        t.title,
        t.area.View(),
    )
}
