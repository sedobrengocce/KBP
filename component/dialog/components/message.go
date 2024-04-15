package dialogComponent

import tea "github.com/charmbracelet/bubbletea"

type Message struct {
    message string
}

func NewMessage(message string) *Message {
    return &Message{message: message}
}

func (m *Message) Update(_ tea.Msg) tea.Cmd {
    return nil
}

func (m Message) Render() string {
    return m.message
}
