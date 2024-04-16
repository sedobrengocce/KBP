package dialogComponent

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type textInput struct {
    input textinput.Model
}

func NewTextInput(prompt, placeholder string) *textInput {
    t := textinput.New()
    t.Prompt = prompt
    t.Placeholder = placeholder
    return &textInput{
        input: t,
    }
}

func (t *textInput) Focus() tea.Cmd {
    cmds := []tea.Cmd{
        t.input.Focus(),
        t.input.Cursor.SetMode(cursor.CursorBlink),
    }
    
    return tea.Batch(cmds...)
}

func (t *textInput) Blur() tea.Cmd {
    t.input.Blur()
    return t.input.Cursor.SetMode(cursor.CursorHide)
}

func (t textInput) IsFocused() bool {
    return t.input.Focused()
}

func (t textInput) GetText() string {
    return t.input.Value()
}

func (t *textInput) Update (msg tea.Msg) tea.Cmd {
    ti, cmd := t.input.Update(msg)
    t.input = ti
    return cmd
}

func (t *textInput) Render() string {
    return t.input.View()
}
