package dialogComponent

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInput struct {
    input textinput.Model
}

func NewTextInput(prompt, placeholder string) *TextInput {
    t := textinput.New()
    t.Prompt = prompt
    t.Placeholder = placeholder
    return &TextInput{
        input: t,
    }
}

func (t *TextInput) Focus() {
    t.input.Focus()
}

func (t *TextInput) Blur() {
    t.input.Blur()
}

func (t TextInput) IsFocused() bool {
    return t.input.Focused()
}

func (t TextInput) GetText() string {
    return t.input.Value()
}

func (t *TextInput) Update (msg tea.Msg) tea.Cmd {
    ti, cmd := t.input.Update(msg)
    t.input = ti
    return cmd
}

func (t *TextInput) Render() string {
    return t.input.View()
}
