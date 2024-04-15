package dialogComponent

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type textArea struct {
    area textarea.Model
}

func NewTextArea(prompt, placeholder string) *textArea {
    area := textarea.New()
    area.Prompt = prompt
    area.Placeholder = placeholder
    return &textArea{
        area: area,
    }
}

func (t *textArea) Focus() {
    t.area.Focus()
    t.area.Cursor.SetMode(cursor.CursorBlink)
}

func (t *textArea) Blur() {
    t.area.Blur()
    t.area.Cursor.SetMode(cursor.CursorHide)
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
    return t.area.View()
}
