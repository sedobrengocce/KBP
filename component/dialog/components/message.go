package dialogComponent

import tea "github.com/charmbracelet/bubbletea"

type Message struct {
    message string
}

func NewMessage(message string) *Message {
    return &Message{message: message}
}

func (m *Message) Focus() tea.Cmd {
    return nil
}

func (m *Message) Blur() tea.Cmd {
    return nil
}

func (m *Message) Update(_ tea.Msg) tea.Cmd {
    return nil
}

func (m Message) Render() string {
    return m.message
}

func (m Message) IsFocused() bool {
    return false
}
