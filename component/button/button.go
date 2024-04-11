package button

import "github.com/charmbracelet/lipgloss"

var (
    activeButtonStyle = lipgloss.NewStyle().Padding(1,2).Background(lipgloss.Color("#FF9800")).Foreground(lipgloss.Color("#000000"))
    inactiveButtonStyle = lipgloss.NewStyle().Padding(1,2).Background(lipgloss.Color("#E3E3E3")).Foreground(lipgloss.Color("#000000"))
)

type Button struct {
    text string
    action func() (any, error)
    isActive bool
}

func NewButton(text string, action func() (any, error)) *Button {
    return &Button{
        text: text,
        action: action,
        isActive: false,
    }
}

func (b *Button) SetActive() {
    b.isActive = true
}

func (b *Button) SetInactive() {
    b.isActive = false
}

func (b *Button) Click() (any, error) {
    return b.action()
}

func (b Button) Render() string {
    if b.isActive {
        return activeButtonStyle.Render(b.text)
    }
    return inactiveButtonStyle.Render(b.text)
}
